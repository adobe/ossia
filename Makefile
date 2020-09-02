APP_NAME    = ossia
VERSION    = $(shell awk '$$2 ~ /Version/ {print $$3}' main.go)
DATE      ?= $(shell date +%FT%T%z)
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

SWAGGER    = $(BIN)/swagger
BINDATA    = $(BIN)/go-bindata

.DEFAULT_GOAL := help

all:

help:
			@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps:
			$(GOGET) github.com/go-swagger/go-swagger/cmd/swagger
			$(GOGET) github.com/go-bindata/go-bindata/v3/go-bindata

swagger:
			$(SWAGGER) generate spec -m -o ./assets/swagger.json

bindata:
			$(BINDATA) -pkg application -o application/bindata.go -fs ./assets/...

prebuild: bindata swagger
			$(MAKE) bindata

build: prebuild ## Build Linux binary (amd64)
			GOOS=linux GOARCH=amd64 $(GO) build -o build/$(APP_NAME)

package: ## Build Debian & RedHat packages (amd64)
			docker build -f Dockerfile -t $(APP_NAME) $(CUR_DIR) ; docker run --rm -i -v $(CUR_DIR):/go/src/$(APP_NAME) $(APP_NAME) --package --version=$(VERSION)

run: prebuild  ## Start OSSIA
			$(GO) run main.go

clean: ## Clean build arftifacts
			-rm -rf build
			-rm -f ossia

.PHONY: all swagger bindata prebuild build package \
			run clean
