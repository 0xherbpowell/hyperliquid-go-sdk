package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// Get user state before transfer
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("User state before transfer:")
	PrintPositions(userState)

	// Transfer 1.23 USDC from perp wallet to spot wallet
	fmt.Println("\nTransferring 1.23 USDC from perp wallet to spot wallet...")
	transferResult, err := exchange.UsdClassTransfer("1.23", false)
	if err != nil {
		log.Printf("Failed to transfer from perp to spot: %v", err)
	} else {
		fmt.Println("Transfer from perp to spot result:")
		PrintOrderResult(transferResult)
	}

	// Transfer 1.23 USDC from spot wallet to perp wallet
	fmt.Println("\nTransferring 1.23 USDC from spot wallet to perp wallet...")
	transferResult2, err := exchange.UsdClassTransfer("1.23", true)
	if err != nil {
		log.Printf("Failed to transfer from spot to perp: %v", err)
	} else {
		fmt.Println("Transfer from spot to perp result:")
		PrintOrderResult(transferResult2)
	}

	fmt.Println("\nNote: UsdClassTransfer moves USDC between your perp and spot wallets within the same account.")
	fmt.Println("- false: moves from perp to spot")
	fmt.Println("- true: moves from spot to perp")
}