package dto

import "time"

type ServiceTicket struct {
	Token      string    `json:"token"`
	ServiceUrl string    `json:"service_url"`
	UserID     string    `json:"user_id"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type CreateServiceTicketRequest struct {
	ServiceUrl string `json:"service_url"`
	UserID     string `json:"user_id"`
}
