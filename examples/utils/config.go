package utils

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	"hyperliquid-go-sdk/pkg/client"
	"hyperliquid-go-sdk/pkg/constants"
	//"hyperliquid-go-sdk/pkg/types"
)

// Config represents the configuration structure
type Config struct {
	SecretKey      string `json:"secret_key"`
	AccountAddress string `json:"account_address"`
	KeystorePath   string `json:"keystore_path"`
}

// Setup initializes the Hyperliquid clients with configuration
func Setup(baseURL string, skipWS bool, perpDexs []string) (string, *client.InfoClient, *client.ExchangeClient, error) {
	// Load configuration
	config, err := loadConfig()
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Get private key
	privateKey, err := getPrivateKey(config)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to get private key: %w", err)
	}

	// Determine account address
	address := config.AccountAddress
	if address == "" {
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			return "", nil, nil, fmt.Errorf("invalid public key type")
		}
		address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	}

	fmt.Printf("Running with account address: %s\n", address)

	// Get wallet address from private key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", nil, nil, fmt.Errorf("invalid public key type")
	}
	walletAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	if address != walletAddress {
		fmt.Printf("Running with agent address: %s\n", walletAddress)
	}

	// Create info client
	infoClient, err := client.NewInfoClient(baseURL, skipWS, nil, nil, perpDexs)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to create info client: %w", err)
	}

	// Check if account has equity
	userState, err := infoClient.UserState(nil, address, "")
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to get user state: %w", err)
	}

	spotUserState, err := infoClient.SpotUserState(nil, address)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to get spot user state: %w", err)
	}

	accountValue := string(userState.MarginSummary.AccountValue)
	if accountValue == "0.0" && len(spotUserState.Balances) == 0 {
		url := strings.Split(baseURL, ".")[1]
		errorMsg := fmt.Sprintf("Not running the example because the provided account has no equity.\n"+
			"No accountValue:\n"+
			"If you think this is a mistake, make sure that %s has a balance on %s.\n"+
			"If address shown is your API wallet address, update the config to specify the address of your account, not the address of the API wallet.",
			address, url)
		return "", nil, nil, fmt.Errorf(errorMsg)
	}

	// Create exchange client
	var accountAddressPtr *string
	if config.AccountAddress != "" {
		accountAddressPtr = &config.AccountAddress
	}

	exchangeClient, err := client.NewExchangeClient(
		privateKey,
		baseURL,
		nil, // meta
		nil, // vault address
		accountAddressPtr,
		nil, // spot meta
		perpDexs,
		nil, // timeout
	)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to create exchange client: %w", err)
	}

	return address, infoClient, exchangeClient, nil
}

// loadConfig loads configuration from config.json
func loadConfig() (*Config, error) {
	configPath := filepath.Join("examples", "config.json")

	// Check if config.json exists, if not, try the example file
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config.json not found. Please copy config.json.example to config.json and fill in your details")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// getPrivateKey gets the private key from configuration
func getPrivateKey(config *Config) (*ecdsa.PrivateKey, error) {
	if config.SecretKey != "" {
		// Remove 0x prefix if present
		secretKey := config.SecretKey
		if strings.HasPrefix(secretKey, "0x") {
			secretKey = secretKey[2:]
		}

		privateKey, err := crypto.HexToECDSA(secretKey)
		if err != nil {
			return nil, fmt.Errorf("invalid secret key: %w", err)
		}
		return privateKey, nil
	}

	if config.KeystorePath != "" {
		// Keystore functionality would need to be implemented
		// This is a placeholder for keystore support
		return nil, fmt.Errorf("keystore support not implemented yet")
	}

	return nil, fmt.Errorf("no secret key or keystore path provided in config")
}

// SetupTestnet creates a testnet setup
func SetupTestnet(skipWS bool, perpDexs []string) (string, *client.InfoClient, *client.ExchangeClient, error) {
	return Setup(constants.TestnetAPIURL, skipWS, perpDexs)
}

// SetupMainnet creates a mainnet setup
func SetupMainnet(skipWS bool, perpDexs []string) (string, *client.InfoClient, *client.ExchangeClient, error) {
	return Setup(constants.MainnetAPIURL, skipWS, perpDexs)
}
