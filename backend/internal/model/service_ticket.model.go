package model

import "github.com/google/uuid"

type ServiceTicket struct {
	Base
	Token      string    `json:"token" gorm:"tinytext"`
	ServiceUrl string    `json:"service_url" gorm:"tinytext"`
	UserID     uuid.UUID `json:"user_id" gorm:"tinytext"`
	ExpiresAt  int       `json:"expires_at" gorm:"type:timestamp"`
}
