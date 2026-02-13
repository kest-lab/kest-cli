#!/bin/bash
# Kest Project Cleanup Script
# 清理项目中不应提交到开源仓库的文件

set -e

echo "🧹 Kest Project Cleanup Script"
echo "================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to remove file/directory if exists
remove_if_exists() {
    if [ -e "$1" ]; then
        echo -e "${YELLOW}Removing:${NC} $1"
        rm -rf "$1"
        echo -e "${GREEN}✓ Removed${NC}"
    else
        echo -e "${GREEN}✓ Not found (already clean):${NC} $1"
    fi
}

echo "Step 1: Removing compiled binaries..."
echo "--------------------------------------"
remove_if_exists "api/main"
remove_if_exists "api/test-server"
remove_if_exists "api/test-server-prod"
remove_if_exists "api/kest-api"
remove_if_exists "api/server"
echo ""

echo "Step 2: Removing log files..."
echo "--------------------------------------"
remove_if_exists "api/server.log"
remove_if_exists "api/test-server.log"
remove_if_exists "api/storage/logs"
remove_if_exists "api/tmp/build-errors.log"
echo ""

echo "Step 3: Removing internal documentation..."
echo "--------------------------------------"
remove_if_exists "API_REVIEW_REPORT.md"
remove_if_exists "FINAL_COMPLETION_REPORT.md"
remove_if_exists "api/AGENTS.md"
remove_if_exists "api/DOCKERFILE_FIX.md"
remove_if_exists "api/PLATFORM_DESIGN_V2.md"
remove_if_exists "api/PLATFORM_IMPLEMENTATION.md"
remove_if_exists "api/WEB_UI_DESIGN.md"
remove_if_exists "api/ZEABUR_DEPLOYMENT.md"
remove_if_exists "web/ADMIN_UI_INTEGRATION.md"
remove_if_exists "web/CLI_INTEGRATION_PLAN.md"
remove_if_exists "web/COMPLETE_SUMMARY.md"
remove_if_exists "web/FINAL_STATUS.md"
remove_if_exists "web/FINAL_TEST_REPORT.md"
remove_if_exists "web/IMPLEMENTATION.md"
remove_if_exists "web/IMPLEMENTATION_SUMMARY.md"
remove_if_exists "web/TEST_AUTH_TOKEN.md"
remove_if_exists "web/TEST_REPORT.md"
echo ""

echo "Step 4: Checking for sensitive files..."
echo "--------------------------------------"
if [ -f "api/.env" ]; then
    echo -e "${RED}⚠️  WARNING: api/.env exists (contains sensitive data)${NC}"
    echo "   This file should NOT be committed. Please verify it's in .gitignore"
else
    echo -e "${GREEN}✓ api/.env not found (good)${NC}"
fi

if [ -f "web/.env" ]; then
    echo -e "${RED}⚠️  WARNING: web/.env exists${NC}"
    echo "   This file should NOT be committed. Please verify it's in .gitignore"
else
    echo -e "${GREEN}✓ web/.env not found (good)${NC}"
fi
echo ""

echo "Step 5: Removing empty directories..."
echo "--------------------------------------"
find . -type d -empty -not -path "./.git/*" -delete 2>/dev/null || true
echo -e "${GREEN}✓ Empty directories removed${NC}"
echo ""

echo "Step 6: Git status check..."
echo "--------------------------------------"
if command -v git &> /dev/null; then
    echo "Files that will be removed from git (if tracked):"
    git ls-files --deleted 2>/dev/null || echo "No deleted files"
    echo ""
    echo "Untracked files:"
    git ls-files --others --exclude-standard 2>/dev/null || echo "No untracked files"
else
    echo "Git not found, skipping git status check"
fi
echo ""

echo "================================"
echo -e "${GREEN}✅ Cleanup completed!${NC}"
echo ""
echo "Next steps:"
echo "1. Review the changes: git status"
echo "2. Check .env files are not tracked: git ls-files | grep .env"
echo "3. Commit the cleanup: git add -A && git commit -m 'chore: cleanup for open source release'"
echo "4. Push to repository: git push"
echo ""
