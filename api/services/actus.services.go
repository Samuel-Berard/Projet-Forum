package services

import (
	"projet-forum/api/models"
)

// ActusService fournit les actualités du forum.
// Les actualités sont codées en dur (pas de table en base de données).
type ActusService struct{}

func InitActusService() *ActusService {
	return &ActusService{}
}

// GetLastActus retourne la liste des actualités codées en dur.
func (s *ActusService) GetLastActus() []models.Actualite {
	return []models.Actualite{
		{
			ID:      1,
			Titre:   "Bienvenue sur le forum !",
			Contenu: "Découvrez les fils de discussion et participez aux échanges de la communauté.",
		},
		{
			ID:      2,
			Titre:   "Les règles de la communauté",
			Contenu: "Restez courtois et respectez les autres membres pour garder un espace agréable.",
		},
		{
			ID:      3,
			Titre:   "Nouveautés à venir",
			Contenu: "De nouvelles fonctionnalités arrivent bientôt sur la plateforme. Restez connectés !",
		},
		{
			ID:      4,
			Titre:   "Astuce du jour",
			Contenu: "Utilisez la recherche pour retrouver rapidement un sujet qui vous intéresse.",
		},
	}
}
