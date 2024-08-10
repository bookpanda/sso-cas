package main

import (
	"fmt"

	"github.com/bookpanda/cas-sso/backend/config"
	"github.com/bookpanda/cas-sso/backend/database"
	"github.com/bookpanda/cas-sso/backend/internal/validator"
	"github.com/bookpanda/cas-sso/backend/logger"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	_, err = database.InitDatabase(&conf.Db, conf.App.IsDevelopment())
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	_ = logger.New(&conf.App)
	_, err = validator.NewDtoValidator()
	if err != nil {
		panic(fmt.Sprintf("Failed to create validator: %v", err))
	}

	// corsHandler := config.MakeCorsConfig(conf)
	// r := router.New(conf, corsHandler, authMiddleware)

	// if err := r.Run(fmt.Sprintf(":%v", conf.App.Port)); err != nil {
	// 	log.Fatal("unable to start server")
	// }

	// r.V1NonAuthGet("/auth/google-url", authHdr.GetGoogleLoginUrl)
	// r.V1NonAuthGet("/auth/verify-google", authHdr.VerifyGoogleLogin)
	// r.V1NonAuthGet("/auth/patient-login", authHdr.PatientLogin)
	// r.V1NonAuthPost("/auth/refresh", authHdr.RefreshToken)

}
