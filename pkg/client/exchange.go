package client

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
)

const DefaultSlippage = 0.05 // 5% default slippage for market orders

// Exchange provides methods for trading operations
type Exchange struct {
	*API
	privateKey     *ecdsa.PrivateKey
	vaultAddress   *string
	accountAddress *string
	info           *Info
	expiresAfter   *int64
}

// NewExchange creates a new Exchange client
func NewExchange(
	privateKey *ecdsa.PrivateKey,
	baseURL string,
	timeout *time.Duration,
	meta *types.Meta,
	vaultAddress *string,
	accountAddress *string,
	spotMeta *types.SpotMeta,
	perpDexs []string,
) (*Exchange, error) {
	api := NewAPI(baseURL, timeout)

	// Create info client with skipWS=true for exchange
	info, err := NewInfo(baseURL, timeout, true, meta, spotMeta, perpDexs)
	if err != nil {
		return nil, fmt.Errorf("failed to create info client: %w", err)
	}

	return &Exchange{
		API:            api,
		privateKey:     privateKey,
		vaultAddress:   vaultAddress,
		accountAddress: accountAddress,
		info:           info,
	}, nil
}

// SetExpiresAfter sets the expiration time for actions
func (e *Exchange) SetExpiresAfter(expiresAfter *int64) {
	e.expiresAfter = expiresAfter
}

// postAction posts an action to the exchange
func (e *Exchange) postAction(action map[string]interface{}, signature map[string]interface{}, nonce int64) (map[string]interface{}, error) {
	var vaultAddress *string
	// Only add vaultAddress for certain action types
	actionType, ok := action["type"].(string)
	if ok && actionType != "usdClassTransfer" && actionType != "sendAsset" {
		vaultAddress = e.vaultAddress
	} else {
		vaultAddress = nil
	}

	payload := map[string]interface{}{
		"action":       action,
		"nonce":        nonce,
		"signature":    signature,
		"vaultAddress": vaultAddress,
		"expiresAfter": e.expiresAfter, // Always include, even if nil (like Python SDK)
	}
	
	// TODO: Check if isFrontend is needed
	// Temporarily removing to test
	// if actionType == "order" || actionType == "cancel" || actionType == "cancelAll" || actionType == "modify" {
	// 	payload["isFrontend"] = true
	// }
	
	// Handle agent mode: if account address differs from wallet address, include user field
	walletAddress := utils.GetAddressFromPrivateKey(e.privateKey)
	if e.accountAddress != nil && *e.accountAddress != walletAddress {
		// Agent mode: signing with agent key for account
		payload["user"] = *e.accountAddress
	}
	
	log.Println("Payload:", payload)
	return e.Post("/exchange", payload)
}

// slippagePrice calculates the price with slippage
func (e *Exchange) slippagePrice(name string, isBuy bool, slippage float64, px *float64) (float64, error) {
	coin, exists := e.info.nameToCoin[name]
	if !exists {
		return 0, fmt.Errorf("coin not found: %s", name)
	}

	var price float64
	if px != nil {
		price = *px
	} else {
		// Get mid price
		mids, err := e.info.AllMids("")
		if err != nil {
			return 0, fmt.Errorf("failed to get mids: %w", err)
		}

		midStr, exists := mids[coin]
		if !exists {
			return 0, fmt.Errorf("mid price not found for coin: %s", coin)
		}

		price, err = strconv.ParseFloat(midStr, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse mid price: %w", err)
		}
	}

	asset, exists := e.info.coinToAsset[coin]
	if !exists {
		return 0, fmt.Errorf("asset not found for coin: %s", coin)
	}

	// spot assets start at 10000
	isSpot := asset >= 10000

	// Calculate slippage
	if isBuy {
		price *= (1 + slippage)
	} else {
		price *= (1 - slippage)
	}

	// Round to appropriate decimal places
	var decimals int
	if isSpot {
		szDecimals, exists := e.info.assetToSzDecimals[asset]
		if exists {
			decimals = 8 - szDecimals
		} else {
			decimals = 8
		}
	} else {
		szDecimals, exists := e.info.assetToSzDecimals[asset]
		if exists {
			decimals = 6 - szDecimals
		} else {
			decimals = 6
		}
	}

	// Round to 5 significant figures and appropriate decimal places
	sigFigs := 5
	magnitude := math.Log10(math.Abs(price))
	roundTo := math.Max(float64(decimals), float64(sigFigs)-magnitude-1)

	multiplier := math.Pow(10, roundTo)
	return math.Round(price*multiplier) / multiplier, nil
}

