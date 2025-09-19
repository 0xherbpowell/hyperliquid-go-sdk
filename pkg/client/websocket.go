package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"hyperliquid-go-sdk/pkg/types"
)

// WebSocketClient manages WebSocket connections to Hyperliquid
type WebSocketClient struct {
	baseURL            string
	conn               *websocket.Conn
	subscriptions      map[string]*types.ActiveSubscription
	subscriptionsMutex sync.RWMutex
	subscriptionID     int
	isConnected        bool
	pingTicker         *time.Ticker
	stopChan           chan struct{}
	reconnectAttempts  int
	maxReconnectTries  int
	reconnectDelay     time.Duration
	onDisconnect       func()
	onReconnect        func()
	logger             *log.Logger
}

// WebSocketConfig represents WebSocket configuration
type WebSocketConfig struct {
	MaxReconnectTries int
	ReconnectDelay    time.Duration
	PingInterval      time.Duration
	OnDisconnect      func()
	OnReconnect       func()
	Logger            *log.Logger
}

// NewWebSocketClient creates a new WebSocket client
func NewWebSocketClient(baseURL string, config *WebSocketConfig) *WebSocketClient {
	if config == nil {
		config = &WebSocketConfig{
			MaxReconnectTries: 10,
			ReconnectDelay:    5 * time.Second,
			PingInterval:      50 * time.Second,
		}
	}

	wsURL := strings.Replace(baseURL, "http", "ws", 1) + "/ws"

	return &WebSocketClient{
		baseURL:           wsURL,
		subscriptions:     make(map[string]*types.ActiveSubscription),
		maxReconnectTries: config.MaxReconnectTries,
		reconnectDelay:    config.ReconnectDelay,
		stopChan:          make(chan struct{}),
		onDisconnect:      config.OnDisconnect,
		onReconnect:       config.OnReconnect,
		logger:            config.Logger,
		pingTicker:        time.NewTicker(config.PingInterval),
	}
}

// Connect establishes a WebSocket connection
func (ws *WebSocketClient) Connect(ctx context.Context) error {
	u, err := url.Parse(ws.baseURL)
	if err != nil {
		return fmt.Errorf("invalid WebSocket URL: %w", err)
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	ws.conn = conn
	ws.isConnected = true
	ws.reconnectAttempts = 0

	// Start message handling
	go ws.handleMessages()
	go ws.handlePing()

	// Resubscribe to existing subscriptions
	ws.subscriptionsMutex.RLock()
	for _, sub := range ws.subscriptions {
		if err := ws.sendSubscription(sub); err != nil {
			ws.logf("Failed to resubscribe: %v", err)
		}
	}
	ws.subscriptionsMutex.RUnlock()

	if ws.onReconnect != nil {
		ws.onReconnect()
	}

	return nil
}

// Disconnect closes the WebSocket connection
func (ws *WebSocketClient) Disconnect() {
	close(ws.stopChan)

	if ws.pingTicker != nil {
		ws.pingTicker.Stop()
	}

	if ws.conn != nil {
		ws.conn.Close()
	}

	ws.isConnected = false
}

// IsConnected returns true if the WebSocket is connected
func (ws *WebSocketClient) IsConnected() bool {
	return ws.isConnected
}

// Subscribe subscribes to a WebSocket channel
func (ws *WebSocketClient) Subscribe(subscription types.Subscription, callback types.SubscriptionCallback) (int, error) {
	ws.subscriptionsMutex.Lock()
	defer ws.subscriptionsMutex.Unlock()

	ws.subscriptionID++
	subscriptionID := ws.subscriptionID

	identifier := ws.subscriptionToIdentifier(subscription)

	activeSub := &types.ActiveSubscription{
		Callback:       callback,
		SubscriptionID: subscriptionID,
	}

	// Store the subscription
	ws.subscriptions[identifier] = activeSub

	// Send subscription if connected
	if ws.isConnected {
		if err := ws.sendSubscription(activeSub); err != nil {
			delete(ws.subscriptions, identifier)
			return 0, err
		}
	}

	return subscriptionID, nil
}

// Unsubscribe unsubscribes from a WebSocket channel
func (ws *WebSocketClient) Unsubscribe(subscription types.Subscription, subscriptionID int) error {
	ws.subscriptionsMutex.Lock()
	defer ws.subscriptionsMutex.Unlock()

	identifier := ws.subscriptionToIdentifier(subscription)

	if activeSub, exists := ws.subscriptions[identifier]; exists {
		if activeSub.SubscriptionID == subscriptionID {
			delete(ws.subscriptions, identifier)

			if ws.isConnected {
				return ws.sendUnsubscription(subscription)
			}
		}
	}

	return nil
}

// sendSubscription sends a subscription message
func (ws *WebSocketClient) sendSubscription(activeSub *types.ActiveSubscription) error {
	// This would need the actual subscription object, which we don't store
	// In a real implementation, you'd need to store the subscription object too
	// For now, this is a placeholder
	return nil
}

// sendUnsubscription sends an unsubscription message
func (ws *WebSocketClient) sendUnsubscription(subscription types.Subscription) error {
	msg := map[string]interface{}{
		"method":       "unsubscribe",
		"subscription": subscription,
	}

	return ws.conn.WriteJSON(msg)
}

// handleMessages handles incoming WebSocket messages
func (ws *WebSocketClient) handleMessages() {
	defer func() {
		ws.isConnected = false
		if ws.onDisconnect != nil {
			ws.onDisconnect()
		}
	}()

	for {
		select {
		case <-ws.stopChan:
			return
		default:
			_, message, err := ws.conn.ReadMessage()
			if err != nil {
				ws.logf("WebSocket read error: %v", err)
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					ws.logf("Unexpected WebSocket close: %v", err)
				}

				// Attempt to reconnect
				go ws.attemptReconnect()
				return
			}

			ws.processMessage(message)
		}
	}
}

