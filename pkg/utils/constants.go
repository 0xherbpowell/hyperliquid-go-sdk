package utils

const (
	// API URLs
	MainnetAPIURL = "https://api.hyperliquid.xyz"
	TestnetAPIURL = "https://api.hyperliquid-testnet.xyz"
	LocalAPIURL   = "http://localhost:3001"

	// WebSocket URLs
	MainnetWSURL = "wss://api.hyperliquid.xyz/ws"
	TestnetWSURL = "wss://api.hyperliquid-testnet.xyz/ws"

	// Chain configurations
	MainnetChainName = "Mainnet"
	TestnetChainName = "Testnet"

	// Signature configurations
	SignatureChainID = "0x66eee"
	EIP712ChainID    = 1337  // EIP712 chain ID for Hyperliquid signing (matches Python SDK)

	// Decimal places
	USDDecimals = 6
	SzDecimals  = 8

	// Default timeouts
	DefaultTimeoutSeconds = 30
)
