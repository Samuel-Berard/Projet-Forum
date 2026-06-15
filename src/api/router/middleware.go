package router

import (
	"context"
	"net/http"
	"projet-forum/src/api/auth"
	"projet-forum/src/api/helper"
	"projet-forum/src/api/utils"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.WriteError(w, http.StatusUnauthorized, "Token manquant")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helper.WriteError(w, http.StatusUnauthorized, "Format du token invalide")
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			helper.WriteError(w, http.StatusUnauthorized, "Token invalide ou expiré")
			return
		}

		ctx := context.WithValue(r.Context(), utils.ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := utils.GetClaims(r)
		if claims == nil || claims.Role != "admin" {
			helper.WriteError(w, http.StatusForbidden, "Accès refusé, rôle administrateur requis")
			return
		}
		next.ServeHTTP(w, r)
	})
}
