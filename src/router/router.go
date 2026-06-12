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
	viewController *controllers.ViewController,
) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", viewController.AfficherAccueil).Methods("GET")
	r.HandleFunc("/forum", viewController.AfficherForum).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	apiPublic := r.PathPrefix("/api").Subrouter()

	apiPublic.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		helper.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Methods("GET")

	apiPublic.HandleFunc("/register", userController.Register).Methods("POST")
	apiPublic.HandleFunc("/login", userController.Login).Methods("POST")

	apiPublic.HandleFunc("/threads", filController.GetFils).Methods("GET")
	apiPublic.HandleFunc("/threads/{id}", filController.GetFil).Methods("GET")
	apiPublic.HandleFunc("/threads/{id}/messages", messageController.GetMessages).Methods("GET")

	apiPrivate := r.PathPrefix("/api").Subrouter()
	apiPrivate.Use(AuthMiddleware)

	apiPrivate.HandleFunc("/threads", filController.CreateFil).Methods("POST")
	apiPrivate.HandleFunc("/threads/{id}", filController.UpdateFil).Methods("PUT")
	apiPrivate.HandleFunc("/threads/{id}", filController.DeleteFil).Methods("DELETE")

	apiPrivate.HandleFunc("/threads/{id}/messages", messageController.CreateMessage).Methods("POST")
	apiPrivate.HandleFunc("/messages/{id}", messageController.UpdateMessage).Methods("PUT")
	apiPrivate.HandleFunc("/messages/{id}", messageController.DeleteMessage).Methods("DELETE")
	apiPrivate.HandleFunc("/messages/{id}/react", messageController.React).Methods("POST")

	admin := apiPrivate.PathPrefix("/admin").Subrouter()
	admin.Use(AdminMiddleware)
	admin.HandleFunc("/threads/{id}/state", filController.ChangeState).Methods("PUT")
	admin.HandleFunc("/messages/{id}", messageController.DeleteMessage).Methods("DELETE")
	admin.HandleFunc("/users/{id}/ban", userController.BanUser).Methods("PUT")

	return r
}
