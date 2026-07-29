package auth

type MockAuthenticator struct {
	Users map[string]AuthResult
}

func NewMockAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{
		Users: map[string]AuthResult{
			"prof1": {UserID: 1, Username: "prof1", Role: "professor"},
			"stu1":  {UserID: 1, Username: "stu1", Role: "student"},
		},
	}
}

func (m *MockAuthenticator) Authenticate(username, password string) (AuthResult, error) {
	if user, ok := m.Users[username]; ok {
		return user, nil
	}
	return AuthResult{}, ErrInvalidCredentials
}
