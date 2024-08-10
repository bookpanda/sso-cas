package logger

import (
	"github.com/bookpanda/cas-sso/backend/config"
	"go.uber.org/zap"
)

func New(conf *config.AppConfig) *zap.Logger {
	var logger *zap.Logger

	if conf.IsDevelopment() {
		logger = zap.Must(zap.NewDevelopment())
	} else {
		logger = zap.Must(zap.NewProduction())
	}

	return logger
}
