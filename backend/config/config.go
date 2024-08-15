package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AppConfig struct {
	Port string
	Env  string
}

type DbConfig struct {
	Url string
}

type OauthConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUri  string
}

type AuthConfig struct {
	STTTL      int
	SessionTTL int
	Services   []string
	IsHTTPS    bool
}

type CorsConfig struct {
	AllowOrigins string
}

type Config struct {
	App   AppConfig
	Db    DbConfig
	Oauth OauthConfig
	Auth  AuthConfig
	Cors  CorsConfig
}

func LoadConfig() (*Config, error) {
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	}

	appConfig := AppConfig{
		Port: os.Getenv("APP_PORT"),
		Env:  os.Getenv("APP_ENV"),
	}

	corsConfig := CorsConfig{
		AllowOrigins: os.Getenv("CORS_ALLOW_ORIGINS"),
	}

	dbConfig := DbConfig{
		Url: os.Getenv("DB_URL"),
	}

	oauthConfig := OauthConfig{
		ClientId:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectUri:  os.Getenv("OAUTH_REDIRECT_URI"),
	}

	STTTL, err := strconv.ParseInt(os.Getenv("AUTH_ST_TTL"), 10, 64)
	if err != nil {
		return nil, err
	}
	sessionTTL, err := strconv.ParseInt(os.Getenv("AUTH_SESSION_TTL"), 10, 64)
	if err != nil {
		return nil, err
	}
	servicesLogoutString := os.Getenv("AUTH_SERVICES_LOGOUT")
	servicesLogout := strings.Split(servicesLogoutString, ",")
	authConfig := AuthConfig{
		STTTL:      int(STTTL),
		SessionTTL: int(sessionTTL),
		Services:   servicesLogout,
		IsHTTPS:    os.Getenv("AUTH_IS_HTTPS") == "true",
	}

	return &Config{
		App:   appConfig,
		Db:    dbConfig,
		Oauth: oauthConfig,
		Auth:  authConfig,
		Cors:  corsConfig,
	}, nil
}

func (ac *AppConfig) IsDevelopment() bool {
	return ac.Env == "development"
}

func LoadOauthConfig(oauth OauthConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     oauth.ClientId,
		ClientSecret: oauth.ClientSecret,
		RedirectURL:  oauth.RedirectUri,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}
}
