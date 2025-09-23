package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"

	"hyperliquid-go-sdk/pkg/utils"
)

func main() {
	fmt.Println("=== Reference SDK Comparison Debug ===")

	// Get environment variables - using testnet
	testnetAgentPrivateKey := "06e10c1cb33b369c878ec8f3d51523f2bdd3a36f02fcb6c29e0867903e17927e"

	// Parse agent key
	agentKey, err := crypto.HexToECDSA(testnetAgentPrivateKey)
	if err != nil {
		log.Fatal("Error parsing agent key:", err)
	}

	agentAddr := crypto.PubkeyToAddress(agentKey.PublicKey).Hex()
	fmt.Printf("Agent address: %s\n", agentAddr)

	// Create the same test action as the failing request
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

	// Use parameters from the failing request
	vaultAddress := ""
	nonce := int64(1758611180225) // Use the same nonce from the log
	var expiresAfter *int64 = nil
	isMainnet := false

	fmt.Printf("Using nonce: %d\n", nonce)
	fmt.Printf("Action: %+v\n", action)

	// Test our signing approach
	fmt.Println("\n=== Our SDK Signing ===")
	ourSig, err := utils.SignL1Action(agentKey, action, vaultAddress, nonce, expiresAfter, isMainnet)
	if err != nil {
		log.Printf("Our signing failed: %v", err)
	} else {
		fmt.Printf("Our signature: R=%s, S=%s, V=%d\n", ourSig.R, ourSig.S, ourSig.V)

		// Test signature recovery
		testSignatureRecovery(agentKey, action, vaultAddress, nonce, expiresAfter, isMainnet, agentAddr)
	}
}

func testSignatureRecovery(privateKey *ecdsa.PrivateKey, action interface{}, vaultAddress string, nonce int64, expiresAfter *int64, isMainnet bool, expectedAddr string) {
	fmt.Println("\n=== Testing Signature Recovery ===")

	// Step 1: Create action hash
	hash, err := utils.ActionHash(action, &vaultAddress, nonce, expiresAfter)
	if err != nil {
		log.Printf("Error creating action hash: %v", err)
		return
	}
	fmt.Printf("Action hash: %x\n", hash)

	// Step 2: Create phantom agent
	phantomAgent := utils.ConstructPhantomAgent(hash, isMainnet)
	fmt.Printf("Phantom agent: %+v\n", phantomAgent)

	// Step 3: Create EIP712 payload
	typedData := utils.L1Payload(phantomAgent)
	fmt.Printf("Domain: %+v\n", typedData.Domain)
	fmt.Printf("Message: %+v\n", typedData.Message)

	// Step 4: Hash the message
	messageHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		log.Printf("Error hashing message: %v", err)
		return
	}
	fmt.Printf("Message hash: %x\n", messageHash)

	// Hash the domain
	domainHash, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		log.Printf("Error hashing domain: %v", err)
		return
	}
	fmt.Printf("Domain hash: %x\n", domainHash)

	// Create final hash
	finalHash := crypto.Keccak256([]byte("\x19\x01"), domainHash, messageHash)
	fmt.Printf("Final hash: %x\n", finalHash)

	// Sign
	signature, err := crypto.Sign(finalHash, privateKey)
	if err != nil {
		log.Printf("Error signing: %v", err)
		return
	}

	fmt.Printf("Raw signature: %x\n", signature)

	// Test recovery
	recoveredPubkey, err := crypto.SigToPub(finalHash, signature)
	if err != nil {
		log.Printf("Error recovering: %v", err)
		return
	}

	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubkey).Hex()
	fmt.Printf("Expected: %s\n", expectedAddr)
	fmt.Printf("Recovered: %s\n", recoveredAddr)
	fmt.Printf("Match: %t\n", recoveredAddr == expectedAddr)

	if recoveredAddr != expectedAddr {
		fmt.Println("❌ Signature recovery failed - addresses don't match")
	} else {
		fmt.Println("✅ Signature recovery successful - addresses match")
	}
}
