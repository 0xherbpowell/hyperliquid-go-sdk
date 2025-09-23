package utils

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/vmihailenco/msgpack/v5"

	"hyperliquid-go-sdk/pkg/types"
)

// EIP712Domain represents the EIP712 domain
type EIP712Domain struct {
	Name              string         `json:"name"`
	Version           string         `json:"version"`
	ChainId           *big.Int       `json:"chainId"`
	VerifyingContract common.Address `json:"verifyingContract"`
}

// PhantomAgent represents a phantom agent
type PhantomAgent struct {
	Source       string      `json:"source"`
	ConnectionId common.Hash `json:"connectionId"`
}

// SignatureResult represents the structured signature result
type SignatureResult struct {
	R string `json:"r"`
	S string `json:"s"`
	V int    `json:"v"`
}

// SignTypes defines the signing type structures
var (
	USDSendSignTypes = []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	}

	SpotTransferSignTypes = []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "token", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	}

	WithdrawSignTypes = []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	}

	USDClassTransferSignTypes = []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "toPerp", Type: "bool"},
		{Name: "nonce", Type: "uint64"},
	}

	SendAssetSignTypes = []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "sourceDex", Type: "string"},
		{Name: "destinationDex", Type: "string"},
		{Name: "token", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "fromSubAccount", Type: "string"},
		{Name: "nonce", Type: "uint64"},
	}

	TokenDelegateTypes = []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "validator", Type: "address"},
		{Name: "wei", Type: "uint64"},
		{Name: "isUndelegate", Type: "bool"},
		{Name: "nonce", Type: "uint64"},
	}

	ConvertToMultiSigUserSignTypes = []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "signers", Type: "string"},
		{Name: "nonce", Type: "uint64"},
	}

	MultiSigEnvelopeSignTypes = []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "multiSigActionHash", Type: "bytes32"},
		{Name: "nonce", Type: "uint64"},
	}
)

// FloatToWire converts a float to wire format string matching Python SDK exactly
func FloatToWire(x float64) (string, error) {
	// Convert to string with 8 decimal places to match Python SDK
	rounded := fmt.Sprintf("%.8f", x)

	// Check for precision loss
	parsed, err := strconv.ParseFloat(rounded, 64)
	if err != nil {
		return "", err
	}

	if abs(parsed-x) >= 1e-12 {
		return "", fmt.Errorf("float_to_wire causes rounding: %f", x)
	}

	// Handle -0 case (must match Python exactly)
	if rounded == "-0.00000000" {
		rounded = "0.00000000"
	}

	// Normalize like Python's Decimal.normalize() - remove trailing zeros
	// Parse as float and format without trailing zeros
	val, err := strconv.ParseFloat(rounded, 64)
	if err != nil {
		return "", err
	}

	// Format without trailing zeros, similar to Python's normalize()
	result := strconv.FormatFloat(val, 'f', -1, 64)

	return result, nil
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// FloatToIntForHashing converts float to int for hashing
func FloatToIntForHashing(x float64) (int64, error) {
	return FloatToInt(x, 8)
}

// FloatToUSDInt converts float to USD int
func FloatToUSDInt(x float64) (int64, error) {
	return FloatToInt(x, 6)
}

// FloatToInt converts float to int with given power
func FloatToInt(x float64, power int) (int64, error) {
	withDecimals := x * pow10(power)
	rounded := round(withDecimals)

	if abs(rounded-withDecimals) >= 1e-3 {
		return 0, fmt.Errorf("float_to_int causes rounding: %f", x)
	}

	return int64(rounded), nil
}

// pow10 returns 10^n
func pow10(n int) float64 {
	result := 1.0
	for i := 0; i < n; i++ {
		result *= 10.0
	}
	return result
}

// round rounds a float64 to the nearest integer
func round(x float64) float64 {
	if x >= 0 {
		return float64(int64(x + 0.5))
	}
	return float64(int64(x - 0.5))
}

// GetTimestampMS returns current timestamp in milliseconds
func GetTimestampMS() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func addressToBytes(address string) []byte {
	address = strings.TrimPrefix(address, "0x")
	bytes, _ := hex.DecodeString(address)
	return bytes
}

// ActionHash computes the hash of an action using same logic as reference SDK
func ActionHash(action interface{}, vaultAddress *string, nonce int64, expiresAfter *int64) []byte {
	// Pack action using msgpack with consistent settings
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	enc.SetSortMapKeys(true)
	enc.UseCompactInts(true)

	err := enc.Encode(action)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal action: %v", err))
	}
	data := buf.Bytes()

	// Add nonce as 8 bytes big endian
	if nonce < 0 {
		panic(fmt.Sprintf("nonce cannot be negative: %d", nonce))
	}
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(nonce))
	data = append(data, nonceBytes...)

	// Add vault address
	if vaultAddress == nil {
		data = append(data, 0x00)
	} else {
		data = append(data, 0x01)
		data = append(data, addressToBytes(*vaultAddress)...)
	}

	// Add expires_after if provided
	if expiresAfter != nil {
		if *expiresAfter < 0 {
			panic(fmt.Sprintf("expiresAfter cannot be negative: %d", *expiresAfter))
		}
		data = append(data, 0x00)
		expiresAfterBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(expiresAfterBytes, uint64(*expiresAfter))
		data = append(data, expiresAfterBytes...)
	}

	// Return keccak256 hash
	hash := crypto.Keccak256(data)
	// fmt.Printf("go action hash: %s\n", hex.EncodeToString(hash))
	return hash
}

