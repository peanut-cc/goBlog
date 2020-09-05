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

	"github.com/peanut-cc/goBlog/internal/app/ent"

	"github.com/peanut-cc/goBlog/internal/app/routers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/peanut-cc/goBlog/pkg/logger"

	"github.com/peanut-cc/goBlog/internal/app/config"
)

var (
	entClient *ent.Client
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
	cfg := config.C.MySQL
	entClient, err := ent.Open("mysql", cfg.DSN())
	if err != nil {
		logger.Errorf(context.Background(), "Ent orm open db error:%v", err.Error())
		return nil, err
	}
	cleanFunc := func() {
		err := entClient.Close()
		if err != nil {
			logger.Errorf(context.Background(), "Ent orm closed error:%v", err.Error())
		}
	}

	// run the auto migration tool
	err = entClient.Schema.Create(context.Background())
	if err != nil {
		logger.Errorf(context.Background(), "Ent orm create schema resources:%v", err)
		return cleanFunc, err
	}
	return cleanFunc, nil
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
