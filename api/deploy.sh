#!/bin/bash

# Deployment script for ZGO
set -e

echo "🚀 Starting deployment of ZGO..."

# Build the Docker image
echo "📦 Building Docker image..."
docker build -t zgo:latest .

# Stop existing containers
echo "🛑 Stopping existing containers..."
docker-compose -f docker-compose.prod.yml down || true

# Start the production environment
echo "🔄 Starting production environment..."
docker-compose -f docker-compose.prod.yml up -d

# Wait for container to be ready
echo "⏳ Waiting for container to be ready..."
sleep 10

# Health check
echo "🏥 Performing health check..."
# Note: Ensure the health check route exists or update this URL
HEALTH_CHECK=$(curl -s http://localhost:8025/v1/health/status || echo "failed")

if [[ $HEALTH_CHECK == *"ok"* ]]; then
    echo "✅ Deployment successful! Server is running at http://localhost:8025"
    echo "📊 Health status: $HEALTH_CHECK"
    echo ""
    echo "📋 Available endpoints:"
    echo "  - Health: http://localhost:8025/v1/health/status"
    echo "  - Register: POST http://localhost:8025/v1/register"
    echo "  - Login: POST http://localhost:8025/v1/login"
else
    echo "❌ Deployment failed! Health check returned: $HEALTH_CHECK"
    echo "📜 Container logs:"
    docker logs zgo-prod
    exit 1
fi
