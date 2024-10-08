package model

import (
	"time"

	"github.com/google/uuid"
)

type ServiceTicket struct {
	Base
	SessionToken string    `json:"session_token" gorm:"tinytext"`
	Token        string    `json:"token" gorm:"tinytext"`
	ServiceUrl   string    `json:"service_url" gorm:"tinytext"`
	UserID       uuid.UUID `json:"user_id" gorm:"tinytext"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"type:timestamp"`
}
