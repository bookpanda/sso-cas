package user

import (
	"context"
	"errors"
	"time"

	"github.com/bookpanda/cas-sso/backend/apperror"
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/model"
	"go.uber.org/zap"

	"gorm.io/gorm"
)

type Service interface {
	FindOne(_ context.Context, id string) (*model.User, *apperror.AppError)
	FindByEmail(_ context.Context, email string) (*model.User, *apperror.AppError)
	Create(_ context.Context, req *dto.CreateUserRequest) (*model.User, *apperror.AppError)
	UpdateProfile(_ context.Context, req *dto.UpdateUserProfileRequest) (*model.User, *apperror.AppError)
}

type serviceImpl struct {
	repo Repository
	log  *zap.Logger
}

func NewService(repo Repository, log *zap.Logger) Service {
	return &serviceImpl{
		repo: repo,
		log:  log,
	}
}

func (s *serviceImpl) FindOne(ctx context.Context, id string) (*model.User, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &model.User{}

	if err := s.repo.FindOne(id, user); err != nil {
		s.log.Named("FindOne").Error("FindOne: ", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError("user not found")
		}
		return nil, apperror.InternalServerError(err.Error())
	}

	return user, nil
}

func (s *serviceImpl) FindByEmail(ctx context.Context, email string) (*model.User, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &model.User{}

	if err := s.repo.FindByEmail(email, user); err != nil {
		s.log.Named("FindByEmail").Error("FindByEmail: ", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError("user not found")
		}
		return nil, apperror.InternalServerError(err.Error())
	}

	return user, nil
}

func (s *serviceImpl) Create(ctx context.Context, req *dto.CreateUserRequest) (*model.User, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	createUser := &model.User{
		Email: req.Email,
	}

	if err := s.repo.Create(createUser); err != nil {
		s.log.Named("Create").Error("Create: ", zap.Error(err))
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.BadRequestError("duplicate email")
		}
		return nil, apperror.InternalServerError(err.Error())
	}

	return createUser, nil
}

func (s *serviceImpl) UpdateProfile(ctx context.Context, req *dto.UpdateUserProfileRequest) (*model.User, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	updateUser, err := UpdateRequestToModel(req)
	if err != nil {
		s.log.Named("Update").Error("UpdateRequestToModel: ", zap.Error(err))
		return nil, apperror.BadRequestError("error parsing request: " + err.Error())
	}

	if err = s.repo.Update(req.ID, updateUser); err != nil {
		s.log.Named("Update").Error("Update: ", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError("user not found")
		}
		return nil, apperror.InternalServerError(err.Error())
	}

	return updateUser, nil
}
