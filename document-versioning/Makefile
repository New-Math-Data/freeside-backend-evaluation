# Document Versioning API - Makefile

.PHONY: help deps generate migrate-up migrate-down run build test docker-build docker-up docker-down clean

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

# Development setup
deps: ## Install development tools (sqlc, goose, oapi-codegen)
	@echo "Installing development tools..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo "Tools installed successfully."

# Code generation
generate: generate-sqlc generate-openapi ## Run all code generation

generate-sqlc: ## Generate Go code from SQL queries
	@echo "Generating sqlc code..."
	@sqlc generate -f migrations/sqlc.yaml
	@echo "sqlc generation complete."

generate-openapi: ## Generate Go server code from OpenAPI spec
	@echo "Generating OpenAPI server code..."
	@oapi-codegen -config openapi/oapi-codegen.yaml openapi/openapi.yaml
	@echo "OpenAPI generation complete."

# Database migrations
migrate-up: ## Apply database migrations
	@echo "Applying migrations..."
	@goose -dir migrations/migrations postgres "postgresql://postgres:postgres@localhost:5432/docstore?sslmode=disable" up
	@echo "Migrations applied."

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@goose -dir migrations/migrations postgres "postgresql://postgres:postgres@localhost:5432/docstore?sslmode=disable" down
	@echo "Migrations rolled back."

migrate-reset: ## Reset all database migrations
	@echo "Resetting migrations..."
	@goose -dir migrations/migrations postgres "postgresql://postgres:postgres@localhost:5432/docstore?sslmode=disable" reset
	@echo "Migrations reset."

# Local development
run: ## Run the server locally
	@echo "Starting server..."
	@go run ./cmd/server

build: ## Build the server binary
	@echo "Building server..."
	@go build -o bin/server ./cmd/server
	@echo "Build complete: bin/server"

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

# Docker commands
docker-build: ## Build the Docker image
	@echo "Building Docker image..."
	@docker compose build api

docker-up: ## Start all services (postgres + api)
	@echo "Starting services..."
	@docker compose up -d postgres
	@echo "Waiting for postgres to be healthy..."
	@sleep 3
	@$(MAKE) migrate-up
	@docker compose up -d api
	@echo "Services started. API available at http://localhost:8080"

docker-down: ## Stop all services
	@echo "Stopping services..."
	@docker compose down
	@echo "Services stopped."

docker-logs: ## View API logs
	@docker compose logs -f api

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "Clean complete."
