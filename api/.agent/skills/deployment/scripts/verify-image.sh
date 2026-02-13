#!/bin/bash

# Deployment Standards Validation Script
# Usage: ./verify-image.sh <path_to_context_dir>

set -e

CONTEXT_PATH=$1

if [ -z "$CONTEXT_PATH" ]; then
    echo "Usage: ./verify-image.sh <path_to_context_dir>"
    exit 1
fi

DOCKERFILE="${CONTEXT_PATH}/Dockerfile"

if [ ! -f "$DOCKERFILE" ]; then
    echo "❌ Dockerfile not found at: $DOCKERFILE"
    exit 1
fi

echo "🔍 Validating Deployment standards for '$CONTEXT_PATH'..."
echo "=============================================="

ERRORS=0
WARNINGS=0

# 1. Check for Multi-Stage Build
if ! grep -q "FROM .* AS builder" "$DOCKERFILE"; then
    echo "❌ Multi-stage build not detected. Production images MUST use multi-stage builds."
    ERRORS=$((ERRORS + 1))
else
    echo "✅ Multi-stage build markers found."
fi

# 2. Check for thin production base image
# Looking for alpine or distroless
if grep "FROM " "$DOCKERFILE" | tail -n 1 | grep -qE "golang|ubuntu|debian|fedora"; then
    echo "⚠️  The final production image seems to be a thick base (golang/ubuntu/etc). Recommended: alpine or distroless."
    WARNINGS=$((WARNINGS + 1))
else
    echo "✅ Lightweight production base image detected."
fi

# 3. Check for .dockerignore
if [ ! -f "${CONTEXT_PATH}/.dockerignore" ]; then
    echo "⚠️  .dockerignore not found. This may result in oversized or insecure images."
    WARNINGS=$((WARNINGS + 1))
else
    echo "✅ .dockerignore found."
    # Check if common sensitive files are ignored
    if ! grep -q ".env" "${CONTEXT_PATH}/.dockerignore"; then
        echo "⚠️  .env not found in .dockerignore. Ensure secrets are not baked into the image."
        WARNINGS=$((WARNINGS + 1))
    fi
fi

# 4. Check for non-root user
if ! grep -q "USER " "$DOCKERFILE"; then
    echo "⚠️  No USER instruction found. Running as root in production is not recommended."
    WARNINGS=$((WARNINGS + 1))
else
    echo "✅ USER instruction detected."
fi

echo "=============================================="
if [ $ERRORS -eq 0 ]; then
    echo "SUCCESS: Deployment standards mostly met ($WARNINGS warnings)."
    exit 0
else
    echo "FAILURE: Found $ERRORS errors. Please fix before proceeding."
    exit 1
fi
