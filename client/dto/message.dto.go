package dto

import "time"

type Message struct {
	ID              int          `json:"id"`
	Contenu         string       `json:"contenu"`
	FilID           int          `json:"fil_id"`
	AuteurID        int          `json:"auteur_id"`
	ScorePopularite int          `json:"score_popularite"`
	UserReaction    string       `json:"user_reaction"`
	CreatedAt       time.Time    `json:"created_at"`
	Auteur          *Utilisateur `json:"auteur,omitempty"`
}

type MessageResponse struct {
	Messages   []Message `json:"messages"`
	TotalPages int       `json:"totalPages"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
}
