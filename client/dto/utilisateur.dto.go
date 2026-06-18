package dto

import "time"

type Utilisateur struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Banned    bool      `json:"banned"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
}
