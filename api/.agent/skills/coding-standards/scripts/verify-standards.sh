#!/bin/bash

# Coding Standards Verification Script
# Usage: ./verify-standards.sh <module_name>

set -e

MODULE=$1
MODULE_DIR="internal/modules/${MODULE}"

if [ -z "$MODULE" ]; then
    echo "Usage: ./verify-standards.sh <module_name>"
    exit 1
fi

if [ ! -d "${MODULE_DIR}" ]; then
    echo "❌ Module directory not found: ${MODULE_DIR}"
    exit 1
fi

echo "🔍 Verifying coding standards for '${MODULE}' module..."
echo "================================================"

# 1. File structure
echo ""
echo "📁 Level 1: Checking file structure..."
if [ -f ".agent/skills/module-creation/scripts/validate-module.sh" ]; then
    .agent/skills/module-creation/scripts/validate-module.sh ${MODULE}
else
    echo "⚠️  Module validation script not found, skipping..."
fi

# 2. Naming conventions
echo ""
echo "📝 Level 2: Checking naming conventions..."

# Check for camelCase JSON tags
CAMEL_CASE=$(grep -rn 'json:"[a-z][a-zA-Z]*[A-Z]' ${MODULE_DIR}/ 2>/dev/null || true)
if [ -n "$CAMEL_CASE" ]; then
    echo "❌ Found camelCase JSON tags (should be snake_case):"
    echo "$CAMEL_CASE"
else
    echo "✅ JSON tags use snake_case"
fi

# Check PO suffix on models
if [ -f "${MODULE_DIR}/model.go" ]; then
    MISSING_PO=$(grep 'type.*struct' ${MODULE_DIR}/model.go | grep -v 'PO\s' || true)
    if [ -n "$MISSING_PO" ]; then
        echo "⚠️  Model types should have PO suffix:"
        echo "$MISSING_PO"
    else
        echo "✅ Model types have PO suffix"
    fi
fi

# 3. Architecture compliance
echo ""
echo "🏗️  Level 3: Checking architecture compliance..."

# Check for handler accessing repo directly
if [ -f "${MODULE_DIR}/handler.go" ]; then
    DIRECT_REPO=$(grep -n 'h\.repo\.' ${MODULE_DIR}/handler.go 2>/dev/null || true)
    if [ -n "$DIRECT_REPO" ]; then
        echo "❌ Handler directly accessing repository (should use service only):"
        echo "$DIRECT_REPO"
    else
        echo "✅ Handler properly uses service layer"
    fi
fi

# Check PO usage in service
if [ -f "${MODULE_DIR}/service.go" ]; then
    PO_IN_SERVICE=$(grep -n 'PO{' ${MODULE_DIR}/service.go 2>/dev/null || true)
    if [ -n "$PO_IN_SERVICE" ]; then
        echo "❌ Service using PO (should use domain entities only):"
        echo "$PO_IN_SERVICE"
    else
        echo "✅ Service uses domain entities"
    fi
fi

# 4. Required content checks
echo ""
echo "📋 Level 4: Checking required content..."

# Check for custom errors in service
if [ -f "${MODULE_DIR}/service.go" ]; then
    CUSTOM_ERRORS=$(grep -n '^var Err' ${MODULE_DIR}/service.go | wc -l)
    if [ "$CUSTOM_ERRORS" -gt 0 ]; then
        echo "✅ Custom errors defined ($CUSTOM_ERRORS found)"
    else
        echo "⚠️  No custom errors defined in service"
    fi
fi

# Check for TableName method
if [ -f "${MODULE_DIR}/model.go" ]; then
    if grep -q 'func.*TableName' ${MODULE_DIR}/model.go; then
        echo "✅ TableName() method exists"
    else
        echo "❌ TableName() method missing in model"
    fi
fi

# Check for soft delete
if [ -f "${MODULE_DIR}/model.go" ]; then
    if grep -q 'DeletedAt.*gorm.DeletedAt' ${MODULE_DIR}/model.go; then
        echo "✅ Soft delete (DeletedAt) enabled"
    else
        echo "❌ Soft delete (DeletedAt) missing"
    fi
fi

# 5. Security checks
echo ""
echo "🔒 Level 5: Checking security..."

# Check password protection
PASSWORD_PROTECTED=$(grep -rn 'Password.*json:"-"' internal/domain/*.go ${MODULE_DIR}/*.go 2>/dev/null || true)
if [ -n "$PASSWORD_PROTECTED" ]; then
    echo "✅ Password fields protected with json:\"-\""
else
    echo "⚠️  Check password field protection (if applicable)"
fi

# Check for binding tags
if [ -f "${MODULE_DIR}/dto.go" ]; then
    BINDING_COUNT=$(grep -o 'binding:"[^"]*"' ${MODULE_DIR}/dto.go | wc -l)
    if [ "$BINDING_COUNT" -gt 0 ]; then
        echo "✅ Input validation with binding tags ($BINDING_COUNT rules)"
    else
        echo "⚠️  No binding validation tags found"
    fi
fi

# 6. Wire DI
echo ""
echo "⚙️  Level 6: Checking Wire DI..."

if [ -f "${MODULE_DIR}/provider.go" ]; then
    if grep -q 'var ProviderSet' ${MODULE_DIR}/provider.go; then
        echo "✅ ProviderSet exported"
    else
        echo "❌ ProviderSet not found"
    fi
    
    WIRE_BINDS=$(grep -c 'wire.Bind' ${MODULE_DIR}/provider.go || true)
    if [ "$WIRE_BINDS" -ge 2 ]; then
        echo "✅ Wire bindings defined ($WIRE_BINDS bindings)"
    else
        echo "⚠️  Check Wire bindings (expected at least 2)"
    fi
fi

# 7. Testing
echo ""
echo "🧪 Level 7: Running tests..."

if [ -f "${MODULE_DIR}/service_test.go" ]; then
    echo "✅ Test file exists"
    
    # Run tests
    if go test ./${MODULE_DIR}/... -v -cover 2>&1; then
        echo "✅ Tests passed"
    else
        echo "❌ Tests failed"
    fi
else
    echo "❌ Test file missing: service_test.go"
fi

# 8. Code quality
echo ""
echo "🔍 Level 8: Checking code quality..."

# Check for TODO comments
TODO_COUNT=$(grep -rn 'TODO\|FIXME\|HACK' ${MODULE_DIR}/ 2>/dev/null | wc -l || true)
if [ "$TODO_COUNT" -eq 0 ]; then
    echo "✅ No TODO/FIXME comments"
else
    echo "⚠️  Found $TODO_COUNT TODO/FIXME comments (create issues instead)"
fi

# Run linter (if available)
if command -v golangci-lint &> /dev/null; then
    echo "Running golangci-lint..."
    if golangci-lint run ./${MODULE_DIR}/... --timeout=2m 2>&1; then
        echo "✅ Linter passed"
    else
        echo "⚠️  Linter found issues (see above)"
    fi
else
    echo "⚠️  golangci-lint not installed, skipping..."
fi

echo ""
echo "================================================"
echo "✅ Standards verification complete!"
echo ""
echo "Summary:"
echo "  - File structure checked"
echo "  - Naming conventions verified"
echo "  - Architecture compliance checked"
echo "  - Security rules verified"
echo "  - Tests executed"
echo "  - Code quality assessed"
