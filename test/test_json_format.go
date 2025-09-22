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

func main() {
	fmt.Println("=== Testing JSON Format ===")

	// Load config
	privateKeyHex := os.Getenv("HYPERLIQUID_PRIVATE_KEY")
	if privateKeyHex == "" {
		if file, err := os.Open("./config.json"); err == nil {
			defer file.Close()
			var config map[string]string
			if json.NewDecoder(file).Decode(&config) == nil {
				privateKeyHex = config["secret_key"]
			}
		}
		if privateKeyHex == "" {
			log.Fatal("Private key not found")
		}
	}

	privateKey, err := utils.ParsePrivateKey(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Test exact payload construction
	timestamp := utils.GetTimestampMS()
	
	action := map[string]interface{}{
		"destination": "0x1234567890123456789012345678901234567890",
		"amount":      "0.01",
		"time":        timestamp,
		"type":        "usdSend",
	}

	signature, err := utils.SignUSDTransferAction(privateKey, action, false)
	if err != nil {
		log.Fatalf("Signing failed: %v", err)
	}

	// Create payload with explicit nil handling
	payload := map[string]interface{}{
		"action":       action,
		"nonce":        timestamp,
		"signature":    signature,
		"vaultAddress": nil,
		"expiresAfter": nil,
	}

	// Check what JSON is actually generated
	jsonBytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal failed: %v", err)
	}

	fmt.Println("Actual JSON being sent:")
	fmt.Println(string(jsonBytes))

	// Compare with working curl format
	fmt.Println("\nWorking curl format (for reference):")
	workingExample := `{
  "action": {
    "grouping": "na",
    "orders": [...],
    "type": "order"
  },
  "expiresAfter": 1758476080506,
  "isFrontend": true,
  "nonce": 1758476066273,
  "signature": {...},
  "vaultAddress": null
}`
	fmt.Println(workingExample)

	// Test sending to API
	timeout := 30 * time.Second
	api := client.NewAPI(utils.TestnetAPIURL, &timeout)
	
	result, err := api.Post("/exchange", payload)
	if err != nil {
		fmt.Printf("\nAPI call failed: %v\n", err)
	} else {
		fmt.Printf("\nAPI call succeeded: %+v\n", result)
	}
}