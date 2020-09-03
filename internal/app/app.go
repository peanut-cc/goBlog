package app

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/peanut-cc/goBlog/pkg/logger"

	"github.com/peanut-cc/goBlog/internal/app/config"
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
	return func() {
		loggerCleanFunc()
	}, nil
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
