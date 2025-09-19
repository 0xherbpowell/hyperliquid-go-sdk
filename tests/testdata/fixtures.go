package testdata

import (
	"hyperliquid-go-sdk/pkg/types"
)

// Test private keys (DO NOT USE IN PRODUCTION)
const (
	TestPrivateKey1 = "0x0123456789012345678901234567890123456789012345678901234567890123"
	TestPrivateKey2 = "0x1234567890123456789012345678901234567890123456789012345678901234"
	TestAddress1    = "0x5e9ee1089755c3435139848e47e6635505d5a13a"
	TestAddress2    = "0xb7b6f3cea3f66bf525f5d8f965f6dbf6d9b017b2"
	VaultAddress1   = "0x1719884eb866cb12b2287399b15f7db5e7d775ea"
)

// Mock user state for testing
var MockUserState = types.UserState{
	AssetPositions: []types.AssetPosition{
		{
			Position: types.Position{
				Coin:           "BTC",
				EntryPx:        stringPtr("28429.3"),
				MarginUsed:     "0.0",
				PositionValue:  "0.0",
				ReturnOnEquity: "0.0",
				Szi:            "0.0",
				UnrealizedPnl:  "0.0",
			},
			Type: "oneWay",
		},
		{
			Position: types.Position{
				Coin:           "ETH",
				EntryPx:        stringPtr("1800.5"),
				MarginUsed:     "100.0",
				PositionValue:  "500.0",
				ReturnOnEquity: "5.0",
				Szi:            "0.25",
				UnrealizedPnl:  "25.0",
			},
			Type: "oneWay",
		},
	},
	MarginSummary: types.MarginSummary{
		AccountValue:    "1182.312496",
		TotalMarginUsed: "100.0",
		TotalNtlPos:     "500.0",
		TotalRawUsd:     "1182.312496",
	},
	CrossMarginSummary: types.MarginSummary{
		AccountValue:    "1182.312496",
		TotalMarginUsed: "100.0",
		TotalNtlPos:     "500.0",
		TotalRawUsd:     "1182.312496",
	},
	Withdrawable: "1082.312496",
}

// Mock spot user state
var MockSpotUserState = types.SpotUserState{
	Balances: []types.SpotBalance{
		{
			Coin:  "USDC",
			Hold:  "0.0",
			Total: "1000.0",
		},
		{
			Coin:  "PURR",
			Hold:  "50.0",
			Total: "100.0",
		},
	},
}

// Mock open orders
var MockOpenOrders = []types.OpenOrder{
	{
		Coin:      "ETH",
		LimitPx:   "1900.0",
		Oid:       12345,
		Side:      types.SideBid,
		Sz:        "0.1",
		Timestamp: 1677777606040,
	},
	{
		Coin:      "BTC",
		LimitPx:   "29000.0",
		Oid:       12346,
		Side:      types.SideAsk,
		Sz:        "0.01",
		Timestamp: 1677777606050,
	},
}

// Mock fills
var MockFills = []types.Fill{
	{
		Coin:          "ETH",
		Px:            "1850.5",
		Sz:            "0.1",
		Side:          types.SideBid,
		Time:          1677777606000,
		StartPosition: "0.0",
		Dir:           "Open Long",
		ClosedPnl:     "0.0",
		Hash:          "0xabcdef1234567890",
		Oid:           12340,
		Crossed:       true,
		Fee:           "1.8505",
		Tid:           1001,
		FeeToken:      "USDC",
	},
	{
		Coin:          "BTC",
		Px:            "28500.0",
		Sz:            "0.005",
		Side:          types.SideAsk,
		Time:          1677777606010,
		StartPosition: "0.005",
		Dir:           "Close Long",
		ClosedPnl:     "50.0",
		Hash:          "0x1234567890abcdef",
		Oid:           12341,
		Crossed:       true,
		Fee:           "1.425",
		Tid:           1002,
		FeeToken:      "USDC",
	},
}

// Mock L2 book data
var MockL2BookData = types.L2BookData{
	Coin: "ETH",
	Levels: [2][]types.L2Level{
		// Bids
		{
			{Px: "1899.0", Sz: "1.5", N: 3},
			{Px: "1898.5", Sz: "2.0", N: 2},
			{Px: "1898.0", Sz: "0.8", N: 1},
		},
		// Asks
		{
			{Px: "1900.0", Sz: "1.2", N: 2},
			{Px: "1900.5", Sz: "1.8", N: 3},
			{Px: "1901.0", Sz: "2.5", N: 4},
		},
	},
	Time: 1677777606000,
}

