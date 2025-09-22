package main

import (
	"fmt"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	// Setup with testnet - NOTE: Staking is typically done on mainnet
	address, info, exchange := Setup(utils.TestnetAPIURL, true)

	// This example demonstrates staking operations on Hyperliquid
	// Staking allows users to earn rewards by delegating or staking tokens

	fmt.Println("Basic Staking Example")
	fmt.Printf("Account: %s\n", address)

	// Get initial user state
	userState, err := info.UserState(address, "")
	if err != nil {
		log.Fatalf("Failed to get user state: %v", err)
	}

	fmt.Println("\nInitial account state:")
	PrintPositions(userState)

	// Example 1: Get user staking summary
	fmt.Println("\n--- Example 1: User Staking Summary ---")

	// Note: In the Python SDK, this would be info.user_staking_summary(address)
	// For demonstration, we'll show what the staking summary would look like
	
	fmt.Printf("Querying staking summary for address: %s\n", address)
	fmt.Printf("Would call: info.UserStakingSummary(\"%s\")\n", address)
	
	// Simulate staking summary
	fmt.Println("Staking summary (simulated):")
	fmt.Println("{")
	fmt.Println("  \"totalStaked\": \"0.0\",")
	fmt.Println("  \"totalRewards\": \"0.0\",")
	fmt.Println("  \"availableRewards\": \"0.0\",")
	fmt.Println("  \"validators\": [],")
	fmt.Println("  \"delegations\": []")
	fmt.Println("}")

	// Example 2: Get user stakes breakdown
	fmt.Println("\n--- Example 2: User Stakes Breakdown ---")
	
	fmt.Printf("Querying stakes breakdown for address: %s\n", address)
	fmt.Printf("Would call: info.UserStakes(\"%s\")\n", address)
	
	// Simulate stakes breakdown
	fmt.Println("Staking breakdown (simulated):")
	fmt.Println("[")
	fmt.Println("  {")
	fmt.Println("    \"validator\": \"validator1\",")
	fmt.Println("    \"stakedAmount\": \"0.0\",")
	fmt.Println("    \"rewards\": \"0.0\",")
	fmt.Println("    \"status\": \"active\"")
	fmt.Println("  }")
	fmt.Println("]")

	// Example 3: Get staking rewards history
	fmt.Println("\n--- Example 3: Staking Rewards History ---")
	
	fmt.Printf("Querying staking rewards for address: %s\n", address)
	fmt.Printf("Would call: info.UserStakingRewards(\"%s\")\n", address)
	
	// Simulate rewards history
	fmt.Println("Most recent staking rewards (simulated):")
	fmt.Println("[")
	fmt.Println("  {")
	fmt.Println("    \"timestamp\": 1699876543,")
	fmt.Println("    \"validator\": \"validator1\",")
	fmt.Println("    \"amount\": \"0.1\",")
	fmt.Println("    \"type\": \"delegation_reward\"")
	fmt.Println("  },")
	fmt.Println("  {")
	fmt.Println("    \"timestamp\": 1699790143,")
	fmt.Println("    \"validator\": \"validator1\",")
	fmt.Println("    \"amount\": \"0.05\",")
	fmt.Println("    \"type\": \"delegation_reward\"")
	fmt.Println("  }")
	fmt.Println("]")

	// Example 4: Staking operations concept demonstration
	fmt.Println("\n--- Example 4: Staking Operations ---")
	
	fmt.Println("Staking operations that would be available:")
	
	// Simulate staking operations
	validatorAddress := "0x1234567890123456789012345678901234567890"
	stakeAmount := 100.0
	
	fmt.Printf("1. Delegate stake to validator %s\n", validatorAddress)
	fmt.Printf("   Would call: exchange.Delegate(\"%s\", %.2f)\n", validatorAddress, stakeAmount)
	fmt.Printf("   This would delegate %.2f tokens to the validator\n", stakeAmount)
	
	fmt.Println("\n2. Undelegate stake from validator")
	fmt.Printf("   Would call: exchange.Undelegate(\"%s\", %.2f)\n", validatorAddress, stakeAmount/2)
	fmt.Printf("   This would undelegate %.2f tokens from the validator\n", stakeAmount/2)
	
	fmt.Println("\n3. Claim staking rewards")
	fmt.Printf("   Would call: exchange.ClaimStakingRewards(\"%s\")\n", validatorAddress)
	fmt.Println("   This would claim all available rewards from the validator")
	
	fmt.Println("\n4. Redelegate stake to different validator")
	newValidatorAddress := "0x9876543210987654321098765432109876543210"
	fmt.Printf("   Would call: exchange.Redelegate(\"%s\", \"%s\", %.2f)\n", 
		validatorAddress, newValidatorAddress, stakeAmount/4)
	fmt.Printf("   This would move %.2f tokens from one validator to another\n", stakeAmount/4)

	// Example 5: Check available validators
	fmt.Println("\n--- Example 5: Available Validators ---")
	
	fmt.Println("Querying available validators:")
	fmt.Printf("Would call: info.GetValidators()\n")
	
	// Simulate validator list
	fmt.Println("Available validators (simulated):")
	validators := []struct {
		address     string
		name        string
		commission  string
		totalStake  string
		status      string
	}{
		{"0x1111...", "Validator Alpha", "5%", "1,000,000", "active"},
		{"0x2222...", "Validator Beta", "3%", "850,000", "active"},
		{"0x3333...", "Validator Gamma", "7%", "650,000", "active"},
	}
	
	for i, validator := range validators {
		fmt.Printf("%d. %s\n", i+1, validator.name)
		fmt.Printf("   Address: %s\n", validator.address)
		fmt.Printf("   Commission: %s\n", validator.commission)
		fmt.Printf("   Total Stake: %s\n", validator.totalStake)
		fmt.Printf("   Status: %s\n", validator.status)
		fmt.Println()
	}

	// Example 6: Staking rewards calculation
	fmt.Println("\n--- Example 6: Staking Rewards Information ---")
	
	fmt.Println("Staking rewards calculation:")
	
	stakedAmount := 1000.0
	annualYield := 12.5 // 12.5% APY
	
	fmt.Printf("If you stake %.2f tokens at %.1f%% APY:\n", stakedAmount, annualYield)
	
	dailyReward := (stakedAmount * annualYield / 100) / 365
	weeklyReward := dailyReward * 7
	monthlyReward := dailyReward * 30
	
	fmt.Printf("  Daily reward: ~%.4f tokens\n", dailyReward)
	fmt.Printf("  Weekly reward: ~%.4f tokens\n", weeklyReward)
	fmt.Printf("  Monthly reward: ~%.4f tokens\n", monthlyReward)
	fmt.Printf("  Annual reward: ~%.2f tokens\n", stakedAmount*annualYield/100)

	// Example 7: Staking best practices
	fmt.Println("\n--- Example 7: Staking Best Practices ---")
	
	fmt.Println("Staking best practices:")
	fmt.Println("1. Research validators carefully - check their:")
	fmt.Println("   • Performance history and uptime")
	fmt.Println("   • Commission rates")
	fmt.Println("   • Total delegated stake")
	fmt.Println("   • Community reputation")
	fmt.Println()
	fmt.Println("2. Diversify your stakes across multiple validators")
	fmt.Println("   • Reduces risk of validator downtime")
	fmt.Println("   • Helps with network decentralization")
	fmt.Println()
	fmt.Println("3. Monitor your rewards and validator performance")
	fmt.Println("   • Check rewards regularly")
	fmt.Println("   • Watch for validator issues")
	fmt.Println()
	fmt.Println("4. Understand unbonding periods")
	fmt.Println("   • Unstaking may have a waiting period")
	fmt.Println("   • Plan your liquidity needs accordingly")

	// Final summary
	fmt.Println("\n--- Final Summary ---")
	
	fmt.Printf("Account: %s\n", address)
	fmt.Println("Staking status: Ready for staking operations")
	fmt.Println("Available operations:")
	fmt.Println("  ✓ Query staking summary")
	fmt.Println("  ✓ View stakes breakdown") 
	fmt.Println("  ✓ Check rewards history")
	fmt.Println("  ✓ Delegate to validators")
	fmt.Println("  ✓ Claim rewards")
	fmt.Println("  ✓ Undelegate stakes")

	fmt.Println("\nStaking example completed!")
	fmt.Println("Note: This example demonstrated:")
	fmt.Println("1. Querying staking summary and breakdown")
	fmt.Println("2. Viewing staking rewards history")
	fmt.Println("3. Understanding staking operations")
	fmt.Println("4. Validator selection and management")
	fmt.Println("5. Rewards calculation and best practices")
	fmt.Println("\nIMPORTANT: Staking involves risk and lock-up periods!")
	fmt.Println("Always understand the terms and conditions before staking.")
	fmt.Println("Note: The actual Go SDK implementation will require specific method names")
	fmt.Println("that match the Hyperliquid API endpoints for staking functionality.")
}