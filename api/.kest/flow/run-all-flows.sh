#!/bin/bash

# Script to run all Kest flow tests
# Usage: ./run-all-flows.sh

set -e

FLOW_DIR=".kest/flow"
FAILED_TESTS=()
PASSED_TESTS=()

echo "🦅 Running Kest Flow Tests"
echo "=========================="
echo ""

# Run each flow file
for flow_file in "$FLOW_DIR"/*.flow.md; do
    if [ -f "$flow_file" ]; then
        filename=$(basename "$flow_file")
        echo "📝 Running: $filename"
        echo "---"
        
        if kest run "$flow_file"; then
            PASSED_TESTS+=("$filename")
            echo "✅ PASSED: $filename"
        else
            FAILED_TESTS+=("$filename")
            echo "❌ FAILED: $filename"
        fi
        
        echo ""
        echo "---"
        echo ""
    fi
done

# Summary
echo "=========================="
echo "📊 Test Summary"
echo "=========================="
echo "✅ Passed: ${#PASSED_TESTS[@]}"
echo "❌ Failed: ${#FAILED_TESTS[@]}"
echo ""

if [ ${#PASSED_TESTS[@]} -gt 0 ]; then
    echo "Passed tests:"
    for test in "${PASSED_TESTS[@]}"; do
        echo "  ✓ $test"
    done
    echo ""
fi

if [ ${#FAILED_TESTS[@]} -gt 0 ]; then
    echo "Failed tests:"
    for test in "${FAILED_TESTS[@]}"; do
        echo "  ✗ $test"
    done
    echo ""
    exit 1
fi

echo "🎉 All tests passed!"
exit 0
