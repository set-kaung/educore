package auth

import (
	"net/http"
	"strings"
)

const TokenCookie = "edu_token"

// CookiePath returns the cookie scoping path for the given base path:
// "/" when serving at the root, otherwise "/<base>/".
func CookiePath(base string) string {
	b := strings.Trim(base, "/")
	if b == "" {
		return "/"
	}
	return "/" + b + "/"
}

// IsHTTPS reports whether the request arrived via TLS, either directly
// or through a reverse proxy that sets X-Forwarded-Proto.
func IsHTTPS(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	return strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Scope  string `json:"scope,omitempty"`
	OID    string `json:"oid,omitempty"`
	Email  string `json:"email,omitempty"`
	Name   string `json:"name,omitempty"`
}

func SetTokenCookie(w http.ResponseWriter, r *http.Request, token, cookiePath string) {
	http.SetCookie(w, &http.Cookie{
		Name:     TokenCookie,
		Value:    token,
		Path:     cookiePath,
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   IsHTTPS(r),
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearTokenCookie(w http.ResponseWriter, r *http.Request, cookiePath string) {
	http.SetCookie(w, &http.Cookie{
		Name:     TokenCookie,
		Value:    "",
		Path:     cookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   IsHTTPS(r),
		SameSite: http.SameSiteLaxMode,
	})
}
