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
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}
	}
	file, err := os.Open(configPath)
	if err != nil {
		return &Config{}
	}
	defer file.Close()
	var config Config
	json.NewDecoder(file).Decode(&config)
	return &config
}

func main() {
	fmt.Println("=== Testing User-Signed Actions WITHOUT Wrapper ===")

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
		log.Fatal("Private key not found")
	}

	privateKey, err := utils.ParsePrivateKey(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	walletAddress := utils.GetAddressFromPrivateKey(privateKey)
	if address == "" {
		address = walletAddress
	}

	fmt.Printf("Testing with account: %s\n", address)

	// Create API client
	timeout := 30 * time.Second
	api := client.NewAPI(utils.TestnetAPIURL, &timeout)

	// Test 1: USD Transfer without wrapper
	fmt.Println("\n=== Test 1: USD Transfer (direct) ===")
	
	timestamp := utils.GetTimestampMS()
	
	// Create action with proper type conversion for signing
	signAction := map[string]interface{}{
		"destination": "0x1234567890123456789012345678901234567890",
		"amount":      "0.01",
		"time":        fmt.Sprintf("%d", timestamp), // String for EIP712
	}

	signature, err := utils.SignUSDTransferAction(privateKey, signAction, false)
	if err != nil {
		log.Printf("USD signing failed: %v", err)
		return
	}

	// Send direct payload (no wrapper)
	directPayload := map[string]interface{}{
		"type":        "usdSend",
		"destination": "0x1234567890123456789012345678901234567890",
		"amount":      "0.01",
		"time":        timestamp, // int64 for API
		"signature":   signature,
	}

	fmt.Println("Direct USD payload:")
	jsonBytes, _ := json.MarshalIndent(directPayload, "", "  ")
	fmt.Println(string(jsonBytes))

	result, err := api.Post("/exchange", directPayload)
	if err != nil {
		fmt.Printf("Direct USD transfer failed: %v\n", err)
	} else {
		fmt.Printf("Direct USD transfer result: %+v\n", result)
	}

	// Test 2: Agent approval without wrapper
	fmt.Println("\n=== Test 2: Agent Approval (direct) ===")
	
	// Generate test agent
	agentPrivateKey, err := utils.CreateRandomWallet()
	if err != nil {
		log.Printf("Failed to create agent wallet: %v", err)
		return
	}
	agentAddress := utils.GetAddressFromPrivateKey(agentPrivateKey)
	
	agentNonce := utils.GetTimestampMS()
	
	// Create action for signing with proper type conversion
	agentSignAction := map[string]interface{}{
		"agentAddress": agentAddress,
		"agentName":    "",
		"nonce":       fmt.Sprintf("%d", agentNonce), // String for EIP712
	}

	agentSignature, err := utils.SignAgent(privateKey, agentSignAction, false)
	if err != nil {
		log.Printf("Agent signing failed: %v", err)
		return
	}

	// Send direct payload (no wrapper)
	directAgentPayload := map[string]interface{}{
		"type":         "approveAgent",
		"agentAddress": agentAddress,
		"nonce":       agentNonce, // int64 for API
		"signature":   agentSignature,
	}
	// Don't include agentName if empty

	fmt.Println("Direct agent payload:")
	jsonBytes2, _ := json.MarshalIndent(directAgentPayload, "", "  ")
	fmt.Println(string(jsonBytes2))

	result2, err := api.Post("/exchange", directAgentPayload)
	if err != nil {
		fmt.Printf("Direct agent approval failed: %v\n", err)
	} else {
		fmt.Printf("Direct agent approval result: %+v\n", result2)
	}
}