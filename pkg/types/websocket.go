package types

// WebSocket subscription types

// AllMidsSubscription represents subscription to all mid prices
type AllMidsSubscription struct {
	Type string `json:"type"` // "allMids"
}

// BboSubscription represents subscription to best bid/offer
type BboSubscription struct {
	Type string `json:"type"` // "bbo"
	Coin string `json:"coin"`
}

// L2BookSubscription represents subscription to L2 order book
type L2BookSubscription struct {
	Type string `json:"type"` // "l2Book"
	Coin string `json:"coin"`
}

// TradesSubscription represents subscription to trades
type TradesSubscription struct {
	Type string `json:"type"` // "trades"
	Coin string `json:"coin"`
}

// UserEventsSubscription represents subscription to user events
type UserEventsSubscription struct {
	Type string `json:"type"` // "userEvents"
	User string `json:"user"`
}

// UserFillsSubscription represents subscription to user fills
type UserFillsSubscription struct {
	Type string `json:"type"` // "userFills"
	User string `json:"user"`
}

// CandleSubscription represents subscription to candle data
type CandleSubscription struct {
	Type     string `json:"type"` // "candle"
	Coin     string `json:"coin"`
	Interval string `json:"interval"`
}

// OrderUpdatesSubscription represents subscription to order updates
type OrderUpdatesSubscription struct {
	Type string `json:"type"` // "orderUpdates"
	User string `json:"user"`
}

// UserFundingsSubscription represents subscription to user funding updates
type UserFundingsSubscription struct {
	Type string `json:"type"` // "userFundings"
	User string `json:"user"`
}

// UserNonFundingLedgerUpdatesSubscription represents subscription to non-funding ledger updates
type UserNonFundingLedgerUpdatesSubscription struct {
	Type string `json:"type"` // "userNonFundingLedgerUpdates"
	User string `json:"user"`
}

// WebData2Subscription represents subscription to web data 2
type WebData2Subscription struct {
	Type string `json:"type"` // "webData2"
	User string `json:"user"`
}

// ActiveAssetCtxSubscription represents subscription to active asset context
type ActiveAssetCtxSubscription struct {
	Type string `json:"type"` // "activeAssetCtx"
	Coin string `json:"coin"`
}

// ActiveAssetDataSubscription represents subscription to active asset data
type ActiveAssetDataSubscription struct {
	Type string `json:"type"` // "activeAssetData"
	User string `json:"user"`
	Coin string `json:"coin"`
}

// Subscription represents any subscription type
type Subscription interface {
	GetType() string
}

func (s AllMidsSubscription) GetType() string                     { return s.Type }
func (s BboSubscription) GetType() string                         { return s.Type }
func (s L2BookSubscription) GetType() string                      { return s.Type }
func (s TradesSubscription) GetType() string                      { return s.Type }
func (s UserEventsSubscription) GetType() string                  { return s.Type }
func (s UserFillsSubscription) GetType() string                   { return s.Type }
func (s CandleSubscription) GetType() string                      { return s.Type }
func (s OrderUpdatesSubscription) GetType() string                { return s.Type }
func (s UserFundingsSubscription) GetType() string                { return s.Type }
func (s UserNonFundingLedgerUpdatesSubscription) GetType() string { return s.Type }
func (s WebData2Subscription) GetType() string                    { return s.Type }
func (s ActiveAssetCtxSubscription) GetType() string              { return s.Type }
func (s ActiveAssetDataSubscription) GetType() string             { return s.Type }

// WebSocket message types

// AllMidsData represents all mid prices data
type AllMidsData struct {
	Mids map[string]FloatString `json:"mids"`
}

// AllMidsMsg represents all mids message
type AllMidsMsg struct {
	Channel string      `json:"channel"` // "allMids"
	Data    AllMidsData `json:"data"`
}

// L2BookMsg represents L2 book message
type L2BookMsg struct {
	Channel string     `json:"channel"` // "l2Book"
	Data    L2BookData `json:"data"`
}

