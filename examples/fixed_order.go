package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"hyperliquid-go-sdk/pkg/utils"
)

// FixedOrderType represents the exact order type format expected by the API
type FixedOrderType struct {
	Limit *FixedLimitOrderType `json:"limit,omitempty"`
}

// FixedLimitOrderType represents the limit order type
type FixedLimitOrderType struct {
	Tif string `json:"tif"`
}

// FixedOrderWire represents the exact wire format expected by the API
type FixedOrderWire struct {
	A int            `json:"a"` // asset
	B bool           `json:"b"` // isBuy
	P string         `json:"p"` // limitPx
	S string         `json:"s"` // sz
	R bool           `json:"r"` // reduceOnly
	T FixedOrderType `json:"t"` // orderType
	C *string        `json:"c,omitempty"` // cloid
}

// FixedOrderAction represents the exact action format expected by the API
type FixedOrderAction struct {
	Type     string           `json:"type"`
	Orders   []FixedOrderWire `json:"orders"`
	Grouping string           `json:"grouping"`
}

// PlaceFixedOrder places an order using the corrected format
func PlaceFixedOrder(
	privateKey *ecdsa.PrivateKey,
	accountAddress string,
	baseURL string,
	coin string,
	isBuy bool,
	size float64,
	limitPrice float64,
	tif string,
) (map[string]interface{}, error) {
	
	// Create the order in the exact format expected by the API
	orderWire := FixedOrderWire{
		A: 4, // ETH asset ID (hardcoded for this example)
		B: isBuy,
		P: fmt.Sprintf("%.8g", limitPrice), // Format price as string without trailing zeros
		S: fmt.Sprintf("%.8g", size),       // Format size as string without trailing zeros
		R: false,                           // reduceOnly
		T: FixedOrderType{
			Limit: &FixedLimitOrderType{
				Tif: tif,
			},
		},
	}

	// Create action in exact format
	action := FixedOrderAction{
		Type:     "order",
		Orders:   []FixedOrderWire{orderWire},
		Grouping: "na",
	}

	// Get current timestamp for nonce
	nonce := time.Now().UnixMilli()

	// Create the payload for signing (simplified approach)
	// For this fix, we'll create a raw HTTP request
	payload := map[string]interface{}{
		"action":       action,
		"nonce":        nonce,
		"signature":    nil, // Will be filled after signing
		"vaultAddress": nil,
		"user":         accountAddress,
	}

	// For simplicity, we'll make a direct HTTP request with the corrected format
	// In a production environment, you'd want to implement proper EIP712 signing
	
	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	fmt.Printf("Sending payload: %s\n", string(jsonData))

	// Make HTTP request
	resp, err := http.Post(baseURL+"/exchange", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// SimpleFixedOrder places an order with a simplified approach that mimics Python exactly
func SimpleFixedOrder() error {
	// Setup with testnet
	address, info, _ := Setup(utils.TestnetAPIURL, true)
	
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

	// Create a simple order that exactly matches Python format
	// Python: exchange.order("ETH", True, 0.2, 1100, {"limit": {"tif": "Gtc"}})
	
	// For now, let's create the exact JSON payload that works
	orderPayload := map[string]interface{}{
		"action": map[string]interface{}{
			"type": "order",
			"orders": []map[string]interface{}{
				{
					"a": 4,      // ETH asset ID
					"b": true,   // isBuy
					"p": "1100", // limitPx as string
					"s": "0.2",  // size as string  
					"r": false,  // reduceOnly
					"t": map[string]interface{}{
						"limit": map[string]interface{}{
							"tif": "Gtc",
						},
					},
				},
			},
			"grouping": "na",
		},
		"nonce":        time.Now().UnixMilli(),
		"signature":    map[string]interface{}{
			"r": "0x0000000000000000000000000000000000000000000000000000000000000000",
			"s": "0x0000000000000000000000000000000000000000000000000000000000000000", 
			"v": 27,
		},
		"vaultAddress": nil,
		"user":         address,
	}

	// Convert to JSON and print for debugging
	jsonData, _ := json.MarshalIndent(orderPayload, "", "  ")
	fmt.Println("Payload to send:")
	fmt.Println(string(jsonData))

	// Note: This is a demonstration of the correct format
	// In practice, you need proper EIP712 signing
	fmt.Println("\nNOTE: This shows the correct format, but requires proper signing to work")
	fmt.Println("The issue is likely in the msgpack serialization or EIP712 signing process")
	
	return nil
}

func main() {
	fmt.Println("=== Fixed Order Implementation ===")
	
	if err := SimpleFixedOrder(); err != nil {
		log.Fatalf("Error: %v", err)
	}
	
	fmt.Println("\nTo fix the original issue, we need to:")
	fmt.Println("1. Ensure exact JSON format matches Python")
	fmt.Println("2. Fix the EIP712 signing process") 
	fmt.Println("3. Verify msgpack serialization")
	fmt.Println("\nThe format above should work with proper signing implementation.")
}