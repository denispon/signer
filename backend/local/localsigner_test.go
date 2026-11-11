package local

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/denispon/signer/sdk/signers"
)

func TestLocalSignerRSA(t *testing.T) {
	// Change to project root for key files
	projectRoot := getProjectRoot(t)
	os.Chdir(projectRoot)
	
	signer := &LocalSigner{}

	request := &signers.SignRequest{
		KeyID:     "rsa-key",
		Algorithm: signers.AlgorithmRSAPSSSHA256,
		Payload:   []byte("Hello, World!"),
	}

	result := signer.Sign(request)

	if result.Status != signers.STATUS_OK {
		t.Errorf("Expected STATUS_OK, got %d", result.Status)
	}

	if len(result.Signature) == 0 {
		t.Error("Expected non-empty signature")
	}

	if result.Algorithm != signers.AlgorithmRSAPSSSHA256 {
		t.Errorf("Expected algorithm %s, got %s", signers.AlgorithmRSAPSSSHA256, result.Algorithm)
	}
}

func TestLocalSignerECDSA(t *testing.T) {
	// Change to project root for key files
	projectRoot := getProjectRoot(t)
	os.Chdir(projectRoot)
	
	signer := &LocalSigner{}

	request := &signers.SignRequest{
		KeyID:     "ecdsa-key",
		Algorithm: signers.AlgorithmECDSAP256SHA256,
		Payload:   []byte("Hello, World!"),
	}

	result := signer.Sign(request)

	if result.Status != signers.STATUS_OK {
		t.Errorf("Expected STATUS_OK, got %d", result.Status)
	}

	if len(result.Signature) == 0 {
		t.Error("Expected non-empty signature")
	}

	if result.Algorithm != signers.AlgorithmECDSAP256SHA256 {
		t.Errorf("Expected algorithm %s, got %s", signers.AlgorithmECDSAP256SHA256, result.Algorithm)
	}
}

func TestLocalSignerInvalidAlgorithm(t *testing.T) {
	projectRoot := getProjectRoot(t)
	os.Chdir(projectRoot)
	
	signer := &LocalSigner{}

	request := &signers.SignRequest{
		KeyID:     "test-key",
		Algorithm: signers.AlgorithmUnknown,
		Payload:   []byte("test"),
	}

	result := signer.Sign(request)

	if result.Status != signers.STATUS_FAILED {
		t.Errorf("Expected STATUS_FAILED, got %d", result.Status)
	}
}

func TestLocalSignerEmptyKeyID(t *testing.T) {
	projectRoot := getProjectRoot(t)
	os.Chdir(projectRoot)
	
	signer := &LocalSigner{}

	request := &signers.SignRequest{
		KeyID:     "",
		Algorithm: signers.AlgorithmRSAPSSSHA256,
		Payload:   []byte("test"),
	}

	result := signer.Sign(request)

	if result.Status != signers.STATUS_FAILED {
		t.Errorf("Expected STATUS_FAILED for empty KeyID, got %d", result.Status)
	}
}

// Helper to find project root
func getProjectRoot(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	
	// Walk up until we find go.mod
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			t.Fatalf("Could not find project root")
		}
		wd = parent
	}
}
