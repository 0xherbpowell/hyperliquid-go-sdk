package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// Get the user state and print out position information
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("Initial user state:")
	PrintPositions(userState)

	// Display current leverage settings
	fmt.Println("Current leverage settings:")
	if assetPositions, ok := userState["assetPositions"].([]interface{}); ok {
		for _, ap := range assetPositions {
			if apMap, ok := ap.(map[string]interface{}); ok {
				if position, ok := apMap["position"].(map[string]interface{}); ok {
					coin := position["coin"].(string)
					leverage := position["leverage"]
					fmt.Printf("  %s: %+v\n", coin, leverage)
				}
			}
		}
	}

	// Adjust leverage for ETH to cross leverage 10x
	fmt.Println("\nSetting ETH leverage to 10x cross...")
	leverageResult, err := exchange.UpdateLeverage("ETH", true, 10) // true = cross leverage
	if err != nil {
		log.Printf("Failed to update leverage: %v", err)
	} else {
		fmt.Println("Leverage adjustment result:")
		PrintOrderResult(leverageResult)
	}

	// Adjust leverage for BTC to isolated leverage 5x
	fmt.Println("\nSetting BTC leverage to 5x isolated...")
	leverageResult2, err := exchange.UpdateLeverage("BTC", false, 5) // false = isolated leverage
	if err != nil {
		log.Printf("Failed to update BTC leverage: %v", err)
	} else {
		fmt.Println("BTC leverage adjustment result:")
		PrintOrderResult(leverageResult2)
	}

	// Get updated user state to verify changes
	updatedUserState, err := info.UserState(address, "")
	if err != nil {
		log.Printf("Failed to get updated user state: %v", err)
		return
	}

	fmt.Println("\nUpdated leverage settings:")
	if assetPositions, ok := updatedUserState["assetPositions"].([]interface{}); ok {
		for _, ap := range assetPositions {
			if apMap, ok := ap.(map[string]interface{}); ok {
				if position, ok := apMap["position"].(map[string]interface{}); ok {
					coin := position["coin"].(string)
					leverage := position["leverage"]
					fmt.Printf("  %s: %+v\n", coin, leverage)
				}
			}
		}
	}

	fmt.Println("\nLeverage adjustment completed!")
	fmt.Println("Note: Leverage changes apply to future trades. Existing positions maintain their original leverage.")
}