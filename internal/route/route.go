package route

import (
	"book-management/internal/controller"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type WebHandler interface {
	RegisterRoute(r *gin.Engine)
}

type MiddleWareSet []gin.HandlerFunc

type WebHandlerSet []WebHandler

var ProviderSet = wire.NewSet(NewGlobalMiddlewareSet, NewWebHandlerSet, NewRouter)

func NewGlobalMiddlewareSet() MiddleWareSet {
	var middleWareSet []gin.HandlerFunc
	return middleWareSet
}

func NewWebHandlerSet(pingCtrl *controller.PingController) WebHandlerSet {
	var webHandlerSet []WebHandler

	webHandlerSet = append(webHandlerSet, pingCtrl)

	return webHandlerSet
}

func NewRouter(globalMiddleWare MiddleWareSet, webHandler WebHandlerSet) *gin.Engine {
	r := gin.Default()
	if len(globalMiddleWare) > 0 {
		r.Use(globalMiddleWare...)
	}

	for _, handler := range webHandler {
		handler.RegisterRoute(r)
	}
	return r
}
