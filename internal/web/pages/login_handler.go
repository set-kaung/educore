package pages

import (
	"errors"
	"net/http"

	"educore/internal"
	"educore/internal/auth"
	"educore/internal/web"
)

type LoginHandler struct {
	Auth      auth.Authenticator
	JWTSecret string
	Render    *web.Renderer
}

func NewLoginHandler(authenticator auth.Authenticator, jwtSecret string, renderer *web.Renderer) *LoginHandler {
	return &LoginHandler{Auth: authenticator, JWTSecret: jwtSecret, Render: renderer}
}

func (h *LoginHandler) Home(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	target := "/login"
	if _, authenticated := auth.TryGetClaims(r); authenticated {
		target = "/students"
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
	return nil
}

func (h *LoginHandler) Show(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	data := web.BaseData(r, "Log in")
	data.Flash = web.PopFlash(w, r)
	h.Render.Page(w, http.StatusOK, "login", data)
	return nil
}

func (h *LoginHandler) Submit(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	if err := r.ParseForm(); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid form submission", Err: err}
	}

	result, err := auth.Login(h.Auth, auth.LoginRequest{
		Username: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	}, h.JWTSecret)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			data := web.BaseData(r, "Log in")
			data.Error = "Invalid username or password."
			h.Render.Page(w, http.StatusUnauthorized, "login", data)
			return nil
		}
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "login failed", Err: err}
	}

	web.SetAuthCookie(w, r, result.Token)
	http.Redirect(w, r, "/students", http.StatusSeeOther)
	return nil
}

func (h *LoginHandler) Logout(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	web.ClearAuthCookie(w)
	web.SetFlash(w, "You have been logged out.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}
