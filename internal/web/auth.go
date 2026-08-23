package web

import (
	"net/http"

	"educore/internal/auth"
)

func SetAuthCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     auth.TokenCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     auth.TokenCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
