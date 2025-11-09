.PHONY: help build test lint lint-fix tidy clean examples tool-verify

# Default target
help:
	@echo "Brow - Browser automation CLI tool"
	@echo ""
	@echo "Build & Run:"
	@echo "  make build           Build the brow binary"
	@echo "  make tidy            Tidy go modules"
	@echo ""
	@echo "Testing:"
	@echo "  make test            Run unit tests"
	@echo "  make examples        Test example scripts"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint            Run golangci-lint"
	@echo "  make lint-fix        Run golangci-lint with auto-fix"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean           Clean build artifacts"
	@echo "  make tool-verify     Verify tool dependencies"

# Build the application
build:
	@echo "Building brow..."
	@go build -o bin/brow .
	@echo "Build complete: ./bin/brow"

# Run unit tests
test:
	@echo "Running unit tests..."
	@go test -v ./...

# Run example scripts to verify functionality
examples:
	@echo "Testing example scripts..."
	@cd examples && ./books.sh
	@echo "Examples completed successfully"

# Linting
lint:
	@echo "Running golangci-lint..."
	@$(shell go env GOPATH)/bin/golangci-lint run --config config/.golangci.yml

# Linting with auto-fix
lint-fix:
	@echo "Running golangci-lint with auto-fix..."
	@$(shell go env GOPATH)/bin/golangci-lint run --config config/.golangci.yml --fix

# Tidy go modules
tidy:
	@echo "Tidying go modules..."
	@go mod tidy
	@echo "Tidy complete"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f bin/brow
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Verify tool dependencies
tool-verify:
	@echo "Verifying tool dependencies..."
	@go mod verify
	@echo "✓ All tool dependencies verified"
	@echo ""
	@echo "Available tools (tracked in go.mod):"
	@echo "  - golangci-lint (code linting)"