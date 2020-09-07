package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peanut-cc/goBlog/internal/app/config"
	"github.com/peanut-cc/goBlog/internal/app/ent/user"
	"github.com/peanut-cc/goBlog/internal/app/global"
	"github.com/peanut-cc/goBlog/pkg/logger"
)

const (
	// 成功
	NOTICE_SUCCESS = "success"
	// 注意
	NOTICE_NOTICE = "notice"
	// 错误
	NOTICE_ERROR = "error"
)

// 全局API
var APIs = make(map[string]func(c *gin.Context))

func init() {
	// 更新帐号信息
	APIs["account"] = apiAccount
	APIs["blog"] = apiBlog
}

func apiAccount(c *gin.Context) {
	email := c.PostForm("email")
	phone := c.PostForm("phoneNumber")
	logger.Debugf(c, "email:%v phone:%v", email, phone)
	_, err := global.EntClient.User.Update().
		Where(
			user.UsernameEQ(config.C.Blog.UserName),
		).
		SetEmail(email).
		SetPhone(phone).Save(c)
	if err != nil {
		logger.Errorf(c, "update user info error:%v", err.Error())
		responseNotice(c, NOTICE_NOTICE, err.Error(), "")
		return
	}
	responseNotice(c, NOTICE_SUCCESS, "更新成功", "")
}

func apiBlog(c *gin.Context) {
	blogName := c.PostForm("blogName")
	bTitle := c.PostForm("bTitle")
	beian := c.PostForm("beiAn")
	subTitle := c.PostForm("subTitle")
	copyRight := c.PostForm("beiAn")
	if blogName == "" || bTitle == "" || copyRight == "" {
		responseNotice(c, NOTICE_NOTICE, "参数错误", "")
		return
	}
	blogInfo, err := global.EntClient.Blog.Query().First(c)
	if err != nil {
		logger.Errorf(c, "query blog info error:%v", err.Error())
		responseNotice(c, NOTICE_NOTICE, err.Error(), "")
		return
	}
	_, err = blogInfo.Update().
		SetBlogName(blogName).
		SetBeian(beian).SetBtitle(bTitle).
		SetSubtitle(subTitle).
		SetBeian(beian).
		Save(c)
	if err != nil {
		logger.Errorf(c, "blog info update error:%v", err.Error())
		responseNotice(c, NOTICE_NOTICE, err.Error(), "")
		return
	}
	responseNotice(c, NOTICE_NOTICE, "更新成功", "")
}

func responseNotice(c *gin.Context, typ, content, hl string) {
	if hl != "" {
		c.SetCookie("notice_highlight", hl, 86400, "/", "", true, false)
	}
	c.SetCookie("notice_type", typ, 86400, "/", "", true, false)
	c.SetCookie("notice", fmt.Sprintf("[\"%s\"]", content), 86400, "/", "", true, false)
	c.Redirect(http.StatusFound, c.Request.Referer())
}
