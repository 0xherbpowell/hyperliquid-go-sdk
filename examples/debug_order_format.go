package main

import (
	"encoding/json"
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, _ := Setup(utils.TestnetAPIURL, true)

	fmt.Printf("Using address: %s\n", address)

	// Create the same order as Python: ETH, True, 0.2, 1100, {"limit": {"tif": "Gtc"}}
	orderType := types.OrderType{
		Limit: &types.LimitOrderType{
			Tif: types.TifGtc,
		},
	}

	// Get asset ID for ETH
	asset, err := info.NameToAsset("ETH")
	if err != nil {
		log.Fatalf("Failed to get asset for ETH: %v", err)
	}

	fmt.Printf("ETH asset ID: %d\n", asset)

	// Create order request
	orderRequest := types.OrderRequest{
		Coin:       "ETH",
		IsBuy:      true,
		Sz:         0.2,
		LimitPx:    1100.0,
		OrderType:  orderType,
		ReduceOnly: false,
		Cloid:      nil,
	}

	// Convert to wire format
	orderWire, err := utils.OrderRequestToOrderWire(orderRequest, asset)
	if err != nil {
		log.Fatalf("Failed to convert to wire format: %v", err)
	}

	fmt.Println("Order wire format:")
	wireJSON, _ := json.MarshalIndent(orderWire, "", "  ")
	fmt.Println(string(wireJSON))

	// Create action
	action := utils.OrderWiresToOrderAction([]types.OrderWire{orderWire}, nil)

	fmt.Println("\nAction format:")
	actionJSON, _ := json.MarshalIndent(action, "", "  ")
	fmt.Println(string(actionJSON))

	// Show what Python would send (for comparison)
	fmt.Println("\nPython equivalent would be:")
	fmt.Println(`{
  "type": "order", 
  "orders": [{
    "a": 4, 
    "b": true, 
    "p": "1100", 
    "s": "0.2", 
    "r": false, 
    "t": {"limit": {"tif": "Gtc"}}
  }], 
  "grouping": "na"
}`)
}