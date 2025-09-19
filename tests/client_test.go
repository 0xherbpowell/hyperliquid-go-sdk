package tests

import (
	//"context"
	"testing"
	"time"

	"hyperliquid-go-sdk/pkg/client"
	"hyperliquid-go-sdk/pkg/constants"
	"hyperliquid-go-sdk/pkg/types"
)

func TestInfoClientCreation(t *testing.T) {
	infoClient, err := client.NewInfoClient(constants.TestnetAPIURL, true, nil, nil, nil)
	if err != nil {
		t.Fatalf("Failed to create info client: %v", err)
	}

	if infoClient == nil {
		t.Error("Info client should not be nil")
	}

	if !infoClient.IsTestnet() {
		t.Error("Client should be connected to testnet")
	}

	if infoClient.IsMainnet() {
		t.Error("Client should not be connected to mainnet")
	}
}

func TestBaseClientURLs(t *testing.T) {
	testCases := []struct {
		baseURL   string
		isMainnet bool
		isTestnet bool
	}{
		{constants.MainnetAPIURL, true, false},
		{constants.TestnetAPIURL, false, true},
		{constants.LocalAPIURL, false, false},
	}

	for _, tc := range testCases {
		baseClient, err := client.NewBaseClient(tc.baseURL, nil)
		if err != nil {
			t.Fatalf("Failed to create base client for URL %s: %v", tc.baseURL, err)
		}

		if baseClient.IsMainnet() != tc.isMainnet {
			t.Errorf("Expected IsMainnet() to be %v for URL %s", tc.isMainnet, tc.baseURL)
		}

		if baseClient.IsTestnet() != tc.isTestnet {
			t.Errorf("Expected IsTestnet() to be %v for URL %s", tc.isTestnet, tc.baseURL)
		}

		if baseClient.GetBaseURL() != tc.baseURL {
			t.Errorf("Expected base URL %s, got %s", tc.baseURL, baseClient.GetBaseURL())
		}
	}
}

func TestClientTimeout(t *testing.T) {
	timeout := 10 * time.Second
	baseClient, err := client.NewBaseClient(constants.TestnetAPIURL, &timeout)
	if err != nil {
		t.Fatalf("Failed to create base client with timeout: %v", err)
	}

	if baseClient == nil {
		t.Error("Base client should not be nil")
	}
}

func TestOrderRequestValidation(t *testing.T) {
	testCases := []struct {
		name         string
		orderRequest types.OrderRequest
		shouldError  bool
	}{
		{
			name: "valid limit order",
			orderRequest: types.OrderRequest{
				Coin:       "ETH",
				IsBuy:      true,
				Sz:         0.1,
				LimitPx:    2000.0,
				OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: constants.TifGtc}},
				ReduceOnly: false,
			},
			shouldError: false,
		},
		{
			name: "valid trigger order",
			orderRequest: types.OrderRequest{
				Coin:    "ETH",
				IsBuy:   true,
				Sz:      0.1,
				LimitPx: 2000.0,
				OrderType: types.OrderType{
					Trigger: &types.TriggerOrderType{
						TriggerPx: 1900.0,
						IsMarket:  true,
						Tpsl:      constants.TpslSl,
					},
				},
				ReduceOnly: true,
			},
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test that the order request structure is valid
			if tc.orderRequest.Coin == "" {
				if !tc.shouldError {
					t.Error("Order request should have a coin specified")
				}
			}
			if tc.orderRequest.Sz <= 0 {
				if !tc.shouldError {
					t.Error("Order request should have a positive size")
				}
			}
			if tc.orderRequest.LimitPx <= 0 {
				if !tc.shouldError {
					t.Error("Order request should have a positive limit price")
				}
			}
		})
	}
}

