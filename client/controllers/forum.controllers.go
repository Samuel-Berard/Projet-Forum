// Package controllers contient la logique de présentation pour les pages du forum.
package controllers

import (
	"io"
	"net/http"
	"net/url"
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
	Utilisateur *dto.Utilisateur      // informations de l'utilisateur connecté
}

// DisplayAccueil affiche la page d'accueil.
func (c *ForumControllers) DisplayAccueil(w http.ResponseWriter, r *http.Request) {
	res, err := c.service.GetFils(1, 10, "", "", 0)
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
		Utilisateur: c.getConnecteUser(r),
	}

	c.template.RenderTemplate(w, r, "index", data)
}

// PageLien représente un numéro de page de la pagination.
// On le prépare côté Go (numéro + page courante) pour garder le template simple.
type PageLien struct {
	Num      int
	Courante bool
}

// ForumData regroupe les données affichées sur la page forum (liste + pagination).
type ForumData struct {
	Fils        []dto.FilDeDiscussion
	Page        int
	Pages       []PageLien
	TotalPages  int
	HasPrev     bool
	PrevPage    int
	HasNext     bool
	NextPage    int
	Connecte    bool
	Utilisateur *dto.Utilisateur
	Erreur      string
	Search      string
	SearchType  string
}

// DisplayForum affiche la liste des fils avec pagination.
func (c *ForumControllers) DisplayForum(w http.ResponseWriter, r *http.Request) {
	page := 1
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	search := strings.TrimSpace(r.URL.Query().Get("q"))
	if search == "" {
		search = strings.TrimSpace(r.URL.Query().Get("recherche"))
	}

	searchType := strings.TrimSpace(r.URL.Query().Get("type"))
	if searchType == "" {
		searchType = "topics"
	}

	categoryID := 0
	if cat, err := strconv.Atoi(r.URL.Query().Get("category")); err == nil && cat > 0 {
		categoryID = cat
	}

	limit := 20
	res, err := c.service.GetFils(page, limit, search, searchType, categoryID)

	var fils []dto.FilDeDiscussion
	totalPages := 1

	if err == nil && res != nil {
		fils = res.Fils
		if res.TotalPages > 0 {
			totalPages = res.TotalPages
		}
	}

	// On prépare les numéros de page côté Go ; le template n'a plus qu'à les afficher.
	var pages []PageLien
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, PageLien{Num: i, Courante: i == page})
	}

	data := ForumData{
		Fils:        fils,
		Page:        page,
		Pages:       pages,
		TotalPages:  totalPages,
		HasPrev:     page > 1,
		PrevPage:    page - 1,
		HasNext:     page < totalPages,
		NextPage:    page + 1,
		Connecte:    estConnecte(r),
		Utilisateur: c.getConnecteUser(r),
		Erreur:      r.URL.Query().Get("erreur"),
		Search:      search,
		SearchType:  searchType,
	}

	c.template.RenderTemplate(w, r, "forum", data)
}

// CreateThread traite la création d'un nouveau sujet (thread) depuis le formulaire.
func (c *ForumControllers) CreateThread(w http.ResponseWriter, r *http.Request) {
	token := tokenDuCookie(r)
	if token == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Redirect(w, r, "/forum?erreur=Le titre ne peut pas être vide.", http.StatusSeeOther)
		return
	}

	categoryStr := r.FormValue("category")
	var categoriesID []int
	if catID, err := strconv.Atoi(categoryStr); err == nil && catID > 0 {
		categoriesID = append(categoriesID, catID)
	} else {
		categoriesID = append(categoriesID, 1) // catégorie "Général" par défaut
	}

	// 1. Créer le fil de discussion
	fil, err := c.service.CreateFil(token, title, categoriesID)
	if err != nil {
		http.Redirect(w, r, "/forum?erreur="+err.Error(), http.StatusSeeOther)
		return
	}

	// 2. Si un premier message est saisi, on le poste dans le fil
	messageContent := strings.TrimSpace(r.FormValue("message_topic"))
	if messageContent != "" {
		_ = c.service.CreateMessage(token, fil.ID, messageContent)
	}

	// Rediriger vers le fil de discussion nouvellement créé
	http.Redirect(w, r, "/threads/"+strconv.Itoa(fil.ID), http.StatusSeeOther)
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

