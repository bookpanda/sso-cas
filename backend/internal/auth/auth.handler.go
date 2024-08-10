package auth

import (
	"github.com/bookpanda/cas-sso/backend/internal/context"
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/validator"
	"go.uber.org/zap"
)

type Handler interface {
	GetGoogleLoginUrl(c context.Ctx)
	VerifyGoogleLogin(c context.Ctx)
}

type handlerImpl struct {
	svc      Service
	validate validator.DtoValidator
	log      *zap.Logger
}

func NewHandler(svc Service, validate validator.DtoValidator, log *zap.Logger) Handler {
	return &handlerImpl{
		svc:      svc,
		validate: validate,
		log:      log,
	}
}

func (h *handlerImpl) GetGoogleLoginUrl(c context.Ctx) {
	serviceUrl := c.Query("service")
	if serviceUrl == "" {
		c.BadRequestError("query parameter 'service' not found")
	}

	res, apperr := h.svc.GetGoogleLoginUrl(c.RequestContext(), serviceUrl)
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	c.JSON(200, res)
}

func (h *handlerImpl) VerifyGoogleLogin(c context.Ctx) {
	code := c.Query("code")
	if code == "" {
		c.BadRequestError("query parameter 'code' not found")
	}

	serviceUrl := c.Query("service")
	if serviceUrl == "" {
		c.BadRequestError("query parameter 'service' not found")
	}

	req := &dto.VerifyGoogleLoginRequest{
		Code:       code,
		ServiceUrl: serviceUrl,
	}

	credential, appErr := h.svc.VerifyGoogleLogin(c.RequestContext(), req)
	if appErr != nil {
		c.ResponseError(appErr)
		return
	}

	c.JSON(200, credential)
}
