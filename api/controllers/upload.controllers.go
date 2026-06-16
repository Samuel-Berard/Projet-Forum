package controllers

import (
	"net/http"

	"projet-forum/api/helper"
	"projet-forum/api/services"
)

type UploadControllers struct {
	service *services.UploadService
}

func InitUploadController(service *services.UploadService) *UploadControllers {
	return &UploadControllers{service: service}
}

// Upload reçoit un fichier en multipart/form-data (champ "file"),
// l'enregistre dans le dossier uploads et renvoie son URL publique.
func (c *UploadControllers) Upload(w http.ResponseWriter, r *http.Request) {
	// Taille maximale en mémoire : 10 Mo.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Fichier trop volumineux ou requête invalide")
		return
	}

	// Récupération du fichier envoyé sous le champ "file".
	file, header, err := r.FormFile("file")
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Aucun fichier reçu (champ 'file' attendu)")
		return
	}
	defer file.Close()

	chemin, err := c.service.SaveImage(file, header)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// URL absolue servie par l'API, utilisable directement dans une balise <img>.
	urlPublique := "http://" + r.Host + chemin

	helper.WriteJSON(w, http.StatusCreated, map[string]string{"url": urlPublique})
}
