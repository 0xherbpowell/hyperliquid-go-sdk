#!/bin/bash

# Script to run individual examples
# Usage: ./run_example.sh <example_name>
# Example: ./run_example.sh basic_order

if [ $# -eq 0 ]; then
    echo "Usage: $0 <example_name>"
    echo "Available examples:"
    echo "  basic_order"
    echo "  basic_trading" 
    echo "  basic_market_order"
    echo "  basic_leverage_adjustment"
    echo "  basic_tpsl"
    echo "  basic_order_modify"
    echo "  basic_order_with_cloid"
    echo "  cancel_open_orders"
    echo "  basic_transfer"
    echo "  basic_withdraw"
    echo "  basic_spot_transfer" 
    echo "  basic_spot_order"
    echo "  websocket_streaming"
    exit 1
fi

EXAMPLE_NAME=$1
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Check if example file exists
if [ ! -f "$SCRIPT_DIR/${EXAMPLE_NAME}.go" ]; then
    echo "Error: Example '$EXAMPLE_NAME' not found at $SCRIPT_DIR/${EXAMPLE_NAME}.go"
    exit 1
fi

# Check if example_utils.go exists
if [ ! -f "$SCRIPT_DIR/example_utils.go" ]; then
    echo "Error: example_utils.go not found at $SCRIPT_DIR/example_utils.go"
    exit 1
fi

echo "Running example: $EXAMPLE_NAME"
echo "Project root: $PROJECT_ROOT"

# Change to project root to ensure proper module resolution
cd "$PROJECT_ROOT"

# Run the example
go run "examples/${EXAMPLE_NAME}.go" "examples/example_utils.go"