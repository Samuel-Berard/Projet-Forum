// Package api fournit un client HTTP pour appeler l'API du forum.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"projet-forum/client/dto"
)

// ForumApi contient l'URL de base utilisée pour les appels à l'API.
type ForumApi struct {
	baseURL string
}

// InitForumApi initialise le client API avec l'URL de base fournie.
func InitForumApi(baseURL string) *ForumApi {
	return &ForumApi{baseURL: baseURL}
}

// executeRequest envoie la requête HTTP et décode la réponse JSON dans result.
func (api *ForumApi) executeRequest(req *http.Request, result interface{}) (int, error) {
	_client := http.Client{
		Timeout: time.Second * 5,
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := _client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return resp.StatusCode, readErr
	}

	if resp.StatusCode >= 400 {
		return resp.StatusCode, fmt.Errorf("Erreur API - %s", string(bytes.TrimSpace(body)))
	}

	if result != nil {
		if err = json.Unmarshal(body, result); err != nil {
			return resp.StatusCode, fmt.Errorf("Erreur décodage JSON - %s", err.Error())
		}
	}

	return resp.StatusCode, nil
}

// GetFils récupère la liste paginée des fils de discussion depuis l'API.
func (api *ForumApi) GetFils(page int, limit int, search string, searchType string, categoryID int) (*dto.ThreadResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/threads?page=%d&limit=%d&search=%s&type=%s&category=%d", api.baseURL, page, limit, url.QueryEscape(search), url.QueryEscape(searchType), categoryID), nil)
	if err != nil {
		return nil, err
	}

	var response dto.ThreadResponse
	_, err = api.executeRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetLastActus récupère les actualités depuis l'API.
func (api *ForumApi) GetLastActus() ([]dto.Actualite, error) {
	req, err := http.NewRequest(http.MethodGet, api.baseURL+"/actus", nil)
	if err != nil {
		return nil, err
	}

	var actus []dto.Actualite
	_, err = api.executeRequest(req, &actus)
	if err != nil {
		return nil, err
	}

	return actus, nil
}

// GetFil récupère un fil de discussion par son identifiant.
func (api *ForumApi) GetFil(id int) (*dto.FilDeDiscussion, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/threads/%d", api.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	var fil dto.FilDeDiscussion
	_, err = api.executeRequest(req, &fil)
	if err != nil {
		return nil, err
	}

	return &fil, nil
}

// GetMessages récupère les messages d'un fil de discussion avec pagination, tri et réactions.
func (api *ForumApi) GetMessages(filID int, page, limit int, sort string, currentUserID int) (*dto.MessageResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/threads/%d/messages?page=%d&limit=%d&sort=%s&current_user_id=%d", api.baseURL, filID, page, limit, url.QueryEscape(sort), currentUserID), nil)
	if err != nil {
		return nil, err
	}

	var response dto.MessageResponse
	_, err = api.executeRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteThread supprime un fil de discussion.
func (api *ForumApi) DeleteThread(token string, id int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/threads/%d", api.baseURL, id), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	_, err = api.executeRequest(req, nil)
	return err
}

// UpdateThread modifie le titre d'un fil.
func (api *ForumApi) UpdateThread(token string, id int, titre string) error {
	corps, err := json.Marshal(map[string]string{"titre": titre})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/threads/%d", api.baseURL, id), bytes.NewBuffer(corps))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	_, err = api.executeRequest(req, nil)
	return err
}

// DeleteMessage supprime un message.
func (api *ForumApi) DeleteMessage(token string, id int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/messages/%d", api.baseURL, id), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	_, err = api.executeRequest(req, nil)
	return err
}

// UpdateMessage modifie le contenu d'un message.
func (api *ForumApi) UpdateMessage(token string, id int, contenu string) error {
	corps, err := json.Marshal(map[string]string{"contenu": contenu})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/messages/%d", api.baseURL, id), bytes.NewBuffer(corps))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	_, err = api.executeRequest(req, nil)
	return err
}

// ReactToMessage ajoute ou modifie une réaction sur un message.
func (api *ForumApi) ReactToMessage(token string, messageID int, reactionType string) error {
	corps, err := json.Marshal(map[string]string{"type": reactionType})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/messages/%d/react", api.baseURL, messageID), bytes.NewBuffer(corps))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	_, err = api.executeRequest(req, nil)
	return err
}

// ChangeThreadState change l'état d'un fil (réservé aux admins).
func (api *ForumApi) ChangeThreadState(token string, id int, state string) error {
	corps, err := json.Marshal(map[string]string{"etat": state})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/threads/%d/state", api.baseURL, id), bytes.NewBuffer(corps))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	_, err = api.executeRequest(req, nil)
	return err
}

// BanUser bannit un utilisateur (réservé aux admins).
func (api *ForumApi) BanUser(token string, userID int) error {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/users/%d/ban", api.baseURL, userID), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	_, err = api.executeRequest(req, nil)
	return err
}

// Register envoie une demande d'inscription (POST) à l'API.
func (api *ForumApi) Register(username, email, password, avatar string) error {
	corps, err := json.Marshal(map[string]string{
		"username": username,
		"email":    email,
		"password": password,
		"avatar":   avatar,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, api.baseURL+"/register", bytes.NewBuffer(corps))
	if err != nil {
		return err
	}

	// executeRequest renvoie une erreur (avec le message de l'API) si le statut est >= 400.
	_, err = api.executeRequest(req, nil)
	return err
}

// Login envoie les identifiants (POST) à l'API et retourne le token JWT en cas de succès.
func (api *ForumApi) Login(identifiant, password string) (string, error) {
	corps, err := json.Marshal(map[string]string{
		"identifiant": identifiant,
		"password":    password,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, api.baseURL+"/login", bytes.NewBuffer(corps))
	if err != nil {
		return "", err
	}

	// L'API répond : {"message": "...", "token": "..."} — on ne garde que le token.
	var reponse struct {
		Token string `json:"token"`
	}
	if _, err := api.executeRequest(req, &reponse); err != nil {
		return "", err
	}

	return reponse.Token, nil
}

// CreateMessage publie un message dans un fil (route protégée → token en Bearer).
func (api *ForumApi) CreateMessage(token string, filID int, contenu string) error {
	corps, err := json.Marshal(map[string]string{"contenu": contenu})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/threads/%d/messages", api.baseURL, filID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(corps))
	if err != nil {
		return err
	}

	// Le token JWT est envoyé dans l'en-tête Authorization au format Bearer (exigé par l'API).
	req.Header.Set("Authorization", "Bearer "+token)

	_, err = api.executeRequest(req, nil)
	return err
}

// UploadFile envoie un fichier à l'API en multipart/form-data (champ "file")
// et retourne l'URL publique de l'image.
func (api *ForumApi) UploadFile(filename string, data []byte) (string, error) {
	// Construction du corps multipart, comme un formulaire de fichier.
	var corps bytes.Buffer
	writer := multipart.NewWriter(&corps)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(data); err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, api.baseURL+"/upload", &corps)
	if err != nil {
		return "", err
	}
	// Indispensable : le Content-Type multipart avec sa frontière (boundary).
	// On ne passe donc pas par executeRequest (qui force application/json).
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("Erreur API - %s", string(bytes.TrimSpace(body)))
	}

	// L'API répond : {"url": "http://..."}.
	var reponse struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &reponse); err != nil {
		return "", fmt.Errorf("Erreur décodage JSON - %s", err.Error())
	}

	return reponse.URL, nil
}

// GetMe récupère les informations de l'utilisateur connecté (route protégée → token Bearer).
func (api *ForumApi) GetMe(token string) (*dto.Utilisateur, error) {
	req, err := http.NewRequest(http.MethodGet, api.baseURL+"/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	var user dto.Utilisateur
	if _, err := api.executeRequest(req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateAvatar enregistre l'URL de l'avatar pour l'utilisateur connecté (route protégée).
func (api *ForumApi) UpdateAvatar(token, avatarURL string) error {
	corps, err := json.Marshal(map[string]string{"avatar": avatarURL})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, api.baseURL+"/me/avatar", bytes.NewBuffer(corps))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	_, err = api.executeRequest(req, nil)
	return err
}

// CreateFil crée un nouveau fil de discussion (route protégée -> token Bearer).
func (api *ForumApi) CreateFil(token string, titre string, categoriesID []int) (*dto.FilDeDiscussion, error) {
	corps, err := json.Marshal(map[string]interface{}{
		"titre":         titre,
		"categories_id": categoriesID,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, api.baseURL+"/threads", bytes.NewBuffer(corps))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	var reponse struct {
		Message string              `json:"message"`
		Fil     dto.FilDeDiscussion `json:"fil"`
	}

	if _, err := api.executeRequest(req, &reponse); err != nil {
		return nil, err
	}

	return &reponse.Fil, nil
}

