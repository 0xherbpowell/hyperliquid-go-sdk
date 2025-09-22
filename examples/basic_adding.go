package main

import (
	"fmt"
	"log"
	"time"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// This example demonstrates "adding" liquidity by placing limit orders
	// that are likely to rest in the order book and provide liquidity

	fmt.Println("Basic Adding Liquidity Example")
	fmt.Printf("Account: %s\n", address)

	// Get initial user state
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("\nInitial account state:")
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

	// Place multiple orders at different price levels to add liquidity
	// These orders are designed to provide liquidity to the market
	orders := []struct {
		side     bool    // true = buy, false = sell
		offset   float64 // price offset as percentage
		size     float64
		desc     string
	}{
		{true, -0.02, 0.01, "Buy order 2% below market"},   // Buy 2% below
		{true, -0.04, 0.015, "Buy order 4% below market"}, // Buy 4% below
		{true, -0.06, 0.02, "Buy order 6% below market"},  // Buy 6% below
		{false, 0.02, 0.01, "Sell order 2% above market"}, // Sell 2% above
		{false, 0.04, 0.015, "Sell order 4% above market"}, // Sell 4% above
		{false, 0.06, 0.02, "Sell order 6% above market"}, // Sell 6% above
	}

	var orderResults []map[string]interface{}
	var orderIds []int

	fmt.Println("\nPlacing liquidity-adding orders...")

	for i, order := range orders {
		price := ethPrice * (1 + order.offset)

		fmt.Printf("Placing %s at price %f\n", order.desc, price)

		result, err := exchange.Order(
			"ETH",                 // coin
			order.side,            // isBuy
			order.size,            // size
			price,                 // limit price
			CreateGtcLimitOrder(), // GTC order type for liquidity
			false,                 // reduce only
			GenerateCloid(),       // unique client order ID
			nil,                   // builder info
		)

		if err != nil {
			log.Printf("Failed to place order %d: %v", i+1, err)
			continue
		}

		fmt.Printf("Order %d result:\n", i+1)
		PrintOrderResult(result)

		orderResults = append(orderResults, result)

		// Extract order ID for later cancellation
		if oid, ok := GetRestingOid(result); ok {
			orderIds = append(orderIds, oid)
		}

		// Small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\nSuccessfully placed %d liquidity-adding orders\n", len(orderResults))

	// Get current order book to see our orders
	orderBook, err := info.L2Book("ETH", "", nil, nil)
	if err != nil {
		log.Printf("Failed to get order book: %v", err)
	} else {
		fmt.Println("\nCurrent ETH order book (showing our liquidity):")
		if book, ok := orderBook["levels"].([]interface{}); ok && len(book) >= 2 {
			if bids, ok := book[0].([]interface{}); ok {
				fmt.Println("Top bids:")
				for i, bid := range bids {
					if i >= 5 {
						break
					} // Show top 5
					if bidMap, ok := bid.(map[string]interface{}); ok {
						fmt.Printf("  %s @ %s\n", bidMap["sz"], bidMap["px"])
					}
				}
			}
			if asks, ok := book[1].([]interface{}); ok {
				fmt.Println("Top asks:")
				for i, ask := range asks {
					if i >= 5 {
						break
					} // Show top 5
					if askMap, ok := ask.(map[string]interface{}); ok {
						fmt.Printf("  %s @ %s\n", askMap["sz"], askMap["px"])
					}
				}
			}
		}
	}

	// Show open orders
	openOrders, err := info.OpenOrders(address, "")
	if err != nil {
		log.Printf("Failed to get open orders: %v", err)
	} else {
		fmt.Println("\nCurrent open orders:")
		if orders, ok := openOrders["orders"].([]interface{}); ok {
			fmt.Printf("Total open orders: %d\n", len(orders))
		}
	}

	// Wait a moment to let the orders potentially interact with the market
	fmt.Println("\nWaiting 5 seconds for potential fills...")
	time.Sleep(5 * time.Second)

	// Check for any fills
	fills, err := info.UserFills(address, "")
	if err != nil {
		log.Printf("Failed to get fills: %v", err)
	} else {
		if fillsData, ok := fills["fills"].([]interface{}); ok {
			fmt.Printf("Recent fills: %d\n", len(fillsData))
			if len(fillsData) > 0 {
				fmt.Println("Some of our liquidity-adding orders were filled!")
			}
		}
	}

	// Cancel remaining orders (cleanup)
	fmt.Println("\nCancelling remaining orders (cleanup)...")

	for i, oid := range orderIds {
		if oid > 0 {
			cancelResult, err := exchange.Cancel("ETH", oid)
			if err != nil {
				log.Printf("Failed to cancel order %d (oid: %d): %v", i+1, oid, err)
			} else {
				fmt.Printf("Cancelled order %d (oid: %d)\n", i+1, oid)
				_ = cancelResult // Suppress unused variable warning
			}
		}
	}

	fmt.Println("\nLiquidity adding example completed!")
	fmt.Println("Note: This example showed how to add liquidity to the market by placing")
	fmt.Println("limit orders at various price levels that are likely to rest in the book.")
}