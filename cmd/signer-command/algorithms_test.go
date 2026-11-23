package signercommand

import (
	"bytes"
	"testing"
)

func TestAlgorithmsCommand(t *testing.T) {
	// Capture output from fmt.Println
	output := CaptureOutput(func() {
		cmd := AlgorithmsCmd
		cmd.SetArgs([]string{})
		_ = cmd.Execute()
	})

	// Check that all required algorithms are listed
	requiredAlgos := []string{
		"ECDSA-P256-SHA256",
		"RSA-PSS-SHA256",
		"Ed25519",
		"Supported Algorithms",
	}

	for _, algo := range requiredAlgos {
		if !bytes.Contains(output, []byte(algo)) {
			t.Errorf("AlgorithmsCommand() output missing '%s'", algo)
		}
	}
}
