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

func main() {
	fmt.Println("=== Debugging Key Derivation ===")

	// Try environment variables
	envKey := os.Getenv("HYPERLIQUID_PRIVATE_KEY")
	envAddress := os.Getenv("HYPERLIQUID_ADDRESS")
	
	fmt.Printf("Environment HYPERLIQUID_PRIVATE_KEY: %s\n", envKey)
	fmt.Printf("Environment HYPERLIQUID_ADDRESS: %s\n", envAddress)

	// Try config file
	var config Config
	if file, err := os.Open("./config.json"); err == nil {
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&config); err == nil {
			fmt.Printf("Config secret_key: %s\n", config.SecretKey)
			fmt.Printf("Config account_address: %s\n", config.AccountAddress)
		} else {
			fmt.Printf("Config decode error: %v\n", err)
		}
	} else {
		fmt.Printf("Config file error: %v\n", err)
	}

	// Use the key we would actually use
	privateKeyHex := envKey
	if privateKeyHex == "" {
		privateKeyHex = config.SecretKey
	}

	if privateKeyHex == "" {
		log.Fatal("No private key found")
	}

	// Parse and derive address
	privateKey, err := utils.ParsePrivateKey(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	derivedAddress := utils.GetAddressFromPrivateKey(privateKey)
	fmt.Printf("Derived address from private key: %s\n", derivedAddress)

	// Compare with expected
	expectedAddress := envAddress
	if expectedAddress == "" {
		expectedAddress = config.AccountAddress
	}
	
	if expectedAddress != "" {
		fmt.Printf("Expected address: %s\n", expectedAddress)
		if derivedAddress == expectedAddress {
			fmt.Println("✅ Addresses match!")
		} else {
			fmt.Println("❌ Address mismatch!")
		}
	}

	// Print first few characters of private key for debugging (be careful!)
	if len(privateKeyHex) > 10 {
		fmt.Printf("Private key prefix: %s...\n", privateKeyHex[:10])
	}
}