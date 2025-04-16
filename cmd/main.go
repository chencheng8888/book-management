package main

import (
	"flag"

	"github.com/gin-gonic/gin"

	"book-management/configs"
	"book-management/pkg/logger"
)

var (
	flagConf string

	port = ":8989"
)

type App struct {
	engine *gin.Engine
}

func newApp(engine *gin.Engine) *App {
	return &App{engine: engine}
}

func (a *App) Run() {
	if err := a.engine.Run(port); err != nil {
		logger.LogPrinter.Info("gin endgine run error:", err)
		return
	}
}

func init() {
	flag.StringVar(&flagConf, "conf", "configs/config.yaml", "config file path")
}

// @title Book Management API
// @version 1.0
// @description This is a sample server for a book management system.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8989
// @BasePath /api/v1
func main() {
	//解析命令行参数
	flag.Parse()

	var appConf configs.AppConfig
	if err := appConf.LoadConfig(flagConf); err != nil {
		logger.LogPrinter.Error("load config error:", err)
		return
	}
	logger.LogPrinter.Info(appConf)

	app, err := InitializeApp(appConf)
	if err != nil {
		panic(err)
	}
	app.Run()
}
