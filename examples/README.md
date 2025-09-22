# Hyperliquid Go SDK Examples

This directory contains example code demonstrating how to use the Hyperliquid Go SDK.

## Setup

1. **Environment Variables** (Recommended):
   ```bash
   export HYPERLIQUID_PRIVATE_KEY="your_private_key_here"
   export HYPERLIQUID_ADDRESS="your_address_here"  # Optional, will be derived from private key
   ```

2. **Config File** (Alternative):
   Create a `config.json` file in the project root:
   ```json
   {
     "secret_key": "your_private_key_here",
     "account_address": "your_address_here"
   }
   ```

## Building and Running Examples

Each example is a standalone program. You can run them individually:

### Method 1: Run directly with go run
```bash
# Run a specific example
go run examples/basic_order.go examples/example_utils.go

# Run another example
go run examples/basic_trading.go examples/example_utils.go
```

### Method 2: Build and run
```bash
# Build a specific example
go build -o basic_order examples/basic_order.go examples/example_utils.go

# Run the built executable
./basic_order
```

### Method 3: Using the run script
```bash
# Make the script executable
chmod +x examples/run_example.sh

# Run a specific example
./examples/run_example.sh basic_order
./examples/run_example.sh basic_trading
```

## Available Examples

### Basic Trading
- `basic_order.go` - Place simple limit and market orders
- `basic_trading.go` - Comprehensive trading example with multiple operations
- `basic_market_order.go` - Market order examples with slippage protection
- `basic_leverage_adjustment.go` - Adjust leverage settings
- `basic_tpsl.go` - Take profit and stop loss orders

### Order Management
- `basic_order_modify.go` - Modify existing orders
- `basic_order_with_cloid.go` - Orders with client IDs and cancellation by client ID
- `cancel_open_orders.go` - Cancel all open orders

### Transfers and Withdrawals
- `basic_transfer.go` - USD transfers between addresses
- `basic_withdraw.go` - Withdraw funds from the exchange
- `basic_spot_transfer.go` - Spot asset transfers
- `basic_spot_order.go` - Spot trading examples

### Advanced Features
- `basic_agent.go` - Agent/sub-account concepts and separate wallet management
- `websocket_streaming.go` - Real-time data streaming

### Testing
- `test_all.sh` - Script to test building all examples
- `run_example.sh` - Script to run individual examples

## Important Notes

1. **Testnet vs Mainnet**: All examples use testnet by default for safety
2. **Account Balance**: Examples check for account balance before trading
3. **Error Handling**: Examples include comprehensive error handling
4. **Cleanup**: Many examples cancel orders after placing them for cleanup

## Safety Features

- All examples use testnet by default
- Account balance validation before trading
- Order size limits for safety
- Comprehensive error handling and logging

## Troubleshooting

### Common Issues

1. **"No equity" error**: Make sure your testnet account has funds
2. **Build errors**: Ensure you're including `example_utils.go` when building
3. **Private key errors**: Check your private key format (should be hex without 0x prefix)

### Getting Testnet Funds

Visit the Hyperliquid testnet to get test funds for your account.

## Contributing

When adding new examples:
1. Follow the existing pattern in `example_utils.go`
2. Include proper error handling
3. Add cleanup code for orders/positions
4. Update this README with the new example