package main

import (
	"fmt"
	"log"
	"time"

	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, _, exchange := Setup(utils.TestnetAPIURL, true)

	fmt.Printf("Setting up multi-sig order for address: %s\n", address)

	// The outer signer is required to be an authorized user or an agent of the authorized user of the multi-sig user.

	// Address of the multi-sig user that the action will be executed for
	// Executing the action requires at least the specified threshold of signatures
	// required for that multi-sig user
	// WARNING: Replace with actual multi-sig user address!
	multiSigUser := "0x0000000000000000000000000000000000000005"

	timestamp := time.Now().UnixMilli()

	// Define the multi-sig inner action
	// This creates a limit order for asset 4 (BTC), buying 0.2 BTC at $1100
	action := map[string]interface{}{
		"type": "order",
		"orders": []map[string]interface{}{
			{
				"a": 4,
				"b": true,
				"p": "1100",
				"s": "0.2",
				"r": false,
				"t": map[string]interface{}{
					"limit": map[string]interface{}{
						"tif": "Gtc",
					},
				},
			},
		},
		"grouping": "na",
	}

	fmt.Printf("Multi-sig user: %s\n", multiSigUser)
	fmt.Printf("Action: %+v\n", action)
	fmt.Printf("Timestamp: %d\n", timestamp)

	// In a real implementation, you would need multiple wallets (private keys)
	// that belong to authorized users of the multi-sig user
	// For this example, we'll show the structure but it will likely fail
	// without proper multi-sig wallet setup

	var signatures []map[string]interface{}

	// NOTE: This is where you would collect signatures from each authorized wallet
	// Each wallet must belong to a user that is authorized for the multi-sig user
	// For demonstration purposes, we'll create a placeholder signature structure

	fmt.Println("\nWARNING: This example requires proper multi-sig wallet setup!")
	fmt.Println("You need to:")
	fmt.Println("1. Have a valid multi-sig user address")
	fmt.Println("2. Have access to the private keys of authorized users")
	fmt.Println("3. Sign the action with each required wallet")
	fmt.Println("4. Collect enough signatures to meet the threshold")

	// Create placeholder signatures (this will fail in actual execution)
	// In real implementation, you would:
	// for _, wallet := range multiSigWallets {
	//     signature := SignMultiSigL1ActionPayload(
	//         wallet,
	//         action,
	//         isMainnet,
	//         nil,
	//         timestamp,
	//         expiresAfter,
	//         multiSigUser,
	//         address,
	//     )
	//     signatures = append(signatures, signature)
	// }

	if len(signatures) == 0 {
		fmt.Println("\nSkipping multi-sig execution due to missing wallet setup.")
		fmt.Println("To use this example:")
		fmt.Println("1. Set up multiple authorized wallets")
		fmt.Println("2. Implement signature collection")
		fmt.Println("3. Update multiSigUser with a real multi-sig address")
		return
	}

	// Execute the multi-sig action with all collected signatures
	// This will only succeed if enough valid signatures are provided
	fmt.Println("\nExecuting multi-sig action...")
	multiSigResult, err := exchange.MultiSig(multiSigUser, action, signatures, timestamp)
	if err != nil {
		log.Printf("Failed to execute multi-sig action: %v", err)
		fmt.Println("\nCommon reasons for failure:")
		fmt.Println("1. Insufficient signatures (below threshold)")
		fmt.Println("2. Invalid signatures")
		fmt.Println("3. Unauthorized signers")
		fmt.Println("4. Invalid multi-sig user address")
		fmt.Println("5. Invalid action format")
		return
	}

	fmt.Println("Multi-sig order result:")
	PrintOrderResult(multiSigResult)

	fmt.Println("\nNote: Multi-sig orders require coordination between multiple authorized users.")
	fmt.Println("Each authorized user must sign the action with their private key.")
}