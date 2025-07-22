#!/bin/bash

set -e

echo "Starting microservices..."

make swagger

# Function to run service in background
run_service() {
    local service=$1
    local port=$2
    echo "Starting $service service on port $port..."
    go run cmd/main.go $service &
    echo $! > ${service}_service.pid
}

# Start services
run_service "user" "50051"
sleep 2
run_service "product" "50052"
sleep 2
run_service "web" "8080"

echo "All services started!"
echo "Gateway API: http://localhost:8080"
echo "Gateway Swagger: http://localhost:8080/swagger/index.html#"
echo "User gRPC: localhost:50051"
echo "Product gRPC: localhost:50052"

echo "To stop services, run: scripts/stop-services.sh"