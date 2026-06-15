package services

import (
	"errors"
	"projet-forum/src/models"
	"projet-forum/src/api/repositories"
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
	
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 
	}

	return s.repo.FindAll(page, limit, search, categoryID)
}

func (s *FilService) GetTotalPages(limit int, search string, categoryID int) (int, error) {
	if limit <= 0 {
		limit = 10
	}
	totalFils, err := s.repo.CountAll(search, categoryID)
	if err != nil {
		return 0, err
	}
	
	totalPages := totalFils / limit
	if totalFils % limit > 0 {
		totalPages++
	}
	
	if totalPages == 0 {
		totalPages = 1
	}
	return totalPages, nil
}

func (s *FilService) GetFilByID(id int) (*models.FilDeDiscussion, error) {
	fil, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	
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

	
	if fil.AuteurID != userID && role != "admin" {
		return errors.New("vous n'êtes pas autorisé à supprimer ce fil")
	}

	return s.repo.Delete(id)
}

func (s *FilService) ChangeFilState(id int, etat string, role string) error {
	
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
