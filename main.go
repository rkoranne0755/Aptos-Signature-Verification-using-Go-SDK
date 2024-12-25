package main

import (
	"Signature-Verification/utils"       // Custom utility package
	"crypto/sha256"                      // For hashing the message
	"fmt"                                // For printing output
	"github.com/aptos-labs/aptos-go-sdk" // Aptos Go SDK
)

// main is the entry point of the application
func main() {
	fmt.Println("Signature Verification Process Started")

	// Step 1: Initialize Aptos Client
	client, err := aptos.NewClient(aptos.DevnetConfig)
	if err != nil {
		panic("Error creating Aptos Client: " + err.Error())
	}

	// Step 2: Read configuration from YAML
	yaml := utils.ReadYAML()

	// Step 3: Retrieve Admin Account Information
	admin := utils.GetAccount(yaml.Profiles.Default.PrivateKey)

	utils.IsCompiled(false, admin, client)

	// Step 4: Construct the message to be signed
	txMsg := utils.MoveMessageStruct{
		From:   admin.Address, // Sender's address
		Amount: 100,           // Transfer amount
		Nonce:  1,             // Transaction nonce
	}
	msgBytes := txMsg.Serialize() // Serialize the message into bytes

	// Step 6: Create a hash of the message
	msgHash := sha256.Sum256(msgBytes)

	// Step 5: Generate private key from the provided hex string
	adminPrivateKey := utils.HexToEd25519PrivateKey(yaml.Profiles.Default.PrivateKey)

	// Step 7: Sign the hashed message using the private key
	adminSignatures := utils.SignMessage(adminPrivateKey, msgHash[:])

	// Step 8: Verify the signature
	isValid := utils.VerifySignature(client, admin.Address, uint64(100), adminSignatures)

	if isValid {
		fmt.Println("Signatures Verified!!!")
	} else {
		fmt.Println("Signatures Not Verified!!!")
	}
	fmt.Println("Signature verification successful!")
}
