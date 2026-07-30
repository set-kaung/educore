package authenticators

import "educore/internal/auth"

type MockAuthenticator struct {
	Users map[string]auth.AuthResult
}

func NewMockAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{
		Users: map[string]auth.AuthResult{
			"prof1": {UserID: 1, Username: "prof1", Role: "professor"},
			"stu1":  {UserID: 1, Username: "stu1", Role: "student"},
		},
	}
}

func (m *MockAuthenticator) Authenticate(username, password string) (auth.AuthResult, error) {
	if user, ok := m.Users[username]; ok {
		return user, nil
	}
	return auth.AuthResult{}, auth.ErrInvalidCredentials
}
