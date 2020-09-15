package routers

import (
	"context"
	"crypto/rand"
	"html/template"

	"github.com/peanut-cc/goBlog/internal/app/middleware"

	"github.com/peanut-cc/goBlog/pkg/logger"

	utils2 "github.com/peanut-cc/goBlog/pkg/utils"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/peanut-cc/goBlog/internal/app/config"
	"github.com/peanut-cc/goBlog/internal/app/controller"
	"github.com/peanut-cc/goBlog/internal/app/iutils"
)

func NewRouter(ctx context.Context) *gin.Engine {
	gin.SetMode(config.C.Server.RunMode)
	app := gin.New()

	// use session
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		logger.Errorf(ctx, "rand read error:%v", err)
		panic(err)
	}
	store := sessions.NewCookieStore(b)
	store.Options(sessions.Options{
		MaxAge:   86400 * 7,
		Path:     "/",
		HttpOnly: true,
	})
	app.Use(sessions.Sessions("su", store))

	// 匹配模版
	controller.Tmpl = template.New("goBlog").Funcs(iutils.TplFuncMap)

	files := utils2.ReadDir("./frontend/template", func(name string) bool {
		if name == ".DS_Store" {
			return true
		}
		return false
	})

	_, err = controller.Tmpl.ParseFiles(files...)
	if err != nil {
		logger.Fatalf(ctx, "tmpl parse file error:%v", err)
		panic(err)
	}
	app.Static("/static", "./frontend/static")

	app.Use(filter())
	app.GET("/admin/login", controller.HandleLogin)
	app.POST("/admin/login", controller.HandleLoginPost)
	admin := app.Group("/admin")
	auth := admin.Use(middleware.AuthFilter())
	{
		auth.GET("/profile", controller.HandleProfile)
		auth.POST("/api/:action", controller.HandleAPI)
		auth.GET("/write-post", controller.HandlePost)
		auth.GET("/manage-posts", controller.HandlePosts)
	}

	return app
}

func filter() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("u")
		if err != nil || cookie == "" {
			c.SetCookie("u", utils2.MustUUID(), 86400*730, "/", "", true, true)
		}
		c.Next()
	}
}
