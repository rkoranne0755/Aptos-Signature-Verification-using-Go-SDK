package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
)

// MODULE_ADDRESS is the account address for the module.
var MODULE_ADDRESS = StringToAccountAddress("93fe88844dd49f95caf0e5b5647ff2300fb32731b896bfdf7d4a77a1d48e7295")

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
	bytes, err := hex.DecodeString(address)
	if err != nil {
		fmt.Println("Error decoding address", err.Error())
		return aptos.AccountAddress{}
	}
	return aptos.AccountAddress(bytes)
}