# Packs Calculator - Makefile
BACKEND_DIR := backend

.PHONY: help setup-dev docker-up docker-down docker-dev dev dev-frontend dev-backend test build clean migrate-up migrate-down migrate-create

help:
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  %-12s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

setup-dev: ## Complete development setup (includes starting all services)
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/air-verse/air@latest
	cd $(BACKEND_DIR) && swag init -g cmd/server/main.go -o ./docs	
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Installing backend dependencies..."
	cd $(BACKEND_DIR) && go mod download && go mod tidy 
	@echo "Installing frontend dependencies..."
	cd frontend && pnpm install
	@echo "Generating Swagger docs..."
	cd $(BACKEND_DIR) && go build -o app cmd/server/main.go
	@echo "Starting all services..."
	docker compose up -d --build
	@echo "✅ Development setup complete!"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:5173"
	@echo "Swagger: http://localhost:8080/swagger/index.html"

docker-up: ## Start all services (database + backend)
	docker compose up -d
	@echo "✅ Services started!"
	@echo "Backend: http://localhost:8080"
	@echo "Swagger: http://localhost:8080/swagger/index.html"

docker-down:
	docker compose down

docker-dev: ## Switch back to full containerized development
	@echo "🐳 Starting backend container..."
	docker compose up -d backend frontend
	@echo "✅ Switched to containerized development!"
	@echo "Backend: http://localhost:8080"

dev: ## Run in development mode with hot-reload (run setup-dev first)
	@echo "Note: Run 'make setup-dev' first to start database and install dependencies"
	@echo "🛑 Stopping backend and frontend containers to avoid port conflicts..."
	docker compose stop backend frontend || true
	@echo "📖 Generating Swagger docs..."
	cd $(BACKEND_DIR) && swag init -g cmd/server/main.go -o ./docs
	@echo "🔥 Starting backend with hot-reload..."
	cd $(BACKEND_DIR) && air &
	@echo "🎨 Starting frontend with hot-reload..."
	cd frontend && pnpm dev

dev-frontend: ## Run only frontend in development mode with hot-reload
	@echo "🎨 Starting frontend with hot-reload..."
	cd frontend && pnpm dev

dev-backend: ## Run only backend in development mode with hot-reload
	@echo "Note: Run 'make setup-dev' first to start database and install dependencies"
	@echo "🛑 Stopping backend container to avoid port conflicts..."
	docker compose stop backend || true
	@echo "📖 Generating Swagger docs..."
	cd $(BACKEND_DIR) && swag init -g cmd/server/main.go -o ./docs
	@echo "🔥 Starting backend with hot-reload..."
	cd $(BACKEND_DIR) && air

migrate-up: ## Run database migrations
	cd $(BACKEND_DIR) && migrate -path migrations -database "postgres://$${POSTGRES_USER:-packer}:$${POSTGRES_PASSWORD:-secret}@$${DB_HOST:-localhost}:$${DB_PORT:-5432}/$${POSTGRES_DB:-packs}?sslmode=disable" up

migrate-down: ## Rollback last migration
	cd $(BACKEND_DIR) && migrate -path migrations -database "postgres://$${POSTGRES_USER:-packer}:$${POSTGRES_PASSWORD:-secret}@$${DB_HOST:-localhost}:$${DB_PORT:-5432}/$${POSTGRES_DB:-packs}?sslmode=disable" down 1

migrate-create: ## Create new migration (usage: make migrate-create NAME=create_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=create_users_table"; \
		exit 1; \
	fi
	cd $(BACKEND_DIR) && migrate create -ext sql -dir migrations -seq $(NAME)

test:
	cd $(BACKEND_DIR) && go test -v ./...

build:
	cd $(BACKEND_DIR) && swag init -g cmd/server/main.go -o ./docs
	cd $(BACKEND_DIR) && go build -o app cmd/server/main.go

clean:
	cd $(BACKEND_DIR) && rm -f app && rm -rf docs/ tmp/ build-errors.log 