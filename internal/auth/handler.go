package auth

import (
	"educore/internal"
	"net/http"
)

type AuthHandler struct{}

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
