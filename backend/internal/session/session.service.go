package session

import (
	"context"
	"errors"
	"time"

	"github.com/bookpanda/cas-sso/backend/apperror"
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/model"
	"github.com/bookpanda/cas-sso/backend/internal/user"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"gorm.io/gorm"
)

type Service interface {
	FindOne(_ context.Context, id string) (*dto.Session, *apperror.AppError)
	Create(_ context.Context, req *dto.CreateSessionRequest) (*dto.Session, *apperror.AppError)
	DeleteByEmail(_ context.Context, req *dto.DeleteByEmailSessionRequest) (*dto.SuccessResponse, *apperror.AppError)
}

type serviceImpl struct {
	repo    Repository
	userSvc user.Service
	log     *zap.Logger
}

func NewService(repo Repository, userSvc user.Service, log *zap.Logger) Service {
	return &serviceImpl{
		repo:    repo,
		userSvc: userSvc,
		log:     log,
	}
}

func (s *serviceImpl) FindOne(ctx context.Context, id string) (*dto.Session, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	session := &model.Session{}

	if err := s.repo.FindOne(id, session); err != nil {
		s.log.Named("FindOne").Error("FindOne: ", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError("session not found")
		}
		return nil, apperror.InternalServerError(err.Error())
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

	createSession := &model.Session{
		UserID: userID,
	}

	if err := s.repo.Create(createSession); err != nil {
		s.log.Named("Create").Error("Create: ", zap.Error(err))
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.BadRequestError("duplicate email")
		}
		return nil, apperror.InternalServerError(err.Error())
	}

	return ModelToDto(createSession), nil
}

func (s *serviceImpl) DeleteByEmail(ctx context.Context, req *dto.DeleteByEmailSessionRequest) (*dto.SuccessResponse, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, apperr := s.userSvc.FindByEmail(ctx, req.Email)
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
