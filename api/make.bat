@echo off
REM ZGO Build Script for Windows
REM Usage: make.bat [command]

setlocal

REM Get version information
for /f "tokens=*" %%i in ('git describe --tags --always --dirty 2^>nul') do set VERSION=%%i
if "%VERSION%"=="" set VERSION=dev

for /f "tokens=*" %%i in ('git rev-parse --short HEAD 2^>nul') do set GIT_COMMIT=%%i
if "%GIT_COMMIT%"=="" set GIT_COMMIT=unknown

for /f "tokens=*" %%i in ('powershell -Command "Get-Date -Format 'yyyy-MM-dd_HH:mm:ss'"') do set BUILD_TIME=%%i

set BINARY_NAME=zgo.exe
set SERVER_NAME=server.exe
set LDFLAGS=-ldflags "-s -w -X main.Version=%VERSION% -X main.GitCommit=%GIT_COMMIT% -X main.BuildTime=%BUILD_TIME%"

REM Default command
if "%1"=="" goto help

REM Route to specific command
if /i "%1"=="build" goto build
if /i "%1"=="build-server" goto build-server
if /i "%1"=="build-all" goto build-all
if /i "%1"=="install" goto install
if /i "%1"=="clean" goto clean
if /i "%1"=="test" goto test
if /i "%1"=="test-race" goto test-race
if /i "%1"=="cover" goto cover
if /i "%1"=="lint" goto lint
if /i "%1"=="lint-fix" goto lint-fix
if /i "%1"=="generate" goto generate
if /i "%1"=="wire" goto wire
if /i "%1"=="docs" goto docs
if /i "%1"=="mock" goto mock
if /i "%1"=="dev" goto dev
if /i "%1"=="server" goto server
if /i "%1"=="tidy" goto tidy
if /i "%1"=="setup" goto setup
if /i "%1"=="help" goto help
goto help

:build
echo Building %BINARY_NAME%...
call :wire
go build %LDFLAGS% -o %BINARY_NAME% cmd/zgo/main.go
goto end

:build-server
echo Building %SERVER_NAME%...
go build %LDFLAGS% -o %SERVER_NAME% cmd/server/main.go
goto end

:build-all
call :build
call :build-server
goto end

:install
echo Installing %BINARY_NAME%...
call :build
go install %LDFLAGS% ./cmd/zgo
goto end

:clean
echo Cleaning...
go clean
if exist %BINARY_NAME% del %BINARY_NAME%
if exist %SERVER_NAME% del %SERVER_NAME%
if exist coverage.txt del coverage.txt
if exist coverage.html del coverage.html
goto end

:test
echo Running tests...
go test ./... -v
goto end

:test-race
echo Running tests with race detection...
go test ./... -race -v
goto end

:cover
echo Running tests with coverage...
go test ./... -coverprofile=coverage.txt -covermode=atomic
go tool cover -html=coverage.txt -o coverage.html
echo Coverage report: coverage.html
goto end

:lint
echo Running linter...
where golangci-lint >nul 2>nul
if %errorlevel% neq 0 (
    echo Installing golangci-lint...
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
)
golangci-lint run ./...
goto end

:lint-fix
echo Fixing lint issues...
golangci-lint run --fix ./...
goto end

:generate
echo Generating code...
go generate ./...
goto end

:wire
echo Running Wire...
where wire >nul 2>nul
if %errorlevel% neq 0 (
    echo Installing wire...
    go install github.com/google/wire/cmd/wire@latest
)
cd internal\wiring
wire
cd ..\..
goto :eof

:docs
echo Generating API documentation...
where swag >nul 2>nul
if %errorlevel% neq 0 (
    echo Installing swag...
    go install github.com/swaggo/swag/cmd/swag@latest
)
swag init -g cmd/server/main.go -o docs/swagger
goto end

:mock
echo Generating mocks...
where mockgen >nul 2>nul
if %errorlevel% neq 0 (
    echo Installing mockgen...
    go install go.uber.org/mock/mockgen@latest
)
go generate ./...
goto end

:dev
echo Starting development server...
go run cmd/server/main.go
goto end

:server
echo Starting server...
go run cmd/server/main.go
goto end

:tidy
echo Tidying dependencies...
go mod tidy
goto end

:setup
echo Setting up development environment...
go mod download
echo Dependencies downloaded
echo Run 'make.bat install' to install zgo CLI globally
goto end

:help
echo Available commands:
echo   make.bat build        - Build the CLI tool
echo   make.bat build-server - Build the server
echo   make.bat build-all    - Build all binaries
echo   make.bat install      - Install CLI to GOPATH/bin
echo   make.bat test         - Run tests
echo   make.bat test-race    - Run tests with race detection
echo   make.bat cover        - Run tests with coverage report
echo   make.bat lint         - Run golangci-lint
echo   make.bat lint-fix     - Fix lint issues automatically
echo   make.bat generate     - Run go generate
echo   make.bat wire         - Run Wire DI generator
echo   make.bat docs         - Generate Swagger documentation
echo   make.bat mock         - Generate mocks
echo   make.bat dev          - Run development server
echo   make.bat server       - Run server
echo   make.bat tidy         - Tidy dependencies
echo   make.bat setup        - Setup development environment
echo   make.bat clean        - Clean build artifacts
echo   make.bat help         - Show this help
goto end

:end
endlocal
