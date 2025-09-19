package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"hyperliquid-go-sdk/pkg/constants"
	"hyperliquid-go-sdk/pkg/errors"
)

// Client represents the HTTP API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	timeout    time.Duration
}

// NewClient creates a new API client
func NewClient(baseURL string, timeout *time.Duration) *Client {
	if baseURL == "" {
		baseURL = constants.MainnetAPIURL
	}

	clientTimeout := 30 * time.Second
	if timeout != nil {
		clientTimeout = *timeout
	}

	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: clientTimeout,
		},
		timeout: clientTimeout,
	}
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error string      `json:"error,omitempty"`
	Code  string      `json:"code,omitempty"`
	Msg   string      `json:"msg,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// Post makes a POST request to the API
func (c *Client) Post(ctx context.Context, urlPath string, payload interface{}) ([]byte, error) {
	url := c.BaseURL + urlPath

	var reqBody []byte
	var err error

	if payload != nil {
		reqBody, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request payload: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := c.handleError(resp.StatusCode, body); err != nil {
		return nil, err
	}

	return body, nil
}

// handleError checks for API errors and returns appropriate error types
func (c *Client) handleError(statusCode int, body []byte) error {
	if statusCode < 400 {
		return nil
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		// If we can't parse the error response, return a generic API error
		return errors.NewAPIError(statusCode, "", string(body), nil)
	}

	// Determine the error message
	message := errorResp.Error
	if message == "" {
		message = errorResp.Msg
	}
	if message == "" {
		message = string(body)
	}

	return errors.NewAPIError(statusCode, errorResp.Code, message, errorResp.Data)
}

// Get makes a GET request to the API (if needed)
func (c *Client) Get(ctx context.Context, urlPath string) ([]byte, error) {
	url := c.BaseURL + urlPath

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := c.handleError(resp.StatusCode, body); err != nil {
		return nil, err
	}

	return body, nil
}
