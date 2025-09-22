package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet using the example utils
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	fmt.Printf("Trading with address: %s\n", address)
	
	// Example 1: Get user state
	fmt.Println("\n=== Getting User State ===")
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Printf("Failed to get user state: %v", err)
	} else {
		fmt.Println("User positions:")
		PrintPositions(userState)
	}
	
	// Example 2: Get all mid prices
	fmt.Println("\n=== Getting All Mid Prices ===")
	mids, err := info.AllMids("")
	if err != nil {
		log.Printf("Failed to get mids: %v", err)
	} else {
		fmt.Printf("Available markets:\n")
		for coin, price := range mids {
			fmt.Printf("  %s: %s\n", coin, price)
		}
	}
	
	// Example 3: Place a limit order
	fmt.Println("\n=== Placing Limit Order ===")
	
	// Debug: Show all available mids
	if len(mids) == 0 {
		log.Printf("No markets found. This may indicate an API error or that testnet has no active markets.")
		return
	}
	
	// Get current ETH price for reference
	ethMid, exists := mids["ETH"]
	if !exists {
		log.Printf("ETH not found in markets. Available markets: %d", len(mids))
		// Try to find a market that exists
		for coin, price := range mids {
			log.Printf("Available: %s = %s", coin, price)
			break
		}
		return
	}
	
	ethPrice, err := utils.ParsePrice(ethMid)
	if err != nil {
		log.Printf("Failed to parse ETH price: %v", err)
		return
	}
	
	fmt.Printf("Current ETH price: %f\n", ethPrice)
	
	// Create a limit order for ETH (5% below market)
	orderPrice := ethPrice * 0.95
	orderType := CreateGtcLimitOrder()
	
	// Generate a unique client order ID
	cloid := GenerateCloid()
	
	result, err := exchange.Order(
		"ETH",      // coin
		true,       // isBuy
		0.01,       // size (0.01 ETH)
		orderPrice, // limit price (5% below market)
		orderType,  // order type
		false,      // reduce only
		cloid,      // client order ID
		nil,        // builder info
	)
	
	if err != nil {
		log.Printf("Failed to place order: %v", err)
	} else {
		fmt.Println("Order result:")
		PrintOrderResult(result)
		
		// Cancel the order for cleanup
		if oid, ok := GetRestingOid(result); ok && oid > 0 {
			fmt.Printf("\nCancelling order (oid: %d)...\n", oid)
			cancelResult, err := exchange.Cancel("ETH", oid)
			if err != nil {
				log.Printf("Failed to cancel order: %v", err)
			} else {
				fmt.Printf("Order cancelled successfully\n")
				_ = cancelResult
			}
		}
	}
	
	// Example 4: Place a market order
	fmt.Println("\n=== Placing Market Order ===")
	
	slippage := 0.01 // 1% slippage
	cloid2 := GenerateCloid()
	
	result, err = exchange.MarketOrder(
		"ETH",      // coin
		false,      // isBuy (sell)
		0.005,      // size (0.005 ETH)
		&slippage,  // slippage
		cloid2,     // client order ID
	)
	
	if err != nil {
		log.Printf("Failed to place market order: %v", err)
	} else {
		fmt.Println("Market order result:")
		PrintOrderResult(result)
	}
	
	// Example 5: Get open orders
	fmt.Println("\n=== Getting Open Orders ===")
	openOrders, err := info.OpenOrders(address, "")
	if err != nil {
		log.Printf("Failed to get open orders: %v", err)
	} else {
		fmt.Printf("Open orders: %+v\n", openOrders)
	}
	
	// Example 6: Get recent fills
	fmt.Println("\n=== Getting Recent Fills ===")
	fills, err := info.UserFills(address, "")
	if err != nil {
		log.Printf("Failed to get fills: %v", err)
	} else {
		fmt.Printf("Recent fills: %+v\n", fills)
	}
	
	// Example 7: Cancel all orders (commented out for safety)
	/*
	fmt.Println("\n=== Canceling All Orders ===")
	result, err = exchange.CancelAll()
	if err != nil {
		log.Printf("Failed to cancel all orders: %v", err)
	} else {
		fmt.Printf("Cancel all result: %+v\n", result)
	}
	*/
	
	fmt.Println("\nExample completed!")
}