PROJECT_NAME := "social-auth-service"

PKG := "."
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

all: build

test: ## Run unittests
	@go test -short ${PKG_LIST}

build: ## Build non statically
	@go build -o $(PROJECT_NAME) -i -v $(PKG)

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)