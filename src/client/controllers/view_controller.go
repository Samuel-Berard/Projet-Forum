package controllers

import (
	"html/template"
	"log"
	"net/http"
	"projet-forum/src/client/services"
	"projet-forum/src/models"
	"strconv"
)

type ViewController struct {
	forumService *services.ForumClientService
}

func NewViewController(forumService *services.ForumClientService) *ViewController {
	return &ViewController{forumService: forumService}
}

type AccueilData struct {
	TopForums []models.FilDeDiscussion
	LastActus []models.FilDeDiscussion
}

func (c *ViewController) AfficherAccueil(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		log.Println("Erreur Template :", err)
		return
	}

	res, err := c.forumService.GetFils(1, 10)
	var fils []models.FilDeDiscussion
	if err == nil && res != nil {
		fils = res.Fils
	}

	lastActus, _ := c.forumService.GetLastActus(10)

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
	res, err := c.forumService.GetFils(page, limit)
	
	var fils []models.FilDeDiscussion
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

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}

func (c *ViewController) AfficherLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page login", http.StatusInternalServerError)
		log.Println("Erreur Template :", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "login", nil)
	if err != nil {
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}

func (c *ViewController) AfficherSignup(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/Signup.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page signup", http.StatusInternalServerError)
		log.Println("Erreur Template :", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "signup", nil)
	if err != nil {
		http.Error(w, "Erreur lors de l'affichage de la page", http.StatusInternalServerError)
	}
}
