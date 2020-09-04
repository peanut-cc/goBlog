package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/peanut-cc/goBlog/internal/app/config"
	"github.com/peanut-cc/goBlog/internal/app/controller"
)

func NewRouter() *gin.Engine {
	gin.SetMode(config.C.Server.RunMode)
	app := gin.New()
	app.GET("/", controller.Index)
	return app

}
