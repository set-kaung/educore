package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Login(authenticator Authenticator, req LoginRequest, jwtSecret string) (LoginResult, error) {
	result, err := authenticator.Authenticate(req.Username, req.Password)
	if err != nil {
		return LoginResult{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": result.UserID,
		"role":    result.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: signed}, nil
}
