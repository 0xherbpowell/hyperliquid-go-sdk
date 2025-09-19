package types

import (
	"encoding/json"
	"time"
)

// Address represents a blockchain address
type Address string

// Hash represents a transaction hash
type Hash string

// FloatString represents a float value as string for precise decimal handling
type FloatString string

// Side represents order side (A for Ask/Sell, B for Bid/Buy)
type Side string

const (
	SideAsk Side = "A"
	SideBid Side = "B"
)

// Tif represents time in force
type Tif string

const (
	TifAlo Tif = "Alo" // Add liquidity only
	TifIoc Tif = "Ioc" // Immediate or cancel
	TifGtc Tif = "Gtc" // Good till cancelled
)

// Tpsl represents take profit / stop loss type
type Tpsl string

const (
	TpslTp Tpsl = "tp" // Take profit
	TpslSl Tpsl = "sl" // Stop loss
)

// AssetInfo represents basic asset information
type AssetInfo struct {
	Name       string `json:"name"`
	SzDecimals int    `json:"szDecimals"`
}

// Meta represents exchange metadata
type Meta struct {
	Universe []AssetInfo `json:"universe"`
}

// SpotAssetInfo represents spot asset information
type SpotAssetInfo struct {
	Name        string `json:"name"`
	Tokens      []int  `json:"tokens"`
	Index       int    `json:"index"`
	IsCanonical bool   `json:"isCanonical"`
}

// SpotTokenInfo represents spot token information
type SpotTokenInfo struct {
	Name        string  `json:"name"`
	SzDecimals  int     `json:"szDecimals"`
	WeiDecimals int     `json:"weiDecimals"`
	Index       int     `json:"index"`
	TokenID     string  `json:"tokenId"`
	IsCanonical bool    `json:"isCanonical"`
	EvmContract *string `json:"evmContract,omitempty"`
	FullName    *string `json:"fullName,omitempty"`
}

// SpotMeta represents spot exchange metadata
type SpotMeta struct {
	Universe []SpotAssetInfo `json:"universe"`
	Tokens   []SpotTokenInfo `json:"tokens"`
}

// SpotAssetCtx represents spot asset context
type SpotAssetCtx struct {
	DayNtlVlm         FloatString  `json:"dayNtlVlm"`
	MarkPx            FloatString  `json:"markPx"`
	MidPx             *FloatString `json:"midPx,omitempty"`
	PrevDayPx         FloatString  `json:"prevDayPx"`
	CirculatingSupply FloatString  `json:"circulatingSupply"`
	Coin              string       `json:"coin"`
}

// SpotMetaAndAssetCtxs represents combined spot metadata and asset contexts
type SpotMetaAndAssetCtxs struct {
	Meta     SpotMeta       `json:"meta"`
	AssetCtx []SpotAssetCtx `json:"assetCtx"`
}

// L2Level represents a level in the order book
type L2Level struct {
	Px FloatString `json:"px"` // Price
	Sz FloatString `json:"sz"` // Size
	N  int         `json:"n"`  // Number of orders
}

// L2BookData represents level 2 order book data
type L2BookData struct {
	Coin   string       `json:"coin"`
	Levels [2][]L2Level `json:"levels"` // [bids, asks]
	Time   int64        `json:"time"`
}

// Trade represents a trade
type Trade struct {
	Coin string      `json:"coin"`
	Side Side        `json:"side"`
	Px   FloatString `json:"px"`
	Sz   int         `json:"sz"`
	Hash Hash        `json:"hash"`
	Time int64       `json:"time"`
}

// CrossLeverage represents cross margin leverage
type CrossLeverage struct {
	Type  string `json:"type"`  // "cross"
	Value int    `json:"value"` // leverage multiplier
}

// IsolatedLeverage represents isolated margin leverage
type IsolatedLeverage struct {
	Type   string      `json:"type"`   // "isolated"
	Value  int         `json:"value"`  // leverage multiplier
	RawUsd FloatString `json:"rawUsd"` // raw USD amount
}

// Leverage represents either cross or isolated leverage
type Leverage interface {
	GetType() string
	GetValue() int
}

func (c CrossLeverage) GetType() string { return c.Type }
func (c CrossLeverage) GetValue() int   { return c.Value }

func (i IsolatedLeverage) GetType() string { return i.Type }
func (i IsolatedLeverage) GetValue() int   { return i.Value }

// PerpAssetCtx represents perpetual asset context
type PerpAssetCtx struct {
	Funding      FloatString     `json:"funding"`
	OpenInterest FloatString     `json:"openInterest"`
	PrevDayPx    FloatString     `json:"prevDayPx"`
	DayNtlVlm    FloatString     `json:"dayNtlVlm"`
	Premium      FloatString     `json:"premium"`
	OraclePx     FloatString     `json:"oraclePx"`
	MarkPx       FloatString     `json:"markPx"`
	MidPx        *FloatString    `json:"midPx,omitempty"`
	ImpactPxs    *[2]FloatString `json:"impactPxs,omitempty"`
	DayBaseVlm   FloatString     `json:"dayBaseVlm"`
}