func TestCloidFunctionality(t *testing.T) {
	// Test string creation
	cloidStr := "0x00000000000000000000000000000001"
	cloid1, err := types.NewCloidFromString(cloidStr)
	if err != nil {
		t.Fatalf("Failed to create cloid from string: %v", err)
	}

	// Test integer creation
	cloid2 := types.NewCloidFromInt(1)

	// Both should be equal
	if cloid1.Raw() != cloid2.Raw() {
		t.Errorf("Cloids should be equal: %s vs %s", cloid1.Raw(), cloid2.Raw())
	}

	// Test string representation
	if cloid1.String() != cloidStr {
		t.Errorf("Expected string %s, got %s", cloidStr, cloid1.String())
	}

	// Test invalid cloids
	invalidCloids := []string{
		"0x123", // Too short
		"123",   // No 0x prefix
		"0x00000000000000000000000000000000000001", // Too long
		"0xgggggggggggggggggggggggggggggggg",       // Invalid hex
	}

	for _, invalid := range invalidCloids {
		_, err := types.NewCloidFromString(invalid)
		if err == nil {
			t.Errorf("Expected error for invalid cloid %s, but got none", invalid)
		}
	}
}

func TestActionTypes(t *testing.T) {
	// Test that all action type constants are defined
	actionTypes := []string{
		constants.ActionOrder,
		constants.ActionCancel,
		constants.ActionCancelByCloid,
		constants.ActionBatchModify,
		constants.ActionScheduleCancel,
		constants.ActionUpdateLeverage,
		constants.ActionUpdateIsolatedMargin,
		constants.ActionSetReferrer,
		constants.ActionCreateSubAccount,
		constants.ActionUsdClassTransfer,
		constants.ActionSendAsset,
		constants.ActionSubAccountTransfer,
		constants.ActionVaultTransfer,
		constants.ActionUsdSend,
		constants.ActionSpotSend,
		constants.ActionTokenDelegate,
		constants.ActionWithdraw3,
		constants.ActionApproveAgent,
		constants.ActionApproveBuilderFee,
		constants.ActionConvertToMultiSigUser,
		constants.ActionSpotDeploy,
		constants.ActionPerpDeploy,
		constants.ActionMultiSig,
		constants.ActionEvmUserModify,
		constants.ActionNoop,
	}

	for _, actionType := range actionTypes {
		if actionType == "" {
			t.Error("Action type should not be empty")
		}
	}
}

func TestTimestampGeneration(t *testing.T) {
	timestamp1 := types.GetTimestampMs()
	time.Sleep(1 * time.Millisecond)
	timestamp2 := types.GetTimestampMs()

	if timestamp2 <= timestamp1 {
		t.Error("Second timestamp should be greater than first")
	}

	// Check that timestamp is reasonable (not too far in past or future)
	now := time.Now().UnixMilli()
	if timestamp2 < now-1000 || timestamp2 > now+1000 {
		t.Errorf("Timestamp %d seems unreasonable compared to now %d", timestamp2, now)
	}
}

func TestSideConstants(t *testing.T) {
	if types.SideAsk != "A" {
		t.Errorf("Expected SideAsk to be 'A', got %s", types.SideAsk)
	}
	if types.SideBid != "B" {
		t.Errorf("Expected SideBid to be 'B', got %s", types.SideBid)
	}
}

func TestOrderTypeConstants(t *testing.T) {
	if constants.TifAlo != "Alo" {
		t.Errorf("Expected TifAlo to be 'Alo', got %s", constants.TifAlo)
	}
	if constants.TifIoc != "Ioc" {
		t.Errorf("Expected TifIoc to be 'Ioc', got %s", constants.TifIoc)
	}
	if constants.TifGtc != "Gtc" {
		t.Errorf("Expected TifGtc to be 'Gtc', got %s", constants.TifGtc)
	}
	if constants.TpslTp != "tp" {
		t.Errorf("Expected TpslTp to be 'tp', got %s", constants.TpslTp)
	}
	if constants.TpslSl != "sl" {
		t.Errorf("Expected TpslSl to be 'sl', got %s", constants.TpslSl)
	}
}

// Benchmark tests
func BenchmarkCloidCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		types.NewCloidFromInt(int64(i))
	}
}

func BenchmarkTimestampGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		types.GetTimestampMs()
	}
}
