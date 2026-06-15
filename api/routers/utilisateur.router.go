package routers

import (
	"net/http"

	"projet-forum/api/controllers"
	"projet-forum/api/middleware"

	"github.com/gorilla/mux"
)

// RegisterUtilisateurRoutes enregistre les routes liees aux utilisateurs et a l'authentification.
func RegisterUtilisateurRoutes(r *mux.Router, utilisateurController *controllers.UtilisateurControllers) {
	// Routes publiques : inscription et connexion.
	r.HandleFunc("/register", utilisateurController.Register).Methods("POST")
	r.HandleFunc("/login", utilisateurController.Login).Methods("POST")

	// Route protegee (admin) : bannir un utilisateur.
	r.Handle("/users/{id}/ban", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(utilisateurController.BanUser)))).Methods("PUT")
}