// Mock metadata
var MockMeta = types.Meta{
	Universe: []types.AssetInfo{
		{Name: "BTC", SzDecimals: 5},
		{Name: "ETH", SzDecimals: 4},
		{Name: "SOL", SzDecimals: 3},
		{Name: "ATOM", SzDecimals: 2},
	},
}

// Mock spot metadata
var MockSpotMeta = types.SpotMeta{
	Universe: []types.SpotAssetInfo{
		{
			Name:        "PURR/USDC",
			Tokens:      []int{1, 0},
			Index:       0,
			IsCanonical: true,
		},
		{
			Name:        "KORILA/USDC",
			Tokens:      []int{8, 0},
			Index:       8,
			IsCanonical: false,
		},
	},
	Tokens: []types.SpotTokenInfo{
		{
			Name:        "USDC",
			SzDecimals:  6,
			WeiDecimals: 6,
			Index:       0,
			TokenID:     "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
			IsCanonical: true,
		},
		{
			Name:        "PURR",
			SzDecimals:  8,
			WeiDecimals: 18,
			Index:       1,
			TokenID:     "0xc4bf3f870c0e9465323c0b6ed28096c2",
			IsCanonical: false,
		},
		{
			Name:        "KORILA",
			SzDecimals:  8,
			WeiDecimals: 18,
			Index:       8,
			TokenID:     "0x1234567890abcdef1234567890abcdef",
			IsCanonical: false,
		},
	},
}

// Mock all mids
var MockAllMids = types.AllMids{
	"BTC":  "28429.3",
	"ETH":  "1850.5",
	"SOL":  "98.45",
	"ATOM": "11.2585",
	"PURR": "0.5234",
}

// Mock order responses
var MockOrderResponseSuccess = map[string]interface{}{
	"status": "ok",
	"response": map[string]interface{}{
		"type": "order",
		"data": map[string]interface{}{
			"statuses": []interface{}{
				map[string]interface{}{
					"resting": map[string]interface{}{
						"oid": float64(12345),
					},
				},
			},
		},
	},
}

var MockOrderResponseFilled = map[string]interface{}{
	"status": "ok",
	"response": map[string]interface{}{
		"type": "order",
		"data": map[string]interface{}{
			"statuses": []interface{}{
				map[string]interface{}{
					"filled": map[string]interface{}{
						"oid":     float64(12346),
						"totalSz": "0.1",
						"avgPx":   "1850.5",
					},
				},
			},
		},
	},
}

var MockOrderResponseError = map[string]interface{}{
	"status": "ok",
	"response": map[string]interface{}{
		"type": "order",
		"data": map[string]interface{}{
			"statuses": []interface{}{
				map[string]interface{}{
					"error": "Insufficient margin",
				},
			},
		},
	},
}

// Mock cancel response
var MockCancelResponse = map[string]interface{}{
	"status": "ok",
	"response": map[string]interface{}{
		"type": "cancel",
		"data": map[string]interface{}{
			"statuses": []interface{}{
				"success",
			},
		},
	},
}

// Mock WebSocket messages
var MockWebSocketMessages = struct {
	AllMids    map[string]interface{}
	L2Book     map[string]interface{}
	UserEvents map[string]interface{}
	Trades     map[string]interface{}
	Pong       map[string]interface{}
}{
	AllMids: map[string]interface{}{
		"channel": "allMids",
		"data": map[string]interface{}{
			"mids": MockAllMids,
		},
	},
	L2Book: map[string]interface{}{
		"channel": "l2Book",
		"data":    MockL2BookData,
	},
	UserEvents: map[string]interface{}{
		"channel": "user",
		"data": map[string]interface{}{
			"fills": MockFills,
		},
	},
	Trades: map[string]interface{}{
		"channel": "trades",
		"data": []types.Trade{
			{
				Coin: "ETH",
				Side: types.SideBid,
				Px:   "1850.5",
				Sz:   100,
				Hash: "0xabcdef1234567890",
				Time: 1677777606000,
			},
		},
	},
	Pong: map[string]interface{}{
		"channel": "pong",
	},
}

