package main

import (
	"fmt"
	"log"
	"strings"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
address, info, _ := Setup(utils.TestnetAPIURL, true)

	// This example demonstrates builder fee operations on Hyperliquid
	// Builder fees are fees paid to block builders/validators in the Hyperliquid ecosystem

	fmt.Println("Basic Builder Fee Example")
	fmt.Printf("Account: %s\n", address)

	// Get initial user state
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("\nInitial account state:")
	PrintPositions(userState)

	// Example 1: Understanding builder fees
	fmt.Println("\n--- Example 1: Builder Fee Concept ---")

	fmt.Println("Builder fees in Hyperliquid:")
	fmt.Println("• Builder fees are paid to validators/block producers")
	fmt.Println("• They help prioritize transactions during congestion")
	fmt.Println("• Higher fees = higher priority and faster execution")
	fmt.Println("• Fees are dynamic based on network conditions")
	fmt.Println("• Optional but recommended for time-sensitive trades")

	// Example 2: Current builder fee rates
	fmt.Println("\n--- Example 2: Current Builder Fee Rates ---")

	fmt.Printf("Querying current builder fee rates...\n")
	fmt.Printf("Would call: info.GetBuilderFeeRates()\n")
	
	// Simulate current rates
	fmt.Println("Current builder fee rates (simulated):")
	feeRates := []struct {
		priority    string
		feeRate     string
		description string
		avgWaitTime string
	}{
		{"Low", "0.01%", "Standard processing", "~10-30 seconds"},
		{"Medium", "0.05%", "Priority processing", "~5-15 seconds"},
		{"High", "0.10%", "Express processing", "~1-5 seconds"},
		{"Ultra", "0.25%", "Immediate processing", "~0-2 seconds"},
	}
	
	fmt.Printf("%-8s %-8s %-20s %-15s\n", "Priority", "Fee Rate", "Description", "Avg Wait Time")
	fmt.Println("---------------------------------------------------------------")
	
	for _, rate := range feeRates {
		fmt.Printf("%-8s %-8s %-20s %-15s\n", 
			rate.priority, rate.feeRate, rate.description, rate.avgWaitTime)
	}

	// Example 3: Network congestion analysis
	fmt.Println("\n--- Example 3: Network Congestion Analysis ---")
	
	fmt.Printf("Analyzing network congestion...\n")
	fmt.Printf("Would call: info.GetNetworkCongestion()\n")
	
	// Simulate network status
	fmt.Println("Current network status (simulated):")
	congestionLevel := "Medium"
	pendingTxs := 1247
	avgBlockTime := 2.3
	recommendedFee := "0.05%"
	
	fmt.Printf("Congestion Level: %s\n", congestionLevel)
	fmt.Printf("Pending Transactions: %d\n", pendingTxs)
	fmt.Printf("Average Block Time: %.1f seconds\n", avgBlockTime)
	fmt.Printf("Recommended Builder Fee: %s\n", recommendedFee)
	
	fmt.Println("\nCongestion impact:")
	if congestionLevel == "High" || congestionLevel == "Ultra" {
		fmt.Println("⚠️  High congestion detected!")
		fmt.Println("• Consider higher builder fees for faster execution")
		fmt.Println("• Delay non-urgent transactions if possible")
		fmt.Println("• Monitor for better conditions")
	} else if congestionLevel == "Medium" {
		fmt.Println("⚡ Moderate congestion")
		fmt.Println("• Standard fees may experience delays")
		fmt.Println("• Consider medium priority for important trades")
	} else {
		fmt.Println("✅ Low congestion")
		fmt.Println("• Standard fees should process quickly")
		fmt.Println("• Good time for routine transactions")
	}

	// Example 4: Builder fee calculation
	fmt.Println("\n--- Example 4: Builder Fee Calculation ---")
	
	fmt.Println("Calculating builder fees for different trade sizes:")
	
	tradeSizes := []float64{100, 500, 1000, 5000, 10000}
	feePercentages := []float64{0.0001, 0.0005, 0.001, 0.0025} // 0.01%, 0.05%, 0.10%, 0.25%
	feeNames := []string{"Low", "Medium", "High", "Ultra"}
	
	fmt.Printf("\n%-12s", "Trade Size")
	for _, name := range feeNames {
		fmt.Printf("%-10s", name)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", 52))
	
	for _, size := range tradeSizes {
		fmt.Printf("$%-11.0f", size)
		for _, feePercent := range feePercentages {
			fee := size * feePercent
			fmt.Printf("$%-9.2f", fee)
		}
		fmt.Println()
	}
	
	fmt.Println("\nFee optimization tips:")
	fmt.Println("• For small trades (<$1000): Use low priority fees")
	fmt.Println("• For medium trades ($1000-$5000): Consider network conditions")
	fmt.Println("• For large trades (>$5000): Higher fees may be cost-effective")
	fmt.Println("• For arbitrage/time-sensitive: Use high/ultra priority")

	// Example 5: Setting builder fees for trades
	fmt.Println("\n--- Example 5: Setting Builder Fees ---")
	
	fmt.Println("Setting builder fees for different trade types:")
	
	// Market order with builder fee
	fmt.Println("\n1. Market order with builder fee:")
	symbol := "ETH"
	size := 1.0
	builderFee := 0.001 // 0.10%
	
	fmt.Printf("Symbol: %s\n", symbol)
	fmt.Printf("Size: %.3f\n", size)
	fmt.Printf("Builder fee: %.2f%%\n", builderFee*100)
	fmt.Printf("Would call: exchange.PlaceMarketOrderWithBuilderFee(\"%s\", \"buy\", %.3f, %.4f)\n", 
		symbol, size, builderFee)
	
	// Limit order with builder fee
	fmt.Println("\n2. Limit order with builder fee:")
	price := 2150.50
	fmt.Printf("Symbol: %s\n", symbol)
	fmt.Printf("Size: %.3f\n", size)
	fmt.Printf("Price: $%.2f\n", price)
	fmt.Printf("Builder fee: %.2f%%\n", builderFee*100)
	fmt.Printf("Would call: exchange.PlaceLimitOrderWithBuilderFee(\"%s\", \"buy\", %.3f, %.2f, %.4f)\n", 
		symbol, size, price, builderFee)
	
	// Conditional order with builder fee
	fmt.Println("\n3. Stop-loss order with builder fee:")
	stopPrice := 2000.00
	fmt.Printf("Symbol: %s\n", symbol)
	fmt.Printf("Size: %.3f\n", size)
	fmt.Printf("Stop price: $%.2f\n", stopPrice)
	fmt.Printf("Builder fee: %.2f%%\n", builderFee*100)
	fmt.Printf("Would call: exchange.PlaceStopLossWithBuilderFee(\"%s\", \"sell\", %.3f, %.2f, %.4f)\n", 
		symbol, size, stopPrice, builderFee)

	// Example 6: Dynamic builder fee adjustment
	fmt.Println("\n--- Example 6: Dynamic Fee Adjustment ---")
	
	fmt.Println("Dynamic builder fee adjustment based on conditions:")
	
	// Simulate trading conditions
	conditions := []struct {
		scenario      string
		urgency       string
		tradeSize     float64
		congestion    string
		recommendedFee float64
		reasoning     string
	}{
		{
			"Arbitrage Opportunity", 
			"Critical", 
			2500.00, 
			"Medium", 
			0.0025,
			"Time-sensitive trade, profits exceed fee cost",
		},
		{
			"DCA Purchase", 
			"Low", 
			500.00, 
			"Low", 
			0.0001,
			"Routine purchase, no time pressure",
		},
		{
			"Risk Management Exit", 
			"High", 
			8000.00, 
			"High", 
			0.0010,
			"Important exit, moderate fee acceptable",
		},
		{
			"Portfolio Rebalance", 
			"Medium", 
			1500.00, 
			"Medium", 
			0.0005,
			"Planned trade, moderate priority",
		},
	}
	
	fmt.Printf("%-20s %-8s %-10s %-11s %-8s %s\n", 
		"Scenario", "Urgency", "Size", "Congestion", "Fee %", "Reasoning")
	fmt.Println(strings.Repeat("-", 85))
	
	for _, cond := range conditions {
		fmt.Printf("%-20s %-8s $%-9.0f %-11s %-7.2f%% %s\n",
			cond.scenario, cond.urgency, cond.tradeSize, 
			cond.congestion, cond.recommendedFee*100, cond.reasoning)
	}
	
	// Auto-adjustment logic
	fmt.Println("\nAuto-adjustment logic example:")
	fmt.Printf("Would call: exchange.EnableDynamicBuilderFees(true, \"balanced\")\n")
	fmt.Println("Settings:")
	fmt.Println("• Low urgency trades: Use minimum fees")
	fmt.Println("• Medium urgency: Adjust based on congestion")
	fmt.Println("• High urgency: Use recommended fees")
	fmt.Println("• Critical urgency: Use maximum fees")

	// Example 7: Builder fee analytics
	fmt.Println("\n--- Example 7: Builder Fee Analytics ---")
	
	fmt.Printf("Analyzing builder fee usage and effectiveness...\n")
	fmt.Printf("Would call: info.GetBuilderFeeAnalytics(\"%s\", \"30d\")\n", address)
	
	// Simulate analytics data
	fmt.Println("Builder fee analytics (last 30 days, simulated):")
	
	analytics := struct {
		totalTrades      int
		tradesWithFees   int
		totalFeePaid     float64
		avgFeePercent    float64
		executionTimes   map[string]float64
		savingsVsSlippage float64
	}{
		totalTrades:    156,
		tradesWithFees: 89,
		totalFeePaid:   127.35,
		avgFeePercent:  0.0008,
		executionTimes: map[string]float64{
			"No Fee":   18.3,
			"Low Fee":  12.7,
			"Med Fee":  6.4,
			"High Fee": 2.1,
		},
		savingsVsSlippage: 234.67,
	}
	
	fmt.Printf("Total trades: %d\n", analytics.totalTrades)
	fmt.Printf("Trades with builder fees: %d (%.1f%%)\n", 
		analytics.tradesWithFees, 
		float64(analytics.tradesWithFees)/float64(analytics.totalTrades)*100)
	fmt.Printf("Total builder fees paid: $%.2f\n", analytics.totalFeePaid)
	fmt.Printf("Average fee rate: %.3f%%\n", analytics.avgFeePercent*100)
	
	fmt.Println("\nExecution time analysis:")
	for feeType, avgTime := range analytics.executionTimes {
		fmt.Printf("  %s: %.1f seconds average\n", feeType, avgTime)
	}
	
	fmt.Printf("\nCost-benefit analysis:\n")
	fmt.Printf("Total fees paid: $%.2f\n", analytics.totalFeePaid)
	fmt.Printf("Estimated slippage savings: $%.2f\n", analytics.savingsVsSlippage)
	fmt.Printf("Net benefit: $%.2f\n", analytics.savingsVsSlippage-analytics.totalFeePaid)
	
	if analytics.savingsVsSlippage > analytics.totalFeePaid {
		fmt.Printf("✅ Builder fees provided net benefit\n")
	} else {
		fmt.Printf("⚠️  Consider adjusting builder fee strategy\n")
	}

	// Example 8: Builder fee best practices
	fmt.Println("\n--- Example 8: Builder Fee Best Practices ---")
	
	fmt.Println("Builder fee best practices and strategies:")
	
	fmt.Println("\n1. Trade Type Considerations:")
	fmt.Println("   • Arbitrage trades: Always use high/ultra fees")
	fmt.Println("   • Liquidation avoidance: Use medium/high fees")
	fmt.Println("   • DCA/regular buys: Use low/no fees")
	fmt.Println("   • Portfolio rebalancing: Use low/medium fees")
	
	fmt.Println("\n2. Market Condition Adjustments:")
	fmt.Println("   • High volatility: Increase fee priority")
	fmt.Println("   • Low liquidity: Use higher fees for better fills")
	fmt.Println("   • Network congestion: Scale fees with conditions")
	fmt.Println("   • Off-peak hours: Lower fees often sufficient")
	
	fmt.Println("\n3. Cost-Benefit Analysis:")
	fmt.Println("   • Calculate potential slippage vs. fee cost")
	fmt.Println("   • Consider time value of execution")
	fmt.Println("   • Factor in opportunity costs")
	fmt.Println("   • Monitor fee effectiveness over time")
	
	fmt.Println("\n4. Automation Strategies:")
	fmt.Printf("   Would call: exchange.SetBuilderFeeRules({\n")
	fmt.Printf("     \"arbitrage\": 0.0025,\n")
	fmt.Printf("     \"risk_management\": 0.0010,\n")
	fmt.Printf("     \"routine\": 0.0001,\n")
	fmt.Printf("     \"congestion_multiplier\": 1.5\n")
	fmt.Printf("   })\n")
	
	fmt.Println("\n5. Monitoring and Optimization:")
	fmt.Println("   • Track execution times by fee level")
	fmt.Println("   • Monitor slippage vs. fee costs")
	fmt.Println("   • Analyze success rates by priority")
	fmt.Println("   • Adjust strategy based on performance")

	// Example 9: Builder fee estimation tool
	fmt.Println("\n--- Example 9: Fee Estimation Tool ---")
	
	fmt.Println("Builder fee estimation for upcoming trade:")
	
	// Trade parameters
	estimateSymbol := "BTC"
	estimateSize := 0.5
	estimateValue := 21500.00 * estimateSize
	urgencyLevel := "Medium"
	currentCongestion := "Medium"
	
	fmt.Printf("Trade details:\n")
	fmt.Printf("  Symbol: %s\n", estimateSymbol)
	fmt.Printf("  Size: %.3f\n", estimateSize)
	fmt.Printf("  Estimated value: $%.2f\n", estimateValue)
	fmt.Printf("  Urgency: %s\n", urgencyLevel)
	fmt.Printf("  Network congestion: %s\n", currentCongestion)
	
	// Fee estimation
	baseFeeRate := 0.0005 // 0.05% for medium urgency
	congestionMultiplier := 1.2 // 20% increase for medium congestion
	recommendedFeeRate := baseFeeRate * congestionMultiplier
	estimatedFee := estimateValue * recommendedFeeRate
	
	fmt.Printf("\nFee estimation:\n")
	fmt.Printf("  Base fee rate: %.2f%%\n", baseFeeRate*100)
	fmt.Printf("  Congestion multiplier: %.1fx\n", congestionMultiplier)
	fmt.Printf("  Recommended fee rate: %.3f%%\n", recommendedFeeRate*100)
	fmt.Printf("  Estimated fee cost: $%.2f\n", estimatedFee)
	
	// Alternative fee levels
	fmt.Println("\nAlternative fee levels:")
	alternatives := []struct {
		level    string
		rate     float64
		fee      float64
		execTime string
	}{
		{"Low", 0.0001 * congestionMultiplier, estimateValue * 0.0001 * congestionMultiplier, "15-25 sec"},
		{"Medium", recommendedFeeRate, estimatedFee, "5-12 sec"},
		{"High", 0.0010 * congestionMultiplier, estimateValue * 0.0010 * congestionMultiplier, "2-6 sec"},
		{"Ultra", 0.0025 * congestionMultiplier, estimateValue * 0.0025 * congestionMultiplier, "0-3 sec"},
	}
	
	fmt.Printf("%-6s %-8s %-8s %-12s\n", "Level", "Rate %", "Fee $", "Exec Time")
	fmt.Println("------------------------------------")
	
	for _, alt := range alternatives {
		fmt.Printf("%-6s %-7.3f%% $%-7.2f %-12s\n", 
			alt.level, alt.rate*100, alt.fee, alt.execTime)
	}

	// Final summary
	fmt.Println("\n--- Final Summary ---")
	
	fmt.Printf("Account: %s\n", address)
	fmt.Println("Builder fee status: Ready for optimized trading")
	fmt.Println("Available builder fee operations:")
	fmt.Println("  ✓ Query current fee rates")
	fmt.Println("  ✓ Analyze network congestion")
	fmt.Println("  ✓ Calculate optimal fees")
	fmt.Println("  ✓ Set dynamic fee adjustment")
	fmt.Println("  ✓ Monitor fee effectiveness")
	fmt.Println("  ✓ Estimate fees for trades")

	fmt.Println("\nBuilder fee example completed!")
	fmt.Println("Note: This example demonstrated:")
	fmt.Println("1. Understanding builder fee concepts")
	fmt.Println("2. Analyzing current rates and congestion")
	fmt.Println("3. Calculating fees for different trade sizes")
	fmt.Println("4. Setting fees for various order types")
	fmt.Println("5. Dynamic fee adjustment strategies")
	fmt.Println("6. Fee analytics and performance tracking")
	fmt.Println("7. Best practices and optimization")
	fmt.Println("8. Fee estimation tools")
	fmt.Println("\nIMPORTANT: Builder fees are optional but can improve execution!")
	fmt.Println("Consider trade urgency, network conditions, and cost-benefit ratio.")
	fmt.Println("Note: The actual Go SDK implementation will require specific method names")
	fmt.Println("that match the Hyperliquid API endpoints for builder fee functionality.")
}