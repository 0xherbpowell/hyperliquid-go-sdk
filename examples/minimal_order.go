package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	fmt.Printf("Testing minimal order for address: %s\n", address)

	// Get current ETH price to place a reasonable order
	mids, err := info.AllMids("")
	if err != nil {
		log.Fatalf("Failed to get mids: %v", err)
	}

	ethPrice := mids["ETH"]
	fmt.Printf("Current ETH price: %s\n", ethPrice)

	// Try different order configurations to see which works

	fmt.Println("\n=== Test 1: Very simple limit order ===")
	orderResult, err := exchange.LimitOrder(
		"ETH",                // coin
		true,                 // isBuy  
		0.01,                 // size
		1000.0,               // limit price (very low to ensure it doesn't fill)
		types.TifGtc,        // time in force
		false,                // reduce only
		nil,                  // cloid
	)
	if err != nil {
		log.Printf("LimitOrder failed: %v", err)
	} else {
		fmt.Println("✓ LimitOrder succeeded!")
		PrintOrderResult(orderResult)
	}

	fmt.Println("\n=== Test 2: Using Order function directly ===")
	orderResult2, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		0.01,                  // size  
		1000.0,                // limit price
		CreateGtcLimitOrder(), // order type
		false,                 // reduce only
		nil,                   // cloid
		nil,                   // builder info
	)
	if err != nil {
		log.Printf("Order function failed: %v", err)
	} else {
		fmt.Println("✓ Order function succeeded!")
		PrintOrderResult(orderResult2)
	}

	fmt.Println("\n=== Test 3: Try BTC instead ===")
	orderResult3, err := exchange.Order(
		"BTC",                 // coin
		true,                  // isBuy
		0.001,                 // size (smaller for BTC)
		30000.0,               // limit price (low for BTC)
		CreateGtcLimitOrder(), // order type
		false,                 // reduce only
		nil,                   // cloid
		nil,                   // builder info
	)
	if err != nil {
		log.Printf("BTC order failed: %v", err)
	} else {
		fmt.Println("✓ BTC order succeeded!")
		PrintOrderResult(orderResult3)
	}

	fmt.Println("\n=== Test 4: Try with different size format ===")
	// Use a size that converts cleanly to string
	orderResult4, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		0.1,                   // size (0.1 should convert cleanly)
		1000.0,                // limit price
		CreateGtcLimitOrder(), // order type
		false,                 // reduce only
		nil,                   // cloid
		nil,                   // builder info
	)
	if err != nil {
		log.Printf("Clean size order failed: %v", err)
	} else {
		fmt.Println("✓ Clean size order succeeded!")
		PrintOrderResult(orderResult4)
	}

	fmt.Println("\n=== Test 5: Try with IOC order type ===")
	orderResult5, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		0.01,                  // size
		1000.0,                // limit price
		CreateIocLimitOrder(), // IOC order type
		false,                 // reduce only
		nil,                   // cloid
		nil,                   // builder info
	)
	if err != nil {
		log.Printf("IOC order failed: %v", err)
	} else {
		fmt.Println("✓ IOC order succeeded!")
		PrintOrderResult(orderResult5)
	}

	fmt.Println("\n=== Order Test Complete ===")
}