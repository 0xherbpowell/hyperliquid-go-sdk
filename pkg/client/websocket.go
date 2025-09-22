package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"hyperliquid-go-sdk/pkg/types"
	"hyperliquid-go-sdk/pkg/utils"
)

// WebsocketManager manages WebSocket connections for real-time data
type WebsocketManager struct {
	baseURL         string
	wsURL           string
	conn            *websocket.Conn
	subscriptions   map[string]func(interface{})
	isRunning       bool
	mutex           sync.RWMutex
	reconnectDelay  time.Duration
	maxReconnects   int
	currentRetries  int
	pingInterval    time.Duration
	pongTimeout     time.Duration
	done            chan struct{}
}

// NewWebsocketManager creates a new WebSocket manager
func NewWebsocketManager(baseURL string) (*WebsocketManager, error) {
	var wsURL string
	
	switch baseURL {
	case utils.MainnetAPIURL:
		wsURL = utils.MainnetWSURL
	case utils.TestnetAPIURL:
		wsURL = utils.TestnetWSURL
	default:
		// For custom URLs, try to convert HTTP to WS
		u, err := url.Parse(baseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid base URL: %w", err)
		}
		
		switch u.Scheme {
		case "http":
			u.Scheme = "ws"
		case "https":
			u.Scheme = "wss"
		default:
			return nil, fmt.Errorf("unsupported URL scheme: %s", u.Scheme)
		}
		
		u.Path = "/ws"
		wsURL = u.String()
	}
	
	return &WebsocketManager{
		baseURL:        baseURL,
		wsURL:          wsURL,
		subscriptions:  make(map[string]func(interface{})),
		reconnectDelay: 5 * time.Second,
		maxReconnects:  10,
		pingInterval:   30 * time.Second,
		pongTimeout:    10 * time.Second,
		done:           make(chan struct{}),
	}, nil
}

// Start starts the WebSocket connection
func (w *WebsocketManager) Start() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	if w.isRunning {
		return fmt.Errorf("WebSocket manager is already running")
	}
	
	if err := w.connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	
	w.isRunning = true
	
	// Start message handling goroutines
	go w.readPump()
	go w.pingPump()
	
	return nil
}

// Stop stops the WebSocket connection
func (w *WebsocketManager) Stop() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	if !w.isRunning {
		return nil
	}
	
	w.isRunning = false
	close(w.done)
	
	if w.conn != nil {
		// Send close frame
		w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		w.conn.Close()
		w.conn = nil
	}
	
	return nil
}

// connect establishes the WebSocket connection
func (w *WebsocketManager) connect() error {
	dialer := websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
	}
	
	conn, _, err := dialer.Dial(w.wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to dial WebSocket: %w", err)
	}
	
	w.conn = conn
	w.currentRetries = 0
	
	// Set read deadline for pong messages
	w.conn.SetReadDeadline(time.Now().Add(w.pongTimeout))
	w.conn.SetPongHandler(func(string) error {
		w.conn.SetReadDeadline(time.Now().Add(w.pongTimeout))
		return nil
	})
	
	return nil
}

// reconnect attempts to reconnect the WebSocket
func (w *WebsocketManager) reconnect() error {
	if w.currentRetries >= w.maxReconnects {
		return fmt.Errorf("maximum reconnection attempts reached")
	}
	
	w.currentRetries++
	log.Printf("WebSocket reconnection attempt %d/%d", w.currentRetries, w.maxReconnects)
	
	time.Sleep(w.reconnectDelay)
	
	if err := w.connect(); err != nil {
		return fmt.Errorf("reconnection failed: %w", err)
	}
	
	// Resubscribe to all active subscriptions
	w.mutex.RLock()
	subscriptions := make([]string, 0, len(w.subscriptions))
	for sub := range w.subscriptions {
		subscriptions = append(subscriptions, sub)
	}
	w.mutex.RUnlock()
	
	for _, sub := range subscriptions {
		var subscription types.Subscription
		if err := json.Unmarshal([]byte(sub), &subscription); err == nil {
			w.sendSubscription(subscription)
		}
	}
	
	log.Printf("WebSocket reconnected successfully")
	return nil
}

// readPump handles incoming WebSocket messages
func (w *WebsocketManager) readPump() {
	defer func() {
		if w.conn != nil {
			w.conn.Close()
		}
	}()
	
	for {
		select {
		case <-w.done:
			return
		default:
			_, message, err := w.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				
				// Try to reconnect if still running
				w.mutex.RLock()
				isRunning := w.isRunning
				w.mutex.RUnlock()
				
				if isRunning {
					if err := w.reconnect(); err != nil {
						log.Printf("Failed to reconnect WebSocket: %v", err)
						return
					}
				} else {
					return
				}
				continue
			}
			
			w.handleMessage(message)
		}
	}
}

