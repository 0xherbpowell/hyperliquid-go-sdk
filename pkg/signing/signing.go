package signing

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"
	"hyperliquid-go-sdk/pkg/constants"
	"hyperliquid-go-sdk/pkg/types"
)

// EIP712Type represents a type in EIP-712 signing
type EIP712Type struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// EIP712Domain represents the EIP-712 domain
type EIP712Domain struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainID           int    `json:"chainId"`
	VerifyingContract string `json:"verifyingContract"`
}

// PhantomAgent represents a phantom agent for L1 actions
type PhantomAgent struct {
	Source       string `json:"source"`
	ConnectionID string `json:"connectionId"`
}

// ActionHash computes the hash of an action for signing
func ActionHash(action interface{}, vaultAddress *string, nonce int64, expiresAfter *int64) ([]byte, error) {
	// Pack action with msgpack
	data, err := msgpack.Marshal(action)
	if err != nil {
		return nil, err
	}

	// Add nonce (8 bytes big endian)
	nonceBytes := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		nonceBytes[i] = byte(nonce & 0xff)
		nonce >>= 8
	}
	data = append(data, nonceBytes...)

	// Add vault address
	if vaultAddress == nil {
		data = append(data, 0x00)
	} else {
		data = append(data, 0x01)
		vaultBytes := AddressToBytes(*vaultAddress)
		data = append(data, vaultBytes...)
	}

	// Add expires after
	if expiresAfter != nil {
		data = append(data, 0x00)
		expiresBytes := make([]byte, 8)
		expires := *expiresAfter
		for i := 7; i >= 0; i-- {
			expiresBytes[i] = byte(expires & 0xff)
			expires >>= 8
		}
		data = append(data, expiresBytes...)
	}

	// Compute keccak256 hash
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil), nil
}

// ConstructPhantomAgent creates a phantom agent for L1 actions
func ConstructPhantomAgent(hash []byte, isMainnet bool) PhantomAgent {
	source := "b" // testnet
	if isMainnet {
		source = "a" // mainnet
	}

	return PhantomAgent{
		Source:       source,
		ConnectionID: "0x" + hex.EncodeToString(hash),
	}
}

// L1Payload creates the EIP-712 payload for L1 actions
func L1Payload(phantomAgent PhantomAgent, action interface{}) apitypes.TypedData {
	actionMap, ok := action.(map[string]interface{})
	if !ok {
		return apitypes.TypedData{}
	}

	// Get chain ID from action
	chainIDHex, ok := actionMap["signatureChainId"].(string)
	if !ok {
		return apitypes.TypedData{}
	}

	// Parse chain ID
	chainID := common.HexToHash(chainIDHex).Big()

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
			Name:              constants.ExchangeDomain,
			Version:           constants.DomainVersion,
			ChainId:           (*math.HexOrDecimal256)(chainID),
			VerifyingContract: constants.VerifyingContract,
		},
		Message: apitypes.TypedDataMessage{
			"source":       phantomAgent.Source,
			"connectionId": phantomAgent.ConnectionID,
		},
	}
}

// UserSignedPayload creates the EIP-712 payload for user-signed actions
func UserSignedPayload(primaryType string, payloadTypes []EIP712Type, action interface{}) (apitypes.TypedData, error) {
	// Convert action to map for EIP-712 message
	actionMap, ok := action.(map[string]interface{})
	if !ok {
		return apitypes.TypedData{}, fmt.Errorf("action must be a map[string]interface{}")
	}

	// Get chain ID from action
	chainIDHex, ok := actionMap["signatureChainId"].(string)
	if !ok {
		return apitypes.TypedData{}, fmt.Errorf("action must contain signatureChainId")
	}

	// Parse chain ID
	chainID := common.HexToHash(chainIDHex).Big()

	// Convert payload types to apitypes.Type
	types := make([]apitypes.Type, len(payloadTypes))
	for i, pt := range payloadTypes {
		types[i] = apitypes.Type{Name: pt.Name, Type: pt.Type}
	}

	return apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			primaryType: types,
		},
		PrimaryType: primaryType,
		Domain: apitypes.TypedDataDomain{
			Name:              constants.HyperliquidDomain,
			Version:           constants.DomainVersion,
			ChainId:           (*math.HexOrDecimal256)(chainID),
			VerifyingContract: constants.VerifyingContract,
		},
		Message: actionMap,
	}, nil
}

