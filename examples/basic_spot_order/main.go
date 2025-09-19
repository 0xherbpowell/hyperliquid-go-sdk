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

const (
	PURR            = "PURR/USDC"
	OTHER_COIN      = "@8"
	OTHER_COIN_NAME = "KORILA/USDC"
)

func main() {
	ctx := context.Background()

	// Setup clients (testnet)
	address, info, exchange, err := utils.Setup(constants.TestnetAPIURL, true, nil)
	if err != nil {
		log.Fatalf("Failed to setup: %v", err)
	}

	// Get and print spot user state and balance information
	spotUserState, err := info.SpotUserState(ctx, address)
	if err != nil {
		log.Fatalf("Failed to get spot user state: %v", err)
	}

	fmt.Println("Spot balances:")
	if len(spotUserState.Balances) > 0 {
		for _, balance := range spotUserState.Balances {
			balanceJSON, _ := json.MarshalIndent(balance, "", "  ")
			fmt.Println(string(balanceJSON))
		}
	} else {
		fmt.Println("No available token balances")
	}

	// Place an order that should rest by setting the price very low
	orderRequest := types.OrderRequest{
		Coin:       PURR,
		IsBuy:      true,
		Sz:         24,
		LimitPx:    0.5,
		OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: constants.TifGtc}},
		ReduceOnly: false,
	}

	fmt.Printf("Placing spot order: %+v\n", orderRequest)
	orderResult, err := exchange.Order(ctx, orderRequest, nil)
	if err != nil {
		log.Fatalf("Failed to place order: %v", err)
	}

	orderJSON, _ := json.MarshalIndent(orderResult, "", "  ")
	fmt.Printf("Order result: %s\n", orderJSON)

	// Parse the order result to get the order ID for PURR/USDC
	var purrOrderID int
	if err := parseOrderResult(orderResult, &purrOrderID); err != nil {
		log.Printf("Failed to parse PURR order result: %v", err)
	} else if purrOrderID > 0 {
		// Query the order status by oid
		orderStatus, err := info.QueryOrderByOid(ctx, address, purrOrderID)
		if err != nil {
			log.Printf("Failed to query order status: %v", err)
		} else {
			statusJSON, _ := json.MarshalIndent(orderStatus, "", "  ")
			fmt.Printf("Order status by oid: %s\n", statusJSON)
		}

		// Cancel the PURR order
		fmt.Printf("Cancelling PURR order %d\n", purrOrderID)
		cancelResult, err := exchange.Cancel(ctx, PURR, purrOrderID)
		if err != nil {
			log.Printf("Failed to cancel PURR order: %v", err)
		} else {
			cancelJSON, _ := json.MarshalIndent(cancelResult, "", "  ")
			fmt.Printf("Cancel result: %s\n", cancelJSON)
		}
	}

	// For other spot assets other than PURR/USDC use @{index} e.g. on testnet @8 is KORILA/USDC
	fmt.Printf("\nPlacing order for %s (%s)\n", OTHER_COIN, OTHER_COIN_NAME)
	otherOrderRequest := types.OrderRequest{
		Coin:       OTHER_COIN,
		IsBuy:      true,
		Sz:         1,
		LimitPx:    12,
		OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: constants.TifGtc}},
		ReduceOnly: false,
	}

	otherOrderResult, err := exchange.Order(ctx, otherOrderRequest, nil)
	if err != nil {
		log.Printf("Failed to place %s order: %v", OTHER_COIN, err)
		return
	}

	otherOrderJSON, _ := json.MarshalIndent(otherOrderResult, "", "  ")
	fmt.Printf("Other order result: %s\n", otherOrderJSON)

	// Parse the order result for the other coin
	var otherOrderID int
	if err := parseOrderResult(otherOrderResult, &otherOrderID); err != nil {
		log.Printf("Failed to parse %s order result: %v", OTHER_COIN, err)
	} else if otherOrderID > 0 {
		// The SDK now also supports using spot names, although be careful as they might not always be unique
		fmt.Printf("Cancelling %s order %d using full name %s\n", OTHER_COIN, otherOrderID, OTHER_COIN_NAME)
		cancelResult, err := exchange.Cancel(ctx, OTHER_COIN_NAME, otherOrderID)
		if err != nil {
			log.Printf("Failed to cancel %s order using name: %v", OTHER_COIN_NAME, err)

			// Try cancelling with the @index format
			fmt.Printf("Trying to cancel with %s format\n", OTHER_COIN)
			cancelResult, err = exchange.Cancel(ctx, OTHER_COIN, otherOrderID)
			if err != nil {
				log.Printf("Failed to cancel %s order: %v", OTHER_COIN, err)
			} else {
				cancelJSON, _ := json.MarshalIndent(cancelResult, "", "  ")
				fmt.Printf("Cancel result with @index: %s\n", cancelJSON)
			}
		} else {
			cancelJSON, _ := json.MarshalIndent(cancelResult, "", "  ")
			fmt.Printf("Cancel result with name: %s\n", cancelJSON)
		}
	}

	fmt.Println("\nSpot order example completed!")
}

// parseOrderResult extracts the order ID from the order result
func parseOrderResult(orderResult interface{}, orderID *int) error {
	orderResultMap, ok := orderResult.(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse order result")
	}

	status, ok := orderResultMap["status"].(string)
	if !ok || status != "ok" {
		return fmt.Errorf("order failed with status: %v", status)
	}

	response, ok := orderResultMap["response"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse order response")
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse order data")
	}

	statuses, ok := data["statuses"].([]interface{})
	if !ok || len(statuses) == 0 {
		return fmt.Errorf("failed to parse order statuses")
	}

	firstStatus, ok := statuses[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse first status")
	}

	// Check if order is resting and get order ID
	if resting, exists := firstStatus["resting"]; exists {
		restingMap, ok := resting.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse resting order")
		}

		oidFloat, ok := restingMap["oid"].(float64)
		if !ok {
			return fmt.Errorf("failed to parse order ID")
		}

		*orderID = int(oidFloat)
		fmt.Printf("Order placed with ID: %d\n", *orderID)
		return nil
	} else if filled, exists := firstStatus["filled"]; exists {
		fmt.Printf("Order was filled immediately: %+v\n", filled)
		return fmt.Errorf("order was filled, no resting order ID")
	} else if errorMsg, exists := firstStatus["error"]; exists {
		return fmt.Errorf("order failed with error: %+v", errorMsg)
	}

	return fmt.Errorf("unknown order status")
}
