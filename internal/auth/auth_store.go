package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func Login(db *gorm.DB, req LoginRequest, jwtSecret string) (LoginResponse, error) {
	var student struct {
		ID       uint
		Password string
	}

	err := db.Table("students").
		Select("id", "password").
		Where("username = ?", req.Username).
		Take(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResponse{}, ErrInvalidCredentials
		}
		return LoginResponse{}, err
	}

	if student.Password != req.Password {
		return LoginResponse{}, ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": student.ID,
		"role":    "student",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{Token: signed}, nil
}
