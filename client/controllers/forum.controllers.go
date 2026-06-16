// Package controllers contient la logique de présentation pour les pages du forum.
package controllers

import (
	"io"
	"net/http"
	"strconv"
	"strings"

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

// CarteForum représente une carte du "Top forums" déjà préparée pour l'affichage :
// le numéro de classement et la mise en avant sont calculés ici (en Go), pas dans le template.
type CarteForum struct {
	Rang    int
	EnAvant bool
	ID      int
	Titre   string
	Etat    string
}

// CarteActu représente une actualité préparée pour l'affichage (la première est mise en avant).
type CarteActu struct {
	EnAvant bool
	ID      int
	Titre   string
}

// AccueilData regroupe les données affichées sur la page d'accueil.
type AccueilData struct {
	CartesForum []CarteForum          // top 5 pour la grille "bento" (numéro + carte en avant)
	TopForums   []dto.FilDeDiscussion // liste complète pour "Le fil de la commu'"
	CartesActu  []CarteActu           // actualités de la colonne de droite
	Connecte    bool                  // vrai si l'utilisateur est connecté (cookie présent)
}

// DisplayAccueil affiche la page d'accueil.
func (c *ForumControllers) DisplayAccueil(w http.ResponseWriter, r *http.Request) {
	res, err := c.service.GetFils(1, 10)
	var fils []dto.FilDeDiscussion
	if err == nil && res != nil {
		fils = res.Fils
	}

	lastActus, _ := c.service.GetLastActus()

	// On prépare les 5 cartes du "Top forums" avec une boucle Go classique :
	// i donne le numéro de classement (i+1) et désigne la carte mise en avant (i == 0).
	var cartesForum []CarteForum
	for i, fil := range fils {
		if i >= 5 {
			break
		}
		cartesForum = append(cartesForum, CarteForum{
			Rang:    i + 1,
			EnAvant: i == 0,
			ID:      fil.ID,
			Titre:   fil.Titre,
			Etat:    fil.Etat,
		})
	}

	// On prépare les actus : la première (i == 0) passe "À la une".
	var cartesActu []CarteActu
	for i, actu := range lastActus {
		cartesActu = append(cartesActu, CarteActu{
			EnAvant: i == 0,
			ID:      actu.ID,
			Titre:   actu.Titre,
		})
	}

	data := AccueilData{
		CartesForum: cartesForum,
		TopForums:   fils,
		CartesActu:  cartesActu,
		Connecte:    estConnecte(r),
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

// LoginData regroupe les données du formulaire de connexion (pour réafficher en cas d'erreur).
type LoginData struct {
	Erreur      string
	Identifiant string
}

// tokenDuCookie retourne la valeur du cookie "token" (chaîne vide si absent).
func tokenDuCookie(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// estConnecte indique si l'utilisateur possède un cookie "token" (donc connecté).
func estConnecte(r *http.Request) bool {
	return tokenDuCookie(r) != ""
}

// DisplayLogin affiche la page de connexion (formulaire vide).
func (c *ForumControllers) DisplayLogin(w http.ResponseWriter, r *http.Request) {
	c.template.RenderTemplate(w, r, "login", LoginData{})
}

// LoginUser traite l'envoi du formulaire de connexion.
func (c *ForumControllers) LoginUser(w http.ResponseWriter, r *http.Request) {
	// 1. Récupérer et nettoyer les champs.
	identifiant := strings.TrimSpace(r.FormValue("identifiant"))
	password := r.FormValue("password")

	// 2. Valider la présence des champs.
	if identifiant == "" || password == "" {
		c.template.RenderTemplate(w, r, "login", LoginData{
			Erreur:      "Identifiant et mot de passe obligatoires.",
			Identifiant: identifiant,
		})
		return
	}

	// 3. Appeler l'API : on récupère le token JWT en cas de succès.
	token, err := c.service.Login(identifiant, password)
	if err != nil {
		c.template.RenderTemplate(w, r, "login", LoginData{
			Erreur:      err.Error(),
			Identifiant: identifiant,
		})
		return
	}

	// 4. On garde le token dans un cookie pour rester connecté sur les pages suivantes.
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   24 * 60 * 60, // 24 heures
	})

	// 5. Succès : redirection vers l'accueil.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout supprime le cookie de connexion puis redirige vers l'accueil.
func (c *ForumControllers) Logout(w http.ResponseWriter, r *http.Request) {
	// MaxAge négatif demande au navigateur de supprimer le cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// SignupData regroupe les données du formulaire d'inscription.
// Elle sert à réafficher le formulaire (message + valeurs déjà saisies) en cas d'erreur.
type SignupData struct {
	Erreur  string
	Noms    string
	Prenoms string
	Email   string
}

// DisplaySignup affiche la page d'inscription (formulaire vide).
func (c *ForumControllers) DisplaySignup(w http.ResponseWriter, r *http.Request) {
	c.template.RenderTemplate(w, r, "signup", SignupData{})
}

// RegisterUser traite l'envoi du formulaire d'inscription.
func (c *ForumControllers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// 1. Récupérer et nettoyer les champs (capsule : normalisation avec TrimSpace).
	noms := strings.TrimSpace(r.FormValue("noms"))
	prenoms := strings.TrimSpace(r.FormValue("prenoms"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	// 2. Valider la présence des champs obligatoires (capsule : validation côté serveur).
	if noms == "" || prenoms == "" || email == "" || password == "" {
		c.template.RenderTemplate(w, r, "signup", SignupData{
			Erreur:  "Tous les champs sont obligatoires.",
			Noms:    noms,
			Prenoms: prenoms,
			Email:   email,
		})
		return
	}

	// 3. La base n'a qu'un champ "username" : on fusionne prénom + nom pour le construire.
	username := prenoms + " " + noms

	// 4. Appeler l'API d'inscription.
	if err := c.service.Register(username, email, password); err != nil {
		c.template.RenderTemplate(w, r, "signup", SignupData{
			Erreur:  err.Error(),
			Noms:    noms,
			Prenoms: prenoms,
			Email:   email,
		})
		return
	}

	// 5. Succès : redirection vers la page de connexion (redirection après POST).
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ThreadData regroupe les données affichées sur la page d'un fil de discussion.
type ThreadData struct {
	Fil      *dto.FilDeDiscussion
	Messages []dto.Message
	Connecte bool
	Erreur   string
}

// renderThread récupère un fil + ses messages et affiche la page (avec un éventuel message d'erreur).
func (c *ForumControllers) renderThread(w http.ResponseWriter, r *http.Request, id int, erreur string) {
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
		Connecte: estConnecte(r),
		Erreur:   erreur,
	})
}

// DisplayThread affiche un fil de discussion et ses messages.
func (c *ForumControllers) DisplayThread(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Identifiant du fil invalide", http.StatusBadRequest)
		return
	}

	c.renderThread(w, r, id, "")
}

// CreateMessage traite l'envoi du formulaire de réponse dans un fil (action protégée).
func (c *ForumControllers) CreateMessage(w http.ResponseWriter, r *http.Request) {
	// Sans token (cookie), on n'est pas connecté → direction la page de connexion.
	token := tokenDuCookie(r)
	if token == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Identifiant du fil invalide", http.StatusBadRequest)
		return
	}

	contenu := strings.TrimSpace(r.FormValue("contenu"))
	if contenu == "" {
		c.renderThread(w, r, id, "Le message ne peut pas être vide.")
		return
	}

	// On envoie le token à l'API (Authorization: Bearer ...).
	if err := c.service.CreateMessage(token, id, contenu); err != nil {
		c.renderThread(w, r, id, err.Error())
		return
	}

	// Succès : on revient sur le fil, qui affiche le nouveau message.
	http.Redirect(w, r, "/threads/"+strconv.Itoa(id), http.StatusSeeOther)
}

// UploadData regroupe le résultat de l'upload (URL de l'image ou message d'erreur).
type UploadData struct {
	URL    string
	Erreur string
}

// DisplayUpload affiche le formulaire d'upload d'image (page de démonstration).
func (c *ForumControllers) DisplayUpload(w http.ResponseWriter, r *http.Request) {
	c.template.RenderTemplate(w, r, "upload", UploadData{})
}

// HandleUpload reçoit le fichier du formulaire, le transmet à l'API et réaffiche la page avec l'URL.
func (c *ForumControllers) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// Taille maximale : 10 Mo.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		c.template.RenderTemplate(w, r, "upload", UploadData{Erreur: "Fichier trop volumineux ou requête invalide"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		c.template.RenderTemplate(w, r, "upload", UploadData{Erreur: "Aucun fichier sélectionné"})
		return
	}
	defer file.Close()

	// On lit le contenu du fichier pour le transmettre à l'API.
	data, err := io.ReadAll(file)
	if err != nil {
		c.template.RenderTemplate(w, r, "upload", UploadData{Erreur: "Erreur de lecture du fichier"})
		return
	}

	url, err := c.service.UploadImage(header.Filename, data)
	if err != nil {
		c.template.RenderTemplate(w, r, "upload", UploadData{Erreur: err.Error()})
		return
	}

	c.template.RenderTemplate(w, r, "upload", UploadData{URL: url})
}

// SettingsData regroupe les données de la page Paramètres (profil + avatar).
type SettingsData struct {
	Utilisateur *dto.Utilisateur
	Erreur      string
}

// DisplaySettings affiche la page Paramètres de l'utilisateur connecté.
func (c *ForumControllers) DisplaySettings(w http.ResponseWriter, r *http.Request) {
	token := tokenDuCookie(r)
	if token == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := c.service.GetMe(token)
	if err != nil {
		// Token invalide ou expiré → on renvoie vers la connexion.
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	c.template.RenderTemplate(w, r, "settings", SettingsData{Utilisateur: user})
}

// UpdateAvatarSettings traite l'upload d'un nouvel avatar depuis la page Paramètres.
func (c *ForumControllers) UpdateAvatarSettings(w http.ResponseWriter, r *http.Request) {
	token := tokenDuCookie(r)
	if token == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// On récupère l'utilisateur pour pouvoir réafficher la page en cas d'erreur.
	user, err := c.service.GetMe(token)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		c.template.RenderTemplate(w, r, "settings", SettingsData{Utilisateur: user, Erreur: "Fichier trop volumineux ou invalide"})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		c.template.RenderTemplate(w, r, "settings", SettingsData{Utilisateur: user, Erreur: "Aucun fichier sélectionné"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.template.RenderTemplate(w, r, "settings", SettingsData{Utilisateur: user, Erreur: "Erreur de lecture du fichier"})
		return
	}

	// 1. On envoie l'image à l'API → on récupère son URL.
	url, err := c.service.UploadImage(header.Filename, data)
	if err != nil {
		c.template.RenderTemplate(w, r, "settings", SettingsData{Utilisateur: user, Erreur: err.Error()})
		return
	}

	// 2. On enregistre cette URL comme avatar de l'utilisateur.
	if err := c.service.UpdateAvatar(token, url); err != nil {
		c.template.RenderTemplate(w, r, "settings", SettingsData{Utilisateur: user, Erreur: err.Error()})
		return
	}

	// 3. Redirection vers la page Paramètres, qui affiche le nouvel avatar.
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
}