// Mock error responses
var MockErrorResponses = struct {
	BadRequest         map[string]interface{}
	Unauthorized       map[string]interface{}
	InsufficientMargin map[string]interface{}
	InvalidOrderSize   map[string]interface{}
	NetworkError       string
}{
	BadRequest: map[string]interface{}{
		"error": "Bad request",
		"code":  "BAD_REQUEST",
		"msg":   "Invalid request format",
	},
	Unauthorized: map[string]interface{}{
		"error": "Unauthorized",
		"code":  "UNAUTHORIZED",
		"msg":   "Invalid signature",
	},
	InsufficientMargin: map[string]interface{}{
		"error": "Insufficient margin",
		"code":  "INSUFFICIENT_MARGIN",
		"msg":   "Not enough margin to place this order",
	},
	InvalidOrderSize: map[string]interface{}{
		"error": "Invalid order size",
		"code":  "INVALID_SIZE",
		"msg":   "Order size is below minimum",
	},
	NetworkError: "network connection failed",
}

// Test order requests
var TestOrderRequests = struct {
	BasicLimit types.OrderRequest
	MarketBuy  types.OrderRequest
	StopLoss   types.OrderRequest
	TakeProfit types.OrderRequest
	WithCloid  types.OrderRequest
	ReduceOnly types.OrderRequest
}{
	BasicLimit: types.OrderRequest{
		Coin:       "ETH",
		IsBuy:      true,
		Sz:         0.1,
		LimitPx:    1900.0,
		OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: "Gtc"}},
		ReduceOnly: false,
	},
	MarketBuy: types.OrderRequest{
		Coin:       "ETH",
		IsBuy:      true,
		Sz:         0.1,
		LimitPx:    2000.0,
		OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: "Ioc"}},
		ReduceOnly: false,
	},
	StopLoss: types.OrderRequest{
		Coin:    "ETH",
		IsBuy:   false,
		Sz:      0.1,
		LimitPx: 1700.0,
		OrderType: types.OrderType{
			Trigger: &types.TriggerOrderType{
				TriggerPx: 1750.0,
				IsMarket:  true,
				Tpsl:      "sl",
			},
		},
		ReduceOnly: true,
	},
	TakeProfit: types.OrderRequest{
		Coin:    "ETH",
		IsBuy:   false,
		Sz:      0.1,
		LimitPx: 2100.0,
		OrderType: types.OrderType{
			Trigger: &types.TriggerOrderType{
				TriggerPx: 2000.0,
				IsMarket:  true,
				Tpsl:      "tp",
			},
		},
		ReduceOnly: true,
	},
	WithCloid: types.OrderRequest{
		Coin:       "ETH",
		IsBuy:      true,
		Sz:         0.1,
		LimitPx:    1900.0,
		OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: "Gtc"}},
		ReduceOnly: false,
		Cloid:      types.NewCloidFromInt(1),
	},
	ReduceOnly: types.OrderRequest{
		Coin:       "ETH",
		IsBuy:      false,
		Sz:         0.1,
		LimitPx:    1800.0,
		OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: "Gtc"}},
		ReduceOnly: true,
	},
}

// Test signatures (these correspond to the test private keys above)
var TestSignatures = struct {
	L1ActionMainnet types.Signature
	L1ActionTestnet types.Signature
	USDTransfer     types.Signature
	WithdrawBridge  types.Signature
}{
	L1ActionMainnet: types.Signature{
		R: "0x53749d5b30552aeb2fca34b530185976545bb22d0b3ce6f62e31be961a59298",
		S: "0x755c40ba9bf05223521753995abb2f73ab3229be8ec921f350cb447e384d8ed8",
		V: 27,
	},
	L1ActionTestnet: types.Signature{
		R: "0x542af61ef1f429707e3c76c5293c80d01f74ef853e34b76efffcb57e574f9510",
		S: "0x17b8b32f086e8cdede991f1e2c529f5dd5297cbe8128500e00cbaf766204a613",
		V: 28,
	},
	USDTransfer: types.Signature{
		R: "0x637b37dd731507cdd24f46532ca8ba6eec616952c56218baeff04144e4a77073",
		S: "0x11a6a24900e6e314136d2592e2f8d502cd89b7c15b198e1bee043c9589f9fad7",
		V: 27,
	},
	WithdrawBridge: types.Signature{
		R: "0x8363524c799e90ce9bc41022f7c39b4e9bdba786e5f9c72b20e43e1462c37cf9",
		S: "0x58b1411a775938b83e29182e8ef74975f9054c8e97ebf5ec2dc8d51bfc893881",
		V: 28,
	},
}

// Test configuration
var TestConfig = types.Config{
	SecretKey:      TestPrivateKey1,
	AccountAddress: TestAddress1,
	KeystorePath:   "",
}

// Helper functions
func stringPtr(s string) *types.FloatString {
	fs := types.FloatString(s)
	return &fs
}

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

// GetTestCloid returns a test client order ID
func GetTestCloid() *types.Cloid {
	return types.NewCloidFromInt(1)
}

