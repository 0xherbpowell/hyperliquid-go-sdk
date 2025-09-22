package main

import (
	"fmt"
	"log"
	"time"

	"hyperliquid-go-sdk/pkg/client"
	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// This example demonstrates the concept of agents/sub-accounts
	// Note: Full agent functionality may require additional API endpoints

	fmt.Println("Basic Agent Example (Conceptual)")
	fmt.Printf("Main account: %s\n", address)

	// Get initial user state
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("\nInitial account state:")
	PrintPositions(userState)

	// Generate a new keypair for demonstration purposes
	agentWallet, err := utils.CreateRandomWallet()
	if err != nil {
		log.Fatalf("Failed to create agent wallet: %v", err)
	}

	agentAddress := utils.GetAddressFromPrivateKey(agentWallet)
	fmt.Printf("\nCreated agent address: %s\n", agentAddress)

	// In a real agent system, you would:
	fmt.Println("\nAgent functionality (conceptual):")
	fmt.Println("1. Set up agent permissions through the main account")
	fmt.Println("2. Configure trading limits and restrictions")
	fmt.Println("3. Allow the agent to trade within those limits")
	fmt.Println("4. Monitor agent activity through the main account")

	// Create a separate client for the agent (for demonstration)
	timeout := 30 * time.Second
	agentInfo, err := client.NewInfo(utils.TestnetAPIURL, &timeout, true, nil, nil, nil)
	if err != nil {
		log.Printf("Failed to create agent info client: %v", err)
		return
	}

	_, err = client.NewExchange(
		agentWallet,
		utils.TestnetAPIURL,
		&timeout,
		nil,
		nil,
		&agentAddress,
		nil,
		nil,
	)
	if err != nil {
		log.Printf("Failed to create agent exchange client: %v", err)
		return
	}

	// Note: In production, the agent exchange would be used for actual trading
	fmt.Printf("Agent exchange client created successfully\n")

	// Demonstrate agent functionality with available SDK features
	fmt.Printf("\nAgent %s created (separate wallet for trading)\n", agentAddress)

	// Note: In a production agent system, the agent would need funds or
	// delegation authority to trade. For this example, we'll just show the concept.

	// Get current ETH price for reference
	mids, err := agentInfo.AllMids("")
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

	// In a real scenario, you might:
	fmt.Println("\nAgent capabilities (with proper setup):")
	fmt.Println("1. Place orders on behalf of the main account")
	fmt.Println("2. Manage risk within predefined limits")
	fmt.Println("3. Execute algorithmic trading strategies")
	fmt.Println("4. Report activity back to the main account")

	// Demonstrate placing an order from the main account (simulating agent activity)
	testOrderPrice := ethPrice * 0.98 // 2% below market
	testOrderSize := 0.01             // Small size for testing

	fmt.Printf("\nDemonstrating order placement (from main account as agent example):\n")
	fmt.Printf("Placing test buy order: %f ETH at %f\n", testOrderSize, testOrderPrice)

	orderResult, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy (buy order)
		testOrderSize,         // size
		testOrderPrice,        // limit price
		CreateGtcLimitOrder(), // GTC order type
		false,                 // reduce only
		GenerateCloid(),       // unique client order ID
		nil,                   // builder info
	)

	if err != nil {
		log.Printf("Failed to place order: %v", err)
	} else {
		fmt.Println("Order result:")
		PrintOrderResult(orderResult)

		// Cancel the test order (cleanup)
		if oid, ok := GetRestingOid(orderResult); ok && oid > 0 {
			fmt.Printf("\nCancelling test order (oid: %d)...\n", oid)

			cancelResult, err := exchange.Cancel("ETH", oid)
			if err != nil {
				log.Printf("Failed to cancel order: %v", err)
			} else {
				fmt.Printf("Order cancelled successfully\n")
				_ = cancelResult // Suppress unused variable warning
			}
		}
	}

	fmt.Println("\nAgent example completed!")
	fmt.Println("Note: This example demonstrated:")
	fmt.Println("1. Creating separate wallet keypairs for agents")
	fmt.Println("2. Setting up separate client instances")
	fmt.Println("3. The conceptual framework for agent trading")
	fmt.Println("4. How to structure agent-based trading systems")
	fmt.Println("\nFor full agent functionality, you would need:")
	fmt.Println("• Proper authorization/delegation mechanisms")
	fmt.Println("• Fund allocation to agent accounts")
	fmt.Println("• Risk management and monitoring systems")
	fmt.Println("• Audit trails and reporting")
}
