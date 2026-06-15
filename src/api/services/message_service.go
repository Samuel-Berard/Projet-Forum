package services

import (
	"errors"
	"projet-forum/src/models"
	"projet-forum/src/api/repositories"
)

type MessageService struct {
	repo    *repositories.MessageRepository
	filRepo *repositories.FilRepository
}

func NewMessageService(repo *repositories.MessageRepository, filRepo *repositories.FilRepository) *MessageService {
	return &MessageService{repo: repo, filRepo: filRepo}
}

func (s *MessageService) CreateMessage(contenu string, filID, auteurID int) (*models.Message, error) {
	if contenu == "" {
		return nil, errors.New("le contenu du message ne peut pas être vide")
	}

	fil, err := s.filRepo.FindByID(filID)
	if err != nil {
		return nil, err
	}

	if fil.Etat != "ouvert" {
		return nil, errors.New("impossible de publier, ce fil n'est plus ouvert")
	}

	msg := &models.Message{
		Contenu:  contenu,
		FilID:    filID,
		AuteurID: auteurID,
	}

	err = s.repo.Create(msg)
	if err != nil {
		return nil, errors.New("erreur lors de la création du message")
	}

	return msg, nil
}

func (s *MessageService) GetMessagesByFil(filID, page, limit int, sortBy string) ([]models.Message, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	return s.repo.FindByFilID(filID, page, limit, sortBy)
}

func (s *MessageService) UpdateMessage(id int, contenu string, userID int) error {
	msg, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if msg.AuteurID != userID {
		return errors.New("vous n'êtes pas autorisé à modifier ce message")
	}

	msg.Contenu = contenu
	return s.repo.Update(msg)
}

func (s *MessageService) DeleteMessage(id int, userID int, role string) error {
	msg, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if msg.AuteurID != userID && role != "admin" {
		return errors.New("vous n'êtes pas autorisé à supprimer ce message")
	}

	return s.repo.Delete(id)
}

func (s *MessageService) ReactToMessage(userID, messageID int, typeReaction string) error {
	if typeReaction != "like" && typeReaction != "dislike" {
		return errors.New("type de réaction invalide")
	}

	return s.repo.React(userID, messageID, typeReaction)
}
