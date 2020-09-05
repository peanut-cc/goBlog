package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peanut-cc/goBlog/internal/app/iutils"
)

func AuthFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !iutils.IsLogin(c) {
			c.Abort()
			c.Redirect(http.StatusFound, "/admin/login")
			return
		}
		c.Next()
	}
}
