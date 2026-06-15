// Package dto contient les structures de données échangées avec l'API et utilisées dans les templates.
package dto

import "time"

type Categorie struct {
	ID  int    `json:"id"`
	Nom string `json:"nom"`
}

type FilDeDiscussion struct {
	ID         int          `json:"id"`
	Titre      string       `json:"titre"`
	Etat       string       `json:"etat"`
	AuteurID   int          `json:"auteur_id"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	Auteur     *Utilisateur `json:"auteur,omitempty"`
	Categories []Categorie  `json:"categories,omitempty"`
}

// ThreadResponse représente la réponse paginée renvoyée par l'API pour les fils.
type ThreadResponse struct {
	Fils       []FilDeDiscussion `json:"fils"`
	TotalPages int               `json:"totalPages"`
	Page       int               `json:"page"`
}
