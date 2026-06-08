package router

import (
	"net/http"

	"projet-forum/src/controllers"
	"projet-forum/src/helper"

	"github.com/gorilla/mux"
)

func NewRouter(
	userController *controllers.UtilisateurController,
	filController *controllers.FilController,
	messageController *controllers.MessageController,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		helper.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Methods("GET")

	r.HandleFunc("/register", userController.Register).Methods("POST")
	r.HandleFunc("/login", userController.Login).Methods("POST")

	r.HandleFunc("/threads", filController.GetFils).Methods("GET")
	r.HandleFunc("/threads/{id}", filController.GetFil).Methods("GET")
	r.HandleFunc("/threads/{id}/messages", messageController.GetMessages).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(AuthMiddleware)

	api.HandleFunc("/threads", filController.CreateFil).Methods("POST")
	api.HandleFunc("/threads/{id}", filController.UpdateFil).Methods("PUT")
	api.HandleFunc("/threads/{id}", filController.DeleteFil).Methods("DELETE")

	api.HandleFunc("/threads/{id}/messages", messageController.CreateMessage).Methods("POST")
	api.HandleFunc("/messages/{id}", messageController.UpdateMessage).Methods("PUT")
	api.HandleFunc("/messages/{id}", messageController.DeleteMessage).Methods("DELETE")
	api.HandleFunc("/messages/{id}/react", messageController.React).Methods("POST")

	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(AdminMiddleware)
	admin.HandleFunc("/threads/{id}/state", filController.ChangeState).Methods("PUT")
	admin.HandleFunc("/messages/{id}", messageController.DeleteMessage).Methods("DELETE")
	admin.HandleFunc("/users/{id}/ban", userController.BanUser).Methods("PUT")

	return r
}
