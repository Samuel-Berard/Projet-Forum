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

// Me retourne les informations de l'utilisateur actuellement connecté.
func (c *UtilisateurControllers) Me(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
		return
	}

	id, err := strconv.Atoi(claims.UserID)
	if err != nil {
		helper.WriteError(w, http.StatusUnauthorized, "Identifiant utilisateur invalide")
		return
	}

	user, err := c.service.GetByID(id)
	if err != nil {
		helper.WriteError(w, http.StatusNotFound, "Utilisateur introuvable")
		return
	}

	helper.WriteJSON(w, http.StatusOK, user)
}

// UpdateAvatar met à jour l'avatar de l'utilisateur connecté.
func (c *UtilisateurControllers) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
		return
	}

	id, err := strconv.Atoi(claims.UserID)
	if err != nil {
		helper.WriteError(w, http.StatusUnauthorized, "Identifiant utilisateur invalide")
		return
	}

	var req struct {
		Avatar string `json:"avatar"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Requête invalide")
		return
	}

	if err := c.service.UpdateAvatar(id, req.Avatar); err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "Erreur lors de la mise à jour de l'avatar")
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Avatar mis à jour",
		"avatar":  req.Avatar,
	})
}

// GetAllUsers retourne la liste de tous les utilisateurs (admin uniquement).
func (c *UtilisateurControllers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
		return
	}

	users, err := c.service.GetAllUsers(claims.Role)
	if err != nil {
		helper.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, users)
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
