package cmd

import (
	"fmt"
	"os"

	"github.com/denispon/signer/cmd/signer-command"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "signer",
	Short: "A CLI tool for cryptographic signing operations",
	Long: `Signer is a command-line tool that provides cryptographic signing capabilities.
It supports multiple algorithms including RSA-PSS, ECDSA, and Ed25519.`,
	Version: "1.0.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(signercommand.SignCmd)
	rootCmd.AddCommand(signercommand.VerifyCmd)
	rootCmd.AddCommand(signercommand.AlgorithmsCmd)
}
