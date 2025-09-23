package main

import (
	"fmt"
	"log"
	"time"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet (WebSocket enabled)
	address, info, _ := Setup(utils.TestnetAPIURL, false)

	fmt.Printf("Starting WebSocket subscriptions for address: %s\n", address)

	// Subscribe to various data streams
	// Some subscriptions do not return snapshots, so you will not receive a message until something happens

	// All mids subscription
	err := info.Subscribe(map[string]interface{}{"type": "allMids"}, func(data interface{}) {
		fmt.Printf("All Mids: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to allMids: %v", err)
	}

	// L2 Book for ETH
	err = info.Subscribe(map[string]interface{}{"type": "l2Book", "coin": "ETH"}, func(data interface{}) {
		fmt.Printf("L2 Book ETH: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to l2Book: %v", err)
	}

	// Trades for PURR/USDC
	err = info.Subscribe(map[string]interface{}{"type": "trades", "coin": "PURR/USDC"}, func(data interface{}) {
		fmt.Printf("Trades PURR/USDC: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to trades: %v", err)
	}

	// User events
	err = info.Subscribe(map[string]interface{}{"type": "userEvents", "user": address}, func(data interface{}) {
		fmt.Printf("User Events: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to userEvents: %v", err)
	}

	// User fills
	err = info.Subscribe(map[string]interface{}{"type": "userFills", "user": address}, func(data interface{}) {
		fmt.Printf("User Fills: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to userFills: %v", err)
	}

	// Candle data for ETH 1m
	err = info.Subscribe(map[string]interface{}{"type": "candle", "coin": "ETH", "interval": "1m"}, func(data interface{}) {
		fmt.Printf("Candle ETH 1m: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to candle: %v", err)
	}

	// Order updates
	err = info.Subscribe(map[string]interface{}{"type": "orderUpdates", "user": address}, func(data interface{}) {
		fmt.Printf("Order Updates: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to orderUpdates: %v", err)
	}

	// User funding updates
	err = info.Subscribe(map[string]interface{}{"type": "userFundings", "user": address}, func(data interface{}) {
		fmt.Printf("User Fundings: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to userFundings: %v", err)
	}

	// User non-funding ledger updates
	err = info.Subscribe(map[string]interface{}{"type": "userNonFundingLedgerUpdates", "user": address}, func(data interface{}) {
		fmt.Printf("User Non-Funding Ledger Updates: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to userNonFundingLedgerUpdates: %v", err)
	}

	// Web data 2
	err = info.Subscribe(map[string]interface{}{"type": "webData2", "user": address}, func(data interface{}) {
		fmt.Printf("Web Data 2: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to webData2: %v", err)
	}

	// Best bid offer for ETH
	err = info.Subscribe(map[string]interface{}{"type": "bbo", "coin": "ETH"}, func(data interface{}) {
		fmt.Printf("BBO ETH: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to bbo: %v", err)
	}

	// Active asset context for BTC (Perp)
	err = info.Subscribe(map[string]interface{}{"type": "activeAssetCtx", "coin": "BTC"}, func(data interface{}) {
		fmt.Printf("Active Asset Ctx BTC: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to activeAssetCtx BTC: %v", err)
	}

	// Active asset context for @1 (Spot)
	err = info.Subscribe(map[string]interface{}{"type": "activeAssetCtx", "coin": "@1"}, func(data interface{}) {
		fmt.Printf("Active Asset Ctx @1: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to activeAssetCtx @1: %v", err)
	}

	// Active asset data for BTC (Perp only)
	err = info.Subscribe(map[string]interface{}{"type": "activeAssetData", "user": address, "coin": "BTC"}, func(data interface{}) {
		fmt.Printf("Active Asset Data BTC: %+v\n", data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to activeAssetData: %v", err)
	}

	fmt.Println("\nWebSocket subscriptions established. Waiting for messages...")
	fmt.Println("Note: Some subscriptions don't return snapshots, so you'll only see messages when events happen.")
	
	// Keep the program running to receive WebSocket messages
	// In a real application, you might want to handle graceful shutdown
	select {
	case <-time.After(60 * time.Second):
		fmt.Println("\nStopping after 60 seconds...")
	}
}