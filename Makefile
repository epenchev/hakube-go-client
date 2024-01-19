GO_BUILD = CGO_ENABLED=0 go build -trimpath -o hakube-client main.go

GO := $(shell which go)

default: build

all: build

.PHONY: fmt
fmt: ## Format source code.
	go fmt ./...

.PHONY: fmt build
build:
	$(GO_BUILD)
