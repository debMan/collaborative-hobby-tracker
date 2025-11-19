.PHONY: help build run test test-unit test-integration clean install dev docker-up docker-down

# Default target
help:
	@echo "Available targets:"
	@echo "  install          - Install dependencies"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application"
	@echo "  dev              - Run in development mode with hot reload"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  clean            - Clean build artifacts"
	@echo "  docker-up        - Start Docker services (MongoDB)"
	@echo "  docker-down      - Stop Docker services"
	@echo "  lint             - Run linter"
	@echo "  fmt              - Format code"

# Install dependencies
install:
	go mod download
	go mod tidy

# Build the application
build:
	go build -o bin/api ./cmd/api

# Run the application
run: build
	./bin/api

# Run in development mode
dev:
	@which air > /dev/null || (echo "Installing air for hot reload..." && go install github.com/cosmtrek/air@latest)
	air

# Run all tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run unit tests only
test-unit:
	go test -v -race -short ./...

# Run integration tests only
test-integration:
	go test -v -race -run Integration ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Run linter
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Format code
fmt:
	go fmt ./...
	gofmt -s -w .

# Start Docker services
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# Start MongoDB only
mongo-up:
	docker run -d --name hobby-tracker-mongo -p 27017:27017 mongo:latest

# Stop MongoDB
mongo-down:
	docker stop hobby-tracker-mongo
	docker rm hobby-tracker-mongo
