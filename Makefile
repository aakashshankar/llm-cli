# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=llm

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

all: build
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -rf bin/*

.PHONY: test
test:
	$(GOTEST) -v ./test/...