package service_ticket

import (
	"context"
	"time"

	"github.com/bookpanda/cas-sso/backend/apperror"
	"github.com/bookpanda/cas-sso/backend/config"
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	Create(_ context.Context, req *dto.CreateServiceTicketRequest) (*dto.ServiceTicket, *apperror.AppError)
}

type serviceImpl struct {
	conf *config.AuthConfig
	repo Repository
	log  *zap.Logger
}

func NewService(conf *config.AuthConfig, repo Repository, log *zap.Logger) Service {
	return &serviceImpl{
		conf: conf,
		repo: repo,
		log:  log,
	}
}

func (s *serviceImpl) Create(ctx context.Context, req *dto.CreateServiceTicketRequest) (*dto.ServiceTicket, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.BadRequestError("invalid user id")
	}

	createST := &model.ServiceTicket{
		ServiceUrl: req.ServiceUrl,
		UserID:     userID,
		ExpiresAt:  s.conf.STTTL,
	}

	if err := s.repo.Create(createST); err != nil {
		s.log.Named("Create").Error("Create: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}

	return ModelToDto(createST), nil
}
