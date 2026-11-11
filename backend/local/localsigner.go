package local

import (
	"crypto"
	"crypto/x509"
	"os"

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

	privKey, _, err := loadP12Key(path+"/"+"signing-key.p12", "password1!")
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
