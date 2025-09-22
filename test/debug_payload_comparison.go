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
	fmt.Println("=== Debugging EIP712 Payload Differences ===")

	// Load config
	privateKeyHex := os.Getenv("HYPERLIQUID_PRIVATE_KEY")
	address := os.Getenv("HYPERLIQUID_ADDRESS")

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

	walletAddress := utils.GetAddressFromPrivateKey(privateKey)
	if address == "" {
		address = walletAddress
	}

	fmt.Printf("Testing with account: %s\n", address)

	// Create exchange client
	timeout := 30 * time.Second
	exchange, err := client.NewExchange(
		privateKey,
		utils.TestnetAPIURL,
		&timeout,
		nil,
		nil,
		&address,
		nil,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to create exchange client: %v", err)
	}

	// Test 1: Create the signing payload manually like Python does
	fmt.Println("\n=== Test 1: Manual EIP712 signing ===")
	
	timestamp := utils.GetTimestampMS()
	
	// Create the action exactly like Python
	action := map[string]interface{}{
		"destination": "0x1234567890123456789012345678901234567890",
		"amount":      "0.01",
		"time":        timestamp,
		"type":        "usdSend",
	}

	fmt.Printf("Action before signing: %+v\n", action)

	// Sign it manually
	signature, err := utils.SignUSDTransferAction(privateKey, action, false)
	if err != nil {
		log.Printf("Signing failed: %v", err)
		return
	}

	fmt.Printf("Signature: %+v\n", signature)
	fmt.Printf("Action after signing: %+v\n", action)

	// Test sending with postAction wrapper
	payload := map[string]interface{}{
		"action":       action,
		"nonce":        timestamp,
		"signature":    signature,
		"vaultAddress": nil,
		"expiresAfter": nil,
	}

	fmt.Println("\nFinal payload structure:")
	jsonBytes, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println(string(jsonBytes))

	result, err := exchange.Post("/exchange", payload)
	if err != nil {
		fmt.Printf("Manual test failed: %v\n", err)
	} else {
		fmt.Printf("Manual test result: %+v\n", result)
	}

	// Test 2: Using our SDK method
	fmt.Println("\n=== Test 2: SDK UsdTransfer method ===")
	
	result2, err := exchange.UsdTransfer("0x1234567890123456789012345678901234567890", "0.01")
	if err != nil {
		fmt.Printf("SDK method failed: %v\n", err)
	} else {
		fmt.Printf("SDK method result: %+v\n", result2)
	}

	// Test 3: Agent approval manual test
	fmt.Println("\n=== Test 3: Manual Agent Approval ===")
	
	// Generate test agent
	agentPrivateKey, err := utils.CreateRandomWallet()
	if err != nil {
		log.Printf("Failed to create agent wallet: %v", err)
		return
	}
	agentAddress := utils.GetAddressFromPrivateKey(agentPrivateKey)
	
	agentNonce := utils.GetTimestampMS()
	
	agentAction := map[string]interface{}{
		"type":         "approveAgent",
		"agentAddress": agentAddress,
		"agentName":    "",
		"nonce":       agentNonce,
	}

	fmt.Printf("Agent action before signing: %+v\n", agentAction)

	agentSignature, err := utils.SignAgent(privateKey, agentAction, false)
	if err != nil {
		log.Printf("Agent signing failed: %v", err)
		return
	}

	fmt.Printf("Agent signature: %+v\n", agentSignature)
	fmt.Printf("Agent action after signing: %+v\n", agentAction)

	// Remove empty agentName like Python does
	if agentAction["agentName"] == "" {
		delete(agentAction, "agentName")
	}

	agentPayload := map[string]interface{}{
		"action":       agentAction,
		"nonce":        agentNonce,
		"signature":    agentSignature,
		"vaultAddress": nil,
		"expiresAfter": nil,
	}

	fmt.Println("\nAgent payload structure:")
	jsonBytes2, _ := json.MarshalIndent(agentPayload, "", "  ")
	fmt.Println(string(jsonBytes2))

	result3, err := exchange.Post("/exchange", agentPayload)
	if err != nil {
		fmt.Printf("Manual agent test failed: %v\n", err)
	} else {
		fmt.Printf("Manual agent test result: %+v\n", result3)
	}
}