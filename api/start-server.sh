#!/bin/bash
export PORT=8025
export GIN_MODE=debug
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=kest
export DB_USERNAME=kest_user
export DB_PASSWORD=kest_password_123
export JWT_SECRET=kest-jwt-secret-key-for-development-only-change-in-production
export JWT_EXPIRATION=24h
export ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
export DB_AUTO_MIGRATE=false

./test-server