// Fill represents a trade fill
type Fill struct {
	Coin          string      `json:"coin"`
	Px            FloatString `json:"px"`
	Sz            FloatString `json:"sz"`
	Side          Side        `json:"side"`
	Time          int64       `json:"time"`
	StartPosition FloatString `json:"startPosition"`
	Dir           string      `json:"dir"`
	ClosedPnl     FloatString `json:"closedPnl"`
	Hash          Hash        `json:"hash"`
	Oid           int         `json:"oid"`
	Crossed       bool        `json:"crossed"`
	Fee           FloatString `json:"fee"`
	Tid           int         `json:"tid"`
	FeeToken      string      `json:"feeToken"`
}

// BuilderInfo represents builder information for orders
type BuilderInfo struct {
	B string `json:"b"` // Builder address
	F int    `json:"f"` // Fee in tenths of basis points
}

// MarginSummary represents margin summary information
type MarginSummary struct {
	AccountValue    FloatString `json:"accountValue"`
	TotalMarginUsed FloatString `json:"totalMarginUsed"`
	TotalNtlPos     FloatString `json:"totalNtlPos"`
	TotalRawUsd     FloatString `json:"totalRawUsd"`
}

// Position represents a trading position
type Position struct {
	Coin           string          `json:"coin"`
	EntryPx        *FloatString    `json:"entryPx"`
	Leverage       json.RawMessage `json:"leverage"` // Can be CrossLeverage or IsolatedLeverage
	LiquidationPx  *FloatString    `json:"liquidationPx"`
	MarginUsed     FloatString     `json:"marginUsed"`
	MaxTradeSzs    [2]FloatString  `json:"maxTradeSzs"`
	PositionValue  FloatString     `json:"positionValue"`
	ReturnOnEquity FloatString     `json:"returnOnEquity"`
	Szi            FloatString     `json:"szi"` // Signed size
	UnrealizedPnl  FloatString     `json:"unrealizedPnl"`
}

// AssetPosition represents an asset position
type AssetPosition struct {
	Position Position `json:"position"`
	Type     string   `json:"type"` // Usually "oneWay"
}

// UserState represents user's trading state
type UserState struct {
	AssetPositions     []AssetPosition `json:"assetPositions"`
	MarginSummary      MarginSummary   `json:"marginSummary"`
	CrossMarginSummary MarginSummary   `json:"crossMarginSummary"`
	Withdrawable       FloatString     `json:"withdrawable"`
}

// SpotBalance represents a spot balance
type SpotBalance struct {
	Coin  string      `json:"coin"`
	Hold  FloatString `json:"hold"`
	Total FloatString `json:"total"`
}

// SpotUserState represents user's spot trading state
type SpotUserState struct {
	Balances []SpotBalance `json:"balances"`
}

// OpenOrder represents an open order
type OpenOrder struct {
	Coin      string      `json:"coin"`
	LimitPx   FloatString `json:"limitPx"`
	Oid       int         `json:"oid"`
	Side      Side        `json:"side"`
	Sz        FloatString `json:"sz"`
	Timestamp int64       `json:"timestamp"`
}

// FrontendOpenOrder represents an open order with additional frontend info
type FrontendOpenOrder struct {
	Children         []any       `json:"children"`
	Coin             string      `json:"coin"`
	IsPositionTpsl   bool        `json:"isPositionTpsl"`
	IsTrigger        bool        `json:"isTrigger"`
	LimitPx          FloatString `json:"limitPx"`
	Oid              int         `json:"oid"`
	OrderType        string      `json:"orderType"`
	OrigSz           FloatString `json:"origSz"`
	ReduceOnly       bool        `json:"reduceOnly"`
	Side             Side        `json:"side"`
	Sz               FloatString `json:"sz"`
	Tif              string      `json:"tif"`
	Timestamp        int64       `json:"timestamp"`
	TriggerCondition string      `json:"triggerCondition"`
	TriggerPx        FloatString `json:"triggerPx"`
}

// AllMids represents all mid prices
type AllMids map[string]FloatString

// Signature represents an EIP-712 signature
type Signature struct {
	R string `json:"r"`
	S string `json:"s"`
	V int    `json:"v"`
}

// Config represents SDK configuration
type Config struct {
	SecretKey      string `json:"secret_key"`
	AccountAddress string `json:"account_address"`
	KeystorePath   string `json:"keystore_path"`
}

// GetTimestampMs returns current timestamp in milliseconds
func GetTimestampMs() int64 {
	return time.Now().UnixMilli()
}