// BboData represents best bid/offer data
type BboData struct {
	Coin string      `json:"coin"`
	Time int64       `json:"time"`
	Bbo  [2]*L2Level `json:"bbo"` // [bid, ask]
}

// BboMsg represents BBO message
type BboMsg struct {
	Channel string  `json:"channel"` // "bbo"
	Data    BboData `json:"data"`
}

// PongMsg represents pong message
type PongMsg struct {
	Channel string `json:"channel"` // "pong"
}

// TradesMsg represents trades message
type TradesMsg struct {
	Channel string  `json:"channel"` // "trades"
	Data    []Trade `json:"data"`
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

// ActiveAssetCtxMsg represents active asset context message
type ActiveAssetCtxMsg struct {
	Channel string         `json:"channel"` // "activeAssetCtx"
	Data    ActiveAssetCtx `json:"data"`
}

// ActiveSpotAssetCtxMsg represents active spot asset context message
type ActiveSpotAssetCtxMsg struct {
	Channel string             `json:"channel"` // "activeSpotAssetCtx"
	Data    ActiveSpotAssetCtx `json:"data"`
}

// ActiveAssetData represents active asset data
type ActiveAssetData struct {
	User             string         `json:"user"`
	Coin             string         `json:"coin"`
	Leverage         Leverage       `json:"leverage"`
	MaxTradeSzs      [2]FloatString `json:"maxTradeSzs"`
	AvailableToTrade [2]FloatString `json:"availableToTrade"`
	MarkPx           FloatString    `json:"markPx"`
}

// ActiveAssetDataMsg represents active asset data message
type ActiveAssetDataMsg struct {
	Channel string          `json:"channel"` // "activeAssetData"
	Data    ActiveAssetData `json:"data"`
}

// UserEventsData represents user events data
type UserEventsData struct {
	Fills []Fill `json:"fills,omitempty"`
}

// UserEventsMsg represents user events message
type UserEventsMsg struct {
	Channel string         `json:"channel"` // "user"
	Data    UserEventsData `json:"data"`
}

// UserFillsData represents user fills data
type UserFillsData struct {
	User       string `json:"user"`
	IsSnapshot bool   `json:"isSnapshot"`
	Fills      []Fill `json:"fills"`
}

// UserFillsMsg represents user fills message
type UserFillsMsg struct {
	Channel string        `json:"channel"` // "userFills"
	Data    UserFillsData `json:"data"`
}

// OtherWsMsg represents other WebSocket messages
type OtherWsMsg struct {
	Channel string      `json:"channel"`
	Data    interface{} `json:"data,omitempty"`
}

// WsMsg represents any WebSocket message
type WsMsg interface {
	GetChannel() string
}

func (m AllMidsMsg) GetChannel() string            { return m.Channel }
func (m BboMsg) GetChannel() string                { return m.Channel }
func (m L2BookMsg) GetChannel() string             { return m.Channel }
func (m TradesMsg) GetChannel() string             { return m.Channel }
func (m UserEventsMsg) GetChannel() string         { return m.Channel }
func (m PongMsg) GetChannel() string               { return m.Channel }
func (m UserFillsMsg) GetChannel() string          { return m.Channel }
func (m ActiveAssetCtxMsg) GetChannel() string     { return m.Channel }
func (m ActiveSpotAssetCtxMsg) GetChannel() string { return m.Channel }
func (m ActiveAssetDataMsg) GetChannel() string    { return m.Channel }
func (m OtherWsMsg) GetChannel() string            { return m.Channel }

// WebSocket request/response types

// WsRequest represents a WebSocket request
type WsRequest struct {
	Method       string       `json:"method"`
	Subscription Subscription `json:"subscription,omitempty"`
}

// SubscriptionCallback represents a callback function for subscription data
type SubscriptionCallback func(data interface{})

// ActiveSubscription represents an active WebSocket subscription
type ActiveSubscription struct {
	Callback       SubscriptionCallback
	SubscriptionID int
}
