package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, _ := Setup(utils.TestnetAPIURL, true)

	// This example demonstrates sub-account operations on Hyperliquid
	// Sub-accounts allow users to segregate funds and manage multiple trading strategies

	fmt.Println("Basic Sub-Account Example")
	fmt.Printf("Main Account: %s\n", address)

	// Get initial user state
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("\nMain account state:")
	PrintPositions(userState)

	// Example 1: List existing sub-accounts
	fmt.Println("\n--- Example 1: Sub-Account Management ---")

	fmt.Printf("Querying sub-accounts for main account: %s\n", address)
	fmt.Printf("Would call: info.GetSubAccounts(\"%s\")\n", address)

	// Simulate sub-account list
	fmt.Println("Existing sub-accounts (simulated):")
	subAccounts := []struct {
		address string
		name    string
		balance string
		equity  string
		status  string
		created string
	}{
		{"0xSub1...", "Arbitrage Bot", "$5,250.00", "$5,387.25", "active", "2023-09-15"},
		{"0xSub2...", "Long-term Holdings", "$10,000.00", "$12,450.80", "active", "2023-08-20"},
		{"0xSub3...", "High-Risk Trading", "$2,500.00", "$1,987.50", "active", "2023-10-01"},
		{"0xSub4...", "Test Account", "$100.00", "$95.75", "suspended", "2023-11-01"},
	}

	totalEquity := 0.0
	for i, sub := range subAccounts {
		fmt.Printf("%d. %s\n", i+1, sub.name)
		fmt.Printf("   Address: %s\n", sub.address)
		fmt.Printf("   Balance: %s\n", sub.balance)
		fmt.Printf("   Equity: %s\n", sub.equity)
		fmt.Printf("   Status: %s\n", sub.status)
		fmt.Printf("   Created: %s\n", sub.created)

		// Parse equity for total calculation (simplified)
		var equity float64
		fmt.Sscanf(sub.equity, "$%f", &equity)
		if sub.status == "active" {
			totalEquity += equity
		}
		fmt.Println()
	}

	fmt.Printf("Total active sub-account equity: $%.2f\n", totalEquity)

	// Example 2: Create a new sub-account
	fmt.Println("\n--- Example 2: Creating New Sub-Account ---")

	newSubAccountName := "DeFi Strategy Bot"
	initialDeposit := 1000.0

	fmt.Printf("Creating new sub-account: %s\n", newSubAccountName)
	fmt.Printf("Initial deposit: $%.2f\n", initialDeposit)
	fmt.Printf("Would call: exchange.CreateSubAccount(\"%s\", %.2f)\n", newSubAccountName, initialDeposit)

	// Simulate sub-account creation
	newSubAddress := "0xSub5..."
	fmt.Println("\nSub-account creation simulation:")
	fmt.Printf("✓ Sub-account created successfully\n")
	fmt.Printf("✓ Sub-account address: %s\n", newSubAddress)
	fmt.Printf("✓ Initial deposit transferred: $%.2f\n", initialDeposit)
	fmt.Printf("✓ Sub-account status: Active\n")
	fmt.Printf("✓ Trading permissions: Enabled\n")

	fmt.Println("\nNew sub-account details:")
	fmt.Printf("  Name: %s\n", newSubAccountName)
	fmt.Printf("  Address: %s\n", newSubAddress)
	fmt.Printf("  Parent: %s\n", address)
	fmt.Printf("  Balance: $%.2f\n", initialDeposit)
	fmt.Printf("  Available margin: $%.2f\n", initialDeposit*0.95) // Assuming 5% reserved

	// Example 3: Transfer funds between accounts
	fmt.Println("\n--- Example 3: Inter-Account Transfers ---")

	fmt.Println("Transfer operations between main account and sub-accounts:")

	// Transfer from main to sub-account
	transferToSub := 500.0
	targetSubAccount := "0xSub1..."

	fmt.Println("\n1. Transfer from main account to sub-account:")
	fmt.Printf("   From: %s (Main Account)\n", address)
	fmt.Printf("   To: %s (Arbitrage Bot)\n", targetSubAccount)
	fmt.Printf("   Amount: $%.2f\n", transferToSub)
	fmt.Printf("   Would call: exchange.TransferToSubAccount(\"%s\", %.2f)\n", targetSubAccount, transferToSub)

	// Transfer from sub-account to main
	transferToMain := 300.0
	sourceSubAccount := "0xSub2..."

	fmt.Println("\n2. Transfer from sub-account to main account:")
	fmt.Printf("   From: %s (Long-term Holdings)\n", sourceSubAccount)
	fmt.Printf("   To: %s (Main Account)\n", address)
	fmt.Printf("   Amount: $%.2f\n", transferToMain)
	fmt.Printf("   Would call: exchange.TransferFromSubAccount(\"%s\", %.2f)\n", sourceSubAccount, transferToMain)

	// Transfer between sub-accounts
	transferBetweenSubs := 200.0
	fromSub := "0xSub1..."
	toSub := "0xSub3..."

	fmt.Println("\n3. Transfer between sub-accounts:")
	fmt.Printf("   From: %s (Arbitrage Bot)\n", fromSub)
	fmt.Printf("   To: %s (High-Risk Trading)\n", toSub)
	fmt.Printf("   Amount: $%.2f\n", transferBetweenSubs)
	fmt.Printf("   Would call: exchange.TransferBetweenSubAccounts(\"%s\", \"%s\", %.2f)\n",
		fromSub, toSub, transferBetweenSubs)

	// Example 4: Sub-account permissions and restrictions
	fmt.Println("\n--- Example 4: Sub-Account Permissions ---")

	fmt.Println("Configuring sub-account permissions and restrictions:")

	targetSub := "0xSub5..."

	fmt.Printf("Setting permissions for sub-account: %s\n", targetSub)

	// Trading permissions
	fmt.Println("\n1. Trading Permissions:")
	fmt.Printf("   Would call: exchange.SetSubAccountTradingPermissions(\"%s\", true, true, false)\n", targetSub)
	fmt.Printf("   • Spot trading: Enabled\n")
	fmt.Printf("   • Derivatives trading: Enabled\n")
	fmt.Printf("   • Withdrawal permissions: Disabled\n")

	// Risk limits
	fmt.Println("\n2. Risk Limits:")
	maxPositionSize := 5000.0
	maxLeverage := 10.0
	fmt.Printf("   Would call: exchange.SetSubAccountRiskLimits(\"%s\", %.2f, %.1fx)\n",
		targetSub, maxPositionSize, maxLeverage)
	fmt.Printf("   • Maximum position size: $%.2f\n", maxPositionSize)
	fmt.Printf("   • Maximum leverage: %.1fx\n", maxLeverage)
	fmt.Printf("   • Daily loss limit: $500.00\n")

	// API access
	fmt.Println("\n3. API Access:")
	fmt.Printf("   Would call: exchange.GenerateSubAccountAPIKeys(\"%s\", [\"trading\", \"read\"])\n", targetSub)
	fmt.Printf("   • API key generated: sub_****\n")
	fmt.Printf("   • Permissions: Trading, Read-only\n")
	fmt.Printf("   • Rate limits: Standard\n")

	// Example 5: Sub-account performance monitoring
	fmt.Println("\n--- Example 5: Performance Monitoring ---")

	fmt.Printf("Monitoring performance across all sub-accounts:\n")
	fmt.Printf("Would call: info.GetSubAccountsPerformance(\"%s\", \"30d\")\n", address)

	// Simulate performance data
	fmt.Println("Sub-account performance (last 30 days, simulated):")

	performances := []struct {
		name         string
		startValue   float64
		currentValue float64
		pnl          float64
		pnlPercent   float64
		trades       int
		winRate      float64
	}{
		{"Arbitrage Bot", 5000.00, 5387.25, 387.25, 7.75, 1247, 68.5},
		{"Long-term Holdings", 10000.00, 12450.80, 2450.80, 24.51, 23, 87.0},
		{"High-Risk Trading", 2500.00, 1987.50, -512.50, -20.50, 89, 42.7},
		{"DeFi Strategy Bot", 1000.00, 1000.00, 0.00, 0.00, 0, 0.0},
	}

	fmt.Printf("%-20s %12s %12s %10s %8s %7s %8s\n",
		"Sub-Account", "Start Value", "Current", "PnL", "PnL %", "Trades", "Win Rate")
	fmt.Println("--------------------------------------------------------------------------------")

	totalPnL := 0.0
	for _, perf := range performances {
		fmt.Printf("%-20s $%10.2f $%10.2f $%8.2f %7.2f%% %6d %7.1f%%\n",
			perf.name, perf.startValue, perf.currentValue,
			perf.pnl, perf.pnlPercent, perf.trades, perf.winRate)
		totalPnL += perf.pnl
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%-20s %12s %12s $%8.2f\n", "TOTAL", "", "", totalPnL)

	// Example 6: Sub-account risk management
	fmt.Println("\n--- Example 6: Risk Management ---")

	fmt.Println("Sub-account risk management and monitoring:")

	// Risk metrics per sub-account
	fmt.Println("\n1. Risk Metrics:")
	riskMetrics := []struct {
		name       string
		leverage   float64
		marginUsed float64
		marginFree float64
		riskScore  string
		alerts     int
	}{
		{"Arbitrage Bot", 2.1, 2850.00, 2537.25, "Low", 0},
		{"Long-term Holdings", 1.3, 4200.00, 8250.80, "Low", 0},
		{"High-Risk Trading", 8.7, 1750.00, 237.50, "High", 3},
		{"DeFi Strategy Bot", 0.0, 0.00, 1000.00, "None", 0},
	}

	fmt.Printf("%-20s %8s %12s %12s %10s %7s\n",
		"Sub-Account", "Leverage", "Margin Used", "Margin Free", "Risk", "Alerts")
	fmt.Println("------------------------------------------------------------------------")

	for _, risk := range riskMetrics {
		fmt.Printf("%-20s %7.1fx $%10.2f $%10.2f %9s %6d\n",
			risk.name, risk.leverage, risk.marginUsed,
			risk.marginFree, risk.riskScore, risk.alerts)
	}

	// Risk alerts
	fmt.Println("\n2. Active Risk Alerts:")
	fmt.Println("   Sub-Account: High-Risk Trading")
	fmt.Println("   ⚠️  Alert: High leverage detected (8.7x)")
	fmt.Println("   ⚠️  Alert: Low margin remaining ($237.50)")
	fmt.Println("   ⚠️  Alert: Daily loss approaching limit (-$450/$500)")

	// Auto-actions
	fmt.Println("\n3. Automated Risk Actions:")
	fmt.Printf("   Would call: exchange.SetAutoLiquidation(\"%s\", true, 0.15)\n", "0xSub3...")
	fmt.Println("   • Auto-liquidation enabled at 15% margin")
	fmt.Println("   • Position size limits enforced")
	fmt.Println("   • Daily trading limits active")

	// Example 7: Sub-account consolidation and management
	fmt.Println("\n--- Example 7: Account Consolidation ---")

	fmt.Println("Sub-account consolidation and cleanup:")

	// Consolidate funds
	fmt.Println("\n1. Consolidating funds to main account:")
	consolidationAccounts := []struct {
		address string
		name    string
		balance float64
	}{
		{"0xSub4...", "Test Account", 95.75},
		{"0xSub3...", "High-Risk Trading", 1987.50},
	}

	totalConsolidation := 0.0
	for _, acc := range consolidationAccounts {
		fmt.Printf("   From: %s (%s) - $%.2f\n", acc.address, acc.name, acc.balance)
		totalConsolidation += acc.balance
		fmt.Printf("   Would call: exchange.TransferFromSubAccount(\"%s\", %.2f)\n",
			acc.address, acc.balance)
	}
	fmt.Printf("   Total consolidated: $%.2f\n", totalConsolidation)

	// Close unused sub-accounts
	fmt.Println("\n2. Closing unused sub-accounts:")
	fmt.Printf("   Closing: 0xSub4... (Test Account)\n")
	fmt.Printf("   Would call: exchange.CloseSubAccount(\"0xSub4...\")\n")
	fmt.Printf("   • All funds transferred to main account\n")
	fmt.Printf("   • Trading permissions revoked\n")
	fmt.Printf("   • API keys deactivated\n")
	fmt.Printf("   • Account status: Closed\n")

	// Example 8: Sub-account reporting and analytics
	fmt.Println("\n--- Example 8: Reporting and Analytics ---")

	fmt.Printf("Generating sub-account reports for main account: %s\n", address)

	// Portfolio summary
	fmt.Println("\n1. Portfolio Summary:")
	mainBalance := 15000.00
	subAccountsTotal := totalEquity
	overallTotal := mainBalance + subAccountsTotal

	fmt.Printf("   Main Account: $%.2f (%.1f%%)\n",
		mainBalance, (mainBalance/overallTotal)*100)
	fmt.Printf("   Sub-Accounts: $%.2f (%.1f%%)\n",
		subAccountsTotal, (subAccountsTotal/overallTotal)*100)
	fmt.Printf("   Total Portfolio: $%.2f\n", overallTotal)

	// Performance attribution
	fmt.Println("\n2. Performance Attribution:")
	fmt.Printf("   Main Account PnL: $%.2f\n", 245.80)
	fmt.Printf("   Sub-Accounts PnL: $%.2f\n", totalPnL)
	fmt.Printf("   Combined PnL: $%.2f\n", 245.80+totalPnL)
	fmt.Printf("   Best Performer: Long-term Holdings (+24.51%%)\n")
	fmt.Printf("   Worst Performer: High-Risk Trading (-20.50%%)\n")

	// Strategy allocation
	fmt.Println("\n3. Strategy Allocation:")
	fmt.Printf("   Arbitrage: 25.8%% ($%.2f)\n", 5387.25)
	fmt.Printf("   Long-term: 59.6%% ($%.2f)\n", 12450.80)
	fmt.Printf("   High-risk: 9.5%% ($%.2f)\n", 1987.50)
	fmt.Printf("   Cash/New: 5.1%% ($%.2f)\n", 1000.00)

	// Final summary
	fmt.Println("\n--- Final Summary ---")

	fmt.Printf("Main Account: %s\n", address)
	fmt.Printf("Active Sub-Accounts: %d\n", 4)
	fmt.Printf("Total Portfolio Value: $%.2f\n", overallTotal)
	fmt.Println("Sub-account operations available:")
	fmt.Println("  ✓ Create and manage sub-accounts")
	fmt.Println("  ✓ Transfer funds between accounts")
	fmt.Println("  ✓ Set permissions and risk limits")
	fmt.Println("  ✓ Monitor performance and risk")
	fmt.Println("  ✓ Consolidate and close accounts")
	fmt.Println("  ✓ Generate reports and analytics")

	fmt.Println("\nSub-account example completed!")
	fmt.Println("Note: This example demonstrated:")
	fmt.Println("1. Managing multiple sub-accounts")
	fmt.Println("2. Creating and configuring new sub-accounts")
	fmt.Println("3. Inter-account fund transfers")
	fmt.Println("4. Setting permissions and risk limits")
	fmt.Println("5. Performance monitoring across accounts")
	fmt.Println("6. Risk management and alerts")
	fmt.Println("7. Account consolidation strategies")
	fmt.Println("8. Comprehensive reporting and analytics")
	fmt.Println("\nIMPORTANT: Sub-accounts provide isolation but require careful management!")
	fmt.Println("Monitor risk limits, permissions, and performance regularly.")
	fmt.Println("Note: The actual Go SDK implementation will require specific method names")
	fmt.Println("that match the Hyperliquid API endpoints for sub-account functionality.")
}