// Order places a single order
func (e *Exchange) Order(
	name string,
	isBuy bool,
	sz float64,
	limitPx float64,
	orderType types.OrderType,
	reduceOnly bool,
	cloid *types.Cloid,
	builder *types.BuilderInfo,
) (map[string]interface{}, error) {
	order := types.OrderRequest{
		Coin:       name,
		IsBuy:      isBuy,
		Sz:         sz,
		LimitPx:    limitPx,
		OrderType:  orderType,
		ReduceOnly: reduceOnly,
		Cloid:      cloid,
	}

	return e.BulkOrders([]types.OrderRequest{order}, builder)
}

// BulkOrders places multiple orders in a single transaction
func (e *Exchange) BulkOrders(orderRequests []types.OrderRequest, builder *types.BuilderInfo) (map[string]interface{}, error) {
	var orderWires []types.OrderWire

	for _, order := range orderRequests {
		asset, err := e.info.NameToAsset(order.Coin)
		if err != nil {
			return nil, fmt.Errorf("failed to get asset for coin %s: %w", order.Coin, err)
		}

		orderWire, err := utils.OrderRequestToOrderWire(order, asset)
		if err != nil {
			return nil, fmt.Errorf("failed to convert order to wire format: %w", err)
		}

		orderWires = append(orderWires, orderWire)
	}

	timestamp := utils.GetTimestampMS()

	// Normalize builder address to lowercase
	if builder != nil {
		builder.B = strings.ToLower(builder.B)
	}

	orderAction := utils.OrderWiresToOrderAction(orderWires, builder)

	// Use SignL1ActionWithAccount to handle agent mode properly
	signature, err := utils.SignL1ActionWithAccount(
		e.privateKey,
		orderAction,
		e.vaultAddress,
		timestamp,
		e.expiresAfter,
		e.IsMainnet(),
		e.accountAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign order action: %w", err)
	}

	return e.postAction(orderAction, signature, timestamp)
}

// MarketOrder places a market order with slippage protection
func (e *Exchange) MarketOrder(
	name string,
	isBuy bool,
	sz float64,
	slippage *float64,
	cloid *types.Cloid,
) (map[string]interface{}, error) {
	if slippage == nil {
		defaultSlippage := DefaultSlippage
		slippage = &defaultSlippage
	}

	limitPx, err := e.slippagePrice(name, isBuy, *slippage, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate slippage price: %w", err)
	}

	orderType := types.OrderType{
		Limit: &types.LimitOrderType{
			Tif: types.TifIoc, // Immediate or cancel for market orders
		},
	}

	return e.Order(name, isBuy, sz, limitPx, orderType, false, cloid, nil)
}

// LimitOrder places a limit order
func (e *Exchange) LimitOrder(
	name string,
	isBuy bool,
	sz float64,
	limitPx float64,
	tif types.Tif,
	reduceOnly bool,
	cloid *types.Cloid,
) (map[string]interface{}, error) {
	orderType := types.OrderType{
		Limit: &types.LimitOrderType{
			Tif: tif,
		},
	}

	return e.Order(name, isBuy, sz, limitPx, orderType, reduceOnly, cloid, nil)
}

// TriggerOrder places a trigger order (stop loss or take profit)
func (e *Exchange) TriggerOrder(
	name string,
	isBuy bool,
	sz float64,
	triggerPx float64,
	isMarket bool,
	tpsl types.Tpsl,
	reduceOnly bool,
	cloid *types.Cloid,
) (map[string]interface{}, error) {
	orderType := types.OrderType{
		Trigger: &types.TriggerOrderType{
			TriggerPx: triggerPx,
			IsMarket:  isMarket,
			Tpsl:      tpsl,
		},
	}

	// For trigger orders, limit price should be the trigger price
	limitPx := triggerPx
	if isMarket {
		// For market trigger orders, use slippage protection
		var err error
		limitPx, err = e.slippagePrice(name, isBuy, DefaultSlippage, &triggerPx)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate trigger slippage price: %w", err)
		}
	}

	return e.Order(name, isBuy, sz, limitPx, orderType, reduceOnly, cloid, nil)
}

// Cancel cancels an order by order ID
func (e *Exchange) Cancel(coin string, oid int) (map[string]interface{}, error) {
	asset, err := e.info.NameToAsset(coin)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset for coin %s: %w", coin, err)
	}

	timestamp := utils.GetTimestampMS()

	action := map[string]interface{}{
		"type": "cancel",
		"cancels": []map[string]interface{}{
			{
				"a":   asset,
				"oid": oid,
			},
		},
	}

	signature, err := utils.SignL1ActionWithAccount(
		e.privateKey,
		action,
		e.vaultAddress,
		timestamp,
		e.expiresAfter,
		e.IsMainnet(),
		e.accountAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign cancel action: %w", err)
	}

	return e.postAction(action, signature, timestamp)
}

