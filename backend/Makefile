.PHONY: help build run test clean docker-build docker-run docker-stop

# Variables
APP_NAME=protein-analyzer
DOCKER_IMAGE=protein-analyzer-backend
DOCKER_TAG=latest
GO_FILES=$(shell find . -name '*.go' -type f)

# Default target
help:
	@echo "Available commands:"
	@echo "  make build        - Build the Go binary"
	@echo "  make run          - Run the application locally"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run with docker-compose"
	@echo "  make docker-stop  - Stop docker-compose services"
	@echo "  make setup        - Setup development environment"

# Build the Go binary
build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) .
	@echo "Build complete: bin/$(APP_NAME)"

# Run the application locally
run:
	@echo "Starting $(APP_NAME)..."
	@go run .

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf uploads/*
	@echo "Clean complete"

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@go mod download
	@go mod tidy
	@mkdir -p uploads
	@cp .env.example .env 2>/dev/null || true
	@echo "Installing Python dependencies..."
	@pip3 install biopython
	@echo "Setup complete"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run:
	@echo "Starting services with docker-compose..."
	@docker-compose up -d
	@echo "Services started. Backend available at http://localhost:8080"

docker-stop:
	@echo "Stopping services..."
	@docker-compose down
	@echo "Services stopped"

docker-logs:
	@docker-compose logs -f backend

# Development with hot reload
dev:
	@echo "Starting development server with hot reload..."
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	@air

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w $(GO_FILES)
	@go mod tidy

# Lint code
lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run

# Generate mocks for testing
mocks:
	@echo "Generating mocks..."
	@which mockgen > /dev/null || go install github.com/golang/mock/mockgen@latest
	@go generate ./...