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

	// Place a market order using IOC (Immediate or Cancel) order type
	// This simulates a market order by using slippage protection
	slippage := 0.03 // 3% slippage
	orderResult, err := exchange.MarketOrder(
		"ETH",     // coin
		true,      // isBuy
		0.1,       // size
		&slippage, // slippage
		GenerateCloid(), // client order ID
	)
	if err != nil {
		log.Printf("Failed to place market order: %v", err)
		return
	}

	fmt.Println("Market order result:")
	PrintOrderResult(orderResult)

	// Alternative approach: Place a limit order with IOC (acts like market order)
	// Get current mid price first
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

	// Get ETH asset ID for tick size calculation
	ethAsset, err := info.NameToAsset("ETH")
	if err != nil {
		log.Printf("Failed to get ETH asset ID: %v", err)
		return
	}

	// Add some slippage to ensure execution
	rawLimitPrice := ethPrice * 1.02 // 2% above mid for buy order
	limitPrice := RoundToTickSize(rawLimitPrice, ethAsset) // Round to proper tick size

	fmt.Printf("ETH mid: %f, Raw limit: %f, Rounded limit: %f\n", ethPrice, rawLimitPrice, limitPrice)

	orderResult2, err := exchange.Order(
		"ETH",               // coin
		true,                // isBuy
		0.05,                // size (smaller size)
		limitPrice,          // limit price with slippage (tick-aligned)
		CreateIocLimitOrder(), // IOC order type
		false,               // reduce only
		GenerateCloid(),     // client order ID
		nil,                 // builder info
	)
	if err != nil {
		log.Printf("Failed to place IOC order: %v", err)
		return
	}

	fmt.Println("IOC order result:")
	PrintOrderResult(orderResult2)
}