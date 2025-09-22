package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// This example demonstrates how to schedule order cancellations
	// Orders can be automatically cancelled after a specified time period

	fmt.Println("Basic Schedule Cancel Example")
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

	// Example 1: Place an order that will be manually scheduled for cancellation
	fmt.Println("\n--- Example 1: Manual Schedule Cancel ---")

	orderPrice := ethPrice * 0.95 // 5% below market
	orderSize := 0.01             // Small size for testing

	fmt.Printf("Placing order that will be cancelled in 30 seconds\n")
	fmt.Printf("Order: %f ETH at %f\n", orderSize, orderPrice)

	orderResult, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		orderSize,             // size
		orderPrice,            // limit price
		CreateGtcLimitOrder(), // GTC order type
		false,                 // reduce only
		GenerateCloid(),       // unique client order ID
		nil,                   // builder info
	)

	if err != nil {
		log.Fatalf("Failed to place order: %v", err)
	}

	fmt.Println("Order placed successfully:")
	PrintOrderResult(orderResult)

	// Extract order ID for cancellation
	var orderOid int
	if oid, ok := GetRestingOid(orderResult); ok {
		orderOid = oid
		fmt.Printf("Order ID: %d\n", orderOid)
	} else {
		fmt.Println("Order was not resting - it may have been filled immediately")
		return
	}

	// Schedule cancellation after 30 seconds
	fmt.Printf("Scheduling cancellation of order %d in 30 seconds...\n", orderOid)

	// Start a goroutine to cancel the order after 30 seconds
	go func(oid int) {
		time.Sleep(30 * time.Second)

		fmt.Printf("\n[Scheduled Cancel] Time elapsed - cancelling order %d\n", oid)

		cancelResult, err := exchange.Cancel("ETH", oid)
		if err != nil {
			log.Printf("[Scheduled Cancel] Failed to cancel order %d: %v", oid, err)
		} else {
			fmt.Printf("[Scheduled Cancel] Successfully cancelled order %d\n", oid)
			_ = cancelResult // Suppress unused variable warning
		}
	}(orderOid)

	// Example 2: Place multiple orders with different cancellation schedules
	fmt.Println("\n--- Example 2: Multiple Orders with Different Schedules ---")

	orderConfigs := []struct {
		priceOffset float64
		size        float64
		cancelAfter time.Duration
		description string
	}{
		{-0.03, 0.015, 15 * time.Second, "3% below market, cancel in 15s"},
		{-0.06, 0.02, 45 * time.Second, "6% below market, cancel in 45s"},
		{0.03, 0.012, 25 * time.Second, "3% above market, cancel in 25s"},
	}

	var scheduledOrders []struct {
		oid         int
		cancelAfter time.Duration
		description string
	}

	for i, config := range orderConfigs {
		price := ethPrice * (1 + config.priceOffset)
		isBuy := config.priceOffset < 0 // Buy if price is below market

		fmt.Printf("Placing order %d: %s\n", i+1, config.description)
		fmt.Printf("  %f ETH at %f (is_buy: %t)\n", config.size, price, isBuy)

		result, err := exchange.Order(
			"ETH",                 // coin
			isBuy,                 // isBuy
			config.size,           // size
			price,                 // limit price
			CreateGtcLimitOrder(), // GTC order type
			false,                 // reduce only
			GenerateCloid(),       // unique client order ID
			nil,                   // builder info
		)

		if err != nil {
			log.Printf("Failed to place order %d: %v", i+1, err)
			continue
		}

		if oid, ok := GetRestingOid(result); ok {
			scheduledOrders = append(scheduledOrders, struct {
				oid         int
				cancelAfter time.Duration
				description string
			}{
				oid:         oid,
				cancelAfter: config.cancelAfter,
				description: config.description,
			})

			fmt.Printf("  Order placed with ID: %d\n", oid)
		} else {
			fmt.Printf("  Order %d was filled immediately\n", i+1)
		}

		time.Sleep(100 * time.Millisecond) // Small delay between orders
	}

	// Schedule cancellations for all orders
	for _, order := range scheduledOrders {
		go func(o struct {
			oid         int
			cancelAfter time.Duration
			description string
		}) {
			time.Sleep(o.cancelAfter)

			fmt.Printf("\n[Scheduled Cancel] Cancelling order %d (%s)\n",
				o.oid, o.description)

			cancelResult, err := exchange.Cancel("ETH", o.oid)
			if err != nil {
				log.Printf("[Scheduled Cancel] Failed to cancel order %d: %v", o.oid, err)
			} else {
				fmt.Printf("[Scheduled Cancel] Successfully cancelled order %d\n", o.oid)
				_ = cancelResult // Suppress unused variable warning
			}
		}(order)
	}

	// Example 3: Schedule cancel based on market conditions (simplified)
	fmt.Println("\n--- Example 3: Conditional Schedule Cancel ---")

	// Place an order and monitor market conditions
	conditionalPrice := ethPrice * 0.97 // 3% below market
	conditionalSize := 0.008

	fmt.Printf("Placing conditional order: %f ETH at %f\n", conditionalSize, conditionalPrice)
	fmt.Println("This order will be cancelled if the market moves too much")

	conditionalResult, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		conditionalSize,       // size
		conditionalPrice,      // limit price
		CreateGtcLimitOrder(), // GTC order type
		false,                 // reduce only
		GenerateCloid(),       // unique client order ID
		nil,                   // builder info
	)

	if err != nil {
		log.Printf("Failed to place conditional order: %v", err)
	} else if conditionalOid, ok := GetRestingOid(conditionalResult); ok {
		fmt.Printf("Conditional order placed with ID: %d\n", conditionalOid)

		// Start monitoring for cancellation conditions
		go func(oid int, originalPrice float64) {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			maxChecks := 12 // Check for 1 minute (12 * 5 seconds)
			checks := 0

			for range ticker.C {
				checks++
				if checks >= maxChecks {
					fmt.Printf("\n[Conditional Cancel] Time limit reached - cancelling order %d\n", oid)
					break
				}

				// Get current market price
				currentMids, err := info.AllMids("")
				if err != nil {
					continue
				}

				if currentEthMid, exists := currentMids["ETH"]; exists {
					if currentPrice, err := utils.ParsePrice(currentEthMid); err == nil {
						priceChange := (currentPrice - originalPrice) / originalPrice

						fmt.Printf("[Conditional Monitor] Current ETH price: %f (change: %.2f%%)\n",
							currentPrice, priceChange*100)

						// Cancel if price moved more than 2% from original
						if abs(priceChange) > 0.02 {
							fmt.Printf("\n[Conditional Cancel] Price moved %.2f%% - cancelling order %d\n",
								priceChange*100, oid)
							break
						}
					}
				}
			}

			// Cancel the order
			cancelResult, err := exchange.Cancel("ETH", oid)
			if err != nil {
				log.Printf("[Conditional Cancel] Failed to cancel order %d: %v", oid, err)
			} else {
				fmt.Printf("[Conditional Cancel] Successfully cancelled order %d\n", oid)
				_ = cancelResult // Suppress unused variable warning
			}
		}(conditionalOid, ethPrice)
	}

	// Show current open orders
	fmt.Println("\n--- Current Open Orders ---")
	openOrders, err := info.OpenOrders(address, "")
	if err != nil {
		log.Printf("Failed to get open orders: %v", err)
	} else {
		if orders, ok := openOrders["orders"].([]interface{}); ok {
			fmt.Printf("Total open orders: %d\n", len(orders))

			for i, order := range orders {
				if orderMap, ok := order.(map[string]interface{}); ok {
					fmt.Printf("Order %d: %s %s @ %s (oid: %v)\n",
						i+1,
						orderMap["sz"],
						orderMap["coin"],
						orderMap["limitPx"],
						orderMap["oid"])
				}
			}
		}
	}

	// Wait for all scheduled cancellations to complete
	fmt.Println("\nWaiting for scheduled cancellations to complete...")
	fmt.Println("(This may take up to 1 minute)")
	time.Sleep(65 * time.Second) // Wait a bit longer than the longest schedule

	// Final check of open orders
	fmt.Println("\n--- Final Open Orders Check ---")
	finalOrders, err := info.OpenOrders(address, "")
	if err != nil {
		log.Printf("Failed to get final open orders: %v", err)
	} else {
		if orders, ok := finalOrders["orders"].([]interface{}); ok {
			fmt.Printf("Remaining open orders: %d\n", len(orders))

			if len(orders) > 0 {
				fmt.Println("Cancelling any remaining orders...")
				for _, order := range orders {
					if orderMap, ok := order.(map[string]interface{}); ok {
						if oid, ok := orderMap["oid"].(float64); ok {
							exchange.Cancel("ETH", int(oid))
						}
					}
				}
			}
		}
	}

	fmt.Println("\nSchedule cancel example completed!")
	fmt.Println("Note: This example demonstrated:")
	fmt.Println("1. Manually scheduling order cancellations after a time delay")
	fmt.Println("2. Multiple orders with different cancellation schedules")
	fmt.Println("3. Conditional cancellation based on market conditions")
	fmt.Println("4. Monitoring and cleanup of scheduled orders")
}

// Helper function to get absolute value
func abs(x float64) float64 {
	return math.Abs(x)
}