// ConstructPhantomAgent creates a phantom agent from hash
func ConstructPhantomAgent(hash []byte, isMainnet bool) map[string]interface{} {
	source := "b"
	if isMainnet {
		source = "a"
	}

	return map[string]interface{}{
		"source":       source,
		"connectionId": hash,
	}
}

// L1Payload constructs the EIP712 payload for L1 actions using same logic as reference SDK
func L1Payload(phantomAgent map[string]interface{}) apitypes.TypedData {
	// Fix: Use direct cast instead of dereferencing to avoid conversion issues
	chainIdValue := big.NewInt(EIP712ChainID)
	chainId := (*math.HexOrDecimal256)(chainIdValue)
	return apitypes.TypedData{
		Domain: apitypes.TypedDataDomain{
			ChainId:           chainId,
			Name:              "Exchange",
			Version:           "1",
			VerifyingContract: "0x0000000000000000000000000000000000000000",
		},
		Types: apitypes.Types{
			"Agent": []apitypes.Type{
				{Name: "source", Type: "string"},
				{Name: "connectionId", Type: "bytes32"},
			},
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
		},
		PrimaryType: "Agent",
		Message:     phantomAgent,
	}
}

// UserSignedPayload constructs the EIP712 payload for user signed actions
func UserSignedPayload(primaryType string, payloadTypes []apitypes.Type, action map[string]interface{}) apitypes.TypedData {
	chainIdStr, ok := action["signatureChainId"].(string)
	if !ok {
		chainIdStr = SignatureChainID
	}

	chainId, _ := big.NewInt(0).SetString(chainIdStr, 0)

	types := apitypes.Types{
		"EIP712Domain": []apitypes.Type{
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
	}
	types[primaryType] = payloadTypes

	// Create field map for precise filtering
	validFields := make(map[string]bool)
	for _, fieldType := range payloadTypes {
		validFields[fieldType.Name] = true
	}

	// Create message with only schema fields (exclude signatureChainId which is used for domain)
	message := make(map[string]interface{})
	for k, v := range action {
		if k != "signatureChainId" && validFields[k] {
			message[k] = v
		}
	}

	return apitypes.TypedData{
		Types:       types,
		PrimaryType: primaryType,
		Domain: apitypes.TypedDataDomain{
			Name:              "HyperliquidSignTransaction",
			Version:           "1",
			ChainId:           (*math.HexOrDecimal256)(chainId),
			VerifyingContract: "0x0000000000000000000000000000000000000000",
		},
		Message: message,
	}
}

func SignL1Action(
	privateKey *ecdsa.PrivateKey,
	action any,
	vaultAddress *string,
	timestamp int64,
	expiresAfter *int64,
	isMainnet bool,
) (SignatureResult, error) {

	hash := ActionHash(action, vaultAddress, timestamp, expiresAfter)

	phantomAgent := ConstructPhantomAgent(hash, isMainnet)

	typedData := L1Payload(phantomAgent)

	return SignInner(privateKey, typedData)
}

//// SignL1ActionWithAccount signs an L1 action with optional account address for agent trading
//// Returns map[string]interface{} for compatibility with existing exchange code
//func SignL1ActionWithAccount(privateKey *ecdsa.PrivateKey, action interface{}, activePool *string, nonce int64, expiresAfter *int64, isMainnet bool) (map[string]interface{}, error) {
//	// For agent trading, the signature is the same, but the system relies on
//	// the agent being pre-authorized to act on behalf of the account
//	// The account address is handled at the exchange/API level, not in the signature
//
//	vaultAddress := ""
//	if activePool != nil {
//		vaultAddress = *activePool
//	}
//
//	sig, err := SignL1Action(privateKey, action, &vaultAddress, nonce, expiresAfter, isMainnet)
//	if err != nil {
//		return nil, err
//	}
//
//	// Convert SignatureResult to map for compatibility
//	return map[string]interface{}{
//		"r": sig.R,
//		"s": sig.S,
//		"v": sig.V,
//	}, nil
//}

// SignUserSignedAction signs a user signed action
func SignUserSignedAction(privateKey *ecdsa.PrivateKey, action map[string]interface{}, payloadTypes []apitypes.Type, primaryType string, isMainnet bool) (map[string]interface{}, error) {
	// Make a copy of the action to avoid modifying the original
	signAction := make(map[string]interface{})
	for k, v := range action {
		signAction[k] = v
	}

	// Add required fields
	signAction["signatureChainId"] = SignatureChainID
	if isMainnet {
		signAction["hyperliquidChain"] = MainnetChainName
	} else {
		signAction["hyperliquidChain"] = TestnetChainName
	}

	data := UserSignedPayload(primaryType, payloadTypes, signAction)
	sig, err := SignInner(privateKey, data)
	if err != nil {
		return nil, err
	}

	// Convert SignatureResult to map for compatibility
	return map[string]interface{}{
		"r": sig.R,
		"s": sig.S,
		"v": sig.V,
	}, nil
}

func SignInner(privateKey *ecdsa.PrivateKey, typedData apitypes.TypedData) (SignatureResult, error) {

	// Create EIP-712 hash
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return SignatureResult{}, fmt.Errorf("failed to hash domain: %w", err)
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return SignatureResult{}, fmt.Errorf("failed to hash typed data: %w", err)
	}
	rawData := []byte{0x19, 0x01}
	rawData = append(rawData, domainSeparator...)
	rawData = append(rawData, typedDataHash...)

	msgHash := crypto.Keccak256Hash(rawData)

	signature, err := crypto.Sign(msgHash.Bytes(), privateKey)
	if err != nil {
		return SignatureResult{}, fmt.Errorf("failed to sign message: %w", err)
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])
	v := int(signature[64]) + 27

	result := SignatureResult{
		R: hexutil.EncodeBig(r),
		S: hexutil.EncodeBig(s),
		V: v,
	}

	return result, nil
}

