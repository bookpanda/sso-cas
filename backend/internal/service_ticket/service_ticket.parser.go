package service_ticket

import (
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/model"
)

func ModelToDto(in *model.ServiceTicket) *dto.ServiceTicket {
	return &dto.ServiceTicket{
		SessionToken: in.SessionToken,
		Token:        in.Token,
		ServiceUrl:   in.ServiceUrl,
		UserID:       in.UserID.String(),
		ExpiresAt:    in.ExpiresAt,
	}
}
