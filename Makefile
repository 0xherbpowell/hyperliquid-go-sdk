# Define the shell to use when executing commands
SHELL := /bin/bash

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Binary names
BINARY_DIR=bin
EXAMPLES_DIR=examples

.PHONY: help build clean test deps fmt lint vet tidy examples run-basic-order run-market-order

help: ## Display this help screen
	@grep -h '^[a-zA-Z]' $(MAKEFILE_LIST) | awk -F ':.*?## ' 'NF==2 {printf "   %-22s%s\n", $$1, $$2}' | sort

build: ## Build all examples
	@echo "Building examples..."
	@mkdir -p $(BINARY_DIR)
	@cd $(EXAMPLES_DIR)/basic_order && $(GOBUILD) -o ../../$(BINARY_DIR)/basic_order .
	@cd $(EXAMPLES_DIR)/basic_market_order && $(GOBUILD) -o ../../$(BINARY_DIR)/basic_market_order .
	@echo "Build complete!"

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)

test: ## Run tests
	$(GOTEST) -v ./...

test-cover: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) verify

tidy: ## Tidy go modules
	$(GOMOD) tidy

fmt: ## Format Go code
	$(GOFMT) -s -w .

lint: ## Run linter
	$(GOLINT) run ./...

vet: ## Run go vet
	$(GOCMD) vet ./...

check: fmt vet lint ## Run all checks (format, vet, lint)

examples: build ## Build all examples (alias for build)

run-basic-order: ## Run basic order example (requires config.json)
	@if [ ! -f examples/config.json ]; then \
		echo "Error: examples/config.json not found. Please copy examples/config.json.example to examples/config.json and fill in your details."; \
		exit 1; \
	fi
	@cd examples/basic_order && go run main.go

run-market-order: ## Run market order example (requires config.json)
	@if [ ! -f examples/config.json ]; then \
		echo "Error: examples/config.json not found. Please copy examples/config.json.example to examples/config.json and fill in your details."; \
		exit 1; \
	fi
	@cd examples/basic_market_order && go run main.go

install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

setup-config: ## Setup configuration file from example
	@if [ ! -f examples/config.json ]; then \
		cp examples/config.json.example examples/config.json; \
		echo "Configuration file created at examples/config.json"; \
		echo "Please edit this file and add your secret key and account details."; \
	else \
		echo "Configuration file already exists at examples/config.json"; \
	fi

all: clean deps tidy fmt vet lint test build ## Run full build pipeline

.DEFAULT_GOAL := help