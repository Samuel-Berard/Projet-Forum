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

type MessageController struct {
	service *services.MessageService
}

func NewMessageController(service *services.MessageService) *MessageController {
	return &MessageController{service: service}
}

// readMessageId lit l'identifiant du message depuis l'URL (méthode du prof)
func readMessageId(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

func (c *MessageController) GetMessages(w http.ResponseWriter, r *http.Request) {
	// L'ID ici c'est celui du fil (route: /threads/{id}/messages)
	filID, err := readMessageId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du fil invalide")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	sortBy := r.URL.Query().Get("sort") // "chronologique" ou "popularite"

	if sortBy == "" {
		sortBy = "chronologique" // par défaut FT-8
	}

	messages, err := c.service.GetMessagesByFil(filID, page, limit, sortBy)
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "Erreur lors de la récupération des messages")
		return
	}

	helper.WriteJSON(w, http.StatusOK, messages)
}

func (c *MessageController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
		return
	}

	// L'ID ici c'est celui du fil (route: /threads/{id}/messages)
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

	msg, err := c.service.CreateMessage(req.Contenu, filID, claims.UserID)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Message publié",
		"data":    msg,
	})
}

func (c *MessageController) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
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

	err = c.service.UpdateMessage(idMessage, req.Contenu, claims.UserID)
	if err != nil {
		helper.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Message mis à jour"})
}

func (c *MessageController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
		return
	}

	idMessage, err := readMessageId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du message invalide")
		return
	}

	err = c.service.DeleteMessage(idMessage, claims.UserID, claims.Role)
	if err != nil {
		helper.WriteError(w, http.StatusForbidden, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Message supprimé"})
}

func (c *MessageController) React(w http.ResponseWriter, r *http.Request) {
	claims := utils.GetClaims(r)
	if claims == nil {
		helper.WriteError(w, http.StatusUnauthorized, "Non autorisé")
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

	err = c.service.ReactToMessage(claims.UserID, idMessage, req.Type)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	helper.WriteJSON(w, http.StatusOK, map[string]string{"message": "Réaction enregistrée"})
}
