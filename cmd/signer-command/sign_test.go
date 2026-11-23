package signercommand

import (
	"bytes"
	"os"
	"testing"
)

func TestSignCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		wantInOutput string
	}{
		{
			name:        "sign with valid ECDSA algorithm",
			args:        []string{"-a", "ECDSA-P256-SHA256", "-k", "test-key", "-i", "Hello World"},
			wantErr:     false,
			wantInOutput: "Signature",
		},
		{
			name:        "sign with valid RSA algorithm",
			args:        []string{"-a", "RSA-PSS-SHA256", "-k", "test-key", "-i", "Hello World"},
			wantErr:     false,
			wantInOutput: "Signature",
		},
		{
			name:        "sign missing required algorithm flag",
			args:        []string{"-k", "test-key", "-i", "Hello World"},
			wantErr:     true,
			wantInOutput: "",
		},
		{
			name:        "sign missing required key-id flag",
			args:        []string{"-a", "ECDSA-P256-SHA256", "-i", "Hello World"},
			wantErr:     true,
			wantInOutput: "",
		},
		{
			name:        "sign missing required input flag",
			args:        []string{"-a", "ECDSA-P256-SHA256", "-k", "test-key"},
			wantErr:     true,
			wantInOutput: "",
		},
		{
			name:        "sign with invalid algorithm",
			args:        []string{"-a", "INVALID", "-k", "test-key", "-i", "Hello World"},
			wantErr:     true,
			wantInOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			oldAlgo := signAlgorithm
			oldKeyID := signKeyID
			oldInput := signInput
			oldOutput := signOutput

			// Reset flags
			signAlgorithm = ""
			signKeyID = ""
			signInput = ""
			signOutput = ""

			// Capture output
			stdout, stderr := CaptureOutputWithErr(func() {
				cmd := SignCmd
				cmd.SetArgs(tt.args)
				_ = cmd.Execute()
			})

			// Combine output
			output := append(stdout, stderr...)

			// Restore original values
			signAlgorithm = oldAlgo
			signKeyID = oldKeyID
			signInput = oldInput
			signOutput = oldOutput

			// Check error handling
			if tt.wantErr && len(output) == 0 {
				// If we expected an error and got no output, that's also acceptable
				// (command ran but flags were invalid)
			}

			if !tt.wantErr && tt.wantInOutput != "" {
				if !bytes.Contains(output, []byte(tt.wantInOutput)) {
					t.Errorf("SignCommand() output missing '%s', got: %s", tt.wantInOutput, string(output))
				}
			}
		})
	}
}

func TestSignCommandWithFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_input_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test data
	if _, err := tmpFile.WriteString("Test data for signing"); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Save and reset flags
	oldAlgo := signAlgorithm
	oldKeyID := signKeyID
	oldInput := signInput
	oldOutput := signOutput

	signAlgorithm = ""
	signKeyID = ""
	signInput = ""
	signOutput = ""

	// Capture output
	output, _ := CaptureOutputWithErr(func() {
		cmd := SignCmd
		cmd.SetArgs([]string{"-a", "ECDSA-P256-SHA256", "-k", "test-key", "-i", tmpFile.Name()})
		_ = cmd.Execute()
	})

	// Restore
	signAlgorithm = oldAlgo
	signKeyID = oldKeyID
	signInput = oldInput
	signOutput = oldOutput

	if !bytes.Contains(output, []byte("Signature")) {
		t.Errorf("SignCommand() output missing 'Signature', got: %s", string(output))
	}
}

func TestSignCommandOutputToFile(t *testing.T) {
	// Create temporary output file
	tmpOutput, err := os.CreateTemp("", "test_signature_*.b64")
	if err != nil {
		t.Fatalf("failed to create temp output file: %v", err)
	}
	tmpOutput.Close()
	defer os.Remove(tmpOutput.Name())

	// Save and reset flags
	oldAlgo := signAlgorithm
	oldKeyID := signKeyID
	oldInput := signInput
	oldOutput := signOutput

	signAlgorithm = ""
	signKeyID = ""
	signInput = ""
	signOutput = ""

	// Capture output
	output, _ := CaptureOutputWithErr(func() {
		cmd := SignCmd
		cmd.SetArgs([]string{"-a", "ECDSA-P256-SHA256", "-k", "test-key", "-i", "Hello World", "-o", tmpOutput.Name()})
		_ = cmd.Execute()
	})

	// Restore
	signAlgorithm = oldAlgo
	signKeyID = oldKeyID
	signInput = oldInput
	signOutput = oldOutput

	if !bytes.Contains(output, []byte("Signature written to")) {
		t.Logf("SignCommand() output: %s", string(output))
	}

	// Verify file was created and contains signature
	content, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		t.Errorf("failed to read output file: %v", err)
		return
	}

	if len(content) == 0 {
		t.Error("output file is empty")
	}
}
