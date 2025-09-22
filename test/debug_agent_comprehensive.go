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
	fmt.Println("Comprehensive debugging of agent approval...")
	
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

	// Let's try to understand the agent approval by comparing with a working method first
	fmt.Println("\n=== Testing UsdTransfer (should work) ===")
	// Note: This will fail because we don't have enough funds, but it should not have deserialization errors
	result, err := exchange.UsdTransfer("0x1234567890123456789012345678901234567890", "0.01")
	if err != nil {
		fmt.Printf("UsdTransfer error (expected): %v\n", err)
	} else {
		fmt.Printf("UsdTransfer unexpectedly succeeded: %+v\n", result)
	}

	// Now let's try to manually construct the agent approval with different approaches
	fmt.Println("\n=== Approach 1: Agent approval through different method ===")
	
	// Maybe agent approval should be treated as an L1 action instead of user-signed action?
	// Let's try creating it manually with L1 signing
	timestamp := utils.GetTimestampMS()
	
	// Generate agent address
	agentPrivateKey, err := utils.CreateRandomWallet()
	if err != nil {
		log.Fatalf("Failed to create agent wallet: %v", err)
	}
	agentAddress := utils.GetAddressFromPrivateKey(agentPrivateKey)
	fmt.Printf("Test agent address: %s\n", agentAddress)

	// Try L1 action approach
	l1Action := map[string]interface{}{
		"type":         "approveAgent",
		"agentAddress": agentAddress,
		"agentName":    "",
		"nonce":       timestamp,
	}

	fmt.Printf("Trying L1 action approach...\n")
	signature, err := utils.SignL1Action(
		privateKey,
		l1Action,
		nil,    // vault address
		timestamp,
		nil,    // expires after
		false,  // testnet
	)
	if err != nil {
		fmt.Printf("L1 signing failed: %v\n", err)
	} else {
		fmt.Printf("L1 signing succeeded, now testing API call...\n")
		
		// Create a custom postAction call to test this approach
		payload := map[string]interface{}{
			"action":       l1Action,
			"nonce":        timestamp,
			"signature":    signature,
			"vaultAddress": nil,
			"expiresAfter": nil,
		}
		
		fmt.Printf("L1 Payload: %+v\n", payload)
		result, err := exchange.Post("/exchange", payload)
		if err != nil {
			fmt.Printf("L1 API call failed: %v\n", err)
		} else {
			fmt.Printf("L1 API call succeeded: %+v\n", result)
		}
	}

	fmt.Println("\n=== Approach 2: User-signed without top-level nonce ===")
	// Maybe the issue is that user-signed actions don't need the top-level nonce
	// Try building the payload manually without the wrapper
	
	userSignedAction := map[string]interface{}{
		"agentAddress": agentAddress,
		"agentName":    "",
		"nonce":       fmt.Sprintf("%d", timestamp),
	}
	
	userSignature, err := utils.SignAgent(privateKey, userSignedAction, false)
	if err != nil {
		fmt.Printf("User-signed approach failed: %v\n", err)
	} else {
		fmt.Printf("User-signed approach succeeded, testing direct API call...\n")
		
		// Try sending just the signed action directly
		directPayload := map[string]interface{}{
			"type":         "approveAgent",
			"agentAddress": agentAddress,
			"agentName":    "",
			"nonce":       timestamp, // Use int64 for API
			"signature":   userSignature,
		}
		
		fmt.Printf("Direct payload: %+v\n", directPayload)
		result, err := exchange.Post("/exchange", directPayload)
		if err != nil {
			fmt.Printf("Direct API call failed: %v\n", err)
		} else {
			fmt.Printf("Direct API call succeeded: %+v\n", result)
		}
	}

	fmt.Println("\n=== Approach 3: Check existing ApproveAgent method ===")
	approveResult, err := exchange.ApproveAgent()
	if err != nil {
		fmt.Printf("ApproveAgent method error: %v\n", err)
	} else {
		fmt.Printf("ApproveAgent method succeeded: %+v\n", approveResult)
	}
}