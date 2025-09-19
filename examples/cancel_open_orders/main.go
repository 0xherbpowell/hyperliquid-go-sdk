package main

import (
	"context"
	"fmt"
	"log"

	"hyperliquid-go-sdk/examples/utils"
	"hyperliquid-go-sdk/pkg/constants"
)

func main() {
	ctx := context.Background()

	// Setup clients (testnet)
	address, info, exchange, err := utils.Setup(constants.TestnetAPIURL, true, nil)
	if err != nil {
		log.Fatalf("Failed to setup: %v", err)
	}

	// Get all open orders
	openOrders, err := info.OpenOrders(ctx, address, "")
	if err != nil {
		log.Fatalf("Failed to get open orders: %v", err)
	}

	if len(openOrders) == 0 {
		fmt.Println("No open orders to cancel")
		return
	}

	fmt.Printf("Found %d open orders to cancel\n", len(openOrders))

	// Cancel each open order
	for _, openOrder := range openOrders {
		fmt.Printf("Cancelling order: Coin=%s, OID=%d, Side=%s, Size=%s, Price=%s\n",
			openOrder.Coin, openOrder.Oid, openOrder.Side, openOrder.Sz, openOrder.LimitPx)

		result, err := exchange.Cancel(ctx, openOrder.Coin, openOrder.Oid)
		if err != nil {
			log.Printf("Failed to cancel order %d: %v", openOrder.Oid, err)
			continue
		}

		// Check if cancellation was successful
		if resultMap, ok := result.(map[string]interface{}); ok {
			if status, ok := resultMap["status"].(string); ok && status == "ok" {
				fmt.Printf("Successfully cancelled order %d\n", openOrder.Oid)
			} else {
				fmt.Printf("Failed to cancel order %d: %+v\n", openOrder.Oid, result)
			}
		} else {
			fmt.Printf("Unexpected response format for order %d: %+v\n", openOrder.Oid, result)
		}
	}

	fmt.Println("Finished cancelling orders")
}
