package services

import (
	"errors"
	"projet-forum/src/models"
	"projet-forum/src/repositories"
)

type FilService struct {
	repo *repositories.FilRepository
}

func NewFilService(repo *repositories.FilRepository) *FilService {
	return &FilService{repo: repo}
}

func (s *FilService) CreateFil(titre string, auteurID int, categoryIDs []int) (*models.FilDeDiscussion, error) {
	if titre == "" {
		return nil, errors.New("le titre ne peut pas être vide")
	}

	fil := &models.FilDeDiscussion{
		Titre:    titre,
		Etat:     "ouvert",
		AuteurID: auteurID,
	}

	err := s.repo.Create(fil, categoryIDs)
	if err != nil {
		return nil, errors.New("erreur lors de la création du fil de discussion")
	}

	return fil, nil
}

func (s *FilService) GetFils(page, limit int, search string, categoryID int) ([]models.FilDeDiscussion, error) {
	// Defaults for pagination
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // Default batch size according to FT-9
	}

	return s.repo.FindAll(page, limit, search, categoryID)
}

func (s *FilService) GetFilByID(id int) (*models.FilDeDiscussion, error) {
	fil, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Un fil archivé ne doit pas être visible, mais s'il est accédé par l'admin, peut-être,
	// pour l'instant la règle FT-4 dit : Les fils archivés ne doivent pas être visibles.
	if fil.Etat == "archive" {
		return nil, errors.New("ce fil a été archivé")
	}

	return fil, nil
}

func (s *FilService) UpdateFil(id int, titre string, auteurID int, role string) error {
	fil, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Seul le propriétaire peut modifier le fil (FT-7)
	// (Admin pourrait le supprimer ou changer l'état mais FT-7 ne dit pas explicitement que l'admin peut modifier le titre)
	if fil.AuteurID != auteurID && role != "admin" {
		return errors.New("vous n'êtes pas autorisé à modifier ce fil")
	}

	fil.Titre = titre
	return s.repo.Update(fil)
}

func (s *FilService) DeleteFil(id int, userID int, role string) error {
	fil, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Admin ou propriétaire peuvent supprimer (FT-7)
	if fil.AuteurID != userID && role != "admin" {
		return errors.New("vous n'êtes pas autorisé à supprimer ce fil")
	}

	return s.repo.Delete(id)
}

func (s *FilService) ChangeFilState(id int, etat string, role string) error {
	// Seul l'admin peut modifier l'état (FT-12)
	if role != "admin" {
		return errors.New("accès refusé, rôle administrateur requis")
	}

	fil, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if etat != "ouvert" && etat != "ferme" && etat != "archive" {
		return errors.New("état invalide")
	}

	fil.Etat = etat
	return s.repo.Update(fil)
}
