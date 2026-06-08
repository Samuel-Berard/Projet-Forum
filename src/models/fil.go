package models

import "time"

type Categorie struct {
	ID  int    `json:"id"`
	Nom string `json:"nom"`
}

type FilDeDiscussion struct {
	ID        int       `json:"id"`
	Titre     string    `json:"titre"`
	Etat      string    `json:"etat"` // ouvert, ferme, archive
	AuteurID  int       `json:"auteur_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Auteur    *Utilisateur `json:"auteur,omitempty"`
	Categories []Categorie `json:"categories,omitempty"`
}

type CreateFilRequest struct {
	Titre        string `json:"titre"`
	CategoriesID []int  `json:"categories_id"`
}

type UpdateFilRequest struct {
	Titre string `json:"titre"`
}
