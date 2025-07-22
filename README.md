# Go Microservice Boilerplate

A production-ready microservice architecture built with Go, featuring gRPC services, REST API gateway, MongoDB, Redis, and comprehensive Swagger documentation.

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Gateway       │    │   User Service   │    │ Product Service │
│   (REST API)    │◄──►│     (gRPC)       │    │     (gRPC)      │
│   Port: 8080    │    │   Port: 50051    │    │   Port: 50052   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                        │                        │
         └────────────────────────┼────────────────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │                           │
            ┌───────▼────────┐         ┌────────▼────────┐
            │   MongoDB      │         │     Redis       │
            │   Port: 27017  │         │   Port: 6379    │
            └────────────────┘         └─────────────────┘
```

## Features

### Core Features
- **Microservice Architecture** - Clean separation of concerns
- **gRPC Communication** - High-performance service-to-service communication
- **REST API Gateway** - Public HTTP API with Swagger documentation
- **Database Layer** - MongoDB with connection pooling
- **Caching Layer** - Redis for high-performance caching
- **Configuration Management** - YAML-based configuration with environment overrides

### Documentation & Development
- **Swagger UI** - Interactive API documentation with authentication
- **Structured Logging** - JSON logging with configurable levels
- **Health Checks** - Service health monitoring endpoints
- **Hot Reload** - Development-friendly with automatic restarts

### Security & Production
- **Authentication** - Environment-based Swagger authentication
- **CORS Support** - Configurable cross-origin resource sharing
- **Graceful Shutdown** - Proper service lifecycle management
- **Error Handling** - Consistent error responses across services

## Quick Start

### Prerequisites

- **Go 1.21+**
- **MongoDB 6.0+**
- **Redis 7.0+**
- **Protocol Buffers** (for gRPC)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/go-microservice-boilerplate.git
   cd go-microservice-boilerplate
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Install development tools:**
   ```bash
   # Install Protocol Buffers compiler
   brew install protobuf  # macOS
   # sudo apt-get install protobuf-compiler  # Ubuntu

   # Install Go tools
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

4. **Setup databases:**
   ```bash
   # Using Docker
   docker-compose up -d mongodb redis

   # Or install locally
   # MongoDB: https://docs.mongodb.com/manual/installation/
   # Redis: https://redis.io/download
   ```

5. **Generate code:**
   ```bash
   make proto    # Generate gRPC code
   make swagger  # Generate Swagger documentation
   ```

## Running the Services

### Development Mode (Recommended)

```bash
# Terminal 1: User Service
make run-dev-user

# Terminal 2: Product Service  
make run-dev-product

# Terminal 3: Gateway Service
make run-dev-web
```

### Production Mode

```bash
# Build and run
make build
./bin/main user &
./bin/main product &
./bin/main web
```

### Using Scripts

```bash
# Start all services
./scripts/run-services.sh

# Stop all services
./scripts/stop-services.sh
```

### Using Docker

```bash
# Start with Docker Compose
docker-compose up --build

# Start only databases
docker-compose up -d mongodb redis
```

## API Documentation

### Swagger UI

Access interactive API documentation at:
- **URL:** `http://localhost:8080/swagger/index.html`
- **Authentication:** Set environment variables for security

```bash
# Enable Swagger authentication
export SWAGGER_USERNAME=admin
export SWAGGER_PASSWORD=swagger123

# Or run without authentication
unset SWAGGER_USERNAME SWAGGER_PASSWORD
```

### API Endpoints

#### Health Check
- `GET /api/v1/health` - Service health status

#### Users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user
- `GET /api/v1/users` - List users (with pagination)

#### Products
- `POST /api/v1/products` - Create product
- `GET /api/v1/products/{id}` - Get product by ID
- `PUT /api/v1/products/{id}` - Update product
- `DELETE /api/v1/products/{id}` - Delete product
- `GET /api/v1/products` - List products (with pagination and filtering)

## Configuration

### Environment Variables

Create a `.env` file or set environment variables:

```bash
# Database Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=microservices_db
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# Service Ports
GATEWAY_PORT=8080
USER_SERVICE_PORT=50051
PRODUCT_SERVICE_PORT=50052

# Swagger Authentication
SWAGGER_USERNAME=admin
SWAGGER_PASSWORD=swagger123

# Security
JWT_SECRET=your-super-secure-jwt-secret

# Logging
LOG_LEVEL=info
```

### Configuration File

Modify `configs/config.yaml`:

