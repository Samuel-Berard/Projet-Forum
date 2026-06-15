package dto

import "time"

type Message struct {
	ID              int          `json:"id"`
	Contenu         string       `json:"contenu"`
	FilID           int          `json:"fil_id"`
	AuteurID        int          `json:"auteur_id"`
	ScorePopularite int          `json:"score_popularite"`
	CreatedAt       time.Time    `json:"created_at"`
	Auteur          *Utilisateur `json:"auteur,omitempty"`
}
