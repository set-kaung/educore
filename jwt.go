package main

import (
	"context"
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
			http.Error(w, `{"status_code":401,"status":"Unauthorized","message":"missing or invalid authorization header"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(j.secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, `{"status_code":401,"status":"Unauthorized","message":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"status_code":401,"status":"Unauthorized","message":"invalid token claims"}`, http.StatusUnauthorized)
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
