.PHONY: build test clean static_analysis lint vet fmt chkfmt

APP         = ubume
GO          = go
GO_BUILD    = $(GO) build
GO_FORMAT   = $(GO) fmt
GOFMT       = gofmt
GO_LIST     = $(GO) list
GO_TEST     = $(GO) test -v
GO_VET      = $(GO) vet
GO_DEP      = $(GO) mod
GO_LDFLAGS  = -ldflags="-s -w"
GOOS        = linux
GO_PKGROOT  = ./...
GO_PACKAGES = $(shell $(GO_LIST) $(GO_PKGROOT))

build:  ## Build binary 
	env GO111MODULE=on GOOS=$(GOOS) $(GO_BUILD) $(GO_LDFLAGS) -o $(APP) main.go

clean: ## Clean project
	-rm -rf $(APP)

test: ## Start the test
	env GOOS=$(GOOS) $(GO_TEST) $(GO_PKGROOT)

vet: ## Start go vet
	$(GO_VET) $(GO_PACKAGES)

fmt: ## Format go source code 
	$(GO_FORMAT) $(GO_PKGROOT)

.DEFAULT_GOAL := help
help:  
	@grep -E '^[0-9a-zA-Z_-]+[[:blank:]]*:.*?## .*$$' $(MAKEFILE_LIST) | sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1;32m%-15s\033[0m %s\n", $$1, $$2}'