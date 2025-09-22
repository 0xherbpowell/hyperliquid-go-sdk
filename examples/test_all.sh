#!/bin/bash

# Script to test building all examples
# This helps ensure all examples compile correctly

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo "Testing all examples..."
echo "Project root: $PROJECT_ROOT"

cd "$PROJECT_ROOT"

# Counter for tracking results
TOTAL=0
PASSED=0
FAILED=0

# Function to test building an example
test_example() {
    local example_name=$1
    local example_file="examples/${example_name}.go"
    
    if [ ! -f "$example_file" ]; then
        echo "❌ $example_name: File not found"
        return 1
    fi
    
    echo -n "🧪 Testing $example_name... "
    
    if go build -o /tmp/test_example "$example_file" "examples/example_utils.go" 2>/dev/null; then
        echo "✅ PASS"
        ((PASSED++))
        rm -f /tmp/test_example
        return 0
    else
        echo "❌ FAIL"
        echo "   Build errors:"
        go build -o /tmp/test_example "$example_file" "examples/example_utils.go" 2>&1 | sed 's/^/   /'
        echo
        ((FAILED++))
        return 1
    fi
}

# List of examples to test
examples=(
    "basic_agent"
    "basic_order_modify" 
    "basic_market_order"
    "basic_withdraw"
    "basic_trading"
    "websocket_streaming"
    "basic_order"
    "basic_tpsl"
    "basic_spot_order"
    "basic_leverage_adjustment"
    "basic_order_with_cloid"
    "basic_transfer"
    "basic_spot_transfer"
    "cancel_open_orders"
)

echo "Found ${#examples[@]} examples to test"
echo "========================================"

# Test each example
for example in "${examples[@]}"; do
    test_example "$example"
    ((TOTAL++))
done

echo "========================================"
echo "Test Summary:"
echo "  Total:  $TOTAL"
echo "  Passed: $PASSED"
echo "  Failed: $FAILED"

if [ $FAILED -eq 0 ]; then
    echo "🎉 All examples built successfully!"
    exit 0
else
    echo "⚠️  Some examples failed to build"
    exit 1
fi