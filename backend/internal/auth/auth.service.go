package auth

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/bookpanda/cas-sso/backend/apperror"
	"github.com/bookpanda/cas-sso/backend/internal/auth/oauth"
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/service_ticket"
	"github.com/bookpanda/cas-sso/backend/internal/user"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Service interface {
	IssueST(ctx context.Context, req *dto.IssueSTRequest) (*dto.ServiceTicket, *apperror.AppError)
	GetGoogleLoginUrl(_ context.Context, serviceUrl string) (string, *apperror.AppError)
	VerifyGoogleLogin(_ context.Context, req *dto.VerifyGoogleLoginRequest) (*dto.VerifyGoogleLoginResponse, *apperror.AppError)
}

type serviceImpl struct {
	oauthConfig *oauth2.Config
	oauthClient oauth.GoogleOauthClient
	userSvc     user.Service
	ticketSvc   service_ticket.Service
	log         *zap.Logger
}

func NewService(oauthConfig *oauth2.Config, oauthClient oauth.GoogleOauthClient, userSvc user.Service, ticketSvc service_ticket.Service, log *zap.Logger) Service {
	return &serviceImpl{
		oauthConfig: oauthConfig,
		oauthClient: oauthClient,
		userSvc:     userSvc,
		ticketSvc:   ticketSvc,
		log:         log,
	}
}

func (s *serviceImpl) IssueST(ctx context.Context, req *dto.IssueSTRequest) (*dto.ServiceTicket, *apperror.AppError) {
	createTicket := &dto.CreateServiceTicketRequest{
		ServiceUrl: req.ServiceUrl,
		UserID:     req.UserID,
	}

	serviceTicket, err := s.ticketSvc.Create(ctx, createTicket)
	if err != nil {
		s.log.Named("VerifyGoogleLogin").Error("Create: ", zap.Error(err))
		return nil, apperror.InternalServerError(err.Error())
	}

	return serviceTicket, nil
}

func (s *serviceImpl) GetGoogleLoginUrl(_ context.Context, serviceUrl string) (string, *apperror.AppError) {
	URL, err := url.Parse(s.oauthConfig.Endpoint.AuthURL)
	if err != nil {
		s.log.Named("GetGoogleLoginUrl").Error("Parse: ", zap.Error(err))
		return "", apperror.InternalServerError("Cannot parse Google OAuth URL")
	}

	parameters := url.Values{}
	parameters.Add("client_id", s.oauthConfig.ClientID)
	parameters.Add("scope", strings.Join(s.oauthConfig.Scopes, " "))
	parameters.Add("redirect_uri", s.oauthConfig.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", serviceUrl) // encode base64
	URL.RawQuery = parameters.Encode()
	loginUrl := URL.String()

	return loginUrl, nil
}

func (s *serviceImpl) VerifyGoogleLogin(ctx context.Context, req *dto.VerifyGoogleLoginRequest) (*dto.VerifyGoogleLoginResponse, *apperror.AppError) {
	if req.Code == "" {
		return nil, apperror.BadRequestError("No code is provided")
	}

	email, err := s.oauthClient.GetUserEmail(req.Code)
	if err != nil {
		s.log.Named("VerifyGoogleLogin").Error("GetUserEmail: ", zap.Error(err))
		switch err.Error() {
		case "Invalid code":
			return nil, apperror.BadRequestError("Invalid code")
		default:
			return nil, apperror.InternalServerError(err.Error())
		}
	}

	user, apperr := s.userSvc.FindByEmail(context.Background(), email)
	if apperr != nil {
		switch apperr.HttpCode {
		case http.StatusNotFound:
			s.log.Named("VerifyGoogleLogin").Info("User not found, creating new user")

			createUser := &dto.CreateUserRequest{
				Email: email,
			}
			user, err := s.userSvc.Create(context.Background(), createUser)
			if err != nil {
				s.log.Named("VerifyGoogleLogin").Error("Create: ", zap.Error(err))
				return nil, err
			}

			authReq := &dto.IssueSTRequest{
				ServiceUrl: req.ServiceUrl,
				UserID:     user.ID,
			}
			serviceTicket, apperr := s.IssueST(ctx, authReq)
			if apperr != nil {
				s.log.Named("VerifyGoogleLogin").Error("IssueST: ", zap.Error(apperr))
				return nil, apperr
			}

			return &dto.VerifyGoogleLoginResponse{
				ServiceTicket: serviceTicket,
				SessionCookie: "session_cookie",
			}, nil

		default:
			s.log.Named("VerifyGoogleLogin").Error("FindByEmail: ", zap.Error(apperr))
			return nil, apperr
		}
	}

	authReq := &dto.IssueSTRequest{
		ServiceUrl: req.ServiceUrl,
		UserID:     user.ID,
	}
	serviceTicket, apperr := s.IssueST(ctx, authReq)
	if apperr != nil {
		s.log.Named("VerifyGoogleLogin").Error("IssueST: ", zap.Error(apperr))
		return nil, apperr
	}

	return &dto.VerifyGoogleLoginResponse{
		ServiceTicket: serviceTicket,
		SessionCookie: "session_cookie",
	}, nil
}
