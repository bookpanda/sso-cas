package config

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CorsHandler gin.HandlerFunc

func makeCorsConfig(conf *Config) gin.HandlerFunc {
	if conf.App.IsDevelopment() {
		return cors.New(cors.Config{
			AllowMethods:     []string{"*"},
			AllowHeaders:     []string{"*", "content-type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
			AllowOriginFunc: func(string) bool {
				return true
			},
		})

	}

	allowOrigins := strings.Split(conf.Cors.AllowOrigins, ",")

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func MakeCorsConfig(conf *Config) CorsHandler {
	return CorsHandler(makeCorsConfig(conf))
}
