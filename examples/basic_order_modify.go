package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// Get the user state and print out position information
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	PrintPositions(userState)

	// Place an initial order that should rest by setting the price very low
	orderResult, err := exchange.Order(
		"ETH",                // coin
		true,                 // isBuy
		0.2,                  // size
		1100.0,               // limit price (low price to ensure it rests)
		CreateGtcLimitOrder(), // order type
		false,                // reduce only
		GenerateCloid(),      // client order ID
		nil,                  // builder info
	)
	if err != nil {
		log.Printf("Failed to place order: %v", err)
		return
	}

	fmt.Println("Initial order result:")
	PrintOrderResult(orderResult)

	// Get the resting order ID
	oid, ok := GetRestingOid(orderResult)
	if !ok {
		fmt.Println("Order was not resting, cannot modify")
		return
	}

	fmt.Printf("Order ID to modify: %d\n", oid)

	// Create a modified order with different price and size
	modifiedOrder := types.OrderRequest{
		Coin:       "ETH",
		IsBuy:      true,
		Sz:         0.3,  // increased size
		LimitPx:    1200.0, // higher price
		OrderType:  CreateGtcLimitOrder(),
		ReduceOnly: false,
		Cloid:      GenerateCloid(), // new client order ID for the modified order
	}

	// Modify the order
	modifyResult, err := exchange.Modify(oid, modifiedOrder)
	if err != nil {
		log.Printf("Failed to modify order: %v", err)
		return
	}

	fmt.Println("Modify order result:")
	PrintOrderResult(modifyResult)

	// Get the new order ID from the modify result
	newOid, ok := GetRestingOid(modifyResult)
	if ok {
		fmt.Printf("New order ID after modification: %d\n", newOid)

		// Query the modified order status
		orderStatus, err := info.OrderStatus(address, newOid, "")
		if err != nil {
			log.Printf("Failed to get modified order status: %v", err)
		} else {
			fmt.Println("Modified order status:")
			PrintOrderResult(orderStatus)
		}

		// Cancel the modified order
		cancelResult, err := exchange.Cancel("ETH", newOid)
		if err != nil {
			log.Printf("Failed to cancel modified order: %v", err)
		} else {
			fmt.Println("Cancel modified order result:")
			PrintOrderResult(cancelResult)
		}
	} else {
		fmt.Println("Modified order may have been filled immediately")
	}
}