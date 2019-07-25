
GOCMD=go
GOBUILD=$(GOCMD) build
ENV_PLUGIN_BINARY_NAME=wob
TAG=v0
TASK_CONFIG=tasks.json

all:build
build:
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(ENV_PLUGIN_BINARY_NAME)-$(TAG) -v
install:
	$(GOBUILD) -o $(ENV_PLUGIN_BINARY_NAME)-$(TAG) -v
	jiminy install .


