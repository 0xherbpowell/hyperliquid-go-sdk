package main

import (
	"fmt"
	"log"
	"strconv"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// This example demonstrates how to send assets (tokens/coins) to another address
	// This includes both USD transfers and spot token transfers

	fmt.Println("Basic Send Asset Example")
	fmt.Printf("From account: %s\n", address)

	// Get initial user state to see available balances
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("\nInitial account balances:")
	PrintBalances(userState)

	// Example destination address (in a real scenario, this would be a different address)
	// Using a different address for demonstration - replace with actual destination
	destinationAddress := "0x1234567890123456789012345678901234567890"
	fmt.Printf("Destination address: %s\n", destinationAddress)

	// Example 1: Send USD (USDC)
	fmt.Println("\n--- Example 1: Send USD ---")

	// Check if we have USD balance
	var usdBalance float64 = 0
	if withdrawable, ok := userState["withdrawable"].(map[string]interface{}); ok {
		if usdStr, exists := withdrawable["USD"]; exists {
			if usdVal, err := utils.ParsePrice(usdStr.(string)); err == nil {
				usdBalance = usdVal
			}
		}
	}

	fmt.Printf("Current USD balance: %f\n", usdBalance)

	if usdBalance > 1.0 { // Only proceed if we have more than 1 USD
		sendAmount := "0.5" // Send 0.5 USD

		fmt.Printf("Sending %s USD to %s\n", sendAmount, destinationAddress)

		// Send USD using UsdTransfer
		sendResult, err := exchange.UsdTransfer(destinationAddress, sendAmount)
		if err != nil {
			log.Printf("Failed to send USD: %v", err)
		} else {
			fmt.Println("USD send result:")
			PrintOrderResult(sendResult)
		}
	} else {
		fmt.Println("Insufficient USD balance for sending (need at least 1.0 USD)")
	}

	// Example 2: Send spot token (e.g., send ETH if available)
	fmt.Println("\n--- Example 2: Send Spot Token (ETH) ---")

	// Check if we have ETH balance in spot
	var ethBalance float64 = 0
	if spotBalances, ok := userState["spotBalances"].([]interface{}); ok {
		for _, balance := range spotBalances {
			if balanceMap, ok := balance.(map[string]interface{}); ok {
				if coin, exists := balanceMap["coin"]; exists && coin == "ETH" {
					if balanceStr, exists := balanceMap["hold"]; exists {
						if balanceVal, err := utils.ParsePrice(balanceStr.(string)); err == nil {
							ethBalance = balanceVal
						}
					}
				}
			}
		}
	}

	fmt.Printf("Current ETH spot balance: %f\n", ethBalance)

	if ethBalance > 0.01 { // Only proceed if we have more than 0.01 ETH
		sendAmount := "0.001" // Send 0.001 ETH

		fmt.Printf("Sending %s ETH to %s\n", sendAmount, destinationAddress)

		// Send spot token (ETH)
		sendResult, err := exchange.SpotTransfer(destinationAddress, "ETH", sendAmount)
		if err != nil {
			log.Printf("Failed to send ETH: %v", err)
		} else {
			fmt.Println("ETH send result:")
			PrintOrderResult(sendResult)
		}
	} else {
		fmt.Println("Insufficient ETH balance for sending (need at least 0.01 ETH)")
	}

	// Example 3: Send other available spot tokens
	fmt.Println("\n--- Example 3: Send Other Available Tokens ---")

	if spotBalances, ok := userState["spotBalances"].([]interface{}); ok {
		fmt.Println("Available spot tokens for transfer:")
		
		for _, balance := range spotBalances {
			if balanceMap, ok := balance.(map[string]interface{}); ok {
				coin := balanceMap["coin"]
				hold := balanceMap["hold"]
				
				if holdVal, err := strconv.ParseFloat(hold.(string), 64); err == nil && holdVal > 0 {
					fmt.Printf("  %s: %s available\n", coin, hold)
					
					// Example: send a small amount of the first available token (other than ETH)
					if coin != "ETH" && holdVal > 0.1 {
						sendAmount := "0.01"
						fmt.Printf("    Sending %s %s to %s\n", sendAmount, coin, destinationAddress)
						
						sendResult, err := exchange.SpotTransfer(destinationAddress, coin.(string), sendAmount)
						if err != nil {
							log.Printf("    Failed to send %s: %v", coin, err)
						} else {
							fmt.Printf("    %s send successful\n", coin)
							_ = sendResult
							break // Only send one token as example
						}
					}
				}
			}
		}
	}

	// Wait a moment for transactions to process
	fmt.Println("\nWaiting for transactions to process...")
	// time.Sleep(3 * time.Second) // Uncomment if needed

	// Check updated balances after sends
	fmt.Println("\n--- Final Balance Check ---")

	updatedUserState, err := info.UserState(address, "")
	if err != nil {
		log.Printf("Failed to get updated user state: %v", err)
	} else {
		fmt.Println("Updated account balances:")
		PrintBalances(updatedUserState)
	}

	// Example of querying transaction history to see our sends
	fmt.Println("\n--- Recent Transaction History ---")

	// Note: This would show recent transfers/sends in the transaction history
	recentTransfers, err := info.UserFills(address, "")
	if err != nil {
		log.Printf("Failed to get recent transfers: %v", err)
	} else {
		// Process and display recent transfers
		if fills, ok := recentTransfers["fills"]; ok {
			// If transfers are included in fills data
			if fillsArray, ok := fills.([]interface{}); ok {
				fmt.Printf("Recent account activity: %d items\n", len(fillsArray))
			}
		}
	}

	fmt.Println("\nSend asset example completed!")
	fmt.Println("Note: This example demonstrated:")
	fmt.Println("1. Sending USD to another address")
	fmt.Println("2. Sending spot tokens (like ETH) to another address")
	fmt.Println("3. Checking available tokens for transfer")
	fmt.Println("4. Checking balances before and after transfers")
	fmt.Println("\nIMPORTANT: Always verify destination addresses before sending!")
	fmt.Println("Lost funds due to incorrect addresses cannot be recovered.")
}

// Helper function to print balances in a readable format
func PrintBalances(userState map[string]interface{}) {
	// Print withdrawable balances (native tokens like USD)
	if withdrawable, ok := userState["withdrawable"].(map[string]interface{}); ok {
		fmt.Println("  Withdrawable balances:")
		for token, balance := range withdrawable {
			fmt.Printf("    %s: %s\n", token, balance)
		}
	}

	// Print spot balances
	if spotBalances, ok := userState["spotBalances"].([]interface{}); ok {
		fmt.Println("  Spot balances:")
		for _, balance := range spotBalances {
			if balanceMap, ok := balance.(map[string]interface{}); ok {
				coin := balanceMap["coin"]
				hold := balanceMap["hold"]
				total := balanceMap["total"]
				fmt.Printf("    %s: hold=%s, total=%s\n", coin, hold, total)
			}
		}
	}

	// Print margin summary if available
	if marginSummary, ok := userState["marginSummary"].(map[string]interface{}); ok {
		if accountValue, exists := marginSummary["accountValue"]; exists {
			fmt.Printf("  Account value: %s\n", accountValue)
		}
	}
}