package iutils

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/peanut-cc/goBlog/internal/app/config"
	"github.com/peanut-cc/goBlog/pkg/logger"
	"github.com/peanut-cc/goBlog/pkg/utils"
)

// 验证密码
func VerifyPasswd(origin, name, input string) bool {
	return origin == utils.EncryptPasswd(name, input)
}

// admin is login
func IsLogin(c *gin.Context) bool {
	session := sessions.Default(c)
	v := session.Get("username")
	logger.Warnf(c, "%v", v)
	if v == nil || v.(string) != config.C.Blog.UserName {
		return false
	}

	return true
}
