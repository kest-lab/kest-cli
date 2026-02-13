#!/bin/bash

# Testing Strategy Validation Script
# Usage: ./validate-tests.sh <path_to_module_or_package>

set -e

TARGET_PATH=$1

if [ -z "$TARGET_PATH" ]; then
    echo "Usage: ./validate-tests.sh <path_to_module_or_package>"
    exit 1
fi

if [ ! -d "$TARGET_PATH" ]; then
    echo "❌ Target path not found: $TARGET_PATH"
    exit 1
fi

echo "🔍 Validating testing standards for '$TARGET_PATH'..."
echo "=============================================="

# Track validation status
ERRORS=0
WARNINGS=0

# 1. Check for test files
TEST_FILES=$(find "$TARGET_PATH" -name "*_test.go")
if [ -z "$TEST_FILES" ]; then
    echo "⚠️  No test files found in $TARGET_PATH"
    WARNINGS=$((WARNINGS + 1))
else
    echo "✅ Found $(echo "$TEST_FILES" | wc -l) test files"
fi

# 2. Check for Table-Driven tests
TABLE_TESTS=$(grep -r "struct {" "$TARGET_PATH" | grep "_test.go" || true)
if [ -z "$TABLE_TESTS" ] && [ -n "$TEST_FILES" ]; then
    echo "⚠️  No table-driven tests detected (look for tests using slice of structs)"
    WARNINGS=$((WARNINGS + 1))
else
    echo "✅ Table-driven tests detected"
fi

# 3. Check for testify/assert usage
if grep -rq "github.com/stretchr/testify/assert" "$TARGET_PATH"; then
    echo "✅ Using testify/assert for clean assertions"
else
    echo "⚠️  Not using testify/assert (recommended for consistency)"
    WARNINGS=$((WARNINGS + 1))
fi

# 4. Check for mocking in service tests
SERVICE_TESTS=$(find "$TARGET_PATH" -name "service_test.go")
if [ -n "$SERVICE_TESTS" ]; then
    if grep -rq "github.com/stretchr/testify/mock" "$TARGET_PATH"; then
        echo "✅ Mocking detected in service tests"
    else
        echo "⚠️  No mocking detected in service tests. Ensure dependencies are mocked."
        WARNINGS=$((WARNINGS + 1))
    fi
fi

# 5. Check for AssertExpectations
MOCK_VERIFICATION=$(grep -r "AssertExpectations" "$TARGET_PATH" | grep "_test.go" || true)
if [ -z "$MOCK_VERIFICATION" ] && grep -rq "github.com/stretchr/testify/mock" "$TARGET_PATH"; then
    echo "❌ Mock expectations are not being verified (missing AssertExpectations)"
    ERRORS=$((ERRORS + 1))
else
    echo "✅ Mock verification detected"
fi

# 6. Check for large test files
LARGE_TEST_FILES=$(find "$TARGET_PATH" -name "*_test.go" -exec wc -l {} + | awk '$1 > 800 {print $2}')
if [ -n "$LARGE_TEST_FILES" ]; then
    echo "⚠️  Large test files detected (>800 lines):"
    echo "$LARGE_TEST_FILES"
    WARNINGS=$((WARNINGS + 1))
fi

echo ""
echo "=============================================="
echo "📊 Test Validation Summary"
echo "=============================================="

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo "✅ Testing standards checks passed!"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo "⚠️  Found $WARNINGS warning(s). Please review the suggestions."
    exit 0
else
    echo "❌ Found $ERRORS error(s) and $WARNINGS warning(s). Standards MUST be met."
    exit 1
fi
