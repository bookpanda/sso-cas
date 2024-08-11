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
	"github.com/bookpanda/cas-sso/backend/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"gorm.io/gorm"
)

type Service interface {
	FindByToken(_ context.Context, token string) (*model.Session, *apperror.AppError)
	Create(_ context.Context, req *dto.CreateSessionRequest) (*model.Session, *apperror.AppError)
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

func (s *serviceImpl) FindByToken(ctx context.Context, token string) (*model.Session, *apperror.AppError) {
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

	localExpire, err := utils.ParseLocalTime(session.ExpiresAt)
	if err != nil {
		s.log.Named("Create").Error("ParseLocalTime: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}

	if localExpire.Before(time.Now()) {
		if err := s.repo.DeleteByUserID(session.UserID.String()); err != nil {
			s.log.Named("FindByToken").Error("DeleteByUserID: ", zap.Error(err))
			return nil, apperror.InternalServerError(err.Error())
		}
		return nil, apperror.NotFoundError("session not found")
	}

	return session, nil
}

func (s *serviceImpl) Create(ctx context.Context, req *dto.CreateSessionRequest) (*model.Session, *apperror.AppError) {
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

	localExpire, err := utils.ParseLocalTime(time.Now().Add(time.Duration(s.conf.SessionTTL) * time.Second))
	if err != nil {
		s.log.Named("Create").Error("ParseLocalTime: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}

	createSession := &model.Session{
		Token:     "session_" + token,
		UserID:    userID,
		ExpiresAt: localExpire,
	}

	if err := s.repo.Create(createSession); err != nil {
		s.log.Named("Create").Error("Create: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}

	return createSession, nil
}

func (s *serviceImpl) DeleteByEmail(ctx context.Context, email string) (*dto.SuccessResponse, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, apperr := s.userSvc.FindByEmail(ctx, email)
	if apperr != nil {
		s.log.Named("DeleteByEmail").Error("FindByEmail: ", zap.Error(apperr))
		return nil, apperr
	}

	if err := s.repo.DeleteByUserID(user.ID.String()); err != nil {
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
