package signercommand

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/denispon/signer/backend/local"
	"github.com/denispon/signer/sdk/signers"
	"github.com/spf13/cobra"
)

var (
	verifyAlgorithm string
	verifyKeyID     string
	verifyInput     string
	verifySignature string
)

var VerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a signature",
	Long: `Verify that a signature is valid for the given data.
The signature should be provided as base64-encoded string or file.`,
	RunE: runVerify,
}

func runVerify(cmd *cobra.Command, args []string) error {
	// Validate inputs
	if verifyAlgorithm == "" {
		return fmt.Errorf("--algorithm is required")
	}
	if verifyKeyID == "" {
		return fmt.Errorf("--key-id is required")
	}
	if verifyInput == "" {
		return fmt.Errorf("--input is required")
	}
	if verifySignature == "" {
		return fmt.Errorf("--signature is required")
	}

	// Validate algorithm
	algo := signers.Algorithm(verifyAlgorithm)
	if !signers.ValidAlgorithm[algo] {
		return fmt.Errorf("invalid algorithm: %s", verifyAlgorithm)
	}

	// Read input data
	var payload []byte
	var err error

	if _, err := os.Stat(verifyInput); err == nil {
		payload, err = os.ReadFile(verifyInput)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
	} else {
		payload = []byte(verifyInput)
	}

	// Read signature
	var sigBytes []byte

	if _, err := os.Stat(verifySignature); err == nil {
		sigBytes, err = os.ReadFile(verifySignature)
		if err != nil {
			return fmt.Errorf("failed to read signature file: %w", err)
		}
	} else {
		sigBytes = []byte(verifySignature)
	}

	// Decode base64 signature
	decodedSig, err := base64.StdEncoding.DecodeString(string(sigBytes))
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	// Create result with signature for verification
	result := &signers.SignResult{
		Status:    signers.STATUS_OK,
		Signature: decodedSig,
		Algorithm: algo,
		KeyID:     verifyKeyID,
	}

	// Verify
	signer := &local.LocalSigner{}
	err = signer.Verify(payload, result)
	if err != nil {
		fmt.Println("Verification failed: signature is invalid")
		return nil
	}

	fmt.Println("✓ Signature verification successful")
	fmt.Printf("Algorithm: %s\n", result.Algorithm)
	fmt.Printf("Key ID: %s\n", result.KeyID)

	return nil
}

func init() {
	VerifyCmd.Flags().StringVarP(&verifyAlgorithm, "algorithm", "a", "", "Signing algorithm (required)")
	VerifyCmd.Flags().StringVarP(&verifyKeyID, "key-id", "k", "", "Key identifier (required)")
	VerifyCmd.Flags().StringVarP(&verifyInput, "input", "i", "", "Input data or file path (required)")
	VerifyCmd.Flags().StringVarP(&verifySignature, "signature", "s", "", "Signature (base64) or file path (required)")

	VerifyCmd.MarkFlagRequired("algorithm")
	VerifyCmd.MarkFlagRequired("key-id")
	VerifyCmd.MarkFlagRequired("input")
	VerifyCmd.MarkFlagRequired("signature")
}
