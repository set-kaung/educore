package web

import (
	"net/http"
	"net/url"
)

const FlashCookie = "edu_flash"

func SetFlash(w http.ResponseWriter, message string) {
	http.SetCookie(w, &http.Cookie{
		Name:     FlashCookie,
		Value:    url.QueryEscape(message),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

func PopFlash(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie(FlashCookie)
	if err != nil || cookie.Value == "" {
		return ""
	}

	http.SetCookie(w, &http.Cookie{
		Name:     FlashCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	message, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return ""
	}
	return message
}
