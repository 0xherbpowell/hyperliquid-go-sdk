package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
address, info, _ := Setup(utils.TestnetAPIURL, true)

	// This example demonstrates vault transfer operations on Hyperliquid
	// Vault transfers allow moving funds between different vaults efficiently

	fmt.Println("Basic Vault Transfer Example")
	fmt.Printf("Account: %s\n", address)

	// Get initial user state
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("\nInitial account state:")
	PrintPositions(userState)

	// Example 1: Get current vault positions
	fmt.Println("\n--- Example 1: Current Vault Positions ---")

	fmt.Printf("Querying vault positions for address: %s\n", address)
	fmt.Printf("Would call: info.GetUserVaultPositions(\"%s\")\n", address)
	
	// Simulate current vault positions
	fmt.Println("Current vault positions (simulated):")
	vaultPositions := []struct {
		address     string
		name        string
		shares      string
		value       string
		withdrawable bool
	}{
		{"0xVault1...", "Conservative Trading Vault", "500.00", "625.50", true},
		{"0xVault2...", "High-Frequency Arbitrage", "250.00", "387.25", true},
		{"0xVault3...", "DeFi Yield Optimizer", "300.00", "318.90", false},
	}
	
	totalValue := 0.0
	for i, vault := range vaultPositions {
		fmt.Printf("%d. %s\n", i+1, vault.name)
		fmt.Printf("   Address: %s\n", vault.address)
		fmt.Printf("   Shares: %s\n", vault.shares)
		fmt.Printf("   Value: $%s\n", vault.value)
		fmt.Printf("   Withdrawable: %t\n", vault.withdrawable)
		
		// Parse value for total calculation (simplified)
		var value float64
		fmt.Sscanf(vault.value, "%f", &value)
		totalValue += value
		fmt.Println()
	}
	
	fmt.Printf("Total vault portfolio value: $%.2f\n", totalValue)

	// Example 2: Direct vault-to-vault transfer
	fmt.Println("\n--- Example 2: Direct Vault Transfer ---")
	
	fromVault := "0xVault1..."
	toVault := "0xVault2..."
	transferAmount := 100.0 // shares to transfer
	
	fmt.Printf("Transferring %.2f shares from one vault to another:\n", transferAmount)
	fmt.Printf("From vault: %s (Conservative Trading Vault)\n", fromVault)
	fmt.Printf("To vault: %s (High-Frequency Arbitrage)\n", toVault)
	fmt.Printf("Transfer amount: %.2f shares\n", transferAmount)
	
	fmt.Println("\nDirect transfer method:")
	fmt.Printf("Would call: exchange.TransferBetweenVaults(\"%s\", \"%s\", %.2f)\n", 
		fromVault, toVault, transferAmount)
	
	// Simulate transfer result
	fmt.Println("\nTransfer simulation:")
	fmt.Printf("✓ Withdrew %.2f shares from Conservative Trading Vault\n", transferAmount)
	fmt.Printf("✓ Converted to $%.2f USDC at current share price ($1.25)\n", transferAmount*1.25)
	fmt.Printf("✓ Deposited $%.2f USDC to High-Frequency Arbitrage\n", transferAmount*1.25)
	fmt.Printf("✓ Received %.2f new shares at current price ($1.55)\n", (transferAmount*1.25)/1.55)
	
	fmt.Println("\nTransaction details:")
	fmt.Println("  Transaction hash: 0x1234567890abcdef...")
	fmt.Println("  Gas used: 85,000")
	fmt.Println("  Status: Confirmed")

	// Example 3: Multi-step vault transfer (withdraw + deposit)
	fmt.Println("\n--- Example 3: Multi-Step Transfer Process ---")
	
	fmt.Println("Alternative approach: Withdraw from one vault, then deposit to another")
	
	// Step 1: Withdraw from source vault
	fmt.Println("\nStep 1: Withdraw from source vault")
	withdrawShares := 150.0
	fmt.Printf("Would call: exchange.WithdrawFromVault(\"%s\", %.2f)\n", fromVault, withdrawShares)
	
	// Simulate withdrawal
	withdrawValue := withdrawShares * 1.25 // Assume $1.25 per share
	fmt.Printf("Withdrawing %.2f shares = $%.2f USDC\n", withdrawShares, withdrawValue)
	
	// Step 2: Deposit to target vault
	fmt.Println("\nStep 2: Deposit to target vault")
	fmt.Printf("Would call: exchange.DepositToVault(\"%s\", %.2f)\n", toVault, withdrawValue)
	
	// Simulate deposit
	newShares := withdrawValue / 1.55 // Assume $1.55 per share in target vault
	fmt.Printf("Depositing $%.2f USDC = %.2f new shares\n", withdrawValue, newShares)

	// Example 4: Bulk vault rebalancing
	fmt.Println("\n--- Example 4: Vault Portfolio Rebalancing ---")
	
	fmt.Println("Current portfolio allocation:")
	conservative := 625.50
	highFreq := 387.25
	defiYield := 318.90
	total := conservative + highFreq + defiYield
	
	fmt.Printf("Conservative Trading: $%.2f (%.1f%%)\n", conservative, (conservative/total)*100)
	fmt.Printf("High-Frequency Arbitrage: $%.2f (%.1f%%)\n", highFreq, (highFreq/total)*100)
	fmt.Printf("DeFi Yield Optimizer: $%.2f (%.1f%%)\n", defiYield, (defiYield/total)*100)
	
	fmt.Println("\nTarget allocation (40% / 35% / 25%):")
	targetConservative := total * 0.40
	targetHighFreq := total * 0.35
	targetDefiYield := total * 0.25
	
	fmt.Printf("Conservative Trading: $%.2f (target)\n", targetConservative)
	fmt.Printf("High-Frequency Arbitrage: $%.2f (target)\n", targetHighFreq)
	fmt.Printf("DeFi Yield Optimizer: $%.2f (target)\n", targetDefiYield)
	
	fmt.Println("\nRequired rebalancing transfers:")
	conservativeDiff := targetConservative - conservative
	highFreqDiff := targetHighFreq - highFreq
	defiYieldDiff := targetDefiYield - defiYield
	
	if conservativeDiff > 0 {
		fmt.Printf("Add $%.2f to Conservative Trading\n", conservativeDiff)
	} else {
		fmt.Printf("Remove $%.2f from Conservative Trading\n", -conservativeDiff)
	}
	
	if highFreqDiff > 0 {
		fmt.Printf("Add $%.2f to High-Frequency Arbitrage\n", highFreqDiff)
	} else {
		fmt.Printf("Remove $%.2f from High-Frequency Arbitrage\n", -highFreqDiff)
	}
	
	if defiYieldDiff > 0 {
		fmt.Printf("Add $%.2f to DeFi Yield Optimizer\n", defiYieldDiff)
	} else {
		fmt.Printf("Remove $%.2f from DeFi Yield Optimizer\n", -defiYieldDiff)
	}

	// Example 5: Transfer with slippage consideration
	fmt.Println("\n--- Example 5: Transfer with Slippage Protection ---")
	
	fmt.Println("Transfer considerations and slippage protection:")
	
	transferValue := 500.0
	fmt.Printf("Transferring $%.2f between vaults\n", transferValue)
	
	// Source vault share price and slippage
	sourceSharePrice := 1.25
	sourceSlippage := 0.002 // 0.2%
	effectiveWithdrawPrice := sourceSharePrice * (1 - sourceSlippage)
	
	// Target vault share price and slippage  
	targetSharePrice := 1.55
	targetSlippage := 0.003 // 0.3%
	effectiveDepositPrice := targetSharePrice * (1 + targetSlippage)
	
	fmt.Printf("\nSlippage analysis:\n")
	fmt.Printf("Source vault:\n")
	fmt.Printf("  Share price: $%.3f\n", sourceSharePrice)
	fmt.Printf("  Withdrawal slippage: %.1f%%\n", sourceSlippage*100)
	fmt.Printf("  Effective price: $%.3f\n", effectiveWithdrawPrice)
	
	fmt.Printf("Target vault:\n")
	fmt.Printf("  Share price: $%.3f\n", targetSharePrice)
	fmt.Printf("  Deposit slippage: %.1f%%\n", targetSlippage*100)
	fmt.Printf("  Effective price: $%.3f\n", effectiveDepositPrice)
	
	// Calculate final shares received
	sharesWithdrawn := transferValue / effectiveWithdrawPrice
	finalShares := transferValue / effectiveDepositPrice
	
	fmt.Printf("\nTransfer result:\n")
	fmt.Printf("Shares withdrawn from source: %.3f\n", sharesWithdrawn)
	fmt.Printf("Shares received in target: %.3f\n", finalShares)
	fmt.Printf("Total slippage cost: $%.2f (%.2f%%)\n", 
		transferValue-(finalShares*targetSharePrice), 
		((transferValue-(finalShares*targetSharePrice))/transferValue)*100)

	// Example 6: Transfer scheduling and timing
	fmt.Println("\n--- Example 6: Transfer Timing Strategy ---")
	
	fmt.Println("Optimal transfer timing considerations:")
	
	fmt.Println("\n1. Market conditions:")
	fmt.Println("   • Low volatility periods reduce slippage")
	fmt.Println("   • Avoid transfers during high-volume events")
	fmt.Println("   • Consider time zones and market hours")
	
	fmt.Println("\n2. Vault-specific factors:")
	fmt.Println("   • Check vault liquidity and capacity")
	fmt.Println("   • Avoid transfers near strategy rebalancing")
	fmt.Println("   • Consider vault performance cycles")
	
	fmt.Println("\n3. Cost optimization:")
	fmt.Println("   • Batch multiple transfers when possible")
	fmt.Println("   • Consider gas costs vs. transfer benefits")
	fmt.Println("   • Use limit orders if available")
	
	// Simulate scheduled transfer
	fmt.Println("\nScheduled transfer example:")
	fmt.Printf("Would call: exchange.ScheduleVaultTransfer(\"%s\", \"%s\", %.2f, \"2023-11-15T10:00:00Z\")\n", 
		fromVault, toVault, 200.0)
	fmt.Println("Transfer scheduled for: 2023-11-15 10:00 UTC (low volatility period)")

	// Example 7: Transfer cost analysis
	fmt.Println("\n--- Example 7: Transfer Cost Analysis ---")
	
	fmt.Println("Cost breakdown for vault transfers:")
	
	transferAmount = 1000.0
	
	// Gas costs
	gasCost := 15.0 // USD
	
	// Slippage costs (estimated)
	withdrawalSlippage := transferAmount * 0.002 // 0.2%
	depositSlippage := transferAmount * 0.003    // 0.3%
	totalSlippage := withdrawalSlippage + depositSlippage
	
	// Opportunity cost (time out of market)
	timeOutOfMarket := 2.0 // minutes
	dailyReturn := 0.0012  // 0.12% daily return
	opportunityCost := transferAmount * dailyReturn * (timeOutOfMarket / 1440) // minutes in a day
	
	totalCost := gasCost + totalSlippage + opportunityCost
	
	fmt.Printf("Transfer amount: $%.2f\n", transferAmount)
	fmt.Printf("Costs breakdown:\n")
	fmt.Printf("  Gas fees: $%.2f (%.3f%%)\n", gasCost, (gasCost/transferAmount)*100)
	fmt.Printf("  Withdrawal slippage: $%.2f (%.3f%%)\n", withdrawalSlippage, (withdrawalSlippage/transferAmount)*100)
	fmt.Printf("  Deposit slippage: $%.2f (%.3f%%)\n", depositSlippage, (depositSlippage/transferAmount)*100)
	fmt.Printf("  Opportunity cost: $%.2f (%.3f%%)\n", opportunityCost, (opportunityCost/transferAmount)*100)
	fmt.Printf("Total transfer cost: $%.2f (%.3f%%)\n", totalCost, (totalCost/transferAmount)*100)
	
	fmt.Println("\nCost optimization suggestions:")
	if totalCost/transferAmount > 0.01 { // 1%
		fmt.Println("  ⚠️  Transfer cost is high (>1%) - consider:")
		fmt.Println("     • Waiting for better market conditions")
		fmt.Println("     • Increasing transfer size to reduce relative costs")
		fmt.Println("     • Using direct vault-to-vault transfer if available")
	} else {
		fmt.Println("  ✅ Transfer cost is reasonable (<1%)")
	}

	// Example 8: Transfer history and tracking
	fmt.Println("\n--- Example 8: Transfer History ---")
	
	fmt.Printf("Querying transfer history for address: %s\n", address)
	fmt.Printf("Would call: info.GetVaultTransferHistory(\"%s\", \"30d\")\n", address)
	
	// Simulate transfer history
	fmt.Println("Recent vault transfers (last 30 days, simulated):")
	transfers := []struct {
		date      string
		fromVault string
		toVault   string
		amount    string
		status    string
	}{
		{"2023-11-01", "Conservative Trading", "High-Freq Arbitrage", "$250.00", "completed"},
		{"2023-10-28", "DeFi Yield Optimizer", "Conservative Trading", "$500.00", "completed"},
		{"2023-10-25", "High-Freq Arbitrage", "DeFi Yield Optimizer", "$150.00", "completed"},
	}
	
	fmt.Printf("%-12s %-20s %-20s %-10s %-10s\n", 
		"Date", "From Vault", "To Vault", "Amount", "Status")
	fmt.Println("--------------------------------------------------------------------------------")
	
	for _, transfer := range transfers {
		fmt.Printf("%-12s %-20s %-20s %-10s %-10s\n",
			transfer.date, transfer.fromVault, transfer.toVault, 
			transfer.amount, transfer.status)
	}
	
	fmt.Println("\nTransfer statistics:")
	fmt.Println("  Total transfers this month: 3")
	fmt.Println("  Total value transferred: $900.00")
	fmt.Println("  Average transfer size: $300.00")
	fmt.Println("  Success rate: 100%")

	// Final summary
	fmt.Println("\n--- Final Summary ---")
	
	fmt.Printf("Account: %s\n", address)
	fmt.Println("Vault transfer status: Ready for transfer operations")
	fmt.Println("Available transfer methods:")
	fmt.Println("  ✓ Direct vault-to-vault transfers")
	fmt.Println("  ✓ Multi-step withdraw and deposit")
	fmt.Println("  ✓ Bulk portfolio rebalancing")
	fmt.Println("  ✓ Scheduled transfers")
	fmt.Println("  ✓ Transfer cost analysis")
	fmt.Println("  ✓ Transfer history tracking")

	fmt.Println("\nVault transfer example completed!")
	fmt.Println("Note: This example demonstrated:")
	fmt.Println("1. Viewing current vault positions")
	fmt.Println("2. Direct vault-to-vault transfers")
	fmt.Println("3. Multi-step transfer processes")
	fmt.Println("4. Portfolio rebalancing strategies")
	fmt.Println("5. Slippage protection and timing")
	fmt.Println("6. Cost analysis and optimization")
	fmt.Println("7. Transfer history and tracking")
	fmt.Println("\nIMPORTANT: Always consider transfer costs and timing!")
	fmt.Println("Transfers may involve slippage, gas fees, and opportunity costs.")
	fmt.Println("Note: The actual Go SDK implementation will require specific method names")
	fmt.Println("that match the Hyperliquid API endpoints for vault transfer functionality.")
}