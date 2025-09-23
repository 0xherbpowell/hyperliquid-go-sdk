package types

import (
	"fmt"
	"strings"
)

// Side represents the side of an order (Buy/Sell)
type Side string

const (
	SideBuy  Side = "A"
	SideSell Side = "B"
)

// Tif represents the time in force for orders
type Tif string

const (
	TifAlo Tif = "Alo" // Add liquidity only
	TifIoc Tif = "Ioc" // Immediate or cancel
	TifGtc Tif = "Gtc" // Good till cancel
)

// Tpsl represents take profit or stop loss
type Tpsl string

const (
	TpslTp Tpsl = "tp" // Take profit
	TpslSl Tpsl = "sl" // Stop loss
)

// Grouping represents order grouping
type Grouping string

const (
	GroupingNa           Grouping = "na"
	GroupingNormalTpsl   Grouping = "normalTpsl"
	GroupingPositionTpsl Grouping = "positionTpsl"
)

// Cloid represents a client order ID (16 bytes hex)
type Cloid struct {
	rawCloid string
}

// NewCloid creates a new Cloid from a hex string
func NewCloid(raw string) (*Cloid, error) {
	if !strings.HasPrefix(raw, "0x") {
		return nil, fmt.Errorf("cloid is not a hex string")
	}
	if len(raw[2:]) != 32 {
		return nil, fmt.Errorf("cloid is not 16 bytes")
	}
	return &Cloid{rawCloid: raw}, nil
}

// NewCloidFromInt creates a new Cloid from an integer
func NewCloidFromInt(cloid int64) *Cloid {
	return &Cloid{rawCloid: fmt.Sprintf("0x%032x", cloid)}
}

// String returns the raw cloid string
func (c *Cloid) String() string {
	if c == nil {
		return ""
	}
	return c.rawCloid
}

// ToRaw returns the raw cloid string
func (c *Cloid) ToRaw() string {
	if c == nil {
		return ""
	}
	return c.rawCloid
}

// MarshalJSON implements the json.Marshaler interface
func (c *Cloid) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte("null"), nil
	}
	return []byte(`"` + c.rawCloid + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (c *Cloid) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		c.rawCloid = ""
		return nil
	}
	// Remove quotes from JSON string
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		c.rawCloid = string(data[1 : len(data)-1])
	} else {
		c.rawCloid = string(data)
	}
	return nil
}

type MarginTier struct {
	LowerBound  string `json:"lowerBound"`
	MaxLeverage int    `json:"maxLeverage"`
}

// AssetInfo represents metadata about an asset
type AssetInfo struct {
	Name       string `json:"name"`
	SzDecimals int    `json:"szDecimals"`
}

type MarginTable struct {
	ID          int
	Description string       `json:"description"`
	MarginTiers []MarginTier `json:"marginTiers"`
}

// Meta represents the universe of assets
type Meta struct {
	Universe     []AssetInfo   `json:"universe"`
	MarginTables []MarginTable `json:"marginTables"`
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
	TokenId     string  `json:"tokenId"`
	IsCanonical bool    `json:"isCanonical"`
	EvmContract *string `json:"evmContract,omitempty"`
	FullName    *string `json:"fullName,omitempty"`
}

// SpotMeta represents spot asset metadata
type SpotMeta struct {
	Universe []SpotAssetInfo `json:"universe"`
	Tokens   []SpotTokenInfo `json:"tokens"`
}

// SpotAssetCtx represents spot asset context
type SpotAssetCtx struct {
	DayNtlVlm         string  `json:"dayNtlVlm"`
	MarkPx            string  `json:"markPx"`
	MidPx             *string `json:"midPx,omitempty"`
	PrevDayPx         string  `json:"prevDayPx"`
	CirculatingSupply string  `json:"circulatingSupply"`
	Coin              string  `json:"coin"`
}

// PerpAssetCtx represents perpetual asset context
type PerpAssetCtx struct {
	Funding      string     `json:"funding"`
	OpenInterest string     `json:"openInterest"`
	PrevDayPx    string     `json:"prevDayPx"`
	DayNtlVlm    string     `json:"dayNtlVlm"`
	Premium      string     `json:"premium"`
	OraclePx     string     `json:"oraclePx"`
	MarkPx       string     `json:"markPx"`
	MidPx        *string    `json:"midPx,omitempty"`
	ImpactPxs    *[2]string `json:"impactPxs,omitempty"`
	DayBaseVlm   string     `json:"dayBaseVlm"`
}

// LimitOrderType represents a limit order
type LimitOrderType struct {
	Tif Tif `json:"tif" msgpack:"tif"`
}

