package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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
	fmt.Println("Testing UsdTransfer signing to compare with agent signing...")
	
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

	fmt.Println("\n=== Testing UsdTransfer signing ===")
	
	// Test UsdTransfer signing manually (similar to what's done in exchange.go)
	timestamp := utils.GetTimestampMS()
	
	action := map[string]interface{}{
		"destination": "0x1234567890123456789012345678901234567890",
		"amount":      "1.0",
		"time":        uint64(timestamp), // This should work if uint64 is handled correctly
	}

	fmt.Printf("UsdTransfer action: %+v\n", action)
	fmt.Printf("timestamp type: %T, value: %v\n", action["time"], action["time"])

	_, err = utils.SignUSDTransferAction(privateKey, action, false)
	if err != nil {
		fmt.Printf("UsdTransfer signing failed: %v\n", err)
	} else {
		fmt.Printf("UsdTransfer signing successful!\n")
	}

	fmt.Println("\n=== Testing Agent signing ===")
	
	// Now test agent signing
	agentPrivateKey, err := utils.CreateRandomWallet()
	if err != nil {
		log.Fatalf("Failed to create agent wallet: %v", err)
	}
	agentAddress := utils.GetAddressFromPrivateKey(agentPrivateKey)
	
	nonce := uint64(timestamp)
	
	agentAction := map[string]interface{}{
		"agentAddress": agentAddress,
		"agentName":    "",
		"nonce":       nonce,
	}

	fmt.Printf("Agent action: %+v\n", agentAction)
	fmt.Printf("nonce type: %T, value: %v\n", agentAction["nonce"], agentAction["nonce"])

	_, err = utils.SignAgent(privateKey, agentAction, false)
	if err != nil {
		fmt.Printf("Agent signing failed: %v\n", err)
	} else {
		fmt.Printf("Agent signing successful!\n")
	}

	// Compare the signing type definitions
	fmt.Println("\n=== Comparing type definitions ===")
	fmt.Println("USDSendSignTypes should include:")
	fmt.Println("  - hyperliquidChain (string)")
	fmt.Println("  - destination (string)")
	fmt.Println("  - amount (string)")
	fmt.Println("  - time (uint64)")
	
	fmt.Println("AgentSignTypes should include:")
	fmt.Println("  - hyperliquidChain (string)")
	fmt.Println("  - agentAddress (address)")
	fmt.Println("  - agentName (string)")
	fmt.Println("  - nonce (uint64)")
	
	fmt.Println("\nBoth use uint64, so if one works and the other doesn't, there may be a difference in how they're processed.")
}