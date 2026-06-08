package utils

import (
	"net/http"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func GetClaims(r *http.Request) *Claims {
	claims, ok := r.Context().Value(ClaimsKey).(*Claims)
	if !ok {
		return nil
	}
	return claims
}
