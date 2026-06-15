package services

import (
	"projet-forum/src/client/api"
	"projet-forum/src/models"
)

type ForumClientService struct {
	api *api.ForumApi
}

func NewForumClientService(api *api.ForumApi) *ForumClientService {
	return &ForumClientService{api: api}
}

func (s *ForumClientService) GetFils(page int, limit int) (*api.ThreadResponse, error) {
	return s.api.GetFils(page, limit)
}

func (s *ForumClientService) GetLastActus(limit int) ([]models.FilDeDiscussion, error) {
	return s.api.GetLastActus(limit)
}
