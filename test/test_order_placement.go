package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"hyperliquid-go-sdk/pkg/client"
	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
)

// Config represents the configuration structure
type Config struct {
	SecretKey      string `json:"secret_key"`
	AccountAddress string `json:"account_address"`
}

// loadConfig loads configuration from config.json file
func loadConfig() *Config {
	configPath := "./config.json"

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("config.json not found. Please set environment variables or create config.json")
		return &Config{}
	}

	file, err := os.Open(configPath)
	if err != nil {
		log.Printf("Error opening config file: %v", err)
		return &Config{}
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Printf("Error decoding config file: %v", err)
		return &Config{}
	}

	return &config
}

func main() {
	// Try to get private key from environment variable first
	privateKeyHex := os.Getenv("HYPERLIQUID_PRIVATE_KEY")
	address := os.Getenv("HYPERLIQUID_ADDRESS")
	
	// If not found in environment, try to read from config file
	if privateKeyHex == "" {
		config := loadConfig()
		privateKeyHex = config.SecretKey
		if address == "" {
			address = config.AccountAddress
		}
	}

	if privateKeyHex == "" {
		log.Fatal("Private key not found. Set HYPERLIQUID_PRIVATE_KEY environment variable or update config.json")
	}

	// Parse private key
	privateKey, err := utils.ParsePrivateKey(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Get address from private key if not provided
	walletAddress := utils.GetAddressFromPrivateKey(privateKey)
	if address == "" {
		address = walletAddress
	}

	fmt.Printf("Running with account address: %s\n", address)
	if address != walletAddress {
		fmt.Printf("Using agent wallet: %s\n", walletAddress)
	}

	// Use testnet
	timeout := 30 * time.Second
	
	// Create exchange client - try with main account first
	exchange, err := client.NewExchange(
		privateKey,
		utils.TestnetAPIURL,
		&timeout,
		nil,      // meta
		nil,      // vault address
		&address,       // account address
		nil,      // spot meta
		nil,      // perp dexs
	)
	if err != nil {
		log.Fatalf("Failed to create exchange client: %v", err)
	}

	fmt.Println("Exchange client created successfully")

	// Create a simple limit order for ETH
	// Using hardcoded asset ID for ETH (typically 0)
	orderType := types.OrderType{
		Limit: &types.LimitOrderType{
			Tif: types.TifGtc,
		},
	}

	fmt.Println("Attempting to place order...")

	// Place order with very low price to avoid accidental execution
	result, err := exchange.Order(
		"ETH",          // coin
		true,           // isBuy
		0.01,           // size
		100.0,          // limit price (very low to ensure it doesn't execute)
		orderType,      // order type
		false,          // reduce only
		nil,            // cloid
		nil,            // builder info
	)

	if err != nil {
		log.Printf("Order placement failed: %v", err)
		// This is expected - let's see what the error is
		fmt.Printf("Error type: %T\n", err)
		fmt.Printf("Error message: %s\n", err.Error())
	} else {
		fmt.Println("Order placed successfully!")
		fmt.Printf("Result: %+v\n", result)
	}

	fmt.Println("Test completed.")
}