package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"projet-forum/src/models"
	"time"
)

type ForumApi struct {
	baseURL string
}

func NewForumApi(baseURL string) *ForumApi {
	return &ForumApi{baseURL: baseURL}
}

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

// ThreadResponse représente la réponse paginée de l'API
type ThreadResponse struct {
	Fils       []models.FilDeDiscussion `json:"fils"`
	TotalPages int                      `json:"totalPages"`
	Page       int                      `json:"page"`
}

// GetFils fetch threads from the API
func (api *ForumApi) GetFils(page int, limit int) (*ThreadResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/threads?page=%d&limit=%d", api.baseURL, page, limit), nil)
	if err != nil {
		return nil, err
	}

	var response ThreadResponse
	_, err = api.executeRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetLastActus fetch recent threads or news
func (api *ForumApi) GetLastActus(limit int) ([]models.FilDeDiscussion, error) {
	// The API doesn't have a distinct Actus endpoint, it just returns fils.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/threads?limit=%d", api.baseURL, limit), nil)
	if err != nil {
		return nil, err
	}

	var response ThreadResponse
	_, err = api.executeRequest(req, &response)
	if err != nil {
		return nil, err
	}

	return response.Fils, nil
}
