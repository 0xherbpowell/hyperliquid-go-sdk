package client

import (
	"time"

	"hyperliquid-go-sdk/internal/api"
	"hyperliquid-go-sdk/pkg/constants"
)

// BaseClient provides the foundation for all Hyperliquid clients
type BaseClient struct {
	apiClient *api.Client
	baseURL   string
}

// NewBaseClient creates a new base client
func NewBaseClient(baseURL string, timeout *time.Duration) (*BaseClient, error) {
	if baseURL == "" {
		baseURL = constants.MainnetAPIURL
	}

	apiClient := api.NewClient(baseURL, timeout)

	return &BaseClient{
		apiClient: apiClient,
		baseURL:   baseURL,
	}, nil
}

// GetBaseURL returns the base URL of the client
func (c *BaseClient) GetBaseURL() string {
	return c.baseURL
}

// IsMainnet returns true if the client is connected to mainnet
func (c *BaseClient) IsMainnet() bool {
	return c.baseURL == constants.MainnetAPIURL
}

// IsTestnet returns true if the client is connected to testnet
func (c *BaseClient) IsTestnet() bool {
	return c.baseURL == constants.TestnetAPIURL
}
