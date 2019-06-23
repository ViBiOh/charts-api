SHELL = /bin/sh

APP_NAME = eponae-api
PACKAGES ?= ./...
GO_FILES ?= */*.go */*/*.go

GOBIN=bin
BINARY_PATH=$(GOBIN)/$(APP_NAME)

SERVER_SOURCE = cmd/api.go
SERVER_RUNNER = go run $(SERVER_SOURCE)
ifeq ($(DEBUG), true)
	SERVER_RUNNER = dlv debug $(SERVER_SOURCE) --
endif

.DEFAULT_GOAL := app

## help: Display list of commands
.PHONY: help
help: Makefile
	@sed -n 's|^##||p' $< | column -t -s ':' | sed -e 's|^| |'

## name: Output name
.PHONY: name
name:
	@echo -n $(APP_NAME)

## dist: Output binary path
.PHONY: dist
dist:
	@echo -n $(BINARY_PATH)

## version: Output sha1 of last commit
.PHONY: version
version:
	@echo -n $(shell git rev-parse --short HEAD)

## author: Output author's name of last commit
.PHONY: author
author:
	@python -c 'import sys; import urllib; sys.stdout.write(urllib.quote_plus(sys.argv[1]))' "$(shell git log --pretty=format:'%an' -n 1)"

## app: Build app with dependencies download
.PHONY: app
app: deps go

## go: Build app
.PHONY: go
go: format lint test bench build

## deps: Download dependencies
.PHONY: deps
deps:
	go get github.com/kisielk/errcheck
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports

## format: Format code
.PHONY: format
format:
	goimports -w $(GO_FILES)
	gofmt -s -w $(GO_FILES)

## lint: Lint code
.PHONY: lint
lint:
	golint $(PACKAGES)
	errcheck -ignoretests $(PACKAGES)
	go vet $(PACKAGES)

## test: Test with coverage
.PHONY: test
test:
	script/coverage

## bench: Benchmark code
.PHONY: bench
bench:
	go test $(PACKAGES) -bench . -benchmem -run Benchmark.*

## build: Build binary
.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(BINARY_PATH) $(SERVER_SOURCE)

## start: Start app
.PHONY: start
start:
	$(SERVER_RUNNER) \
		-csp "default-src 'self'; base-uri 'self'; script-src 'self' 'unsafe-inline' unpkg.com/swagger-ui-dist@3/; style-src 'self' 'unsafe-inline' unpkg.com/swagger-ui-dist@3/; img-src 'self' data:; connect-src 'self' api.eponae.fr" \
		-dbHost ${EPONAE_DATABASE_HOST} \
		-dbUser ${EPONAE_DATABASE_USER} \
		-dbPass ${EPONAE_DATABASE_PASS} \
		-dbName ${EPONAE_DATABASE_NAME}
