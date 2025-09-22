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

	// Get the user state and print out position information
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	PrintPositions(userState)

	// Get current ETH price for reference
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

	// First, place a position order to have something to set TP/SL for
	// Place a small market buy order
	slippage := 0.02
	positionResult, err := exchange.MarketOrder(
		"ETH",             // coin
		true,              // isBuy
		0.05,              // small size
		&slippage,         // slippage
		GenerateCloid(),   // client order ID
	)
	if err != nil {
		log.Printf("Failed to place position order: %v", err)
		return
	}

	fmt.Println("Position order result:")
	PrintOrderResult(positionResult)

	// Set Take Profit order - 5% above current price
	takeProfitPrice := ethPrice * 1.05
	
	tpOrderResult, err := exchange.TriggerOrder(
		"ETH",                       // coin
		false,                       // isBuy (sell for TP)
		0.05,                        // size (same as position)
		takeProfitPrice,             // trigger price
		true,                        // isMarket (market order when triggered)
		types.TpslTp,                // take profit
		true,                        // reduce only
		GenerateCloid(),             // client order ID
	)
	if err != nil {
		log.Printf("Failed to place take profit order: %v", err)
	} else {
		fmt.Printf("Take Profit order at %f:\n", takeProfitPrice)
		PrintOrderResult(tpOrderResult)
	}

	// Set Stop Loss order - 3% below current price
	stopLossPrice := ethPrice * 0.97

	slOrderResult, err := exchange.TriggerOrder(
		"ETH",                       // coin
		false,                       // isBuy (sell for SL)
		0.05,                        // size (same as position)
		stopLossPrice,               // trigger price
		true,                        // isMarket (market order when triggered)
		types.TpslSl,                // stop loss
		true,                        // reduce only
		GenerateCloid(),             // client order ID
	)
	if err != nil {
		log.Printf("Failed to place stop loss order: %v", err)
	} else {
		fmt.Printf("Stop Loss order at %f:\n", stopLossPrice)
		PrintOrderResult(slOrderResult)
	}

	// Display the current open orders
	openOrders, err := info.OpenOrders(address, "")
	if err != nil {
		log.Printf("Failed to get open orders: %v", err)
	} else {
		fmt.Println("Current open orders:")
		PrintOrderResult(openOrders)
	}

	// Optional: Cancel the TP/SL orders (uncomment if needed)
	/*
	// Cancel Take Profit order if it's resting
	if tpOid, ok := GetRestingOid(tpOrderResult); ok {
		cancelTP, err := exchange.Cancel("ETH", tpOid)
		if err != nil {
			log.Printf("Failed to cancel TP order: %v", err)
		} else {
			fmt.Println("Cancelled Take Profit order:")
			PrintOrderResult(cancelTP)
		}
	}

	// Cancel Stop Loss order if it's resting
	if slOid, ok := GetRestingOid(slOrderResult); ok {
		cancelSL, err := exchange.Cancel("ETH", slOid)
		if err != nil {
			log.Printf("Failed to cancel SL order: %v", err)
		} else {
			fmt.Println("Cancelled Stop Loss order:")
			PrintOrderResult(cancelSL)
		}
	}
	*/

	fmt.Println("\nTP/SL example completed!")
	fmt.Println("Note: In a real scenario, monitor the position and these orders will trigger automatically when price reaches the specified levels.")
}