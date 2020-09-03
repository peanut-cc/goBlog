package main

import (
	"context"
	"os"

	"github.com/peanut-cc/goBlog/internal/app"

	"github.com/peanut-cc/goBlog/pkg/logger"
	"github.com/urfave/cli/v2"
)

// VERSION 版本号，可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "0.0.1"

func main() {
	logger.SetVersion(VERSION)
	ctx := logger.NewTraceIDContext(context.Background(), "main")

	app := cli.NewApp()
	app.Name = "goBlog"
	app.Version = VERSION
	app.Usage = "GIN + Ent ORM Blog"
	app.Commands = []*cli.Command{
		newWebCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Errorf(ctx, err.Error())
	}
}

func newWebCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "运行web服务",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "配置文件(.json,.yaml,.toml)",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			return app.Run(ctx,
				app.SetConfigFile(c.String("conf")))
		},
	}
}
