package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const sessionTokenTTL = 24 * time.Hour
const setupTokenTTL = 15 * time.Minute

func IssueSessionToken(jwtSecret string, userID uint, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(sessionTokenTTL).Unix(),
	})
	return token.SignedString([]byte(jwtSecret))
}

func IssueSetupToken(jwtSecret, oid, email, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"scope": "setup",
		"oid":   oid,
		"email": email,
		"name":  name,
		"exp":   time.Now().Add(setupTokenTTL).Unix(),
	})
	return token.SignedString([]byte(jwtSecret))
}
