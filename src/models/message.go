package models

import "time"

type Message struct {
	ID              int          `json:"id"`
	Contenu         string       `json:"contenu"`
	FilID           int          `json:"fil_id"`
	AuteurID        int          `json:"auteur_id"`
	ScorePopularite int          `json:"score_popularite"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	Auteur          *Utilisateur `json:"auteur,omitempty"`
}

type CreateMessageRequest struct {
	Contenu string `json:"contenu"`
}

type UpdateMessageRequest struct {
	Contenu string `json:"contenu"`
}

type ReactionRequest struct {
	Type string `json:"type"` // "like" ou "dislike"
}

type Reaction struct {
	UtilisateurID int    `json:"utilisateur_id"`
	MessageID     int    `json:"message_id"`
	Type          string `json:"type"`
}
