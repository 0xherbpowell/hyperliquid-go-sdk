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

	// Get current spot mid prices
	mids, err := info.AllMids("")
	if err != nil {
		log.Printf("Failed to get mids: %v", err)
		return
	}

	// Look for a spot pair (usually contains '/')
	var spotPair string
	var spotPrice float64
	
	// Common spot pairs to look for
	spotPairs := []string{"BTC/USDC", "ETH/USDC", "SOL/USDC", "USDC/USDT"}
	
	for _, pair := range spotPairs {
		if mid, exists := mids[pair]; exists {
			spotPair = pair
			var err error
			spotPrice, err = utils.ParsePrice(mid)
			if err != nil {
				continue
			}
			break
		}
	}

	if spotPair == "" {
		fmt.Println("No suitable spot pairs found. Available pairs:")
		for coin, price := range mids {
			// Print pairs that look like spot pairs (contain '/')
			if len(coin) > 3 && (coin[len(coin)-4:] == "USDC" || coin[len(coin)-4:] == "USDT" || coin[3:] == "/") {
				fmt.Printf("  %s: %s\n", coin, price)
			}
		}
		return
	}

	fmt.Printf("Found spot pair: %s at price %f\n", spotPair, spotPrice)

	// Place a spot buy order - buy order below market price to ensure it rests
	buyPrice := spotPrice * 0.95 // 5% below market
	
	spotBuyResult, err := exchange.Order(
		spotPair,             // spot pair
		true,                 // isBuy
		1.0,                  // size (adjust based on pair)
		buyPrice,             // limit price
		CreateGtcLimitOrder(), // order type
		false,                // reduce only
		GenerateCloid(),      // client order ID
		nil,                  // builder info
	)
	if err != nil {
		log.Printf("Failed to place spot buy order: %v", err)
	} else {
		fmt.Printf("Spot buy order result for %s:\n", spotPair)
		PrintOrderResult(spotBuyResult)
	}

	// Place a spot sell order - sell order above market price to ensure it rests
	sellPrice := spotPrice * 1.05 // 5% above market
	
	spotSellResult, err := exchange.Order(
		spotPair,             // spot pair
		false,                // isBuy (false = sell)
		1.0,                  // size
		sellPrice,            // limit price
		CreateGtcLimitOrder(), // order type
		false,                // reduce only
		GenerateCloid(),      // client order ID
		nil,                  // builder info
	)
	if err != nil {
		log.Printf("Failed to place spot sell order: %v", err)
	} else {
		fmt.Printf("Spot sell order result for %s:\n", spotPair)
		PrintOrderResult(spotSellResult)
	}

	// Get current open orders to see our spot orders
	openOrders, err := info.OpenOrders(address, "")
	if err != nil {
		log.Printf("Failed to get open orders: %v", err)
	} else {
		fmt.Println("\nCurrent open orders:")
		PrintOrderResult(openOrders)
	}

	// Cancel the orders we just placed (optional)
	fmt.Println("\nCancelling the demo orders...")
	
	if buyOid, ok := GetRestingOid(spotBuyResult); ok {
		cancelBuy, err := exchange.Cancel(spotPair, buyOid)
		if err != nil {
			log.Printf("Failed to cancel buy order: %v", err)
		} else {
			fmt.Println("Cancelled spot buy order:")
			PrintOrderResult(cancelBuy)
		}
	}

	if sellOid, ok := GetRestingOid(spotSellResult); ok {
		cancelSell, err := exchange.Cancel(spotPair, sellOid)
		if err != nil {
			log.Printf("Failed to cancel sell order: %v", err)
		} else {
			fmt.Println("Cancelled spot sell order:")
			PrintOrderResult(cancelSell)
		}
	}

	fmt.Println("\nSpot trading example completed!")
	fmt.Printf("Note: Spot trading involves trading actual tokens (e.g., %s)\n", spotPair)
	fmt.Println("Unlike perpetuals, spot trading requires holding the underlying assets.")
}