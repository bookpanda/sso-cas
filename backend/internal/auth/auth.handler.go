package auth

import (
	"github.com/bookpanda/cas-sso/backend/config"
	"github.com/bookpanda/cas-sso/backend/internal/context"
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/service_ticket"
	_session "github.com/bookpanda/cas-sso/backend/internal/session"
	"github.com/bookpanda/cas-sso/backend/internal/validator"
	"go.uber.org/zap"
)

type Handler interface {
	CheckSession(c context.Ctx)
	ValidateST(c context.Ctx)
	GetGoogleLoginUrl(c context.Ctx)
	VerifyGoogleLogin(c context.Ctx)
	Signout(c context.Ctx)
}

type handlerImpl struct {
	conf       *config.AuthConfig
	svc        Service
	sessionSvc _session.Service
	ticketSvc  service_ticket.Service
	validate   validator.DtoValidator
	log        *zap.Logger
}

func NewHandler(conf *config.AuthConfig, svc Service, sessionSvc _session.Service, ticketSvc service_ticket.Service, validate validator.DtoValidator, log *zap.Logger) Handler {
	return &handlerImpl{
		conf:       conf,
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

	//check if the session is valid + not expired
	session, apperr := h.sessionSvc.FindByToken(c.RequestContext(), token)
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	serviceTicket, apperr := h.ticketSvc.Create(c.RequestContext(), &dto.CreateServiceTicketRequest{
		SessionToken: session.Token,
		ServiceUrl:   serviceUrl,
		UserID:       session.UserID.String(),
	})
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	c.JSON(200, &dto.ServiceTicketToken{
		ServiceTicket: serviceTicket.Token,
	})
}

func (h *handlerImpl) ValidateST(c context.Ctx) {
	serviceTicketToken := c.Query("ticket")
	if serviceTicketToken == "" {
		h.log.Error("ValidateST: query parameter 'ticket' not found")
		c.BadRequestError("query parameter 'ticket' not found")
		return
	}

	service := c.Query("service")
	if serviceTicketToken == "" {
		h.log.Error("ValidateST: query parameter 'service' not found")
		c.BadRequestError("query parameter 'service' not found")
		return
	}

	serviceTicket, apperr := h.ticketSvc.FindByToken(c.RequestContext(), serviceTicketToken)
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	h.log.Info("ValidateST: ", zap.String("service", service), zap.String("serviceTicket.ServiceUrl", serviceTicket.ServiceUrl))
	if service != serviceTicket.ServiceUrl {
		h.log.Error("ValidateST: 'service' query parameter does not match the service ticket")
		c.UnauthorizedError("'service' query parameter does not match the service ticket")
		return
	}

	_, apperr = h.ticketSvc.DeleteByToken(c.RequestContext(), serviceTicket.Token)
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	session, apperr := h.sessionSvc.FindByToken(c.RequestContext(), serviceTicket.SessionToken)
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	c.JSON(200, _session.ModelToDto(session))
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

	session, apperr := h.sessionSvc.Create(c.RequestContext(), &dto.CreateSessionRequest{
		UserID: res.User.ID,
	})
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	serviceTicket, apperr := h.ticketSvc.Create(c.RequestContext(), &dto.CreateServiceTicketRequest{
		SessionToken: session.Token,
		ServiceUrl:   serviceUrl,
		UserID:       res.User.ID,
	})
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	c.SetCookie("CASTGC", session.Token, h.conf.SessionTTL, "/", "localhost", false, true)

	c.JSON(200, &dto.ServiceTicketToken{
		ServiceTicket: serviceTicket.Token,
	})
}

func (h *handlerImpl) Signout(c context.Ctx) {
	token, err := c.Cookie("CASTGC")
	if err != nil {
		h.log.Error("Signout: ", zap.Error(err))
		c.UnauthorizedError("'CASTGC' HTTP only cookie not found")
		return
	}

	session, apperr := h.sessionSvc.FindByToken(c.RequestContext(), token)
	if apperr != nil {
		c.ResponseError(apperr)
		return
	}

	// remove cookie
	c.SetCookie("CASTGC", "", -1, "/", "localhost", false, true)

	for _, service := range h.conf.Services {
		apperr := h.svc.SignoutService(c.RequestContext(), service, session.UserID.String())
		if apperr != nil {
			h.log.Error("Signout: ", zap.Error(apperr))
			c.ResponseError(apperr)
			return
		}
	}

	c.JSON(200, nil)
}
