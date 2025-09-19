package tests

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"hyperliquid-go-sdk/pkg/signing"
	"hyperliquid-go-sdk/pkg/types"
)

func TestPhantomAgentCreation(t *testing.T) {
	timestamp := int64(1677777606040)

	orderRequest := types.OrderRequest{
		Coin:       "ETH",
		IsBuy:      true,
		Sz:         0.0147,
		LimitPx:    1670.1,
		ReduceOnly: false,
		OrderType:  types.OrderType{Limit: &types.LimitOrderType{Tif: "Ioc"}},
	}

	orderWire, err := signing.OrderRequestToOrderWire(orderRequest, 4)
	if err != nil {
		t.Fatalf("Failed to convert order request to wire: %v", err)
	}

	orderAction := signing.OrderWiresToOrderAction([]types.OrderWire{orderWire}, nil)

	hash, err := signing.ActionHash(orderAction, nil, timestamp, nil)
	if err != nil {
		t.Fatalf("Failed to compute action hash: %v", err)
	}

	phantomAgent := signing.ConstructPhantomAgent(hash, true)
	expectedConnectionID := "0x0fcbeda5ae3c4950a548021552a4fea2226858c4453571bf3f24ba017eac2908"

	if phantomAgent.ConnectionID != expectedConnectionID {
		t.Errorf("Expected connection ID %s, got %s", expectedConnectionID, phantomAgent.ConnectionID)
	}
}

func TestL1ActionSigning(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("0123456789012345678901234567890123456789012345678901234567890123")
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}

	action := map[string]interface{}{
		"type": "dummy",
		"num":  int64(100000000000),
	}

	// Test mainnet signing
	signature, err := signing.SignL1Action(privateKey, action, nil, 0, nil, true)
	if err != nil {
		t.Fatalf("Failed to sign L1 action: %v", err)
	}

	expectedR := "0x53749d5b30552aeb2fca34b530185976545bb22d0b3ce6f62e31be961a59298"
	expectedS := "0x755c40ba9bf05223521753995abb2f73ab3229be8ec921f350cb447e384d8ed8"
	expectedV := 27

	if signature.R != expectedR {
		t.Errorf("Expected R %s, got %s", expectedR, signature.R)
	}
	if signature.S != expectedS {
		t.Errorf("Expected S %s, got %s", expectedS, signature.S)
	}
	if signature.V != expectedV {
		t.Errorf("Expected V %d, got %d", expectedV, signature.V)
	}

	// Test testnet signing
	signatureTestnet, err := signing.SignL1Action(privateKey, action, nil, 0, nil, false)
	if err != nil {
		t.Fatalf("Failed to sign L1 action for testnet: %v", err)
	}

	expectedRTestnet := "0x542af61ef1f429707e3c76c5293c80d01f74ef853e34b76efffcb57e574f9510"
	expectedSTestnet := "0x17b8b32f086e8cdede991f1e2c529f5dd5297cbe8128500e00cbaf766204a613"
	expectedVTestnet := 28

	if signatureTestnet.R != expectedRTestnet {
		t.Errorf("Expected R %s, got %s", expectedRTestnet, signatureTestnet.R)
	}
	if signatureTestnet.S != expectedSTestnet {
		t.Errorf("Expected S %s, got %s", expectedSTestnet, signatureTestnet.S)
	}
	if signatureTestnet.V != expectedVTestnet {
		t.Errorf("Expected V %d, got %d", expectedVTestnet, signatureTestnet.V)
	}
}

func TestFloatToIntForHashing(t *testing.T) {
	testCases := []struct {
		input    float64
		expected int64
	}{
		{123123123123, 12312312312300000000},
		{0.00001231, 1231},
		{1.033, 103300000},
	}

	for _, tc := range testCases {
		result, err := signing.FloatToIntForHashing(tc.input)
		if err != nil {
			t.Errorf("Failed to convert %f: %v", tc.input, err)
			continue
		}
		if result != tc.expected {
			t.Errorf("Expected %d for input %f, got %d", tc.expected, tc.input, result)
		}
	}

	// Test error case
	_, err := signing.FloatToIntForHashing(0.000012312312)
	if err == nil {
		t.Error("Expected error for precision that causes rounding, but got none")
	}
}

func TestFloatToWire(t *testing.T) {
	testCases := []struct {
		input    float64
		expected string
	}{
		{1234.5, "1234.5"},
		{0.001234, "0.001234"},
		{123456.0, "123456"},
	}

	for _, tc := range testCases {
		result, err := signing.FloatToWire(tc.input)
		if err != nil {
			t.Errorf("Failed to convert %f to wire: %v", tc.input, err)
			continue
		}
		if result != tc.expected {
			t.Errorf("Expected %s for input %f, got %s", tc.expected, tc.input, result)
		}
	}
}

func TestCloidCreation(t *testing.T) {
	// Test creating from string
	cloidStr := "0x00000000000000000000000000000001"
	cloid, err := types.NewCloidFromString(cloidStr)
	if err != nil {
		t.Fatalf("Failed to create cloid from string: %v", err)
	}

	if cloid.Raw() != cloidStr {
		t.Errorf("Expected %s, got %s", cloidStr, cloid.Raw())
	}

	// Test creating from int
	cloidInt := types.NewCloidFromInt(1)
	expectedStr := "0x00000000000000000000000000000001"
	if cloidInt.Raw() != expectedStr {
		t.Errorf("Expected %s, got %s", expectedStr, cloidInt.Raw())
	}

	// Test invalid cloid
	_, err = types.NewCloidFromString("0x123") // Too short
	if err == nil {
		t.Error("Expected error for invalid cloid, but got none")
	}

	_, err = types.NewCloidFromString("123") // No 0x prefix
	if err == nil {
		t.Error("Expected error for invalid cloid without 0x prefix, but got none")
	}
}
