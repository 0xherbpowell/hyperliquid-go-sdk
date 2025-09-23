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

	// Get ETH asset ID to ensure we're using the correct one
	ethAsset, err := info.NameToAsset("ETH")
	if err != nil {
		log.Fatalf("Failed to get ETH asset ID: %v", err)
	}
	fmt.Printf("ETH Asset ID: %d\n", ethAsset)

	// Get current ETH price to place order below market
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

	// Place an order well below market price to ensure it rests
	orderPrice := ethPrice * 0.5 // 50% below market to ensure it rests
	orderResult, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		0.01,                  // size (smaller, safer size)
		orderPrice,            // limit price (well below market)
		CreateGtcLimitOrder(), // order type
		false,                 // reduce only
		nil,                   // cloid
		nil,                   // builder info
	)
	if err != nil {
		log.Printf("Failed to place order: %v", err)
		// Try to get more detailed error info
		fmt.Printf("Error details: %+v\n", err)
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