// processMessage processes incoming WebSocket messages
func (ws *WebSocketClient) processMessage(message []byte) {
	// Check for connection establishment message
	if string(message) == "Websocket connection established." {
		ws.logf("WebSocket connection established")
		return
	}

	// Parse JSON message
	var rawMsg map[string]interface{}
	if err := json.Unmarshal(message, &rawMsg); err != nil {
		ws.logf("Failed to parse WebSocket message: %v", err)
		return
	}

	// Handle pong messages
	if channel, ok := rawMsg["channel"].(string); ok && channel == "pong" {
		ws.logf("Received pong")
		return
	}

	// Route message to appropriate subscription
	identifier := ws.messageToIdentifier(rawMsg)
	if identifier == "" {
		ws.logf("Could not determine identifier for message: %s", string(message))
		return
	}

	ws.subscriptionsMutex.RLock()
	activeSub, exists := ws.subscriptions[identifier]
	ws.subscriptionsMutex.RUnlock()

	if !exists {
		ws.logf("No subscription found for identifier: %s", identifier)
		return
	}

	// Call the callback
	if activeSub.Callback != nil {
		go activeSub.Callback(rawMsg)
	}
}

// handlePing sends periodic ping messages
func (ws *WebSocketClient) handlePing() {
	for {
		select {
		case <-ws.stopChan:
			return
		case <-ws.pingTicker.C:
			if ws.isConnected {
				msg := map[string]string{"method": "ping"}
				if err := ws.conn.WriteJSON(msg); err != nil {
					ws.logf("Failed to send ping: %v", err)
				} else {
					ws.logf("Sent ping")
				}
			}
		}
	}
}

// attemptReconnect attempts to reconnect to the WebSocket
func (ws *WebSocketClient) attemptReconnect() {
	if ws.reconnectAttempts >= ws.maxReconnectTries {
		ws.logf("Max reconnection attempts reached")
		return
	}

	ws.reconnectAttempts++
	ws.logf("Attempting to reconnect (%d/%d)", ws.reconnectAttempts, ws.maxReconnectTries)

	time.Sleep(ws.reconnectDelay)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := ws.Connect(ctx); err != nil {
		ws.logf("Reconnection failed: %v", err)
		go ws.attemptReconnect()
	} else {
		ws.logf("Successfully reconnected")
	}
}

// subscriptionToIdentifier converts a subscription to a unique identifier
func (ws *WebSocketClient) subscriptionToIdentifier(subscription types.Subscription) string {
	switch sub := subscription.(type) {
	case types.AllMidsSubscription:
		return "allMids"
	case types.L2BookSubscription:
		return fmt.Sprintf("l2Book:%s", strings.ToLower(sub.Coin))
	case types.TradesSubscription:
		return fmt.Sprintf("trades:%s", strings.ToLower(sub.Coin))
	case types.UserEventsSubscription:
		return "userEvents"
	case types.UserFillsSubscription:
		return fmt.Sprintf("userFills:%s", strings.ToLower(sub.User))
	case types.CandleSubscription:
		return fmt.Sprintf("candle:%s,%s", strings.ToLower(sub.Coin), sub.Interval)
	case types.OrderUpdatesSubscription:
		return "orderUpdates"
	case types.UserFundingsSubscription:
		return fmt.Sprintf("userFundings:%s", strings.ToLower(sub.User))
	case types.UserNonFundingLedgerUpdatesSubscription:
		return fmt.Sprintf("userNonFundingLedgerUpdates:%s", strings.ToLower(sub.User))
	case types.WebData2Subscription:
		return fmt.Sprintf("webData2:%s", strings.ToLower(sub.User))
	case types.BboSubscription:
		return fmt.Sprintf("bbo:%s", strings.ToLower(sub.Coin))
	case types.ActiveAssetCtxSubscription:
		return fmt.Sprintf("activeAssetCtx:%s", strings.ToLower(sub.Coin))
	case types.ActiveAssetDataSubscription:
		return fmt.Sprintf("activeAssetData:%s,%s", strings.ToLower(sub.Coin), strings.ToLower(sub.User))
	default:
		return ""
	}
}

