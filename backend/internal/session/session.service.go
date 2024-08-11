package session

import (
	"context"
	"errors"
	"time"

	"github.com/bookpanda/cas-sso/backend/apperror"
	"github.com/bookpanda/cas-sso/backend/config"
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/model"
	"github.com/bookpanda/cas-sso/backend/internal/token"
	"github.com/bookpanda/cas-sso/backend/internal/user"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"gorm.io/gorm"
)

type Service interface {
	FindByToken(_ context.Context, token string) (*dto.Session, *apperror.AppError)
	Create(_ context.Context, req *dto.CreateSessionRequest) (*dto.Session, *apperror.AppError)
	DeleteByEmail(_ context.Context, email string) (*dto.SuccessResponse, *apperror.AppError)
}

type serviceImpl struct {
	conf     *config.AuthConfig
	repo     Repository
	userSvc  user.Service
	tokenSvc token.Service
	log      *zap.Logger
}

func NewService(conf *config.AuthConfig, repo Repository, userSvc user.Service, tokenSvc token.Service, log *zap.Logger) Service {
	return &serviceImpl{
		conf:     conf,
		repo:     repo,
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
		log:      log,
	}
}

func (s *serviceImpl) FindByToken(ctx context.Context, token string) (*dto.Session, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	session := &model.Session{}

	if err := s.repo.FindByToken(token, session); err != nil {
		s.log.Named("FindByToken").Error("FindByToken: ", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError("session not found")
		}
		return nil, apperror.InternalServerError(err.Error())
	}

	if session.ExpiresAt.Before(time.Now()) {
		if err := s.repo.DeleteByUserID(session.UserID.String()); err != nil {
			s.log.Named("FindByToken").Error("DeleteByUserID: ", zap.Error(err))
			return nil, apperror.InternalServerError(err.Error())
		}
		return nil, apperror.NotFoundError("session not found")
	}

	return ModelToDto(session), nil
}

func (s *serviceImpl) Create(ctx context.Context, req *dto.CreateSessionRequest) (*dto.Session, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, apperror.BadRequestError("invalid user id")
	}

	token, apperr := s.tokenSvc.GenerateOpaqueToken(ctx, 32)
	if apperr != nil {
		s.log.Named("Create").Error("GenerateOpaqueToken: ", zap.Error(apperr))
		return nil, apperr
	}

	createSession := &model.Session{
		Token:      "session_" + token,
		ServiceUrl: req.ServiceUrl,
		UserID:     userID,
		ExpiresAt:  time.Now().Add(time.Duration(s.conf.SessionTTL) * time.Second),
	}

	if err := s.repo.Create(createSession); err != nil {
		s.log.Named("Create").Error("Create: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}

	return ModelToDto(createSession), nil
}

func (s *serviceImpl) DeleteByEmail(ctx context.Context, email string) (*dto.SuccessResponse, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, apperr := s.userSvc.FindByEmail(ctx, email)
	if apperr != nil {
		s.log.Named("DeleteByEmail").Error("FindByEmail: ", zap.Error(apperr))
		return nil, apperr
	}

	if err := s.repo.DeleteByUserID(user.ID); err != nil {
		s.log.Named("DeleteByEmail").Error("DeleteByUserID: ", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError("session not found")
		}
		return nil, apperror.InternalServerError(err.Error())
	}

	return &dto.SuccessResponse{
		Message: "session deleted",
	}, nil
}
