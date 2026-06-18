// Package routers configure les routes de l'application cliente.
package routers

import (
	"projet-forum/client/controllers"

	"github.com/gorilla/mux"
)

// RegisterForumRoutes enregistre les routes des pages du forum.
func RegisterForumRoutes(r *mux.Router, forumController *controllers.ForumControllers) {
	r.HandleFunc("/", forumController.DisplayAccueil).Methods("GET")
	r.HandleFunc("/login", forumController.DisplayLogin).Methods("GET")
	r.HandleFunc("/login", forumController.LoginUser).Methods("POST")
	r.HandleFunc("/logout", forumController.Logout).Methods("GET")
	r.HandleFunc("/signup", forumController.DisplaySignup).Methods("GET")
	r.HandleFunc("/signup", forumController.RegisterUser).Methods("POST")
	r.HandleFunc("/forum", forumController.DisplayForum).Methods("GET")
	r.HandleFunc("/threads/{id}", forumController.DisplayThread).Methods("GET")
	r.HandleFunc("/threads/{id}/messages", forumController.CreateMessage).Methods("POST")
	r.HandleFunc("/upload", forumController.DisplayUpload).Methods("GET")
	r.HandleFunc("/upload", forumController.HandleUpload).Methods("POST")
	r.HandleFunc("/settings", forumController.DisplaySettings).Methods("GET")
	r.HandleFunc("/settings", forumController.UpdateAvatarSettings).Methods("POST")
	r.HandleFunc("/threads/new", forumController.CreateThread).Methods("POST")

	// Routes d'action modération/édition
	r.HandleFunc("/threads/{id}/delete", forumController.DeleteThread).Methods("POST")
	r.HandleFunc("/threads/{id}/edit", forumController.EditThread).Methods("POST")
	r.HandleFunc("/messages/{id}/delete", forumController.DeleteMessage).Methods("POST")
	r.HandleFunc("/messages/{id}/edit", forumController.EditMessage).Methods("POST")
	r.HandleFunc("/messages/{id}/react", forumController.ReactToMessage).Methods("POST")
	r.HandleFunc("/threads/{id}/state", forumController.ChangeThreadState).Methods("POST")
	r.HandleFunc("/users/{id}/ban", forumController.BanUser).Methods("POST")

	// Tableau de bord d'administration (FT-12) : accès réservé aux admins (vérifié dans le contrôleur).
	r.HandleFunc("/admin", forumController.DisplayAdmin).Methods("GET")
}
