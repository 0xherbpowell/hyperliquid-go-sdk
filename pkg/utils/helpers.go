package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

// ParsePrivateKey parses a private key from hex string
func ParsePrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	// Remove 0x prefix if present
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}
	
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key hex: %w", err)
	}
	
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	
	return privateKey, nil
}

// GetAddressFromPrivateKey gets the Ethereum address from a private key
func GetAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return ""
	}
	
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address.Hex()
}

// FormatPrice formats a price with appropriate decimal places
func FormatPrice(price float64, decimals int) string {
	format := fmt.Sprintf("%%.%df", decimals)
	formatted := fmt.Sprintf(format, price)
	
	// Remove trailing zeros
	if strings.Contains(formatted, ".") {
		formatted = strings.TrimRight(formatted, "0")
		formatted = strings.TrimRight(formatted, ".")
	}
	
	return formatted
}

// ParsePrice parses a price string to float64
func ParsePrice(priceStr string) (float64, error) {
	return strconv.ParseFloat(priceStr, 64)
}

// FormatSize formats a size with appropriate decimal places
func FormatSize(size float64, decimals int) string {
	return FormatPrice(size, decimals)
}

// ParseSize parses a size string to float64
func ParseSize(sizeStr string) (float64, error) {
	return strconv.ParseFloat(sizeStr, 64)
}

// ValidateAddress validates an Ethereum address
func ValidateAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	
	if len(address) != 42 {
		return false
	}
	
	// Check if it's valid hex
	_, err := hex.DecodeString(address[2:])
	return err == nil
}

// NormalizeAddress normalizes an address to lowercase
func NormalizeAddress(address string) string {
	return strings.ToLower(address)
}

// IsSpotAsset checks if an asset ID represents a spot asset
func IsSpotAsset(asset int) bool {
	return asset >= 10000
}

// IsPerpAsset checks if an asset ID represents a perpetual asset
func IsPerpAsset(asset int) bool {
	return asset < 10000 || asset >= 110000
}

// GetDecimalPower returns 10^n as float64
func GetDecimalPower(n int) float64 {
	return pow10(n)
}

// TruncateFloat truncates a float to a specified number of decimal places
func TruncateFloat(f float64, decimals int) float64 {
	multiplier := pow10(decimals)
	return float64(int64(f*multiplier)) / multiplier
}

// RoundToSignificantFigures rounds a number to a specified number of significant figures
func RoundToSignificantFigures(f float64, sigFigs int) float64 {
	if f == 0 {
		return 0
	}
	
	// Find the magnitude
	magnitude := 0
	absF := abs(f)
	
	if absF >= 1 {
		for absF >= 10 {
			absF /= 10
			magnitude++
		}
	} else {
		for absF < 1 {
			absF *= 10
			magnitude--
		}
	}
	
	// Calculate the rounding factor
	factor := pow10(sigFigs - 1 - magnitude)
	
	// Round and return
	return round(f*factor) / factor
}

// CalculateSlippagePrice calculates price with slippage
func CalculateSlippagePrice(price float64, slippage float64, isBuy bool) float64 {
	if isBuy {
		return price * (1 + slippage)
	}
	return price * (1 - slippage)
}

// IsValidSlippage checks if slippage is within reasonable bounds
func IsValidSlippage(slippage float64) bool {
	return slippage >= 0 && slippage <= 1.0 // 0% to 100%
}

// IsValidLeverage checks if leverage is within reasonable bounds
func IsValidLeverage(leverage int) bool {
	return leverage >= 1 && leverage <= 100 // 1x to 100x
}

// ConvertBasisPointsToDecimal converts basis points to decimal (e.g., 100 bp = 0.01)
func ConvertBasisPointsToDecimal(basisPoints int) float64 {
	return float64(basisPoints) / 10000.0
}

// ConvertDecimalToBasisPoints converts decimal to basis points (e.g., 0.01 = 100 bp)
func ConvertDecimalToBasisPoints(decimal float64) int {
	return int(decimal * 10000)
}

// FormatLeverage formats leverage for display
func FormatLeverage(leverage int, isIsolated bool) string {
	if isIsolated {
		return fmt.Sprintf("%dx (Isolated)", leverage)
	}
	return fmt.Sprintf("%dx (Cross)", leverage)
}

// SanitizeInput sanitizes user input by trimming whitespace and converting to lowercase
func SanitizeInput(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

// ValidateCoinName validates a coin name format
func ValidateCoinName(coin string) bool {
	if coin == "" {
		return false
	}
	
	// Basic validation - should contain only alphanumeric characters and some special chars
	for _, r := range coin {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || 
			 (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '/') {
			return false
		}
	}
	
	return true
}

// GetOrderSideString returns human-readable order side
func GetOrderSideString(isBuy bool) string {
	if isBuy {
		return "BUY"
	}
	return "SELL"
}

// ParseOrderSide parses order side string to boolean
func ParseOrderSide(side string) (bool, error) {
	switch strings.ToUpper(strings.TrimSpace(side)) {
	case "BUY", "B", "1":
		return true, nil
	case "SELL", "S", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid order side: %s", side)
	}
}

// CalculatePnL calculates PnL given entry price, current price, size, and side
func CalculatePnL(entryPrice, currentPrice, size float64, isBuy bool) float64 {
	if isBuy {
		return (currentPrice - entryPrice) * size
	}
	return (entryPrice - currentPrice) * size
}

// CalculateROE calculates return on equity as a percentage
func CalculateROE(pnl, margin float64) float64 {
	if margin == 0 {
		return 0
	}
	return (pnl / margin) * 100
}

// EstimateGasPrice provides a rough estimate for gas price (in Gwei)
// Note: This is a placeholder - in production, you'd fetch from gas oracle
func EstimateGasPrice() int64 {
	return 20 // 20 Gwei default
}

// FormatDuration formats duration in milliseconds to human readable format
func FormatDuration(durationMs int64) string {
	if durationMs < 1000 {
		return fmt.Sprintf("%dms", durationMs)
	}
	
	seconds := durationMs / 1000
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	
	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	
	hours := minutes / 60
	return fmt.Sprintf("%dh", hours)
}

// CreateRandomWallet creates a new random wallet for testing
func CreateRandomWallet() (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return privateKey, nil
}
