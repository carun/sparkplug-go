# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=sparkplug-publisher
PROTO_FILES=$(wildcard proto/*.proto)

# Build information
VERSION?=1.0.0
BUILD_DIR=build
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

.PHONY: all build clean test coverage deps proto run help

all: clean deps proto build test ## Run clean, deps, proto, build, and test

build: ## Build the binary
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/main.go

clean: ## Clean build directory
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f proto/*.pb.*

test: ## Run tests
	$(GOTEST) -v ./...

coverage: ## Run tests with coverage
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

deps: ## Get dependencies
	$(GOMOD) download
	$(GOMOD) tidy

proto: $(PROTO_FILES) ## Generate protobuf files
	protoc --go_out=. --go_opt=paths=source_relative $(PROTO_FILES)

run: build ## Run the application
	./$(BUILD_DIR)/$(BINARY_NAME)

lint: ## Run linters
	golangci-lint run

vet: ## Run go vet
	$(GOCMD) vet ./...

fmt: ## Run go fmt
	$(GOCMD) fmt ./...

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
