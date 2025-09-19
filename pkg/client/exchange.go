package client

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"hyperliquid-go-sdk/pkg/constants"
	"hyperliquid-go-sdk/pkg/signing"
	"hyperliquid-go-sdk/pkg/types"
)

// ExchangeClient provides trading functionality for Hyperliquid
type ExchangeClient struct {
	*BaseClient
	privateKey     *ecdsa.PrivateKey
	walletAddress  string
	vaultAddress   *string
	accountAddress *string
	info           *InfoClient
	expiresAfter   *int64
}

// NewExchangeClient creates a new exchange client
func NewExchangeClient(
	privateKey *ecdsa.PrivateKey,
	baseURL string,
	meta *types.Meta,
	vaultAddress *string,
	accountAddress *string,
	spotMeta *types.SpotMeta,
	perpDexs []string,
	timeout *time.Duration,
) (*ExchangeClient, error) {
	baseClient, err := NewBaseClient(baseURL, timeout)
	if err != nil {
		return nil, err
	}

	// Get wallet address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key type")
	}
	walletAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	// Initialize info client
	infoClient, err := NewInfoClient(baseURL, true, meta, spotMeta, perpDexs)
	if err != nil {
		return nil, err
	}

	return &ExchangeClient{
		BaseClient:     baseClient,
		privateKey:     privateKey,
		walletAddress:  walletAddress,
		vaultAddress:   vaultAddress,
		accountAddress: accountAddress,
		info:           infoClient,
	}, nil
}

// SetExpiresAfter sets the expiration time for actions
func (c *ExchangeClient) SetExpiresAfter(expiresAfter *int64) {
	c.expiresAfter = expiresAfter
}

// postAction posts an action to the exchange
func (c *ExchangeClient) postAction(ctx context.Context, action interface{}, signature types.Signature, nonce int64) (interface{}, error) {
	vaultAddr := c.vaultAddress
	// Special handling for certain action types
	if actionMap, ok := action.(map[string]interface{}); ok {
		actionType, _ := actionMap["type"].(string)
		if actionType == constants.ActionUsdClassTransfer || actionType == constants.ActionSendAsset {
			vaultAddr = nil
		}
	}

	payload := map[string]interface{}{
		"action":    action,
		"nonce":     nonce,
		"signature": signature,
	}

	if vaultAddr != nil {
		payload["vaultAddress"] = *vaultAddr
	} else {
		payload["vaultAddress"] = nil
	}

	if c.expiresAfter != nil {
		payload["expiresAfter"] = *c.expiresAfter
	}

	resp, err := c.apiClient.Post(ctx, "/exchange", payload)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// slippagePrice calculates the slippage price for market orders
func (c *ExchangeClient) slippagePrice(ctx context.Context, name string, isBuy bool, slippage float64, px *float64) (float64, error) {
	var price float64
	if px != nil {
		price = *px
	} else {
		// Get mid price
		allMids, err := c.info.AllMids(ctx, "")
		if err != nil {
			return 0, err
		}

		coin := c.info.nameToCoin[name]
		if midStr, exists := allMids[coin]; exists {
			var err error
			price, err = strconv.ParseFloat(string(midStr), 64)
			if err != nil {
				return 0, fmt.Errorf("failed to parse mid price: %w", err)
			}
		} else {
			return 0, fmt.Errorf("mid price not found for coin %s", coin)
		}
	}

	asset := c.info.NameToAsset(name)
	isSpot := asset >= 10000

	// Calculate slippage
	if isBuy {
		price *= (1 + slippage)
	} else {
		price *= (1 - slippage)
	}

	// Round price to appropriate decimals
	maxDecimals := constants.MaxDecimals
	if isSpot {
		maxDecimals = constants.SpotMaxDecimals
	}

	szDecimals := c.info.assetToSzDecimals[asset]
	decimalPlaces := maxDecimals - szDecimals

	// Format to 5 significant figures and appropriate decimal places
	formatted := fmt.Sprintf("%.5g", price)
	price, _ = strconv.ParseFloat(formatted, 64)

	// Round to decimal places
	multiplier := 1.0
	for i := 0; i < decimalPlaces; i++ {
		multiplier *= 10
	}
	price = float64(int(price*multiplier+0.5)) / multiplier

	return price, nil
}

// Order places a single order
func (c *ExchangeClient) Order(ctx context.Context, order types.OrderRequest, builder *types.BuilderInfo) (interface{}, error) {
	return c.BulkOrders(ctx, []types.OrderRequest{order}, builder)
}

// BulkOrders places multiple orders
func (c *ExchangeClient) BulkOrders(ctx context.Context, orderRequests []types.OrderRequest, builder *types.BuilderInfo) (interface{}, error) {
	var orderWires []types.OrderWire
	for _, order := range orderRequests {
		asset := c.info.NameToAsset(order.Coin)
		if asset == 0 {
			return nil, fmt.Errorf("unknown coin: %s", order.Coin)
		}

		orderWire, err := signing.OrderRequestToOrderWire(order, asset)
		if err != nil {
			return nil, err
		}
		orderWires = append(orderWires, orderWire)
	}

	timestamp := types.GetTimestampMs()

	if builder != nil {
		builder.B = strings.ToLower(builder.B)
	}

	orderAction := signing.OrderWiresToOrderAction(orderWires, builder)

	signature, err := signing.SignL1Action(
		c.privateKey,
		orderAction,
		c.vaultAddress,
		timestamp,
		c.expiresAfter,
		c.IsMainnet(),
	)
	if err != nil {
		return nil, err
	}

	return c.postAction(ctx, orderAction, signature, timestamp)
}

// MarketOpen places a market order to open a position
func (c *ExchangeClient) MarketOpen(ctx context.Context, name string, isBuy bool, sz float64, px *float64, slippage *float64, cloid *types.Cloid, builder *types.BuilderInfo) (interface{}, error) {
	if slippage == nil {
		defaultSlippage := constants.DefaultSlippage
		slippage = &defaultSlippage
	}

	price, err := c.slippagePrice(ctx, name, isBuy, *slippage, px)
	if err != nil {
		return nil, err
	}

	order := types.OrderRequest{
		Coin:       name,
		IsBuy:      isBuy,
		Sz:         sz,
		LimitPx:    price,
		OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: constants.TifIoc}},
		ReduceOnly: false,
		Cloid:      cloid,
	}

	return c.Order(ctx, order, builder)
}

