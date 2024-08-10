package service_ticket

import (
	"context"
	"time"

	"github.com/bookpanda/cas-sso/backend/apperror"
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"go.uber.org/zap"
)

type Service interface {
	Create(_ context.Context, req *dto.CreateServiceTicketRequest) (*dto.ServiceTicket, *apperror.AppError)
}

type serviceImpl struct {
	log *zap.Logger
}

func NewService(log *zap.Logger) Service {
	return &serviceImpl{
		log: log,
	}
}

func (s *serviceImpl) Create(ctx context.Context, req *dto.CreateServiceTicketRequest) (*dto.ServiceTicket, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return &dto.ServiceTicket{
		ServiceUrl: req.ServiceUrl,
	}, nil
}
