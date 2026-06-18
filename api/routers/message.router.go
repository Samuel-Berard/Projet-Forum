package routers

import (
	"net/http"

	"projet-forum/api/controllers"
	"projet-forum/api/middleware"

	"github.com/gorilla/mux"
)

// RegisterMessageRoutes enregistre les routes liees aux messages et aux reactions.
func RegisterMessageRoutes(r *mux.Router, messageController *controllers.MessageControllers) {
	// Route publique : consultation des messages d'un fil.
	r.HandleFunc("/threads/{id}/messages", messageController.GetMessages).Methods("GET")

	// Routes protegees : publier, modifier, supprimer, reagir (utilisateur connecte).
	r.Handle("/threads/{id}/messages", middleware.AuthMiddleware(http.HandlerFunc(messageController.CreateMessage))).Methods("POST")
	r.Handle("/messages/{id}", middleware.AuthMiddleware(http.HandlerFunc(messageController.UpdateMessage))).Methods("PUT")
	r.Handle("/messages/{id}", middleware.AuthMiddleware(http.HandlerFunc(messageController.DeleteMessage))).Methods("DELETE")
	r.Handle("/messages/{id}/react", middleware.AuthMiddleware(http.HandlerFunc(messageController.React))).Methods("POST")
}
