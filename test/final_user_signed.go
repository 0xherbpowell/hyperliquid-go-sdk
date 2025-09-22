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
	fmt.Println("=== Final User-Signed Action Test ===")

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

	timeout := 30 * time.Second
	api := client.NewAPI(utils.TestnetAPIURL, &timeout)

	// Test different formats for user-signed actions

	// Format 1: Direct action (no wrapper) - like Python might do for user-signed
	fmt.Println("\n=== Format 1: Direct USD Transfer ===")
	timestamp1 := utils.GetTimestampMS()
	
	usdAction := map[string]interface{}{
		"destination": "0x1234567890123456789012345678901234567890",
		"amount":      "0.01", 
		"time":        fmt.Sprintf("%d", timestamp1), // String for EIP712
	}

	usdSig, err := utils.SignUSDTransferAction(privateKey, usdAction, false)
	if err != nil {
		log.Printf("USD signing failed: %v", err)
		return
	}

	directUsdPayload := map[string]interface{}{
		"type":        "usdSend",
		"destination": "0x1234567890123456789012345678901234567890", 
		"amount":      "0.01",
		"time":        timestamp1, // int64 for API
		"signature":   usdSig,
	}

	fmt.Println("Direct USD payload JSON:")
	jsonBytes1, _ := json.MarshalIndent(directUsdPayload, "", "  ")
	fmt.Println(string(jsonBytes1))

	result1, err := api.Post("/exchange", directUsdPayload)
	if err != nil {
		fmt.Printf("Direct USD failed: %v\n", err)
	} else {
		fmt.Printf("Direct USD succeeded: %+v\n", result1)
	}

	// Format 2: Try with different action type name
	fmt.Println("\n=== Format 2: USD Transfer with different type ===")
	timestamp2 := utils.GetTimestampMS()
	
	usdAction2 := map[string]interface{}{
		"destination": "0x1234567890123456789012345678901234567890",
		"amount":      "0.01", 
		"time":        fmt.Sprintf("%d", timestamp2),
	}

	usdSig2, err := utils.SignUSDTransferAction(privateKey, usdAction2, false)
	if err != nil {
		log.Printf("USD signing 2 failed: %v", err)
		return
	}

	// Try with "usdTransfer" instead of "usdSend"
	directUsdPayload2 := map[string]interface{}{
		"type":        "usdTransfer", // Different type name
		"destination": "0x1234567890123456789012345678901234567890", 
		"amount":      "0.01",
		"time":        timestamp2,
		"signature":   usdSig2,
	}

	fmt.Println("USD Transfer (different type) JSON:")
	jsonBytes2, _ := json.MarshalIndent(directUsdPayload2, "", "  ")
	fmt.Println(string(jsonBytes2))

	result2, err := api.Post("/exchange", directUsdPayload2)
	if err != nil {
		fmt.Printf("USD Transfer (alt type) failed: %v\n", err)
	} else {
		fmt.Printf("USD Transfer (alt type) succeeded: %+v\n", result2)
	}

	// Format 3: Agent approval direct
	fmt.Println("\n=== Format 3: Direct Agent Approval ===")
	
	agentPrivateKey, err := utils.CreateRandomWallet()
	if err != nil {
		log.Printf("Failed to create agent: %v", err)
		return
	}
	agentAddress := utils.GetAddressFromPrivateKey(agentPrivateKey)
	timestamp3 := utils.GetTimestampMS()
	
	agentAction := map[string]interface{}{
		"agentAddress": agentAddress,
		"agentName":    "", 
		"nonce":       fmt.Sprintf("%d", timestamp3), // String for EIP712
	}

	agentSig, err := utils.SignAgent(privateKey, agentAction, false)
	if err != nil {
		log.Printf("Agent signing failed: %v", err)
		return
	}

	directAgentPayload := map[string]interface{}{
		"type":         "approveAgent",
		"agentAddress": agentAddress,
		"nonce":       timestamp3, // int64 for API
		"signature":   agentSig,
		// Don't include empty agentName
	}

	fmt.Println("Direct agent payload JSON:")
	jsonBytes3, _ := json.MarshalIndent(directAgentPayload, "", "  ")
	fmt.Println(string(jsonBytes3))

	result3, err := api.Post("/exchange", directAgentPayload)
	if err != nil {
		fmt.Printf("Direct agent failed: %v\n", err)
	} else {
		fmt.Printf("Direct agent succeeded: %+v\n", result3)
	}
}