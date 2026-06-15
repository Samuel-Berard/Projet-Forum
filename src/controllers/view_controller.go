package controllers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"projet-forum/src/models"
	"projet-forum/src/services"
)

type ViewController struct {
	filService *services.FilService
}

func NewViewController(filService *services.FilService) *ViewController {
	return &ViewController{filService: filService}
}

type AccueilData struct {
	TopForums []models.FilDeDiscussion
	LastActus []models.Actualite
}

func (c *ViewController) AfficherAccueil(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		log.Println("Erreur Template :", err)
		return
	}

	fils, _ := c.filService.GetFils(1, 10, "", 0)
	lastActus, _ := c.filService.GetLastActus(10)


	data := AccueilData{
		TopForums: fils,
		LastActus: lastActus,
	}

	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}

type ForumData struct {
	Fils       []models.FilDeDiscussion
	Page       int
	Pages      []int
	TotalPages int
	HasPrev    bool
	PrevPage   int
	HasNext    bool
	NextPage   int
}

func (c *ViewController) AfficherForum(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/forum.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page forum", http.StatusInternalServerError)
		log.Println("Erreur Template :", err)
		return
	}

	page := 1
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	limit := 20
	fils, _ := c.filService.GetFils(page, limit, "", 0)
	totalPages, _ := c.filService.GetTotalPages(limit, "", 0)

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

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
