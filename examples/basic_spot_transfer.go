package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// Get spot user state (for spot balances)
	// Note: This is a placeholder - you'd need to implement spot user state
	// For now, we'll use the regular user state
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("User state before spot transfer:")
	PrintPositions(userState)

	// Destination address for transfer
	// WARNING: Make sure this is a valid address you control!
	destinationAddress := "0x0000000000000000000000000000000000000000" // Replace with actual destination

	fmt.Printf("\nTransferring spot assets from %s to %s\n", address, destinationAddress)

	// Transfer spot USDC (assuming it's available)
	// Common spot tokens: USDC, USDT, WETH, etc.
	token := "USDC"
	transferAmount := "1.0" // 1 USDC

	transferResult, err := exchange.SpotTransfer(destinationAddress, token, transferAmount)
	if err != nil {
		log.Printf("Failed to transfer spot %s: %v", token, err)

		// If transfer fails, it might be because:
		// 1. Insufficient balance
		// 2. Invalid destination address
		// 3. Invalid token
		// 4. Network issues
		// 5. API limitations on testnet

		fmt.Println("\nSpot transfer failed. This might be expected if:")
		fmt.Println("1. You don't have sufficient spot balance")
		fmt.Println("2. The token doesn't exist or isn't available")
		fmt.Println("3. Spot transfers are limited on testnet")
		fmt.Println("Error details:")
		fmt.Printf("  %v\n", err)

		// Try with a different common token
		token2 := "USDT"
		fmt.Printf("\nTrying with %s instead...\n", token2)
		
		transferResult2, err2 := exchange.SpotTransfer(destinationAddress, token2, transferAmount)
		if err2 != nil {
			fmt.Printf("Transfer with %s also failed: %v\n", token2, err2)
			return
		} else {
			fmt.Printf("Transfer with %s succeeded:\n", token2)
			PrintOrderResult(transferResult2)
		}
		
		return
	}

	fmt.Printf("Spot %s transfer result:\n", token)
	PrintOrderResult(transferResult)

	fmt.Printf("\nSpot transfer completed!")
	fmt.Printf("Note: Transferred %s %s from %s to %s\n", transferAmount, token, address, destinationAddress)
	fmt.Println("Remember: Spot transfers move tokens between addresses, not between spot/perp accounts.")
}