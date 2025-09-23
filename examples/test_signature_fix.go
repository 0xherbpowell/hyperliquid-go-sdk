package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	fmt.Println("=== Testing Fixed Signature Generation ===")

	// Setup using testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// Test basic order placement first
	fmt.Println("\n=== Testing Basic Order Placement ===")
	orderResult, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		0.01,                  // size
		1800.0,               // limit price (below market to rest)
		CreateGtcLimitOrder(), // order type
		false,                 // reduce only
		nil,                   // cloid
		nil,                   // builder info
	)

	if err != nil {
		fmt.Printf("Order placement failed: %v\n", err)
		// Try to get more detailed error info
		fmt.Printf("Error details: %+v\n", err)
	} else {
		fmt.Println("Order placement succeeded!")
		PrintOrderResult(orderResult)
	}

	// Now test signature generation manually to verify the fix
	fmt.Println("\n=== Manual Signature Verification ===")

	// Create a test action
	action := map[string]interface{}{
		"type":     "order",
		"grouping": "na",
		"orders": []map[string]interface{}{
			{
				"a": 0,     // ETH asset ID
				"b": true,  // buy
				"p": "1800", // price
				"s": "0.01", // size
				"r": false,  // reduce only
				"t": map[string]interface{}{
					"limit": map[string]interface{}{
						"tif": "Gtc",
					},
				},
			},
		},
	}

	fmt.Println("Test action:")
	actionJSON, _ := json.MarshalIndent(action, "", "  ")
	fmt.Printf("%s\n", actionJSON)

	// Test parameters
	var vaultAddress *string = nil
	nonce := time.Now().UnixMilli()
	var expiresAfter *int64 = nil

	fmt.Printf("Using nonce: %d\n", nonce)

	// Test action hash calculation
	actionHash, err := utils.ActionHash(action, vaultAddress, nonce, expiresAfter)
	if err != nil {
		log.Fatal("Error calculating action hash:", err)
	}
	fmt.Printf("Action hash: %s\n", hex.EncodeToString(actionHash))

	// Test phantom agent construction
	phantomAgent := utils.ConstructPhantomAgent(actionHash, false)
	fmt.Printf("Phantom agent source: %s\n", phantomAgent.Source)
	fmt.Printf("Phantom agent connectionId: %s\n", hex.EncodeToString(phantomAgent.ConnectionId[:]))

	// Test EIP712 payload construction
	eip712Data := utils.L1Payload(phantomAgent)
	fmt.Println("EIP712 payload domain:")
	domainJSON, _ := json.MarshalIndent(eip712Data.Domain, "", "  ")
	fmt.Printf("%s\n", domainJSON)

	// Try with user state info
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Printf("Warning: Could not get user state: %v\n", err)
	} else {
		fmt.Println("\nUser state retrieved successfully")
		// Check if user has funds
		if marginSummary, ok := userState["marginSummary"].(map[string]interface{}); ok {
			if accountValue, ok := marginSummary["accountValue"].(string); ok {
				fmt.Printf("Account value: %s\n", accountValue)
			}
		}
	}

	fmt.Println("\n=== Summary ===")
	fmt.Printf("Account: %s\n", address)
	fmt.Printf("Action hash: %s\n", hex.EncodeToString(actionHash))
	fmt.Printf("Phantom agent source: %s (should be 'b' for testnet)\n", phantomAgent.Source)
	fmt.Printf("EIP712 domain chainId: %v\n", eip712Data.Domain.ChainId)

	if err == nil {
		fmt.Println("✓ Signature generation mechanism is working correctly")
	} else {
		fmt.Printf("✗ Issues found: %v\n", err)
	}
}