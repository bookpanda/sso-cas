package main

import (
	"fmt"

	"github.com/bookpanda/cas-sso/backend/config"
	"github.com/bookpanda/cas-sso/backend/database"
	"github.com/bookpanda/cas-sso/backend/internal/auth"
	"github.com/bookpanda/cas-sso/backend/internal/auth/oauth"
	"github.com/bookpanda/cas-sso/backend/internal/router"
	"github.com/bookpanda/cas-sso/backend/internal/service_ticket"
	"github.com/bookpanda/cas-sso/backend/internal/session"
	"github.com/bookpanda/cas-sso/backend/internal/token"
	"github.com/bookpanda/cas-sso/backend/internal/user"
	"github.com/bookpanda/cas-sso/backend/internal/validator"
	"github.com/bookpanda/cas-sso/backend/logger"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	db, err := database.InitDatabase(&conf.Db, conf.App.IsDevelopment())
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	log := logger.New(&conf.App)
	validate, err := validator.NewDtoValidator()
	if err != nil {
		panic(fmt.Sprintf("Failed to create validator: %v", err))
	}

	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, log.Named("userSvc"))

	tokenSvc := token.NewService(log.Named("tokenSvc"))

	ticketRepo := service_ticket.NewRepository(db)
	ticketSvc := service_ticket.NewService(&conf.Auth, ticketRepo, tokenSvc, log.Named("ticketSvc"))

	sessionRepo := session.NewRepository(db)
	sessionSvc := session.NewService(&conf.Auth, sessionRepo, userSvc, tokenSvc, log.Named("sessionSvc"))

	oauthConfig := config.LoadOauthConfig(conf.Oauth)
	oauthClient := oauth.NewGoogleOauthClient(oauthConfig, log.Named("oauthClient"))
	authSvc := auth.NewService(oauthConfig, oauthClient, userSvc, ticketSvc, log.Named("authSvc"))
	authHdr := auth.NewHandler(authSvc, sessionSvc, ticketSvc, validate, log)

	corsHandler := config.MakeCorsConfig(conf)
	r := router.New(conf, corsHandler)

	r.V1Get("/auth/check-session", authHdr.CheckSession)
	r.V1Get("/auth/validate-st", authHdr.ValidateST)
	r.V1Get("/auth/google-url", authHdr.GetGoogleLoginUrl)
	r.V1Get("/auth/verify-google", authHdr.VerifyGoogleLogin)

	if err := r.Run(fmt.Sprintf(":%v", conf.App.Port)); err != nil {
		log.Fatal("unable to start server")
	}

}
