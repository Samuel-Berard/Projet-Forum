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
