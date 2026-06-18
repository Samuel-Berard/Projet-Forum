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
			Titre:   "Lancement de la Saison 2 de Crown of Embers",
			Contenu: "La mise à jour majeure apporte un nouveau biome volcanique, des boss redoutables et un rééquilibrage de la magie de feu. Donnez vos avis sur le forum !",
		},
		{
			ID:      2,
			Titre:   "Guide ultime : Optimiser son build dans Elden Ring DLC",
			Contenu: "Vous bloquez sur les boss du royaume des ombres ? Notre guide collaboratif regroupe les meilleures configurations d'armes et de talismans.",
		},
		{
			ID:      3,
			Titre:   "Nintendo Switch 2 : Ce que l'on sait officiellement",
			Contenu: "Rétrocompatibilité, puissance graphique en hausse, écran OLED de 8 pouces... Récapitulatif complet des annonces officielles et des rumeurs crédibles.",
		},
		{
			ID:      4,
			Titre:   "Le patch 1.1 de Velocity GP est disponible",
			Contenu: "Correction des bugs physiques lors des collisions, ajout de 3 nouveaux circuits urbains et amélioration générale du matchmaking classé.",
		},
		{
			ID:      5,
			Titre:   "Pixel Tactics : Le jeu de stratégie indé à faire absolument en 2026",
			Contenu: "Découvrez notre test de cette pépite tactique au tour par tour en pixel art, disponible dès maintenant sur PC et consoles.",
		},
		{
			ID:      6,
			Titre:   "GTA 6 : Une nouvelle bande-annonce analysée en détail",
			Contenu: "Les moindres détails du dernier trailer décryptés : cycle jour/nuit dynamique, système de gestion des foules et détails physiques des véhicules.",
		},
	}
}
