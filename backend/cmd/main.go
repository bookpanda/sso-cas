package main

import (
	"fmt"

	"github.com/MrRehab/rehab-backend/config"
	"github.com/MrRehab/rehab-backend/database"
	"github.com/MrRehab/rehab-backend/internal/auth"
	"github.com/MrRehab/rehab-backend/internal/auth/oauth"
	"github.com/MrRehab/rehab-backend/internal/auth/token"
	"github.com/MrRehab/rehab-backend/internal/auth/token/jwt"
	"github.com/MrRehab/rehab-backend/internal/cache"
	"github.com/MrRehab/rehab-backend/internal/doctor"
	"github.com/MrRehab/rehab-backend/internal/middleware"
	"github.com/MrRehab/rehab-backend/internal/patient"
	"github.com/MrRehab/rehab-backend/internal/router"
	"github.com/MrRehab/rehab-backend/internal/validator"
	"github.com/MrRehab/rehab-backend/logger"
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

	redis, err := database.InitRedis(&conf.Redis)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to redis: %v", err))
	}

	log := logger.New(&conf.App)
	validate, err := validator.NewDtoValidator()
	if err != nil {
		panic(fmt.Sprintf("Failed to create validator: %v", err))
	}

	cacheRepo := cache.NewRepository(redis)

	doctorRepo := doctor.NewRepository(db)
	doctorSvc := doctor.NewService(doctorRepo, log.Named("doctorSvc"))

	patientRepo := patient.NewRepository(db)
	patientSvc := patient.NewService(patientRepo, log.Named("patientSvc"))

	jwtSvc := jwt.NewService(conf.Jwt, jwt.NewJwtStrategy(conf.Jwt.Secret), jwt.NewJwtUtils(), log.Named("jwtSvc"))
	tokenSvc := token.NewService(jwtSvc, cacheRepo, token.NewTokenUtils(), log.Named("tokenSvc"))
	oauthConfig := config.LoadOauthConfig(conf.Oauth)
	oauthClient := oauth.NewGoogleOauthClient(oauthConfig, log.Named("oauthClient"))
	authSvc := auth.NewService(oauthConfig, oauthClient, doctorSvc, patientSvc, tokenSvc, log.Named("authSvc"))
	authHdr := auth.NewHandler(authSvc, validate, log.Named("authHdr"))

	corsHandler := config.MakeCorsConfig(conf)
	authMiddleware := middleware.NewAuthMiddleware(authSvc)
	r := router.New(conf, corsHandler, authMiddleware)

	if err := r.Run(fmt.Sprintf(":%v", conf.App.Port)); err != nil {
		log.Fatal("unable to start server")
	}

	r.V1NonAuthGet("/auth/google-url", authHdr.GetGoogleLoginUrl)
	r.V1NonAuthGet("/auth/verify-google", authHdr.VerifyGoogleLogin)
	r.V1NonAuthGet("/auth/patient-login", authHdr.PatientLogin)
	r.V1NonAuthPost("/auth/refresh", authHdr.RefreshToken)

}
