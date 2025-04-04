package route

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type WebHandler interface {
	RegisterRoute(r *gin.Engine)
}

type MiddleWare gin.HandlerFunc

var ProviderSet = wire.NewSet(NewRouter)

func NewRouter(middlewares []MiddleWare, webHandlers []WebHandler) *gin.Engine {
	r := gin.Default()
	for _, middleware := range middlewares {
		r.Use(gin.HandlerFunc(middleware))
	}

	for _, handler := range webHandlers {
		handler.RegisterRoute(r)
	}
	return r
}
