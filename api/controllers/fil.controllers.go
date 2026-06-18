package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"projet-forum/api/helper"
	"projet-forum/api/models"
	"projet-forum/api/services"
	"projet-forum/api/utils"

	"github.com/gorilla/mux"
)

type FilControllers struct {
	service *services.FilService
}

func InitFilController(service *services.FilService) *FilControllers {
	return &FilControllers{service: service}
}

func readFilId(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

func (c *FilControllers) GetFils(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	search := r.URL.Query().Get("search")
	searchType := r.URL.Query().Get("type")
	if searchType == "" {
		searchType = "topics"
	}
	categoryID, _ := strconv.Atoi(r.URL.Query().Get("category"))

	fils, err := c.service.GetFils(page, limit, search, searchType, categoryID)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "Erreur lors de la récupération des fils de discussion")
		return
	}

	totalPages, err := c.service.GetTotalPages(limit, search, searchType, categoryID)
	if err != nil {
		totalPages = 1
	}

	response := map[string]interface{}{
		"fils":       fils,
		"totalPages": totalPages,
		"page":       page,
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (c *FilControllers) GetFil(w http.ResponseWriter, r *http.Request) {
	idFil, err := readFilId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du fil invalide")
		return
	}

	fil, err := c.service.GetFilByID(idFil)
	if err != nil {
		helper.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, fil)
}

func (c *FilControllers) CreateFil(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
		return
	}

	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		helper.WriteError(w, http.StatusUnauthorized, "Identifiant utilisateur invalide")
		return
	}

	var req models.CreateFilRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Requête invalide")
		return
	}

	fil, err := c.service.CreateFil(req.Titre, userID, req.CategoriesID)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Fil de discussion créé avec succès",
		"fil":     fil,
	})
}

func (c *FilControllers) UpdateFil(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
		return
	}

	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		helper.WriteError(w, http.StatusUnauthorized, "Identifiant utilisateur invalide")
		return
	}

	idFil, err := readFilId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du fil invalide")
		return
	}

	var req models.UpdateFilRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Requête invalide")
		return
	}

	err = c.service.UpdateFil(idFil, req.Titre, userID, claims.Role)
	if err != nil {
		helper.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Fil de discussion mis à jour"})
}

func (c *FilControllers) DeleteFil(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
		return
	}

	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		helper.WriteError(w, http.StatusUnauthorized, "Identifiant utilisateur invalide")
		return
	}

	idFil, err := readFilId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du fil invalide")
		return
	}

	err = c.service.DeleteFil(idFil, userID, claims.Role)
	if err != nil {
		helper.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Fil de discussion supprimé"})
}

func (c *FilControllers) ChangeState(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
		return
	}

	idFil, err := readFilId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du fil invalide")
		return
	}

	var req struct {
		Etat string `json:"etat"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Requête invalide")
		return
	}

	err = c.service.ChangeFilState(idFil, req.Etat, claims.Role)
	if err != nil {
		helper.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "État du fil mis à jour"})
}
