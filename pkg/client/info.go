package client

import (
	"fmt"
	"time"

	"hyperliquid-go-sdk/pkg/types"
)

// Info provides methods to query market data and information
type Info struct {
	*API
	coinToAsset       map[string]int
	nameToCoin        map[string]string
	assetToSzDecimals map[int]int
	wsManager         *WebsocketManager
}

// NewInfo creates a new Info client
func NewInfo(baseURL string, timeout *time.Duration, skipWS bool, meta *types.Meta, spotMeta *types.SpotMeta, perpDexs []string) (*Info, error) {
	api := NewAPI(baseURL, timeout)

	info := &Info{
		API:               api,
		coinToAsset:       make(map[string]int),
		nameToCoin:        make(map[string]string),
		assetToSzDecimals: make(map[int]int),
	}

	// Initialize WebSocket manager if not skipped
	if !skipWS {
		wsManager, err := NewWebsocketManager(api.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create websocket manager: %w", err)
		}
		info.wsManager = wsManager
		if err := info.wsManager.Start(); err != nil {
			return nil, fmt.Errorf("failed to start websocket manager: %w", err)
		}
	}

	// Initialize spot meta
	if spotMeta == nil {
		var err error
		spotMeta, err = info.SpotMeta()
		if err != nil {
			return nil, fmt.Errorf("failed to get spot meta: %w", err)
		}
	}

	// Initialize spot assets (start at 10000)
	for _, spotInfo := range spotMeta.Universe {
		asset := spotInfo.Index + 10000
		info.coinToAsset[spotInfo.Name] = asset
		info.nameToCoin[spotInfo.Name] = spotInfo.Name

		if len(spotInfo.Tokens) >= 2 {
			base := spotInfo.Tokens[0]
			quote := spotInfo.Tokens[1]

			if base < len(spotMeta.Tokens) && quote < len(spotMeta.Tokens) {
				baseInfo := spotMeta.Tokens[base]
				quoteInfo := spotMeta.Tokens[quote]
				info.assetToSzDecimals[asset] = baseInfo.SzDecimals

				name := fmt.Sprintf("%s/%s", baseInfo.Name, quoteInfo.Name)
				if _, exists := info.nameToCoin[name]; !exists {
					info.nameToCoin[name] = spotInfo.Name
				}
			}
		}
	}

	// Initialize perp dex mappings
	perpDexToOffset := map[string]int{"": 0}

	if perpDexs == nil {
		perpDexs = []string{""}
	} else {
		perpDexsList, err := info.PerpDexs()
		if err != nil {
			return nil, fmt.Errorf("failed to get perp dexs: %w", err)
		}

		for i, perpDex := range perpDexsList[1:] {
			// builder-deployed perp dexs start at 110000
			if perpDexMap, ok := perpDex.(map[string]interface{}); ok {
				if name, ok := perpDexMap["name"].(string); ok {
					perpDexToOffset[name] = 110000 + i*10000
				}
			}
		}
	}

	// Initialize perp assets
	for _, perpDex := range perpDexs {
		offset := perpDexToOffset[perpDex]

		var perpMeta *types.Meta
		var err error

		if perpDex == "" && meta != nil {
			perpMeta = meta
		} else {
			perpMeta, err = info.Meta(perpDex)
			if err != nil {
				return nil, fmt.Errorf("failed to get meta for dex %s: %w", perpDex, err)
			}
		}

		info.setPerpMeta(perpMeta, offset)
	}

	return info, nil
}

// setPerpMeta sets the perpetual asset metadata
func (i *Info) setPerpMeta(meta *types.Meta, offset int) {
	for asset, assetInfo := range meta.Universe {
		actualAsset := asset + offset
		i.coinToAsset[assetInfo.Name] = actualAsset
		i.nameToCoin[assetInfo.Name] = assetInfo.Name
		i.assetToSzDecimals[actualAsset] = assetInfo.SzDecimals
	}
}

// DisconnectWebsocket disconnects the WebSocket connection
func (i *Info) DisconnectWebsocket() error {
	if i.wsManager == nil {
		return fmt.Errorf("cannot call disconnect_websocket since skip_ws was used")
	}
	return i.wsManager.Stop()
}

