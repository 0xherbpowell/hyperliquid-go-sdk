package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"hyperliquid-go-sdk/pkg/utils"
)

// API represents the HTTP client for interacting with Hyperliquid API
type API struct {
	BaseURL    string
	HTTPClient *http.Client
	timeout    time.Duration
}

// NewAPI creates a new API client
func NewAPI(baseURL string, timeout *time.Duration) *API {
	if baseURL == "" {
		baseURL = utils.MainnetAPIURL
	}

	clientTimeout := time.Duration(utils.DefaultTimeoutSeconds) * time.Second
	if timeout != nil {
		clientTimeout = *timeout
	}

	return &API{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: clientTimeout,
		},
		timeout: clientTimeout,
	}
}

// Post makes a POST request to the API
func (a *API) Post(urlPath string, payload interface{}) (map[string]interface{}, error) {
	if payload == nil {
		payload = map[string]interface{}{}
	}

	url := a.BaseURL + urlPath

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle HTTP errors
	if err := a.handleException(resp, body); err != nil {
		return nil, err
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("Could not parse JSON: %s", string(body)),
		}, nil
	}

	return result, nil
}

// handleException handles HTTP errors and creates appropriate error types
func (a *API) handleException(resp *http.Response, body []byte) error {
	statusCode := resp.StatusCode
	if statusCode < 400 {
		return nil
	}

	if statusCode >= 400 && statusCode < 500 {
		var errResp map[string]interface{}
		if err := json.Unmarshal(body, &errResp); err != nil {
			return utils.NewClientError(statusCode, nil, string(body), resp.Header, nil)
		}

		if errResp == nil {
			return utils.NewClientError(statusCode, nil, string(body), resp.Header, nil)
		}

		var code *string
		var msg string
		var data interface{}

		if codeVal, ok := errResp["code"].(string); ok {
			code = &codeVal
		}

		if msgVal, ok := errResp["msg"].(string); ok {
			msg = msgVal
		} else {
			msg = string(body)
		}

		if dataVal, ok := errResp["data"]; ok {
			data = dataVal
		}

		return utils.NewClientError(statusCode, code, msg, resp.Header, data)
	}

	return utils.NewServerError(statusCode, string(body))
}

// IsMainnet returns true if the client is connected to mainnet
func (a *API) IsMainnet() bool {
	return a.BaseURL == utils.MainnetAPIURL
}

// IsTestnet returns true if the client is connected to testnet
func (a *API) IsTestnet() bool {
	return a.BaseURL == utils.TestnetAPIURL
}
