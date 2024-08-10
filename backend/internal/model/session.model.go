package model

type Session struct {
	Base
	UserID    string `json:"user_id" gorm:"tinytext"`
	ExpiresAt int64  `json:"expires_at" gorm:"type:timestamp"`
}
