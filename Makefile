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
PLUGIN_BINARY_NAME_1=wob-v0
PLUGIN_BINARY_UNIX=$(PLUGIN_BINARY_NAME)_unix
PLUGIN_FOLDER_1=./plugin-go-grpc
PLUGIN_FOLDER_2=./plugin-get-dom
PLUGIN_BINARY_NAME_2=wob-v1
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
	rm -rf $(PLUGIN_FOLDER_1)/$(PLUGIN_BINARY_NAME_1)
	rm -rf $(PLUGIN_FOLDER_1)/$(PLUGIN_BINARY_NAME_1).zip
	rm -rf $(PLUGIN_FOLDER_2)/$(PLUGIN_BINARY_NAME_2)
	rm -rf $(PLUGIN_FOLDER_2)/$(PLUGIN_BINARY_NAME_2).zip
run:
	export ENV_PLUGIN="./$(PLUGIN_BINARY_NAME)"
	./$(CONTROLLER_BINARY_NAME)

plugin-v1:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(PLUGIN_FOLDER_1)/$(PLUGIN_BINARY_NAME_1) -v $(PLUGIN_FOLDER_1)

plugin-v2:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(PLUGIN_FOLDER_2)/$(PLUGIN_BINARY_NAME_2) -v $(PLUGIN_FOLDER_2)

install:
	$(GOBUILD) -o $(GOPATH)/bin/$(CONTROLLER_BINARY_NAME)

docker:
	#env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(PLUGIN_FOLDER_1)/$(PLUGIN_BINARY_NAME_1) -v $(PLUGIN_FOLDER_1)
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(PLUGIN_FOLDER_2)/$(PLUGIN_BINARY_NAME_2) -v $(PLUGIN_FOLDER_2)
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(CONTROLLER_BINARY_NAME) -v
	#jiminy zip plugin-go-grpc/
	jiminy zip plugin-get-dom/
	docker build . -t sibeshkar/jiminy-env:$(VERSION) --force-rm

docker-d:
	#env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(PLUGIN_FOLDER_1)/$(PLUGIN_BINARY_NAME_1) -v $(PLUGIN_FOLDER_1)
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(PLUGIN_FOLDER_2)/$(PLUGIN_BINARY_NAME_2) -v $(PLUGIN_FOLDER_2)
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(CONTROLLER_BINARY_NAME) -v
	jiminy zip plugin-go-grpc/
	jiminy zip plugin-get-dom/
	docker build . -t sibeshkar/jiminy-env:detached -f Dockerfile.detached --force-rm

docker-run:
	docker run -it --rm -p 5901:5900 -p 15901:15900 --memory="500m" -p 84:6080 sibeshkar/jiminy-env:$(VERSION)

docker-run-d:
	docker run -it --rm -p 5901:5900 -p 15901:15900 --memory="500m" -p 84:6080 sibeshkar/jiminy-env:detached

docker-record:
	docker run -it --rm -p 5901:5901 -p 15901:15900 --memory="500m" -p 84:6080 sibeshkar/jiminy-env:$(VERSION)

docker-push:
	docker push sibeshkar/jiminy-env:$(VERSION)
	#docker push sibeshkar/jiminy-env:detached


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
