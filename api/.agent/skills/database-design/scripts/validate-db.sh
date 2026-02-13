#!/bin/bash

# Database Standards Validation Script
# Usage: ./validate-db.sh <path_to_model_file>

set -e

MODEL_FILE=$1

if [ -z "$MODEL_FILE" ]; then
    echo "Usage: ./validate-db.sh <path_to_model_file>"
    exit 1
fi

if [ ! -f "$MODEL_FILE" ]; then
    echo "❌ Model file not found: $MODEL_FILE"
    exit 1
fi

echo "🔍 Validating Database standards for '$(basename $MODEL_FILE)'..."
echo "=============================================="

ERRORS=0
WARNINGS=0

# 1. Check for PO suffix on struct names
if grep -q "type.*struct" "$MODEL_FILE" && ! grep -q "PO struct" "$MODEL_FILE"; then
    echo "❌ Struct definitions in model.go should have 'PO' suffix (e.g., UserPO)."
    ERRORS=$((ERRORS + 1))
else
    echo "✅ Struct naming conventions passed."
fi

# 2. Check for TableName() method
if ! grep -q "func.*TableName().*string" "$MODEL_FILE"; then
    echo "⚠️  Missing TableName() method. Explicit table names are recommended."
    WARNINGS=$((WARNINGS + 1))
else
    echo "✅ TableName() method detected."
fi

# 3. Check for mandatory audit fields
for field in "ID" "CreatedAt" "UpdatedAt" "DeletedAt"; do
    if ! grep -q "$field" "$MODEL_FILE"; then
        echo "❌ Missing mandatory field: $field"
        ERRORS=$((ERRORS + 1))
    fi
done

if [ $ERRORS -eq 0 ]; then
    echo "✅ All mandatory audit fields present."
fi

# 4. Check for snake_case in gorm labels
if grep "gorm:\"" "$MODEL_FILE" | grep -v "primaryKey" | grep -q "[A-Z]"; then
    # This is a very simple check that might have false positives, but it targets mixedCase in tags
    echo "⚠️  Detected possible non-snake_case naming in GORM tags. Check column names."
    WARNINGS=$((WARNINGS + 1))
fi

echo "=============================================="
if [ $ERRORS -eq 0 ]; then
    echo "SUCCESS: Database standards mostly met ($WARNINGS warnings)."
    exit 0
else
    echo "FAILURE: Found $ERRORS errors. Please fix before proceeding."
    exit 1
fi