// CancelByCloid cancels an order by client order ID
func (e *Exchange) CancelByCloid(coin string, cloid *types.Cloid) (map[string]interface{}, error) {
	asset, err := e.info.NameToAsset(coin)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset for coin %s: %w", coin, err)
	}

	timestamp := utils.GetTimestampMS()

	action := map[string]interface{}{
		"type": "cancelByCloid",
		"cancels": []map[string]interface{}{
			{
				"asset": asset,
				"cloid": cloid.ToRaw(),
			},
		},
	}

	signature, err := utils.SignL1ActionWithAccount(
		e.privateKey,
		action,
		e.vaultAddress,
		timestamp,
		e.expiresAfter,
		e.IsMainnet(),
		e.accountAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign cancel by cloid action: %w", err)
	}

	return e.postAction(action, signature, timestamp)
}

// Modify modifies an existing order
func (e *Exchange) Modify(oid int, orderRequest types.OrderRequest) (map[string]interface{}, error) {
	asset, err := e.info.NameToAsset(orderRequest.Coin)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset for coin %s: %w", orderRequest.Coin, err)
	}

	orderWire, err := utils.OrderRequestToOrderWire(orderRequest, asset)
	if err != nil {
		return nil, fmt.Errorf("failed to convert order to wire format: %w", err)
	}

	timestamp := utils.GetTimestampMS()

	// Convert OrderWire to map for proper JSON serialization
	orderMap := map[string]interface{}{
		"a": orderWire.A,
		"b": orderWire.B,
		"p": orderWire.P,
		"s": orderWire.S,
		"r": orderWire.R,
		"t": utils.ConvertOrderTypeWireToMap(orderWire.T),
	}
	if orderWire.C != nil {
		orderMap["c"] = *orderWire.C
	}
	
	action := map[string]interface{}{
		"type": "modify",
		"modifies": []map[string]interface{}{
			{
				"oid":   oid,
				"order": orderMap,
			},
		},
	}

	signature, err := utils.SignL1ActionWithAccount(
		e.privateKey,
		action,
		e.vaultAddress,
		timestamp,
		e.expiresAfter,
		e.IsMainnet(),
		e.accountAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign modify action: %w", err)
	}

	return e.postAction(action, signature, timestamp)
}

// CancelAll cancels all open orders
func (e *Exchange) CancelAll() (map[string]interface{}, error) {
	timestamp := utils.GetTimestampMS()

	action := map[string]interface{}{
		"type": "cancelAll",
	}

	signature, err := utils.SignL1ActionWithAccount(
		e.privateKey,
		action,
		e.vaultAddress,
		timestamp,
		e.expiresAfter,
		e.IsMainnet(),
		e.accountAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign cancel all action: %w", err)
	}

	return e.postAction(action, signature, timestamp)
}

// UpdateLeverage updates the leverage for a coin
func (e *Exchange) UpdateLeverage(coin string, isCross bool, leverage int) (map[string]interface{}, error) {
	asset, err := e.info.NameToAsset(coin)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset for coin %s: %w", coin, err)
	}

	timestamp := utils.GetTimestampMS()

	action := map[string]interface{}{
		"type":     "updateLeverage",
		"asset":    asset,
		"isCross":  isCross,
		"leverage": leverage,
	}

	signature, err := utils.SignL1ActionWithAccount(
		e.privateKey,
		action,
		e.vaultAddress,
		timestamp,
		e.expiresAfter,
		e.IsMainnet(),
		e.accountAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign update leverage action: %w", err)
	}

	return e.postAction(action, signature, timestamp)
}

// UpdateIsolatedMargin updates the isolated margin for a coin
func (e *Exchange) UpdateIsolatedMargin(coin string, isBuy bool, ntli int64) (map[string]interface{}, error) {
	asset, err := e.info.NameToAsset(coin)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset for coin %s: %w", coin, err)
	}

	timestamp := utils.GetTimestampMS()

	action := map[string]interface{}{
		"type":  "updateIsolatedMargin",
		"asset": asset,
		"isBuy": isBuy,
		"ntli":  ntli,
	}

	signature, err := utils.SignL1ActionWithAccount(
		e.privateKey,
		action,
		e.vaultAddress,
		timestamp,
		e.expiresAfter,
		e.IsMainnet(),
		e.accountAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign update isolated margin action: %w", err)
	}

	return e.postAction(action, signature, timestamp)
}

