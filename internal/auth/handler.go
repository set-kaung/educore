package auth

import (
	"educore/internal"
	"encoding/json"
	"errors"
	"net/http"
)

type AuthHandler struct {
	Authenticator Authenticator
	JWTSecret     string
}

// HandleLogin godoc
// @Summary      Log in
// @Description  Authenticate with username and password, sets the session cookie and returns a JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  LoginRequest  true  "Credentials"
// @Success      200   {object}  internal.ResponseBody{data=LoginResult}
// @Failure      401   {object}  internal.ResponseBody
// @Router       /api/login [post]
func (ah AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid request body", Err: err}
	}

	res, err := Login(ah.Authenticator, req, ah.JWTSecret)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return &internal.HTTPError{StatusCode: http.StatusUnauthorized, Message: "invalid credentials", Err: err}
		}
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "login failed", Err: err}
	}

	SetTokenCookie(w, r, res.Token)
	internal.WriteData(w, "", res, nil)
	return nil
}

// HandleLogout godoc
// @Summary      Log out
// @Description  Clears the session cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  internal.ResponseBody
// @Router       /api/logout [post]
func (ah AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	ClearTokenCookie(w)
	internal.WriteData(w, "logged out", nil, nil)
	return nil
}

// HandleSession godoc
// @Summary      Current session
// @Description  Returns the authenticated user's id and role based on the session cookie or bearer token
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  internal.ResponseBody{data=Claims}
// @Failure      401  {object}  internal.ResponseBody
// @Router       /api/session [get]
func (ah AuthHandler) HandleSession(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	claims := GetClaims(r)
	internal.WriteData(w, "", claims, nil)
	return nil
}
