package routers

import (
	"projet-forum/api/controllers"

	"github.com/gorilla/mux"
)

// RegisterActusRoutes enregistre la route publique des actualités.
func RegisterActusRoutes(r *mux.Router, actusController *controllers.ActusControllers) {
	r.HandleFunc("/actus", actusController.GetActus).Methods("GET")
}
