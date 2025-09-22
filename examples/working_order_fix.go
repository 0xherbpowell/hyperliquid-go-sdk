package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
	"github.com/vmihailenco/msgpack/v5"
)

// WorkingOrder demonstrates the exact approach that should work
func WorkingOrder() error {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)
	
	fmt.Printf("Using address: %s\n", address)

	// Get current ETH price for reference
	mids, err := info.AllMids("")
	if err != nil {
		return fmt.Errorf("failed to get mids: %w", err)
	}
	
	ethPrice := "unknown"
	if price, exists := mids["ETH"]; exists {
		ethPrice = price
	}
	fmt.Printf("Current ETH price: %s\n", ethPrice)

	// Create the order using the exact same approach as the working examples
	// but with corrected parameters
	orderResult, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		0.01,                  // size (smaller size to be safe)
		3800.0,                // limit price (conservative price below market)
		CreateGtcLimitOrder(), // order type
		false,                 // reduce only
		nil,                   // cloid
		nil,                   // builder info
	)
	
	fmt.Printf("Order attempt result: %+v\n", orderResult)
	
	if err != nil {
		fmt.Printf("Order failed with error: %v\n", err)
		
		// Let's debug by examining what the SDK is actually sending
		fmt.Println("\n=== Debugging the issue ===")
		
		// Create order manually to see the exact format
		orderType := types.OrderType{
			Limit: &types.LimitOrderType{
				Tif: types.TifGtc,
			},
		}

		asset, _ := info.NameToAsset("ETH")
		orderRequest := types.OrderRequest{
			Coin:       "ETH",
			IsBuy:      true,
			Sz:         0.01,
			LimitPx:    3800.0,
			OrderType:  orderType,
			ReduceOnly: false,
			Cloid:      nil,
		}

		orderWire, err := utils.OrderRequestToOrderWire(orderRequest, asset)
		if err != nil {
			return fmt.Errorf("failed to convert to wire: %v", err)
		}

		fmt.Println("OrderWire structure:")
		wireJSON, _ := json.MarshalIndent(orderWire, "", "  ")
		fmt.Println(string(wireJSON))

		// Test msgpack serialization
		fmt.Println("\nMsgpack serialization test:")
		msgpackData, err := msgpack.Marshal(orderWire)
		if err != nil {
			fmt.Printf("Msgpack error: %v\n", err)
		} else {
			fmt.Printf("Msgpack length: %d bytes\n", len(msgpackData))
		}

		return err
	}

	fmt.Println("Order placed successfully!")
	return nil
}

// TryAlternativeFormat attempts an alternative format that might work
func TryAlternativeFormat() error {
	_, _, _ = Setup(utils.TestnetAPIURL, true)
	
	// Try to manually construct the exact payload that Python would send
	_ = time.Now().UnixMilli()
	
	// Create action exactly as Python does
	action := map[string]interface{}{
		"type": "order",
		"orders": []map[string]interface{}{
			{
				"a": 4,      // ETH asset ID
				"b": true,   // isBuy
				"p": "3800", // price as string
				"s": "0.01", // size as string
				"r": false,  // reduceOnly
				"t": map[string]interface{}{
					"limit": map[string]interface{}{
						"tif": "Gtc",
					},
				},
			},
		},
		"grouping": "na",
	}

	fmt.Println("Alternative action format:")
	actionJSON, _ := json.MarshalIndent(action, "", "  ")
	fmt.Println(string(actionJSON))

	// Test msgpack serialization of this format
	msgpackData, err := msgpack.Marshal(action)
	if err != nil {
		fmt.Printf("Msgpack serialization failed: %v\n", err)
	} else {
		fmt.Printf("Msgpack serialization successful: %d bytes\n", len(msgpackData))
	}

	fmt.Println("This format should work with proper EIP712 signing")
	return nil
}

func main() {
	fmt.Println("=== Working Order Fix ===\n")
	
	fmt.Println("1. Trying with the existing SDK approach:")
	if err := WorkingOrder(); err != nil {
		fmt.Printf("SDK approach failed: %v\n\n", err)
	}
	
	fmt.Println("2. Testing alternative format:")
	if err := TryAlternativeFormat(); err != nil {
		log.Fatalf("Alternative format test failed: %v", err)
	}

	fmt.Println("\n=== Analysis ===")
	fmt.Println("The issue is most likely one of these:")
	fmt.Println("1. The msgpack serialization format doesn't match what the API expects")
	fmt.Println("2. The EIP712 signing implementation has a bug")
	fmt.Println("3. The API has changed and the SDK is outdated")
	fmt.Println("4. There's a specific field ordering requirement in the JSON/msgpack")
	
	fmt.Println("\n=== Recommended Solutions ===")
	fmt.Println("1. Continue using Python SDK for orders (most reliable)")
	fmt.Println("2. Use Go SDK for read-only operations (works fine)")
	fmt.Println("3. File an issue with the Go SDK maintainers")
	fmt.Println("4. Try to find a newer version of the Go SDK")
}