package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/peanut-cc/goBlog/internal/app/config"
	"github.com/peanut-cc/goBlog/internal/app/ent"
	"github.com/peanut-cc/goBlog/internal/app/ent/post"
	"github.com/peanut-cc/goBlog/internal/app/ent/user"
	"github.com/peanut-cc/goBlog/internal/app/global"
	"github.com/peanut-cc/goBlog/internal/app/iutils"
	"github.com/peanut-cc/goBlog/pkg/logger"
	"github.com/peanut-cc/goBlog/pkg/utils"
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
	APIs["password"] = apiPassword
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

func apiPassword(c *gin.Context) {
	old := c.PostForm("old")
	new := c.PostForm("new")
	confirm := c.PostForm("confirm")

	if new != confirm {
		responseNotice(c, NOTICE_NOTICE, "两次密码输入不一致", "")
		return
	}
	admin, err := global.EntClient.User.Query().Where(user.UsernameEQ(config.C.Blog.UserName)).Only(c)
	if err != nil {
		responseNotice(c, NOTICE_NOTICE, "未找到用户", "")
		return
	}
	if !iutils.VerifyPasswd(admin.Password, admin.Username, old) {
		responseNotice(c, NOTICE_NOTICE, "原始密码不正确", "")
		return
	}
	newPwd := utils.EncryptPasswd(admin.Username, new)
	_, err = admin.Update().SetPassword(newPwd).Save(c)
	if err != nil {
		logger.StartSpan(c, logger.SetSpanFuncName("apiPassword")).Errorf("admin update password error:%v", err.Error())
		responseNotice(c, NOTICE_NOTICE, err.Error(), "")
		return
	}
	responseNotice(c, NOTICE_SUCCESS, "更新成功", "")
}

func responseNotice(c *gin.Context, typ, content, hl string) {
	if hl != "" {
		c.SetCookie("notice_highlight", hl, 86400, "/", "", true, false)
	}
	c.SetCookie("notice_type", typ, 86400, "/", "", true, false)
	c.SetCookie("notice", fmt.Sprintf("[\"%s\"]", content), 86400, "/", "", true, false)
	c.Redirect(http.StatusFound, c.Request.Referer())
}

func apiPostAdd(c *gin.Context) {
	var (
		err error
		do  string
		cid int
	)
	do = c.PostForm("do")
	slug := c.PostForm("slug")
	title := c.PostForm("title")
	text := c.PostForm("text")
	category := c.PostForm("serie")
	tag := c.PostForm("tags")
	update := c.PostForm("update")
	date := utils.CheckDate(c.PostForm("date"))
	if slug == "" || title == "" || text == "" {
		err = errors.New("参数错误")
		return
	}
	var tags []string
	if tag != "" {
		tags = strings.Split(tag, ",")
	}
	cid, err = strconv.Atoi(c.PostForm("cid"))
	//  表示新文章
	if err != nil || cid < 1 {
		global.EntClient.Post.Create().
			SetAuthor("peanut").
			SetBody(text).
			SetTitle(title).
			SetCreatedTime(date).
			Save(c)
	}

}

func UpdateMultiTags(ctx context.Context, newTags []string, postID int) {
	post, err := global.EntClient.Post.Query().Where(post.IDEQ(postID)).WithTags().First(ctx)
	if err != nil {
		logger.StartSpan(ctx, logger.SetSpanFuncName("UpdateMultiTags")).Fatalf("query post error:%v", err.Error())
	}
	var needToDelTagID []int
	for _, originTag := range post.Edges.Tags {
		if !IsInArray(originTag.Name, newTags) {
			needToDelTagID = append(needToDelTagID, originTag.ID)
		}
	}
	post.Update().RemoveTagIDs(needToDelTagID...).Save(ctx)
	var oldTagNames []string
	for _, oldTag := range post.Edges.Tags {
		oldTagNames = append(oldTagNames, oldTag.Name)
	}
	var needToAddTag []*ent.TagCreate
	for _, newTag := range newTags {
		if !IsInArray(newTag, oldTagNames) {
			needToAddTag = append(needToAddTag, global.EntClient.Tag.Create().SetName(newTag))
		}
	}
	addTags, err := global.EntClient.Tag.CreateBulk(needToAddTag...).Save(ctx)
	if err != nil {
		logger.StartSpan(ctx, logger.SetSpanFuncName("UpdateMultiTags")).Fatalf("tag create build error:%v", err.Error())
	}
	post.Update().AddTags(addTags...).Save(ctx)
}

// 判断tag是否在array中
func IsInArray(name string, tagNameArray []string) bool {
	for _, v := range tagNameArray {
		if name == v {
			return true
		}
	}
	return false
}