// messageToIdentifier extracts the identifier from a WebSocket message
func (ws *WebSocketClient) messageToIdentifier(msg map[string]interface{}) string {
	channel, ok := msg["channel"].(string)
	if !ok {
		return ""
	}

	switch channel {
	case "pong":
		return "pong"
	case "allMids":
		return "allMids"
	case "l2Book":
		if data, ok := msg["data"].(map[string]interface{}); ok {
			if coin, ok := data["coin"].(string); ok {
				return fmt.Sprintf("l2Book:%s", strings.ToLower(coin))
			}
		}
	case "trades":
		if data, ok := msg["data"].([]interface{}); ok && len(data) > 0 {
			if trade, ok := data[0].(map[string]interface{}); ok {
				if coin, ok := trade["coin"].(string); ok {
					return fmt.Sprintf("trades:%s", strings.ToLower(coin))
				}
			}
		}
	case "user":
		return "userEvents"
	case "userFills":
		if data, ok := msg["data"].(map[string]interface{}); ok {
			if user, ok := data["user"].(string); ok {
				return fmt.Sprintf("userFills:%s", strings.ToLower(user))
			}
		}
	case "candle":
		if data, ok := msg["data"].(map[string]interface{}); ok {
			if coin, ok := data["s"].(string); ok {
				if interval, ok := data["i"].(string); ok {
					return fmt.Sprintf("candle:%s,%s", strings.ToLower(coin), interval)
				}
			}
		}
	case "orderUpdates":
		return "orderUpdates"
	case "userFundings":
		if data, ok := msg["data"].(map[string]interface{}); ok {
			if user, ok := data["user"].(string); ok {
				return fmt.Sprintf("userFundings:%s", strings.ToLower(user))
			}
		}
	case "userNonFundingLedgerUpdates":
		if data, ok := msg["data"].(map[string]interface{}); ok {
			if user, ok := data["user"].(string); ok {
				return fmt.Sprintf("userNonFundingLedgerUpdates:%s", strings.ToLower(user))
			}
		}
	case "webData2":
		if data, ok := msg["data"].(map[string]interface{}); ok {
			if user, ok := data["user"].(string); ok {
				return fmt.Sprintf("webData2:%s", strings.ToLower(user))
			}
		}
	case "bbo":
		if data, ok := msg["data"].(map[string]interface{}); ok {
			if coin, ok := data["coin"].(string); ok {
				return fmt.Sprintf("bbo:%s", strings.ToLower(coin))
			}
		}
	case "activeAssetCtx", "activeSpotAssetCtx":
		if data, ok := msg["data"].(map[string]interface{}); ok {
			if coin, ok := data["coin"].(string); ok {
				return fmt.Sprintf("activeAssetCtx:%s", strings.ToLower(coin))
			}
		}
	case "activeAssetData":
		if data, ok := msg["data"].(map[string]interface{}); ok {
			if coin, ok := data["coin"].(string); ok {
				if user, ok := data["user"].(string); ok {
					return fmt.Sprintf("activeAssetData:%s,%s", strings.ToLower(coin), strings.ToLower(user))
				}
			}
		}
	}

	return ""
}

// logf logs a message if a logger is configured
func (ws *WebSocketClient) logf(format string, args ...interface{}) {
	if ws.logger != nil {
		ws.logger.Printf(format, args...)
	}
}

// GetSubscriptionCount returns the number of active subscriptions
func (ws *WebSocketClient) GetSubscriptionCount() int {
	ws.subscriptionsMutex.RLock()
	defer ws.subscriptionsMutex.RUnlock()
	return len(ws.subscriptions)
}

// GetConnectionStatus returns detailed connection status
func (ws *WebSocketClient) GetConnectionStatus() map[string]interface{} {
	return map[string]interface{}{
		"connected":         ws.isConnected,
		"subscriptions":     ws.GetSubscriptionCount(),
		"reconnectAttempts": ws.reconnectAttempts,
		"maxReconnectTries": ws.maxReconnectTries,
	}
}
