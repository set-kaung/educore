package auth

import (
	"context"
	"educore/internal"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ClaimsKey contextKey = "claims"

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

type JWTAuth struct {
	secret string
}

func NewJWTAuth(secret string) *JWTAuth {
	return &JWTAuth{secret: secret}
}

func (j *JWTAuth) Middleware() func(http.Handler) http.Handler {
	return j.validate
}

func (j *JWTAuth) validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			internal.WriteError(w, http.StatusUnauthorized, "missing or invalid authorization header", nil)
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(j.secret), nil
		})
		if err != nil || !token.Valid {
			internal.WriteError(w, http.StatusUnauthorized, "invalid or expired token", nil)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			internal.WriteError(w, http.StatusUnauthorized, "invalid token claims", nil)
			return
		}

		c := Claims{
			UserID: uint(claims["user_id"].(float64)),
			Role:   claims["role"].(string),
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaims(r *http.Request) Claims {
	return r.Context().Value(ClaimsKey).(Claims)
}
