package services

import (
	"errors"
	"projet-forum/src/models"
	"projet-forum/src/repositories"
)

type ActualiteService struct {
	repo *repositories.ActualiteRepository
}

func NewActualiteService(repo *repositories.ActualiteRepository) *ActualiteService {
	return &ActualiteService{repo: repo}
}

func (s *ActualiteService) GetLastActus(limit int) ([]models.Actualite, error) {
	return s.repo.FindLast(limit)
}

func (s *ActualiteService) CreateActualite(titre, contenu string, auteurID int) (*models.Actualite, error) {
	actualite := &models.Actualite{
		Titre:    titre,
		Contenu:  contenu,
		AuteurID: auteurID,
	}

	err := s.repo.Create(actualite)
	if err != nil {
		return nil, err
	}

	return actualite, nil
}

func (s *ActualiteService) UpdateActualite(id int, titre, contenu string, auteurID int) error {
	actualite, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if actualite.AuteurID != auteurID {
		return errors.New("vous n'êtes pas autorisé à modifier cette actualité")
	}

	actualite.Titre = titre
	actualite.Contenu = contenu
	return s.repo.Update(actualite)
}

func (s *ActualiteService) DeleteActualite(id int, auteurID int) error {
	actualite, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if actualite.AuteurID != auteurID {
		return errors.New("vous n'êtes pas autorisé à supprimer cette actualité")
	}

	return s.repo.Delete(id)
}

func (s *ActualiteService) GetActualiteByID(id int) (*models.Actualite, error) {
	return s.repo.FindByID(id)
}
func (s *ActualiteService) GetActualites(page, limit int) ([]models.Actualite, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	return s.repo.FindAll(page, limit)
}

func (s *ActualiteService) GetTotalPages(limit int) (int, error) {
	if limit <= 0 {
		limit = 10
	}

	totalActus, err := s.repo.CountAll()
	if err != nil {
		return 0, err
	}

	totalPages := totalActus / limit
	if totalActus%limit > 0 {
		totalPages++
	}

	if totalPages == 0 {
		totalPages = 1
	}
	return totalPages, nil
}
