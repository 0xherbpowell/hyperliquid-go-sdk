package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// FormatFloat formats a float with the specified precision
func FormatFloat(value float64, precision int) string {
	return strconv.FormatFloat(value, 'f', precision, 64)
}

// ParseFloat parses a string to float64 with error handling
func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ParseFloatWithDefault parses a string to float64 with a default value
func ParseFloatWithDefault(s string, defaultValue float64) float64 {
	if f, err := ParseFloat(s); err == nil {
		return f
	}
	return defaultValue
}

// ParseInt parses a string to int64 with error handling
func ParseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// ParseIntWithDefault parses a string to int64 with a default value
func ParseIntWithDefault(s string, defaultValue int64) int64 {
	if i, err := ParseInt(s); err == nil {
		return i
	}
	return defaultValue
}

// RoundToSignificantFigures rounds a number to the specified significant figures
func RoundToSignificantFigures(value float64, sigFigs int) float64 {
	if value == 0 {
		return 0
	}

	power := math.Pow(10, float64(sigFigs)-math.Ceil(math.Log10(math.Abs(value))))
	return math.Round(value*power) / power
}

// RoundToDecimals rounds a number to the specified decimal places
func RoundToDecimals(value float64, decimals int) float64 {
	power := math.Pow(10, float64(decimals))
	return math.Round(value*power) / power
}

// IsValidPrice validates if a price follows Hyperliquid pricing rules
func IsValidPrice(price float64, szDecimals int, isSpot bool) bool {
	maxDecimals := 6 // For perps
	if isSpot {
		maxDecimals = 8 // For spot
	}

	// Check if it's an integer (integers are always allowed)
	if price == math.Floor(price) {
		return true
	}

	// Check significant figures (max 5)
	sigFigs := CountSignificantFigures(price)
	if sigFigs > 5 {
		return false
	}

	// Check decimal places
	decimalPlaces := CountDecimalPlaces(price)
	maxAllowedDecimals := maxDecimals - szDecimals
	if decimalPlaces > maxAllowedDecimals {
		return false
	}

	return true
}

// CountSignificantFigures counts the number of significant figures in a float
func CountSignificantFigures(value float64) int {
	if value == 0 {
		return 1
	}

	str := strconv.FormatFloat(math.Abs(value), 'g', -1, 64)
	str = strings.Replace(str, ".", "", 1)
	str = strings.Replace(str, "e", "", 1)
	str = strings.Replace(str, "+", "", 1)
	str = strings.Replace(str, "-", "", 1)

	// Remove leading zeros
	str = strings.TrimLeft(str, "0")

	return len(str)
}

// CountDecimalPlaces counts the number of decimal places in a float
func CountDecimalPlaces(value float64) int {
	str := strconv.FormatFloat(value, 'f', -1, 64)
	if !strings.Contains(str, ".") {
		return 0
	}

	parts := strings.Split(str, ".")
	if len(parts) != 2 {
		return 0
	}

	return len(parts[1])
}

// NormalizeDecimal normalizes a decimal using shopspring/decimal for precision
func NormalizeDecimal(value string) (string, error) {
	dec, err := decimal.NewFromString(value)
	if err != nil {
		return "", err
	}
	return dec.String(), nil
}

// ConvertToWei converts a float amount to wei (with specified decimal places)
func ConvertToWei(amount float64, decimals int) (int64, error) {
	multiplier := math.Pow(10, float64(decimals))
	wei := amount * multiplier

	// Check for precision loss
	if math.Abs(math.Round(wei)-wei) >= 1e-3 {
		return 0, fmt.Errorf("precision loss when converting %f to wei with %d decimals", amount, decimals)
	}

	return int64(math.Round(wei)), nil
}

// ConvertFromWei converts wei to float amount (with specified decimal places)
func ConvertFromWei(wei int64, decimals int) float64 {
	divisor := math.Pow(10, float64(decimals))
	return float64(wei) / divisor
}

// GetCurrentTimestamp returns the current Unix timestamp in milliseconds
func GetCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// IsZeroOrEmpty checks if a string represents zero or is empty
func IsZeroOrEmpty(s string) bool {
	if s == "" {
		return true
	}

	value, err := ParseFloat(s)
	if err != nil {
		return true
	}

	return value == 0
}

// CompareFloatStrings compares two float strings numerically
func CompareFloatStrings(a, b string) int {
	valA, errA := ParseFloat(a)
	valB, errB := ParseFloat(b)

	if errA != nil && errB != nil {
		return strings.Compare(a, b)
	}
	if errA != nil {
		return -1
	}
	if errB != nil {
		return 1
	}

	if valA < valB {
		return -1
	} else if valA > valB {
		return 1
	}
	return 0
}

// FormatPercentage formats a decimal as a percentage
func FormatPercentage(value float64, decimals int) string {
	percentage := value * 100
	return fmt.Sprintf("%."+strconv.Itoa(decimals)+"f%%", percentage)
}

// ClampFloat clamps a float value between min and max
func ClampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Contains checks if a slice contains a specific value
func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Filter filters a slice based on a predicate function
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map applies a function to each element of a slice
func Map[T, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// Reduce reduces a slice to a single value using an accumulator function
func Reduce[T, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// UniqueStrings returns unique strings from a slice
func UniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// MergeMaps merges multiple maps into a single map
func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// DeepCopy performs a deep copy of a struct using JSON marshaling/unmarshaling
func DeepCopy(src, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

// IsNil safely checks if an interface is nil
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

// CoalesceString returns the first non-empty string
func CoalesceString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// CoalesceFloat returns the first non-zero float
func CoalesceFloat(values ...float64) float64 {
	for _, value := range values {
		if value != 0 {
			return value
		}
	}
	return 0
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// FormatDuration formats a duration in a human-readable format
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	} else {
		return fmt.Sprintf("%.1fd", d.Hours()/24)
	}
}

// BatchProcess processes items in batches
func BatchProcess[T any](items []T, batchSize int, processor func([]T) error) error {
	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		if err := processor(batch); err != nil {
			return fmt.Errorf("error processing batch %d-%d: %w", i, end-1, err)
		}
	}

	return nil
}

// Retry executes a function with retry logic
func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}

		if i == attempts-1 {
			break
		}

		time.Sleep(delay)
	}

	return fmt.Errorf("failed after %d attempts: %w", attempts, err)
}

// SafeDivide performs division with zero check
func SafeDivide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

// CalculatePercentageChange calculates percentage change between two values
func CalculatePercentageChange(oldValue, newValue float64) float64 {
	if oldValue == 0 {
		return 0
	}
	return ((newValue - oldValue) / oldValue) * 100
}

// FormatBytes formats bytes in a human-readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// ValidateEmail performs basic email validation
func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// SanitizeString removes potentially dangerous characters from a string
func SanitizeString(s string) string {
	// Remove null bytes and control characters
	s = strings.ReplaceAll(s, "\x00", "")
	s = strings.Map(func(r rune) rune {
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			return -1
		}
		return r
	}, s)

	return strings.TrimSpace(s)
}
