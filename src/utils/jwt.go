package utils

import (
	"errors"
	"projet-forum/src/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, username, role string) (string, error) {
	secretKey := config.GetEnvWithDefault("JWT_SECRET", "super-secret-key")

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateJWT(tokenString string) (*Claims, error) {
	secretKey := config.GetEnvWithDefault("JWT_SECRET", "super-secret-key")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token invalide")
	}

	return claims, nil
}
