package controllers

import (
	"net/http"

	"projet-forum/api/helper"
	"projet-forum/api/services"
)

type ActusControllers struct {
	service *services.ActusService
}

func InitActusController(service *services.ActusService) *ActusControllers {
	return &ActusControllers{service: service}
}

// GetActus retourne la liste des actualités.
func (c *ActusControllers) GetActus(w http.ResponseWriter, r *http.Request) {
	actus := c.service.GetLastActus()
	helper.WriteJSON(w, http.StatusOK, actus)
}
