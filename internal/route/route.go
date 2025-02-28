package route

import (
	"book-management/internal/controller/ping"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type WebHandler interface {
	RegisterRoute(r *gin.Engine)
}

type MiddleWareSet []gin.HandlerFunc

type WebHandlerSet []WebHandler

var ProviderSet = wire.NewSet(NewGlobalMiddlewareSet,NewWebHandlerSet,NewRouter,

)

func NewGlobalMiddlewareSet() MiddleWareSet {
	var middleWareSet []gin.HandlerFunc
	return middleWareSet
}

func NewWebHandlerSet(pingCtrl *ping.PingController) WebHandlerSet {
	var webHandlerSet []WebHandler
	
	webHandlerSet = append(webHandlerSet,pingCtrl)

	return webHandlerSet
}


func NewRouter(globalMiddleWare MiddleWareSet,webHandler WebHandlerSet) *gin.Engine {
	r := gin.Default()
	if len(globalMiddleWare) > 0 {
		r.Use(globalMiddleWare...)
	}
	
	for _, handler := range webHandler {
		handler.RegisterRoute(r)
	}
	return r
}


// 定义一个空的 WebHandler 实现
type EmptyWebHandler struct{}

func (h *EmptyWebHandler) RegisterRoute(r *gin.Engine) {
    // 空实现
}