// GetTestBuilder returns a test builder info
func GetTestBuilder() *types.BuilderInfo {
	return &types.BuilderInfo{
		B: TestAddress2,
		F: 10, // 1 basis point
	}
}

// MockPositions for testing position calculations
var MockPositions = []types.DetailedPosition{
	{
		Coin:           "BTC",
		EntryPx:        stringPtr("28000.0"),
		Leverage:       types.LeverageInfo{Type: "cross", Value: 10},
		LiquidationPx:  stringPtr("25000.0"),
		MarginUsed:     "1000.0",
		MaxTradeSzs:    [2]types.FloatString{"0.1", "0.1"},
		PositionValue:  "5000.0",
		ReturnOnEquity: "5.0",
		Szi:            "0.17857", // Positive = long
		UnrealizedPnl:  "250.0",
	},
	{
		Coin:           "ETH",
		EntryPx:        stringPtr("1800.0"),
		Leverage:       types.LeverageInfo{Type: "isolated", Value: 5, RawUsd: stringPtr("500.0")},
		LiquidationPx:  stringPtr("1500.0"),
		MarginUsed:     "500.0",
		MaxTradeSzs:    [2]types.FloatString{"1.0", "1.0"},
		PositionValue:  "2500.0",
		ReturnOnEquity: "10.0",
		Szi:            "-1.35135", // Negative = short
		UnrealizedPnl:  "-100.0",
	},
	{
		Coin:           "SOL",
		EntryPx:        nil, // No position
		Leverage:       types.LeverageInfo{Type: "cross", Value: 1},
		LiquidationPx:  nil,
		MarginUsed:     "0.0",
		MaxTradeSzs:    [2]types.FloatString{"0.0", "0.0"},
		PositionValue:  "0.0",
		ReturnOnEquity: "0.0",
		Szi:            "0.0", // No position
		UnrealizedPnl:  "0.0",
	},
}

// Test timestamps
const (
	TestTimestamp1 = int64(1677777606040)
	TestTimestamp2 = int64(1677777606050)
	TestTimestamp3 = int64(1677777606060)
)

// Mock market data for different assets
var MockMarketData = map[string]map[string]interface{}{
	"BTC": {
		"markPx":       "28429.3",
		"midPx":        "28430.0",
		"oraclePx":     "28445.0",
		"funding":      "-0.0000886",
		"openInterest": "1234567.89",
		"dayNtlVlm":    "3559323.53447",
		"prevDayPx":    "29368.0",
	},
	"ETH": {
		"markPx":       "1850.5",
		"midPx":        "1851.0",
		"oraclePx":     "1855.0",
		"funding":      "0.0001234",
		"openInterest": "987654.32",
		"dayNtlVlm":    "1234567.89",
		"prevDayPx":    "1890.0",
	},
}

// MockWebSocketSubscriptions for testing
var MockWebSocketSubscriptions = struct {
	AllMids         types.AllMidsSubscription
	L2Book          types.L2BookSubscription
	Trades          types.TradesSubscription
	UserEvents      types.UserEventsSubscription
	UserFills       types.UserFillsSubscription
	Candle          types.CandleSubscription
	OrderUpdates    types.OrderUpdatesSubscription
	UserFundings    types.UserFundingsSubscription
	BBO             types.BboSubscription
	ActiveAssetCtx  types.ActiveAssetCtxSubscription
	ActiveAssetData types.ActiveAssetDataSubscription
}{
	AllMids:         types.AllMidsSubscription{Type: "allMids"},
	L2Book:          types.L2BookSubscription{Type: "l2Book", Coin: "ETH"},
	Trades:          types.TradesSubscription{Type: "trades", Coin: "ETH"},
	UserEvents:      types.UserEventsSubscription{Type: "userEvents", User: TestAddress1},
	UserFills:       types.UserFillsSubscription{Type: "userFills", User: TestAddress1},
	Candle:          types.CandleSubscription{Type: "candle", Coin: "ETH", Interval: "1m"},
	OrderUpdates:    types.OrderUpdatesSubscription{Type: "orderUpdates", User: TestAddress1},
	UserFundings:    types.UserFundingsSubscription{Type: "userFundings", User: TestAddress1},
	BBO:             types.BboSubscription{Type: "bbo", Coin: "ETH"},
	ActiveAssetCtx:  types.ActiveAssetCtxSubscription{Type: "activeAssetCtx", Coin: "ETH"},
	ActiveAssetData: types.ActiveAssetDataSubscription{Type: "activeAssetData", User: TestAddress1, Coin: "ETH"},
}
