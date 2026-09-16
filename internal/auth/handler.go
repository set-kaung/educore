package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"educore/internal"
)

type AuthHandler struct {
	BasePath string
}

// HandleLogout godoc
// @Summary      Log out
// @Description  Clears the session cookie
// @Tags         auth
// @Produce      json
// @Success      200  {object}  internal.ResponseBody
// @Router       /api/logout [post]
func (ah AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	cookiePath := CookiePath(ah.BasePath)

	ClearTokenCookie(w, r, cookiePath)
	if cookiePath != "/" {
		ClearTokenCookie(w, r, "/")
	}
	internal.WriteData(w, "logged out", nil, nil)
	return nil
}

// cookieOccurrences counts how many times a cookie name (e.g. "edu_token=")
// appears in a Cookie header. Multiple occurrences mean the browser holds
// the same-named cookie under different Path/Domain scopes.
func cookieOccurrences(header, name string) int {
	if header == "" {
		return 0
	}
	count := 0
	for _, part := range strings.Split(header, ";") {
		kv := strings.TrimSpace(part)
		if strings.HasPrefix(kv, name+"=") {
			count++
		}
	}
	return count
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
	count := cookieOccurrences(r.Header.Get("Cookie"), TokenCookie)
	slog.Info("session", "method", r.Method, "path", r.URL.Path, "tokenCookieCount", count, "userID", claims.UserID, "role", claims.Role)
	internal.WriteData(w, "", claims, nil)
	return nil
}
