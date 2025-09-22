package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hyperliquid-go-sdk/pkg/client"
	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Create info client with WebSocket enabled
	timeout := 30 * time.Second
	info, err := client.NewInfo(
		utils.TestnetAPIURL, // Use TestnetAPIURL for testing
		&timeout,
		false, // skipWS = false to enable WebSocket
		nil,   // meta
		nil,   // spotMeta
		nil,   // perpDexs
	)
	if err != nil {
		log.Fatalf("Failed to create info client: %v", err)
	}
	
	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	defer func() {
		if err := info.DisconnectWebsocket(); err != nil {
			log.Printf("Error disconnecting WebSocket: %v", err)
		}
	}()
	
	// Example 1: Subscribe to all mid prices
	fmt.Println("=== Subscribing to All Mids ===")
	allMidsSubscription := []types.Subscription{
		{Type: "allMids"},
	}
	
	allMidsCallback := func(data interface{}) {
		if msgData, ok := data.(map[string]interface{}); ok {
			if dataField, ok := msgData["data"].(map[string]interface{}); ok {
				if mids, ok := dataField["mids"].(map[string]interface{}); ok {
					fmt.Printf("Mid prices update - showing first 5:\n")
					count := 0
					for coin, price := range mids {
						if count >= 5 {
							break
						}
						fmt.Printf("  %s: %v\n", coin, price)
						count++
					}
					if len(mids) > 5 {
						fmt.Printf("  ... and %d more\n", len(mids)-5)
					}
					fmt.Println()
				}
			}
		}
	}
	
	if err := info.Subscribe(allMidsSubscription, allMidsCallback); err != nil {
		log.Printf("Failed to subscribe to all mids: %v", err)
	}
	
	// Example 2: Subscribe to ETH order book
	fmt.Println("=== Subscribing to ETH L2 Book ===")
	l2BookSubscription := []types.Subscription{
		{Type: "l2Book", Coin: "ETH"},
	}
	
	l2BookCallback := func(data interface{}) {
		if msgData, ok := data.(map[string]interface{}); ok {
			if dataField, ok := msgData["data"].(map[string]interface{}); ok {
				if coin, ok := dataField["coin"].(string); ok {
					fmt.Printf("Order book update for %s:\n", coin)
					
					if levels, ok := dataField["levels"].([]interface{}); ok && len(levels) >= 2 {
						// Bids (index 0) and Asks (index 1)
						if bids, ok := levels[0].([]interface{}); ok && len(bids) > 0 {
							fmt.Printf("  Best bid: ")
							if bid, ok := bids[0].(map[string]interface{}); ok {
								fmt.Printf("%v @ %v\n", bid["sz"], bid["px"])
							}
						}
						
						if asks, ok := levels[1].([]interface{}); ok && len(asks) > 0 {
							fmt.Printf("  Best ask: ")
							if ask, ok := asks[0].(map[string]interface{}); ok {
								fmt.Printf("%v @ %v\n", ask["sz"], ask["px"])
							}
						}
					}
					fmt.Println()
				}
			}
		}
	}
	
	if err := info.Subscribe(l2BookSubscription, l2BookCallback); err != nil {
		log.Printf("Failed to subscribe to L2 book: %v", err)
	}
	
	// Example 3: Subscribe to ETH trades
	fmt.Println("=== Subscribing to ETH Trades ===")
	tradesSubscription := []types.Subscription{
		{Type: "trades", Coin: "ETH"},
	}
	
	tradesCallback := func(data interface{}) {
		if msgData, ok := data.(map[string]interface{}); ok {
			if dataField, ok := msgData["data"].([]interface{}); ok {
				for _, tradeInterface := range dataField {
					if trade, ok := tradeInterface.(map[string]interface{}); ok {
						fmt.Printf("Trade: %s %v @ %v (side: %v)\n",
							trade["coin"],
							trade["sz"],
							trade["px"],
							trade["side"])
					}
				}
				fmt.Println()
			}
		}
	}
	
	if err := info.Subscribe(tradesSubscription, tradesCallback); err != nil {
		log.Printf("Failed to subscribe to trades: %v", err)
	}
	
	// Example 4: Subscribe to user events (requires address)
	// Uncomment and set your address to test user-specific subscriptions
	/*
	userAddress := "0x..." // Your wallet address
	fmt.Printf("=== Subscribing to User Events for %s ===\n", userAddress)
	userSubscription := []types.Subscription{
		{Type: "userEvents", User: userAddress},
	}
	
	userCallback := func(data interface{}) {
		fmt.Printf("User event: %+v\n", data)
	}
	
	if err := info.Subscribe(userSubscription, userCallback); err != nil {
		log.Printf("Failed to subscribe to user events: %v", err)
	}
	*/
	
	fmt.Println("WebSocket subscriptions active. Press Ctrl+C to exit...")
	
	// Wait for interrupt signal
	<-c
	fmt.Println("\nShutting down...")
	
	// Optionally unsubscribe before closing (automatic on disconnect)
	fmt.Println("Unsubscribing from streams...")
	if err := info.Unsubscribe(allMidsSubscription); err != nil {
		log.Printf("Error unsubscribing from all mids: %v", err)
	}
	if err := info.Unsubscribe(l2BookSubscription); err != nil {
		log.Printf("Error unsubscribing from L2 book: %v", err)
	}
	if err := info.Unsubscribe(tradesSubscription); err != nil {
		log.Printf("Error unsubscribing from trades: %v", err)
	}
	
	fmt.Println("WebSocket example completed!")
}