package model

import "github.com/google/uuid"

type Session struct {
	Base
	UserID    uuid.UUID `json:"user_id" gorm:"tinytext"`
	ExpiresAt int64     `json:"expires_at" gorm:"type:timestamp"`
}
