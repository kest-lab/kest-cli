#!/bin/bash

# Module Validation Script
# Usage: ./validate-module.sh <module_name>

set -e

MODULE_NAME=$1
MODULE_DIR="internal/modules/${MODULE_NAME}"

echo "🔍 Validating module: ${MODULE_NAME}"
echo "======================================="

# Check if module directory exists
if [ ! -d "${MODULE_DIR}" ]; then
    echo "❌ Module directory not found: ${MODULE_DIR}"
    exit 1
fi

echo "✓ Module directory exists"

# Check for 8 required files
REQUIRED_FILES=(
    "model.go"
    "dto.go"
    "repository.go"
    "service.go"
    "handler.go"
    "routes.go"
    "provider.go"
    "service_test.go"
)

MISSING_FILES=()

for file in "${REQUIRED_FILES[@]}"; do
    if [ ! -f "${MODULE_DIR}/${file}" ]; then
        MISSING_FILES+=("${file}")
    fi
done

if [ ${#MISSING_FILES[@]} -eq 0 ]; then
    echo "✓ All 8 required files present"
else
    echo "❌ Missing files:"
    for file in "${MISSING_FILES[@]}"; do
        echo "  - ${file}"
    done
    exit 1
fi

# Check for package declaration
echo ""
echo "Checking package declarations..."
for file in "${MODULE_DIR}"/*.go; do
    if ! grep -q "^package ${MODULE_NAME}$" "${file}"; then
        echo "❌ Incorrect package name in: $(basename ${file})"
        echo "   Expected: package ${MODULE_NAME}"
        exit 1
    fi
done
echo "✓ Package declarations correct"

# Check for provider.go content
echo ""
echo "Checking provider.go..."
if ! grep -q "var ProviderSet = wire.NewSet" "${MODULE_DIR}/provider.go"; then
    echo "❌ provider.go missing ProviderSet declaration"
    exit 1
fi
echo "✓ ProviderSet declared"

# Check for repository interface
echo ""
echo "Checking repository interface..."
if ! grep -q "type Repository interface" "${MODULE_DIR}/repository.go"; then
    echo "❌ Repository interface not found"
    exit 1
fi
echo "✓ Repository interface defined"

# Check for service interface
echo ""
echo "Checking service interface..."
if ! grep -q "type Service interface" "${MODULE_DIR}/service.go"; then
    echo "❌ Service interface not found"
    exit 1
fi
echo "✓ Service interface defined"

# Check for handler struct
echo ""
echo "Checking handler..."
if ! grep -q "type Handler struct" "${MODULE_DIR}/handler.go"; then
    echo "❌ Handler struct not found"
    exit 1
fi
echo "✓ Handler struct defined"

# Check for routes registration
echo ""
echo "Checking routes..."
if ! grep -q "func RegisterRoutes" "${MODULE_DIR}/routes.go"; then
    echo "❌ RegisterRoutes function not found"
    exit 1
fi
echo "✓ RegisterRoutes function defined"

# Try to build the module
echo ""
echo "Building module..."
if ! go build ./internal/modules/${MODULE_NAME}/...; then
    echo "❌ Module build failed"
    exit 1
fi
echo "✓ Module builds successfully"

# Run tests
echo ""
echo "Running tests..."
if ! go test ./internal/modules/${MODULE_NAME}/... -v; then
    echo "⚠️  Some tests failed (review output above)"
else
    echo "✓ All tests passed"
fi

echo ""
echo "======================================="
echo "✅ Module validation complete!"
echo ""
echo "Next steps:"
echo "1. Run 'cd internal/wiring && wire' to generate DI code"
echo "2. Register routes in routes/api.go"
echo "3. Create database migration"
echo "4. Test API endpoints"
