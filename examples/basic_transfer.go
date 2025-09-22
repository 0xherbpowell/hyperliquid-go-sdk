package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// Get the user state before transfer
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("User state before transfer:")
	if marginSummary, ok := userState["marginSummary"].(map[string]interface{}); ok {
		fmt.Printf("Account Value: %v\n", marginSummary["accountValue"])
		fmt.Printf("Total Margin Used: %v\n", marginSummary["totalMarginUsed"])
		fmt.Printf("Total NTL POS: %v\n", marginSummary["totalNtlPos"])
		fmt.Printf("Total RAW USD: %v\n", marginSummary["totalRawUsd"])
	}

	// Destination address for transfer
	// WARNING: Make sure this is a valid address you control!
	// For testing, you might want to transfer to another wallet you own
	destinationAddress := "0x0000000000000000000000000000000000000000" // Replace with actual destination
	
	fmt.Printf("\nTransferring USD from %s to %s\n", address, destinationAddress)
	
	// Transfer a small amount (in USD)
	transferAmount := "1.0" // $1.00 USD
	
	transferResult, err := exchange.UsdTransfer(destinationAddress, transferAmount)
	if err != nil {
		log.Printf("Failed to transfer USD: %v", err)
		
		// If transfer fails, it might be because:
		// 1. Insufficient balance
		// 2. Invalid destination address
		// 3. Network issues
		// 4. API limitations on testnet
		
		fmt.Println("\nTransfer failed. This might be expected on testnet or with insufficient balance.")
		fmt.Println("Error details:")
		fmt.Printf("  %v\n", err)
		
		return
	}

	fmt.Println("Transfer result:")
	PrintOrderResult(transferResult)

	// Get updated user state to see the change
	updatedUserState, err := info.UserState(address, "")
	if err != nil {
		log.Printf("Failed to get updated user state: %v", err)
		return
	}

	fmt.Println("\nUser state after transfer:")
	if marginSummary, ok := updatedUserState["marginSummary"].(map[string]interface{}); ok {
		fmt.Printf("Account Value: %v\n", marginSummary["accountValue"])
		fmt.Printf("Total Margin Used: %v\n", marginSummary["totalMarginUsed"])
		fmt.Printf("Total NTL POS: %v\n", marginSummary["totalNtlPos"])
		fmt.Printf("Total RAW USD: %v\n", marginSummary["totalRawUsd"])
	}

	fmt.Println("\nTransfer completed!")
	fmt.Printf("Note: Transferred %s USD from %s to %s\n", transferAmount, address, destinationAddress)
}