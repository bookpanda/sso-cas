package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Base
	Token     string    `json:"token" gorm:"tinytext"`
	UserID    uuid.UUID `json:"user_id" gorm:"tinytext"`
	Payload   string    `json:"payload" gorm:"text"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:timestamp"`
}
