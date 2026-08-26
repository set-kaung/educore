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

type JWTAuth struct {
	secret string
}

func NewJWTAuth(secret string) *JWTAuth {
	return &JWTAuth{secret: secret}
}

func (j *JWTAuth) Middleware() func(http.Handler) http.Handler {
	return j.middleware(j.writeJSONError)
}

func (j *JWTAuth) PageMiddleware(loginPath string) func(http.Handler) http.Handler {
	return j.middleware(func(w http.ResponseWriter, r *http.Request, _ int, _ string) {
		http.Redirect(w, r, loginPath, http.StatusSeeOther)
	})
}

func (j *JWTAuth) OptionalMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if tokenStr := j.extractToken(r); tokenStr != "" {
				if c, status, _ := j.parseToken(tokenStr); status == 0 {
					r = r.WithContext(context.WithValue(r.Context(), ClaimsKey, c))
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (j *JWTAuth) middleware(onError func(http.ResponseWriter, *http.Request, int, string)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := j.extractToken(r)
			if tokenStr == "" {
				onError(w, r, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}

			c, statusCode, message := j.parseToken(tokenStr)
			if statusCode != 0 {
				onError(w, r, statusCode, message)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, c)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (j *JWTAuth) parseToken(tokenStr string) (Claims, int, string) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil || !token.Valid {
		return Claims{}, http.StatusUnauthorized, "invalid or expired token"
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, http.StatusUnauthorized, "invalid token claims"
	}

	return Claims{
		UserID: uint(mapClaims["user_id"].(float64)),
		Role:   mapClaims["role"].(string),
	}, 0, ""
}

func (j *JWTAuth) extractToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	if cookie, err := r.Cookie(TokenCookie); err == nil {
		return cookie.Value
	}
	return ""
}

func (j *JWTAuth) writeJSONError(w http.ResponseWriter, _ *http.Request, statusCode int, message string) {
	internal.WriteError(w, statusCode, message, nil)
}

func GetClaims(r *http.Request) Claims {
	claims, ok := TryGetClaims(r)
	if !ok {
		panic("auth claims not found in request context")
	}
	return claims
}

func TryGetClaims(r *http.Request) (Claims, bool) {
	claims, ok := r.Context().Value(ClaimsKey).(Claims)
	return claims, ok
}
