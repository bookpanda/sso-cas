package dto

import "time"

type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	Payload   string    `json:"payload"`
	ExpiresAt time.Time `json:"expires_at"`
}

type SessionPayload struct {
}

type CreateSessionRequest struct {
	UserID string `json:"user_id"`
}

type DeleteByEmailSessionRequest struct {
	Email string `json:"email"`
}
