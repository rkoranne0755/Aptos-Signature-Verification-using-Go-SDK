# Aptos Signature Verification using Go SDK

This repository demonstrates how to perform signature verification for the Aptos blockchain.

## Getting Started

### Clone the Repository

Clone the repository using the following command:

```bash
git clone https://github.com/rkoranne0755/Aptos-Signature-Verification-using-Go-SDK.git
```

### Change Directory

Move into the project directory:

```bash
cd move/
```

### Initialize Aptos

Run the following command to initialize Aptos:

```bash
aptos init
```

1. Select `devnet` as the network.
2. Press `Enter` to confirm.

This will create a new wallet configuration file at `.aptos/config.yaml`.

### Configure Account Address

1. Locate your **Account Address** from the `.aptos/config.yaml` file.
2. Open the file `utils/config.go`.
3. Replace the placeholder with your **MODULE_ADDRESS**.

### Run the Program

After updating the wallet information, execute the program:

```bash
go run main.go
```

## Notes

- Make sure Go and the Aptos CLI are installed on your system.
- The `go.mod` file contains all required dependencies. Run `go mod tidy` if any dependencies are missing.

Enjoy using this repository for Aptos signature verification!
