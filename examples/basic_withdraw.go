package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// Get the user state before withdrawal
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("User state before withdrawal:")
	if marginSummary, ok := userState["marginSummary"].(map[string]interface{}); ok {
		fmt.Printf("Account Value: %v\n", marginSummary["accountValue"])
		fmt.Printf("Total RAW USD: %v\n", marginSummary["totalRawUsd"])
	}

	// Destination address for withdrawal
	// WARNING: Make sure this is a valid Ethereum address you control!
	// This should be your Ethereum wallet address where you want to receive the funds
	destinationAddress := "0x0000000000000000000000000000000000000000" // Replace with your actual Ethereum address

	fmt.Printf("\nWithdrawing USDC from Hyperliquid to Ethereum address: %s\n", destinationAddress)

	// Withdraw a small amount (in USDC)
	withdrawAmount := "1.0" // $1.00 USDC

	withdrawResult, err := exchange.WithdrawFromBridge(destinationAddress, withdrawAmount)
	if err != nil {
		log.Printf("Failed to withdraw: %v", err)

		// If withdrawal fails, it might be because:
		// 1. Insufficient balance
		// 2. Invalid destination address (must be valid Ethereum address)
		// 3. Minimum withdrawal amount not met
		// 4. Network issues
		// 5. Bridge limitations on testnet
		// 6. Withdrawal temporarily disabled

		fmt.Println("\nWithdrawal failed. This might be expected because:")
		fmt.Println("1. Insufficient balance to cover withdrawal + gas fees")
		fmt.Println("2. Invalid Ethereum destination address")
		fmt.Println("3. Withdrawal amount below minimum threshold")
		fmt.Println("4. Bridge functionality may be limited on testnet")
		fmt.Println("5. Withdrawal requires sufficient account value")
		fmt.Println("Error details:")
		fmt.Printf("  %v\n", err)

		return
	}

	fmt.Println("Withdrawal result:")
	PrintOrderResult(withdrawResult)

	// Get updated user state to see the change
	updatedUserState, err := info.UserState(address, "")
	if err != nil {
		log.Printf("Failed to get updated user state: %v", err)
		return
	}

	fmt.Println("\nUser state after withdrawal:")
	if marginSummary, ok := updatedUserState["marginSummary"].(map[string]interface{}); ok {
		fmt.Printf("Account Value: %v\n", marginSummary["accountValue"])
		fmt.Printf("Total RAW USD: %v\n", marginSummary["totalRawUsd"])
	}

	fmt.Println("\nWithdrawal initiated!")
	fmt.Printf("Note: Initiated withdrawal of %s USDC to Ethereum address %s\n", withdrawAmount, destinationAddress)
	fmt.Println("Bridge withdrawals may take some time to process. Check your Ethereum wallet for the funds.")
	fmt.Println("WARNING: Always double-check the destination address before withdrawing!")
}