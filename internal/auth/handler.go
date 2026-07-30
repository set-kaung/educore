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
// @Description  Authenticate with username and password, returns JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  auth.LoginRequest  true  "Credentials"
// @Success      200   {object}  internal.ResponseBody{data=auth.LoginResult}
// @Failure      401   {object}  internal.ResponseBody
// @Router       /login [post]
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

	internal.WriteData(w, "", res, nil)
	return nil
}
