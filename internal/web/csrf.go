package web

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"log/slog"
	"net/http"
)

const (
	CSRFCookie = "edu_csrf"
	CSRFField  = "_csrf"
)

type csrfKeyType struct{}

var csrfKey csrfKeyType

type CSRF struct {
	secret []byte
}

func NewCSRF(secret string) *CSRF {
	return &CSRF{secret: []byte(secret)}
}

func (c *CSRF) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := c.ensureCookie(w, r)
		r = r.WithContext(context.WithValue(r.Context(), csrfKey, c.sign(value)))

		if !safeMethod(r.Method) && !c.valid(r) {
			http.Error(w, "invalid or missing CSRF token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CSRFToken(r *http.Request) string {
	token, _ := r.Context().Value(csrfKey).([]byte)
	if token == nil {
		return ""
	}
	return string(token)
}

func (c *CSRF) ensureCookie(w http.ResponseWriter, r *http.Request) string {
	if cookie, err := r.Cookie(CSRFCookie); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		slog.Error("generate csrf token", "error", err)
		return ""
	}

	value := hex.EncodeToString(raw)
	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookie,
		Value:    value,
		Path:     "/",
		MaxAge:   43200,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
	return value
}

func (c *CSRF) valid(r *http.Request) bool {
	expected, _ := r.Context().Value(csrfKey).([]byte)
	token := r.PostFormValue(CSRFField)
	return expected != nil &&
		token != "" &&
		subtle.ConstantTimeCompare(expected, []byte(token)) == 1
}

func (c *CSRF) sign(value string) []byte {
	mac := hmac.New(sha256.New, c.secret)
	mac.Write([]byte(value))
	return []byte(hex.EncodeToString(mac.Sum(nil)))
}

func safeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}
