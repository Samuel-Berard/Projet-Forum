package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("hjgkhukgh")

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, role string) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    "product-api",
			Audience:  []string{"product-front"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
