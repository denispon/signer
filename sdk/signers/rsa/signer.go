package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

type RsaSigner struct {
}

func (rsaSigner *RsaSigner) Sign(privateKey crypto.PrivateKey, payload []byte) ([]byte, error) {
	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("incorect key type")
	}
	hash := crypto.SHA256
	hashed := hash.New()
	hashed.Write(payload)
	digest := hashed.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, rsaKey, hash, digest, nil)

	return signature, err
}
