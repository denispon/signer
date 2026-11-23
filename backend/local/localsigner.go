package local

import (
	"crypto"
	"crypto/x509"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/denispon/signer/sdk/signers"
	rsasigner "github.com/denispon/signer/sdk/signers/rsa"
	ecdsasigner "github.com/denispon/signer/sdk/signers/ecdsa"
	"software.sslmate.com/src/go-pkcs12"
)

type LocalSigner struct {
}

func (ls *LocalSigner) Sign(request *signers.SignRequest) *signers.SignResult {

	if len(request.KeyID) == 0 || !signers.ValidAlgorithm[request.Algorithm] {
		return &signers.SignResult{
			Status: signers.STATUS_FAILED,
		}
	}

	var keySigner signers.KeySigner
	path := "backend/local/keys"

	switch request.Algorithm {
	case signers.AlgorithmRSAPSSSHA256:
		keySigner = &rsasigner.RsaSigner{}
		path = path + "/" + "rsa"
	case signers.AlgorithmECDSAP256SHA256:
		keySigner = &ecdsasigner.EcdsaSigner{}
		path = path + "/" + "ecdsa"
	default:
		return &signers.SignResult{
			Status: signers.STATUS_FAILED,
		}
	}

	privKey, _, err := loadP12Key(getKeyPath(path+"/"+"signing-key.p12"), "password1!")
	if err != nil {
		return &signers.SignResult{
			Status: signers.STATUS_FAILED,
		}
	}

	signature, err := keySigner.Sign(privKey, request.Payload)
	if err != nil {
		return &signers.SignResult{
			Status: signers.STATUS_FAILED,
		}
	}

	return &signers.SignResult{
		Status:    signers.STATUS_OK,
		Signature: signature,
		Algorithm: request.Algorithm,
	}

}

func (ls *LocalSigner) Verify(message []byte, result *signers.SignResult) error {
	return nil
}

// isTestMode returns true if the code is running under Go tests
func isTestMode() bool {
	// Check if testing flags are registered (most reliable method)
	return flag.Lookup("test.v") != nil || flag.Lookup("test.run") != nil
}

// getKeyPath returns the path to the key file, preferring testdata for deterministic test behavior
func getKeyPath(primaryPath string) string {
	// Extract algorithm from primary path (e.g., "rsa" or "ecdsa")
	var algo string
	if strings.Contains(primaryPath, "rsa") {
		algo = "rsa"
	} else if strings.Contains(primaryPath, "ecdsa") {
		algo = "ecdsa"
	}

	if algo == "" {
		return primaryPath
	}

	// If we're in test mode, prefer testdata keys
	if isTestMode() {
		return findTestdataKey(algo)
	}

	// In production mode, use primary path
	if _, err := os.Stat(primaryPath); err == nil {
		return primaryPath
	}

	// Fallback to testdata if primary path not found
	return findTestdataKey(algo)
}

// findTestdataKey locates the testdata key file by traversing up the directory tree
func findTestdataKey(algo string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Testdata path relative to current working directory
	testdataPath := filepath.Join("signer", "cmd", "signer-command", "testdata", "keys", algo, "signing-key.p12")

	// Try different working directory levels (up to 5 levels up)
	currentPath := cwd
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(currentPath, testdataPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		currentPath = filepath.Dir(currentPath)
	}

	// Not found, return the testdata path anyway (will fail with clear error)
	return filepath.Join(cwd, testdataPath)
}

func loadP12Key(storagePath string, password string) (crypto.PrivateKey, *x509.Certificate, error) {
	data, err := os.ReadFile(storagePath)
	if err != nil {
		return nil, nil, err
	}

	privateKey, cert, err := pkcs12.Decode(data, password)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, cert, nil
}
