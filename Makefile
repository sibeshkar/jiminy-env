# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
CONTROLLER_BINARY_NAME=jiminy
CONTROLLER_BINARY_UNIX=$(CONTROLLER_BINARY_NAME)_unix
PLUGIN_BINARY_NAME=env-go-grpc
PLUGIN_BINARY_UNIX=$(PLUGIN_BINARY_NAME)_unix
PLUGIN_FOLDER=./plugin-go-grpc

all:build
build:
	$(GOBUILD) -o $(CONTROLLER_BINARY_NAME) -v

proto:
	protoc -I proto/ proto/env.proto --go_out=plugins=grpc:proto/
test:
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(PLUGIN_BINARY_NAME)
run:
	export ENV_PLUGIN="./$(PLUGIN_BINARY_NAME)"
	./$(CONTROLLER_BINARY_NAME)

install:
	$(GOBUILD) -o $(GOPATH)/bin/$(CONTROLLER_BINARY_NAME)


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
