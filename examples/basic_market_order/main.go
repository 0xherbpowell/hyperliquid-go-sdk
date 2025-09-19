package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"hyperliquid-go-sdk/examples/utils"
	"hyperliquid-go-sdk/pkg/constants"
)

func main() {
	ctx := context.Background()

	// Setup clients (testnet)
	_, _, exchange, err := utils.Setup(constants.TestnetAPIURL, true, nil)
	if err != nil {
		log.Fatalf("Failed to setup: %v", err)
	}

	coin := "ETH"
	isBuy := false
	sz := 0.05
	slippage := 0.01 // 1% slippage

	fmt.Printf("Attempting to Market %s %f %s with %f%% slippage\n",
		map[bool]string{true: "Buy", false: "Sell"}[isBuy], sz, coin, slippage*100)

	// Place market order to open position
	orderResult, err := exchange.MarketOpen(ctx, coin, isBuy, sz, nil, &slippage, nil, nil)
	if err != nil {
		log.Fatalf("Failed to place market order: %v", err)
	}

	orderResultMap, ok := orderResult.(map[string]interface{})
	if !ok {
		log.Fatal("Failed to parse order result")
	}

	if status, ok := orderResultMap["status"].(string); ok && status == "ok" {
		if response, ok := orderResultMap["response"].(map[string]interface{}); ok {
			if data, ok := response["data"].(map[string]interface{}); ok {
				if statuses, ok := data["statuses"].([]interface{}); ok {
					for _, statusInterface := range statuses {
						if statusMap, ok := statusInterface.(map[string]interface{}); ok {
							if filled, exists := statusMap["filled"]; exists {
								if filledMap, ok := filled.(map[string]interface{}); ok {
									oid := filledMap["oid"]
									totalSz := filledMap["totalSz"]
									avgPx := filledMap["avgPx"]
									fmt.Printf("Order #%v filled %v @%v\n", oid, totalSz, avgPx)
								}
							} else if errorMsg, exists := statusMap["error"]; exists {
								fmt.Printf("Error: %v\n", errorMsg)
							}
						}
					}
				}
			}
		}

		fmt.Println("Waiting 2 seconds before closing position...")
		time.Sleep(2 * time.Second)

		fmt.Printf("Attempting to Market Close all %s position\n", coin)
		closeResult, err := exchange.MarketClose(ctx, coin, nil, nil, &slippage, nil, nil)
		if err != nil {
			log.Printf("Failed to close position: %v", err)
		} else {
			closeResultMap, ok := closeResult.(map[string]interface{})
			if !ok {
				log.Printf("Failed to parse close result")
				return
			}

			if status, ok := closeResultMap["status"].(string); ok && status == "ok" {
				if response, ok := closeResultMap["response"].(map[string]interface{}); ok {
					if data, ok := response["data"].(map[string]interface{}); ok {
						if statuses, ok := data["statuses"].([]interface{}); ok {
							for _, statusInterface := range statuses {
								if statusMap, ok := statusInterface.(map[string]interface{}); ok {
									if filled, exists := statusMap["filled"]; exists {
										if filledMap, ok := filled.(map[string]interface{}); ok {
											oid := filledMap["oid"]
											totalSz := filledMap["totalSz"]
											avgPx := filledMap["avgPx"]
											fmt.Printf("Order #%v filled %v @%v\n", oid, totalSz, avgPx)
										}
									} else if errorMsg, exists := statusMap["error"]; exists {
										fmt.Printf("Error: %v\n", errorMsg)
									}
								}
							}
						}
					}
				}
			} else {
				closeJSON, _ := json.MarshalIndent(closeResult, "", "  ")
				fmt.Printf("Close failed with result: %s\n", closeJSON)
			}
		}
	} else {
		orderJSON, _ := json.MarshalIndent(orderResult, "", "  ")
		fmt.Printf("Order failed with result: %s\n", orderJSON)
	}
}
