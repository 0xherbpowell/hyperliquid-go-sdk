package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

// This example shows how to switch an account to use big blocks on the EVM
func main() {
	// Setup with testnet
	address, _, exchange := Setup(utils.TestnetAPIURL, true)

	fmt.Printf("Managing big blocks setting for address: %s\n", address)

	// Enable big blocks
	fmt.Println("\nEnabling big blocks...")
	result1, err := exchange.UseBigBlocks(true)
	if err != nil {
		log.Printf("Failed to enable big blocks: %v", err)
	} else {
		fmt.Println("Enable big blocks result:")
		PrintOrderResult(result1)
	}

	// Disable big blocks
	fmt.Println("\nDisabling big blocks...")
	result2, err := exchange.UseBigBlocks(false)
	if err != nil {
		log.Printf("Failed to disable big blocks: %v", err)
	} else {
		fmt.Println("Disable big blocks result:")
		PrintOrderResult(result2)
	}

	fmt.Println("\nNote: Big blocks setting affects how your transactions are processed on the EVM.")
	fmt.Println("- true: Use big blocks (may provide different gas/processing characteristics)")
	fmt.Println("- false: Use normal blocks")
}