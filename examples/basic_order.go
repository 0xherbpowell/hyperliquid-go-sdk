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

	// Place an order that should rest by setting the price very low
	orderResult, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		0.2,                   // size
		1100.0,                // limit price (low price to ensure it rests)
		CreateGtcLimitOrder(), // order type
		false,                 // reduce only
		nil,                   // cloid
		nil,                   // builder info
	)
	if err != nil {
		log.Printf("Failed to place order: %v", err)
		return
	}

	fmt.Println("Order result:")
	PrintOrderResult(orderResult)

	// Query the order status by oid
	if oid, ok := GetRestingOid(orderResult); ok {
		orderStatus, err := info.OrderStatus(address, oid, "")
		if err != nil {
			log.Printf("Failed to get order status: %v", err)
		} else {
			fmt.Printf("Order status by oid: %+v\n", orderStatus)
		}

		// Cancel the order
		cancelResult, err := exchange.Cancel("ETH", oid)
		if err != nil {
			log.Printf("Failed to cancel order: %v", err)
		} else {
			fmt.Println("Cancel result:")
			PrintOrderResult(cancelResult)
		}
	}
}
