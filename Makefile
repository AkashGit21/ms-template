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
	$(GOBUILD) -tags "$(TAGS)" -o $(BINARY_NAME) -v ./cmd/ms-project

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)


.PHONY: docker-run
docker-run:
	echo -e "\nBuilding docker Image... \n\nThis may take a few minutes!\n"
	docker build -t ms-template .

	echo -e "\n\nRunning the container..."
	docker run -d --rm -p 8081:8081/tcp -p 8082:8082/udp ms-template
	echo -e "\n\nContainer is running successfully in the background! "


.PHONY: gen
gen:
	protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=. --openapiv2_out=:swagger --proto_path=. internal/proto-files/*.proto --experimental_allow_proto3_optional


.PHONY: help
help:
	go run ./cmd/ms-project help

.PHONY: run
run:
	go run ./cmd/ms-project run
