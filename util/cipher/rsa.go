package cipher

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"

	"github.com/juju/errors"
)

// EncryptRSA ...
func EncryptRSA(rsaPubKey *rsa.PublicKey, data string) (string, error) {
	if rsaPubKey == nil || data == "" {
		return "", errors.Annotate(errors.New("Invalid parameter"), "cipher EncryptRSA")
	}

	cipherData, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPubKey, []byte(data))
	if err != nil {
		return "", errors.Annotate(err, "cipher EncryptRSA")
	}

	encData := base64.StdEncoding.EncodeToString(cipherData)

	return string(encData), nil
}

// DecryptRSA ...
func DecryptRSA(rsaPriKey *rsa.PrivateKey, data string) (string, error) {
	if rsaPriKey == nil || data == "" {
		return "", errors.Annotate(errors.New("Invalid parameter"), "cipher DecryptRSA")
	}

	encData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.Annotate(err, "cipher DecryptRSA")
	}

	plainData, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPriKey, encData)
	if err != nil {
		return "", errors.Annotate(err, "cipher DecryptRSA")
	}

	return string(plainData), nil
}