// NameToAsset converts asset name to asset ID
func (i *Info) NameToAsset(name string) (int, error) {
	if coin, exists := i.nameToCoin[name]; exists {
		if asset, exists := i.coinToAsset[coin]; exists {
			return asset, nil
		}
	}
	return 0, fmt.Errorf("asset not found: %s", name)
}

// UserState retrieves trading details about a user
func (i *Info) UserState(address string, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "clearinghouseState",
		"user": address,
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// OpenOrders retrieves a user's open orders
func (i *Info) OpenOrders(address string, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "openOrders",
		"user": address,
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// FrontendOpenOrders retrieves a user's open orders with additional frontend data
func (i *Info) FrontendOpenOrders(address string, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "frontendOpenOrders",
		"user": address,
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// UserFills retrieves a user's fills
func (i *Info) UserFills(address string, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "userFills",
		"user": address,
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// UserFillsByTime retrieves a user's fills within a time range
func (i *Info) UserFillsByTime(address string, startTime int64, endTime *int64, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type":      "userFillsByTime",
		"user":      address,
		"startTime": startTime,
	}

	if endTime != nil {
		payload["endTime"] = *endTime
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// UserNonFundingLedgerUpdates retrieves a user's non-funding ledger updates
func (i *Info) UserNonFundingLedgerUpdates(address string, startTime int64, endTime *int64, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type":      "userNonFundingLedgerUpdates",
		"user":      address,
		"startTime": startTime,
	}

	if endTime != nil {
		payload["endTime"] = *endTime
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// UserFunding retrieves a user's funding history
func (i *Info) UserFunding(address string, startTime int64, endTime *int64, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type":      "userFunding",
		"user":      address,
		"startTime": startTime,
	}

	if endTime != nil {
		payload["endTime"] = *endTime
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// UserRateLimit retrieves a user's rate limit information
func (i *Info) UserRateLimit(address string, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "userRateLimit",
		"user": address,
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// OrderStatus retrieves the status of an order
func (i *Info) OrderStatus(address string, oid int, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "orderStatus",
		"user": address,
		"oid":  oid,
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// L2Book retrieves the L2 order book for an asset
func (i *Info) L2Book(coin string, dex string, nSigFigs *int, mantissa *int) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "l2Book",
		"coin": coin,
	}

	if dex != "" {
		payload["dex"] = dex
	}

	if nSigFigs != nil {
		payload["nSigFigs"] = *nSigFigs
	}

	if mantissa != nil {
		payload["mantissa"] = *mantissa
	}

	return i.Post("/info", payload)
}

// RecentTrades retrieves recent trades for an asset
func (i *Info) RecentTrades(coin string, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "recentTrades",
		"coin": coin,
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// AllMids retrieves mid prices for all assets
func (i *Info) AllMids(dex string) (map[string]string, error) {
	payload := map[string]interface{}{
		"type": "allMids",
	}

	if dex != "" {
		payload["dex"] = dex
	}

	result, err := i.Post("/info", payload)
	if err != nil {
		return nil, err
	}

	mids := make(map[string]string)
	// The API response directly contains the price data, not wrapped in a 'mids' key
	for k, v := range result {
		if str, ok := v.(string); ok {
			mids[k] = str
		}
	}

	return mids, nil
}

// UserTradesHistory retrieves a user's trade history
func (i *Info) UserTradesHistory(address string, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type": "userTradesHistory",
		"user": address,
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// Meta retrieves the universe of perpetual assets
func (i *Info) Meta(dex string) (*types.Meta, error) {
	payload := map[string]interface{}{
		"type": "meta",
	}

	if dex != "" {
		payload["dex"] = dex
	}

	result, err := i.Post("/info", payload)
	if err != nil {
		return nil, err
	}

	var meta types.Meta
	if universe, ok := result["universe"].([]interface{}); ok {
		for _, item := range universe {
			if assetMap, ok := item.(map[string]interface{}); ok {
				var asset types.AssetInfo
				if name, ok := assetMap["name"].(string); ok {
					asset.Name = name
				}
				if szDecimals, ok := assetMap["szDecimals"].(float64); ok {
					asset.SzDecimals = int(szDecimals)
				}
				meta.Universe = append(meta.Universe, asset)
			}
		}
	}

	return &meta, nil
}

// SpotMeta retrieves the universe of spot assets
func (i *Info) SpotMeta() (*types.SpotMeta, error) {
	payload := map[string]interface{}{
		"type": "spotMeta",
	}

	result, err := i.Post("/info", payload)
	if err != nil {
		return nil, err
	}

	var spotMeta types.SpotMeta

	// Parse universe
	if universe, ok := result["universe"].([]interface{}); ok {
		for _, item := range universe {
			if assetMap, ok := item.(map[string]interface{}); ok {
				var asset types.SpotAssetInfo

				if name, ok := assetMap["name"].(string); ok {
					asset.Name = name
				}
				if index, ok := assetMap["index"].(float64); ok {
					asset.Index = int(index)
				}
				if isCanonical, ok := assetMap["isCanonical"].(bool); ok {
					asset.IsCanonical = isCanonical
				}
				if tokens, ok := assetMap["tokens"].([]interface{}); ok {
					for _, token := range tokens {
						if tokenInt, ok := token.(float64); ok {
							asset.Tokens = append(asset.Tokens, int(tokenInt))
						}
					}
				}

				spotMeta.Universe = append(spotMeta.Universe, asset)
			}
		}
	}

	// Parse tokens
	if tokens, ok := result["tokens"].([]interface{}); ok {
		for _, item := range tokens {
			if tokenMap, ok := item.(map[string]interface{}); ok {
				var token types.SpotTokenInfo

				if name, ok := tokenMap["name"].(string); ok {
					token.Name = name
				}
				if szDecimals, ok := tokenMap["szDecimals"].(float64); ok {
					token.SzDecimals = int(szDecimals)
				}
				if weiDecimals, ok := tokenMap["weiDecimals"].(float64); ok {
					token.WeiDecimals = int(weiDecimals)
				}
				if index, ok := tokenMap["index"].(float64); ok {
					token.Index = int(index)
				}
				if tokenId, ok := tokenMap["tokenId"].(string); ok {
					token.TokenId = tokenId
				}
				if isCanonical, ok := tokenMap["isCanonical"].(bool); ok {
					token.IsCanonical = isCanonical
				}
				if evmContract, ok := tokenMap["evmContract"].(string); ok && evmContract != "" {
					token.EvmContract = &evmContract
				}
				if fullName, ok := tokenMap["fullName"].(string); ok && fullName != "" {
					token.FullName = &fullName
				}

				spotMeta.Tokens = append(spotMeta.Tokens, token)
			}
		}
	}

	return &spotMeta, nil
}

// PerpDexs retrieves the list of perpetual dexes
func (i *Info) PerpDexs() ([]interface{}, error) {
	payload := map[string]interface{}{
		"type": "perpDexs",
	}

	result, err := i.Post("/info", payload)
	if err != nil {
		return nil, err
	}

	if dexs, ok := result["dexs"].([]interface{}); ok {
		return dexs, nil
	}

	return []interface{}{}, nil
}

// ClearinghouseState retrieves clearinghouse state
func (i *Info) ClearinghouseState(address string, dex string) (map[string]interface{}, error) {
	return i.UserState(address, dex)
}

// BatchUserStates retrieves user states for multiple addresses
func (i *Info) BatchUserStates(addresses []string, dex string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"type":  "batchUserStates",
		"users": addresses,
	}

	if dex != "" {
		payload["dex"] = dex
	}

	return i.Post("/info", payload)
}

// Subscribe subscribes to WebSocket channels (if WebSocket is enabled)
func (i *Info) Subscribe(subscriptions []types.Subscription, callback func(interface{})) error {
	if i.wsManager == nil {
		return fmt.Errorf("WebSocket manager not available (skip_ws was used)")
	}

	return i.wsManager.Subscribe(subscriptions, callback)
}

// Unsubscribe unsubscribes from WebSocket channels (if WebSocket is enabled)
func (i *Info) Unsubscribe(subscriptions []types.Subscription) error {
	if i.wsManager == nil {
		return fmt.Errorf("WebSocket manager not available (skip_ws was used)")
	}

	return i.wsManager.Unsubscribe(subscriptions)
}
