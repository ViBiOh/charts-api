SHELL = /bin/sh

APP_NAME ?= eponae-api
VERSION ?= $(shell git rev-parse --short HEAD)
AUTHOR ?= $(shell git log --pretty=format:'%an' -n 1)

PACKAGES ?= ./...
APP_PACKAGES = $(shell go list -e $(PACKAGES) | grep -v vendor | grep -v node_modules)

GOBIN=bin
BINARY_PATH=$(GOBIN)/$(APP_NAME)

SERVER_SOURCE = cmd/api.go
SERVER_RUNNER = go run $(SERVER_SOURCE)
ifeq ($(DEBUG), true)
	SERVER_RUNNER = dlv debug $(SERVER_SOURCE) --
endif

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
	@echo -n $(VERSION)

## author: Output author's name of last commit
.PHONY: author
author:
	@python -c 'import sys; import urllib; sys.stdout.write(urllib.quote_plus(sys.argv[1]))' "$(AUTHOR)"

## $(APP_NAME): Build app with dependencies download
.PHONY: $(APP_NAME)
$(APP_NAME): deps go

## go: Build app
.PHONY: go
go: format lint test bench build

## doc: Build doc
.PHONY: doc
doc:
	docker run -it --rm -v "${PWD}/doc:/doc" quay.io/bukalapak/snowboard html -o api.html api.apib

## deps: Download dependencies
.PHONY: deps
deps:
	go get github.com/golang/dep/cmd/dep
	go get github.com/kisielk/errcheck
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports
	dep ensure

## format: Format code
.PHONY: format
format:
	goimports -w */*.go */*/*.go
	gofmt -s -w */*.go */*/*.go

## lint: Lint code
.PHONY: lint
lint:
	golint $(APP_PACKAGES)
	errcheck -ignoretests $(APP_PACKAGES)
	go vet $(APP_PACKAGES)

## test: Test with coverage
.PHONY: test
test:
	script/coverage

## bench: Benchmark code
.PHONY: bench
bench:
	go test $(APP_PACKAGES) -bench . -benchmem -run Benchmark.*

## build: Build binary
.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o $(BINARY_PATH) $(SERVER_SOURCE)

## start: Start app
.PHONY: start
start:
	$(SERVER_RUNNER) \
		-dbHost ${EPONAE_DATABASE_HOST} \
		-dbUser ${EPONAE_DATABASE_USER} \
		-dbPass ${EPONAE_DATABASE_PASS} \
		-dbName ${EPONAE_DATABASE_NAME}
