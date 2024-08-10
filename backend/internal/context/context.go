package context

import (
	"context"

	"github.com/bookpanda/cas-sso/backend/apperror"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Ctx interface {
	JSON(statusCode int, obj interface{})
	ResponseError(err *apperror.AppError)
	BadRequestError(err string)
	ForbiddenError(err string)
	InternalServerError(err string)
	NewUUID() uuid.UUID
	Bind(obj interface{}) error
	Param(key string) string
	Query(key string) string
	PostForm(key string) string
	GetString(key string) string
	GetHeader(key string) string
	RequestContext() context.Context
	Abort()
	Next()
}

type contextImpl struct {
	*gin.Context
}

func NewContext(c *gin.Context) Ctx {
	return &contextImpl{
		Context: c,
	}
}

func (c *contextImpl) RequestContext() context.Context {
	return c.Context.Request.Context()
}

func (c *contextImpl) JSON(statusCode int, obj interface{}) {
	c.Context.JSON(statusCode, obj)
}

func (c *contextImpl) ResponseError(err *apperror.AppError) {
	c.JSON(err.HttpCode, gin.H{"error": err.Error()})
}

func (c *contextImpl) BadRequestError(err string) {
	c.JSON(apperror.BadRequest.HttpCode, gin.H{"error": err})
}

func (c *contextImpl) UnauthorizedError(err string) {
	c.JSON(apperror.Unauthorized.HttpCode, gin.H{"error": err})
}

func (c *contextImpl) ForbiddenError(err string) {
	c.JSON(apperror.Forbidden.HttpCode, gin.H{"error": err})
}

func (c *contextImpl) NotFoundError(err string) {
	c.JSON(apperror.NotFound.HttpCode, gin.H{"error": err})
}

func (c *contextImpl) InternalServerError(err string) {
	c.JSON(apperror.InternalServer.HttpCode, gin.H{"error": err})
}

func (c *contextImpl) ServiceUnavailableError(err string) {
	c.JSON(apperror.ServiceUnavailable.HttpCode, gin.H{"error": err})
}

func (c *contextImpl) NewUUID() uuid.UUID {
	return uuid.New()
}

func (c *contextImpl) Bind(obj interface{}) error {
	return c.Context.Bind(obj)
}

func (c *contextImpl) Param(key string) string {
	return c.Context.Param(key)
}

func (c *contextImpl) Query(key string) string {
	return c.Context.Query(key)
}

func (c *contextImpl) PostForm(key string) string {
	return c.Context.PostForm(key)
}

func (c *contextImpl) GetString(key string) string {
	return c.Context.GetString(key)
}

func (c *contextImpl) GetHeader(key string) string {
	return c.Context.GetHeader(key)
}

func (c *contextImpl) Abort() {
	c.Context.Abort()
}

func (c *contextImpl) Next() {
	c.Context.Next()
}