// MarketClose places a market order to close a position
func (c *ExchangeClient) MarketClose(ctx context.Context, coin string, sz *float64, px *float64, slippage *float64, cloid *types.Cloid, builder *types.BuilderInfo) (interface{}, error) {
	address := c.walletAddress
	if c.accountAddress != nil {
		address = *c.accountAddress
	}
	if c.vaultAddress != nil {
		address = *c.vaultAddress
	}

	userState, err := c.info.UserState(ctx, address, "")
	if err != nil {
		return nil, err
	}

	for _, assetPosition := range userState.AssetPositions {
		if coin != assetPosition.Position.Coin {
			continue
		}

		szi, err := strconv.ParseFloat(string(assetPosition.Position.Szi), 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse position size: %w", err)
		}

		size := sz
		if size == nil {
			absSize := szi
			if absSize < 0 {
				absSize = -absSize
			}
			size = &absSize
		}

		isBuy := szi < 0

		if slippage == nil {
			defaultSlippage := constants.DefaultSlippage
			slippage = &defaultSlippage
		}

		price, err := c.slippagePrice(ctx, coin, isBuy, *slippage, px)
		if err != nil {
			return nil, err
		}

		order := types.OrderRequest{
			Coin:       coin,
			IsBuy:      isBuy,
			Sz:         *size,
			LimitPx:    price,
			OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: constants.TifIoc}},
			ReduceOnly: true,
			Cloid:      cloid,
		}

		return c.Order(ctx, order, builder)
	}

	return nil, fmt.Errorf("no position found for coin %s", coin)
}

// Cancel cancels an order by order ID
func (c *ExchangeClient) Cancel(ctx context.Context, name string, oid int) (interface{}, error) {
	return c.BulkCancel(ctx, []types.CancelRequest{{Coin: name, Oid: oid}})
}

// CancelByCloid cancels an order by client order ID
func (c *ExchangeClient) CancelByCloid(ctx context.Context, name string, cloid *types.Cloid) (interface{}, error) {
	return c.BulkCancelByCloid(ctx, []types.CancelByCloidRequest{{Coin: name, Cloid: cloid}})
}

// BulkCancel cancels multiple orders by order ID
func (c *ExchangeClient) BulkCancel(ctx context.Context, cancelRequests []types.CancelRequest) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	var cancels []map[string]interface{}
	for _, cancel := range cancelRequests {
		asset := c.info.NameToAsset(cancel.Coin)
		if asset == 0 {
			return nil, fmt.Errorf("unknown coin: %s", cancel.Coin)
		}

		cancels = append(cancels, map[string]interface{}{
			"a": asset,
			"o": cancel.Oid,
		})
	}

	//cancelAction := types.CancelAction{
	//	Type: constants.ActionCancel,
	//}

	// Set cancels field manually since we need the specific structure
	cancelActionMap := map[string]interface{}{
		"type":    constants.ActionCancel,
		"cancels": cancels,
	}

	signature, err := signing.SignL1Action(
		c.privateKey,
		cancelActionMap,
		c.vaultAddress,
		timestamp,
		c.expiresAfter,
		c.IsMainnet(),
	)
	if err != nil {
		return nil, err
	}

	return c.postAction(ctx, cancelActionMap, signature, timestamp)
}

