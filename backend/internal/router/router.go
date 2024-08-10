package router

import (
	"fmt"
	"time"

	"github.com/bookpanda/cas-sso/backend/config"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
	V1 *gin.RouterGroup
}

func New(conf *config.Config, corsHandler config.CorsHandler) *Router {
	if !conf.App.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.HandlerFunc(corsHandler))
	v1 := r.Group("/api/v1")
	v1.Use(gin.LoggerWithFormatter(logRequest))

	return &Router{r, v1}
}

func logRequest(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[%s] \"%d %s %s %s %s\"\n",
		param.TimeStamp.Format(time.RFC1123),
		param.StatusCode,
		param.Method,
		param.Path,
		param.Latency,
		param.ErrorMessage,
	)
}
