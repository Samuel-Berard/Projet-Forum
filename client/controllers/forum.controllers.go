// Package controllers contient la logique de présentation pour les pages du forum.
package controllers

import (
	"net/http"
	"strconv"

	"projet-forum/client/dto"
	"projet-forum/client/services"
	"projet-forum/client/templates"

	"github.com/gorilla/mux"
)

// ForumControllers gère les actions web liées au forum.
type ForumControllers struct {
	service  *services.ForumService
	template *templates.TemplateManager
}

// InitForumController crée un contrôleur forum avec son service et son moteur de templates.
func InitForumController(service *services.ForumService, template *templates.TemplateManager) *ForumControllers {
	return &ForumControllers{service: service, template: template}
}

// AccueilData regroupe les données affichées sur la page d'accueil.
type AccueilData struct {
	TopForums []dto.FilDeDiscussion
	LastActus []dto.Actualite
}

// DisplayAccueil affiche la page d'accueil.
func (c *ForumControllers) DisplayAccueil(w http.ResponseWriter, r *http.Request) {
	res, err := c.service.GetFils(1, 10)
	var fils []dto.FilDeDiscussion
	if err == nil && res != nil {
		fils = res.Fils
	}

	lastActus, _ := c.service.GetLastActus()

	data := AccueilData{
		TopForums: fils,
		LastActus: lastActus,
	}

	c.template.RenderTemplate(w, r, "index", data)
}

// ForumData regroupe les données affichées sur la page forum (liste + pagination).
type ForumData struct {
	Fils       []dto.FilDeDiscussion
	Page       int
	Pages      []int
	TotalPages int
	HasPrev    bool
	PrevPage   int
	HasNext    bool
	NextPage   int
}

// DisplayForum affiche la liste des fils avec pagination.
func (c *ForumControllers) DisplayForum(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	limit := 20
	res, err := c.service.GetFils(page, limit)

	var fils []dto.FilDeDiscussion
	totalPages := 1

	if err == nil && res != nil {
		fils = res.Fils
		if res.TotalPages > 0 {
			totalPages = res.TotalPages
		}
	}

	var pages []int
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, i)
	}

	data := ForumData{
		Fils:       fils,
		Page:       page,
		Pages:      pages,
		TotalPages: totalPages,
		HasPrev:    page > 1,
		PrevPage:   page - 1,
		HasNext:    page < totalPages,
		NextPage:   page + 1,
	}

	c.template.RenderTemplate(w, r, "forum", data)
}

// DisplayLogin affiche la page de connexion.
func (c *ForumControllers) DisplayLogin(w http.ResponseWriter, r *http.Request) {
	c.template.RenderTemplate(w, r, "login", nil)
}

// DisplaySignup affiche la page d'inscription.
func (c *ForumControllers) DisplaySignup(w http.ResponseWriter, r *http.Request) {
	c.template.RenderTemplate(w, r, "signup", nil)
}

// ThreadData regroupe les données affichées sur la page d'un fil de discussion.
type ThreadData struct {
	Fil      *dto.FilDeDiscussion
	Messages []dto.Message
}

// DisplayThread affiche un fil de discussion et ses messages.
func (c *ForumControllers) DisplayThread(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Identifiant du fil invalide", http.StatusBadRequest)
		return
	}

	fil, err := c.service.GetFil(id)
	if err != nil || fil == nil {
		http.Error(w, "Fil de discussion introuvable", http.StatusNotFound)
		return
	}

	// On récupère les messages ; en cas d'erreur on affiche le fil sans messages.
	messages, _ := c.service.GetMessages(id)

	c.template.RenderTemplate(w, r, "thread", ThreadData{
		Fil:      fil,
		Messages: messages,
	})
}
