// Package routers configure les routes HTTP de l'application cliente.
package routers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterAssetsRoutes expose le répertoire ./static/ pour servir les fichiers statiques.
func RegisterAssetsRoutes(r *mux.Router) {
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}
