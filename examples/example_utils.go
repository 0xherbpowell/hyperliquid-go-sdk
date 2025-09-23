package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"hyperliquid-go-sdk/pkg/client"
	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
)

// Config represents the configuration structure
// Note: examples should not prompt for env files or private keys; this is only for
// optional local convenience when present.
type Config struct {
	SecretKey      string `json:"secret_key"`
	AccountAddress string `json:"account_address"`
}

// Setup initializes the exchange and info clients for examples
func Setup(baseURL string, skipWS bool) (string, *client.Info, *client.Exchange) {
	// Read optional environment values without logging them
	privateKeyHex := os.Getenv("HYPERLIQUID_PRIVATE_KEY")
	address := os.Getenv("HYPERLIQUID_ADDRESS")

	// Optional: fall back to config.json if present, but do not require it
	if privateKeyHex == "" {
		config := loadConfig()
		privateKeyHex = config.SecretKey
		if address == "" {
			address = config.AccountAddress
		}
	}

	if privateKeyHex == "" {
		log.Fatal("No signing key configured; cannot place orders. Configure a signer in your environment or code.")
	}

	// Parse private key
	privateKey, err := utils.ParsePrivateKey(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	debugSignature(privateKey, address)
	// Get address from private key - this is who will be signing
	walletAddress := utils.GetAddressFromPrivateKey(privateKey)

	// Handle agent vs direct wallet scenarios
	if address == "" {
		// No account address specified, use the wallet address directly
		address = walletAddress
		fmt.Printf("Direct wallet mode: %s\n", address)
	} else if !strings.EqualFold(address, walletAddress) {
		// Agent mode: wallet signs for the account
		fmt.Printf("Agent mode: Account %s, Agent wallet %s\n", address, walletAddress)
		// Keep the original account address - the SDK will handle agent signing
	} else {
		// Addresses match
		fmt.Printf("Direct wallet mode: %s\n", address)
	}

	fmt.Printf("Running with account address: %s\n", address)

	// Create info client
	timeout := 30 * time.Second
	info, err := client.NewInfo(baseURL, &timeout, skipWS, nil, nil, nil)
	if err != nil {
		log.Fatalf("Failed to create info client: %v", err)
	}

	// Check if account has equity
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	// Check margin summary
	if marginSummary, ok := userState["marginSummary"].(map[string]interface{}); ok {
		if accountValue, ok := marginSummary["accountValue"].(string); ok {
			if accountValue == "0" {
				log.Fatal("Not running the example because the provided account has no equity.")
			}
		}
	}

	// Create exchange client
	exchange, err := client.NewExchange(
		privateKey,
		baseURL,
		&timeout,
		nil,      // meta
		nil,      // vault address
		&address, // account address
		nil,      // spot meta
		nil,      // perp dexs
	)
	if err != nil {
		log.Fatalf("Failed to create exchange client: %v", err)
	}

	return address, info, exchange
}

func debugSignature(privateKey *ecdsa.PrivateKey, accountAddress string) {
	// Verify the address derivation
	derivedAddress := utils.GetAddressFromPrivateKey(privateKey)
	fmt.Printf("Private key derives to: %s\n", derivedAddress)
	fmt.Printf("Account address: %s\n", accountAddress)
	fmt.Printf("Addresses match: %t\n", strings.EqualFold(derivedAddress, accountAddress))

	// Test a simple signature to verify the private key works
	testMessage := []byte("test message")
	testHash := crypto.Keccak256Hash(testMessage)
	signature, err := crypto.Sign(testHash.Bytes(), privateKey)
	if err != nil {
		fmt.Printf("Failed to create test signature: %v\n", err)
		return
	}

	// Recover the public key from signature
	recoveredPubKey, err := crypto.SigToPub(testHash.Bytes(), signature)
	if err != nil {
		fmt.Printf("Failed to recover public key: %v\n", err)
		return
	}

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)
	fmt.Printf("Recovered address from test signature: %s\n", recoveredAddress.Hex())
	fmt.Printf("Recovery matches derived: %t\n", strings.EqualFold(recoveredAddress.Hex(), derivedAddress))
}

// loadConfig loads configuration from config.json file
func loadConfig() *Config {
	configPath := "./config.json"

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}
	}

	file, err := os.Open(configPath)
	if err != nil {
		return &Config{}
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return &Config{}
	}

	return &config
}

// PrintPositions prints user positions in a formatted way
func PrintPositions(userState map[string]interface{}) {
	if assetPositions, ok := userState["assetPositions"].([]interface{}); ok {
		positions := []map[string]interface{}{}

		for _, ap := range assetPositions {
			if apMap, ok := ap.(map[string]interface{}); ok {
				if position, ok := apMap["position"].(map[string]interface{}); ok {
					positions = append(positions, position)
				}
			}
		}

		if len(positions) > 0 {
			fmt.Println("positions:")
			for _, position := range positions {
				jsonData, _ := json.MarshalIndent(position, "", "  ")
				fmt.Println(string(jsonData))
			}
		} else {
			fmt.Println("no open positions")
		}
	}
}

// CreateGtcLimitOrder creates a Good Till Cancel limit order type
func CreateGtcLimitOrder() types.OrderType {
	return types.OrderType{
		Limit: &types.LimitOrderType{
			Tif: types.TifGtc,
		},
	}
}

// CreateIocLimitOrder creates an Immediate or Cancel limit order type
func CreateIocLimitOrder() types.OrderType {
	return types.OrderType{
		Limit: &types.LimitOrderType{
			Tif: types.TifIoc,
		},
	}
}

// CreateTpslOrder creates a take profit or stop loss order type
func CreateTpslOrder(triggerPx float64, isMarket bool, tpsl types.Tpsl) types.OrderType {
	return types.OrderType{
		Trigger: &types.TriggerOrderType{
			TriggerPx: triggerPx,
			IsMarket:  isMarket,
			Tpsl:      tpsl,
		},
	}
}

// GenerateCloid generates a unique client order ID
func GenerateCloid() *types.Cloid {
	return types.NewCloidFromInt(time.Now().Unix())
}

// PrintOrderResult prints the order result in a formatted way
func PrintOrderResult(result map[string]interface{}) {
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
}

// ParsePrice parses a price string to float64
func ParsePrice(priceStr string) (float64, error) {
	return strconv.ParseFloat(priceStr, 64)
}

// GetRestingOid extracts the resting order ID from order result
func GetRestingOid(orderResult map[string]interface{}) (int, bool) {
	if status, ok := orderResult["status"].(string); ok && status == "ok" {
		if response, ok := orderResult["response"].(map[string]interface{}); ok {
			if data, ok := response["data"].(map[string]interface{}); ok {
				if statuses, ok := data["statuses"].([]interface{}); ok && len(statuses) > 0 {
					if statusMap, ok := statuses[0].(map[string]interface{}); ok {
						if resting, ok := statusMap["resting"].(map[string]interface{}); ok {
							if oidFloat, ok := resting["oid"].(float64); ok {
								return int(oidFloat), true
							}
						}
					}
				}
			}
		}
	}
	return 0, false
}
