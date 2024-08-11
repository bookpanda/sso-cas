package token

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/bookpanda/cas-sso/backend/apperror"
	"go.uber.org/zap"
)

type Service interface {
	GenerateOpaqueToken(_ context.Context, length int) (string, *apperror.AppError)
}

type serviceImpl struct {
	log *zap.Logger
}

func NewService(log *zap.Logger) Service {
	return &serviceImpl{
		log: log,
	}
}

func (s *serviceImpl) GenerateOpaqueToken(ctx context.Context, length int) (string, *apperror.AppError) {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", apperror.InternalServerError(err.Error())
	}
	return hex.EncodeToString(bytes), nil
}
