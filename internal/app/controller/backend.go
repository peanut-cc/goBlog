package controller

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleLogin(c *gin.Context) {
	c.Status(http.StatusOK)
	RenderHTMLBack(c, "login.html", gin.H{
		"BTitle": "GoBlog",
	})
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
