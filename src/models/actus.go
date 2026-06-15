package models

import "time"

type Actualite struct {
	ID        int          `json:"id"`
	Titre     string       `json:"titre"`
	Contenu   string       `json:"contenu"`
	AuteurID  int          `json:"auteur_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Auteur    *Utilisateur `json:"auteur,omitempty"`
}
