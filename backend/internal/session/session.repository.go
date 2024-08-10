package session

import (
	"github.com/bookpanda/cas-sso/backend/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindOne(id string, session *model.Session) error
	Create(session *model.Session) error
	DeleteByUserID(userID string) error
}

type repositoryImpl struct {
	Db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{Db: db}
}

func (r *repositoryImpl) FindOne(id string, session *model.Session) error {
	return r.Db.Model(session).First(session, "id = ?", id).Error
}

func (r *repositoryImpl) Create(session *model.Session) error {
	return r.Db.Create(session).Error
}

func (r *repositoryImpl) DeleteByUserID(userID string) error {
	return r.Db.Where("user_id = ?", userID).Delete(&model.Session{}).Error
}
