.PHONY: help build run test lint clean install

# Default target
help:
	@echo "Pokedex CLI - Available Commands"
	@echo "================================="
	@echo ""
	@echo "  make build      - Build the binary"
	@echo "  make run        - Run the application"
	@echo "  make test       - Run all tests"
	@echo "  make lint       - Run linters (golangci-lint)"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make install    - Install binary to \$$GOPATH/bin"
	@echo "  make help       - Show this help message"
	@echo ""

# Build the binary
build:
	@echo "Building Pokedex CLI..."
	go build -o ./bin/pokedex ./cmd/pokedex

# Run the application
run: build
	./bin/pokedex

# Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run ./... || go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f ./bin/pokedex
	go clean
	go mod tidy

# Install binary to $GOPATH/bin
install: build
	@echo "Installing Pokedex CLI..."
	go install ./cmd/pokedex
