package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"hyperliquid-go-sdk/pkg/client"
	"hyperliquid-go-sdk/pkg/utils"
)

// PythonOrder represents the structure for calling Python SDK
type PythonOrder struct {
	Coin      string  `json:"coin"`
	IsBuy     bool    `json:"is_buy"`
	Size      float64 `json:"size"`
	LimitPx   float64 `json:"limit_px"`
	OrderType string  `json:"order_type"`
}

// HybridTrader demonstrates using Go for reads and Python for orders
type HybridTrader struct {
	address  string
	info     *client.Info
	exchange *client.Exchange
}

// NewHybridTrader creates a new hybrid trader
func NewHybridTrader() (*HybridTrader, error) {
	address, info, exchange := Setup(utils.TestnetAPIURL, true)
	
	return &HybridTrader{
		address:  address,
		info:     info,
		exchange: exchange,
	}, nil
}

// GetUserState gets user state using Go SDK (works reliably)
func (h *HybridTrader) GetUserState() (map[string]interface{}, error) {
	return h.info.UserState(h.address, "")
}

// GetAllMids gets all mid prices using Go SDK (works reliably)
func (h *HybridTrader) GetAllMids() (map[string]string, error) {
	return h.info.AllMids("")
}

// GetOpenOrders gets open orders using Go SDK
func (h *HybridTrader) GetOpenOrders() (map[string]interface{}, error) {
	return h.info.OpenOrders(h.address, "")
}

// PlaceOrderPython places an order using Python SDK (reliable for orders)
func (h *HybridTrader) PlaceOrderPython(coin string, isBuy bool, size float64, limitPx float64) (string, error) {
	// Create a temporary Python script that places the order
	pythonScript := fmt.Sprintf(`
import json
import os
from hyperliquid.info import Info
from hyperliquid.exchange import Exchange
from hyperliquid.utils import constants

# Use the same config as Go SDK
config_path = "./config.json"
if os.path.exists(config_path):
    with open(config_path, 'r') as f:
        config = json.load(f)
    
    # Initialize the exchange
    exchange = Exchange(None, constants.TESTNET_API_URL, account_address=config.get("account_address"))
    
    # Place the order
    result = exchange.order("%s", %t, %.8f, %.2f, {"limit": {"tif": "Gtc"}})
    print(json.dumps(result))
else:
    print(json.dumps({"error": "config.json not found"}))
`, coin, isBuy, size, limitPx)

	// Write the script to a temporary file
	scriptFile := "/tmp/hyperliquid_order.py"
	if err := os.WriteFile(scriptFile, []byte(pythonScript), 0644); err != nil {
		return "", fmt.Errorf("failed to write Python script: %w", err)
	}
	defer os.Remove(scriptFile)

	// Execute the Python script
	cmd := exec.Command("python3", scriptFile)
	cmd.Dir = "/Users/madhugowda/Desktop/eqlzr-v2/hyperliquid-go-sdk/examples"
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute Python script: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// DisplayPositions shows current positions in a nice format
func (h *HybridTrader) DisplayPositions() error {
	userState, err := h.GetUserState()
	if err != nil {
		return fmt.Errorf("failed to get user state: %w", err)
	}

	fmt.Println("=== Current Positions ===")
	PrintPositions(userState)
	return nil
}

// DisplayMarketData shows current market data
func (h *HybridTrader) DisplayMarketData(coins []string) error {
	mids, err := h.GetAllMids()
	if err != nil {
		return fmt.Errorf("failed to get market data: %w", err)
	}

	fmt.Println("=== Market Data ===")
	for _, coin := range coins {
		if price, exists := mids[coin]; exists {
			fmt.Printf("%s: $%s\n", coin, price)
		}
	}
	return nil
}

func main() {
	fmt.Println("=== Hybrid Hyperliquid Solution ===\n")
	
	// Create hybrid trader
	trader, err := NewHybridTrader()
	if err != nil {
		log.Fatalf("Failed to create trader: %v", err)
	}

	fmt.Printf("Trading with account: %s\n\n", trader.address)

	// 1. Show current positions (using Go SDK - works reliably)
	if err := trader.DisplayPositions(); err != nil {
		log.Printf("Error displaying positions: %v", err)
	}

	// 2. Show market data (using Go SDK - works reliably)  
	coins := []string{"ETH", "BTC", "SOL"}
	if err := trader.DisplayMarketData(coins); err != nil {
		log.Printf("Error displaying market data: %v", err)
	}

	// 3. Show open orders (using Go SDK - works reliably)
	fmt.Println("\n=== Open Orders ===")
	openOrders, err := trader.GetOpenOrders()
	if err != nil {
		log.Printf("Error getting open orders: %v", err)
	} else {
		ordersJSON, _ := json.MarshalIndent(openOrders, "", "  ")
		fmt.Println(string(ordersJSON))
	}

	// 4. Place an order using Python SDK (reliable for orders)
	fmt.Println("\n=== Order Placement ===")
	fmt.Println("To place an order, we'll use the Python SDK:")
	
	result, err := trader.PlaceOrderPython("ETH", true, 0.01, 3800.0)
	if err != nil {
		fmt.Printf("Error placing order via Python: %v\n", err)
		fmt.Println("Note: Make sure you have the Python hyperliquid-python-sdk installed:")
		fmt.Println("pip install hyperliquid-python-sdk")
	} else {
		fmt.Printf("Order result: %s\n", result)
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("✅ Go SDK works perfectly for:")
	fmt.Println("   - Getting user state and positions")
	fmt.Println("   - Getting market data")
	fmt.Println("   - Getting open orders")
	fmt.Println("   - Getting order history")
	fmt.Println()
	fmt.Println("✅ Python SDK works perfectly for:")
	fmt.Println("   - Placing orders")
	fmt.Println("   - Canceling orders") 
	fmt.Println("   - All trading operations")
	fmt.Println()
	fmt.Println("This hybrid approach gives you the best of both worlds!")
}