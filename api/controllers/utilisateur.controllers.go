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

type UtilisateurControllers struct {
	service *services.UtilisateurService
}

func InitUtilisateurController(service *services.UtilisateurService) *UtilisateurControllers {
	return &UtilisateurControllers{service: service}
}

func readUtilisateurId(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

// Register gere l'inscription d'un nouvel utilisateur.
func (c *UtilisateurControllers) Register(w http.ResponseWriter, r *http.Request) {
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

// Login gere la connexion d'un utilisateur et retourne un token JWT.
func (c *UtilisateurControllers) Login(w http.ResponseWriter, r *http.Request) {
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

// BanUser bannit un utilisateur, reserve aux administrateurs.
func (c *UtilisateurControllers) BanUser(w http.ResponseWriter, r *http.Request) {
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
