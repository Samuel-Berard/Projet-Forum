package routers

import (
	"projet-forum/api/controllers"

	"github.com/gorilla/mux"
)

// RegisterUploadRoutes enregistre la route d'upload de fichier.
func RegisterUploadRoutes(r *mux.Router, uploadController *controllers.UploadControllers) {
	r.HandleFunc("/upload", uploadController.Upload).Methods("POST")
}
