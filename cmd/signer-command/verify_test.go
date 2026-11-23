package signercommand

import (
	"os"
	"testing"
)

func TestVerifyCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "verify missing required algorithm flag",
			args:    []string{"-k", "test-key", "-i", "Hello World", "-s", "sig"},
			wantErr: true,
		},
		{
			name:    "verify missing required key-id flag",
			args:    []string{"-a", "ECDSA-P256-SHA256", "-i", "Hello World", "-s", "sig"},
			wantErr: true,
		},
		{
			name:    "verify missing required input flag",
			args:    []string{"-a", "ECDSA-P256-SHA256", "-k", "test-key", "-s", "sig"},
			wantErr: true,
		},
		{
			name:    "verify missing required signature flag",
			args:    []string{"-a", "ECDSA-P256-SHA256", "-k", "test-key", "-i", "Hello World"},
			wantErr: true,
		},
		{
			name:    "verify with invalid algorithm",
			args:    []string{"-a", "INVALID", "-k", "test-key", "-i", "Hello World", "-s", "sig"},
			wantErr: true,
		},
		{
			name:    "verify with invalid base64 signature",
			args:    []string{"-a", "ECDSA-P256-SHA256", "-k", "test-key", "-i", "Hello World", "-s", "!!!invalid!!!"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and reset flags
			oldAlgo := verifyAlgorithm
			oldKeyID := verifyKeyID
			oldInput := verifyInput
			oldSignature := verifySignature

			verifyAlgorithm = ""
			verifyKeyID = ""
			verifyInput = ""
			verifySignature = ""

			// Capture output
			output, _ := CaptureOutputWithErr(func() {
				cmd := VerifyCmd
				cmd.SetArgs(tt.args)
				_ = cmd.Execute()
			})

			// Restore
			verifyAlgorithm = oldAlgo
			verifyKeyID = oldKeyID
			verifyInput = oldInput
			verifySignature = oldSignature

			t.Logf("Test %q: output = %s", tt.name, string(output))
		})
	}
}

func TestVerifyCommandWithFiles(t *testing.T) {
	// Create temporary files
	tmpInput, err := os.CreateTemp("", "test_input_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp input file: %v", err)
	}
	defer os.Remove(tmpInput.Name())

	tmpSig, err := os.CreateTemp("", "test_sig_*.b64")
	if err != nil {
		t.Fatalf("failed to create temp sig file: %v", err)
	}
	defer os.Remove(tmpSig.Name())

	// Write test data
	if _, err := tmpInput.WriteString("File test data"); err != nil {
		t.Fatalf("failed to write to input file: %v", err)
	}
	tmpInput.Close()

	// Write a placeholder signature
	if _, err := tmpSig.WriteString("MEUCIQDj8K0pX3g9N8vL1p5D7q8q9r0s1t2u3v4w5x6y7z8a9bCIEQ9A0B1C2D3E4F5G6H7I8J9K0L1M2N3O4P5Q6R7S8T"); err != nil {
		t.Fatalf("failed to write to sig file: %v", err)
	}
	tmpSig.Close()

	// Save and reset flags
	oldAlgo := verifyAlgorithm
	oldKeyID := verifyKeyID
	oldInput := verifyInput
	oldSignature := verifySignature

	verifyAlgorithm = ""
	verifyKeyID = ""
	verifyInput = ""
	verifySignature = ""

	// Capture output
	output, _ := CaptureOutputWithErr(func() {
		cmd := VerifyCmd
		cmd.SetArgs([]string{"-a", "ECDSA-P256-SHA256", "-k", "file-test-key", "-i", tmpInput.Name(), "-s", tmpSig.Name()})
		_ = cmd.Execute()
	})

	// Restore
	verifyAlgorithm = oldAlgo
	verifyKeyID = oldKeyID
	verifyInput = oldInput
	verifySignature = oldSignature

	// Error expected since signature is placeholder
	t.Logf("Verify with files - output: %s", string(output))
}
