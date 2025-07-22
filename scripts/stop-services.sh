#!/bin/bash

echo "Stopping microservices..."

# Stop services using PID files
for service in user product web; do
    if [ -f "${service}_service.pid" ]; then
        pid=$(cat "${service}_service.pid")
        if kill -0 $pid 2>/dev/null; then
            echo "Stopping $service service (PID: $pid)..."
            kill $pid
        fi
        rm -f "${service}_service.pid"
    fi
done

echo "All services stopped!"