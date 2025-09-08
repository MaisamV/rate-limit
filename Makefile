# Go Clean Architecture Template - Makefile
# This file provides convenient commands for development and deployment tasks

# Variables
GO_VERSION := 1.24.4
APP_NAME := go-clean-template
MIGRATION_DIR := ./scripts/migrations
DATABASE_URL ?= postgres://postgres:password@localhost:5432/go_clean_db?sslmode=disable

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)Available commands:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'

# =============================================================================
# Development Commands
# =============================================================================

.PHONY: dev
dev: ## Start development environment with docker-compose
	@echo "$(BLUE)Starting development environment...$(RESET)"
	docker-compose up -d

.PHONY: dev-down
dev-down: ## Stop development environment
	@echo "$(BLUE)Stopping development environment...$(RESET)"
	docker-compose down

.PHONY: dev-logs
dev-logs: ## Show development environment logs
	docker-compose logs -f

# =============================================================================
# Build Commands
# =============================================================================

.PHONY: build
build: ## Build the application
	@echo "$(BLUE)Building application...$(RESET)"
	go build -o bin/$(APP_NAME) ./cmd/app

.PHONY: run
run: ## Run the application locally
	@echo "$(BLUE)Running application...$(RESET)"
	go run ./cmd/app

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	rm -rf bin/
	go clean

# =============================================================================
# Testing Commands
# =============================================================================

.PHONY: test
test: ## Run all tests
	@echo "$(BLUE)Running tests...$(RESET)"
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "$(BLUE)Running tests with coverage...$(RESET)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(RESET)"

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "$(BLUE)Running integration tests...$(RESET)"
	go test -v -tags=integration ./test/...

# =============================================================================
# Database Migration Commands
# =============================================================================

.PHONY: migrate-install
migrate-install: ## Install golang-migrate tool
	@echo "$(BLUE)Installing golang-migrate...$(RESET)"
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "$(YELLOW)Installing migrate tool...$(RESET)"; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	else \
		echo "$(GREEN)migrate tool already installed$(RESET)"; \
	fi

.PHONY: migrate-up
migrate-up: migrate-install ## Apply all pending migrations
	@echo "$(BLUE)Applying migrations...$(RESET)"
	migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" up
	@echo "$(GREEN)Migrations applied successfully$(RESET)"

.PHONY: migrate-down
migrate-down: migrate-install ## Rollback last migration
	@echo "$(YELLOW)Rolling back last migration...$(RESET)"
	@read -p "Are you sure you want to rollback the last migration? [y/N]: " confirm && \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" down 1; \
		echo "$(GREEN)Migration rolled back successfully$(RESET)"; \
	else \
		echo "$(YELLOW)Migration rollback cancelled$(RESET)"; \
	fi

.PHONY: migrate-status
migrate-status: migrate-install ## Check migration status
	@echo "$(BLUE)Checking migration status...$(RESET)"
	migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" version

.PHONY: migrate-force
migrate-force: migrate-install ## Force migration to specific version (use with caution)
	@echo "$(RED)WARNING: This will force the migration version without running migrations!$(RESET)"
	@read -p "Enter version number to force to: " version && \
	read -p "Are you absolutely sure? This can corrupt your database! [y/N]: " confirm && \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" force $$version; \
		echo "$(GREEN)Migration version forced to $$version$(RESET)"; \
	else \
		echo "$(YELLOW)Force migration cancelled$(RESET)"; \
	fi

.PHONY: migrate-create
migrate-create: migrate-install ## Create new migration files (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME is required. Usage: make migrate-create NAME=migration_name$(RESET)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Creating new migration: $(NAME)$(RESET)"
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)
	@echo "$(GREEN)Migration files created successfully$(RESET)"

.PHONY: migrate-reset
migrate-reset: migrate-install ## Reset database (drop all tables and reapply migrations)
	@echo "$(RED)WARNING: This will drop all tables and data!$(RESET)"
	@read -p "Are you absolutely sure? This will destroy all data! [y/N]: " confirm && \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" drop; \
		migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" up; \
		echo "$(GREEN)Database reset successfully$(RESET)"; \
	else \
		echo "$(YELLOW)Database reset cancelled$(RESET)"; \
	fi

# =============================================================================
# Code Quality Commands
# =============================================================================

.PHONY: lint
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(RESET)"; \
	fi

.PHONY: fmt
fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(RESET)"
	go fmt ./...
	goimports -w .

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(RESET)"
	go vet ./...

# =============================================================================
# Docker Commands
# =============================================================================

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(RESET)"
	docker build -t $(APP_NAME):latest .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "$(BLUE)Running Docker container...$(RESET)"
	docker run --rm -p 8080:8080 $(APP_NAME):latest

# =============================================================================
# Dependency Management
# =============================================================================

.PHONY: deps
deps: ## Download dependencies
	@echo "$(BLUE)Downloading dependencies...$(RESET)"
	go mod download

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "$(BLUE)Updating dependencies...$(RESET)"
	go get -u ./...
	go mod tidy

.PHONY: deps-verify
deps-verify: ## Verify dependencies
	@echo "$(BLUE)Verifying dependencies...$(RESET)"
	go mod verify

# =============================================================================
# Default target
# =============================================================================

.DEFAULT_GOAL := help