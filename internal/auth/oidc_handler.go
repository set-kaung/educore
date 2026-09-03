package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"educore/internal"

	"gorm.io/gorm"
)

const oauthCookieTTL = 10 * time.Minute

const stateCookie = "edu_oauth_state"
const nonceCookie = "edu_oauth_nonce"

type OIDCHandler struct {
	Provider  *OIDCProvider
	DB        *gorm.DB
	JWTSecret string
}

type adClaims struct {
	OID               string `json:"oid"`
	Email             string `json:"email"`
	DisplayName       string `json:"displayName"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
}

func (h *OIDCHandler) HandleADLogin(w http.ResponseWriter, r *http.Request) {
	state, err := randomHex()
	if err != nil {
		http.Error(w, "could not start sign-in", http.StatusInternalServerError)
		return
	}
	nonce, err := randomHex()
	if err != nil {
		http.Error(w, "could not start sign-in", http.StatusInternalServerError)
		return
	}

	setOAuthCookie(w, r, stateCookie, state)
	setOAuthCookie(w, r, nonceCookie, nonce)

	http.Redirect(w, r, h.Provider.AuthCodeURL(state, nonce), http.StatusSeeOther)
}

func (h *OIDCHandler) HandleADCallback(w http.ResponseWriter, r *http.Request) {
	clearOAuthCookie(w, r, stateCookie)
	clearOAuthCookie(w, r, nonceCookie)

	if msg := r.URL.Query().Get("error"); msg != "" {
		slog.Warn("oidc authorization failed", "error", msg, "description", r.URL.Query().Get("error_description"))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	stateValue, err := r.Cookie(stateCookie)
	if err != nil || code == "" || state == "" || stateValue.Value != state {
		http.Error(w, "invalid sign-in state", http.StatusBadRequest)
		return
	}

	nonceValue, err := r.Cookie(nonceCookie)
	if err != nil {
		http.Error(w, "invalid sign-in state", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	oauth2Token, err := h.Provider.Exchange(ctx, code)
	if err != nil {
		slog.Error("oidc code exchange failed", "error", err)
		http.Error(w, "could not complete sign-in", http.StatusBadGateway)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		slog.Error("oidc token response missing id_token")
		http.Error(w, "could not complete sign-in", http.StatusBadGateway)
		return
	}

	idToken, err := h.Provider.VerifyIDToken(ctx, rawIDToken)
	if err != nil {
		slog.Error("oidc id token verification failed", "error", err)
		http.Error(w, "could not complete sign-in", http.StatusUnauthorized)
		return
	}

	if idToken.Nonce != nonceValue.Value {
		http.Error(w, "invalid sign-in nonce", http.StatusBadRequest)
		return
	}

	var claims adClaims
	if err := idToken.Claims(&claims); err != nil {
		slog.Error("oidc id token claims unreadable", "error", err)
		http.Error(w, "could not complete sign-in", http.StatusBadGateway)
		return
	}
	if claims.OID == "" {
		slog.Error("oidc id token missing oid claim")
		http.Error(w, "could not complete sign-in", http.StatusBadGateway)
		return
	}
	if claims.Email == "" {
		claims.Email = claims.PreferredUsername
	}
	if claims.DisplayName == "" {
		claims.DisplayName = claims.Name
	}

	userID, role, found, err := FindUserByADObjectID(h.DB, claims.OID)
	if err != nil {
		slog.Error("ad account lookup failed", "error", err)
		http.Error(w, "could not complete sign-in", http.StatusInternalServerError)
		return
	}

	if found {
		h.finishSignIn(w, r, userID, role)
		return
	}

	setupToken, err := IssueSetupToken(h.JWTSecret, claims.OID, claims.Email, claims.DisplayName)
	if err != nil {
		slog.Error("could not issue setup token", "error", err)
		http.Error(w, "could not complete sign-in", http.StatusInternalServerError)
		return
	}
	setSetupTokenCookie(w, r, setupToken)

	http.Redirect(w, r, "/setup", http.StatusSeeOther)
}

// HandleSetupContext godoc
// @Summary      Account setup context
// @Description  Prefill data (email, suggested name) for first-time account setup. Requires a setup token.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  internal.ResponseBody{data=SetupContext}
// @Failure      401  {object}  internal.ResponseBody
// @Router       /api/setup/context [get]
func (h *OIDCHandler) HandleSetupContext(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	claims := GetClaims(r)
	internal.WriteData(w, "", SetupContext{Email: claims.Email, Name: claims.Name}, nil)
	return nil
}

// HandleSetup godoc
// @Summary      Complete account setup
// @Description  Create the Student account for the authenticated directory user and issue a session. Requires a setup token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  SetupRequest  true  "Account details"
// @Success      200  {object}  internal.ResponseBody{data=SetupResult}
// @Failure      400  {object}  internal.ResponseBody
// @Failure      404  {object}  internal.ResponseBody
// @Failure      409  {object}  internal.ResponseBody
// @Router       /api/setup [post]
func (h *OIDCHandler) HandleSetup(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	claims := GetClaims(r)

	var req SetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid request body", Err: err}
	}

	req.Name = strings.TrimSpace(req.Name)
	req.StudentID = strings.TrimSpace(req.StudentID)

	if req.Name == "" {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "name is required", Err: nil}
	}
	if req.StudentID == "" {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "student id is required", Err: nil}
	}

	var department internal.Department
	if err := h.DB.First(&department, req.DepartmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &internal.HTTPError{StatusCode: http.StatusNotFound, Message: "department not found", Err: err}
		}
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not validate department", Err: err}
	}

	if userID, role, found, err := FindUserByADObjectID(h.DB, claims.OID); err == nil && found {
		h.issueSession(w, r, userID, role)
		internal.WriteData(w, "", SetupResult{UserID: userID, Role: role}, nil)
		return nil
	}

	username := claims.Email
	if username == "" {
		username = claims.OID
	}
	user, err := CreateUserForOIDC(h.DB, req.Name, username, req.StudentID, department.ID, claims.OID)
	if err != nil {
		if existingID, existingRole, found, lookupErr := FindUserByADObjectID(h.DB, claims.OID); lookupErr == nil && found {
			h.issueSession(w, r, existingID, existingRole)
			internal.WriteData(w, "", SetupResult{UserID: existingID, Role: existingRole}, nil)
			return nil
		}
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "could not create account", Err: err}
	}

	h.issueSession(w, r, user.ID, user.Role)
	internal.WriteData(w, "", SetupResult{UserID: user.ID, Role: user.Role}, nil)
	return nil
}

func (h *OIDCHandler) finishSignIn(w http.ResponseWriter, r *http.Request, userID uint, role string) {
	if err := h.issueSession(w, r, userID, role); err != nil {
		slog.Error("could not issue session token", "error", err)
		http.Error(w, "could not complete sign-in", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *OIDCHandler) issueSession(w http.ResponseWriter, r *http.Request, userID uint, role string) error {
	token, err := IssueSessionToken(h.JWTSecret, userID, role)
	if err != nil {
		return err
	}
	SetTokenCookie(w, r, token)
	return nil
}

func randomHex() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func setOAuthCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(oauthCookieTTL.Seconds()),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearOAuthCookie(w http.ResponseWriter, r *http.Request, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

func setSetupTokenCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     TokenCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   int(setupTokenTTL.Seconds()),
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}