// getConnecteUser retourne les informations de l'utilisateur connecté s'il existe.
func (c *ForumControllers) getConnecteUser(r *http.Request) *dto.Utilisateur {
	token := tokenDuCookie(r)
	if token == "" {
		return nil
	}
	user, err := c.service.GetMe(token)
	if err != nil {
		return nil
	}
	return user
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

	// 3. Avatar (optionnel) : si un fichier est envoyé, on l'upload via l'API et on récupère son URL.
	avatar := ""
	if file, header, errFile := r.FormFile("avatar"); errFile == nil {
		defer file.Close()
		data, errLecture := io.ReadAll(file)
		if errLecture != nil {
			c.template.RenderTemplate(w, r, "signup", SignupData{
				Erreur: "Erreur de lecture de l'avatar.", Noms: noms, Prenoms: prenoms, Email: email,
			})
			return
		}
		url, errUpload := c.service.UploadImage(header.Filename, data)
		if errUpload != nil {
			c.template.RenderTemplate(w, r, "signup", SignupData{
				Erreur: errUpload.Error(), Noms: noms, Prenoms: prenoms, Email: email,
			})
			return
		}
		avatar = url
	}

	// 4. La base n'a qu'un champ "username" : on fusionne prénom + nom pour le construire.
	username := prenoms + " " + noms

	// 5. Appeler l'API d'inscription (avec l'URL de l'avatar, vide si aucun fichier).
	if err := c.service.Register(username, email, password, avatar); err != nil {
		c.template.RenderTemplate(w, r, "signup", SignupData{
			Erreur:  err.Error(),
			Noms:    noms,
			Prenoms: prenoms,
			Email:   email,
		})
		return
	}

	// 6. Connexion automatique après inscription.
	token, err := c.service.Login(email, password)
	if err == nil && token != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   24 * 60 * 60, // 24 heures
		})
	}

	// Redirection après succès vers l'accueil.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ThreadData regroupe les données affichées sur la page d'un fil de discussion.
type ThreadData struct {
	Fil         *dto.FilDeDiscussion
	Messages    []dto.Message
	Connecte    bool
	Erreur      string
	Utilisateur *dto.Utilisateur
	Page        int
	Limit       int
	Sort        string
	TotalPages  int
	Pages       []PageLien
	HasPrev     bool
	PrevPage    int
	HasNext     bool
	NextPage    int
}

// renderThread récupère un fil + ses messages et affiche la page (avec un éventuel message d'erreur).
func (c *ForumControllers) renderThread(w http.ResponseWriter, r *http.Request, id int, page, limit int, sort string, erreur string) {
	fil, err := c.service.GetFil(id)
	if err != nil || fil == nil {
		http.Error(w, "Fil de discussion introuvable", http.StatusNotFound)
		return
	}

	var currentUserID int
	user := c.getConnecteUser(r)
	if user != nil {
		currentUserID = user.ID
	}

	// On récupère les messages ; en cas d'erreur on affiche le fil sans messages.
	res, err := c.service.GetMessages(id, page, limit, sort, currentUserID)
	var messages []dto.Message
	totalPages := 1
	if err == nil && res != nil {
		messages = res.Messages
		totalPages = res.TotalPages
	}

	var pages []PageLien
	for i := 1; i <= totalPages; i++ {
		pages = append(pages, PageLien{Num: i, Courante: i == page})
	}

	c.template.RenderTemplate(w, r, "thread", ThreadData{
		Fil:         fil,
		Messages:    messages,
		Connecte:    estConnecte(r),
		Erreur:      erreur,
		Utilisateur: user,
		Page:        page,
		Limit:       limit,
		Sort:        sort,
		TotalPages:  totalPages,
		Pages:       pages,
		HasPrev:     page > 1,
		PrevPage:    page - 1,
		HasNext:     page < totalPages,
		NextPage:    page + 1,
	})
}

