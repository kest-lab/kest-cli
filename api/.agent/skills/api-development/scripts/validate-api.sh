#!/bin/bash

# API Standards Validation Script
# Usage: ./validate-api.sh <module_name>

set -e

MODULE=$1
MODULE_DIR="internal/modules/${MODULE}"

if [ -z "$MODULE" ]; then
    echo "Usage: ./validate-api.sh <module_name>"
    exit 1
fi

if [ ! -d "${MODULE_DIR}" ]; then
    echo "❌ Module directory not found: ${MODULE_DIR}"
    exit 1
fi

echo "🔍 Validating API standards for '${MODULE}' module..."
echo "=============================================="

HANDLER_FILE="${MODULE_DIR}/handler.go"

if [ ! -f "${HANDLER_FILE}" ]; then
    echo "❌ Handler file not found: ${HANDLER_FILE}"
    exit 1
fi

# Track validation status
ERRORS=0

# =============================================================================
# 1. Check for Pagination in List Methods
# =============================================================================

echo ""
echo "📋 Level 1: Checking pagination..."

# Check if there's a List function
if grep -q "func.*List.*gin.Context" "${HANDLER_FILE}"; then
    # Check if pagination is used
    if grep -q "pagination.PaginateFromContext\|pagination.Paginate" "${HANDLER_FILE}"; then
        echo "✅ Pagination detected in list endpoints"
    else
        echo "❌ List endpoint found but no pagination detected!"
        echo "   Required: pagination.PaginateFromContext[T](c, db)"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "⚠️  No List function found (may not be needed)"
fi

# =============================================================================
# 2. Check for Unified Error Responses
# =============================================================================

echo ""
echo "🚫 Level 2: Checking error handling..."

# Check for manual JSON responses with status codes (forbidden)
MANUAL_ERRORS=$(grep -n 'c\.JSON([0-9][0-9][0-9]' "${HANDLER_FILE}" 2>/dev/null | grep -v "//\|response.Success\|response.Created" || true)

if [ -n "$MANUAL_ERRORS" ]; then
    echo "❌ Found manual status code usage (should use response.* functions):"
    echo "$MANUAL_ERRORS"
    ERRORS=$((ERRORS + 1))
else
    echo "✅ No manual status codes found"
fi

# Check for response package usage
if grep -q "response\." "${HANDLER_FILE}"; then
    echo "✅ Using response package for responses"
else
    echo "⚠️  response package not detected"
fi

# =============================================================================
# 3. Check for Proper HTTP Methods
# =============================================================================

echo ""
echo "📡 Level 3: Checking HTTP method usage..."

# Check DELETE returns NoContent
if grep -q "func.*Delete.*gin.Context" "${HANDLER_FILE}"; then
    if grep -q "response.NoContent" "${HANDLER_FILE}"; then
        echo "✅ DELETE returns 204 No Content"
    else
        echo "❌ DELETE should return response.NoContent() (204)"
        ERRORS=$((ERRORS + 1))
    fi
fi

# Check POST/Create returns Created
if grep -q "func.*Create.*gin.Context" "${HANDLER_FILE}"; then
    if grep -q "response.Created" "${HANDLER_FILE}"; then
        echo "✅ CREATE returns 201 Created"
    else
        echo "⚠️  CREATE should return response.Created() (201)"
    fi
fi

# =============================================================================
# 4. Check for Request Validation
# =============================================================================

echo ""
echo "✅ Level 4: Checking request validation..."

# Check for BindJSON usage
if grep -q "handler.BindJSON" "${HANDLER_FILE}"; then
    echo "✅ Using handler.BindJSON() for validation"
else
    echo "⚠️  handler.BindJSON() not detected (may not be needed)"
fi

# Check DTOs have validation tags
DTO_FILE="${MODULE_DIR}/dto.go"
if [ -f "${DTO_FILE}" ]; then
    if grep -q 'binding:' "${DTO_FILE}"; then
        echo "✅ Validation tags found in DTOs"
    else
        echo "⚠️  No binding validation tags found in DTOs"
    fi
fi

# =============================================================================
# 5. Check for RESTful Patterns
# =============================================================================

echo ""
echo "🌐 Level 5: Checking RESTful patterns..."

ROUTES_FILE="${MODULE_DIR}/routes.go"
if [ -f "${ROUTES_FILE}" ]; then
    # Check for RESTful route registration
    if grep -q "GET\|POST\|PATCH\|DELETE" "${ROUTES_FILE}"; then
        echo "✅ RESTful HTTP methods detected"
    fi
    
    # Look for verb-based routes (anti-pattern)
    VERB_ROUTES=$(grep -i 'create\|update\|delete\|get' "${ROUTES_FILE}" | grep -v 'func\|//' | grep '\"/' || true)
    if [ -n "$VERB_ROUTES" ]; then
        echo "⚠️  Possible verb-based routes detected (should use HTTP methods):"
        echo "$VERB_ROUTES"
    fi
fi

# =============================================================================
# 6. Check for Common Anti-Patterns
# =============================================================================

echo ""
echo "🚨 Level 6: Checking for anti-patterns..."

# Check for c.AbortWithStatusJSON (should use response package)
if grep -q "c.AbortWithStatusJSON" "${HANDLER_FILE}"; then
    echo "❌ Found c.AbortWithStatusJSON() - use response.* functions instead"
    ERRORS=$((ERRORS + 1))
fi

# Check for gin.H usage (prefer struct or response package)
GIN_H_COUNT=$(grep -c 'gin.H{' "${HANDLER_FILE}" 2>/dev/null || echo "0")
if [ "$GIN_H_COUNT" -gt 0 ]; then
    echo "⚠️  Found $GIN_H_COUNT gin.H{} usage - prefer typed structures"
fi

# =============================================================================
# 7. Generate Report
# =============================================================================

echo ""
echo "=============================================="
echo "📊 Validation Summary"
echo "=============================================="

if [ $ERRORS -eq 0 ]; then
    echo "✅ All API standards checks passed!"
    echo ""
    echo "Module '${MODULE}' follows ZGO API development standards."
    exit 0
else
    echo "❌ Found $ERRORS error(s)"
    echo ""
    echo "Please fix the issues above before submitting."
    echo "Refer to: .agent/skills/api-development/SKILL.md"
    exit 1
fi
