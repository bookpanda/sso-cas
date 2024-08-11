package service_ticket

import (
	"context"
	"errors"
	"time"

	"github.com/bookpanda/cas-sso/backend/apperror"
	"github.com/bookpanda/cas-sso/backend/config"
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/model"
	"github.com/bookpanda/cas-sso/backend/internal/token"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	FindByToken(_ context.Context, token string) (*dto.ServiceTicket, *apperror.AppError)
	Create(_ context.Context, req *dto.CreateServiceTicketRequest) (*dto.ServiceTicket, *apperror.AppError)
}

type serviceImpl struct {
	conf     *config.AuthConfig
	repo     Repository
	tokenSvc token.Service
	log      *zap.Logger
}

func NewService(conf *config.AuthConfig, repo Repository, tokenSvc token.Service, log *zap.Logger) Service {
	return &serviceImpl{
		conf:     conf,
		repo:     repo,
		tokenSvc: tokenSvc,
		log:      log,
	}
}

func (s *serviceImpl) FindByToken(ctx context.Context, token string) (*dto.ServiceTicket, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	serviceTicket := &model.ServiceTicket{}

	if err := s.repo.FindByToken(token, serviceTicket); err != nil {
		s.log.Named("FindByToken").Error("FindByToken: ", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError("session not found")
		}
		return nil, apperror.InternalServerError(err.Error())
	}

	if serviceTicket.ExpiresAt.Before(time.Now()) {
		if err := s.repo.DeleteByUserID(serviceTicket.UserID.String()); err != nil {
			s.log.Named("FindByToken").Error("DeleteByUserID: ", zap.Error(err))
			return nil, apperror.InternalServerError(err.Error())
		}
		return nil, apperror.NotFoundError("session not found")
	}

	return ModelToDto(serviceTicket), nil
}

func (s *serviceImpl) Create(ctx context.Context, req *dto.CreateServiceTicketRequest) (*dto.ServiceTicket, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.BadRequestError("invalid user id")
	}

	token, apperr := s.tokenSvc.GenerateOpaqueToken(ctx, 16)
	if apperr != nil {
		s.log.Named("Create").Error("GenerateOpaqueToken: ", zap.Error(apperr))
		return nil, apperr
	}

	createST := &model.ServiceTicket{
		Token:      "st_" + token,
		ServiceUrl: req.ServiceUrl,
		UserID:     userID,
		ExpiresAt:  time.Now().Add(time.Duration(s.conf.STTTL) * time.Second),
	}

	if err := s.repo.Create(createST); err != nil {
		s.log.Named("Create").Error("Create: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}

	return ModelToDto(createST), nil
}
