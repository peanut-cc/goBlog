package controller

import (
	"bytes"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/peanut-cc/goBlog/internal/app/config"
	"github.com/peanut-cc/goBlog/internal/app/ent"

	"github.com/peanut-cc/goBlog/internal/app/global"

	"github.com/peanut-cc/goBlog/internal/app/iutils"
	"github.com/peanut-cc/goBlog/pkg/logger"

	ipost "github.com/peanut-cc/goBlog/internal/app/ent/post"
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
		c.Redirect(http.StatusFound, "/admin/profile")
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
		logger.StartSpan(c, logger.SetSpanFuncName("HandleLoginPost")).Errorf("用户名或密码参数错误:username=%v,password=%s", user, pwd)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
	admin, err := global.EntClient.User.Query().Where(iuser.UsernameEQ(user)).Only(c)
	if err != nil {
		logger.StartSpan(c, logger.SetSpanFuncName("HandleLoginPost")).Errorf("ent orm query user name is: %v error:%v", user, err)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
	if !iutils.VerifyPasswd(admin.Password, user, pwd) {
		logger.StartSpan(c, logger.SetSpanFuncName("HandleLoginPost")).Errorf("用户名或密码参数错误:username=%v,password=%s", user, pwd)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
	session := sessions.Default(c)
	session.Set("username", user)
	err = session.Save()
	_, err = admin.Update().SetLoginTime(time.Now()).Save(c)
	if err != nil {
		logger.StartSpan(c, logger.SetSpanFuncName("HandleLoginPost")).Warnf("admin user update login time error: %v", err)
	}
	c.Redirect(http.StatusFound, "/admin/profile")

}

func HandleProfile(c *gin.Context) {
	admin, err := global.EntClient.User.Query().Where(iuser.UsernameEQ(config.C.Blog.UserName)).Only(c)
	if err != nil {
		logger.Errorf(c, "ent orm query user name is: %v error:%v", config.C.Blog.UserName, err)
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}
	blogInfo := global.EntClient.Blog.Query().FirstX(c)

	c.Status(http.StatusOK)
	RenderHTMLBack(c, "admin-profile", gin.H{
		"Console":  true,
		"Path":     c.Request.URL.Path,
		"Title":    "个人配置 | " + config.C.Blog.BTitle,
		"Account":  admin,
		"BlogInfo": blogInfo,
	})
}

func HandlePost(c *gin.Context) {
	h := gin.H{}
	blogInfo, err := global.EntClient.Blog.Query().First(c)
	if err != nil {
		logger.Errorf(c, "ent orm query blog info error:%v", err.Error())
		c.Redirect(http.StatusFound, "/admin/manage-posts")
		return
	}
	id, err := strconv.Atoi(c.Query("cid"))
	if err == nil && id > 0 {
		post, err2 := global.EntClient.Post.Query().WithCategory().Where(ipost.IDEQ(id)).Only(c)
		if err2 != nil {
			// logger.Warnf(c, "not found post err:%v", err2.Error())
			logger.Errorf(c, "ent orm query post info error:%v", err.Error())
			c.Redirect(http.StatusFound, "/admin/manage-posts")
			return
		}
		h["Title"] = "编辑文章 | " + blogInfo.Btitle
		h["Edit"] = post
		var postTags []string
		for _, tag := range post.QueryTags().AllX(c) {
			postTags = append(postTags, tag.Name)
		}
		h["PostTags"] = postTags
	}
	if h["Title"] == nil {
		h["Title"] = "撰写文章 | " + blogInfo.Btitle
	}
	tags, err := global.EntClient.Tag.Query().All(c)
	if err != nil {
		logger.Errorf(c, "ent orm query tag error:%v", err.Error())
		c.Redirect(http.StatusFound, "/admin/manage-posts")
		return
	}
	h["Tags"] = tags
	categories, err := global.EntClient.Category.Query().All(c)
	if err != nil {
		logger.Errorf(c, "ent orm query categories error:%v", err.Error())
		c.Redirect(http.StatusFound, "/admin/manage-posts")
		return
	}

	h["Categories"] = categories
	h["Path"] = c.Request.URL.Path
	h["Domain"] = config.C.Server.Domain
	c.Status(http.StatusOK)
	RenderHTMLBack(c, "admin-post", h)
}

// post list
func HandlePosts(c *gin.Context) {
	h := gin.H{}
	blogInfo, err := global.EntClient.Blog.Query().First(c)
	tmp := c.Query("serie")
	se, err := strconv.Atoi(tmp)
	if err != nil {
		logger.Warnf(c, "error:%v", err)
		// logger.Errorf(c, "serie args error:%v", err.Error())
		// c.Redirect(http.StatusFound, "/admin/profile")
		// return
	}
	pg, err := strconv.Atoi(c.Query("page"))
	if err != nil || pg < 1 {
		pg = 1
	}
	allPosts, err := global.EntClient.Post.Query().WithCategory().All(c)
	if err != nil {
		logger.Errorf(c, "query posts error:%v", err.Error())
		c.Redirect(http.StatusFound, "/admin/manage-posts")
		return
	}
	pagination := &Pagination{
		CurrentPage: pg,
		PerPage:     5,
		Total:       len(allPosts),
	}
	start := (pg - 1) * 5
	var end int
	if start+5 > len(allPosts) {
		end = len(allPosts)
	} else {
		end = start + 5
	}

	perPosts := allPosts[start:end]
	h["Console"] = true
	h["Path"] = c.Request.URL.Path
	h["Title"] = "个人配置 | " + blogInfo.Btitle
	h["Serie"] = se
	h["Posts"] = perPosts
	h["PostCount"] = len(allPosts)
	h["Pagination"] = pagination
	c.Status(http.StatusOK)
	RenderHTMLBack(c, "admin-posts", h)
}

func HandleCategories(c *gin.Context) {
	blogInfo, err := global.EntClient.Blog.Query().First(c)
	if err != nil {
		logger.Errorf(c, "ent orm query blog info error:%v", err.Error())
		c.Redirect(http.StatusFound, "/admin/profile")
		return
	}
	h := gin.H{}
	h["Manage"] = true
	h["Path"] = c.Request.URL.Path
	h["Title"] = "专题管理 | " + blogInfo.Btitle
	categories, err := global.EntClient.Category.Query().All(c)
	if err != nil {
		logger.Errorf(c, "ent orm query category info error:%v", err.Error())
		c.Redirect(http.StatusFound, "/admin/profile")
		return
	}
	h["Categories"] = categories
	c.Status(http.StatusOK)
	RenderHTMLBack(c, "admin-series", h)
}

func HandleCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("mid"))
	blogInfo, err := global.EntClient.Blog.Query().First(c)
	if err != nil {
		logger.Errorf(c, "ent orm query blog info error:%v", err.Error())
		c.Redirect(http.StatusFound, "/admin/manage-series")
		return
	}
	h := gin.H{}
	category, err := global.EntClient.Category.Get(c, id)
	if ent.IsNotFound(err) {
		h["Title"] = "新增分类 | " + blogInfo.Btitle
	} else if err != nil {
		logger.Errorf(c, "ent orm query category error:%v", err.Error())
		c.Redirect(http.StatusFound, "/admin/manage-series")
		return
	} else {
		h["Title"] = "编辑分类 | " + blogInfo.Btitle
		h["Category"] = category
	}
	h["Manage"] = true
	h["Path"] = c.Request.URL.Path
	c.Status(http.StatusOK)
	RenderHTMLBack(c, "admin-serie", h)
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

// api
func HandleAPI(c *gin.Context) {
	action := c.Param("action")
	logger.Debugf(c, "handle api action is: %v", action)
	api := APIs[action]
	if api == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid API Request",
		})
		return
	}
	api(c)
}
