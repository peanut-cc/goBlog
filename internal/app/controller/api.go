package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peanut-cc/goBlog/internal/app/config"
	"github.com/peanut-cc/goBlog/internal/app/ent"
	"github.com/peanut-cc/goBlog/internal/app/ent/post"
	"github.com/peanut-cc/goBlog/internal/app/ent/tag"
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

const (
	SUCCESS = iota
	FAIL
)

// 全局API
var APIs = make(map[string]func(c *gin.Context))

func init() {
	// 更新帐号信息
	APIs["account"] = apiAccount
	APIs["blog"] = apiBlog
	APIs["password"] = apiPassword
	APIs["post-add"] = apiPostAdd
	APIs["serie-add"] = apiSerieAdd
	APIs["serie-delete"] = apiCategoryDelete
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
	if blogName == "" || bTitle == "" {
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
		Save(c)
	if err != nil {
		logger.Errorf(c, "blog info update error:%v", err.Error())
		responseNotice(c, NOTICE_NOTICE, err.Error(), "")
		return
	}
	global.BlogInfo.BlogName = blogName
	global.BlogInfo.BTitle = bTitle
	global.BlogInfo.BeiAn = beian
	global.BlogInfo.SubTitle = subTitle
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

func apiPostAdd(c *gin.Context) {
	var (
		err error
		do  string
		cid string
	)
	defer func() {
		switch do {
		case "auto":
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"fail": FAIL, "time": time.Now().Format("15:04:05 PM"), "cid": cid})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": SUCCESS, "time": time.Now().Format("15:04:05 PM"), "cid": cid})
		case "save", "publish":
			if err != nil {
				if err != nil {
					responseNotice(c, NOTICE_NOTICE, err.Error(), "")
					return
				}
			}
			uri := "/admin/manage-draft"
			if do == "publish" {
				uri = "/admin/profile"
			}
			c.Redirect(http.StatusFound, uri)
		}
	}()
	do = c.PostForm("do")
	slug := c.PostForm("slug")
	title := c.PostForm("title")
	text := c.PostForm("text")
	category := c.PostForm("serie")
	tag := c.PostForm("tags")
	update := c.PostForm("update")
	date := utils.CheckDate(c.PostForm("date"))
	cid = c.PostForm("cid")
	if slug == "" || title == "" || text == "" {
		err = errors.New("参数错误")
		return
	}
	var tags []string
	if tag != "" {
		tags = strings.Split(tag, ",")
	}

	oldPost, err := global.EntClient.Post.Query().Where(post.TitleEQ(title)).Only(c)
	if err == nil {
		// 已经存在的文章
		oldPost.Update().SetBody(text).Save(c)
		UpdateMultiTags(c, tags, oldPost.ID)
		if utils.CheckBool(update) {
			oldPost.Update().SetModifiedTime(time.Now()).Save(c)
		}
		logger.Warnf(c, "aa:%v", category)
		if category != "" {
			categoryID, err2 := strconv.Atoi(category)
			logger.Warnf(c, "bb:%v", categoryID)
			if err2 == nil {
				oldPost.Update().SetCategoryID(categoryID).Save(c)
			}
		}
		return
	}
	//  表示新文章
	if cid == "" || err != nil {
		newPost, err := global.EntClient.Post.Create().
			SetAuthor("peanut").
			SetBody(text).
			SetTitle(title).
			SetCreatedTime(date).
			SetIsDraft(do != "publish").
			Save(c)
		if err != nil {
			logger.StartSpan(c, logger.SetSpanFuncName("apiPostAdd")).Errorf("post create error:%v", err.Error())
		} else {
			UpdateMultiTags(c, tags, newPost.ID)
			if category != "" {
				categoryID, err := strconv.Atoi(category)
				if err == nil {
					newPost.Update().SetCategoryID(categoryID).Save(c)
				}
			}
		}
		return
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
			_tag, err := global.EntClient.Tag.Query().Where(tag.NameEQ(newTag)).Only(ctx)
			if ent.IsNotFound(err) {
				needToAddTag = append(needToAddTag, global.EntClient.Tag.Create().SetName(newTag))
			} else {
				post.Update().AddTags(_tag).Save(ctx)
			}
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

// 分类的新增和删除
func apiSerieAdd(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		responseNotice(c, NOTICE_NOTICE, "参数错误", "")
		return
	}
	mid, err := strconv.Atoi(c.PostForm("mid"))
	if err == nil && mid > 0 {
		// 更新分类名称
		category, err := global.EntClient.Category.Get(c, mid)
		if err != nil {
			responseNotice(c, NOTICE_NOTICE, "未找到数据", "")
			return
		}
		_, err = category.Update().SetName(name).Save(c)
		if err != nil {
			logger.Errorf(c, "category update name error:", err.Error())
			responseNotice(c, NOTICE_NOTICE, err.Error(), "")
			return
		}
	} else {
		// 新增分类
		_, err = global.EntClient.Category.Create().SetName(name).Save(c)
		if err != nil {
			logger.Errorf(c, "category create error:", err.Error())
			responseNotice(c, NOTICE_NOTICE, err.Error(), "")
			return
		}
	}
	c.Redirect(http.StatusFound, "/admin/manage-series")
	return
}

func apiCategoryDelete(c *gin.Context) {
	for _, v := range c.PostFormArray("mid[]") {
		id, err := strconv.Atoi(v)
		if err != nil || id < 1 {
			responseNotice(c, NOTICE_NOTICE, err.Error(), "")
			return
		}
		err = global.EntClient.Category.DeleteOneID(id).Exec(c)
		if err != nil {
			logger.Errorf(c, "api category delete error:", err.Error())
			responseNotice(c, NOTICE_NOTICE, err.Error(), "")
			return
		}
	}
	responseNotice(c, NOTICE_SUCCESS, "删除成功", "")
}

func responseNotice(c *gin.Context, typ, content, hl string) {
	if hl != "" {
		c.SetCookie("notice_highlight", hl, 86400, "/", "", true, false)
	}
	c.SetCookie("notice_type", typ, 86400, "/", "", true, false)
	c.SetCookie("notice", fmt.Sprintf("[\"%s\"]", content), 86400, "/", "", true, false)
	c.Redirect(http.StatusFound, c.Request.Referer())
}
