#!/bin/bash

# ===========================================
# Kest Platform Build Script
# ===========================================
# Builds Vite frontend and embeds into Go binary

set -e  # Exit on error

echo "🏗️  Building Kest Platform..."
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# ===========================================
# 1. Build Frontend (Vite)
# ===========================================
echo -e "${BLUE}📦 Building frontend (Vite + React)...${NC}"
cd "$PROJECT_ROOT/web"

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "  📥 Installing dependencies..."
    npm install
fi

# Build Vite
echo "  ⚙️  Compiling Vite..."
npm run build

# Check if build succeeded
if [ ! -d "dist" ]; then
    echo -e "${YELLOW}⚠️  Error: dist/ directory not found${NC}"
    exit 1
fi

echo -e "${GREEN}  ✅ Frontend built successfully${NC}"
echo ""

# ===========================================
# 2. Build Backend (Go + embedded frontend)
# ===========================================
echo -e "${BLUE}🔨 Building backend (Go + embedded frontend)...${NC}"
cd "$PROJECT_ROOT/api"

# Download Go dependencies
echo "  📥 Downloading Go modules..."
go mod download

# Build binary
echo "  ⚙️  Compiling Go binary..."
cd "$PROJECT_ROOT"
CGO_ENABLED=0 go build -ldflags="-s -w" -o kest-server cmd/main.go

# Check if build succeeded
if [ ! -f "$PROJECT_ROOT/kest-server" ]; then
    echo -e "${YELLOW}⚠️  Error: Binary not created${NC}"
    exit 1
fi

echo -e "${GREEN}  ✅ Backend built successfully${NC}"
echo ""

# ===========================================
# Summary
# ===========================================
echo -e "${GREEN}🎉 Build completed!${NC}"
echo ""
echo "📦 Binary: $PROJECT_ROOT/kest-server"
echo "📏 Size: $(du -h "$PROJECT_ROOT/kest-server" | cut -f1)"
echo ""
echo "🚀 To run:"
echo "   cd $PROJECT_ROOT"
echo "   ./kest-server"
echo ""
echo "   Server will start on http://localhost:8080"