// TriggerOrderType represents a trigger order
type TriggerOrderType struct {
	TriggerPx float64 `json:"triggerPx"`
	IsMarket  bool    `json:"isMarket"`
	Tpsl      Tpsl    `json:"tpsl"`
}

// TriggerOrderTypeWire represents a trigger order for wire format
type TriggerOrderTypeWire struct {
	IsMarket  bool   `json:"isMarket" msgpack:"isMarket"`
	TriggerPx string `json:"triggerPx" msgpack:"triggerPx"`
	Tpsl      Tpsl   `json:"tpsl" msgpack:"tpsl"`
}

// OrderType represents the type of order
type OrderType struct {
	Limit   *LimitOrderType   `json:"limit,omitempty"`
	Trigger *TriggerOrderType `json:"trigger,omitempty"`
}

// OrderTypeWire represents the wire format of OrderType
type OrderTypeWire struct {
	Limit   *LimitOrderType       `json:"limit,omitempty" msgpack:"limit,omitempty"`
	Trigger *TriggerOrderTypeWire `json:"trigger,omitempty" msgpack:"trigger,omitempty"`
}

// OrderRequest represents a request to place an order
type OrderRequest struct {
	Coin       string    `json:"coin"`
	IsBuy      bool      `json:"is_buy"`
	Sz         float64   `json:"sz"`
	LimitPx    float64   `json:"limit_px"`
	OrderType  OrderType `json:"order_type"`
	ReduceOnly bool      `json:"reduce_only"`
	Cloid      *Cloid    `json:"cloid,omitempty"`
}

// OrderWire represents the wire format of an order
type OrderWire struct {
	A int           `json:"a" msgpack:"a"`                     // asset
	B bool          `json:"b" msgpack:"b"`                     // isBuy
	P string        `json:"p" msgpack:"p"`                     // limitPx
	S string        `json:"s" msgpack:"s"`                     // sz
	R bool          `json:"r" msgpack:"r"`                     // reduceOnly
	T OrderTypeWire `json:"t" msgpack:"t"`                     // orderType
	C *string       `json:"c,omitempty" msgpack:"c,omitempty"` // cloid
}

// Order represents an order
type Order struct {
	Asset      int     `json:"asset"`
	IsBuy      bool    `json:"isBuy"`
	LimitPx    float64 `json:"limitPx"`
	Sz         float64 `json:"sz"`
	ReduceOnly bool    `json:"reduceOnly"`
	Cloid      *Cloid  `json:"cloid,omitempty"`
}

// ModifyRequest represents a request to modify an order
type ModifyRequest struct {
	Oid   interface{}  `json:"oid"` // Can be int or Cloid
	Order OrderRequest `json:"order"`
}

// ModifyWire represents the wire format of a modify request
type ModifyWire struct {
	Oid   int       `json:"oid"`
	Order OrderWire `json:"order"`
}

// CancelRequest represents a request to cancel an order
type CancelRequest struct {
	Coin string `json:"coin"`
	Oid  int    `json:"oid"`
}

// CancelByCloidRequest represents a request to cancel an order by cloid
type CancelByCloidRequest struct {
	Coin  string `json:"coin"`
	Cloid *Cloid `json:"cloid"`
}

// CrossLeverage represents cross leverage
type CrossLeverage struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

// IsolatedLeverage represents isolated leverage
type IsolatedLeverage struct {
	Type   string `json:"type"`
	Value  int    `json:"value"`
	RawUsd string `json:"rawUsd"`
}

// Leverage represents leverage (cross or isolated)
type Leverage struct {
	Type   string `json:"type"`
	Value  int    `json:"value"`
	RawUsd string `json:"rawUsd,omitempty"`
}

// L2Level represents a level 2 order book entry
type L2Level struct {
	Px string `json:"px"`
	Sz string `json:"sz"`
	N  int    `json:"n"`
}

// Trade represents a trade
type Trade struct {
	Coin string `json:"coin"`
	Side Side   `json:"side"`
	Px   string `json:"px"`
	Sz   string `json:"sz"`
	Hash string `json:"hash"`
	Time int64  `json:"time"`
}

// Fill represents a fill
type Fill struct {
	Coin          string `json:"coin"`
	Px            string `json:"px"`
	Sz            string `json:"sz"`
	Side          Side   `json:"side"`
	Time          int64  `json:"time"`
	StartPosition string `json:"startPosition"`
	Dir           string `json:"dir"`
	ClosedPnl     string `json:"closedPnl"`
	Hash          string `json:"hash"`
	Oid           int    `json:"oid"`
	Crossed       bool   `json:"crossed"`
	Fee           string `json:"fee"`
	Tid           int    `json:"tid"`
	FeeToken      string `json:"feeToken"`
}

