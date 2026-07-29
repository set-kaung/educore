package auth

import "log/slog"

type ADAuthenticator struct {
	ClientID     string
	ClientSecret string
	TenantID     string
}

func NewADAuthenticator(clientID, clientSecret, tenantID string) *ADAuthenticator {
	slog.Warn("ADAuthenticator is a placeholder — not yet implemented")
	return &ADAuthenticator{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TenantID:     tenantID,
	}
}

func (a *ADAuthenticator) Authenticate(username, password string) (AuthResult, error) {
	return AuthResult{}, ErrAuthUnavailable
}
