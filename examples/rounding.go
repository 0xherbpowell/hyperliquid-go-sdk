package main

import (
	"fmt"
	"log"
	"math"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// This example demonstrates proper rounding and precision handling
	// when working with Hyperliquid orders and calculations

	fmt.Println("Rounding and Precision Example")
	fmt.Printf("Account: %s\n", address)

	// Get initial user state
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("\nInitial account state:")
	PrintPositions(userState)

	// Example 1: Price rounding
	fmt.Println("\n--- Example 1: Price Rounding ---")

	// Get current ETH price
	mids, err := info.AllMids("")
	if err != nil {
		log.Printf("Failed to get mids: %v", err)
		return
	}

	ethMid, exists := mids["ETH"]
	if !exists {
		log.Printf("ETH mid price not found")
		return
	}

	ethPrice, err := utils.ParsePrice(ethMid)
	if err != nil {
		log.Printf("Failed to parse ETH price: %v", err)
		return
	}

	fmt.Printf("Raw ETH price: %.10f\n", ethPrice)

	// Demonstrate different rounding strategies for prices
	priceOffsets := []float64{0.001234567, -0.002345678, 0.000987654}

	for i, offset := range priceOffsets {
		targetPrice := ethPrice * (1 + offset)

		fmt.Printf("\nPrice calculation %d:\n", i+1)
		fmt.Printf("  Target price (unrounded): %.10f\n", targetPrice)

		// Round to different precisions commonly used in trading
		rounded2 := roundToDecimalPlaces(targetPrice, 2)
		rounded4 := roundToDecimalPlaces(targetPrice, 4)
		rounded6 := roundToDecimalPlaces(targetPrice, 6)

		fmt.Printf("  Rounded to 2 decimals: %.2f\n", rounded2)
		fmt.Printf("  Rounded to 4 decimals: %.4f\n", rounded4)
		fmt.Printf("  Rounded to 6 decimals: %.6f\n", rounded6)

		// Show tick size rounding (common in exchanges)
		tickSize := 0.01 // Example tick size
		tickRounded := roundToTickSize(targetPrice, tickSize)
		fmt.Printf("  Rounded to tick size %.2f: %.2f\n", tickSize, tickRounded)
	}

	// Example 2: Size rounding
	fmt.Println("\n--- Example 2: Size Rounding ---")

	baseSizes := []float64{0.123456789, 1.987654321, 10.555555555}

	for i, baseSize := range baseSizes {
		fmt.Printf("\nSize calculation %d:\n", i+1)
		fmt.Printf("  Raw size: %.10f\n", baseSize)

		// Different rounding strategies for order sizes
		roundedDown := math.Floor(baseSize*1000) / 1000    // Floor to 3 decimals
		roundedUp := math.Ceil(baseSize*1000) / 1000       // Ceil to 3 decimals
		roundedNearest := math.Round(baseSize*1000) / 1000 // Round to nearest 3 decimals

		fmt.Printf("  Floor (3 decimals): %.3f\n", roundedDown)
		fmt.Printf("  Ceil (3 decimals): %.3f\n", roundedUp)
		fmt.Printf("  Round (3 decimals): %.3f\n", roundedNearest)

		// Minimum size enforcement
		minSize := 0.001
		enforcedSize := math.Max(roundedNearest, minSize)
		fmt.Printf("  Enforced minimum (%.3f): %.3f\n", minSize, enforcedSize)
	}

	// Example 3: Placing orders with proper rounding
	fmt.Println("\n--- Example 3: Orders with Proper Rounding ---")

	// Calculate order parameters with proper rounding
	rawOrderPrice := ethPrice * 0.995123456789 // 0.5% below market with extra precision
	rawOrderSize := 0.012345678901             // Size with extra precision

	fmt.Printf("Raw order parameters:\n")
	fmt.Printf("  Price: %.10f\n", rawOrderPrice)
	fmt.Printf("  Size: %.10f\n", rawOrderSize)

	// Apply proper rounding for Hyperliquid
	properPrice := roundToDecimalPlaces(rawOrderPrice, 6) // 6 decimal precision for price
	properSize := roundToDecimalPlaces(rawOrderSize, 4)   // 4 decimal precision for size

	fmt.Printf("\nProperly rounded parameters:\n")
	fmt.Printf("  Price: %.6f\n", properPrice)
	fmt.Printf("  Size: %.4f\n", properSize)

	// Place order with rounded values
	fmt.Printf("\nPlacing order with rounded values:\n")

	orderResult, err := exchange.Order(
		"ETH",                 // coin
		true,                  // isBuy
		properSize,            // properly rounded size
		properPrice,           // properly rounded price
		CreateGtcLimitOrder(), // GTC order type
		false,                 // reduce only
		GenerateCloid(),       // unique client order ID
		nil,                   // builder info
	)

	if err != nil {
		log.Printf("Failed to place order: %v", err)
	} else {
		fmt.Println("Order placed successfully:")
		PrintOrderResult(orderResult)
	}

	// Example 4: Portfolio calculations with rounding
	fmt.Println("\n--- Example 4: Portfolio Calculations ---")

	// Simulate portfolio calculations that require careful rounding
	positions := []struct {
		coin  string
		size  float64
		price float64
	}{
		{"ETH", 1.23456789, ethPrice},
		{"BTC", 0.05678901, ethPrice * 20},  // Assuming BTC is ~20x ETH price
		{"SOL", 12.34567890, ethPrice * 0.1}, // Assuming SOL is ~0.1x ETH price
	}

	fmt.Println("Portfolio value calculations:")
	totalValue := 0.0

	for _, pos := range positions {
		rawValue := pos.size * pos.price
		roundedValue := roundToDecimalPlaces(rawValue, 2) // Round to cents

		fmt.Printf("  %s: %.8f × %.4f = %.8f → %.2f\n",
			pos.coin, pos.size, pos.price, rawValue, roundedValue)

		totalValue += roundedValue
	}

	fmt.Printf("Total portfolio value: $%.2f\n", totalValue)

	// Example 5: Percentage calculations
	fmt.Println("\n--- Example 5: Percentage Calculations ---")

	// Calculate percentage changes with proper rounding
	initialPrice := 2000.0
	finalPrices := []float64{2001.23, 1998.76, 2050.45, 1950.12}

	fmt.Printf("Price change calculations (initial: $%.2f):\n", initialPrice)

	for i, finalPrice := range finalPrices {
		rawChange := (finalPrice - initialPrice) / initialPrice * 100
		roundedChange := roundToDecimalPlaces(rawChange, 4) // 4 decimal precision for percentages

		fmt.Printf("  Change %d: $%.2f → %.8f%% → %.4f%%\n",
			i+1, finalPrice, rawChange, roundedChange)
	}

	// Example 6: Common rounding errors to avoid
	fmt.Println("\n--- Example 6: Common Rounding Pitfalls ---")

	fmt.Println("Demonstrating common rounding issues:")

	// Floating point precision issues
	value1 := 0.1 + 0.2
	fmt.Printf("0.1 + 0.2 = %.20f (should be 0.3)\n", value1)
	fmt.Printf("Properly rounded: %.1f\n", roundToDecimalPlaces(value1, 1))

	// Cumulative rounding errors
	sum := 0.0
	for i := 0; i < 10; i++ {
		sum += 0.1
	}
	fmt.Printf("Sum of 0.1 ten times: %.20f (should be 1.0)\n", sum)
	fmt.Printf("Properly rounded: %.1f\n", roundToDecimalPlaces(sum, 1))

	// Banker's rounding vs standard rounding
	testValues := []float64{2.5, 3.5, 4.5, 5.5}
	fmt.Println("\nRounding .5 values:")
	for _, val := range testValues {
		standardRound := math.Round(val)
		fmt.Printf("  %.1f → %.0f (standard rounding)\n", val, standardRound)
	}

	// Cleanup: cancel the order if it's still resting
	if oid, ok := GetRestingOid(orderResult); ok && oid > 0 {
		fmt.Printf("\nCancelling test order (oid: %d)...\n", oid)
		cancelResult, err := exchange.Cancel("ETH", oid)
		if err != nil {
			log.Printf("Failed to cancel order: %v", err)
		} else {
			fmt.Printf("Order cancelled successfully\n")
			_ = cancelResult
		}
	}

	fmt.Println("\nRounding example completed!")
	fmt.Println("Key rounding principles:")
	fmt.Println("1. Always round prices and sizes to appropriate precision")
	fmt.Println("2. Be consistent with rounding methods throughout your application")
	fmt.Println("3. Account for exchange-specific tick sizes and minimum increments")
	fmt.Println("4. Watch out for cumulative rounding errors in calculations")
	fmt.Println("5. Use proper decimal arithmetic for financial calculations when needed")
	fmt.Println("6. Test rounding behavior with edge cases (.5 values, etc.)")
	fmt.Println("\nImportant notes:")
	fmt.Println("• Hyperliquid may reject orders with inappropriate precision")
	fmt.Println("• Always check exchange specifications for precision requirements")
	fmt.Println("• Consider using decimal libraries for critical financial calculations")
}

// Helper function to round to specific decimal places
func roundToDecimalPlaces(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(value*multiplier) / multiplier
}

// Helper function to round to tick size
func roundToTickSize(value float64, tickSize float64) float64 {
	return math.Round(value/tickSize) * tickSize
}

// Helper function to round down to decimal places
func floorToDecimalPlaces(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Floor(value*multiplier) / multiplier
}

// Helper function to round up to decimal places
func ceilToDecimalPlaces(value float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Ceil(value*multiplier) / multiplier
}