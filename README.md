# hyperliquid-go-sdk

<div align="center">

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/doc/install)
[![Go Report Card](https://goreportcard.com/badge/github.com/hyperliquid-dex/hyperliquid-go-sdk)](https://goreportcard.com/report/github.com/hyperliquid-dex/hyperliquid-go-sdk)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.md)

SDK for Hyperliquid API trading with Go.

</div>

## Installation

```bash
go get github.com/hyperliquid-dex/hyperliquid-go-sdk
```

## Quick Start

### Configuration

1. Copy the example configuration file:
```bash
cp examples/config.json.example examples/config.json
```

2. Edit `examples/config.json` and add your private key:
```json
{
    "secret_key": "0x1234567890abcdef...",
    "account_address": "0x1234567890abcdef..."
}
```

**Note:**
- Set the public key as the `account_address`
- Set your private key as the `secret_key`
- If using an API wallet, set the API wallet's private key as `secret_key` and the main wallet's public address as `account_address`

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/hyperliquid-dex/hyperliquid-go-sdk/pkg/client"
    "github.com/hyperliquid-dex/hyperliquid-go-sdk/pkg/constants"
    "github.com/hyperliquid-dex/hyperliquid-go-sdk/pkg/types"
)

func main() {
    ctx := context.Background()
    
    // Create info client (read-only)
    infoClient, err := client.NewInfoClient(constants.TestnetAPIURL, true, nil, nil, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Get user state
    userState, err := infoClient.UserState(ctx, "0xcd5051944f780a621ee62e39e493c489668acf4d", "")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Account Value: %s\n", userState.MarginSummary.AccountValue)
}
```

### Trading Example

```go
package main

import (
    "context"
    "crypto/ecdsa"
    "log"

    "github.com/ethereum/go-ethereum/crypto"
    "github.com/hyperliquid-dex/hyperliquid-go-sdk/pkg/client"
    "github.com/hyperliquid-dex/hyperliquid-go-sdk/pkg/constants"
    "github.com/hyperliquid-dex/hyperliquid-go-sdk/pkg/types"
)

func main() {
    ctx := context.Background()
    
    // Load private key
    privateKey, err := crypto.HexToECDSA("your_private_key_without_0x")
    if err != nil {
        log.Fatal(err)
    }

    // Create exchange client
    exchangeClient, err := client.NewExchangeClient(
        privateKey,
        constants.TestnetAPIURL,
        nil, // meta
        nil, // vault address  
        nil, // account address (will be derived from private key)
        nil, // spot meta
        nil, // perp dexs
        nil, // timeout
    )
    if err != nil {
        log.Fatal(err)
    }

    // Place a limit order
    order := types.OrderRequest{
        Coin:       "ETH",
        IsBuy:      true,
        Sz:         0.1,
        LimitPx:    2000.0,
        OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: constants.TifGtc}},
        ReduceOnly: false,
    }

    result, err := exchangeClient.Order(ctx, order, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Order result: %+v\n", result)
}
```

## Examples

The repository includes several example programs:

### Available Examples

- **basic_order**: Place and cancel a basic limit order
- **basic_market_order**: Place market orders to open and close positions
- **cancel_open_orders**: Cancel all open orders for an account

### Running Examples

1. Set up your configuration:
```bash
make setup-config
# Edit examples/config.json with your details
```

2. Run an example:
```bash
make run-basic-order
# or
make run-market-order
```

3. Or run directly:
```bash
cd examples/basic_order
go run main.go
```

## API Reference

### Info Client (Read-only operations)

```go
infoClient, err := client.NewInfoClient(baseURL, skipWS, meta, spotMeta, perpDexs)

// Get user trading state
userState, err := infoClient.UserState(ctx, address, dex)

// Get open orders
openOrders, err := infoClient.OpenOrders(ctx, address, dex)

// Get all mid prices
allMids, err := infoClient.AllMids(ctx, dex)

// Get user fills
fills, err := infoClient.UserFills(ctx, address)

// Get L2 order book
l2Book, err := infoClient.L2Snapshot(ctx, "ETH")

// Query order by ID
orderStatus, err := infoClient.QueryOrderByOid(ctx, address, orderId)
```

### Exchange Client (Trading operations)

```go
exchangeClient, err := client.NewExchangeClient(privateKey, baseURL, meta, vaultAddress, accountAddress, spotMeta, perpDexs, timeout)

// Place orders
result, err := exchangeClient.Order(ctx, orderRequest, builder)
result, err := exchangeClient.BulkOrders(ctx, orderRequests, builder)

// Market orders
result, err := exchangeClient.MarketOpen(ctx, "ETH", true, 0.1, nil, &slippage, cloid, builder)
result, err := exchangeClient.MarketClose(ctx, "ETH", nil, nil, &slippage, cloid, builder)

// Cancel orders
result, err := exchangeClient.Cancel(ctx, "ETH", orderId)
result, err := exchangeClient.CancelByCloid(ctx, "ETH", cloid)
result, err := exchangeClient.BulkCancel(ctx, cancelRequests)

// Update leverage
result, err := exchangeClient.UpdateLeverage(ctx, 10, "ETH", true)

// Transfers
result, err := exchangeClient.USDTransfer(ctx, 100.0, "0x...")
result, err := exchangeClient.USDClassTransfer(ctx, 100.0, true) // perp to spot
```

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile commands)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/hyperliquid-dex/hyperliquid-go-sdk.git
cd hyperliquid-go-sdk
```

2. Install dependencies:
```bash
make deps
```

3. Install development tools:
```bash
make install-tools
```

### Available Make Commands

```bash
make help          # Show available commands
make build         # Build all examples
make test          # Run tests
make test-cover    # Run tests with coverage
make fmt           # Format code
make lint          # Run linter
make vet           # Run go vet
make check         # Run all checks (fmt, vet, lint)
make clean         # Clean build artifacts
make tidy          # Tidy go modules
make all           # Run full build pipeline
```

### Project Structure

```
hyperliquid-go-sdk/
├── pkg/
│   ├── client/           # Client implementations
│   │   ├── client.go     # Base client
│   │   ├── info.go       # Info client (read-only)
│   │   └── exchange.go   # Exchange client (trading)
│   ├── types/            # Type definitions
│   │   ├── common.go     # Common types
│   │   ├── orders.go     # Order-related types  
│   │   ├── cloid.go      # Client order ID type
│   │   └── websocket.go  # WebSocket types
│   ├── signing/          # Cryptographic signing
│   │   ├── signing.go    # Core signing logic
│   │   ├── actions.go    # Action type definitions
│   │   └── utils.go      # Utility functions
│   ├── constants/        # Constants and enums
│   └── errors/           # Error types
├── internal/
│   └── api/              # Internal API client
├── examples/             # Example programs
├── tests/                # Test files
└── docs/                 # Documentation
```

## Testing

Run tests with:
```bash
make test
```

Run tests with coverage:
```bash
make test-cover
```

## WebSocket Support

WebSocket functionality is available but currently requires manual implementation. The SDK provides the necessary types and interfaces in `pkg/types/websocket.go`.

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and add tests
4. Run the full test suite: `make all`
5. Commit your changes: `git commit -am 'Add feature'`
6. Push to the branch: `git push origin feature-name`
7. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.

## Support

- [Documentation](https://hyperliquid.gitbook.io/hyperliquid-docs/)
- [Discord](https://discord.gg/hyperliquid)
- [GitHub Issues](https://github.com/hyperliquid-dex/hyperliquid-go-sdk/issues)

## Disclaimer

This SDK is provided as-is and is not officially supported by Hyperliquid. Use at your own risk. Always test on testnet before using with real funds.