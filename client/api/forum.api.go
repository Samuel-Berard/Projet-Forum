// Package api fournit un client HTTP pour appeler l'API du forum.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
func (api *ForumApi) GetFils(page int, limit int) (*dto.ThreadResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/threads?page=%d&limit=%d", api.baseURL, page, limit), nil)
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

// GetMessages récupère les messages d'un fil de discussion.
func (api *ForumApi) GetMessages(filID int) ([]dto.Message, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/threads/%d/messages", api.baseURL, filID), nil)
	if err != nil {
		return nil, err
	}

	var messages []dto.Message
	_, err = api.executeRequest(req, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// Register envoie une demande d'inscription (POST) à l'API.
func (api *ForumApi) Register(username, email, password string) error {
	corps, err := json.Marshal(map[string]string{
		"username": username,
		"email":    email,
		"password": password,
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
