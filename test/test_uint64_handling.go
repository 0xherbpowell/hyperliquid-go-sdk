package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

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
	fmt.Println("Testing different ways to handle uint64 values in EIP712 signing...")
	
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

	// Generate agent for testing
	agentPrivateKey, err := utils.CreateRandomWallet()
	if err != nil {
		log.Fatalf("Failed to create agent wallet: %v", err)
	}
	agentAddress := utils.GetAddressFromPrivateKey(agentPrivateKey)
	
	timestamp := utils.GetTimestampMS()

	// Test 1: uint64 (current approach - failing)
	fmt.Println("\n=== Test 1: uint64 value ===")
	action1 := map[string]interface{}{
		"agentAddress": agentAddress,
		"agentName":    "",
		"nonce":       uint64(timestamp),
	}
	fmt.Printf("Action: %+v\n", action1)
	_, err = utils.SignAgent(privateKey, action1, false)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Success!\n")
	}

	// Test 2: int64 value 
	fmt.Println("\n=== Test 2: int64 value ===")
	action2 := map[string]interface{}{
		"agentAddress": agentAddress,
		"agentName":    "",
		"nonce":       timestamp, // int64
	}
	fmt.Printf("Action: %+v\n", action2)
	_, err = utils.SignAgent(privateKey, action2, false)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Success!\n")
	}

	// Test 3: string value 
	fmt.Println("\n=== Test 3: string value ===")
	action3 := map[string]interface{}{
		"agentAddress": agentAddress,
		"agentName":    "",
		"nonce":       strconv.FormatInt(timestamp, 10), // string
	}
	fmt.Printf("Action: %+v\n", action3)
	_, err = utils.SignAgent(privateKey, action3, false)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Success!\n")
	}

	// Test 4: Small uint64 value 
	fmt.Println("\n=== Test 4: Small uint64 value (123) ===")
	action4 := map[string]interface{}{
		"agentAddress": agentAddress,
		"agentName":    "",
		"nonce":       uint64(123),
	}
	fmt.Printf("Action: %+v\n", action4)
	_, err = utils.SignAgent(privateKey, action4, false)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Success!\n")
	}

	// Test 5: JSON Number (float64) 
	fmt.Println("\n=== Test 5: float64 value ===")
	action5 := map[string]interface{}{
		"agentAddress": agentAddress,
		"agentName":    "",
		"nonce":       float64(timestamp),
	}
	fmt.Printf("Action: %+v\n", action5)
	_, err = utils.SignAgent(privateKey, action5, false)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Success!\n")
	}
}