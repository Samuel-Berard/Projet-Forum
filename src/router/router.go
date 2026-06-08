package router

import (
	"database/sql"
	"net/http"

	"projet-forum/src/controllers"
	"projet-forum/src/helper"
	"projet-forum/src/repositories"
	"projet-forum/src/services"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	userRepo := repositories.NewUtilisateurRepository(db)
	filRepo := repositories.NewFilRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	userService := services.NewUtilisateurService(userRepo)
	filService := services.NewFilService(filRepo)
	messageService := services.NewMessageService(messageRepo, filRepo)

	userController := controllers.NewUtilisateurController(userService)
	filController := controllers.NewFilController(filService)
	messageController := controllers.NewMessageController(messageService)

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
