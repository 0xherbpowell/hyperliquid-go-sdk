package client

import (
	"context"
	"encoding/json"
	"fmt"

	"hyperliquid-go-sdk/pkg/constants"
	"hyperliquid-go-sdk/pkg/types"
)

// InfoClient provides read-only access to Hyperliquid API
type InfoClient struct {
	*BaseClient
	coinToAsset       map[string]int
	nameToCoin        map[string]string
	assetToSzDecimals map[int]int
}

// NewInfoClient creates a new info client
func NewInfoClient(baseURL string, skipWS bool, meta *types.Meta, spotMeta *types.SpotMeta, perpDexs []string) (*InfoClient, error) {
	baseClient, err := NewBaseClient(baseURL, nil)
	if err != nil {
		return nil, err
	}

	client := &InfoClient{
		BaseClient:        baseClient,
		coinToAsset:       make(map[string]int),
		nameToCoin:        make(map[string]string),
		assetToSzDecimals: make(map[int]int),
	}

	// Initialize metadata
	if err := client.initializeMetadata(meta, spotMeta, perpDexs); err != nil {
		return nil, err
	}

	// Initialize WebSocket if not skipped
	if !skipWS {
		if err := client.initializeWebSocket(); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// initializeMetadata initializes the metadata mappings
func (c *InfoClient) initializeMetadata(meta *types.Meta, spotMeta *types.SpotMeta, perpDexs []string) error {
	// Initialize spot metadata
	if spotMeta == nil {
		var err error
		spotMeta, err = c.SpotMeta(context.Background())
		if err != nil {
			return fmt.Errorf("failed to fetch spot meta: %w", err)
		}
	}

	// Spot assets start at 10000
	for _, spotInfo := range spotMeta.Universe {
		asset := spotInfo.Index + 10000
		c.coinToAsset[spotInfo.Name] = asset
		c.nameToCoin[spotInfo.Name] = spotInfo.Name

		if len(spotInfo.Tokens) >= 1 {
			baseTokenIdx := spotInfo.Tokens[0]
			if baseTokenIdx < len(spotMeta.Tokens) {
				baseInfo := spotMeta.Tokens[baseTokenIdx]
				c.assetToSzDecimals[asset] = baseInfo.SzDecimals

				// Create name mapping for full pair name
				if len(spotInfo.Tokens) >= 2 {
					quoteTokenIdx := spotInfo.Tokens[1]
					if quoteTokenIdx < len(spotMeta.Tokens) {
						quoteInfo := spotMeta.Tokens[quoteTokenIdx]
						fullName := fmt.Sprintf("%s/%s", baseInfo.Name, quoteInfo.Name)
						if _, exists := c.nameToCoin[fullName]; !exists {
							c.nameToCoin[fullName] = spotInfo.Name
						}
					}
				}
			}
		}
	}

	// Initialize perp metadata
	perpDexToOffset := map[string]int{"": 0}
	if perpDexs == nil {
		perpDexs = []string{""}
	} else {
		perpDexsList, err := c.PerpDexs(context.Background())
		if err == nil && len(perpDexsList) > 1 {
			for i, perpDex := range perpDexsList[1:] {
				if dexInfo, ok := perpDex.(map[string]interface{}); ok {
					if name, ok := dexInfo["name"].(string); ok {
						perpDexToOffset[name] = 110000 + i*10000
					}
				}
			}
		}
	}

	for _, perpDex := range perpDexs {
		offset := perpDexToOffset[perpDex]
		var perpMeta *types.Meta
		var err error

		if perpDex == "" && meta != nil {
			perpMeta = meta
		} else {
			perpMeta, err = c.MetaWithDex(context.Background(), perpDex)
			if err != nil {
				continue // Skip if we can't fetch meta for this dex
			}
		}

		c.setPerpMeta(perpMeta, offset)
	}

	return nil
}

// setPerpMeta sets the perp metadata with the given offset
func (c *InfoClient) setPerpMeta(meta *types.Meta, offset int) {
	for asset, assetInfo := range meta.Universe {
		assetID := asset + offset
		c.coinToAsset[assetInfo.Name] = assetID
		c.nameToCoin[assetInfo.Name] = assetInfo.Name
		c.assetToSzDecimals[assetID] = assetInfo.SzDecimals
	}
}

// UserState retrieves trading details about a user
func (c *InfoClient) UserState(ctx context.Context, address string, dex string) (*types.UserState, error) {
	req := map[string]interface{}{
		"type": constants.InfoClearinghouseState,
		"user": address,
	}
	if dex != "" {
		req["dex"] = dex
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var userState types.UserState
	if err := json.Unmarshal(resp, &userState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user state: %w", err)
	}

	return &userState, nil
}

// SpotUserState retrieves user's spot trading state
func (c *InfoClient) SpotUserState(ctx context.Context, address string) (*types.SpotUserState, error) {
	req := map[string]interface{}{
		"type": constants.InfoSpotClearinghouseState,
		"user": address,
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var spotUserState types.SpotUserState
	if err := json.Unmarshal(resp, &spotUserState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spot user state: %w", err)
	}

	return &spotUserState, nil
}

// OpenOrders retrieves a user's open orders
func (c *InfoClient) OpenOrders(ctx context.Context, address string, dex string) ([]types.OpenOrder, error) {
	req := map[string]interface{}{
		"type": constants.InfoOpenOrders,
		"user": address,
	}
	if dex != "" {
		req["dex"] = dex
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var openOrders []types.OpenOrder
	if err := json.Unmarshal(resp, &openOrders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal open orders: %w", err)
	}

	return openOrders, nil
}

// FrontendOpenOrders retrieves a user's open orders with additional frontend info
func (c *InfoClient) FrontendOpenOrders(ctx context.Context, address string, dex string) ([]types.FrontendOpenOrder, error) {
	req := map[string]interface{}{
		"type": constants.InfoFrontendOpenOrders,
		"user": address,
	}
	if dex != "" {
		req["dex"] = dex
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var frontendOpenOrders []types.FrontendOpenOrder
	if err := json.Unmarshal(resp, &frontendOpenOrders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal frontend open orders: %w", err)
	}

	return frontendOpenOrders, nil
}

// AllMids retrieves all mids for all actively traded coins
func (c *InfoClient) AllMids(ctx context.Context, dex string) (types.AllMids, error) {
	req := map[string]interface{}{
		"type": constants.InfoAllMids,
	}
	if dex != "" {
		req["dex"] = dex
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var allMids types.AllMids
	if err := json.Unmarshal(resp, &allMids); err != nil {
		return nil, fmt.Errorf("failed to unmarshal all mids: %w", err)
	}

	return allMids, nil
}

// UserFills retrieves a given user's fills
func (c *InfoClient) UserFills(ctx context.Context, address string) ([]types.Fill, error) {
	req := map[string]interface{}{
		"type": constants.InfoUserFills,
		"user": address,
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var fills []types.Fill
	if err := json.Unmarshal(resp, &fills); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user fills: %w", err)
	}

	return fills, nil
}

// UserFillsByTime retrieves a user's fills by time
func (c *InfoClient) UserFillsByTime(ctx context.Context, address string, startTime int64, endTime *int64, aggregateByTime *bool) ([]types.Fill, error) {
	req := map[string]interface{}{
		"type":      constants.InfoUserFillsByTime,
		"user":      address,
		"startTime": startTime,
	}
	if endTime != nil {
		req["endTime"] = *endTime
	}
	if aggregateByTime != nil {
		req["aggregateByTime"] = *aggregateByTime
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var fills []types.Fill
	if err := json.Unmarshal(resp, &fills); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user fills by time: %w", err)
	}

	return fills, nil
}

// Meta retrieves exchange perp metadata
func (c *InfoClient) Meta(ctx context.Context) (*types.Meta, error) {
	return c.MetaWithDex(ctx, "")
}

// MetaWithDex retrieves exchange perp metadata for a specific dex
func (c *InfoClient) MetaWithDex(ctx context.Context, dex string) (*types.Meta, error) {
	req := map[string]interface{}{
		"type": constants.InfoMeta,
	}
	if dex != "" {
		req["dex"] = dex
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var meta types.Meta
	if err := json.Unmarshal(resp, &meta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meta: %w", err)
	}

	return &meta, nil
}

// MetaAndAssetCtxs retrieves exchange metadata and asset contexts
func (c *InfoClient) MetaAndAssetCtxs(ctx context.Context) (interface{}, error) {
	req := map[string]interface{}{
		"type": constants.InfoMetaAndAssetCtxs,
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meta and asset ctxs: %w", err)
	}

	return result, nil
}

// PerpDexs retrieves available perp dexs
func (c *InfoClient) PerpDexs(ctx context.Context) ([]interface{}, error) {
	req := map[string]interface{}{
		"type": constants.InfoPerpDexs,
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var perpDexs []interface{}
	if err := json.Unmarshal(resp, &perpDexs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal perp dexs: %w", err)
	}

	return perpDexs, nil
}

// SpotMeta retrieves exchange spot metadata
func (c *InfoClient) SpotMeta(ctx context.Context) (*types.SpotMeta, error) {
	req := map[string]interface{}{
		"type": constants.InfoSpotMeta,
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var spotMeta types.SpotMeta
	if err := json.Unmarshal(resp, &spotMeta); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spot meta: %w", err)
	}

	return &spotMeta, nil
}

// SpotMetaAndAssetCtxs retrieves exchange spot metadata and asset contexts
func (c *InfoClient) SpotMetaAndAssetCtxs(ctx context.Context) (*types.SpotMetaAndAssetCtxs, error) {
	req := map[string]interface{}{
		"type": constants.InfoSpotMetaAndAssetCtxs,
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var result types.SpotMetaAndAssetCtxs
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spot meta and asset ctxs: %w", err)
	}

	return &result, nil
}

// FundingHistory retrieves funding history for a given coin
func (c *InfoClient) FundingHistory(ctx context.Context, name string, startTime int64, endTime *int64) (interface{}, error) {
	coin := c.nameToCoin[name]
	req := map[string]interface{}{
		"type":      constants.InfoFundingHistory,
		"coin":      coin,
		"startTime": startTime,
	}
	if endTime != nil {
		req["endTime"] = *endTime
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal funding history: %w", err)
	}

	return result, nil
}

// UserFundingHistory retrieves a user's funding history
func (c *InfoClient) UserFundingHistory(ctx context.Context, user string, startTime int64, endTime *int64) (interface{}, error) {
	req := map[string]interface{}{
		"type":      constants.InfoUserFunding,
		"user":      user,
		"startTime": startTime,
	}
	if endTime != nil {
		req["endTime"] = *endTime
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user funding history: %w", err)
	}

	return result, nil
}

// L2Snapshot retrieves L2 snapshot for a given coin
func (c *InfoClient) L2Snapshot(ctx context.Context, name string) (*types.L2BookData, error) {
	coin := c.nameToCoin[name]
	req := map[string]interface{}{
		"type": constants.InfoL2Book,
		"coin": coin,
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var l2BookData types.L2BookData
	if err := json.Unmarshal(resp, &l2BookData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal L2 snapshot: %w", err)
	}

	return &l2BookData, nil
}

// CandlesSnapshot retrieves candles snapshot for a given coin
func (c *InfoClient) CandlesSnapshot(ctx context.Context, name, interval string, startTime, endTime int64) (interface{}, error) {
	coin := c.nameToCoin[name]
	req := map[string]interface{}{
		"type": constants.InfoCandleSnapshot,
		"req": map[string]interface{}{
			"coin":      coin,
			"interval":  interval,
			"startTime": startTime,
			"endTime":   endTime,
		},
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal candles snapshot: %w", err)
	}

	return result, nil
}

// UserFees retrieves the volume of trading activity associated with a user
func (c *InfoClient) UserFees(ctx context.Context, address string) (interface{}, error) {
	req := map[string]interface{}{
		"type": constants.InfoUserFees,
		"user": address,
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user fees: %w", err)
	}

	return result, nil
}

// UserStakingSummary retrieves the staking summary associated with a user
func (c *InfoClient) UserStakingSummary(ctx context.Context, address string) (interface{}, error) {
	req := map[string]interface{}{
		"type": constants.InfoDelegatorSummary,
		"user": address,
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user staking summary: %w", err)
	}

	return result, nil
}

// QueryOrderByOid retrieves order status by order ID
func (c *InfoClient) QueryOrderByOid(ctx context.Context, user string, oid int) (interface{}, error) {
	req := map[string]interface{}{
		"type": constants.InfoOrderStatus,
		"user": user,
		"oid":  oid,
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order status: %w", err)
	}

	return result, nil
}

// QueryOrderByCloid retrieves order status by client order ID
func (c *InfoClient) QueryOrderByCloid(ctx context.Context, user string, cloid *types.Cloid) (interface{}, error) {
	req := map[string]interface{}{
		"type": constants.InfoOrderStatus,
		"user": user,
		"oid":  cloid.Raw(),
	}

	resp, err := c.apiClient.Post(ctx, "/info", req)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order status: %w", err)
	}

	return result, nil
}

// NameToAsset converts name to asset ID
func (c *InfoClient) NameToAsset(name string) int {
	if coin, exists := c.nameToCoin[name]; exists {
		if asset, exists := c.coinToAsset[coin]; exists {
			return asset
		}
	}
	return 0
}

// Helper methods for WebSocket subscriptions (if WebSocket is enabled)

// initializeWebSocket initializes the WebSocket connection
func (c *InfoClient) initializeWebSocket() error {
	// WebSocket implementation would go here
	// For now, we'll skip this as it's a complex implementation
	return nil
}