// OrderTypeToWire converts OrderType to wire format
func OrderTypeToWire(orderType types.OrderType) (types.OrderTypeWire, error) {
	var wire types.OrderTypeWire

	if orderType.Limit != nil {
		wire.Limit = orderType.Limit
	} else if orderType.Trigger != nil {
		triggerPxWire, err := FloatToWire(orderType.Trigger.TriggerPx)
		if err != nil {
			return wire, err
		}

		wire.Trigger = &types.TriggerOrderTypeWire{
			IsMarket:  orderType.Trigger.IsMarket,
			TriggerPx: triggerPxWire,
			Tpsl:      orderType.Trigger.Tpsl,
		}
	} else {
		return wire, fmt.Errorf("invalid order type")
	}

	return wire, nil
}

// OrderRequestToOrderWire converts OrderRequest to wire format
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

	// Create OrderWire with fields in the expected order
	wire := types.OrderWire{
		A: asset,            // asset ID
		B: order.IsBuy,      // is buy order
		P: limitPxWire,      // limit price
		S: szWire,           // size
		R: order.ReduceOnly, // reduce only
		T: orderTypeWire,    // order type
	}

	if order.Cloid != nil {
		cloidStr := order.Cloid.ToRaw()
		wire.C = &cloidStr
	}

	return wire, nil
}

