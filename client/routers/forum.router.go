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
}
