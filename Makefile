# Go parameters
GOCMD=go
GOOS=linux
GOARCH=amd64
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
CONTROLLER_BINARY_NAME=jiminy
CONTROLLER_BINARY_UNIX=$(CONTROLLER_BINARY_NAME)_unix
PLUGIN_BINARY_NAME=wob-v0
PLUGIN_BINARY_UNIX=$(PLUGIN_BINARY_NAME)_unix
PLUGIN_FOLDER=./plugin-go-grpc
VERSION=0.1.0

all:build
build:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(CONTROLLER_BINARY_NAME) -v

proto:
	protoc -I proto/ proto/env.proto --go_out=plugins=grpc:proto/
test:
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf $(CONTROLLER_BINARY_NAME)
	rm -rf $(PLUGIN_FOLDER)/$(PLUGIN_BINARY_NAME)
	rm -rf $(PLUGIN_FOLDER)/$(PLUGIN_BINARY_NAME).zip
run:
	export ENV_PLUGIN="./$(PLUGIN_BINARY_NAME)"
	./$(CONTROLLER_BINARY_NAME)

install:
	$(GOBUILD) -o $(GOPATH)/bin/$(CONTROLLER_BINARY_NAME)

docker:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(CONTROLLER_BINARY_NAME) -v
	jiminy zip plugin-go-grpc/
	docker build . -t sibeshkar/jiminy-env:$(VERSION) --force-rm

docker-run:
	docker run -it --rm -p 5901:5901 -p 15901:15900 sibeshkar/jiminy-env:$(VERSION)

docker-run-d:
	docker run -it --rm -p 5901:5901 -p 15901:15900 sibeshkar/jiminy-env:detached

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
