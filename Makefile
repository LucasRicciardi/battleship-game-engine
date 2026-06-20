# Battleship Game Engine Makefile

.PHONY: help build test lint run clean docker-build docker-run docker-stop

help: ## Show this help message
	@echo "Battleship Game Engine - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

build: ## Build the application
	@echo "Building Battleship Game Engine..."
	go build -o battleship-game-engine ./cmd/main.go

test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep -v "total:"
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	golangci-lint run ./...

lint-fix: ## Run linters with auto-fix
	@echo "Running linters with auto-fix..."
	golangci-lint run ./... --fix

run: build ## Run the application
	@echo "Running Battleship Game Engine..."
	./battleship-game-engine

clean: ## Clean build artifacts
	@echo "Cleaning..."
	go clean -cache -testcache
	rm -f battleship-game-engine
	rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t battleship-game-engine:latest .

docker-run: ## Run Docker container (requires docker-compose up)
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-stop: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs: ## Show Docker container logs
	@echo "Showing Docker logs..."
	docker-compose logs -f

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	go run ./cmd/migrate up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	go run ./cmd/migrate down

migrate-status: ## Show migration status
	@echo "Showing migration status..."
	go run ./cmd/migrate status

.PHONY: help build test test-coverage lint lint-fix run clean docker-build docker-run docker-stop docker-logs migrate-up migrate-down migrate-status
