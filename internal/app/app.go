package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/peanut-cc/goBlog/internal/app/ent/user"
	"github.com/peanut-cc/goBlog/pkg/utils"

	"github.com/peanut-cc/goBlog/internal/app/ent"
	"github.com/peanut-cc/goBlog/internal/app/ent/migrate"
	"github.com/peanut-cc/goBlog/internal/app/routers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/peanut-cc/goBlog/pkg/logger"

	"github.com/peanut-cc/goBlog/internal/app/config"
	"github.com/peanut-cc/goBlog/internal/app/global"
)

type options struct {
	ConfigFile string
	Version    string
}

// Option 定义配置项
type Option func(*options)

// SetConfigFile 设定配置文件
func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

// SetVersion 设定版本号
func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	config.MustLoad(o.ConfigFile)
	config.PrintWithJSON()

	logger.Printf(ctx, "服务启动，运行模式：%s，版本号：%s，进程号：%d", config.C.Server.RunMode, o.Version, os.Getpid())
	// 初始化日志模块
	loggerCleanFunc, err := InitLogger()
	if err != nil {
		return nil, err
	}

	router := routers.NewRouter(ctx)
	httpServerCleanFunc := InitHttpServer(ctx, router)

	entOrmCleanFunc, err := InitEntOrm()
	if err != nil {
		loggerCleanFunc()
		return nil, err
	}
	loadAdminUser(ctx)
	return func() {
		entOrmCleanFunc()
		loggerCleanFunc()
		httpServerCleanFunc()
	}, nil
}

// 初始化 http server
func InitHttpServer(ctx context.Context, handle http.Handler) func() {
	cfg := config.C.Server
	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: handle,
	}
	go func() {
		logger.Printf(ctx, "HTTP server is running at %s.", addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()
		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf(ctx, err.Error())
		}
	}
}

func InitEntOrm() (func(), error) {
	var err error
	cfg := config.C.MySQL
	global.EntClient, err = ent.Open("mysql", cfg.DSN())
	if err != nil {
		logger.Errorf(context.Background(), "Ent orm open db error:%v", err.Error())
		return nil, err
	}
	cleanFunc := func() {
		err := global.EntClient.Close()
		if err != nil {
			logger.Errorf(context.Background(), "Ent orm closed error:%v", err.Error())
		}
	}

	// run the auto migration tool
	err = global.EntClient.Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	if err != nil {
		logger.Errorf(context.Background(), "Ent orm create schema resources:%v", err)
		return cleanFunc, err
	}
	return cleanFunc, nil
}

// 初始化admin用户的信息到数据库
func loadAdminUser(ctx context.Context) {
	admin, err := global.EntClient.User.Query().Where(user.UsernameEQ(config.C.Blog.UserName)).Only(ctx)
	if err != nil {
		logger.Errorf(ctx, "Initializing admin user:%v", config.C.Blog.UserName)
		global.EntClient.User.Create().
			SetUsername(config.C.Blog.UserName).
			SetPassword(utils.EncryptPasswd(config.C.Blog.UserName, config.C.Blog.Password)).
			SetEmail(config.C.Blog.Email).
			SetPhone(config.C.Blog.Phone).SaveX(ctx)
	} else {
		logger.Printf(ctx, "admin user is %v", admin.Username)
	}

	// load blog info
	// TODO 关于这里的异常判断可以处理的更加合理
	_, err = global.EntClient.Blog.Query().First(ctx)
	if err != nil {
		logger.Errorf(ctx, "Query blog info err:%v ", err.Error())
		global.EntClient.Blog.Create().
			SetDefaultPageNum(config.C.Blog.DefaultPageNum).
			SetBlogName(config.C.Blog.BlogName).
			SetBtitle(config.C.Blog.BTitle).
			SetCopyRight(config.C.Blog.Copyright).
			SetSubtitle(config.C.Blog.SubTitle).
			SetBeian(config.C.Blog.BeiAn).
			Save(ctx)
	}
}

func Run(ctx context.Context, opts ...Option) error {
	var state int32 = 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return err
	}
EXIT:
	for {
		sig := <-sc
		logger.Printf(ctx, "接收到信号[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			atomic.CompareAndSwapInt32(&state, 1, 0)
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	logger.Printf(ctx, "服务退出")
	time.Sleep(time.Second)
	os.Exit(int(atomic.LoadInt32(&state)))
	return nil
}
