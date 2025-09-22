package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"hyperliquid-go-sdk/pkg/client"
	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
)

// Config represents the configuration structure
type Config struct {
	SecretKey      string `json:"secret_key"`
	AccountAddress string `json:"account_address"`
}

// loadConfig loads configuration from config.json file
func loadConfig() *Config {
	configPath := "./config.json"

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("config.json not found. Please set environment variables or create config.json")
		return &Config{}
	}

	file, err := os.Open(configPath)
	if err != nil {
		log.Printf("Error opening config file: %v", err)
		return &Config{}
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Printf("Error decoding config file: %v", err)
		return &Config{}
	}

	return &config
}

func main() {
	/*
		Go equivalent of Python basic_agent.py
		
		This demonstrates:
		1. Setting up an environment for testing purposes by creating an agent
		2. The agent can place trades on behalf of the account
		3. The agent does not have permission to transfer or withdraw funds
		4. Shows how to create temporary and persistent agents
		5. Demonstrates placing and canceling orders with agents
	*/

	// Set up the environment (exchange, account info, etc.) for testing purposes.
	fmt.Println("Setting up environment...")
	
	// Try to get private key from environment variable first
	privateKeyHex := os.Getenv("HYPERLIQUID_PRIVATE_KEY")
	address := os.Getenv("HYPERLIQUID_ADDRESS")
	
	// If not found in environment, try to read from config file
	if privateKeyHex == "" {
		config := loadConfig()
		privateKeyHex = config.SecretKey
		if address == "" {
			address = config.AccountAddress
		}
	}

	if privateKeyHex == "" {
		log.Fatal("Private key not found. Set HYPERLIQUID_PRIVATE_KEY environment variable or update config.json")
	}

	// Parse private key
	privateKey, err := utils.ParsePrivateKey(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Get address from private key if not provided
	walletAddress := utils.GetAddressFromPrivateKey(privateKey)
	if address == "" {
		address = walletAddress
	}

	fmt.Printf("Main account: %s\n", address)
	if address != walletAddress {
		fmt.Printf("Wallet: %s\n", walletAddress)
	}

	// Ensure that the wallet address and account address are the same.
	// If these are not the same then we need to use the wallet address as the main account.
	if address != walletAddress {
		fmt.Printf("Account address differs from wallet address. Using wallet address as main account.\n")
		address = walletAddress
	}

	// Create clients
	timeout := 30 * time.Second
	
	// Create info client (not used in this test but could be useful)
	_, err = client.NewInfo(utils.TestnetAPIURL, &timeout, true, nil, nil, nil)
	if err != nil {
		log.Fatalf("Failed to create info client: %v", err)
	}

	// Create exchange client
	exchange, err := client.NewExchange(
		privateKey,
		utils.TestnetAPIURL,
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

	fmt.Println("Exchange and info clients created successfully")

	// Step 1: Approve an agent
	fmt.Println("\n=== Approving Agent ===")
	approveResult, err := exchange.ApproveAgent()
	if err != nil {
		log.Fatalf("Failed to approve agent: %v", err)
	}

	// Check if the agent approval was successful
	fmt.Printf("Agent approval result: %+v\n", approveResult.Result)
	
	status, ok := approveResult.Result["status"].(string)
	if !ok || status != "ok" {
		log.Fatalf("Agent approval failed: %+v", approveResult.Result)
	}

	fmt.Printf("Agent approved successfully!\n")
	fmt.Printf("Agent private key: %s\n", approveResult.AgentKey)

	// Step 2: Create the agent's exchange client
	fmt.Println("\n=== Creating Agent Exchange Client ===")
	
	// Parse the agent's private key
	agentPrivateKey, err := utils.ParsePrivateKey(approveResult.AgentKey)
	if err != nil {
		log.Fatalf("Failed to parse agent private key: %v", err)
	}

	// Get agent address
	agentAddress := utils.GetAddressFromPrivateKey(agentPrivateKey)
	fmt.Printf("Agent address: %s\n", agentAddress)

	// Create agent exchange client - use agent's private key but main account's address
	agentExchange, err := client.NewExchange(
		agentPrivateKey,
		utils.TestnetAPIURL,
		&timeout,
		nil,      // meta
		nil,      // vault address
		&address, // account address (main account, not agent address)
		nil,      // spot meta
		nil,      // perp dexs
	)
	if err != nil {
		log.Fatalf("Failed to create agent exchange client: %v", err)
	}

	fmt.Println("Agent exchange client created successfully!")

	// Step 3: Place a test order with the agent
	fmt.Println("\n=== Placing Order with Agent ===")
	
	// Create order type
	orderType := types.OrderType{
		Limit: &types.LimitOrderType{
			Tif: types.TifGtc,
		},
	}

	// Place a test order with very low price so it rests in the order book
	orderResult, err := agentExchange.Order(
		"ETH",     // coin
		true,      // isBuy
		0.2,       // size
		1000.0,    // limit price (very low to ensure it rests)
		orderType, // order type
		false,     // reduce only
		nil,       // cloid
		nil,       // builder info
	)
	
	if err != nil {
		log.Printf("Failed to place agent order: %v", err)
	} else {
		fmt.Printf("Agent order result: %+v\n", orderResult)
		
		// If the order was placed successfully and is resting, cancel it
		if status, ok := orderResult["status"].(string); ok && status == "ok" {
			if response, ok := orderResult["response"].(map[string]interface{}); ok {
				if data, ok := response["data"].(map[string]interface{}); ok {
					if statuses, ok := data["statuses"].([]interface{}); ok && len(statuses) > 0 {
						if statusMap, ok := statuses[0].(map[string]interface{}); ok {
							if resting, ok := statusMap["resting"].(map[string]interface{}); ok {
								if oidFloat, ok := resting["oid"].(float64); ok {
									oid := int(oidFloat)
									fmt.Printf("Order placed and resting with oid: %d\n", oid)
									
									// Cancel the order
									fmt.Println("Canceling the order...")
									cancelResult, err := agentExchange.Cancel("ETH", oid)
									if err != nil {
										log.Printf("Failed to cancel order: %v", err)
									} else {
										fmt.Printf("Cancel result: %+v\n", cancelResult)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Step 4: Create a persistent agent
	fmt.Println("\n=== Creating Persistent Agent ===")
	
	persistentResult, err := exchange.ApproveAgent("persistent_agent")
	if err != nil {
		log.Printf("Failed to approve persistent agent: %v", err)
	} else {
		fmt.Printf("Persistent agent approval result: %+v\n", persistentResult.Result)
		
		if status, ok := persistentResult.Result["status"].(string); ok && status == "ok" {
			fmt.Printf("Persistent agent created successfully!\n")
			fmt.Printf("Persistent agent private key: %s\n", persistentResult.AgentKey)
			
			// Parse persistent agent key
			persistentAgentKey, err := utils.ParsePrivateKey(persistentResult.AgentKey)
			if err != nil {
				log.Printf("Failed to parse persistent agent key: %v", err)
			} else {
				persistentAgentAddress := utils.GetAddressFromPrivateKey(persistentAgentKey)
				fmt.Printf("Persistent agent address: %s\n", persistentAgentAddress)

				// Create persistent agent exchange client
				persistentExchange, err := client.NewExchange(
					persistentAgentKey,
					utils.TestnetAPIURL,
					&timeout,
					nil,      // meta
					nil,      // vault address
					&address, // account address
					nil,      // spot meta
					nil,      // perp dexs
				)
				if err != nil {
					log.Printf("Failed to create persistent agent exchange client: %v", err)
				} else {
					fmt.Println("Persistent agent exchange client created!")
					
					// Place order with persistent agent
					fmt.Println("Placing order with persistent agent...")
					persistentOrderResult, err := persistentExchange.Order(
						"ETH",     // coin
						true,      // isBuy
						0.2,       // size
						1000.0,    // limit price
						orderType, // order type
						false,     // reduce only
						nil,       // cloid
						nil,       // builder info
					)
					
					if err != nil {
						log.Printf("Failed to place persistent agent order: %v", err)
					} else {
						fmt.Printf("Persistent agent order result: %+v\n", persistentOrderResult)
						
						// Cancel if resting
						if status, ok := persistentOrderResult["status"].(string); ok && status == "ok" {
							if response, ok := persistentOrderResult["response"].(map[string]interface{}); ok {
								if data, ok := response["data"].(map[string]interface{}); ok {
									if statuses, ok := data["statuses"].([]interface{}); ok && len(statuses) > 0 {
										if statusMap, ok := statuses[0].(map[string]interface{}); ok {
											if resting, ok := statusMap["resting"].(map[string]interface{}); ok {
												if oidFloat, ok := resting["oid"].(float64); ok {
													oid := int(oidFloat)
													fmt.Printf("Persistent agent order resting with oid: %d\n", oid)
													
													// Cancel the order
													fmt.Println("Canceling persistent agent order...")
													cancelResult, err := persistentExchange.Cancel("ETH", oid)
													if err != nil {
														log.Printf("Failed to cancel persistent agent order: %v", err)
													} else {
														fmt.Printf("Persistent agent cancel result: %+v\n", cancelResult)
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	fmt.Println("\n=== Agent Workflow Test Completed ===")
	fmt.Println("This test demonstrated:")
	fmt.Println("1. Creating and approving agents for trading")
	fmt.Println("2. Using agent keys to create exchange clients")
	fmt.Println("3. Placing orders with agents on behalf of main account")
	fmt.Println("4. Creating both temporary and persistent agents")
	fmt.Println("5. Managing orders placed by agents")
}