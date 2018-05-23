default: api docker

api: deps go

go: format lint tst bench build

docker: docker-deps docker-build

deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u golang.org/x/tools/cmd/goimports
	dep ensure

format:
	goimports -w */*.go */*/*.go
	gofmt -s -w */*.go */*/*.go

lint:
	golint `go list ./... | grep -v vendor`
	errcheck -ignoretests `go list ./... | grep -v vendor`
	go vet ./...

tst:
	script/coverage

bench:
	go test ./... -bench . -benchmem -run Benchmark.*

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o bin/api cmd/api.go

docker-deps:
	curl -s -o cacert.pem https://curl.haxx.se/ca/cacert.pem

docker-build:
	docker run -it --rm -v `pwd`/doc:/doc bukalapak/snowboard html -o api.html api.apib
	docker build -t $(DOCKER_USER)/eponae-api .

docker-push:
	echo $(DOCKER_PASS) | docker login -u $(DOCKER_USER) --password-stdin
	docker push $(DOCKER_USER)/eponae-api

start-deps:
	go get -u github.com/ViBiOh/auth/cmd/bcrypt

start-api:
	go run -race cmd/api.go \
		-tls=false \
		-readingsAuthUsers admin:admin \
		-readingsBasicUsers 1:admin:`bcrypt password`

.PHONY: go dev docker deps format lint tst bench build docker-deps docker-build docker-push start-deps start-api
