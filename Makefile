-include .env

VERSION := $(shell git rev-parse --short HEAD)
PROJECT_NAME=go_exec

GOPATH := $(PWD)

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.LogDir=$(LOG_DIR)"

build:
	@GOPATH=$(GOPATH) go build $(LDFLAGS) -o bin/$(PROJECT_NAME) src/main.go

run:
	@GOPATH=$(GOPATH) go run $(LDFLAGS) src/main.go

clean:
	rm -rf bin/

all:
	clean
	build
