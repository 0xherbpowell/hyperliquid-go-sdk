package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

// Example script to deploy HIP-1 and HIP-2 assets
// See https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/deploying-hip-1-and-hip-2-assets
// for the spec.
//
// IMPORTANT: Replace any arguments for the exchange calls below to match your deployment requirements.

const (
	// Set to true to enable freeze functionality for the deployed token
	// See step 2-a below for more details on freezing.
	ENABLE_FREEZE_PRIVILEGE = false
	
	// Set to true to set the deployer trading fee share
	// See step 6 below for more details on setting the deployer trading fee share.
	SET_DEPLOYER_TRADING_FEE_SHARE = false
	
	// See step 7 below for more details on enabling quote token.
	ENABLE_QUOTE_TOKEN = false
	
	DUMMY_USER = "0x0000000000000000000000000000000000000001"
)

func main() {
	// Setup with testnet
	address, _, exchange := Setup(utils.TestnetAPIURL, true)

	fmt.Printf("Starting spot deployment process for address: %s\n", address)

	// Step 1: Registering the Token
	//
	// Takes part in the spot deploy auction and if successful, registers token "TEST0"
	// with sz_decimals 2 and wei_decimals 8.
	// The max gas is 10,000 HYPE and represents the max amount to be paid for the spot deploy auction.
	fmt.Println("\n=== Step 1: Registering the Token ===")
	
	registerTokenResult, err := exchange.SpotDeployRegisterToken("TEST0", 2, 8, "1000000000000", "Test token example")
	if err != nil {
		log.Printf("Failed to register token: %v", err)
		return
	}
	
	fmt.Println("Register token result:")
	PrintOrderResult(registerTokenResult)

	// If registration is successful, a token index will be returned. This token index is required for
	// later steps in the spot deploy process.
	var token interface{}
	if status, ok := registerTokenResult["status"].(string); ok && status == "ok" {
		if response, ok := registerTokenResult["response"].(map[string]interface{}); ok {
			if data, ok := response["data"]; ok {
				token = data
				fmt.Printf("Token registered successfully with index: %v\n", token)
			}
		}
	} else {
		fmt.Println("Token registration failed")
		return
	}

	// Step 2: User Genesis
	//
	// User genesis can be called multiple times to associate balances to specific users and/or
	// tokens for genesis.
	fmt.Println("\n=== Step 2: User Genesis ===")
	
	// Associate 100000000000000 wei with user and hyperliquidity
	userGenesisBalances := [][]string{
		{DUMMY_USER, "100000000000000"},
		{"0xffffffffffffffffffffffffffffffffffffffff", "100000000000000"},
	}
	tokenDistribution := [][]interface{}{} // Empty for this call

	userGenesisResult, err := exchange.SpotDeployUserGenesis(token, userGenesisBalances, tokenDistribution)
	if err != nil {
		log.Printf("Failed user genesis: %v", err)
	} else {
		fmt.Println("User genesis result:")
		PrintOrderResult(userGenesisResult)
	}

	// No-op example
	userGenesisResult2, err := exchange.SpotDeployUserGenesis(token, [][]string{}, [][]interface{}{})
	if err != nil {
		log.Printf("Failed user genesis (no-op): %v", err)
	} else {
		fmt.Println("User genesis (no-op) result:")
		PrintOrderResult(userGenesisResult2)
	}

	// Distribute 100000000000000 wei on a weighted basis to all holders of token with index 1
	tokenDistribution = [][]interface{}{{1, "100000000000000"}}
	userGenesisResult3, err := exchange.SpotDeployUserGenesis(token, [][]string{}, tokenDistribution)
	if err != nil {
		log.Printf("Failed user genesis (token distribution): %v", err)
	} else {
		fmt.Println("User genesis (token distribution) result:")
		PrintOrderResult(userGenesisResult3)
	}

	if ENABLE_FREEZE_PRIVILEGE {
		// Step 2-a: Enables the deployer to freeze/unfreeze users. Freezing a user means
		// that user cannot trade, send, or receive this token.
		fmt.Println("\n=== Step 2-a: Enable Freeze Privilege ===")
		
		enableFreezeResult, err := exchange.SpotDeployEnableFreezePrivilege(token)
		if err != nil {
			log.Printf("Failed to enable freeze privilege: %v", err)
		} else {
			fmt.Println("Enable freeze privilege result:")
			PrintOrderResult(enableFreezeResult)
		}

		// Freeze user for token
		freezeResult, err := exchange.SpotDeployFreezeUser(token, DUMMY_USER, true)
		if err != nil {
			log.Printf("Failed to freeze user: %v", err)
		} else {
			fmt.Println("Freeze user result:")
			PrintOrderResult(freezeResult)
		}

		// Unfreeze user for token
		unfreezeResult, err := exchange.SpotDeployFreezeUser(token, DUMMY_USER, false)
		if err != nil {
			log.Printf("Failed to unfreeze user: %v", err)
		} else {
			fmt.Println("Unfreeze user result:")
			PrintOrderResult(unfreezeResult)
		}
	}

	// Step 3: Genesis
	//
	// Finalize genesis. The max supply of 300000000000000 wei needs to match the total
	// allocation above from user genesis.
	//
	// "noHyperliquidity" can also be set to disable hyperliquidity. In that case, no balance
	// should be associated with hyperliquidity from step 2 (user genesis).
	fmt.Println("\n=== Step 3: Genesis ===")
	
	genesisResult, err := exchange.SpotDeployGenesis(token, "300000000000000", false)
	if err != nil {
		log.Printf("Failed genesis: %v", err)
		return
	}
	
	fmt.Println("Genesis result:")
	PrintOrderResult(genesisResult)

	// Step 4: Register Spot
	//
	// Register the spot pair (TEST0/USDC) given base and quote token indices. 0 represents USDC.
	// The base token is the first token in the pair and the quote token is the second token.
	fmt.Println("\n=== Step 4: Register Spot ===")
	
	registerSpotResult, err := exchange.SpotDeployRegisterSpot(token, 0)
	if err != nil {
		log.Printf("Failed to register spot: %v", err)
		return
	}
	
	fmt.Println("Register spot result:")
	PrintOrderResult(registerSpotResult)

	// If registration is successful, a spot index will be returned. This spot index is required for
	// registering hyperliquidity.
	var spot interface{}
	if status, ok := registerSpotResult["status"].(string); ok && status == "ok" {
		if response, ok := registerSpotResult["response"].(map[string]interface{}); ok {
			if data, ok := response["data"]; ok {
				spot = data
				fmt.Printf("Spot registered successfully with index: %v\n", spot)
			}
		}
	} else {
		fmt.Println("Spot registration failed")
		return
	}

	// Step 5: Register Hyperliquidity
	//
	// Registers hyperliquidity for the spot pair. In this example, hyperliquidity is registered
	// with a starting price of $2, an order size of 4, and 100 total orders.
	//
	// This step is required even if "noHyperliquidity" was set to True.
	// If "noHyperliquidity" was set to True during step 3 (genesis), then "n_orders" is required to be 0.
	fmt.Println("\n=== Step 5: Register Hyperliquidity ===")
	
	registerHyperliquidityResult, err := exchange.SpotDeployRegisterHyperliquidity(spot, 2.0, 4.0, 100, nil)
	if err != nil {
		log.Printf("Failed to register hyperliquidity: %v", err)
	} else {
		fmt.Println("Register hyperliquidity result:")
		PrintOrderResult(registerHyperliquidityResult)
	}

	if SET_DEPLOYER_TRADING_FEE_SHARE {
		// Step 6
		//
		// Note that the deployer trading fee share cannot increase.
		// The default is already 100% and the smallest increment is 0.001%.
		fmt.Println("\n=== Step 6: Set Deployer Trading Fee Share ===")
		
		setFeeShareResult, err := exchange.SpotDeploySetDeployerTradingFeeShare(token, "100%")
		if err != nil {
			log.Printf("Failed to set deployer trading fee share: %v", err)
		} else {
			fmt.Println("Set deployer trading fee share result:")
			PrintOrderResult(setFeeShareResult)
		}
	}

	if ENABLE_QUOTE_TOKEN {
		// Step 7
		//
		// Note that deployer trading fee share must be zero.
		// The quote token must also be allowed.
		fmt.Println("\n=== Step 7: Enable Quote Token ===")
		
		enableQuoteResult, err := exchange.SpotDeployEnableQuoteToken(token)
		if err != nil {
			log.Printf("Failed to enable quote token: %v", err)
		} else {
			fmt.Println("Enable quote token result:")
			PrintOrderResult(enableQuoteResult)
		}
	}

	fmt.Println("\n=== Spot Deployment Complete ===")
	fmt.Printf("Token: %v\n", token)
	fmt.Printf("Spot: %v\n", spot)
	fmt.Println("\nIMPORTANT: This is a testnet deployment.")
	fmt.Println("For mainnet, ensure all parameters match your requirements.")
	fmt.Println("Spot deployment involves real costs and should be done carefully.")
}