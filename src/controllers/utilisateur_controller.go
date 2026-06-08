package controllers

import (
	"encoding/json"
	"net/http"
	"projet-forum/src/helper"
	"projet-forum/src/models"
	"projet-forum/src/services"
	"projet-forum/src/utils"
	"strconv"

	"github.com/gorilla/mux"
)

type UtilisateurController struct {
	service *services.UtilisateurService
}

func NewUtilisateurController(service *services.UtilisateurService) *UtilisateurController {
	return &UtilisateurController{service: service}
}

// readUtilisateurId lit l'identifiant de l'utilisateur depuis l'URL (méthode du prof)
func readUtilisateurId(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

func (c *UtilisateurController) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Requête invalide")
		return
	}

	user, err := c.service.Register(&req)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Inscription réussie",
		"user":    user,
	})
}

func (c *UtilisateurController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Requête invalide")
		return
	}

	token, err := c.service.Login(&req)
	if err != nil {
		helper.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Connexion réussie",
		"token":   token,
	})
}

func (c *UtilisateurController) BanUser(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
		return
	}

	idUtilisateur, err := readUtilisateurId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant de l'utilisateur invalide")
		return
	}

	err = c.service.BanUser(idUtilisateur, claims.Role)
	if err != nil {
		helper.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Utilisateur banni"})
}
