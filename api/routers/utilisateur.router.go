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

	// Routes protegees (utilisateur connecte) : profil courant et mise a jour de l'avatar.
	r.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(utilisateurController.Me))).Methods("GET")
	r.Handle("/me/avatar", middleware.AuthMiddleware(http.HandlerFunc(utilisateurController.UpdateAvatar))).Methods("PUT")

	// Routes protegees (admin) : lister les utilisateurs et bannir un compte.
	r.Handle("/users", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(utilisateurController.GetAllUsers)))).Methods("GET")
	r.Handle("/users/{id}/ban", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(utilisateurController.BanUser)))).Methods("PUT")
}
