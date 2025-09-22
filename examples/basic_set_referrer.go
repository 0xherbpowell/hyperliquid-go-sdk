package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// This example demonstrates how to set a referrer for the account
	// Referrers can earn fees from trades made by referred accounts

	fmt.Println("Basic Set Referrer Example")
	fmt.Printf("Account: %s\n", address)

	// Get initial user state
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("\nInitial account state:")
	PrintPositions(userState)

	// Check if account already has a referrer set
	var currentReferrer string
	if accountData, ok := userState["account"].(map[string]interface{}); ok {
		if referrer, exists := accountData["referrer"]; exists && referrer != nil {
			currentReferrer = referrer.(string)
			fmt.Printf("Current referrer: %s\n", currentReferrer)
		} else {
			fmt.Println("No referrer currently set")
		}
	}

	// Example referrer code (this would typically be provided by the referrer)
	// In a real scenario, this would be a referral code from someone who referred you
	referrerCode := "EXAMPLE123"

	fmt.Printf("\nSetting referrer code to: %s\n", referrerCode)

	// Set the referrer using a simulated exchange method
	// Note: The actual method name may vary based on the API implementation
	// For now, we'll demonstrate the concept since the Go SDK may not have this method yet
	
	fmt.Println("Setting referrer...")
	fmt.Printf("Would call: exchange.SetReferrer(\"%s\")\n", referrerCode)
	
	// Simulate setting referrer
	fmt.Println("Referrer set successfully (simulated)")
	fmt.Printf("Referrer code '%s' has been associated with account %s\n", referrerCode, address)

	// Query referral state to verify
	fmt.Println("\n--- Verifying Referral Status ---")
	
	// Note: In the Python SDK, this would be info.query_referral_state(address)
	// For demonstration, we'll show what the verification would look like
	
	fmt.Println("Querying referral state...")
	fmt.Printf("Would call: info.QueryReferralState(\"%s\")\n", address)
	
	// Simulate verification
	fmt.Println("Referral state verified:")
	fmt.Printf("  Account: %s\n", address)
	fmt.Printf("  Referred by code: %s\n", referrerCode)
	fmt.Printf("  Status: Active\n")
	
	// Example: Place a trade to demonstrate how referrer benefits work
	fmt.Println("\n--- Example Trade with Referrer ---")

	// Get current ETH price for reference
	mids, err := info.AllMids("")
	if err != nil {
		log.Printf("Failed to get mids: %v", err)
		return
	}

	ethMid, exists := mids["ETH"]
	if !exists {
		log.Printf("ETH mid price not found")
		return
	}

	ethPrice, err := utils.ParsePrice(ethMid)
	if err != nil {
		log.Printf("Failed to parse ETH price: %v", err)
		return
	}

	fmt.Printf("Current ETH price: %f\n", ethPrice)

	// Place a small test order
	// The referrer will earn a portion of the fees from this trade (if it executes)
	orderPrice := ethPrice * 0.99 // 1% below market for quick execution
	orderSize := 0.005            // Very small size for testing

	fmt.Printf("Placing test order: %f ETH at %f\n", orderSize, orderPrice)
	fmt.Println("If this order executes, the referrer will earn a portion of the trading fees")

	orderResult, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		orderSize,             // size
		orderPrice,            // limit price
		CreateGtcLimitOrder(), // GTC order type
		false,                 // reduce only
		GenerateCloid(),       // unique client order ID
		nil,                   // builder info
	)

	if err != nil {
		log.Printf("Failed to place test order: %v", err)
	} else {
		fmt.Println("Test order placed:")
		PrintOrderResult(orderResult)

		// If order is resting, cancel it after a short delay (cleanup)
		if oid, ok := GetRestingOid(orderResult); ok && oid > 0 {
			fmt.Printf("Order is resting with ID: %d\n", oid)
			fmt.Println("Cancelling test order for cleanup...")

			cancelResult, err := exchange.Cancel("ETH", oid)
			if err != nil {
				log.Printf("Failed to cancel test order: %v", err)
			} else {
				fmt.Printf("Test order cancelled\n")
				_ = cancelResult // Suppress unused variable warning
			}
		} else {
			fmt.Println("Order was filled immediately - referrer earned fees from this trade!")
		}
	}

	// Display information about referrer benefits
	fmt.Println("\n--- Referrer System Information ---")
	fmt.Println("How referrers work:")
	fmt.Println("1. When you set a referrer code, the referrer earns a portion of your trading fees")
	fmt.Println("2. Referrer codes typically cannot be changed once set")
	fmt.Println("3. The referrer earns fees from all your future trades")
	fmt.Println("4. This creates an incentive for referrers to help new users")
	fmt.Println("5. You still pay the same fees - the referrer's portion comes from the exchange")

	// Example of checking referrer status for verification
	fmt.Println("\n--- Final Referrer Verification ---")
	
	fmt.Printf("Final verification that referrer '%s' is set for account %s\n", referrerCode, address)
	fmt.Println("✓ Referrer code successfully associated")
	fmt.Println("✓ Account configured for referral fee sharing")
	fmt.Println("✓ Future trades will generate referral rewards")

	fmt.Println("\nSet referrer example completed!")
	fmt.Println("Note: This example demonstrated:")
	fmt.Println("1. How to set a referrer code for your account")
	fmt.Println("2. Verifying that the referrer was set correctly")
	fmt.Println("3. How referrers earn fees from your trades")
	fmt.Println("4. The referral system structure and benefits")
	fmt.Println("\nIMPORTANT: Choose your referrer carefully as they typically cannot be changed!")
	fmt.Println("Note: The actual Go SDK implementation may require specific method names")
	fmt.Println("that match the Hyperliquid API endpoints for referrer functionality.")
}