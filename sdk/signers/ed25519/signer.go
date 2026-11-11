package ed25519

import (
	"crypto"
)

type Ed25519Signer struct {
}

func (dd25519Signer *Ed25519Signer) Sign(privateKey crypto.PrivateKey, payload []byte) ([]byte, error) {
	return []byte{}, nil
}
