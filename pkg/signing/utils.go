package signing

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"hyperliquid-go-sdk/pkg/errors"
	"hyperliquid-go-sdk/pkg/types"
)

// FloatToWire converts a float to a wire format string with proper precision
func FloatToWire(x float64) (string, error) {
	rounded := fmt.Sprintf("%.8f", x)

	// Check for rounding errors
	if parsedFloat, err := strconv.ParseFloat(rounded, 64); err != nil {
		return "", err
	} else if math.Abs(parsedFloat-x) >= 1e-12 {
		return "", fmt.Errorf("%w: FloatToWire causes rounding for value %f", errors.ErrFloatPrecision, x)
	}

	// Handle "-0" case
	if rounded == "-0.00000000" {
		rounded = "0.00000000"
	}

	// Use decimal for normalization
	dec, err := decimal.NewFromString(rounded)
	if err != nil {
		return "", err
	}

	return dec.String(), nil
}

// FloatToIntForHashing converts a float to integer for hashing (8 decimal places)
func FloatToIntForHashing(x float64) (int64, error) {
	return FloatToInt(x, 8)
}

// FloatToUsdInt converts a float to USD integer (6 decimal places)
func FloatToUsdInt(x float64) (int64, error) {
	return FloatToInt(x, 6)
}

// FloatToInt converts a float to integer with specified decimal places
func FloatToInt(x float64, power int) (int64, error) {
	withDecimals := x * math.Pow(10, float64(power))
	rounded := math.Round(withDecimals)

	if math.Abs(rounded-withDecimals) >= 1e-3 {
		return 0, fmt.Errorf("%w: FloatToInt causes rounding for value %f", errors.ErrFloatPrecision, x)
	}

	return int64(rounded), nil
}

// OrderTypeToWire converts OrderType to OrderTypeWire
func OrderTypeToWire(orderType types.OrderType) (types.OrderTypeWire, error) {
	var wire types.OrderTypeWire

	if orderType.Limit != nil {
		wire.Limit = orderType.Limit
		return wire, nil
	}

	if orderType.Trigger != nil {
		triggerPxWire, err := FloatToWire(orderType.Trigger.TriggerPx)
		if err != nil {
			return wire, err
		}

		wire.Trigger = &types.TriggerOrderTypeWire{
			TriggerPx: types.FloatString(triggerPxWire),
			IsMarket:  orderType.Trigger.IsMarket,
			Tpsl:      orderType.Trigger.Tpsl,
		}
		return wire, nil
	}

	return wire, fmt.Errorf("invalid order type: must have either limit or trigger")
}

// OrderRequestToOrderWire converts OrderRequest to OrderWire
func OrderRequestToOrderWire(order types.OrderRequest, asset int) (types.OrderWire, error) {
	limitPxWire, err := FloatToWire(order.LimitPx)
	if err != nil {
		return types.OrderWire{}, err
	}

	szWire, err := FloatToWire(order.Sz)
	if err != nil {
		return types.OrderWire{}, err
	}

	orderTypeWire, err := OrderTypeToWire(order.OrderType)
	if err != nil {
		return types.OrderWire{}, err
	}

	wire := types.OrderWire{
		A: asset,
		B: order.IsBuy,
		P: types.FloatString(limitPxWire),
		S: types.FloatString(szWire),
		R: order.ReduceOnly,
		T: orderTypeWire,
	}

	if order.Cloid != nil {
		cloidStr := order.Cloid.Raw()
		wire.C = &cloidStr
	}

	return wire, nil
}

// OrderWiresToOrderAction converts order wires to order action
func OrderWiresToOrderAction(orderWires []types.OrderWire, builder *types.BuilderInfo) types.OrderAction {
	action := types.OrderAction{
		Type:     "order",
		Orders:   orderWires,
		Grouping: "na",
	}

	if builder != nil {
		action.Builder = builder
	}

	return action
}

// AddressToBytes converts hex address string to bytes
func AddressToBytes(address string) []byte {
	// Remove 0x prefix if present
	if strings.HasPrefix(address, "0x") {
		address = address[2:]
	}

	// Convert hex string to bytes
	bytes := make([]byte, len(address)/2)
	for i := 0; i < len(address); i += 2 {
		b, _ := strconv.ParseUint(address[i:i+2], 16, 8)
		bytes[i/2] = byte(b)
	}

	return bytes
}

// ValidateAddress validates if the address is in correct format
func ValidateAddress(address string) error {
	if !strings.HasPrefix(address, "0x") {
		return errors.ErrInvalidAddress
	}

	if len(address) != 42 {
		return errors.ErrInvalidAddress
	}

	// Check if it's valid hex
	if _, err := strconv.ParseInt(address[2:], 16, 64); err != nil {
		return errors.ErrInvalidAddress
	}

	return nil
}
