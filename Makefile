.PHONY: help build test bench run clean install release

# Variables
BINARY_NAME=logcleaner
GO=go
GOFLAGS=-v

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) cmd/logcleaner/main.go

test: ## Run tests
	$(GO) test $(GOFLAGS) ./...

test-coverage: ## Run tests with coverage
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

bench: ## Run benchmarks
	$(GO) test -bench=. -benchmem ./internal/cleaner/

run: ## Run the application
	$(GO) run cmd/logcleaner/main.go

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -rf dist/

install: ## Install the binary to $GOPATH/bin
	$(GO) install cmd/logcleaner/main.go

lint: ## Run linter
	golangci-lint run ./...

fmt: ## Format code
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

mod-tidy: ## Tidy go modules
	$(GO) mod tidy

mod-download: ## Download go modules
	$(GO) mod download

release-test: ## Test release build with goreleaser
	goreleaser release --snapshot --clean

dev: ## Run in development mode with live reload (requires air)
	air

# Example targets
example-basic: build ## Run basic example
	./$(BINARY_NAME)

example-test: build ## Test with example log file
	@echo "Testing with examples/logs/sample.log"
	@mkdir -p /tmp/logcleaner-test
	@cp examples/logs/sample.log /tmp/logcleaner-test/
	@echo "Sample log file copied to /tmp/logcleaner-test/sample.log"

.DEFAULT_GOAL := help
