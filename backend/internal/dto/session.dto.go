package dto

import (
	"time"
)

type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CreateSessionRequest struct {
	UserID string `json:"user_id"`
}

type DeleteByEmailSessionRequest struct {
	Email string `json:"email"`
}
