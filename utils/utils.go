package utils

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
	"gopkg.in/yaml.v3"
)

// IsCompiled checks if the contract is compiled. If not, it compiles, publishes,
// and submits the transaction for the contract.
func IsCompiled(isCompiled bool, account *aptos.Account, client *aptos.Client) {
	if !isCompiled {
		// Define named addresses for the contract compilation
		NamedAddress := []NamedAddress{{Name: "Admin", Address: account.Address.String()}}

		// Compile the Move package
		err := compilePackage("./move", "./move/Signature-Verification.json", NamedAddress)
		if err != nil {
			panic("Error Compiling Contract: " + err.Error())
		}
		fmt.Println("Compilation completed successfully.")

		// Prepare the package for publishing
		filePath := "./move/Signature-Verification.json" // Path to compiled package metadata
		metadataBytes, byteCode, err := getPackageBytesToPublish(filePath)
		if err != nil {
			fmt.Println("Error during package preparation:", err)
			return
		}

		// Create a transaction payload for publishing
		tx, err := aptos.PublishPackagePayloadFromJsonFile(metadataBytes, byteCode)
		if err != nil {
			fmt.Println("Error creating transaction payload:", err)
			return
		}

		// Build the raw transaction
		rawTx, err := client.BuildTransaction(account.Address, *tx)
		if err != nil {
			panic("Failed to build transaction: " + err.Error())
		}

		// Simulate the transaction
		simulationResult, err := client.SimulateTransaction(rawTx, account)
		if err != nil {
			panic("Failed to simulate transaction: " + err.Error())
		}
		fmt.Printf("\n=== Simulation ===\n")
		fmt.Printf("Gas unit price: %d\n", simulationResult[0].GasUnitPrice)
		fmt.Printf("Gas used: %d\n", simulationResult[0].GasUsed)
		fmt.Printf("Total gas fee: %d\n", simulationResult[0].GasUsed*simulationResult[0].GasUnitPrice)
		fmt.Printf("Status: %s\n", simulationResult[0].VmStatus)

		// Sign the transaction
		signedTxn, err := rawTx.SignedTransaction(account)
		if err != nil {
			fmt.Println("Error Signing Transaction:", err.Error())
			return
		}

		// Submit the transaction
		submitResult, err := client.SubmitTransaction(signedTxn)
		if err != nil {
			fmt.Println("Error Submitting Transaction:", err.Error())
			return
		}
		txnHash := submitResult.Hash
		fmt.Println("Transaction submitted with hash:", txnHash)

		// Wait for the transaction to be committed
		committedTx, err := client.WaitForTransaction(txnHash)
		if err != nil {
			fmt.Println("Error Waiting for the Transaction:", err.Error())
			return
		}

		// Check the transaction status
		if !committedTx.Success {
			fmt.Println("Transaction failed with VM status:", committedTx.VmStatus)
			return
		}

		fmt.Println("Contract Published Successfully.\nTransaction Successful:", committedTx.Success)
	}
}

// ReadYAML reads the YAML configuration file and parses it into a struct.
func ReadYAML() *YAML {
	data, err := os.ReadFile("./move/.aptos/config.yaml")
	if err != nil {
		panic("Error reading YAML file: " + err.Error())
	}
	var config YAML
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic("Error parsing YAML file: " + err.Error())
	}
	return &config
}

// GetAccount creates an Aptos account using the provided private key in hex format.
func GetAccount(privateKeyHex string) *aptos.Account {
	key, _ := crypto.FormatPrivateKey(privateKeyHex, crypto.PrivateKeyVariantEd25519)
	privateKey := &crypto.Ed25519PrivateKey{}
	err := privateKey.FromHex(key)
	if err != nil {
		panic("Failed to parse private key:" + err.Error())
	}
	account, err := aptos.NewAccountFromSigner(privateKey)
	if err != nil {
		panic("Failed to generate account from private key:" + err.Error())
	}
	return account
}

// SignMessage signs a message hash using the provided Ed25519 private key.
func SignMessage(privateKey ed25519.PrivateKey, messageHash []byte) []byte {
	return ed25519.Sign(privateKey, messageHash)
}

// HexToEd25519PrivateKey converts a hex-encoded string to an Ed25519 private key.
func HexToEd25519PrivateKey(hexKey string) ed25519.PrivateKey {
	// Remove "0x" prefix if present
	hexKey = strings.TrimPrefix(hexKey, "0x")

	// Decode the hex string into bytes
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		log.Fatalf("Invalid hex key: %v", err)
	}

	// Ensure the key is at least 32 bytes
	if len(keyBytes) < 32 {
		log.Fatalf("Hex key is too short, must be at least 32 bytes: got %d bytes", len(keyBytes))
	}

	return ed25519.NewKeyFromSeed(keyBytes[:32])
}

// parseHexToBytes converts a hex string to a byte slice.
func parseHexToBytes(value interface{}) ([]byte, error) {
	strValue, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf("expected string, got %T", value)
	}
	return hex.DecodeString(strValue[2:]) // Remove "0x" prefix before decoding
}

// parseHexArrayToBytesArray converts an array of hex strings to a slice of byte slices.
func parseHexArrayToBytesArray(value interface{}) ([][]byte, error) {
	arrValue, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", value)
	}

	result := make([][]byte, len(arrValue))
	for i, v := range arrValue {
		strValue, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("expected string in array, got %T", v)
		}
		decoded, err := hex.DecodeString(strValue[2:]) // Remove "0x" prefix before decoding
		if err != nil {
			return nil, fmt.Errorf("failed to decode hex string: %v", err)
		}
		result[i] = decoded
	}
	return result, nil
}