// SignL1Action signs an L1 action
func SignL1Action(privateKey *ecdsa.PrivateKey, action interface{}, vaultAddress *string, nonce int64, expiresAfter *int64, isMainnet bool) (types.Signature, error) {
	hash, err := ActionHash(action, vaultAddress, nonce, expiresAfter)
	if err != nil {
		return types.Signature{}, err
	}

	phantomAgent := ConstructPhantomAgent(hash, isMainnet)
	payload := L1Payload(phantomAgent, action)

	return SignInner(privateKey, payload)
}

// SignUserSignedAction signs a user-signed action
func SignUserSignedAction(privateKey *ecdsa.PrivateKey, action interface{}, payloadTypes []EIP712Type, primaryType string, isMainnet bool) (types.Signature, error) {
	// Add chain information to action
	actionMap, ok := action.(map[string]interface{})
	if !ok {
		return types.Signature{}, fmt.Errorf("action must be a map[string]interface{}")
	}

	actionMap["signatureChainId"] = constants.SignatureChainID
	if isMainnet {
		actionMap["hyperliquidChain"] = constants.MainnetChain
	} else {
		actionMap["hyperliquidChain"] = constants.TestnetChain
	}

	payload, err := UserSignedPayload(primaryType, payloadTypes, actionMap)
	if err != nil {
		return types.Signature{}, err
	}

	return SignInner(privateKey, payload)
}

// SignInner performs the actual EIP-712 signing
func SignInner(privateKey *ecdsa.PrivateKey, data apitypes.TypedData) (types.Signature, error) {
	// Encode the typed data
	domainSeparator, err := data.HashStruct("EIP712Domain", data.Domain.Map())
	if err != nil {
		return types.Signature{}, err
	}

	typedDataHash, err := data.HashStruct(data.PrimaryType, data.Message)
	if err != nil {
		return types.Signature{}, err
	}

	// Create the final hash
	rawData := []byte{0x19, 0x01}
	rawData = append(rawData, domainSeparator...)
	rawData = append(rawData, typedDataHash...)
	hash := crypto.Keccak256(rawData)

	// Sign the hash
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return types.Signature{}, err
	}

	// Convert to R, S, V format
	r := hex.EncodeToString(signature[0:32])
	s := hex.EncodeToString(signature[32:64])
	v := int(signature[64]) + 27

	return types.Signature{
		R: "0x" + r,
		S: "0x" + s,
		V: v,
	}, nil
}

// RecoverAgentOrUserFromL1Action recovers the signer address from an L1 action
func RecoverAgentOrUserFromL1Action(action interface{}, signature types.Signature, vaultAddress *string, nonce int64, expiresAfter *int64, isMainnet bool) (common.Address, error) {
	hash, err := ActionHash(action, vaultAddress, nonce, expiresAfter)
	if err != nil {
		return common.Address{}, err
	}

	phantomAgent := ConstructPhantomAgent(hash, isMainnet)
	payload := L1Payload(phantomAgent, action)

	return RecoverAddress(payload, signature)
}

// RecoverUserFromUserSignedAction recovers the signer address from a user-signed action
func RecoverUserFromUserSignedAction(action interface{}, signature types.Signature, payloadTypes []EIP712Type, primaryType string, isMainnet bool) (common.Address, error) {
	// Add chain information to action
	actionMap, ok := action.(map[string]interface{})
	if !ok {
		return common.Address{}, fmt.Errorf("action must be a map[string]interface{}")
	}

	if isMainnet {
		actionMap["hyperliquidChain"] = constants.MainnetChain
	} else {
		actionMap["hyperliquidChain"] = constants.TestnetChain
	}

	payload, err := UserSignedPayload(primaryType, payloadTypes, actionMap)
	if err != nil {
		return common.Address{}, err
	}

	return RecoverAddress(payload, signature)
}

