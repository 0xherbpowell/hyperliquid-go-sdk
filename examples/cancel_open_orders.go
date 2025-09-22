package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// Get current open orders
	openOrders, err := info.OpenOrders(address, "")
	if err != nil {
		log.Fatalf("Failed to get open orders: %v", err)
	}

	fmt.Println("Current open orders:")
	PrintOrderResult(openOrders)

	// Count the number of open orders
	orderCount := 0
	if orders, ok := openOrders["orders"].([]interface{}); ok {
		orderCount = len(orders)
	}

	if orderCount == 0 {
		fmt.Println("No open orders to cancel.")
		
		// Place a few orders first so we have something to cancel
		fmt.Println("Placing a few orders to demonstrate cancellation...")
		
		// Place order 1
		_, err1 := exchange.Order(
			"ETH", true, 0.1, 1000.0, // Low price to ensure it rests
			CreateGtcLimitOrder(), false, GenerateCloid(), nil,
		)
		if err1 != nil {
			log.Printf("Failed to place order 1: %v", err1)
		}

		// Place order 2
		_, err2 := exchange.Order(
			"BTC", true, 0.01, 20000.0, // Low price to ensure it rests
			CreateGtcLimitOrder(), false, GenerateCloid(), nil,
		)
		if err2 != nil {
			log.Printf("Failed to place order 2: %v", err2)
		}

		// Get updated open orders
		openOrders, err = info.OpenOrders(address, "")
		if err != nil {
			log.Printf("Failed to get updated open orders: %v", err)
			return
		}

		fmt.Println("\nOpen orders after placing new orders:")
		PrintOrderResult(openOrders)
	}

	// Cancel all open orders using the bulk cancel method
	fmt.Println("\nCancelling all open orders...")
	
	cancelResult, err := exchange.CancelAll()
	if err != nil {
		log.Printf("Failed to cancel all orders: %v", err)
		return
	}

	fmt.Println("Cancel all orders result:")
	PrintOrderResult(cancelResult)

	// Verify that orders were cancelled
	finalOpenOrders, err := info.OpenOrders(address, "")
	if err != nil {
		log.Printf("Failed to get final open orders: %v", err)
		return
	}

	fmt.Println("\nOpen orders after cancellation:")
	PrintOrderResult(finalOpenOrders)

	// Count remaining orders
	remainingCount := 0
	if orders, ok := finalOpenOrders["orders"].([]interface{}); ok {
		remainingCount = len(orders)
	}

	if remainingCount == 0 {
		fmt.Println("✅ All orders successfully cancelled!")
	} else {
		fmt.Printf("⚠️  %d orders still remain (they may have been filled or partially filled)\n", remainingCount)
	}

	fmt.Println("\nCancel open orders example completed!")
}