package auth

import "errors"

var (
	ErrAuthUnavailable     = errors.New("authenticator unavailable")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)

type AuthResult struct {
	UserID   uint
	Username string
	Role     string
}

type Authenticator interface {
	Authenticate(username, password string) (AuthResult, error)
}
