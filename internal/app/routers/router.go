package routers

import (
	"context"
	"html/template"

	"github.com/peanut-cc/goBlog/pkg/logger"

	utils2 "github.com/peanut-cc/goBlog/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/peanut-cc/goBlog/internal/app/config"
	"github.com/peanut-cc/goBlog/internal/app/controller"
	"github.com/peanut-cc/goBlog/internal/app/utils"
)

func NewRouter(ctx context.Context) *gin.Engine {
	gin.SetMode(config.C.Server.RunMode)
	app := gin.New()
	// 匹配模版
	controller.Tmpl = template.New("goBlog").Funcs(utils.TplFuncMap)

	files := utils2.ReadDir("./frontend/template", func(name string) bool {
		if name == ".DS_Store" {
			return true
		}
		return false
	})

	_, err := controller.Tmpl.ParseFiles(files...)
	if err != nil {
		logger.Fatalf(ctx, "tmpl parse file error:%v", err)
	}

	app.Static("/static", "./frontend/static")
	app.GET("/admin/login", controller.HandleLogin)
	app.GET("/", controller.Index)
	return app
}
