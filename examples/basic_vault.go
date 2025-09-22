package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// This example demonstrates vault operations on Hyperliquid
	// Vaults allow users to deposit funds with professional traders/strategies

	fmt.Println("Basic Vault Example")
	fmt.Printf("Account: %s\n", address)

	// Get initial user state
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("\nInitial account state:")
	PrintPositions(userState)

	// Example 1: List available vaults
	fmt.Println("\n--- Example 1: Available Vaults ---")

	fmt.Println("Querying available vaults...")
	fmt.Printf("Would call: info.GetVaults()\n")
	
	// Simulate vault list
	fmt.Println("Available vaults (simulated):")
	vaults := []struct {
		address     string
		name        string
		description string
		tvl         string
		apy         string
		fees        string
		status      string
	}{
		{
			"0xVault1...", 
			"Conservative Trading Vault", 
			"Low-risk trading strategies focusing on market making",
			"$2,500,000",
			"15.2%",
			"20% performance + 2% management",
			"active",
		},
		{
			"0xVault2...", 
			"High-Frequency Arbitrage", 
			"Automated arbitrage strategies across multiple exchanges",
			"$1,800,000",
			"28.7%",
			"25% performance + 3% management",
			"active",
		},
		{
			"0xVault3...", 
			"DeFi Yield Optimizer", 
			"Optimized yield farming across DeFi protocols",
			"$950,000",
			"12.5%",
			"15% performance + 1.5% management",
			"active",
		},
	}
	
	for i, vault := range vaults {
		fmt.Printf("%d. %s\n", i+1, vault.name)
		fmt.Printf("   Address: %s\n", vault.address)
		fmt.Printf("   Description: %s\n", vault.description)
		fmt.Printf("   TVL: %s\n", vault.tvl)
		fmt.Printf("   APY: %s\n", vault.apy)
		fmt.Printf("   Fees: %s\n", vault.fees)
		fmt.Printf("   Status: %s\n", vault.status)
		fmt.Println()
	}

	// Example 2: Get vault details
	fmt.Println("\n--- Example 2: Vault Details ---")
	
	selectedVaultAddress := "0xVault1..."
	fmt.Printf("Getting details for vault: %s\n", selectedVaultAddress)
	fmt.Printf("Would call: info.GetVaultDetails(\"%s\")\n", selectedVaultAddress)
	
	// Simulate vault details
	fmt.Printf("Vault details for %s (simulated):\n", selectedVaultAddress)
	fmt.Println("{")
	fmt.Println("  \"address\": \"0xVault1...\",")
	fmt.Println("  \"name\": \"Conservative Trading Vault\",")
	fmt.Println("  \"manager\": \"0xManager1...\",")
	fmt.Println("  \"tvl\": \"2500000.0\",")
	fmt.Println("  \"sharePrice\": \"1.152\",")
	fmt.Println("  \"totalShares\": \"2173913.04\",")
	fmt.Println("  \"maxCapacity\": \"10000000.0\",")
	fmt.Println("  \"depositEnabled\": true,")
	fmt.Println("  \"withdrawEnabled\": true,")
	fmt.Println("  \"performanceFee\": \"0.20\",")
	fmt.Println("  \"managementFee\": \"0.02\",")
	fmt.Println("  \"lockupPeriod\": \"7 days\",")
	fmt.Println("  \"inception\": \"2023-01-15\",")
	fmt.Println("  \"strategy\": \"Market making and delta-neutral strategies\"")
	fmt.Println("}")

	// Example 3: Get user vault positions
	fmt.Println("\n--- Example 3: User Vault Positions ---")
	
	fmt.Printf("Querying vault positions for address: %s\n", address)
	fmt.Printf("Would call: info.GetUserVaultPositions(\"%s\")\n", address)
	
	// Simulate user vault positions
	fmt.Println("User vault positions (simulated):")
	fmt.Println("[")
	fmt.Println("  {")
	fmt.Println("    \"vaultAddress\": \"0xVault1...\",")
	fmt.Println("    \"vaultName\": \"Conservative Trading Vault\",")
	fmt.Println("    \"shares\": \"434.78\",")
	fmt.Println("    \"currentValue\": \"500.98\",")
	fmt.Println("    \"depositValue\": \"400.00\",")
	fmt.Println("    \"unrealizedPnl\": \"100.98\",")
	fmt.Println("    \"unrealizedPnlPercent\": \"25.25%\",")
	fmt.Println("    \"depositDate\": \"2023-10-15\",")
	fmt.Println("    \"withdrawable\": true")
	fmt.Println("  },")
	fmt.Println("  {")
	fmt.Println("    \"vaultAddress\": \"0xVault3...\",")
	fmt.Println("    \"vaultName\": \"DeFi Yield Optimizer\",")
	fmt.Println("    \"shares\": \"266.67\",")
	fmt.Println("    \"currentValue\": \"300.13\",")
	fmt.Println("    \"depositValue\": \"300.00\",")
	fmt.Println("    \"unrealizedPnl\": \"0.13\",")
	fmt.Println("    \"unrealizedPnlPercent\": \"0.04%\",")
	fmt.Println("    \"depositDate\": \"2023-11-01\",")
	fmt.Println("    \"withdrawable\": false")
	fmt.Println("  }")
	fmt.Println("]")

	// Example 4: Vault operations
	fmt.Println("\n--- Example 4: Vault Operations ---")
	
	fmt.Println("Vault operations that would be available:")
	
	depositAmount := 1000.0
	withdrawShares := 100.0
	
	fmt.Println("\n1. Deposit to vault")
	fmt.Printf("   Would call: exchange.DepositToVault(\"%s\", %.2f)\n", selectedVaultAddress, depositAmount)
	fmt.Printf("   This would deposit %.2f USDC to the vault\n", depositAmount)
	
	fmt.Println("\n2. Withdraw from vault")
	fmt.Printf("   Would call: exchange.WithdrawFromVault(\"%s\", %.2f)\n", selectedVaultAddress, withdrawShares)
	fmt.Printf("   This would withdraw %.2f shares from the vault\n", withdrawShares)
	
	fmt.Println("\n3. Request withdrawal (if lock-up applies)")
	fmt.Printf("   Would call: exchange.RequestVaultWithdrawal(\"%s\", %.2f)\n", selectedVaultAddress, withdrawShares)
	fmt.Printf("   This would request withdrawal of %.2f shares (subject to lock-up)\n", withdrawShares)

	// Example 5: Vault performance history
	fmt.Println("\n--- Example 5: Vault Performance History ---")
	
	fmt.Printf("Getting performance history for vault: %s\n", selectedVaultAddress)
	fmt.Printf("Would call: info.GetVaultPerformance(\"%s\", \"30d\")\n", selectedVaultAddress)
	
	// Simulate performance history
	fmt.Println("Vault performance history (last 30 days, simulated):")
	fmt.Println("[")
	fmt.Println("  { \"date\": \"2023-11-01\", \"sharePrice\": \"1.120\", \"tvl\": \"2300000\", \"pnl\": \"45000\" },")
	fmt.Println("  { \"date\": \"2023-11-02\", \"sharePrice\": \"1.125\", \"tvl\": \"2350000\", \"pnl\": \"23000\" },")
	fmt.Println("  { \"date\": \"2023-11-03\", \"sharePrice\": \"1.118\", \"tvl\": \"2320000\", \"pnl\": \"-18000\" },")
	fmt.Println("  { \"date\": \"2023-11-04\", \"sharePrice\": \"1.131\", \"tvl\": \"2410000\", \"pnl\": \"67000\" },")
	fmt.Println("  { \"date\": \"2023-11-05\", \"sharePrice\": \"1.152\", \"tvl\": \"2500000\", \"pnl\": \"89000\" }")
	fmt.Println("]")
	
	fmt.Println("\nPerformance summary:")
	fmt.Println("  30-day return: +2.86%")
	fmt.Println("  Best day: +3.95%")
	fmt.Println("  Worst day: -1.62%")
	fmt.Println("  Sharpe ratio: 2.34")
	fmt.Println("  Max drawdown: -2.1%")

	// Example 6: Calculate vault returns
	fmt.Println("\n--- Example 6: Return Calculations ---")
	
	fmt.Println("Vault investment return calculator:")
	
	// Example calculation for existing position
	initialDeposit := 400.0
	currentValue := 500.98
	daysInvested := 45
	
	totalReturn := currentValue - initialDeposit
	returnPercent := (totalReturn / initialDeposit) * 100
	annualizedReturn := ((currentValue / initialDeposit) - 1) * (365.0 / float64(daysInvested)) * 100
	
	fmt.Printf("Investment analysis:\n")
	fmt.Printf("  Initial deposit: $%.2f\n", initialDeposit)
	fmt.Printf("  Current value: $%.2f\n", currentValue)
	fmt.Printf("  Days invested: %d\n", daysInvested)
	fmt.Printf("  Total return: $%.2f (%.2f%%)\n", totalReturn, returnPercent)
	fmt.Printf("  Annualized return: %.2f%%\n", annualizedReturn)
	
	// Calculate fees impact
	fmt.Println("\nFee impact analysis:")
	managementFeeAnnual := 0.02 // 2%
	performanceFee := 0.20      // 20%
	
	managementFeeCost := initialDeposit * managementFeeAnnual * (float64(daysInvested) / 365.0)
	performanceFeeCost := totalReturn * performanceFee
	totalFees := managementFeeCost + performanceFeeCost
	
	fmt.Printf("  Management fee (%.1f%% annual): $%.2f\n", managementFeeAnnual*100, managementFeeCost)
	fmt.Printf("  Performance fee (%.1f%% of gains): $%.2f\n", performanceFee*100, performanceFeeCost)
	fmt.Printf("  Total fees: $%.2f\n", totalFees)
	fmt.Printf("  Net return after fees: $%.2f (%.2f%%)\n", totalReturn-performanceFeeCost, ((totalReturn-performanceFeeCost)/initialDeposit)*100)

	// Example 7: Vault comparison
	fmt.Println("\n--- Example 7: Vault Comparison ---")
	
	fmt.Println("Comparing vault strategies:")
	
	strategies := []struct {
		name           string
		apy            float64
		maxDrawdown    float64
		volatility     float64
		sharpeRatio    float64
		managementFee  float64
		performanceFee float64
	}{
		{"Conservative Trading", 15.2, 2.1, 8.5, 2.34, 2.0, 20.0},
		{"High-Freq Arbitrage", 28.7, 5.8, 12.3, 2.89, 3.0, 25.0},
		{"DeFi Yield Optimizer", 12.5, 1.8, 6.2, 2.12, 1.5, 15.0},
	}
	
	fmt.Printf("%-22s %8s %12s %11s %8s %8s %8s\n", 
		"Strategy", "APY", "Max DD", "Volatility", "Sharpe", "Mgmt Fee", "Perf Fee")
	fmt.Println("---------------------------------------------------------------------------------")
	
	for _, strategy := range strategies {
		fmt.Printf("%-22s %7.1f%% %10.1f%% %9.1f%% %8.2f %7.1f%% %7.1f%%\n",
			strategy.name, strategy.apy, strategy.maxDrawdown, strategy.volatility,
			strategy.sharpeRatio, strategy.managementFee, strategy.performanceFee)
	}

	// Example 8: Risk management
	fmt.Println("\n--- Example 8: Vault Risk Management ---")
	
	fmt.Println("Vault investment risk considerations:")
	fmt.Println("\n1. Diversification:")
	fmt.Println("   • Don't put all funds in a single vault")
	fmt.Println("   • Consider different strategy types")
	fmt.Println("   • Monitor correlation between vaults")
	
	fmt.Println("\n2. Due diligence:")
	fmt.Println("   • Research vault manager track record")
	fmt.Println("   • Understand the strategy being used")
	fmt.Println("   • Review historical performance and drawdowns")
	fmt.Println("   • Check fee structure and terms")
	
	fmt.Println("\n3. Lock-up periods:")
	fmt.Println("   • Understand withdrawal restrictions")
	fmt.Println("   • Plan for liquidity needs")
	fmt.Println("   • Consider emergency fund requirements")
	
	fmt.Println("\n4. Performance monitoring:")
	fmt.Println("   • Regular performance review")
	fmt.Println("   • Compare to benchmarks")
	fmt.Println("   • Monitor risk metrics")
	fmt.Println("   • Stay informed about strategy changes")

	// Final summary
	fmt.Println("\n--- Final Summary ---")
	
	fmt.Printf("Account: %s\n", address)
	fmt.Println("Vault status: Ready for vault operations")
	fmt.Println("Available operations:")
	fmt.Println("  ✓ Browse available vaults")
	fmt.Println("  ✓ View vault details and performance")
	fmt.Println("  ✓ Check user vault positions")
	fmt.Println("  ✓ Deposit to vaults")
	fmt.Println("  ✓ Withdraw from vaults")
	fmt.Println("  ✓ Monitor performance and returns")

	fmt.Println("\nVault example completed!")
	fmt.Println("Note: This example demonstrated:")
	fmt.Println("1. Browsing and comparing available vaults")
	fmt.Println("2. Understanding vault details and fees")
	fmt.Println("3. Monitoring vault positions and performance")
	fmt.Println("4. Calculating returns and fee impact")
	fmt.Println("5. Risk management considerations")
	fmt.Println("\nIMPORTANT: Vault investments carry risk of loss!")
	fmt.Println("Past performance does not guarantee future results.")
	fmt.Println("Always understand the strategy, fees, and risks before investing.")
	fmt.Println("Note: The actual Go SDK implementation will require specific method names")
	fmt.Println("that match the Hyperliquid API endpoints for vault functionality.")
}