package ecdsa

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
)

type EcdsaSigner struct {
}

func (ecdsaSigner *EcdsaSigner) Sign(privateKey crypto.PrivateKey, payload []byte) ([]byte, error) {
	rsaKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("incorect key type")
	}
	hash := crypto.SHA256
	hashed := hash.New()
	hashed.Write(payload)
	digest := hashed.Sum(nil)

	signature, err := ecdsa.SignASN1(rand.Reader, rsaKey, digest)

	return signature, err
}