// ConvertOrderTypeWireToMap converts OrderTypeWire to map format for JSON serialization
func ConvertOrderTypeWireToMap(orderType types.OrderTypeWire) map[string]interface{} {
	if orderType.Limit != nil {
		return map[string]interface{}{
			"limit": map[string]interface{}{
				"tif": string(orderType.Limit.Tif),
			},
		}
	} else if orderType.Trigger != nil {
		return map[string]interface{}{
			"trigger": map[string]interface{}{
				"triggerPx": orderType.Trigger.TriggerPx, // Already a string
				"isMarket":  orderType.Trigger.IsMarket,
				"tpsl":      string(orderType.Trigger.Tpsl),
			},
		}
	}
	return map[string]interface{}{}
}

// OrderWiresToOrderAction converts order wires to order action
func OrderWiresToOrderAction(orderWires []types.OrderWire, builder *types.BuilderInfo) map[string]interface{} {
	// Convert OrderWires to maps to ensure proper JSON serialization
	// This matches the TypeScript SDK format exactly
	orderMaps := make([]map[string]interface{}, len(orderWires))
	for i, wire := range orderWires {
		orderMap := map[string]interface{}{
			"a": wire.A,                            // asset (number)
			"b": wire.B,                            // isBuy (boolean)
			"p": wire.P,                            // price (string)
			"s": wire.S,                            // size (string)
			"r": wire.R,                            // reduceOnly (boolean)
			"t": ConvertOrderTypeWireToMap(wire.T), // orderType (object)
		}

		// Add cloid if present
		if wire.C != nil {
			orderMap["c"] = *wire.C
		}

		orderMaps[i] = orderMap
	}

	// Create action with proper structure matching TypeScript SDK
	action := map[string]interface{}{
		"type":     "order",
		"orders":   orderMaps, // Now using maps instead of structs
		"grouping": "na",
	}

	if builder != nil {
		action["builder"] = map[string]interface{}{
			"b": builder.B,
			"f": builder.F,
		}
	}

	return action
}

// Sign functions for different action types

// SignUSDTransferAction signs a USD transfer action
func SignUSDTransferAction(privateKey *ecdsa.PrivateKey, action map[string]interface{}, isMainnet bool) (map[string]interface{}, error) {
	// Create a copy of the action for signing with proper time field handling
	signAction := make(map[string]interface{})
	for k, v := range action {
		// Skip the 'type' field - it's not part of the EIP712 schema
		if k == "type" {
			continue
		}
		if k == "time" {
			// Convert time string to uint64 for EIP712 signing
			// EIP712 uint64 values need to be provided as *big.Int
			if timeStr, ok := v.(string); ok {
				if timestamp, err := strconv.ParseUint(timeStr, 10, 64); err == nil {
					signAction[k] = new(big.Int).SetUint64(timestamp)
				} else {
					signAction[k] = v
				}
			} else if timestamp, ok := v.(int64); ok {
				if timestamp >= 0 {
					signAction[k] = new(big.Int).SetUint64(uint64(timestamp))
				} else {
					signAction[k] = v
				}
			} else if timestamp, ok := v.(uint64); ok {
				signAction[k] = new(big.Int).SetUint64(timestamp)
			} else {
				signAction[k] = v
			}
		} else {
			signAction[k] = v
		}
	}

	return SignUserSignedAction(privateKey, signAction, USDSendSignTypes, "HyperliquidTransaction:UsdSend", isMainnet)
}

