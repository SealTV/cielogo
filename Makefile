.PHONY: help test coverage test-integration test-short lint fmt build clean install-tools ci vet

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "Available targets:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
	@echo ""

## test: Run all unit tests with race detector
test:
	@echo "Running unit tests..."
	go test -v -race ./...

## test-short: Run unit tests in short mode (skip slow tests)
test-short:
	@echo "Running unit tests in short mode..."
	go test -v -short ./...

## coverage: Run tests with coverage report
coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo ""
	@echo "Coverage summary:"
	@go tool cover -func=coverage.out | grep total
	@echo ""
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## test-integration: Run integration tests (requires CIELO_API_KEY)
test-integration:
	@echo "Running integration tests..."
	@if [ -f .env.test ]; then \
		export $$(cat .env.test | xargs) && go test -v -tags=integration ./...; \
	else \
		echo "Warning: .env.test file not found. Using environment variables."; \
		go test -v -tags=integration ./...; \
	fi

## lint: Run golangci-lint
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run 'make install-tools' to install it."; \
		exit 1; \
	fi

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## fmt: Format code and tidy modules
fmt:
	@echo "Formatting code..."
	gofmt -s -w .
	@echo "Tidying modules..."
	go mod tidy

## build: Build all packages
build:
	@echo "Building packages..."
	go build -v ./...

## clean: Clean build artifacts and test cache
clean:
	@echo "Cleaning..."
	rm -f coverage.out coverage.html
	go clean -cache -testcache -modcache

## install-tools: Install development tools
install-tools:
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		echo "golangci-lint installed successfully"; \
	else \
		echo "golangci-lint is already installed"; \
	fi

## ci: Run all checks (like CI pipeline)
ci: fmt vet lint test
	@echo ""
	@echo "================================"
	@echo "All checks passed successfully!"
	@echo "================================"
	@echo ""
