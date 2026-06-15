package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("hjgkhukgh")

// Claims represente les informations stockees dans le JWT.
// Les champs personnalises permettent d'identifier l'utilisateur et son role.
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken cree un JWT signe pour un utilisateur donne.
func GenerateToken(userID string, role string) (string, error) {
	now := time.Now()

	// On prepare les donnees qui seront encodees dans le token.
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// Subject identifie le proprietaire du token.
			Subject:  userID,
			Issuer:   "product-api",
			Audience: []string{"product-front"},
			IssuedAt: jwt.NewNumericDate(now),
			// Le token expire apres 15 minutes pour limiter les risques en cas de vol.
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}

	// On signe le token avec l'algorithme HS256 et la cle secrete.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
