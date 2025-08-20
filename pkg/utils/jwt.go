// pkg/utils/jwt.go

package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GANTI STRUCT INI
type Claims struct {
	UserID int      `json:"user_id"`
	Roles  []string `json:"roles"` // Ubah dari Role string menjadi Roles []string
	jwt.RegisteredClaims
}

func GenerateJWTToken(userID int, role string, secret string, expireTime time.Time) (string, error) {
	claims := &Claims{
		UserID: userID,
		Roles:  []string{role}, // Masukkan ke dalam slice
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseJWTToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}