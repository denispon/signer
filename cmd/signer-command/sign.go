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
	signAlgorithm string
	signKeyID     string
	signInput     string
	signOutput    string
)

var SignCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign data using a private key",
	Long: `Sign data using a private key stored locally.
The input can be provided as a file path or directly as a string.`,
	RunE: runSign,
}

func runSign(cmd *cobra.Command, args []string) error {
	// Validate inputs
	if signAlgorithm == "" {
		return fmt.Errorf("--algorithm is required")
	}
	if signKeyID == "" {
		return fmt.Errorf("--key-id is required")
	}
	if signInput == "" {
		return fmt.Errorf("--input is required")
	}

	// Validate algorithm
	algo := signers.Algorithm(signAlgorithm)
	if !signers.ValidAlgorithm[algo] {
		return fmt.Errorf("invalid algorithm: %s", signAlgorithm)
	}

	// Read input data
	var payload []byte

	// Try reading as file first
	if _, err := os.Stat(signInput); err == nil {
		var readErr error
		payload, readErr = os.ReadFile(signInput)
		if readErr != nil {
			return fmt.Errorf("failed to read input file: %w", readErr)
		}
	} else {
		// Use as raw string/data
		payload = []byte(signInput)
	}

	// Create signer and sign
	signer := &local.LocalSigner{}
	result := signer.Sign(&signers.SignRequest{
		KeyID:     signKeyID,
		Payload:   payload,
		Algorithm: algo,
	})

	if result.Status != signers.STATUS_OK {
		return fmt.Errorf("signing failed")
	}

	// Output signature
	sigBase64 := base64.StdEncoding.EncodeToString(result.Signature)

	if signOutput != "" {
		// Write to file
		err := os.WriteFile(signOutput, []byte(sigBase64), 0644)
		if err != nil {
			return fmt.Errorf("failed to write signature to file: %w", err)
		}
		fmt.Printf("Signature written to: %s\n", signOutput)
	} else {
		// Print to stdout
		fmt.Printf("Signature (base64): %s\n", sigBase64)
	}

	fmt.Printf("Algorithm: %s\n", result.Algorithm)
	fmt.Printf("Key ID: %s\n", result.KeyID)

	return nil
}

func init() {
	SignCmd.Flags().StringVarP(&signAlgorithm, "algorithm", "a", "", "Signing algorithm (required)")
	SignCmd.Flags().StringVarP(&signKeyID, "key-id", "k", "", "Key identifier (required)")
	SignCmd.Flags().StringVarP(&signInput, "input", "i", "", "Input data or file path (required)")
	SignCmd.Flags().StringVarP(&signOutput, "output", "o", "", "Output file for signature (optional)")

	SignCmd.MarkFlagRequired("algorithm")
	SignCmd.MarkFlagRequired("key-id")
	SignCmd.MarkFlagRequired("input")
}
