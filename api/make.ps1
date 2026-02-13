# ZGO Build Script for Windows PowerShell
# Usage: .\make.ps1 [command]

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

$ErrorActionPreference = "Stop"

# Configuration
$BINARY_NAME = "zgo.exe"
$SERVER_NAME = "server.exe"

# Get version information
try {
    $VERSION = (git describe --tags --always --dirty 2>$null) | Out-String
    $VERSION = $VERSION.Trim()
    if ([string]::IsNullOrEmpty($VERSION)) { $VERSION = "dev" }
} catch {
    $VERSION = "dev"
}

try {
    $GIT_COMMIT = (git rev-parse --short HEAD 2>$null) | Out-String
    $GIT_COMMIT = $GIT_COMMIT.Trim()
    if ([string]::IsNullOrEmpty($GIT_COMMIT)) { $GIT_COMMIT = "unknown" }
} catch {
    $GIT_COMMIT = "unknown"
}

$BUILD_TIME = Get-Date -Format "yyyy-MM-dd_HH:mm:ss"
$LDFLAGS = "-ldflags `"-s -w -X main.Version=$VERSION -X main.GitCommit=$GIT_COMMIT -X main.BuildTime=$BUILD_TIME`""

# Helper function to check if a command exists
function Test-CommandExists {
    param($Command)
    $null = Get-Command $Command -ErrorAction SilentlyContinue
    return $?
}

# Build CLI tool
function Build {
    Write-Host "Building $BINARY_NAME..." -ForegroundColor Green
    Wire
    Invoke-Expression "go build $LDFLAGS -o $BINARY_NAME cmd/zgo/main.go"
}

# Build server
function Build-Server {
    Write-Host "Building $SERVER_NAME..." -ForegroundColor Green
    Invoke-Expression "go build $LDFLAGS -o $SERVER_NAME cmd/server/main.go"
}

# Build all
function Build-All {
    Build
    Build-Server
}

# Install CLI to GOPATH/bin
function Install {
    Write-Host "Installing $BINARY_NAME..." -ForegroundColor Green
    Build
    Invoke-Expression "go install $LDFLAGS ./cmd/zgo"
}

# Clean build artifacts
function Clean {
    Write-Host "Cleaning..." -ForegroundColor Green
    go clean
    if (Test-Path $BINARY_NAME) { Remove-Item $BINARY_NAME }
    if (Test-Path $SERVER_NAME) { Remove-Item $SERVER_NAME }
    if (Test-Path "coverage.txt") { Remove-Item "coverage.txt" }
    if (Test-Path "coverage.html") { Remove-Item "coverage.html" }
}

# Run tests
function Test {
    Write-Host "Running tests..." -ForegroundColor Green
    go test ./... -v
}

# Run tests with race detection
function Test-Race {
    Write-Host "Running tests with race detection..." -ForegroundColor Green
    go test ./... -race -v
}

# Run tests with coverage
function Cover {
    Write-Host "Running tests with coverage..." -ForegroundColor Green
    go test ./... -coverprofile=coverage.txt -covermode=atomic
    go tool cover -html=coverage.txt -o coverage.html
    Write-Host "Coverage report: coverage.html" -ForegroundColor Cyan
}

# Run linter
function Lint {
    Write-Host "Running linter..." -ForegroundColor Green
    if (-not (Test-CommandExists "golangci-lint")) {
        Write-Host "Installing golangci-lint..." -ForegroundColor Yellow
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    }
    golangci-lint run ./...
}

# Fix lint issues
function Lint-Fix {
    Write-Host "Fixing lint issues..." -ForegroundColor Green
    golangci-lint run --fix ./...
}

# Generate code
function Generate {
    Write-Host "Generating code..." -ForegroundColor Green
    go generate ./...
}

# Run Wire
function Wire {
    Write-Host "Running Wire..." -ForegroundColor Green
    if (-not (Test-CommandExists "wire")) {
        Write-Host "Installing wire..." -ForegroundColor Yellow
        go install github.com/google/wire/cmd/wire@latest
    }
    Push-Location internal/wiring
    wire
    Pop-Location
}

# Generate API documentation
function Docs {
    Write-Host "Generating API documentation..." -ForegroundColor Green
    if (-not (Test-CommandExists "swag")) {
        Write-Host "Installing swag..." -ForegroundColor Yellow
        go install github.com/swaggo/swag/cmd/swag@latest
    }
    swag init -g cmd/server/main.go -o docs/swagger
}

# Generate mocks
function Mock {
    Write-Host "Generating mocks..." -ForegroundColor Green
    if (-not (Test-CommandExists "mockgen")) {
        Write-Host "Installing mockgen..." -ForegroundColor Yellow
        go install go.uber.org/mock/mockgen@latest
    }
    go generate ./...
}

# Run development server
function Dev {
    Write-Host "Starting development server..." -ForegroundColor Green
    go run cmd/server/main.go
}

# Run server
function Server {
    Write-Host "Starting server..." -ForegroundColor Green
    go run cmd/server/main.go
}

# Run with Air (hot reload)
function Air {
    Write-Host "Starting with Air (hot reload)..." -ForegroundColor Green
    if (-not (Test-CommandExists "air")) {
        Write-Host "Installing air..." -ForegroundColor Yellow
        go install github.com/air-verse/air@latest
    }
    air
}

# Tidy dependencies
function Tidy {
    Write-Host "Tidying dependencies..." -ForegroundColor Green
    go mod tidy
}

# Update dependencies
function Update {
    Write-Host "Updating dependencies..." -ForegroundColor Green
    go get -u ./...
    go mod tidy
}

# Check for vulnerabilities
function Vuln {
    Write-Host "Checking for vulnerabilities..." -ForegroundColor Green
    if (-not (Test-CommandExists "govulncheck")) {
        Write-Host "Installing govulncheck..." -ForegroundColor Yellow
        go install golang.org/x/vuln/cmd/govulncheck@latest
    }
    govulncheck ./...
}

# Setup development environment
function Setup {
    Write-Host "Setting up development environment..." -ForegroundColor Green
    go mod download
    Write-Host "✅ Dependencies downloaded" -ForegroundColor Green
    Write-Host "Run '.\make.ps1 install' to install zgo CLI globally" -ForegroundColor Cyan
}

# Show help
function Help {
    Write-Host @"
Available commands:
  .\make.ps1 build        - Build the CLI tool
  .\make.ps1 build-server - Build the server
  .\make.ps1 build-all    - Build all binaries
  .\make.ps1 install      - Install CLI to GOPATH/bin
  .\make.ps1 test         - Run tests
  .\make.ps1 test-race    - Run tests with race detection
  .\make.ps1 cover        - Run tests with coverage report
  .\make.ps1 lint         - Run golangci-lint
  .\make.ps1 lint-fix     - Fix lint issues automatically
  .\make.ps1 generate     - Run go generate
  .\make.ps1 wire         - Run Wire DI generator
  .\make.ps1 docs         - Generate Swagger documentation
  .\make.ps1 mock         - Generate mocks
  .\make.ps1 dev          - Run development server
  .\make.ps1 air          - Run with hot reload (Air)
  .\make.ps1 server       - Run server
  .\make.ps1 tidy         - Tidy dependencies
  .\make.ps1 update       - Update dependencies
  .\make.ps1 vuln         - Check for vulnerabilities
  .\make.ps1 setup        - Setup development environment
  .\make.ps1 clean        - Clean build artifacts
  .\make.ps1 help         - Show this help
"@ -ForegroundColor Cyan
}

# Route to appropriate function based on command
switch ($Command.ToLower()) {
    "build"        { Build }
    "build-server" { Build-Server }
    "build-all"    { Build-All }
    "install"      { Install }
    "clean"        { Clean }
    "test"         { Test }
    "test-race"    { Test-Race }
    "cover"        { Cover }
    "lint"         { Lint }
    "lint-fix"     { Lint-Fix }
    "generate"     { Generate }
    "wire"         { Wire }
    "docs"         { Docs }
    "mock"         { Mock }
    "dev"          { Dev }
    "server"       { Server }
    "air"          { Air }
    "tidy"         { Tidy }
    "update"       { Update }
    "vuln"         { Vuln }
    "setup"        { Setup }
    "help"         { Help }
    default        { Help }
}
