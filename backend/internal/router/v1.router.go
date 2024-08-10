package router

import (
	"github.com/bookpanda/cas-sso/backend/internal/context"
	"github.com/gin-gonic/gin"
)

func (r *Router) V1Get(path string, handler func(c context.Ctx)) {
	r.V1.GET(path, func(c *gin.Context) {
		handler(context.NewContext(c))
	})
}

func (r *Router) V1Post(path string, handler func(c context.Ctx)) {
	r.V1.POST(path, func(c *gin.Context) {
		handler(context.NewContext(c))
	})
}

func (r *Router) V1Put(path string, handler func(c context.Ctx)) {
	r.V1.PUT(path, func(c *gin.Context) {
		handler(context.NewContext(c))
	})
}

func (r *Router) V1Patch(path string, handler func(c context.Ctx)) {
	r.V1.PATCH(path, func(c *gin.Context) {
		handler(context.NewContext(c))
	})
}

func (r *Router) V1Delete(path string, handler func(c context.Ctx)) {
	r.V1.DELETE(path, func(c *gin.Context) {
		handler(context.NewContext(c))
	})
}
