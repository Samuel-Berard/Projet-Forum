package models

import "time"

type Utilisateur struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	Banned       bool      `json:"banned"`
	Avatar       string    `json:"avatar"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

type LoginRequest struct {
	Identifiant string `json:"identifiant"`
	Password    string `json:"password"`
}
