.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

APP = blog
SERVER_BIN = ./cmd/${APP}/${APP}
RELEASE_ROOT = release
RELEASE_SERVER = release/${APP}


all: start

start:
	go run cmd/${APP}/main.go web -c ./configs/config.toml