// BuilderInfo represents builder information
type BuilderInfo struct {
	B string `json:"b"` // Public address of the builder
	F int    `json:"f"` // Amount of fee in tenths of basis points
}

// ScheduleCancelAction represents a schedule cancel action
type ScheduleCancelAction struct {
	Type string `json:"type"`
	Time *int64 `json:"time,omitempty"`
}

// ActiveAssetCtx represents active asset context
type ActiveAssetCtx struct {
	Coin string       `json:"coin"`
	Ctx  PerpAssetCtx `json:"ctx"`
}

// ActiveSpotAssetCtx represents active spot asset context
type ActiveSpotAssetCtx struct {
	Coin string       `json:"coin"`
	Ctx  SpotAssetCtx `json:"ctx"`
}

// ActiveAssetData represents active asset data
type ActiveAssetData struct {
	User             string    `json:"user"`
	Coin             string    `json:"coin"`
	Leverage         Leverage  `json:"leverage"`
	MaxTradeSzs      [2]string `json:"maxTradeSzs"`
	AvailableToTrade [2]string `json:"availableToTrade"`
	MarkPx           string    `json:"markPx"`
}

// Subscription represents a WebSocket subscription
type Subscription struct {
	Type     string `json:"type"`
	Coin     string `json:"coin,omitempty"`
	User     string `json:"user,omitempty"`
	Interval string `json:"interval,omitempty"`
}

// AllMidsData represents all mids data
type AllMidsData struct {
	Mids map[string]string `json:"mids"`
}

// AllMidsMsg represents an all mids message
type AllMidsMsg struct {
	Channel string      `json:"channel"`
	Data    AllMidsData `json:"data"`
}

// L2BookData represents level 2 book data
type L2BookData struct {
	Coin   string       `json:"coin"`
	Levels [2][]L2Level `json:"levels"`
	Time   int64        `json:"time"`
}

// L2BookMsg represents a level 2 book message
type L2BookMsg struct {
	Channel string     `json:"channel"`
	Data    L2BookData `json:"data"`
}

// BboData represents best bid offer data
type BboData struct {
	Coin string      `json:"coin"`
	Time int64       `json:"time"`
	Bbo  [2]*L2Level `json:"bbo"`
}

// BboMsg represents a BBO message
type BboMsg struct {
	Channel string  `json:"channel"`
	Data    BboData `json:"data"`
}

// TradesMsg represents a trades message
type TradesMsg struct {
	Channel string  `json:"channel"`
	Data    []Trade `json:"data"`
}

// UserEventsData represents user events data
type UserEventsData struct {
	Fills []Fill `json:"fills,omitempty"`
}

// UserEventsMsg represents a user events message
type UserEventsMsg struct {
	Channel string         `json:"channel"`
	Data    UserEventsData `json:"data"`
}

// UserFillsData represents user fills data
type UserFillsData struct {
	User       string `json:"user"`
	IsSnapshot bool   `json:"isSnapshot"`
	Fills      []Fill `json:"fills"`
}

// UserFillsMsg represents a user fills message
type UserFillsMsg struct {
	Channel string        `json:"channel"`
	Data    UserFillsData `json:"data"`
}

// PongMsg represents a pong message
type PongMsg struct {
	Channel string `json:"channel"`
}

// ActiveAssetCtxMsg represents an active asset context message
type ActiveAssetCtxMsg struct {
	Channel string         `json:"channel"`
	Data    ActiveAssetCtx `json:"data"`
}

// ActiveSpotAssetCtxMsg represents an active spot asset context message
type ActiveSpotAssetCtxMsg struct {
	Channel string             `json:"channel"`
	Data    ActiveSpotAssetCtx `json:"data"`
}

// ActiveAssetDataMsg represents an active asset data message
type ActiveAssetDataMsg struct {
	Channel string          `json:"channel"`
	Data    ActiveAssetData `json:"data"`
}

// OtherWsMsg represents other WebSocket messages
type OtherWsMsg struct {
	Channel string      `json:"channel"`
	Data    interface{} `json:"data,omitempty"`
}

// PerpDexSchemaInput represents perp dex schema input
type PerpDexSchemaInput struct {
	FullName        string  `json:"fullName"`
	CollateralToken int     `json:"collateralToken"`
	OracleUpdater   *string `json:"oracleUpdater,omitempty"`
}
