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
	"github.com/bookpanda/cas-sso/backend/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	FindByToken(_ context.Context, token string) (*model.ServiceTicket, *apperror.AppError)
	Create(_ context.Context, req *dto.CreateServiceTicketRequest) (*model.ServiceTicket, *apperror.AppError)
	DeleteByToken(_ context.Context, token string) (*dto.SuccessResponse, *apperror.AppError)
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

func (s *serviceImpl) FindByToken(ctx context.Context, token string) (*model.ServiceTicket, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	serviceTicket := &model.ServiceTicket{}

	if err := s.repo.FindByToken(token, serviceTicket); err != nil {
		s.log.Named("FindByToken").Error("FindByToken: ", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError("service ticket not found")
		}
		return nil, apperror.InternalServerError(err.Error())
	}

	localExpire, err := utils.ParseLocalTime(serviceTicket.ExpiresAt)
	if err != nil {
		s.log.Named("Create").Error("ParseLocalTime localExpire: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}

	localNow, err := utils.ParseLocalTime(time.Now())
	if err != nil {
		s.log.Named("Create").Error("ParseLocalTime localNow: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}
	if localExpire.Before(localNow) {
		if err := s.repo.DeleteByToken(serviceTicket.Token); err != nil {
			s.log.Named("FindByToken").Error("DeleteByUserID: ", zap.Error(err))
			return nil, apperror.InternalServerError(err.Error())
		}
		return nil, apperror.NotFoundError("service ticket not found")
	}

	return serviceTicket, nil
}

func (s *serviceImpl) Create(ctx context.Context, req *dto.CreateServiceTicketRequest) (*model.ServiceTicket, *apperror.AppError) {
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

	localNow, err := utils.ParseLocalTime(time.Now())
	if err != nil {
		s.log.Named("Create").Error("ParseLocalTime localNow: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}
	localExpire, err := utils.ParseLocalTime(localNow.Add(time.Duration(s.conf.STTTL) * time.Second))
	if err != nil {
		s.log.Named("Create").Error("ParseLocalTime localExpire: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}

	createST := &model.ServiceTicket{
		SessionToken: req.SessionToken,
		Token:        "st_" + token,
		ServiceUrl:   req.ServiceUrl,
		UserID:       userID,
		ExpiresAt:    localExpire,
	}

	if err := s.repo.Create(createST); err != nil {
		s.log.Named("Create").Error("Create: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}

	return createST, nil
}

func (s *serviceImpl) DeleteByToken(ctx context.Context, token string) (*dto.SuccessResponse, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.repo.DeleteByToken(token); err != nil {
		s.log.Named("DeleteByToken").Error("DeleteByToken: ", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError("service ticket not found")
		}
		return nil, apperror.InternalServerError(err.Error())
	}

	return &dto.SuccessResponse{
		Message: "Service ticket deleted",
	}, nil
}
