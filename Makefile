# Makefile for UAS Backend

.PHONY: test test-verbose test-coverage test-service test-middleware test-integration build clean help

# Default target
help:
	@echo "Available targets:"
	@echo "  test              - Run all tests"
	@echo "  test-verbose      - Run all tests with verbose output"
	@echo "  test-coverage     - Run tests with coverage report"
	@echo "  test-service      - Run service layer tests only"
	@echo "  test-middleware   - Run middleware tests only"
	@echo "  test-integration  - Run integration tests only"
	@echo "  build             - Build the application"
	@echo "  clean             - Clean build artifacts"
	@echo "  run               - Run the application"

# Run all tests
test:
	go test ./tests/... -short

# Run all tests with verbose output
test-verbose:
	go test ./tests/... -v

# Run tests with coverage
test-coverage:
	go test ./tests/... -v -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run service layer tests only
test-service:
	go test ./tests/service/... -v

# Run middleware tests only
test-middleware:
	go test ./tests/middleware/... -v

# Run integration tests only
test-integration:
	go test ./tests/integration/... -v

# Build the application
build:
	go build -o bin/uas-backend .

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run the application
run:
	go run .

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Run tests in watch mode (requires entr)
test-watch:
	find . -name "*.go" | entr -r make test-verbose