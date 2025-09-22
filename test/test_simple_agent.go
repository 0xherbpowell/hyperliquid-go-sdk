package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"hyperliquid-go-sdk/pkg/client"
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
	// Load configuration
	fmt.Println("Testing agent approval with simple nonce...")
	
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

	// Use wallet address as the account if they differ
	if address != walletAddress {
		fmt.Println("Using wallet address as account address")
		address = walletAddress
	}

	fmt.Printf("Account: %s\n", address)

	// Create exchange client
	timeout := 30 * time.Second
	exchange, err := client.NewExchange(
		privateKey,
		utils.TestnetAPIURL,
		&timeout,
		nil,      // meta
		nil,      // vault address
		&address, // account address
		nil,      // spot meta
		nil,      // perp dexs
	)
	if err != nil {
		log.Fatalf("Failed to create exchange client: %v", err)
	}

	fmt.Println("\n=== Testing with simple nonce values ===")
	
	// Generate agent address
	agentPrivateKey, err := utils.CreateRandomWallet()
	if err != nil {
		log.Fatalf("Failed to create agent wallet: %v", err)
	}
	agentAddress := utils.GetAddressFromPrivateKey(agentPrivateKey)
	fmt.Printf("Agent address: %s\n", agentAddress)

	// Test with simple nonce values
	testNonces := []uint64{1, 100, 1000, 12345}
	
	for _, nonce := range testNonces {
		fmt.Printf("\n--- Testing with nonce: %d ---\n", nonce)
		
		action := map[string]interface{}{
			"agentAddress": agentAddress,
			"agentName":    "",
			"nonce":       nonce,
		}

		_, err := utils.SignAgent(privateKey, action, false)
		if err != nil {
			fmt.Printf("   Signing failed: %v\n", err)
		} else {
			fmt.Printf("   Signing successful!\n")
			
			// Try the actual ApproveAgent call now
			fmt.Printf("   Testing full workflow...\n")
			
			result, err := exchange.ApproveAgent()
			if err != nil {
				fmt.Printf("   ApproveAgent failed: %v\n", err)
			} else {
				fmt.Printf("   ApproveAgent successful!\n")
				fmt.Printf("   Result status: %v\n", result.Result["status"])
				return // Exit after first successful attempt
			}
			return // Exit after first successful signing
		}
	}

	fmt.Println("\nAll simple nonces failed. The issue may be elsewhere.")
}