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
			Titre:   "GTA 6 : Rockstar confirme le calendrier de sortie pour l'automne 2025",
			Contenu: "Take-Two a réaffirmé sa fenêtre de sortie pour le très attendu Grand Theft Auto VI lors de son dernier bilan financier.",
		},
		{
			ID:      2,
			Titre:   "Nintendo Switch 2 : Ce que l'on sait sur les caractéristiques de la console",
			Contenu: "Rumeurs insistantes sur l'écran, la rétrocompatibilité et la puissance de la future console hybride.",
		},
		{
			ID:      3,
			Titre:   "Hollow Knight Silksong : Nouvelle fuite d'une classification par âge",
			Contenu: "L'attente touche-t-elle à sa fin ? Le jeu a reçu une classification officielle, relançant les spéculations.",
		},
		{
			ID:      4,
			Titre:   "Elden Ring : Un nouveau record de ventes franchi à travers le monde",
			Contenu: "FromSoftware annonce fièrement plus de 25 millions d'exemplaires vendus avant la sortie du DLC.",
		},
	}
}
