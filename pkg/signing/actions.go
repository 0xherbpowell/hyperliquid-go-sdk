package signing

import (
	"hyperliquid-go-sdk/pkg/types"
)

// Sign type definitions for different actions

var (
	// USDSendSignTypes represents the signing types for USD transfer
	USDSendSignTypes = []EIP712Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	}

	// SpotTransferSignTypes represents the signing types for spot transfer
	SpotTransferSignTypes = []EIP712Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "token", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	}

	// WithdrawSignTypes represents the signing types for withdrawal
	WithdrawSignTypes = []EIP712Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	}

	// USDClassTransferSignTypes represents the signing types for USD class transfer
	USDClassTransferSignTypes = []EIP712Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "toPerp", Type: "bool"},
		{Name: "nonce", Type: "uint64"},
	}

	// SendAssetSignTypes represents the signing types for asset transfer
	SendAssetSignTypes = []EIP712Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "sourceDex", Type: "string"},
		{Name: "destinationDex", Type: "string"},
		{Name: "token", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "fromSubAccount", Type: "string"},
		{Name: "nonce", Type: "uint64"},
	}

	// TokenDelegateTypes represents the signing types for token delegation
	TokenDelegateTypes = []EIP712Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "validator", Type: "address"},
		{Name: "wei", Type: "uint64"},
		{Name: "isUndelegate", Type: "bool"},
		{Name: "nonce", Type: "uint64"},
	}

	// ConvertToMultiSigUserSignTypes represents the signing types for multi-sig conversion
	ConvertToMultiSigUserSignTypes = []EIP712Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "signers", Type: "string"},
		{Name: "nonce", Type: "uint64"},
	}

	// MultiSigEnvelopeSignTypes represents the signing types for multi-sig envelope
	MultiSigEnvelopeSignTypes = []EIP712Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "multiSigActionHash", Type: "bytes32"},
		{Name: "nonce", Type: "uint64"},
	}

	// ApproveAgentSignTypes represents the signing types for agent approval
	ApproveAgentSignTypes = []EIP712Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "agentAddress", Type: "address"},
		{Name: "agentName", Type: "string"},
		{Name: "nonce", Type: "uint64"},
	}

	// ApproveBuilderFeeSignTypes represents the signing types for builder fee approval
	ApproveBuilderFeeSignTypes = []EIP712Type{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "maxFeeRate", Type: "string"},
		{Name: "builder", Type: "address"},
		{Name: "nonce", Type: "uint64"},
	}
)

// Action type definitions

// USDTransferAction represents a USD transfer action
type USDTransferAction struct {
	Type             string `json:"type"`
	SignatureChainID string `json:"signatureChainId"`
	HyperliquidChain string `json:"hyperliquidChain"`
	Destination      string `json:"destination"`
	Amount           string `json:"amount"`
	Time             int64  `json:"time"`
}

// SpotTransferAction represents a spot transfer action
type SpotTransferAction struct {
	Type             string `json:"type"`
	SignatureChainID string `json:"signatureChainId"`
	HyperliquidChain string `json:"hyperliquidChain"`
	Destination      string `json:"destination"`
	Token            string `json:"token"`
	Amount           string `json:"amount"`
	Time             int64  `json:"time"`
}

// WithdrawAction represents a withdrawal action
type WithdrawAction struct {
	Type             string `json:"type"`
	SignatureChainID string `json:"signatureChainId"`
	HyperliquidChain string `json:"hyperliquidChain"`
	Destination      string `json:"destination"`
	Amount           string `json:"amount"`
	Time             int64  `json:"time"`
}

// USDClassTransferAction represents a USD class transfer action
type USDClassTransferAction struct {
	Type             string `json:"type"`
	SignatureChainID string `json:"signatureChainId"`
	HyperliquidChain string `json:"hyperliquidChain"`
	Amount           string `json:"amount"`
	ToPerp           bool   `json:"toPerp"`
	Nonce            int64  `json:"nonce"`
}

// SendAssetAction represents a send asset action
type SendAssetAction struct {
	Type             string `json:"type"`
	SignatureChainID string `json:"signatureChainId"`
	HyperliquidChain string `json:"hyperliquidChain"`
	Destination      string `json:"destination"`
	SourceDex        string `json:"sourceDex"`
	DestinationDex   string `json:"destinationDex"`
	Token            string `json:"token"`
	Amount           string `json:"amount"`
	FromSubAccount   string `json:"fromSubAccount"`
	Nonce            int64  `json:"nonce"`
}

// TokenDelegateAction represents a token delegation action
type TokenDelegateAction struct {
	Type             string `json:"type"`
	SignatureChainID string `json:"signatureChainId"`
	HyperliquidChain string `json:"hyperliquidChain"`
	Validator        string `json:"validator"`
	Wei              int64  `json:"wei"`
	IsUndelegate     bool   `json:"isUndelegate"`
	Nonce            int64  `json:"nonce"`
}

// ConvertToMultiSigUserAction represents a multi-sig conversion action
type ConvertToMultiSigUserAction struct {
	Type             string `json:"type"`
	SignatureChainID string `json:"signatureChainId"`
	HyperliquidChain string `json:"hyperliquidChain"`
	Signers          string `json:"signers"`
	Nonce            int64  `json:"nonce"`
}

// ApproveAgentAction represents an agent approval action
type ApproveAgentAction struct {
	Type             string `json:"type"`
	SignatureChainID string `json:"signatureChainId"`
	HyperliquidChain string `json:"hyperliquidChain"`
	AgentAddress     string `json:"agentAddress"`
	AgentName        string `json:"agentName"`
	Nonce            int64  `json:"nonce"`
}

// ApproveBuilderFeeAction represents a builder fee approval action
type ApproveBuilderFeeAction struct {
	Type             string `json:"type"`
	SignatureChainID string `json:"signatureChainId"`
	HyperliquidChain string `json:"hyperliquidChain"`
	MaxFeeRate       string `json:"maxFeeRate"`
	Builder          string `json:"builder"`
	Nonce            int64  `json:"nonce"`
}

// MultiSigAction represents a multi-sig action
type MultiSigAction struct {
	Type             string            `json:"type"`
	SignatureChainID string            `json:"signatureChainId"`
	Signatures       []types.Signature `json:"signatures"`
	Payload          MultiSigPayload   `json:"payload"`
}

// MultiSigPayload represents the payload for multi-sig actions
type MultiSigPayload struct {
	MultiSigUser string      `json:"multiSigUser"`
	OuterSigner  string      `json:"outerSigner"`
	Action       interface{} `json:"action"`
}

// MultiSigEnvelopeAction represents a multi-sig envelope action
type MultiSigEnvelopeAction struct {
	SignatureChainID   string `json:"signatureChainId"`
	HyperliquidChain   string `json:"hyperliquidChain"`
	MultiSigActionHash []byte `json:"multiSigActionHash"`
	Nonce              int64  `json:"nonce"`
}