// BulkCancelByCloid cancels multiple orders by client order ID
func (c *ExchangeClient) BulkCancelByCloid(ctx context.Context, cancelRequests []types.CancelByCloidRequest) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	var cancels []map[string]interface{}
	for _, cancel := range cancelRequests {
		asset := c.info.NameToAsset(cancel.Coin)
		if asset == 0 {
			return nil, fmt.Errorf("unknown coin: %s", cancel.Coin)
		}

		cancels = append(cancels, map[string]interface{}{
			"asset": asset,
			"cloid": cancel.Cloid.Raw(),
		})
	}

	cancelActionMap := map[string]interface{}{
		"type":    constants.ActionCancelByCloid,
		"cancels": cancels,
	}

	signature, err := signing.SignL1Action(
		c.privateKey,
		cancelActionMap,
		c.vaultAddress,
		timestamp,
		c.expiresAfter,
		c.IsMainnet(),
	)
	if err != nil {
		return nil, err
	}

	return c.postAction(ctx, cancelActionMap, signature, timestamp)
}

// ScheduleCancel schedules order cancellation
func (c *ExchangeClient) ScheduleCancel(ctx context.Context, time *int64) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	scheduleCancelAction := map[string]interface{}{
		"type": constants.ActionScheduleCancel,
	}
	if time != nil {
		scheduleCancelAction["time"] = *time
	}

	signature, err := signing.SignL1Action(
		c.privateKey,
		scheduleCancelAction,
		c.vaultAddress,
		timestamp,
		c.expiresAfter,
		c.IsMainnet(),
	)
	if err != nil {
		return nil, err
	}

	return c.postAction(ctx, scheduleCancelAction, signature, timestamp)
}

// UpdateLeverage updates leverage for an asset
func (c *ExchangeClient) UpdateLeverage(ctx context.Context, leverage int, name string, isCross bool) (interface{}, error) {
	timestamp := types.GetTimestampMs()
	asset := c.info.NameToAsset(name)
	if asset == 0 {
		return nil, fmt.Errorf("unknown coin: %s", name)
	}

	updateLeverageAction := map[string]interface{}{
		"type":     constants.ActionUpdateLeverage,
		"asset":    asset,
		"isCross":  isCross,
		"leverage": leverage,
	}

	signature, err := signing.SignL1Action(
		c.privateKey,
		updateLeverageAction,
		c.vaultAddress,
		timestamp,
		c.expiresAfter,
		c.IsMainnet(),
	)
	if err != nil {
		return nil, err
	}

	return c.postAction(ctx, updateLeverageAction, signature, timestamp)
}

// UpdateIsolatedMargin updates isolated margin for an asset
func (c *ExchangeClient) UpdateIsolatedMargin(ctx context.Context, amount float64, name string) (interface{}, error) {
	timestamp := types.GetTimestampMs()
	asset := c.info.NameToAsset(name)
	if asset == 0 {
		return nil, fmt.Errorf("unknown coin: %s", name)
	}

	amountInt, err := signing.FloatToUsdInt(amount)
	if err != nil {
		return nil, err
	}

	updateIsolatedMarginAction := map[string]interface{}{
		"type":  constants.ActionUpdateIsolatedMargin,
		"asset": asset,
		"isBuy": true,
		"ntli":  amountInt,
	}

	signature, err := signing.SignL1Action(
		c.privateKey,
		updateIsolatedMarginAction,
		c.vaultAddress,
		timestamp,
		c.expiresAfter,
		c.IsMainnet(),
	)
	if err != nil {
		return nil, err
	}

	return c.postAction(ctx, updateIsolatedMarginAction, signature, timestamp)
}

