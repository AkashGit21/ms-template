SHELL := /bin/bash

GOCMD=go
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GORUN=$(GOCMD) run

BINARY_NAME=ms-server
MAIN_DIR=./cmd/ms-project
TAGS=coprocess grpc goplugin

TEST_REGEX=.
TEST_COUNT=1

BENCH_REGEX=.
BENCH_RUN=NONE


.PHONY: build
build:
	$(GOBUILD) -tags "$(TAGS)" -o $(BINARY_NAME) -v $(MAIN_DIR)

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)


.PHONY: docker-build
docker-build:
	echo -e "\nBuilding docker Image... \n\nThis may take a few minutes!\n"
	docker build -t ms-template .

.PHONY: docker-run
docker-run:
	echo -e "\n\nRunning the container..."
	docker run -d --rm -p 8081:8081/tcp -p 8082:8082/udp --name $(BINARY_NAME) ms-template
	echo -e "\n\nContainer is running successfully in the background! "

.PHONY: docker-stop
docker-stop:
	echo -e "\n\nStopping the docker container..."
	docker stop $(BINARY_NAME)

.PHONY: gen
gen:
	protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=. --openapiv2_out=./docs/swagger --proto_path=. internal/proto-files/*.proto --experimental_allow_proto3_optional
	cp docs/swagger/internal/proto-files/*.json docs/swagger/
	rm -rf docs/swagger/internal

.PHONY: help
help:
	$(GORUN) $(MAIN_DIR) help

.PHONY: run
run:
	$(GORUN) $(MAIN_DIR) run

.PHONY: test
test:
	echo -e "\n\n Testing..."
	$(GOTEST) -v -race -cover --shuffle=on ./...