// SignSpotTransferAction signs a spot transfer action
func SignSpotTransferAction(privateKey *ecdsa.PrivateKey, action map[string]interface{}, isMainnet bool) (map[string]interface{}, error) {
	return SignUserSignedAction(privateKey, action, SpotTransferSignTypes, "HyperliquidTransaction:SpotSend", isMainnet)
}

// SignWithdrawFromBridgeAction signs a withdraw from bridge action
func SignWithdrawFromBridgeAction(privateKey *ecdsa.PrivateKey, action map[string]interface{}, isMainnet bool) (map[string]interface{}, error) {
	return SignUserSignedAction(privateKey, action, WithdrawSignTypes, "HyperliquidTransaction:Withdraw", isMainnet)
}

// SignUSDClassTransferAction signs a USD class transfer action
func SignUSDClassTransferAction(privateKey *ecdsa.PrivateKey, action map[string]interface{}, isMainnet bool) (map[string]interface{}, error) {
	return SignUserSignedAction(privateKey, action, USDClassTransferSignTypes, "HyperliquidTransaction:UsdClassTransfer", isMainnet)
}

// SignSendAssetAction signs a send asset action
func SignSendAssetAction(privateKey *ecdsa.PrivateKey, action map[string]interface{}, isMainnet bool) (map[string]interface{}, error) {
	return SignUserSignedAction(privateKey, action, SendAssetSignTypes, "HyperliquidTransaction:SendAsset", isMainnet)
}

// SignConvertToMultiSigUserAction signs a convert to multi-sig user action
func SignConvertToMultiSigUserAction(privateKey *ecdsa.PrivateKey, action map[string]interface{}, isMainnet bool) (map[string]interface{}, error) {
	return SignUserSignedAction(privateKey, action, ConvertToMultiSigUserSignTypes, "HyperliquidTransaction:ConvertToMultiSigUser", isMainnet)
}

// SignAgent signs an agent action
func SignAgent(privateKey *ecdsa.PrivateKey, action map[string]interface{}, isMainnet bool) (map[string]interface{}, error) {
	agentSignTypes := []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "agentAddress", Type: "address"},
		{Name: "agentName", Type: "string"},
		{Name: "nonce", Type: "uint64"},
	}

	// Create a copy of the action for signing with proper type conversions
	signAction := make(map[string]interface{})
	for k, v := range action {
		if k == "nonce" {
			// Convert nonce to *big.Int for EIP712 signing
			if nonceStr, ok := v.(string); ok {
				if nonce, err := strconv.ParseUint(nonceStr, 10, 64); err == nil {
					signAction[k] = new(big.Int).SetUint64(nonce)
				} else {
					signAction[k] = v
				}
			} else if nonce, ok := v.(int64); ok {
				if nonce >= 0 {
					signAction[k] = new(big.Int).SetUint64(uint64(nonce))
				} else {
					signAction[k] = v
				}
			} else if nonce, ok := v.(uint64); ok {
				signAction[k] = new(big.Int).SetUint64(nonce)
			} else {
				signAction[k] = v
			}
		} else {
			signAction[k] = v
		}
	}

	return SignUserSignedAction(privateKey, signAction, agentSignTypes, "HyperliquidTransaction:ApproveAgent", isMainnet)
}

// SignApproveBuilderFee signs an approve builder fee action
func SignApproveBuilderFee(privateKey *ecdsa.PrivateKey, action map[string]interface{}, isMainnet bool) (map[string]interface{}, error) {
	builderFeeSignTypes := []apitypes.Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "maxFeeRate", Type: "string"},
		{Name: "builder", Type: "address"},
		{Name: "nonce", Type: "uint64"},
	}
	return SignUserSignedAction(privateKey, action, builderFeeSignTypes, "HyperliquidTransaction:ApproveBuilderFee", isMainnet)
}

// SignTokenDelegateAction signs a token delegate action
func SignTokenDelegateAction(privateKey *ecdsa.PrivateKey, action map[string]interface{}, isMainnet bool) (map[string]interface{}, error) {
	return SignUserSignedAction(privateKey, action, TokenDelegateTypes, "HyperliquidTransaction:TokenDelegate", isMainnet)
}
