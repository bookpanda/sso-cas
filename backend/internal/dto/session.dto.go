package dto

type Session struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
}

type CreateSessionRequest struct {
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
}

type DeleteByEmailSessionRequest struct {
	Email string `json:"email"`
}
