package auth

import (
	"github.com/bookpanda/cas-sso/backend/internal/context"
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/service_ticket"
	"github.com/bookpanda/cas-sso/backend/internal/session"
	"github.com/bookpanda/cas-sso/backend/internal/validator"
	"go.uber.org/zap"
)

type Handler interface {
	CheckSession(c context.Ctx)
	ValidateST(c context.Ctx)
	GetGoogleLoginUrl(c context.Ctx)
	VerifyGoogleLogin(c context.Ctx)
}

type handlerImpl struct {
	svc        Service
	sessionSvc session.Service
	ticketSvc  service_ticket.Service
	validate   validator.DtoValidator
	log        *zap.Logger
}

func NewHandler(svc Service, sessionSvc session.Service, ticketSvc service_ticket.Service, validate validator.DtoValidator, log *zap.Logger) Handler {
	return &handlerImpl{
		svc:        svc,
		sessionSvc: sessionSvc,
		ticketSvc:  ticketSvc,
		validate:   validate,
		log:        log,
	}
}

func (h *handlerImpl) CheckSession(c context.Ctx) {
	token, err := c.Cookie("CASTGC")
	if err != nil {
		h.log.Error("CheckSession: ", zap.Error(err))
		c.UnauthorizedError("'CASTGC' HTTP only cookie not found")
		return
	}

	serviceUrl := c.Query("service")
	if serviceUrl == "" {
		h.log.Error("CheckSession: query parameter 'service' not found")
		c.BadRequestError("query parameter 'service' not found")
		return
	}

	//check if the cookie is valid + not expired
	session, apperr := h.sessionSvc.FindByToken(c.RequestContext(), token)
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	serviceTicket, apperr := h.ticketSvc.Create(c.RequestContext(), &dto.CreateServiceTicketRequest{
		ServiceUrl: serviceUrl,
		UserID:     session.UserID,
	})
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	c.JSON(200, serviceTicket)
}

func (h *handlerImpl) ValidateST(c context.Ctx) {
	serviceTicket := c.Query("ticket")
	if serviceTicket == "" {
		h.log.Error("ValidateST: query parameter 'ticket' not found")
		c.BadRequestError("query parameter 'ticket' not found")
		return
	}

	res, apperr := h.ticketSvc.FindByToken(c.RequestContext(), serviceTicket)
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	c.JSON(200, res)
}

func (h *handlerImpl) GetGoogleLoginUrl(c context.Ctx) {
	serviceUrl := c.Query("service")
	if serviceUrl == "" {
		h.log.Error("GetGoogleLoginUrl: query parameter 'service' not found")
		c.BadRequestError("query parameter 'service' not found")
		return
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
		h.log.Error("VerifyGoogleLogin: query parameter 'code' not found")
		c.BadRequestError("query parameter 'code' not found")
		return
	}

	serviceUrl := c.Query("state")
	if serviceUrl == "" {
		h.log.Error("VerifyGoogleLogin: query parameter 'state' not found")
		c.BadRequestError("query parameter 'state' not found")
		return
	}

	req := &dto.VerifyGoogleLoginRequest{
		Code:       code,
		ServiceUrl: serviceUrl,
	}

	res, apperr := h.svc.VerifyGoogleLogin(c.RequestContext(), req)
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	serviceTicket, apperr := h.ticketSvc.Create(c.RequestContext(), &dto.CreateServiceTicketRequest{
		ServiceUrl: serviceUrl,
		UserID:     res.User.ID,
	})
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	session, apperr := h.sessionSvc.Create(c.RequestContext(), &dto.CreateSessionRequest{
		UserID:     res.User.ID,
		ServiceUrl: serviceUrl,
	})
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	c.SetCookie("CASTGC", session.Token, 0, "/", "", false, true)

	c.JSON(200, &dto.VerifyGoogleLoginResponse{
		ServiceTicket: serviceTicket.Token,
	})
}
