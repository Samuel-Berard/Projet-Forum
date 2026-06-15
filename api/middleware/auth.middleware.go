package middleware

import (
	"context"
	"net/http"
	"strings"

	"projet-forum/api/auth"
	"projet-forum/api/helper"
	"projet-forum/api/utils"
)

// AuthMiddleware protege une route en exigeant un JWT valide.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Le token doit etre envoye dans le header Authorization.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.WriteError(w, http.StatusUnauthorized, "Token manquant")
			return
		}

		// Le format attendu est : Bearer <token>.
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			helper.WriteError(w, http.StatusUnauthorized, "Format du token invalide")
			return
		}

		// On valide le token et on recupere les claims de l'utilisateur.
		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			helper.WriteError(w, http.StatusUnauthorized, "Token invalide ou expiré")
			return
		}

		// Les claims sont ajoutes au contexte pour etre reutilises par les handlers suivants.
		ctx := context.WithValue(r.Context(), utils.ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminMiddleware verifie que l'utilisateur authentifie possede le role administrateur.
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
