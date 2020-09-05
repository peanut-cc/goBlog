package controller

import (
	"bytes"
	"html/template"
	"net/http"
	"time"

	"github.com/peanut-cc/goBlog/internal/app/global"

	"github.com/peanut-cc/goBlog/internal/app/iutils"
	"github.com/peanut-cc/goBlog/pkg/logger"

	iuser "github.com/peanut-cc/goBlog/internal/app/ent/user"

	"github.com/gin-gonic/contrib/sessions"

	"github.com/gin-gonic/gin"
)

func HandleLogin(c *gin.Context) {
	logout := c.Query("logout")
	if logout == "true" {
		session := sessions.Default(c)
		session.Delete("username")
		session.Save()
	} else if iutils.IsLogin(c) {
		c.Redirect(http.StatusFound, "/admin/index")
		return
	}
	c.Status(http.StatusOK)
	RenderHTMLBack(c, "login.html", gin.H{
		"BTitle": "GoBlog",
	})
}

func HandleLoginPost(c *gin.Context) {
	user := c.PostForm("user")
	pwd := c.PostForm("password")
	if user == "" || pwd == "" {
		logger.Errorf(c, "用户名或密码参数错误:username=%v,password=%s", user, pwd)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
	admin, err := global.EntClient.User.Query().Where(iuser.UsernameEQ(user)).Only(c)
	if err != nil {
		logger.Errorf(c, "ent orm query user name is: %v error:%v", user, err)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
	if !iutils.VerifyPasswd(admin.Password, user, pwd) {
		logger.Errorf(c, "用户名或密码参数错误:username=%v,password=%s", user, pwd)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
	session := sessions.Default(c)
	session.Set("username", user)
	err = session.Save()
	logger.Warnf(c, "---%v", err)
	_, err = admin.Update().SetLoginTime(time.Now()).Save(c)
	if err != nil {
		logger.Warnf(c, "admin user update login time error: %v", err)
	}
	c.Redirect(http.StatusFound, "/admin/index")

}

// 渲染 html
func RenderHTMLBack(c *gin.Context, name string, data gin.H) {
	if name == "login.html" {
		err := Tmpl.ExecuteTemplate(c.Writer, name, data)
		if err != nil {
			panic(err)
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		return
	}
	var buf bytes.Buffer
	err := Tmpl.ExecuteTemplate(&buf, name, data)
	if err != nil {
		panic(err)
	}
	data["LayoutContent"] = template.HTML(buf.String())
	err = Tmpl.ExecuteTemplate(c.Writer, "backLayout.html", data)
	if err != nil {
		panic(err)
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
}
