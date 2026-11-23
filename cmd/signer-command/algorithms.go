package signercommand

import (
	"fmt"

	"github.com/denispon/signer/sdk/signers"
	"github.com/spf13/cobra"
)

var AlgorithmsCmd = &cobra.Command{
	Use:   "algorithms",
	Short: "List supported signing algorithms",
	Long:  `Display all supported cryptographic signing algorithms.`,
	RunE:  runAlgorithms,
}

func runAlgorithms(cmd *cobra.Command, args []string) error {
	algorithms := []struct {
		name        string
		description string
	}{
		{
			name:        string(signers.AlgorithmECDSAP256SHA256),
			description: "ECDSA with P-256 curve and SHA-256 hash",
		},
		{
			name:        string(signers.AlgorithmRSAPSSSHA256),
			description: "RSA-PSS with SHA-256 hash",
		},
		{
			name:        string(signers.AlgorithmEd25519),
			description: "Ed25519 elliptic curve signature scheme",
		},
	}

	fmt.Println("Supported Algorithms:")
	fmt.Println("=====================")
	for i, algo := range algorithms {
		fmt.Printf("%d. %s\n", i+1, algo.name)
		fmt.Printf("   %s\n\n", algo.description)
	}

	return nil
}