// ApproveAgent approves an agent
func (c *ExchangeClient) ApproveAgent(ctx context.Context, name *string) (interface{}, string, error) {
	// Generate a random agent key
	agentKeyBytes := make([]byte, 32)
	if _, err := rand.Read(agentKeyBytes); err != nil {
		return nil, "", fmt.Errorf("failed to generate agent key: %w", err)
	}
	agentKey := "0x" + hex.EncodeToString(agentKeyBytes)

	// Create agent account from the generated key
	agentPrivateKey, err := crypto.HexToECDSA(agentKey[2:])
	if err != nil {
		return nil, "", fmt.Errorf("failed to create agent private key: %w", err)
	}

	agentPublicKey := agentPrivateKey.Public()
	agentPublicKeyECDSA, ok := agentPublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, "", fmt.Errorf("invalid agent public key type")
	}
	agentAddress := crypto.PubkeyToAddress(*agentPublicKeyECDSA).Hex()

	timestamp := types.GetTimestampMs()

	approveAgentAction := signing.ApproveAgentAction{
		Type:         constants.ActionApproveAgent,
		AgentAddress: agentAddress,
		Nonce:        timestamp,
	}

	if name != nil {
		approveAgentAction.AgentName = *name
	} else {
		approveAgentAction.AgentName = ""
	}

	signature, err := signing.SignApproveAgentAction(c.privateKey, approveAgentAction, c.IsMainnet())
	if err != nil {
		return nil, "", err
	}

	// Convert to map for API call
	actionMap := map[string]interface{}{
		"type":         approveAgentAction.Type,
		"agentAddress": approveAgentAction.AgentAddress,
		"nonce":        approveAgentAction.Nonce,
	}
	if name != nil {
		actionMap["agentName"] = approveAgentAction.AgentName
	}

	result, err := c.postAction(ctx, actionMap, signature, timestamp)
	return result, agentKey, err
}

// USDTransfer transfers USD to another address
func (c *ExchangeClient) USDTransfer(ctx context.Context, amount float64, destination string) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	action := signing.USDTransferAction{
		Type:        constants.ActionUsdSend,
		Destination: destination,
		Amount:      fmt.Sprintf("%.6f", amount),
		Time:        timestamp,
	}

	signature, err := signing.SignUSDTransferAction(c.privateKey, action, c.IsMainnet())
	if err != nil {
		return nil, err
	}

	// Convert to map for API call
	actionMap := map[string]interface{}{
		"type":        action.Type,
		"destination": action.Destination,
		"amount":      action.Amount,
		"time":        action.Time,
	}

	return c.postAction(ctx, actionMap, signature, timestamp)
}

// SpotTransfer transfers spot tokens to another address
func (c *ExchangeClient) SpotTransfer(ctx context.Context, amount float64, destination, token string) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	action := signing.SpotTransferAction{
		Type:        constants.ActionSpotSend,
		Destination: destination,
		Amount:      fmt.Sprintf("%.8f", amount),
		Token:       token,
		Time:        timestamp,
	}

	signature, err := signing.SignSpotTransferAction(c.privateKey, action, c.IsMainnet())
	if err != nil {
		return nil, err
	}

	// Convert to map for API call
	actionMap := map[string]interface{}{
		"type":        action.Type,
		"destination": action.Destination,
		"amount":      action.Amount,
		"token":       action.Token,
		"time":        action.Time,
	}

	return c.postAction(ctx, actionMap, signature, timestamp)
}

// WithdrawFromBridge withdraws funds from the bridge
func (c *ExchangeClient) WithdrawFromBridge(ctx context.Context, amount float64, destination string) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	action := signing.WithdrawAction{
		Type:        "withdraw3",
		Destination: destination,
		Amount:      fmt.Sprintf("%.6f", amount),
		Time:        timestamp,
	}

	signature, err := signing.SignWithdrawAction(c.privateKey, action, c.IsMainnet())
	if err != nil {
		return nil, err
	}

	// Convert to map for API call
	actionMap := map[string]interface{}{
		"type":        action.Type,
		"destination": action.Destination,
		"amount":      action.Amount,
		"time":        action.Time,
	}

	return c.postAction(ctx, actionMap, signature, timestamp)
}

// USDClassTransfer transfers USD between perp and spot
func (c *ExchangeClient) USDClassTransfer(ctx context.Context, amount float64, toPerp bool) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	amountStr := fmt.Sprintf("%.6f", amount)
	if c.vaultAddress != nil {
		amountStr += fmt.Sprintf(" subaccount:%s", *c.vaultAddress)
	}

	action := signing.USDClassTransferAction{
		Type:   constants.ActionUsdClassTransfer,
		Amount: amountStr,
		ToPerp: toPerp,
		Nonce:  timestamp,
	}

	signature, err := signing.SignUSDClassTransferAction(c.privateKey, action, c.IsMainnet())
	if err != nil {
		return nil, err
	}

	// Convert to map for API call
	actionMap := map[string]interface{}{
		"type":   action.Type,
		"amount": action.Amount,
		"toPerp": action.ToPerp,
		"nonce":  action.Nonce,
	}

	return c.postAction(ctx, actionMap, signature, timestamp)
}