// pingPump sends ping messages to keep the connection alive
func (w *WebsocketManager) pingPump() {
	ticker := time.NewTicker(w.pingInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			w.mutex.RLock()
			conn := w.conn
			w.mutex.RUnlock()
			
			if conn != nil {
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("WebSocket ping failed: %v", err)
					return
				}
			}
		case <-w.done:
			return
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (w *WebsocketManager) handleMessage(message []byte) {
	var msgData map[string]interface{}
	if err := json.Unmarshal(message, &msgData); err != nil {
		log.Printf("Failed to unmarshal WebSocket message: %v", err)
		return
	}
	
	channel, ok := msgData["channel"].(string)
	if !ok {
		log.Printf("WebSocket message missing channel field")
		return
	}
	
	// Call all matching callbacks
	w.mutex.RLock()
	for subKey, callback := range w.subscriptions {
		var subscription types.Subscription
		if err := json.Unmarshal([]byte(subKey), &subscription); err == nil {
			if w.matchesSubscription(subscription, channel, msgData) {
				go callback(msgData)
			}
		}
	}
	w.mutex.RUnlock()
}

// matchesSubscription checks if a message matches a subscription
func (w *WebsocketManager) matchesSubscription(sub types.Subscription, channel string, msgData map[string]interface{}) bool {
	switch sub.Type {
	case "allMids":
		return channel == "allMids"
	case "bbo":
		if channel == "bbo" {
			if data, ok := msgData["data"].(map[string]interface{}); ok {
				if coin, ok := data["coin"].(string); ok {
					return coin == sub.Coin
				}
			}
		}
	case "l2Book":
		if channel == "l2Book" {
			if data, ok := msgData["data"].(map[string]interface{}); ok {
				if coin, ok := data["coin"].(string); ok {
					return coin == sub.Coin
				}
			}
		}
	case "trades":
		if channel == "trades" {
			if data, ok := msgData["data"].([]interface{}); ok && len(data) > 0 {
				if trade, ok := data[0].(map[string]interface{}); ok {
					if coin, ok := trade["coin"].(string); ok {
						return coin == sub.Coin
					}
				}
			}
		}
	case "userEvents", "userFills", "orderUpdates", "userFundings", "userNonFundingLedgerUpdates", "webData2":
		if channel == "user" || channel == sub.Type {
			if data, ok := msgData["data"].(map[string]interface{}); ok {
				if user, ok := data["user"].(string); ok {
					return user == sub.User
				}
			}
		}
	case "candle":
		return channel == "candle" // Additional filtering may be needed for coin and interval
	case "activeAssetCtx":
		if channel == "activeAssetCtx" {
			if data, ok := msgData["data"].(map[string]interface{}); ok {
				if coin, ok := data["coin"].(string); ok {
					return coin == sub.Coin
				}
			}
		}
	case "activeAssetData":
		if channel == "activeAssetData" {
			if data, ok := msgData["data"].(map[string]interface{}); ok {
				if user, ok := data["user"].(string); ok {
					if coin, ok := data["coin"].(string); ok {
						return user == sub.User && coin == sub.Coin
					}
				}
			}
		}
	}
	
	return false
}

// Subscribe subscribes to WebSocket channels
func (w *WebsocketManager) Subscribe(subscriptions []types.Subscription, callback func(interface{})) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	if !w.isRunning {
		return fmt.Errorf("WebSocket manager is not running")
	}
	
	for _, sub := range subscriptions {
		// Convert coin names to proper format if needed
		if sub.Coin != "" {
			// Handle name to coin conversion if needed
			sub.Coin = sub.Coin
		}
		
		subKey, err := json.Marshal(sub)
		if err != nil {
			return fmt.Errorf("failed to marshal subscription: %w", err)
		}
		
		w.subscriptions[string(subKey)] = callback
		
		if err := w.sendSubscription(sub); err != nil {
			return fmt.Errorf("failed to send subscription: %w", err)
		}
	}
	
	return nil
}

// Unsubscribe unsubscribes from WebSocket channels
func (w *WebsocketManager) Unsubscribe(subscriptions []types.Subscription) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	if !w.isRunning {
		return fmt.Errorf("WebSocket manager is not running")
	}
	
	for _, sub := range subscriptions {
		subKey, err := json.Marshal(sub)
		if err != nil {
			continue
		}
		
		delete(w.subscriptions, string(subKey))
		
		if err := w.sendUnsubscription(sub); err != nil {
			log.Printf("Failed to send unsubscription: %v", err)
		}
	}
	
	return nil
}

// sendSubscription sends a subscription message
func (w *WebsocketManager) sendSubscription(sub types.Subscription) error {
	message := map[string]interface{}{
		"method": "subscribe",
		"subscription": sub,
	}
	
	return w.conn.WriteJSON(message)
}

// sendUnsubscription sends an unsubscription message
func (w *WebsocketManager) sendUnsubscription(sub types.Subscription) error {
	message := map[string]interface{}{
		"method": "unsubscribe",
		"subscription": sub,
	}
	
	return w.conn.WriteJSON(message)
}

// IsConnected returns true if the WebSocket is connected
func (w *WebsocketManager) IsConnected() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	
	return w.isRunning && w.conn != nil
}

// GetSubscriptions returns a copy of current subscriptions
func (w *WebsocketManager) GetSubscriptions() []types.Subscription {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	
	var subscriptions []types.Subscription
	for subKey := range w.subscriptions {
		var sub types.Subscription
		if err := json.Unmarshal([]byte(subKey), &sub); err == nil {
			subscriptions = append(subscriptions, sub)
		}
	}
	
	return subscriptions
}