package service_ticket

import (
	"github.com/bookpanda/cas-sso/backend/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindByToken(token string, serviceTicket *model.ServiceTicket) error
	Create(serviceTicket *model.ServiceTicket) error
	DeleteByToken(token string) error
}

type repositoryImpl struct {
	Db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryImpl{Db: db}
}

func (r *repositoryImpl) FindByToken(token string, serviceTicket *model.ServiceTicket) error {
	return r.Db.Model(serviceTicket).First(serviceTicket, "token = ?", token).Error
}

func (r *repositoryImpl) Create(serviceTicket *model.ServiceTicket) error {
	return r.Db.Create(serviceTicket).Error
}

func (r *repositoryImpl) DeleteByToken(token string) error {
	return r.Db.Where("token = ?", token).Delete(&model.ServiceTicket{}).Error
}
