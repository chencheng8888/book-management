package ping

import "github.com/gin-gonic/gin"

var pingController *PingController

func init() {
	pingController = NewPingController()
}



type PingController struct{}

func NewPingController() *PingController {
	return &PingController{}
}

func (p *PingController) Ping(ctx *gin.Context) {
	ctx.JSON(200,gin.H{
		"msg":"pong",
	})
}

func (p *PingController) RegisterRoute(r *gin.Engine) {
	r.GET("/ping",pingController.Ping)
}
