.PHONY: help dev build run test test-coverage clean docker-up docker-down migrate swagger lint format install

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo '$(BLUE)Available targets:$(NC)'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

install: ## Install dependencies
	@echo "$(BLUE)Installing dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

dev: ## Run the application in development mode
	@echo "$(BLUE)Starting development server...$(NC)"
	go run main.go

build: ## Build the application binary
	@echo "$(BLUE)Building application...$(NC)"
	go build -o bin/todos main.go
	@echo "$(GREEN)✓ Build complete: bin/todos$(NC)"

run: build ## Build and run the application
	@echo "$(BLUE)Running application...$(NC)"
	./bin/todos

test: ## Run tests
	@echo "$(BLUE)Running tests...$(NC)"
	go test -v ./tests/...

test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	go test -v -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

clean: ## Clean build artifacts and caches
	@echo "$(YELLOW)Cleaning...$(NC)"
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out coverage.html
	go clean -cache -testcache
	@echo "$(GREEN)✓ Cleaned$(NC)"

docker-up: ## Start Docker containers (database + API)
	@echo "$(BLUE)Starting Docker containers...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)✓ Docker containers started$(NC)"

docker-down: ## Stop Docker containers
	@echo "$(YELLOW)Stopping Docker containers...$(NC)"
	docker-compose down
	@echo "$(GREEN)✓ Docker containers stopped$(NC)"

docker-logs: ## View Docker container logs
	docker-compose logs -f

docker-rebuild: ## Rebuild and restart Docker containers
	@echo "$(BLUE)Rebuilding Docker containers...$(NC)"
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d
	@echo "$(GREEN)✓ Docker containers rebuilt and started$(NC)"

swagger: ## Generate Swagger documentation
	@echo "$(BLUE)Generating Swagger docs...$(NC)"
	swag init
	@echo "$(GREEN)✓ Swagger docs generated$(NC)"

lint: ## Run linter (golangci-lint)
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --timeout 5m; \
	else \
		echo "$(RED)golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

format: ## Format code with gofmt and goimports
	@echo "$(BLUE)Formatting code...$(NC)"
	go fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	else \
		echo "$(YELLOW)goimports not installed. Install with: go install golang.org/x/tools/cmd/goimports@latest$(NC)"; \
	fi
	@echo "$(GREEN)✓ Code formatted$(NC)"

migrate-up: ## Run database migrations (if using migrate tool)
	@echo "$(BLUE)Running migrations...$(NC)"
	@echo "$(YELLOW)Note: This project uses GORM AutoMigrate. Migrations run automatically on startup.$(NC)"

db-reset: ## Reset database (WARNING: deletes all data)
	@echo "$(RED)⚠ WARNING: This will delete all data in the database!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		docker-compose up -d postgres; \
		sleep 2; \
		echo "$(GREEN)✓ Database reset complete$(NC)"; \
	else \
		echo "$(YELLOW)Cancelled$(NC)"; \
	fi

watch: ## Watch for changes and auto-reload (requires air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(RED)air not installed. Install with: go install github.com/cosmtrek/air@latest$(NC)"; \
		echo "$(BLUE)Falling back to: make dev$(NC)"; \
		make dev; \
	fi

env-setup: ## Copy .env.example to .env
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(GREEN)✓ .env file created from .env.example$(NC)"; \
		echo "$(YELLOW)⚠ Please update .env with your configuration$(NC)"; \
	else \
		echo "$(YELLOW).env file already exists$(NC)"; \
	fi

# Default target
.DEFAULT_GOAL := help
