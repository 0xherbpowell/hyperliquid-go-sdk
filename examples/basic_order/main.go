package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"hyperliquid-go-sdk/examples/utils"
	"hyperliquid-go-sdk/pkg/constants"
	"hyperliquid-go-sdk/pkg/types"
)

func main() {
	ctx := context.Background()

	// Setup clients (testnet)
	address, info, exchange, err := utils.Setup(constants.TestnetAPIURL, true, nil)
	if err != nil {
		log.Fatalf("Failed to setup: %v", err)
	}

	// Get and print user state and position information
	userState, err := info.UserState(ctx, address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("Positions:")
	if len(userState.AssetPositions) > 0 {
		for _, assetPosition := range userState.AssetPositions {
			positionJSON, _ := json.MarshalIndent(assetPosition.Position, "", "  ")
			fmt.Println(string(positionJSON))
		}
	} else {
		fmt.Println("No open positions")
	}

	// Place an order that should rest by setting the price very low
	orderRequest := types.OrderRequest{
		Coin:       "ETH",
		IsBuy:      true,
		Sz:         0.2,
		LimitPx:    1100,
		OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: constants.TifGtc}},
		ReduceOnly: false,
	}

	fmt.Printf("Placing order: %+v\n", orderRequest)
	orderResult, err := exchange.Order(ctx, orderRequest, nil)
	if err != nil {
		log.Fatalf("Failed to place order: %v", err)
	}

	orderJSON, _ := json.MarshalIndent(orderResult, "", "  ")
	fmt.Printf("Order result: %s\n", orderJSON)

	// Parse the order result to get the order ID
	orderResultMap, ok := orderResult.(map[string]interface{})
	if !ok {
		log.Fatal("Failed to parse order result")
	}

	status, ok := orderResultMap["status"].(string)
	if !ok || status != "ok" {
		fmt.Printf("Order failed with result: %s\n", orderJSON)
		return
	}

	response, ok := orderResultMap["response"].(map[string]interface{})
	if !ok {
		log.Fatal("Failed to parse order response")
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		log.Fatal("Failed to parse order data")
	}

	statuses, ok := data["statuses"].([]interface{})
	if !ok || len(statuses) == 0 {
		log.Fatal("Failed to parse order statuses")
	}

	firstStatus, ok := statuses[0].(map[string]interface{})
	if !ok {
		log.Fatal("Failed to parse first status")
	}

	// Check if order is resting and get order ID
	if resting, exists := firstStatus["resting"]; exists {
		restingMap, ok := resting.(map[string]interface{})
		if !ok {
			log.Fatal("Failed to parse resting order")
		}

		oidFloat, ok := restingMap["oid"].(float64)
		if !ok {
			log.Fatal("Failed to parse order ID")
		}

		oid := int(oidFloat)
		fmt.Printf("Order placed with ID: %d\n", oid)

		// Query the order status by oid
		orderStatus, err := info.QueryOrderByOid(ctx, address, oid)
		if err != nil {
			log.Printf("Failed to query order status: %v", err)
		} else {
			statusJSON, _ := json.MarshalIndent(orderStatus, "", "  ")
			fmt.Printf("Order status by oid: %s\n", statusJSON)
		}

		// Cancel the order
		fmt.Printf("Cancelling order %d\n", oid)
		cancelResult, err := exchange.Cancel(ctx, "ETH", oid)
		if err != nil {
			log.Printf("Failed to cancel order: %v", err)
		} else {
			cancelJSON, _ := json.MarshalIndent(cancelResult, "", "  ")
			fmt.Printf("Cancel result: %s\n", cancelJSON)
		}
	} else if filled, exists := firstStatus["filled"]; exists {
		fmt.Printf("Order was filled: %+v\n", filled)
	} else if errorMsg, exists := firstStatus["error"]; exists {
		fmt.Printf("Order failed with error: %+v\n", errorMsg)
	}
}
