APP_NAME    = ossia
DATE      ?= $(shell date +%FT%T%z)
VERSION   ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)
BIN        = $(GOPATH)/bin
BASE       = $(GOPATH)/src/$(APP_NAME)
BUILD_ARGS = $(shell env GOOS=linux GOARCH=amd64) 
PKGS       = $(or $(PKG),$(shell cd $(BASE) && env GOPATH=$(GOPATH) $(GO) list ./... | grep -v "^$(APP_NAME)/vendor/"))
TESTPKGS   = $(shell env GOPATH=$(GOPATH) $(GO) list -f '{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' $(PKGS))
CUR_DIR    = $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GO         = go
GOGET      = $(GO) get -u
GODOC      = godoc
GOFMT      = gofmt
SRC        = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

DEP        = $(BIN)/dep
SWAGGER    = $(BIN)/swagger
BINDATA    = $(BIN)/go-bindata
VERSION    = $(shell awk '$$2 ~ /Version/ {print $$3}' main.go)

.DEFAULT_GOAL := help

all:

help:
			@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps:
			$(GOGET) github.com/golang/dep/cmd/dep
			$(GOGET) github.com/shuLhan/go-bindata/...

vendor:
			$(DEP) ensure

swagger:
			$(SWAGGER) generate spec -m -o ./assets/swagger.json

bindata:
			$(BINDATA) -pkg application -o application/bindata.go ./assets/...

prebuild: deps vendor bindata swagger
			$(MAKE) bindata

build: prebuild ## Build Linux binary (amd64)
			GOOS=linux GOARCH=amd64 $(GO) build -o build/$(APP_NAME)

package: ## Build Debian & RedHat packages (amd64)
			docker build -f Dockerfile -t $(APP_NAME) $(CUR_DIR) ; docker run --rm -i -v $(CUR_DIR):/go/src/$(APP_NAME) $(APP_NAME) --package --version=$(VERSION)

run: prebuild  ## Start OSSIA
			$(GO) run main.go

clean: ## Clean build arftifacts
			-rm -rf build
			-rm -rf vendor/*

.PHONY: all vendor swagger bindata prebuild build debian \
			run test clean