// DisplayThread affiche un fil de discussion et ses messages.
func (c *ForumControllers) DisplayThread(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Identifiant du fil invalide", http.StatusBadRequest)
		return
	}

	page := 1
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	limit := 10
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "recent"
	}

	c.renderThread(w, r, id, page, limit, sort, r.URL.Query().Get("erreur"))
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
		c.renderThread(w, r, id, 1, 10, "recent", "Le message ne peut pas être vide.")
		return
	}

	// On envoie le token à l'API (Authorization: Bearer ...).
	if err := c.service.CreateMessage(token, id, contenu); err != nil {
		c.renderThread(w, r, id, 1, 10, "recent", err.Error())
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

// AdminData regroupe les données affichées sur le tableau de bord d'administration.
type AdminData struct {
	Utilisateur *dto.Utilisateur
	Users       []dto.Utilisateur
	Fils        []dto.FilDeDiscussion
	Erreur      string
}

// DisplayAdmin affiche le tableau de bord d'administration (réservé aux administrateurs).
func (c *ForumControllers) DisplayAdmin(w http.ResponseWriter, r *http.Request) {
	token := tokenDuCookie(r)
	if token == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := c.service.GetMe(token)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Seul un administrateur peut accéder au tableau de bord.
	if user.Role != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// On récupère les utilisateurs et les fils (tous états confondus) via l'API.
	users, _ := c.service.GetUsers(token)
	fils, _ := c.service.GetAllThreadsAdmin(token)

	c.template.RenderTemplate(w, r, "admin", AdminData{
		Utilisateur: user,
		Users:       users,
		Fils:        fils,
		Erreur:      r.URL.Query().Get("erreur"),
	})
}

// DeleteThread supprime un fil de discussion (créateur ou admin).
func (c *ForumControllers) DeleteThread(w http.ResponseWriter, r *http.Request) {
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
	// redirect : page de retour après l'action (ex. "/admin"), sinon le forum.
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/forum"
	}
	if err := c.service.DeleteThread(token, id); err != nil {
		http.Redirect(w, r, redirect+"?erreur="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

// EditThread modifie le titre d'un fil (créateur ou admin).
func (c *ForumControllers) EditThread(w http.ResponseWriter, r *http.Request) {
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
	titre := strings.TrimSpace(r.FormValue("title"))
	if titre == "" {
		http.Redirect(w, r, "/threads/"+strconv.Itoa(id)+"?erreur="+url.QueryEscape("Le titre ne peut pas être vide."), http.StatusSeeOther)
		return
	}
	if err := c.service.UpdateThread(token, id, titre); err != nil {
		http.Redirect(w, r, "/threads/"+strconv.Itoa(id)+"?erreur="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/threads/"+strconv.Itoa(id), http.StatusSeeOther)
}

// DeleteMessage supprime un message (auteur ou admin).
func (c *ForumControllers) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	token := tokenDuCookie(r)
	if token == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"]) // message id
	if err != nil {
		http.Error(w, "Identifiant du message invalide", http.StatusBadRequest)
		return
	}
	filIDStr := r.URL.Query().Get("fil_id")
	if err := c.service.DeleteMessage(token, id); err != nil {
		http.Redirect(w, r, "/threads/"+filIDStr+"?erreur="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/threads/"+filIDStr, http.StatusSeeOther)
}

// EditMessage modifie un message (auteur ou admin).
func (c *ForumControllers) EditMessage(w http.ResponseWriter, r *http.Request) {
	token := tokenDuCookie(r)
	if token == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"]) // message id
	if err != nil {
		http.Error(w, "Identifiant du message invalide", http.StatusBadRequest)
		return
	}
	contenu := strings.TrimSpace(r.FormValue("contenu"))
	filIDStr := r.URL.Query().Get("fil_id")
	if contenu == "" {
		http.Redirect(w, r, "/threads/"+filIDStr+"?erreur="+url.QueryEscape("Le message ne peut pas être vide."), http.StatusSeeOther)
		return
	}
	if err := c.service.UpdateMessage(token, id, contenu); err != nil {
		http.Redirect(w, r, "/threads/"+filIDStr+"?erreur="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/threads/"+filIDStr, http.StatusSeeOther)
}

// ReactToMessage gère la réaction AJAX (like/dislike).
func (c *ForumControllers) ReactToMessage(w http.ResponseWriter, r *http.Request) {
	token := tokenDuCookie(r)
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"erreur":"Vous devez être connecté pour réagir."}`))
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"]) // message id
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"erreur":"Identifiant du message invalide."}`))
		return
	}
	reactionType := r.URL.Query().Get("type")
	if reactionType != "like" && reactionType != "dislike" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"erreur":"Type de réaction invalide."}`))
		return
	}
	if err := c.service.ReactToMessage(token, id, reactionType); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"erreur":"` + err.Error() + `"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`))
}

// ChangeThreadState change l'état du fil (admin uniquement).
func (c *ForumControllers) ChangeThreadState(w http.ResponseWriter, r *http.Request) {
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
	state := r.FormValue("state")
	// redirect : page de retour après l'action (ex. "/admin"), sinon le fil concerné.
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/threads/" + strconv.Itoa(id)
	}
	if err := c.service.ChangeThreadState(token, id, state); err != nil {
		http.Redirect(w, r, redirect+"?erreur="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}

// BanUser bannit un utilisateur (admin uniquement).
func (c *ForumControllers) BanUser(w http.ResponseWriter, r *http.Request) {
	token := tokenDuCookie(r)
	if token == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Identifiant de l'utilisateur invalide", http.StatusBadRequest)
		return
	}
	// redirect : page de retour après l'action (ex. "/admin"), sinon le fil concerné.
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = "/threads/" + r.URL.Query().Get("fil_id")
	}
	if err := c.service.BanUser(token, userID); err != nil {
		http.Redirect(w, r, redirect+"?erreur="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
