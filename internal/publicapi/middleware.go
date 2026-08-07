package publicapi

import (
	"educore/internal"
	"net/http"

	"gorm.io/gorm"
)

func APIKeyMiddleware(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("x-api-key")
			if key == "" {
				internal.WriteError(w, http.StatusUnauthorized, "missing x-api-key header", nil)
				return
			}

			valid, err := ValidateAPIKey(db, key)
			if err != nil {
				internal.WriteError(w, http.StatusInternalServerError, "could not validate api key", nil)
				return
			}
			if !valid {
				internal.WriteError(w, http.StatusUnauthorized, "invalid or revoked api key", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
