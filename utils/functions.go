package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
)

func compilePackage(packageDir, outputFile string, namedAddresses []NamedAddress) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)
	}
	fmt.Println("CWD:", cwd)

	fmt.Println("In order to run compilation, you must have the `aptos` CLI installed.")

	// Check if the `aptos` CLI is installed
	cmd := exec.Command("aptos", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("aptos is not installed. Please install it from the instructions on aptos.dev")
	}

	// Generate the named addresses argument
	var addressArgs []string
	for _, namedAddress := range namedAddresses {
		addressArgs = append(addressArgs, fmt.Sprintf("%s=%s", namedAddress.Name, namedAddress.Address))
	}
	addressArg := strings.Join(addressArgs, " ")

	// Construct the compile command
	compileCommand := fmt.Sprintf(
		"aptos move build-publish-payload --json-output-file %s --package-dir %s --named-addresses %s --assume-yes --skip-fetch-latest-git-deps --move-2",
		outputFile,
		packageDir,
		addressArg,
	)

	fmt.Println("Running the compilation locally, in a real situation you may want to compile this ahead of time.")
	fmt.Println("Command:", compileCommand)

	// Execute the compile command
	cmd = exec.Command("sh", "-c", compileCommand) // Use "sh -c" to execute the full command string
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute compile command: %v", err)
	}

	return nil
}

// Function to read and parse the JSON file
func getPackageBytesToPublish(filePath string) ([]byte, [][]byte, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current working directory: %v", err)
	}
	fmt.Println("CWD:", cwd)

	// Construct the path to the JSON file
	modulePath := filepath.Join(cwd, filePath)
	fmt.Println("ModulePath:", modulePath)

	// Read the JSON file
	fileContent, err := ioutil.ReadFile(modulePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Parse the JSON into the struct
	var jsonData JsonData
	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		return nil, nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Ensure there are at least two arguments in the array
	if len(jsonData.Args) < 2 {
		return nil, nil, fmt.Errorf("invalid JSON structure, expected at least 2 args")
	}

	// Extract metadataBytes as []byte
	metadataBytes, err := parseHexToBytes(jsonData.Args[0].Value)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse metadataBytes: %v", err)
	}

	// Extract byteCode as [][]byte
	byteCode, err := parseHexArrayToBytesArray(jsonData.Args[1].Value)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse byteCode: %v", err)
	}

	return metadataBytes, byteCode, nil
}

func VerifySignature(client *aptos.Client, from aptos.AccountAddress, amount uint64, signature []byte) bool {
	addressBytes, _ := bcs.Serialize(&from)
	amountBytes, _ := bcs.SerializeU64(100)
	sigatureBytes, _ := bcs.SerializeBytes(signature)

	viewResponce, err := client.View(&aptos.ViewPayload{
		Module: aptos.ModuleId{
			Address: MODULE_ADDRESS,
			Name:    MODULE_NAME},
		Function: "verify_signatures",
		Args: [][]byte{
			addressBytes,
			amountBytes,
			sigatureBytes,
		},
	})

	if err != nil {
		panic("Error Verifying Signatures: " + err.Error())
	}

	// Assuming viewResponce[0] is a boolean value
	return viewResponce[0].(bool)
}

func GetPublicKey(client *aptos.Client) any {
	viewResponce, err := client.View(&aptos.ViewPayload{
		Module: aptos.ModuleId{
			Address: MODULE_ADDRESS,
			Name:    MODULE_NAME},
		Function: "get_pub_key",
		Args:     [][]byte{},
	})

	if err != nil {
		panic("Error Verifying Signatures: " + err.Error())
	}

	// Assuming viewResponce[0] is a boolean value
	return viewResponce[0]
}

func UpdatePublicKey(client *aptos.Client, account *aptos.Account, publicKey crypto.PublicKey) {
	pubKeyBytes, _ := bcs.Serialize(publicKey)
	rawTx, err := client.BuildTransaction(account.Address, aptos.TransactionPayload{
		Payload: &aptos.EntryFunction{
			Module: aptos.ModuleId{
				Address: MODULE_ADDRESS,
				Name:    MODULE_NAME,
			},
			Function: "update_key",
			ArgTypes: []aptos.TypeTag{},
			Args: [][]byte{
				pubKeyBytes,
			},
		},
	})

	if err != nil {
		fmt.Println("Error Building Transaction: ", err)
	}

	// Add code to submit the transaction if needed
	// 3. Sign transaction
	signedTxn, err := rawTx.SignedTransaction(account)
	if err != nil {
		fmt.Println("Error Signing Transaction: ", err.Error())
	}

	// 4. Submit transaction
	submitResult, err := client.SubmitTransaction(signedTxn)
	if err != nil {
		fmt.Println("Error Submitting Transaction: ", err.Error())
	}
	txnHash := submitResult.Hash
	fmt.Println("transaction: ", txnHash)

	// 5. Wait for the transaction to complete
	committedTx, err := client.WaitForTransaction(txnHash)
	if err != nil {
		fmt.Println("Error Waiting for the Transaction: ", err.Error())
	}

	if !committedTx.Success {
		fmt.Println("Error Submitting  Transaction: ", committedTx.VmStatus)
	}

	fmt.Println("Transaction Successful: ", committedTx.Success)

}
