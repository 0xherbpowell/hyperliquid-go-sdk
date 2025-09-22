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

	PrintPositions(userState)

	// Create a unique client order ID
	cloid := GenerateCloid()
	fmt.Printf("Using client order ID: %s\n", cloid.String())

	// Place an order with client order ID
	orderResult, err := exchange.Order(
		"ETH",                // coin
		true,                 // isBuy
		0.2,                  // size
		1100.0,               // limit price (low price to ensure it rests)
		CreateGtcLimitOrder(), // order type
		false,                // reduce only
		cloid,                // client order ID
		nil,                  // builder info
	)
	if err != nil {
		log.Printf("Failed to place order: %v", err)
		return
	}

	fmt.Println("Order result:")
	PrintOrderResult(orderResult)

	// Wait a moment for the order to be processed
	// time.Sleep(1 * time.Second)

	// Query the order status using the client order ID
	// Note: The Python SDK has a query_order_by_cloid method, but it's not always available
	// We'll use the regular order status check with OID if we can get it
	if oid, ok := GetRestingOid(orderResult); ok {
		orderStatus, err := info.OrderStatus(address, oid, "")
		if err != nil {
			log.Printf("Failed to get order status: %v", err)
		} else {
			fmt.Printf("Order status: %+v\n", orderStatus)
		}

		// Cancel the order using client order ID
		cancelResult, err := exchange.CancelByCloid("ETH", cloid)
		if err != nil {
			log.Printf("Failed to cancel order by cloid: %v", err)
		} else {
			fmt.Println("Cancel by cloid result:")
			PrintOrderResult(cancelResult)
		}
	} else {
		fmt.Println("Order may have been filled immediately or rejected")
	}
}