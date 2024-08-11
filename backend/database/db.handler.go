package database

import (
	"github.com/bookpanda/cas-sso/backend/internal/context"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler interface {
	CleanDb(c context.Ctx)
}

type handlerImpl struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewHandler(db *gorm.DB, log *zap.Logger) Handler {
	return &handlerImpl{
		db:  db,
		log: log,
	}
}

func (h *handlerImpl) CleanDb(c context.Ctx) {
	err := h.db.Exec("TRUNCATE TABLE users, sessions, service_tickets RESTART IDENTITY CASCADE").Error
	if err != nil {
		h.log.Named("CleanDb").Error("Failed to clean database", zap.Error(err))
		c.InternalServerError("Failed to clean database")
		return
	}

	c.JSON(200, "Database cleaned successfully")
}
