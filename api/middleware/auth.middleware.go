package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"projet-forum/api/auth"
	"projet-forum/api/helper"
	"projet-forum/api/utils"
)

var DB *sql.DB

// Init initialise la connexion de base de données pour le middleware.
func Init(db *sql.DB) {
	DB = db
}

// AuthMiddleware protege une route en exigeant un JWT valide et en verifiant que l'utilisateur n'est pas banni.
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

		// On verifie en temps reel si l'utilisateur est banni en base de donnees.
		if DB != nil {
			var banned bool
			err := DB.QueryRow("SELECT banned FROM utilisateurs WHERE id = ?", claims.UserID).Scan(&banned)
			if err != nil {
				if err == sql.ErrNoRows {
					helper.WriteError(w, http.StatusUnauthorized, "Utilisateur introuvable")
				} else {
					helper.WriteError(w, http.StatusInternalServerError, "Erreur serveur lors de la vérification de l'utilisateur")
				}
				return
			}
			if banned {
				helper.WriteError(w, http.StatusForbidden, "Ce compte a été banni")
				return
			}
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
