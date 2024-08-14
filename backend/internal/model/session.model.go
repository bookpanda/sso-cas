package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Base
	Token     string    `json:"token" gorm:"tinytext"`
	UserID    uuid.UUID `json:"user_id" gorm:"tinytext"`
	Email     string    `json:"email" gorm:"tinytext"`
	Role      string    `json:"role" gorm:"tinytext"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:timestamp"`
}