// SendAsset sends assets between different dexs
func (c *ExchangeClient) SendAsset(ctx context.Context, destination, sourceDex, destinationDex, token string, amount float64) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	fromSubAccount := ""
	if c.vaultAddress != nil {
		fromSubAccount = *c.vaultAddress
	}

	action := signing.SendAssetAction{
		Type:           constants.ActionSendAsset,
		Destination:    destination,
		SourceDex:      sourceDex,
		DestinationDex: destinationDex,
		Token:          token,
		Amount:         fmt.Sprintf("%.8f", amount),
		FromSubAccount: fromSubAccount,
		Nonce:          timestamp,
	}

	signature, err := signing.SignSendAssetAction(c.privateKey, action, c.IsMainnet())
	if err != nil {
		return nil, err
	}

	// Convert to map for API call
	actionMap := map[string]interface{}{
		"type":           action.Type,
		"destination":    action.Destination,
		"sourceDex":      action.SourceDex,
		"destinationDex": action.DestinationDex,
		"token":          action.Token,
		"amount":         action.Amount,
		"fromSubAccount": action.FromSubAccount,
		"nonce":          action.Nonce,
	}

	return c.postAction(ctx, actionMap, signature, timestamp)
}

// SetReferrer sets a referrer code
func (c *ExchangeClient) SetReferrer(ctx context.Context, code string) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	setReferrerAction := map[string]interface{}{
		"type": constants.ActionSetReferrer,
		"code": code,
	}

	signature, err := signing.SignL1Action(
		c.privateKey,
		setReferrerAction,
		nil, // Set referrer doesn't use vault address
		timestamp,
		c.expiresAfter,
		c.IsMainnet(),
	)
	if err != nil {
		return nil, err
	}

	return c.postAction(ctx, setReferrerAction, signature, timestamp)
}

// CreateSubAccount creates a new sub-account
func (c *ExchangeClient) CreateSubAccount(ctx context.Context, name string) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	createSubAccountAction := map[string]interface{}{
		"type": constants.ActionCreateSubAccount,
		"name": name,
	}

	signature, err := signing.SignL1Action(
		c.privateKey,
		createSubAccountAction,
		nil, // Create sub-account doesn't use vault address
		timestamp,
		c.expiresAfter,
		c.IsMainnet(),
	)
	if err != nil {
		return nil, err
	}

	return c.postAction(ctx, createSubAccountAction, signature, timestamp)
}

// SubAccountTransfer transfers funds to/from a sub-account
func (c *ExchangeClient) SubAccountTransfer(ctx context.Context, subAccountUser string, isDeposit bool, usd int64) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	subAccountTransferAction := map[string]interface{}{
		"type":           constants.ActionSubAccountTransfer,
		"subAccountUser": subAccountUser,
		"isDeposit":      isDeposit,
		"usd":            usd,
	}

	signature, err := signing.SignL1Action(
		c.privateKey,
		subAccountTransferAction,
		nil, // Sub-account transfer doesn't use vault address
		timestamp,
		c.expiresAfter,
		c.IsMainnet(),
	)
	if err != nil {
		return nil, err
	}

	return c.postAction(ctx, subAccountTransferAction, signature, timestamp)
}

// VaultTransfer transfers funds to/from a vault
func (c *ExchangeClient) VaultTransfer(ctx context.Context, vaultAddress string, isDeposit bool, usd int64) (interface{}, error) {
	timestamp := types.GetTimestampMs()

	vaultTransferAction := map[string]interface{}{
		"type":         constants.ActionVaultTransfer,
		"vaultAddress": vaultAddress,
		"isDeposit":    isDeposit,
		"usd":          usd,
	}

	signature, err := signing.SignL1Action(
		c.privateKey,
		vaultTransferAction,
		nil, // Vault transfer doesn't use vault address
		timestamp,
		c.expiresAfter,
		c.IsMainnet(),
	)
	if err != nil {
		return nil, err
	}

	return c.postAction(ctx, vaultTransferAction, signature, timestamp)
}

// GetWalletAddress returns the wallet address
func (c *ExchangeClient) GetWalletAddress() string {
	return c.walletAddress
}

// GetAccountAddress returns the account address (could be different from wallet if using agent)
func (c *ExchangeClient) GetAccountAddress() string {
	if c.accountAddress != nil {
		return *c.accountAddress
	}
	return c.walletAddress
}

// GetVaultAddress returns the vault address if set
func (c *ExchangeClient) GetVaultAddress() *string {
	return c.vaultAddress
}

// GetInfo returns the info client
func (c *ExchangeClient) GetInfo() *InfoClient {
	return c.info
}
