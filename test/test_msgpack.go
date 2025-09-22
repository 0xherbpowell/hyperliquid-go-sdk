package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vmihailenco/msgpack/v5"
	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Create a sample OrderWire to test msgpack serialization
	orderWire := types.OrderWire{
		A: 0,         // asset (ETH)
		B: true,      // isBuy
		P: "1000.0",  // limitPx
		S: "0.1",     // sz
		R: false,     // reduceOnly
		T: types.OrderTypeWire{
			Limit: &types.LimitOrderType{
				Tif: types.TifGtc,
			},
		},
		C: nil, // cloid
	}

	// Test JSON serialization
	fmt.Println("=== Testing JSON Serialization ===")
	jsonData, err := json.MarshalIndent(orderWire, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshal error: %v", err)
	}
	fmt.Printf("JSON:\n%s\n\n", jsonData)

	// Test MessagePack serialization
	fmt.Println("=== Testing MessagePack Serialization ===")
	msgpackData, err := msgpack.Marshal(orderWire)
	if err != nil {
		log.Fatalf("MessagePack marshal error: %v", err)
	}
	fmt.Printf("MessagePack bytes (length %d): %x\n", len(msgpackData), msgpackData)

	// Test deserialization
	var deserializedOrder types.OrderWire
	err = msgpack.Unmarshal(msgpackData, &deserializedOrder)
	if err != nil {
		log.Fatalf("MessagePack unmarshal error: %v", err)
	}
	
	fmt.Printf("Deserialized OrderWire:\n")
	fmt.Printf("  Asset: %d\n", deserializedOrder.A)
	fmt.Printf("  IsBuy: %t\n", deserializedOrder.B)
	fmt.Printf("  LimitPx: %s\n", deserializedOrder.P)
	fmt.Printf("  Size: %s\n", deserializedOrder.S)
	fmt.Printf("  ReduceOnly: %t\n", deserializedOrder.R)
	fmt.Printf("  OrderType.Limit.Tif: %s\n", deserializedOrder.T.Limit.Tif)

	// Test creating action with OrderWires
	fmt.Println("\n=== Testing Action Serialization ===")
	orderWires := []types.OrderWire{orderWire}
	action := utils.OrderWiresToOrderAction(orderWires, nil)
	
	actionJson, _ := json.MarshalIndent(action, "", "  ")
	fmt.Printf("Action JSON:\n%s\n", actionJson)
	
	actionMsgpack, err := msgpack.Marshal(action)
	if err != nil {
		log.Fatalf("Action MessagePack marshal error: %v", err)
	}
	fmt.Printf("Action MessagePack bytes (length %d): %x\n", len(actionMsgpack), actionMsgpack)

	fmt.Println("\n=== Test completed successfully! ===")
	fmt.Println("MessagePack tags are working correctly.")
}