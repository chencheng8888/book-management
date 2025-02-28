package main

import (
	"book-management/configs"
	"book-management/pkg/logger"
	"flag"

	"github.com/gin-gonic/gin"
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

func (a *App) Run()  {
	if err :=a.engine.Run(port);err!=nil {
		logger.LogPrinter.Info("gin endgine run error:",err)
		return
	}
}

func init() {
	flag.StringVar(&flagConf, "conf", "configs/config.yaml", "config file path")
}

func main() {
	//解析命令行参数
	flag.Parse()


	if err := configs.LoadConfigs(flagConf);err!=nil {
		panic(err)
	}

	app := InitializeApp()
	app.Run()
}