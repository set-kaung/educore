package main

import (
	"educore/internal"
	"educore/internal/auth"
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"
)

type AuthHandler struct {
	db        *gorm.DB
	jwtSecret string
}

func (ah AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) *internal.HTTPError {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return &internal.HTTPError{StatusCode: http.StatusBadRequest, Message: "invalid request body", Err: err}
	}

	res, err := auth.Login(ah.db, req, ah.jwtSecret)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return &internal.HTTPError{StatusCode: http.StatusUnauthorized, Message: "invalid credentials", Err: err}
		}
		return &internal.HTTPError{StatusCode: http.StatusInternalServerError, Message: "login failed", Err: err}
	}

	internal.WriteData(w, "login successful", res, nil)
	return nil
}
