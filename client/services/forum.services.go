// Package services contient la logique métier côté client avant appel à l'API.
package services

import (
	"projet-forum/client/api"
	"projet-forum/client/dto"
)

// ForumService encapsule le client API pour les opérations du forum.
type ForumService struct {
	forumApi *api.ForumApi
}

// InitForumService initialise le service forum avec le client API.
func InitForumService(forumApi *api.ForumApi) *ForumService {
	return &ForumService{forumApi: forumApi}
}

// GetFils retourne la liste paginée des fils depuis l'API.
func (s *ForumService) GetFils(page int, limit int) (*dto.ThreadResponse, error) {
	return s.forumApi.GetFils(page, limit)
}

// GetLastActus retourne les actualités depuis l'API.
func (s *ForumService) GetLastActus() ([]dto.Actualite, error) {
	return s.forumApi.GetLastActus()
}

// GetFil retourne un fil de discussion par son identifiant.
func (s *ForumService) GetFil(id int) (*dto.FilDeDiscussion, error) {
	return s.forumApi.GetFil(id)
}

// GetMessages retourne les messages d'un fil de discussion.
func (s *ForumService) GetMessages(id int) ([]dto.Message, error) {
	return s.forumApi.GetMessages(id)
}

// Register transmet la demande d'inscription à l'API.
func (s *ForumService) Register(username, email, password string) error {
	return s.forumApi.Register(username, email, password)
}
