package session

import (
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/model"
)

func ModelToDto(in *model.Session) *dto.Session {
	return &dto.Session{
		ID:        in.ID.String(),
		UserID:    in.UserID.String(),
		ExpiresAt: in.ExpiresAt,
	}
}
