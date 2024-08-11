package session

import (
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/model"
)

func ModelToDto(in *model.Session) *dto.Session {
	return &dto.Session{
		Token:     in.Token,
		UserID:    in.UserID.String(),
		Payload:   in.Payload,
		ExpiresAt: in.ExpiresAt,
	}
}
