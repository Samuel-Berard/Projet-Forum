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

// Register transmet la demande d'inscription à l'API (avatar = URL de l'image, vide si aucun).
func (s *ForumService) Register(username, email, password, avatar string) error {
	return s.forumApi.Register(username, email, password, avatar)
}

// Login transmet les identifiants à l'API et retourne le token.
func (s *ForumService) Login(identifiant, password string) (string, error) {
	return s.forumApi.Login(identifiant, password)
}

// CreateMessage publie un message dans un fil via l'API (nécessite le token).
func (s *ForumService) CreateMessage(token string, filID int, contenu string) error {
	return s.forumApi.CreateMessage(token, filID, contenu)
}

// UploadImage envoie une image à l'API et retourne son URL publique.
func (s *ForumService) UploadImage(filename string, data []byte) (string, error) {
	return s.forumApi.UploadFile(filename, data)
}

// GetMe retourne les informations de l'utilisateur connecté.
func (s *ForumService) GetMe(token string) (*dto.Utilisateur, error) {
	return s.forumApi.GetMe(token)
}

// UpdateAvatar enregistre l'URL de l'avatar de l'utilisateur connecté.
func (s *ForumService) UpdateAvatar(token, avatarURL string) error {
	return s.forumApi.UpdateAvatar(token, avatarURL)
}