// UsdTransfer transfers USD to another address
func (e *Exchange) UsdTransfer(destination string, amount string) (map[string]interface{}, error) {
	timestamp := utils.GetTimestampMS()

	// Create action for signing (without type field)
	signAction := map[string]interface{}{
		"destination": strings.ToLower(destination),
		"amount":      amount,
		"time":        fmt.Sprintf("%d", timestamp), // String for EIP712
	}

	signature, err := utils.SignUSDTransferAction(e.privateKey, signAction, e.IsMainnet())
	if err != nil {
		return nil, fmt.Errorf("failed to sign USD transfer action: %w", err)
	}

	// Send direct payload without wrapper (user-signed actions don't use postAction wrapper)
	payload := map[string]interface{}{
		"type":        "usdSend",
		"destination": strings.ToLower(destination),
		"amount":      amount,
		"time":        timestamp, // int64 for API
		"signature":   signature,
	}

	return e.Post("/exchange", payload)
}

// SpotTransfer transfers spot assets to another address
func (e *Exchange) SpotTransfer(destination string, token string, amount string) (map[string]interface{}, error) {
	timestamp := utils.GetTimestampMS()

	// Create action for signing (EIP712 expects time as string)
	signAction := map[string]interface{}{
		"destination": strings.ToLower(destination),
		"token":       token,
		"amount":      amount,
		"time":        fmt.Sprintf("%d", timestamp), // uint64 as string for EIP712
	}

	signature, err := utils.SignSpotTransferAction(e.privateKey, signAction, e.IsMainnet())
	if err != nil {
		return nil, fmt.Errorf("failed to sign spot transfer action: %w", err)
	}

	// Send direct payload (user-signed actions don't use postAction wrapper)
	payload := map[string]interface{}{
		"type":        "spotSend",
		"destination": strings.ToLower(destination),
		"token":       token,
		"amount":      amount,
		"time":        timestamp, // int64 for API
		"signature":   signature,
	}

	return e.Post("/exchange", payload)
}

// WithdrawFromBridge withdraws assets from the bridge
func (e *Exchange) WithdrawFromBridge(destination string, amount string) (map[string]interface{}, error) {
	timestamp := utils.GetTimestampMS()

	// Create action for signing (EIP712 expects time as string)
	signAction := map[string]interface{}{
		"destination": strings.ToLower(destination),
		"amount":      amount,
		"time":        fmt.Sprintf("%d", timestamp), // uint64 as string for EIP712
	}

	signature, err := utils.SignWithdrawFromBridgeAction(e.privateKey, signAction, e.IsMainnet())
	if err != nil {
		return nil, fmt.Errorf("failed to sign withdraw action: %w", err)
	}

	// Send direct payload (user-signed actions don't use postAction wrapper)
	payload := map[string]interface{}{
		"type":        "withdraw",
		"destination": strings.ToLower(destination),
		"amount":      amount,
		"time":        timestamp, // int64 for API
		"signature":   signature,
	}

	return e.Post("/exchange", payload)
}

// ApproveAgentResult represents the result of approving an agent
type ApproveAgentResult struct {
	Result   map[string]interface{} `json:"result"`
	AgentKey string                 `json:"agent_key"`
}

// ApproveAgent creates and approves an agent for trading on behalf of the account
// agentName is optional - if empty, a temporary agent is created
// Returns the API response and the agent's private key
func (e *Exchange) ApproveAgent(agentName ...string) (*ApproveAgentResult, error) {
	// Generate a random wallet for the agent
	agentPrivateKey, err := utils.CreateRandomWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to create agent wallet: %w", err)
	}

	// Get the agent's address
	agentAddress := utils.GetAddressFromPrivateKey(agentPrivateKey)

	// Determine the agent name
	var name string
	if len(agentName) > 0 && agentName[0] != "" {
		name = agentName[0]
	} else {
		// Use empty string for temporary agents
		name = ""
	}

	// Get nonce
	nonce := utils.GetTimestampMS()

	// Create action for signing (without type field)
	signAction := map[string]interface{}{
		"agentAddress": strings.ToLower(agentAddress),
		"agentName":    name,
		"nonce":       fmt.Sprintf("%d", nonce), // String for EIP712
	}

	// Sign the action
	signature, err := utils.SignAgent(e.privateKey, signAction, e.IsMainnet())
	if err != nil {
		return nil, fmt.Errorf("failed to sign agent approval: %w", err)
	}

	// Send direct payload without wrapper (user-signed actions don't use postAction wrapper)
	payload := map[string]interface{}{
		"type":         "approveAgent",
		"agentAddress": strings.ToLower(agentAddress),
		"nonce":       nonce, // int64 for API
		"signature":   signature,
	}

	// Only include agentName if not empty
	if name != "" {
		payload["agentName"] = name
	}

	result, err := e.Post("/exchange", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to approve agent: %w", err)
	}

	// Return both the result and the agent's private key
	return &ApproveAgentResult{
		Result:   result,
		AgentKey: fmt.Sprintf("%#x", agentPrivateKey.D),
	}, nil
}
