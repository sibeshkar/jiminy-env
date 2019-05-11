
GOCMD=go
GOBUILD=$(GOCMD) build
ENV_PLUGIN_BINARY_NAME=wob
TAG=v0
TASK_CONFIG=tasks.json

all:build
build:
	$(GOBUILD) -o $(ENV_PLUGIN_BINARY_NAME)-$(TAG) -v
install:
	jiminy install .


