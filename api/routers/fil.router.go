package routers

import (
	"net/http"

	"projet-forum/api/controllers"
	"projet-forum/api/middleware"

	"github.com/gorilla/mux"
)

// RegisterFilRoutes enregistre les routes liees aux fils de discussion.
func RegisterFilRoutes(r *mux.Router, filController *controllers.FilControllers) {
	// Routes publiques : consultation des fils.
	r.HandleFunc("/threads", filController.GetFils).Methods("GET")
	r.HandleFunc("/threads/{id}", filController.GetFil).Methods("GET")

	// Routes protegees : creation, modification, suppression (utilisateur connecte).
	r.Handle("/threads", middleware.AuthMiddleware(http.HandlerFunc(filController.CreateFil))).Methods("POST")
	r.Handle("/threads/{id}", middleware.AuthMiddleware(http.HandlerFunc(filController.UpdateFil))).Methods("PUT")
	r.Handle("/threads/{id}", middleware.AuthMiddleware(http.HandlerFunc(filController.DeleteFil))).Methods("DELETE")

	// Route protegee (admin) : changer l'etat d'un fil.
	r.Handle("/threads/{id}/state", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(filController.ChangeState)))).Methods("PUT")
}
