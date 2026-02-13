#!/bin/bash

# ===========================================
# Kest Security Scan Script
# ===========================================
# This script scans the repository for sensitive information
# before committing to ensure no secrets are exposed.

echo "🔍 Scanning for sensitive information..."
echo "=========================================="

# Colors for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Track if any issues were found
ISSUES_FOUND=0

# ===========================================
# 1. Check for .env files (except .env.example)
# ===========================================
echo ""
echo "📋 Checking for .env files..."
ENV_FILES=$(find . -name ".env*" -type f ! -name ".env.example" ! -name ".env.*.example" 2>/dev/null)

if [ -n "$ENV_FILES" ]; then
    echo -e "${RED}⚠️  Found potentially sensitive .env files:${NC}"
    echo "$ENV_FILES"
    ISSUES_FOUND=1
else
    echo -e "${GREEN}✅ No sensitive .env files found${NC}"
fi

# ===========================================
# 2. Check for common secret patterns
# ===========================================
echo ""
echo "🔐 Scanning for hardcoded secrets..."

# Patterns to search for
PATTERNS=(
    "password.*=.*['\"][^'\"]{8,}"
    "secret.*=.*['\"][^'\"]{16,}"
    "api[_-]?key.*=.*['\"][^'\"]{16,}"
    "private[_-]?key.*=.*['\"][^'\"]{16,}"
    "token.*=.*['\"][^'\"]{16,}"
    "jwt.*=.*['\"][^'\"]{16,}"
)

for pattern in "${PATTERNS[@]}"; do
    # Search in common file types, excluding node_modules and vendor
    MATCHES=$(grep -r -i -E "$pattern" \
        --include="*.go" \
        --include="*.ts" \
        --include="*.tsx" \
        --include="*.js" \
        --include="*.jsx" \
        --include="*.json" \
        --include="*.yml" \
        --include="*.yaml" \
        --exclude-dir=node_modules \
        --exclude-dir=vendor \
        --exclude-dir=.git \
        --exclude-dir=.next \
        --exclude-dir=dist \
        --exclude-dir=build \
        . 2>/dev/null || true)
    
    if [ -n "$MATCHES" ]; then
        echo -e "${YELLOW}⚠️  Found potential secret pattern: $pattern${NC}"
        echo "$MATCHES" | head -n 5
        echo "..."
        ISSUES_FOUND=1
    fi
done

if [ $ISSUES_FOUND -eq 0 ]; then
    echo -e "${GREEN}✅ No hardcoded secrets detected${NC}"
fi

# ===========================================
# 3. Check for database files
# ===========================================
echo ""
echo "💾 Checking for database files..."

DB_FILES=$(find . -type f \( -name "*.db" -o -name "*.sqlite" -o -name "*.sqlite3" \) 2>/dev/null)

if [ -n "$DB_FILES" ]; then
    echo -e "${RED}⚠️  Found database files (should be in .gitignore):${NC}"
    echo "$DB_FILES"
    ISSUES_FOUND=1
else
    echo -e "${GREEN}✅ No database files found${NC}"
fi

# ===========================================
# 4. Check for certificate/key files
# ===========================================
echo ""
echo "🔑 Checking for certificates and keys..."

CERT_FILES=$(find . -type f \( -name "*.pem" -o -name "*.key" -o -name "*.cert" -o -name "*.pfx" \) \
    ! -path "*/node_modules/*" \
    ! -path "*/.git/*" \
    2>/dev/null)

if [ -n "$CERT_FILES" ]; then
    echo -e "${RED}⚠️  Found certificate/key files:${NC}"
    echo "$CERT_FILES"
    ISSUES_FOUND=1
else
    echo -e "${GREEN}✅ No certificate/key files found${NC}"
fi

# ===========================================
# Summary
# ===========================================
echo ""
echo "=========================================="
if [ $ISSUES_FOUND -eq 0 ]; then
    echo -e "${GREEN}✅ Security scan passed! Repository is clean.${NC}"
    exit 0
else
    echo -e "${RED}⚠️  Security issues detected! Please review and fix before committing.${NC}"
    echo ""
    echo "Recommended actions:"
    echo "1. Remove or add sensitive files to .gitignore"
    echo "2. Use .env.example files with placeholder values"
    echo "3. Never commit real passwords, tokens, or API keys"
    echo "4. Use environment variables for sensitive data"
    exit 1
fi
