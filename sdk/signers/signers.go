package signers

import "crypto"

type Signer interface {
	Sign(SignRequest) SignResult
	Verify(message []byte, result SignResult) error
}

type SignRequest struct {
	KeyID     string
	Payload   []byte
	Algorithm Algorithm
}

type SignResult struct {
	Status    Status
	Signature []byte
	Algorithm Algorithm
	KeyID     string
}

type Algorithm string

const (
	AlgorithmUnknown         Algorithm = "UNKNOWN"
	AlgorithmECDSAP256SHA256 Algorithm = "ECDSA-P256-SHA256"
	AlgorithmRSAPSSSHA256    Algorithm = "RSA-PSS-SHA256"
	AlgorithmEd25519         Algorithm = "Ed25519"
)

var ValidAlgorithm = map[Algorithm]bool{
	AlgorithmECDSAP256SHA256: true,
	AlgorithmEd25519:         true,
	AlgorithmRSAPSSSHA256:    true,
}

type Status int8

const (
	STATUS_OK     = iota //0
	STATUS_FAILED        //1
)

type KeySigner interface {
	Sign(privateKey crypto.PrivateKey, payload []byte) ([]byte, error)
}
