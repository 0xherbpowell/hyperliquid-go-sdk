package main

import (
	"fmt"
	"hyperliquid-go-sdk/pkg/utils"
	"log"
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

	// Get BTC asset ID to ensure we're using the correct one
	btcAsset, err := info.NameToAsset("BTC")
	if err != nil {
		log.Fatalf("Failed to get BTC asset ID: %v", err)
	}
	fmt.Printf("BTC Asset ID: %d\n", btcAsset)

	// Get current BTC price to place order below market
	mids, err := info.AllMids("")
	if err != nil {
		log.Printf("Failed to get mids: %v", err)
		return
	}

	btcMid, exists := mids["BTC"]
	if !exists {
		log.Printf("BTC mid price not found")
		return
	}

	btcPrice, err := utils.ParsePrice(btcMid)
	if err != nil {
		log.Printf("Failed to parse BTC price: %v", err)
		return
	}

	fmt.Printf("Current BTC price: %f\n", btcPrice)

	// Place an order well below market price to ensure it rests
	rawOrderPrice := btcPrice * 0.8 // 20% below market to ensure it rests
	orderPrice := RoundToTickSize(rawOrderPrice, "BTC", info) // Round to proper tick size
	
	fmt.Printf("Raw order price: %f, Rounded price: %f\n", rawOrderPrice, orderPrice)
	
	orderResult, err := exchange.Order(
		"BTC",                 // coin
		true,                  // isBuy
		0.001,                 // size (smaller size for BTC)
		orderPrice,            // limit price (well below market, tick-aligned)
		CreateGtcLimitOrder(), // order type
		false,                 // reduce only
		GenerateCloid(),       // unique client order ID
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
		cancelResult, err := exchange.Cancel("BTC", oid)
		if err != nil {
			log.Printf("Failed to cancel order: %v", err)
		} else {
			fmt.Println("Cancel result:")
			PrintOrderResult(cancelResult)
		}
	}
}
