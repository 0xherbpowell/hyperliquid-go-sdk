package utils

import (
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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

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

// FloatToWire converts a float to wire format string
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

	// Remove trailing zeros (must match Python behavior exactly)
	rounded = strings.TrimRight(rounded, "0")
	rounded = strings.TrimRight(rounded, ".")

	// Ensure we never return an empty string
	if rounded == "" {
		rounded = "0"
	}

	return rounded, nil
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
	return float64(int64(x + 0.5))
}

// GetTimestampMS returns current timestamp in milliseconds
func GetTimestampMS() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// AddressToBytes converts hex address string to bytes
func AddressToBytes(address string) ([]byte, error) {
	// Always lowercase the address as per Hyperliquid docs
	address = strings.ToLower(address)
	if strings.HasPrefix(address, "0x") {
		address = address[2:]
	}
	return hex.DecodeString(address)
}

// ActionHash computes the hash of an action
func ActionHash(action interface{}, vaultAddress *string, nonce int64, expiresAfter *int64) ([]byte, error) {
	// Pack action with msgpack
	data, err := msgpack.Marshal(action)
	if err != nil {
		return nil, err
	}

	// Add nonce (8 bytes big endian)
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(nonce))
	data = append(data, nonceBytes...)

	// Add vault address
	if vaultAddress == nil {
		data = append(data, 0x00)
	} else {
		data = append(data, 0x01)
		vaultBytes, err := AddressToBytes(*vaultAddress)
		if err != nil {
			return nil, err
		}
		data = append(data, vaultBytes...)
	}

	// Add expires after
	if expiresAfter != nil {
		data = append(data, 0x01)
		expiresBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(expiresBytes, uint64(*expiresAfter))
		data = append(data, expiresBytes...)
	} else {
		data = append(data, 0x00)
	}

	// Return keccak hash
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)
	return hasher.Sum(nil), nil
}

// ConstructPhantomAgent creates a phantom agent from hash
func ConstructPhantomAgent(hash []byte, isMainnet bool) PhantomAgent {
	source := "b"
	if isMainnet {
		source = "a"
	}

	var connectionId common.Hash
	copy(connectionId[:], hash)

	return PhantomAgent{
		Source:       source,
		ConnectionId: connectionId,
	}
}

// L1Payload constructs the EIP712 payload for L1 actions
func L1Payload(phantomAgent PhantomAgent) apitypes.TypedData {
	return apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Agent": []apitypes.Type{
				{Name: "source", Type: "string"},
				{Name: "connectionId", Type: "bytes32"},
			},
		},
		PrimaryType: "Agent",
		Domain: apitypes.TypedDataDomain{
			Name:              "Exchange",
			Version:           "1",
			ChainId:           (*math.HexOrDecimal256)(big.NewInt(EIP712ChainID)),
			VerifyingContract: "0x0000000000000000000000000000000000000000",
		},
		Message: apitypes.TypedDataMessage{
			"source":       phantomAgent.Source,
			"connectionId": phantomAgent.ConnectionId,
		},
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

// SignL1Action signs an L1 action
func SignL1Action(privateKey *ecdsa.PrivateKey, action interface{}, activePool *string, nonce int64, expiresAfter *int64, isMainnet bool) (map[string]interface{}, error) {
	return SignL1ActionWithAccount(privateKey, action, activePool, nonce, expiresAfter, isMainnet, nil)
}

// SignL1ActionWithAccount signs an L1 action with optional account address for agent trading
func SignL1ActionWithAccount(privateKey *ecdsa.PrivateKey, action interface{}, activePool *string, nonce int64, expiresAfter *int64, isMainnet bool, accountAddress *string) (map[string]interface{}, error) {
	// For agent trading, the signature is the same, but the system relies on
	// the agent being pre-authorized to act on behalf of the account
	// The account address is handled at the exchange/API level, not in the signature

	hash, err := ActionHash(action, activePool, nonce, expiresAfter)
	if err != nil {
		return nil, err
	}

	phantomAgent := ConstructPhantomAgent(hash, isMainnet)
	data := L1Payload(phantomAgent)

	return SignInner(privateKey, data)
}

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
	return SignInner(privateKey, data)
}

// SignInner performs the actual signing
func SignInner(privateKey *ecdsa.PrivateKey, data apitypes.TypedData) (map[string]interface{}, error) {
	// Hash the typed data
	hash, err := data.HashStruct(data.PrimaryType, data.Message)
	if err != nil {
		return nil, err
	}

	// Create the final hash with domain separator
	domainSeparator, err := data.HashStruct("EIP712Domain", data.Domain.Map())
	if err != nil {
		return nil, err
	}

	finalHash := crypto.Keccak256(
		[]byte("\x19\x01"),
		domainSeparator,
		hash,
	)

	// Sign the hash
	signature, err := crypto.Sign(finalHash, privateKey)
	if err != nil {
		return nil, err
	}

	// Parse signature components
	r := signature[:32]
	s := signature[32:64]
	v := signature[64] + 27 // Convert to Ethereum format

	return map[string]interface{}{
		"r": "0x" + hex.EncodeToString(r),
		"s": "0x" + hex.EncodeToString(s),
		"v": int(v),
	}, nil
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
		A: asset,        // asset ID
		B: order.IsBuy,  // is buy order
		P: limitPxWire,  // limit price
		S: szWire,       // size
		R: order.ReduceOnly, // reduce only
		T: orderTypeWire, // order type
	}

	if order.Cloid != nil {
		cloidStr := order.Cloid.ToRaw()
		wire.C = &cloidStr
	}

	return wire, nil
}

// OrderWiresToOrderAction converts order wires to order action
func OrderWiresToOrderAction(orderWires []types.OrderWire, builder *types.BuilderInfo) map[string]interface{} {
	// Create ordered map to ensure consistent field ordering for msgpack
	// Field order MUST match Python SDK exactly per Hyperliquid docs
	action := make(map[string]interface{})
	
	// Add fields in the EXACT order used by Python SDK
	// Python uses: {"type": "order", "orders": [...], "grouping": "na"}
	action["type"] = "order"
	action["orders"] = orderWires
	action["grouping"] = "na"
	
	if builder != nil {
		action["builder"] = builder
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
