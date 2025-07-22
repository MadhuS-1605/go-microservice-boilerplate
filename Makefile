.PHONY: help build run-web run-user run-product proto docker-up docker-down clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/main cmd/main.go

proto: ## Generate protobuf files
	@echo "Generating protobuf files..."
	protoc --go_out=. --go-grpc_out=. internal/proto/common/common.proto
	protoc --go_out=. --go-grpc_out=. internal/proto/user/user.proto
	protoc --go_out=. --go-grpc_out=. internal/proto/product/product.proto

run-web: build ## Run the web gateway service
	./bin/main web

run-user: build ## Run the user service
	./bin/main user

run-product: build ## Run the product service
	./bin/main product

run-dev-web: ## Run web gateway in development mode
	go run cmd/main.go web

run-dev-user: ## Run user service in development mode
	go run cmd/main.go user

run-dev-product: ## Run product service in development mode
	go run cmd/main.go product

docker-build: ## Build docker image
	docker build -t go-microservices .

docker-up: ## Start all services with docker-compose
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## View docker logs
	docker-compose logs -f

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

deps: ## Download dependencies
	go mod download
	go mod tidy

lint: ## Run linter
	golangci-lint run

setup: deps proto ## Setup the project
	@echo "Setting up the project..."
	@echo "Project setup complete!"

# Development shortcuts
dev: ## Run all services in development mode (requires multiple terminals)
	@echo "Run these commands in separate terminals:"
	@echo "make run-dev-user"
	@echo "make run-dev-product"
	@echo "make run-dev-web"

swagger: ## Generate swagger documentation
	swag init -g cmd/main.go -o docs/

.PHONY: swagger