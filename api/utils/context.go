package utils

import (
	"net/http"

	"projet-forum/api/auth"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func GetClaims(r *http.Request) *auth.Claims {
	claims, ok := r.Context().Value(ClaimsKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}
