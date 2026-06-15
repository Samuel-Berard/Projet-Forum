package services

import (
	"errors"
	"regexp"
	"strconv"

	"projet-forum/api/auth"
	"projet-forum/api/models"
	"projet-forum/api/repositories"
	"projet-forum/api/utils"
)

type UtilisateurService struct {
	repo *repositories.UtilisateurRepository
}

func InitUtilisateurService(repo *repositories.UtilisateurRepository) *UtilisateurService {
	return &UtilisateurService{repo: repo}
}

// Register valide les donnees, hache le mot de passe et cree un nouvel utilisateur.
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

// Login verifie les identifiants et retourne un token JWT si la connexion est valide.
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

// isValidPassword controle la robustesse du mot de passe.
func isValidPassword(password string) bool {
	if len(password) < 12 {
		return false
	}
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(password)
	return hasUppercase && hasSpecial
}

// BanUser bannit un utilisateur, reserve aux administrateurs.
func (s *UtilisateurService) BanUser(id int, role string) error {
	if role != "admin" {
		return errors.New("accès refusé, rôle administrateur requis")
	}

	return s.repo.BanUser(id)
}
