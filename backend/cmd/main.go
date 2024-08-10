package main

import (
	"fmt"

	"github.com/bookpanda/cas-sso/backend/config"
	"github.com/bookpanda/cas-sso/backend/database"
	"github.com/bookpanda/cas-sso/backend/internal/auth"
	"github.com/bookpanda/cas-sso/backend/internal/auth/oauth"
	"github.com/bookpanda/cas-sso/backend/internal/router"
	"github.com/bookpanda/cas-sso/backend/internal/service_ticket"
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

	corsHandler := config.MakeCorsConfig(conf)
	r := router.New(conf, corsHandler)

	if err := r.Run(fmt.Sprintf(":%v", conf.App.Port)); err != nil {
		log.Fatal("unable to start server")
	}

	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo, log.Named("userSvc"))

	ticketSvc := service_ticket.NewService(log.Named("ticketSvc"))

	oauthConfig := config.LoadOauthConfig(conf.Oauth)
	oauthClient := oauth.NewGoogleOauthClient(oauthConfig, log.Named("oauthClient"))
	authSvc := auth.NewService(oauthConfig, oauthClient, userSvc, ticketSvc, log.Named("authSvc"))
	authHdr := auth.NewHandler(authSvc, validate, log)

	r.V1Get("/auth/google-url", authHdr.GetGoogleLoginUrl)
	r.V1Get("/auth/verify-google", authHdr.VerifyGoogleLogin)

}
