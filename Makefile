SHELL := /bin/bash

GOCMD=go
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

BINARY_NAME=ms-server
BINARY_LINUX=tyk
TAGS=coprocess grpc goplugin
CONF=tyk.conf

TEST_REGEX=.
TEST_COUNT=1

BENCH_REGEX=.
BENCH_RUN=NONE

.PHONY: build
build:
	$(GOBUILD) -tags "$(TAGS)" -o $(BINARY_NAME) -v .

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)


.PHONY: gen
gen: 
	protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=. --proto_path=. internal/proto-files/*.proto
	

.PHONY: run
run:
	go run main.go -config ./lib/configuration/config.json
