package web

import (
	"net/http"

	"educore/internal/auth"
)

func RedirectIfAuthenticated(to string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := auth.TryGetClaims(r); ok {
				http.Redirect(w, r, to, http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
