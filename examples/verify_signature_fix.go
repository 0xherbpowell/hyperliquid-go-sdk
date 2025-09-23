package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	fmt.Println("=== Verifying Signature Fix ===")

	// Setup using testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	fmt.Printf("Account: %s\n", address)

	// Get user state to check account status
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Printf("Warning: Could not get user state: %v", err)
	} else {
		PrintPositions(userState)
		
		// Check account value
		if marginSummary, ok := userState["marginSummary"].(map[string]interface{}); ok {
			if accountValue, ok := marginSummary["accountValue"].(string); ok {
				fmt.Printf("Account value: %s\n", accountValue)
				if accountValue == "0" {
					fmt.Println("Account has no equity - using a small test order anyway")
				}
			}
		}
	}

	// Try to place a simple order to test if signature works
	fmt.Println("\n=== Testing Order Placement ===")
	
	orderResult, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		0.001,                 // very small size to minimize risk
		1500.0,               // price well below market to avoid execution
		CreateGtcLimitOrder(), // order type
		false,                 // reduce only
		nil,                   // cloid
		nil,                   // builder info
	)

	if err != nil {
		fmt.Printf("❌ Order placement failed: %v\n", err)
		fmt.Printf("Error type: %T\n", err)
		
		// Check if this is still a signature-related error
		errorStr := fmt.Sprintf("%v", err)
		if contains(errorStr, "does not exist") {
			fmt.Println("🔍 This appears to be a signature verification error - wallet address mismatch")
		} else if contains(errorStr, "422") {
			fmt.Println("🔍 This appears to be a JSON deserialization error")
		} else {
			fmt.Println("🔍 This appears to be a different type of error")
		}
		
		return
	}

	fmt.Println("✅ Order placement succeeded!")
	fmt.Println("🎉 Signature fix appears to be working!")
	
	PrintOrderResult(orderResult)

	// If we got an order ID, try to cancel it
	if oid, ok := GetRestingOid(orderResult); ok {
		fmt.Printf("\nCancelling order %d...\n", oid)
		cancelResult, err := exchange.Cancel("ETH", oid)
		if err != nil {
			fmt.Printf("Cancel failed: %v\n", err)
		} else {
			fmt.Println("Order cancelled successfully")
			PrintOrderResult(cancelResult)
		}
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr || 
		   (len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}