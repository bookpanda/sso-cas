package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type GoogleOauthClient interface {
	GetUserEmail(code string) (string, error)
}

type googleOauthClientImpl struct {
	oauthConfig *oauth2.Config
	log         *zap.Logger
}

func NewGoogleOauthClient(oauthConfig *oauth2.Config, log *zap.Logger) GoogleOauthClient {
	return &googleOauthClientImpl{
		oauthConfig,
		log,
	}
}

var (
	ErrInvalidCode   = errors.New("invalid code")
	ErrHttpError     = errors.New("unable to get user info")
	ErrIOError       = errors.New("unable to read google response")
	ErrInvalidFormat = errors.New("google sent unexpected format")
)

func (c *googleOauthClientImpl) GetUserEmail(code string) (string, error) {
	token, err := c.oauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		c.log.Named("GetUserEmail").Error("Exchange: ", zap.Error(err))
		return "", ErrInvalidCode
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		c.log.Named("GetUserEmail").Error("Get: ", zap.Error(err))
		return "", ErrHttpError
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Named("GetUserEmail").Error("ReadAll: ", zap.Error(err))
		return "", ErrIOError
	}

	var parsedResponse dto.GoogleUserEmailResponse
	if err = json.Unmarshal(response, &parsedResponse); err != nil {
		c.log.Named("GetUserEmail").Error("Unmarshal: ", zap.Error(err))
		return "", ErrInvalidFormat
	}

	return parsedResponse.Email, nil
}
