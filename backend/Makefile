# aitestos Makefile
.PHONY: all build test run clean lint fmt vet dev docker-up docker-down

APP_NAME := aitestos
BUILD_DIR := bin
MAIN_PATH := ./cmd/server

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOVET := $(GOCMD) vet
GOFMT := gofmt
GOLINT := golangci-lint

# Default target
all: fmt vet lint test build

# Build the application
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Run unit tests
test:
	@echo "Running unit tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./tests/integration/...

# Run the application
run: build
	@echo "Running server..."
	./$(BUILD_DIR)/$(APP_NAME)

# Development mode with dependencies
dev:
	@echo "Starting development environment..."
	docker-compose up -d
	@sleep 5
	$(MAKE) run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out

# Run linter
lint:
	@echo "Running linter..."
	$(GOLINT) run ./...

# Format code
fmt:
	@echo "Formatting..."
	$(GOFMT) -s -w .

# Run go vet
vet:
	@echo "Running vet..."
	$(GOVET) ./...

# Tidy go modules
tidy:
	@echo "Tidying modules..."
	$(GOCMD) mod tidy

# Generate code (wire, mocks, etc.)
generate:
	@echo "Running go generate..."
	$(GOCMD) generate ./...

# Database migrations
migrate-up:
	@echo "Running migrations..."
	@./scripts/migrate.sh up

migrate-down:
	@echo "Rolling back migrations..."
	@./scripts/migrate.sh down

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/vektra/mockery/v2@v2.42.0

# Check for issues
check: fmt vet lint test
	@echo "All checks passed!"

# Help
help:
	@echo "Available targets:"
	@echo "  all            - Format, vet, lint, test and build"
	@echo "  build          - Build the application"
	@echo "  test           - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  run            - Build and run the application"
	@echo "  dev            - Start dev environment with Docker"
	@echo "  clean          - Remove build artifacts"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format code with gofmt"
	@echo "  vet            - Run go vet"
	@echo "  tidy           - Run go mod tidy"
	@echo "  generate       - Run go generate"
	@echo "  migrate-up     - Run database migrations"
	@echo "  migrate-down   - Rollback database migrations"
	@echo "  docker-up      - Start Docker containers"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  docker-logs    - View Docker logs"
	@echo "  install-tools  - Install development tools"
	@echo "  check          - Run all checks"