// RecoverAddress recovers the signer address from EIP-712 data and signature
func RecoverAddress(data apitypes.TypedData, signature types.Signature) (common.Address, error) {
	// Encode the typed data
	domainSeparator, err := data.HashStruct("EIP712Domain", data.Domain.Map())
	if err != nil {
		return common.Address{}, err
	}

	typedDataHash, err := data.HashStruct(data.PrimaryType, data.Message)
	if err != nil {
		return common.Address{}, err
	}

	// Create the final hash
	rawData := []byte{0x19, 0x01}
	rawData = append(rawData, domainSeparator...)
	rawData = append(rawData, typedDataHash...)
	hash := crypto.Keccak256(rawData)

	// Convert signature to bytes
	r := common.HexToHash(signature.R).Bytes()
	s := common.HexToHash(signature.S).Bytes()
	v := byte(signature.V - 27)

	sigBytes := append(r, s...)
	sigBytes = append(sigBytes, v)

	// Recover public key
	pubKey, err := crypto.SigToPub(hash, sigBytes)
	if err != nil {
		return common.Address{}, err
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}

// Convenience functions for specific action types

// SignUSDTransferAction signs a USD transfer action
func SignUSDTransferAction(privateKey *ecdsa.PrivateKey, action USDTransferAction, isMainnet bool) (types.Signature, error) {
	return SignUserSignedAction(privateKey, action, USDSendSignTypes, "HyperliquidTransaction:UsdSend", isMainnet)
}

// SignSpotTransferAction signs a spot transfer action
func SignSpotTransferAction(privateKey *ecdsa.PrivateKey, action SpotTransferAction, isMainnet bool) (types.Signature, error) {
	return SignUserSignedAction(privateKey, action, SpotTransferSignTypes, "HyperliquidTransaction:SpotSend", isMainnet)
}

// SignWithdrawAction signs a withdraw action
func SignWithdrawAction(privateKey *ecdsa.PrivateKey, action WithdrawAction, isMainnet bool) (types.Signature, error) {
	return SignUserSignedAction(privateKey, action, WithdrawSignTypes, "HyperliquidTransaction:Withdraw", isMainnet)
}

// SignUSDClassTransferAction signs a USD class transfer action
func SignUSDClassTransferAction(privateKey *ecdsa.PrivateKey, action USDClassTransferAction, isMainnet bool) (types.Signature, error) {
	return SignUserSignedAction(privateKey, action, USDClassTransferSignTypes, "HyperliquidTransaction:UsdClassTransfer", isMainnet)
}

// SignSendAssetAction signs a send asset action
func SignSendAssetAction(privateKey *ecdsa.PrivateKey, action SendAssetAction, isMainnet bool) (types.Signature, error) {
	return SignUserSignedAction(privateKey, action, SendAssetSignTypes, "HyperliquidTransaction:SendAsset", isMainnet)
}

// SignTokenDelegateAction signs a token delegate action
func SignTokenDelegateAction(privateKey *ecdsa.PrivateKey, action TokenDelegateAction, isMainnet bool) (types.Signature, error) {
	return SignUserSignedAction(privateKey, action, TokenDelegateTypes, "HyperliquidTransaction:TokenDelegate", isMainnet)
}

// SignConvertToMultiSigUserAction signs a convert to multi-sig user action
func SignConvertToMultiSigUserAction(privateKey *ecdsa.PrivateKey, action ConvertToMultiSigUserAction, isMainnet bool) (types.Signature, error) {
	return SignUserSignedAction(privateKey, action, ConvertToMultiSigUserSignTypes, "HyperliquidTransaction:ConvertToMultiSigUser", isMainnet)
}

// SignApproveAgentAction signs an approve agent action
func SignApproveAgentAction(privateKey *ecdsa.PrivateKey, action ApproveAgentAction, isMainnet bool) (types.Signature, error) {
	return SignUserSignedAction(privateKey, action, ApproveAgentSignTypes, "HyperliquidTransaction:ApproveAgent", isMainnet)
}

// SignApproveBuilderFeeAction signs an approve builder fee action
func SignApproveBuilderFeeAction(privateKey *ecdsa.PrivateKey, action ApproveBuilderFeeAction, isMainnet bool) (types.Signature, error) {
	return SignUserSignedAction(privateKey, action, ApproveBuilderFeeSignTypes, "HyperliquidTransaction:ApproveBuilderFee", isMainnet)
}
