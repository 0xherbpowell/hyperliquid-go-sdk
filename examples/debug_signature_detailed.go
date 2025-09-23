package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	fmt.Println("=== Detailed Signature Debug ===")

	// Get environment variables - using testnet
	testnetPrivateKey := os.Getenv("HYPERLIQUID_TESTNET_PRIVATE_KEY")
	testnetAgentPrivateKey := os.Getenv("HYPERLIQUID_TESTNET_AGENT_PRIVATE_KEY")
	userAddress := os.Getenv("HYPERLIQUID_TESTNET_USER_ADDRESS")

	if testnetPrivateKey == "" {
		log.Fatal("HYPERLIQUID_TESTNET_PRIVATE_KEY environment variable is required")
	}
	if testnetAgentPrivateKey == "" {
		log.Fatal("HYPERLIQUID_TESTNET_AGENT_PRIVATE_KEY environment variable is required")
	}
	if userAddress == "" {
		log.Fatal("HYPERLIQUID_TESTNET_USER_ADDRESS environment variable is required")
	}

	// Parse keys
	userKey, err := crypto.HexToECDSA(testnetPrivateKey)
	if err != nil {
		log.Fatal("Error parsing user key:", err)
	}

	agentKey, err := crypto.HexToECDSA(testnetAgentPrivateKey)
	if err != nil {
		log.Fatal("Error parsing agent key:", err)
	}

	userAddr := crypto.PubkeyToAddress(userKey.PublicKey).Hex()
	agentAddr := crypto.PubkeyToAddress(agentKey.PublicKey).Hex()

	fmt.Printf("User private key: %s\n", testnetPrivateKey)
	fmt.Printf("Agent private key: %s\n", testnetAgentPrivateKey)
	fmt.Printf("Expected user address: %s\n", userAddress)
	fmt.Printf("Derived user address: %s\n", userAddr)
	fmt.Printf("Derived agent address: %s\n", agentAddr)

	// Create test action exactly like the failing request
	action := map[string]interface{}{
		"type":     "order",
		"grouping": "na",
		"orders": []map[string]interface{}{
			{
				"a": 4,       // ETH asset ID (as seen in the log)
				"b": true,    // buy
				"p": "1500",  // price as string
				"s": "0.001", // size as string
				"r": false,   // reduce only
				"t": map[string]interface{}{
					"limit": map[string]interface{}{
						"tif": "Gtc",
					},
				},
			},
		},
	}

	// Use same parameters as the failing request
	var vaultAddress *string = nil
	nonce := int64(1758570668100) // Use the same nonce from the log
	var expiresAfter *int64 = nil
	isMainnet := false

	fmt.Printf("\nUsing nonce: %d\n", nonce)
	fmt.Println("Action:", action)

	// Step 1: Calculate action hash manually
	fmt.Println("\n=== Step 1: Action Hash Calculation ===")
	
	// Msgpack the action
	actionBytes, err := msgpack.Marshal(action)
	if err != nil {
		log.Fatal("Error marshaling action:", err)
	}
	fmt.Printf("Msgpack bytes: %x\n", actionBytes)

	// Add nonce (8 bytes big endian)
	nonceBytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		nonceBytes[7-i] = byte(nonce >> (i * 8))
	}
	fmt.Printf("Nonce bytes: %x\n", nonceBytes)
	
	data := append(actionBytes, nonceBytes...)

	// Add vault address (nil = 0x00)
	data = append(data, 0x00)
	fmt.Printf("After vault: %x\n", data)

	// Add expires after (nil = add nothing based on our fix)
	fmt.Printf("After expires (no change): %x\n", data)

	// Keccak hash
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(data)
	actionHash := hasher.Sum(nil)
	fmt.Printf("Manual action hash: %x\n", actionHash)

	// Compare with utils function
	utilsHash, err := utils.ActionHash(action, vaultAddress, nonce, expiresAfter)
	if err != nil {
		log.Fatal("Error with utils.ActionHash:", err)
	}
	fmt.Printf("Utils action hash: %x\n", utilsHash)
	fmt.Printf("Hashes match: %t\n", hex.EncodeToString(actionHash) == hex.EncodeToString(utilsHash))

	// Step 2: Phantom agent
	fmt.Println("\n=== Step 2: Phantom Agent ===")
	phantomAgent := utils.ConstructPhantomAgent(actionHash, isMainnet)
	fmt.Printf("Phantom agent source: %s\n", phantomAgent["source"])
	fmt.Printf("Phantom agent connectionId: %x\n", phantomAgent["connectionId"])

	// Step 3: EIP712 payload
	fmt.Println("\n=== Step 3: EIP712 Payload ===")
	eip712Data := utils.L1Payload(phantomAgent)
	fmt.Printf("Domain: %+v\n", eip712Data.Domain)
	fmt.Printf("Message: %+v\n", eip712Data.Message)
	
	// Step 4: Manual signing with both keys
	fmt.Println("\n=== Step 4: Manual Signing ===")
	
	keys := []*ecdsa.PrivateKey{userKey, agentKey}
	keyNames := []string{"user", "agent"}
	keyAddrs := []string{userAddr, agentAddr}

	for i, key := range keys {
		fmt.Printf("\n--- Testing %s key (%s) ---\n", keyNames[i], keyAddrs[i])
		
		// Hash the message
		messageHash, err := eip712Data.HashStruct(eip712Data.PrimaryType, eip712Data.Message)
		if err != nil {
			log.Printf("Error hashing message: %v", err)
			continue
		}
		fmt.Printf("Message hash: %x\n", messageHash)

		// Hash the domain
		domainHash, err := eip712Data.HashStruct("EIP712Domain", eip712Data.Domain.Map())
		if err != nil {
			log.Printf("Error hashing domain: %v", err)
			continue
		}
		fmt.Printf("Domain hash: %x\n", domainHash)

		// Create final hash
		finalHash := crypto.Keccak256([]byte("\x19\x01"), domainHash, messageHash)
		fmt.Printf("Final hash: %x\n", finalHash)

		// Sign
		signature, err := crypto.Sign(finalHash, key)
		if err != nil {
			log.Printf("Error signing: %v", err)
			continue
		}

		r := signature[:32]
		s := signature[32:64]
		v := signature[64] + 27

		fmt.Printf("Raw signature: %x\n", signature)
		fmt.Printf("r: %x\n", r)
		fmt.Printf("s: %x\n", s)
		fmt.Printf("v: %d\n", v)

		// Test recovery
		recoveredPubkey, err := crypto.SigToPub(finalHash, signature)
		if err != nil {
			log.Printf("Error recovering: %v", err)
			continue
		}

		recoveredAddr := crypto.PubkeyToAddress(*recoveredPubkey).Hex()
		fmt.Printf("Expected: %s\n", keyAddrs[i])
		fmt.Printf("Recovered: %s\n", recoveredAddr)
		fmt.Printf("Match: %t\n", recoveredAddr == keyAddrs[i])
		
		// Test with SDK signing function
		fmt.Printf("\n--- SDK signing with %s key ---\n", keyNames[i])
		var accountAddr *string = nil
		if i == 1 { // agent key
			accountAddr = &userAddress
		}
		
		sdkSig, err := utils.SignL1ActionWithAccount(key, action, vaultAddress, nonce, expiresAfter, isMainnet, accountAddr)
		if err != nil {
			log.Printf("SDK signing error: %v", err)
			continue
		}
		
		fmt.Printf("SDK signature: %+v\n", sdkSig)
		
		// Parse and test recovery of SDK signature
		rHex := sdkSig["r"].(string)
		sHex := sdkSig["s"].(string)
		vInt := sdkSig["v"].(int)
		
		rBytes, _ := hex.DecodeString(rHex[2:])
		sBytes, _ := hex.DecodeString(sHex[2:])
		
		sdkSigBytes := make([]byte, 65)
		copy(sdkSigBytes[:32], rBytes)
		copy(sdkSigBytes[32:64], sBytes)
		sdkSigBytes[64] = byte(vInt - 27)
		
		sdkRecoveredPubkey, err := crypto.SigToPub(finalHash, sdkSigBytes)
		if err != nil {
			log.Printf("SDK recovery error: %v", err)
			continue
		}
		
		sdkRecoveredAddr := crypto.PubkeyToAddress(*sdkRecoveredPubkey).Hex()
		fmt.Printf("SDK recovered: %s\n", sdkRecoveredAddr)
		fmt.Printf("SDK match: %t\n", sdkRecoveredAddr == keyAddrs[i])
	}
}