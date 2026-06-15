package services

import (
	"errors"
	"projet-forum/src/api/auth"
	"projet-forum/src/models"
	"projet-forum/src/api/repositories"
	"projet-forum/src/api/utils"
	"regexp"
	"strconv"
)

type UtilisateurService struct {
	repo *repositories.UtilisateurRepository
}

func NewUtilisateurService(repo *repositories.UtilisateurRepository) *UtilisateurService {
	return &UtilisateurService{repo: repo}
}

func (s *UtilisateurService) Register(req *models.RegisterRequest) (*models.Utilisateur, error) {
	
	if !isValidPassword(req.Password) {
		return nil, errors.New("le mot de passe doit contenir au moins 12 caractères, une majuscule et un caractère spécial")
	}


	hashedPassword := utils.HashPasswordSHA512(req.Password)

	user := &models.Utilisateur{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "user",
	}

	err := s.repo.Create(user)
	if err != nil {
		return nil, errors.New("erreur lors de l'inscription, nom d'utilisateur ou email peut-être déjà utilisé")
	}

	return user, nil
}

func (s *UtilisateurService) Login(req *models.LoginRequest) (string, error) {
	user, err := s.repo.FindByUsernameOrEmail(req.Identifiant)
	if err != nil {
		return "", errors.New("identifiants incorrects")
	}

	if user.Banned {
		return "", errors.New("ce compte a été banni")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", errors.New("identifiants incorrects")
	}

	token, err := auth.GenerateToken(strconv.Itoa(user.ID), user.Role)
	if err != nil {
		return "", errors.New("erreur lors de la génération du token")
	}

	return token, nil
}

func isValidPassword(password string) bool {
	if len(password) < 12 {
		return false
	}
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(password)
	return hasUppercase && hasSpecial
}

func (s *UtilisateurService) BanUser(id int, role string) error {
	if role != "admin" {
		return errors.New("accès refusé, rôle administrateur requis")
	}

	return s.repo.BanUser(id)
}
