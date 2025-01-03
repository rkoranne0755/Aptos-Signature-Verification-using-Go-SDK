package utils

import (
	"github.com/aptos-labs/aptos-go-sdk"
)

// MODULE_ADDRESS is the account address for the module.
var MODULE_ADDRESS = StringToAccountAddress("Add Wallet Address Here.")

// MODULE_NAME is the name of the module.
const MODULE_NAME = "verify_signature"

// Argument represents the JSON data structure for an argument.
type Argument struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// JsonData represents the JSON structure for data passed to the module.
type JsonData struct {
	FunctionID string     `json:"function_id"`
	TypeArgs   []string   `json:"type_args"`
	Args       []Argument `json:"args"`
}

// NamedAddress represents a named account address.
type NamedAddress struct {
	Name    string
	Address string
}

// YAML represents the configuration structure in YAML format.
type YAML struct {
	Profiles struct {
		Default struct {
			Network    string `yaml:"network"`
			PrivateKey string `yaml:"private_key"`
			PublicKey  string `yaml:"public_key"`
			Account    string `yaml:"account"`
			RestURL    string `yaml:"rest_url"`
			FaucetURL  string `yaml:"faucet_url"`
		} `yaml:"default"`
		// Additional profiles can be uncommented and added here.
		// Acc1 struct {
		// 	Network    string `yaml:"network"`
		// 	PrivateKey string `yaml:"private_key"`
		// 	PublicKey  string `yaml:"public_key"`
		// 	Account    string `yaml:"account"`
		// 	RestURL    string `yaml:"rest_url"`
		// 	FaucetURL  string `yaml:"faucet_url"`
		// } `yaml:"acc1"`
		// Acc2 struct {
		// 	Network    string `yaml:"network"`
		// 	PrivateKey string `yaml:"private_key"`
		// 	PublicKey  string `yaml:"public_key"`
		// 	Account    string `yaml:"account"`
		// 	RestURL    string `yaml:"rest_url"`
		// 	FaucetURL  string `yaml:"faucet_url"`
		// } `yaml:"acc2"`
	} `yaml:"profiles"`
}

// StringToAccountAddress converts a hexadecimal string into an aptos.AccountAddress.
func StringToAccountAddress(address string) aptos.AccountAddress {

	accAdd := &aptos.AccountAddress{}
	err := accAdd.ParseStringRelaxed("0x197327f4b6da2bb794cc1f8136eeddd3851d7c65e78df013520e2505ff34c1f4")
	if err != nil {
		panic("Failed to parse address:" + err.Error())
	}

	// bytes, err := hex.DecodeString(address)
	// if err != nil {
	// 	fmt.Println("Error decoding address", err.Error())
	// 	return aptos.AccountAddress{}
	// }
	// return aptos.AccountAddress(bytes)

	return *accAdd
}
