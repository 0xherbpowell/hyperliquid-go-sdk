package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/client"
	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
	"time"
)

func main() {
	fmt.Println("=== Fixed Nonce Signature Test ===")
	
	// Use the exact same setup as your examples
	privateKeyHex := "06e10c1cb33b369c878ec8f3d51523f2bdd3a36f02fcb6c29e0867903e17927e"
	privateKey, err := utils.ParsePrivateKey(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	
	// Create info client to get asset ID
	timeout := 30 * time.Second
	info, err := client.NewInfo(utils.TestnetAPIURL, &timeout, true, nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	
	// Get ETH asset ID
	asset, err := info.NameToAsset("ETH")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("ETH Asset ID: %d\n", asset)
	
	// Create the order request exactly like in examples
	orderRequest := types.OrderRequest{
		Coin:       "ETH",
		IsBuy:      true,
		Sz:         0.2,
		LimitPx:    1100.0,
		OrderType:  createGtcLimitOrder(),
		ReduceOnly: false,
		Cloid:      generateCloid(),
	}
	
	// Convert to wire format
	orderWire, err := utils.OrderRequestToOrderWire(orderRequest, asset)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("OrderWire: %+v\n", orderWire)
	
	// Create action
	orderAction := utils.OrderWiresToOrderAction([]types.OrderWire{orderWire}, nil)
	
	fmt.Printf("OrderAction: %+v\n", orderAction)
	
	// Use a fixed nonce for comparison
	fixedNonce := int64(1234567890)
	
	// Sign the action
	signature, err := utils.SignL1Action(privateKey, orderAction, nil, fixedNonce, nil, false)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Fixed nonce signature:\n")
	fmt.Printf("r: %s\n", signature["r"])
	fmt.Printf("s: %s\n", signature["s"])
	fmt.Printf("v: %d\n", signature["v"])
	
	// Now let's try to verify this matches our expected address
	expectedAddr := utils.GetAddressFromPrivateKey(privateKey)
	fmt.Printf("Expected address: %s\n", expectedAddr)
}

func createGtcLimitOrder() types.OrderType {
	return types.OrderType{
		Limit: &types.LimitOrderType{
			Tif: types.TifGtc,
		},
	}
}

func generateCloid() *types.Cloid {
	return types.NewCloidFromInt(time.Now().Unix())
}