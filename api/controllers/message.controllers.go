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

type MessageControllers struct {
	service *services.MessageService
}

func InitMessageController(service *services.MessageService) *MessageControllers {
	return &MessageControllers{service: service}
}

func readMessageId(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

func (c *MessageControllers) GetMessages(w http.ResponseWriter, r *http.Request) {
	filID, err := readMessageId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du fil invalide")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 10
	}
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "recent"
	}

	currentUserID, _ := strconv.Atoi(r.URL.Query().Get("current_user_id"))

	messages, err := c.service.GetMessagesByFil(filID, page, limit, sortBy, currentUserID)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "Erreur lors de la récupération des messages")
		return
	}

	totalPages, err := c.service.GetTotalPagesForFil(filID, limit)
	if err != nil {
		totalPages = 1
	}

	response := map[string]interface{}{
		"messages":   messages,
		"totalPages": totalPages,
		"page":       page,
		"limit":      limit,
	}

	helper.WriteJSON(w, http.StatusOK, response)
}

func (c *MessageControllers) CreateMessage(w http.ResponseWriter, r *http.Request) {
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

	filID, err := readMessageId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du fil invalide")
		return
	}

	var req models.CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Requête invalide")
		return
	}

	msg, err := c.service.CreateMessage(req.Contenu, filID, userID)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Message publié",
		"data":    msg,
	})
}

func (c *MessageControllers) UpdateMessage(w http.ResponseWriter, r *http.Request) {
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

	idMessage, err := readMessageId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du message invalide")
		return
	}

	var req models.UpdateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Requête invalide")
		return
	}

	err = c.service.UpdateMessage(idMessage, req.Contenu, userID)
	if err != nil {
		helper.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Message mis à jour"})
}

func (c *MessageControllers) DeleteMessage(w http.ResponseWriter, r *http.Request) {
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

	idMessage, err := readMessageId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du message invalide")
		return
	}

	err = c.service.DeleteMessage(idMessage, userID, claims.Role)
	if err != nil {
		helper.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Message supprimé"})
}

func (c *MessageControllers) React(w http.ResponseWriter, r *http.Request) {
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

	idMessage, err := readMessageId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du message invalide")
		return
	}

	var req models.ReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Requête invalide")
		return
	}

	err = c.service.ReactToMessage(userID, idMessage, req.Type)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Réaction enregistrée"})
}
