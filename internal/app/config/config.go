package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/peanut-cc/goBlog/pkg/utils"

	"github.com/koding/multiconfig"
)

var (
	// C 全局配置(需要先执行MustLoad，否则拿不到配置)
	C    = new(Config)
	once sync.Once
)

// MustLoad 加载配置
func MustLoad(fpaths ...string) {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{},
		}

		for _, fpath := range fpaths {
			if strings.HasSuffix(fpath, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "json") {
				loaders = append(loaders, &multiconfig.JSONLoader{Path: fpath})
			}
			if strings.HasSuffix(fpath, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: fpath})
			}
		}

		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}
		m.MustLoad(C)
	})
}

// PrintWithJSON 基于JSON格式输出配置
func PrintWithJSON() {
	if C.PrintConfig {
		b, err := utils.JSONMarshalIndent(C, "", " ")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

type Config struct {
	PrintConfig bool
	Server      Server
	Blog        Blog
	Log         Log
	MySQL       MySQL
}

type Server struct {
	RunMode         string
	HTTPPort        int
	ShutdownTimeout int
	Domain 			string
}

type Blog struct {
	DefaultPageNum int
	UserName       string
	Password       string
	Email          string
	Phone          string
	BlogName       string
	SubTitle       string
	BeiAn          string
	BTitle         string
	Copyright      string
}

type Log struct {
	Level      int
	Format     string
	Output     string
	OutputFile string
}

// MySQL mysql配置参数
type MySQL struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	Parameters string
}

// DSN 数据库连接串
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}
