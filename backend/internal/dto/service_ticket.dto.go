package dto

import "time"

type ServiceTicket struct {
	SessionToken string    `json:"session_token"`
	Token        string    `json:"token"`
	ServiceUrl   string    `json:"service_url"`
	UserID       string    `json:"user_id"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type CreateServiceTicketRequest struct {
	SessionToken string `json:"session_token"`
	ServiceUrl   string `json:"service_url"`
	UserID       string `json:"user_id"`
}
