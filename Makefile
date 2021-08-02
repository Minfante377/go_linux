-include .env
-include .secrets

VERSION := $(shell git rev-parse --short HEAD)
PROJECT_NAME=go_exec

GOPATH := $(PWD)
PKGS := $(shell ls src)

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Debug=$(DEBUG) -X=main.LogDir=$(LOG_DIR) -X=main.Pass=$(PASS) -X=main.Scripts=$(SCRIPTS) -X=main.TelegramToken=$(TELEGRAM_TOKEN) -X=main.User=$(USER)"

install:
	go get -d ./...

build:
	@GOPATH=$(GOPATH) go build $(LDFLAGS) -o bin/$(PROJECT_NAME) src/main.go

run:
	@GOPATH=$(GOPATH) go run $(LDFLAGS) src/main.go

test:
	@GOPATH=$(GOPATH) go test -v cmd_helper db_helper

clean:
	rm -rf bin/

all:
	clean
	install
	build
