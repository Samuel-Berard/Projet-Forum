package controllers

import (
	"encoding/json"
	"net/http"
	"projet-forum/src/api/helper"
	"projet-forum/src/models"
	"projet-forum/src/api/services"
	"projet-forum/src/api/utils"
	"strconv"

	"github.com/gorilla/mux"
)

type MessageController struct {
	service *services.MessageService
}

func NewMessageController(service *services.MessageService) *MessageController {
	return &MessageController{service: service}
}


func readMessageId(r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)["id"])
}

func (c *MessageController) GetMessages(w http.ResponseWriter, r *http.Request) {
	
	filID, err := readMessageId(r)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Identifiant du fil invalide")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	sortBy := r.URL.Query().Get("sort") 

	if sortBy == "" {
		sortBy = "chronologique" 
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

func (c *MessageController) UpdateMessage(w http.ResponseWriter, r *http.Request) {
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

func (c *MessageController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
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

func (c *MessageController) React(w http.ResponseWriter, r *http.Request) {
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