```yaml
app:
  name: "Go Microservice Boilerplate"
  version: "1.0.0"
  environment: "development"

services:
  gateway:
    port: "8080"
    host: "0.0.0.0"
  user:
    port: "50051"
    host: "0.0.0.0"
  product:
    port: "50052"
    host: "0.0.0.0"

database:
  mongodb:
    uri: "mongodb://localhost:27017"
    database: "microservices_db"
  redis:
    addr: "localhost:6379"
swagger:
  enabled: true
  auth:
    enabled: true
    username: "admin"
    password: "swagger123"
  title: "Go Microservice API"
  version: "1.0.0"
  description: "Microservice API for user and product management"
```

## Development

### Available Make Commands

```bash
make help           # Show all available commands
make build          # Build the application
make proto          # Generate protobuf files
make swagger        # Generate Swagger documentation
make lint           # Run linter
make clean          # Clean build artifacts
make docker-build   # Build Docker image
make docker-up      # Start with Docker Compose
make docker-down    # Stop Docker containers
```

### Project Structure

```
go-microservice-boilerplate/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── database/               # Database connections
│   ├── middleware/             # HTTP middleware
│   ├── proto/                  # Protocol buffer definitions
│   ├── services/
│   │   ├── gateway/            # REST API gateway
│   │   ├── user/               # User microservice
│   │   └── product/            # Product microservice
│   └── utils/                  # Shared utilities
├── configs/                    # Configuration files
├── docs/                       # Generated Swagger docs
├── scripts/                    # Development scripts
├── deployments/                # Deployment configurations
├── docker-compose.yml          # Docker Compose configuration
├── Dockerfile                  # Docker image definition
├── Makefile                    # Build automation
└── README.md                   # This file
```

### Adding a New Service

1. **Create service directory:**
   ```bash
   mkdir -p internal/services/newservice/{handler,model,repository,service}
   ```

2. **Define protobuf:**
   ```bash
   touch internal/proto/newservice/newservice.proto
   ```

3. **Generate code:**
   ```bash
   make proto
   ```

4. **Implement service logic:**
    - Repository layer (database operations)
    - Service layer (business logic)
    - Handler layer (gRPC endpoints)

5. **Update gateway:**
    - Add client in `internal/services/gateway/client/`
    - Add routes in gateway handler



### API Testing

Use the Swagger UI or tools like curl:

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123"}'

# List users
curl "http://localhost:8080/api/v1/users?page=1&limit=10"
```

## Deployment

### Docker Deployment

```bash
# Build and deploy with Docker Compose
docker-compose up --build -d

# Scale services
docker-compose up --scale user-service=3 --scale product-service=2
```

### Production Environment

1. **Set production environment variables:**
   ```bash
   export ENVIRONMENT=production
   export MONGODB_URI=mongodb://prod-mongodb:27017
   export REDIS_ADDR=prod-redis:6379
   export JWT_SECRET=your-production-jwt-secret
   ```

2. **Build for production:**
   ```bash
   make build
   ```

3. **Deploy binaries:**
   ```bash
   # Copy binaries to production server
   scp bin/main user@production-server:/opt/microservices/
   ```

### Kubernetes (Optional)

Example Kubernetes manifests are available in the `deployments/k8s/` directory.

## Troubleshooting

### Common Issues

1. **Port already in use:**
   ```bash
   # Kill process on port 8080
   sudo lsof -t -i tcp:8080 | xargs kill -9
   ```

2. **MongoDB connection failed:**
   ```bash
   # Check MongoDB status
   brew services list | grep mongodb  # macOS
   sudo systemctl status mongod       # Linux
   ```

3. **gRPC generation failed:**
   ```bash
   # Reinstall protobuf tools
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

4. **Swagger not loading:**
   ```bash
   # Regenerate Swagger docs
   make swagger
   
   # Check if docs directory exists
   ls -la docs/
   ```

### Debug Mode

Enable debug logging:

```bash
export LOG_LEVEL=debug
make run-dev-web
```

## Contributing

1. **Fork the repository**
2. **Create a feature branch:**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Commit your changes:**
   ```bash
   git commit -m 'Add some amazing feature'
   ```
4. **Push to the branch:**
   ```bash
   git push origin feature/amazing-feature
   ```
5. **Open a Pull Request**

### Development Guidelines

- Follow Go best practices and idioms
- Add tests for new functionality
- Update documentation for API changes
- Use conventional commit messages
- Ensure all tests pass before submitting PR

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [gRPC-Go](https://github.com/grpc/grpc-go)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [Redis Go Client](https://github.com/go-redis/redis)
- [Swaggo](https://github.com/swaggo/swag)

## Support

If you have any questions or need help, please:

1. Check the [documentation](#-api-documentation)
2. Search through existing [issues](https://github.com/MadhuS-1605/go-microservice-boilerplate/issues)
3. Create a new [issue](https://github.com/MadhuS-1605/go-microservice-boilerplate/issues/new)

---

**Made with ❤️ by [Madhu S Gowda](https://github.com/MadhuS-1605)**