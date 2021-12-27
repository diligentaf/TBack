package cipher

// https://8gwifi.org/docs/go-aes.jsp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/juju/errors"
)

// AESTestGcm ...
func AESTestGcm() {

	// Must Kept Secret No Hardcoding , This is for Demo purpose.
	//      12345678901234567890123456789012
	key := "youwastedtimereadingthissentence"
	plainText := "0987654321"
	fmt.Printf("Original Text:  %s\n", plainText)
	fmt.Println()

	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	fmt.Println("====GCM Encryption/ Decryption Without AAD====")

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	iv := []byte("itssensitive")
	ciphertext, err := AESEncryptGCM(key, plainText, iv, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("GCM Encrypted Text:  %s\n", ciphertext)
	ret, err := AESDecryptGCM(key, ciphertext, iv, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("GCM Decrypted Text:  %s\n", ret)
	fmt.Println()

	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	fmt.Println("====GCM Encryption/ Decryption Using AAD====")

	// Never Use Same IV or Nonce
	iv = make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err.Error())
	}
	additionalData := "Not Secret AAD Value"
	ciphertext, err = AESEncryptGCM(key, plainText, iv, []byte(additionalData))
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("GCM Encrypted Text:  %s\n", ciphertext)
	ret, err = AESDecryptGCM(key, ciphertext, iv, []byte(additionalData))
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("GCM Decrypted Text:  %s\n", ret)

	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	fmt.Println("====GCM Encryption/ Decryption Without IV ====")
	ciphertext, err = AESEncryptGCM(key, plainText, nil, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("GCM Encrypted Text:  %s\n", ciphertext)
	ret, err = AESDecryptGCM(key, ciphertext, nil, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("GCM Decrypted Text:  %s\n", ret)
	fmt.Println()
}

// AESEncryptGCM ...
func AESEncryptGCM(k string, plaintext string, iv []byte, additionalData []byte) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	key := k // key, _ := base64.StdEncoding.DecodeString(k) // we do not need base64-decoding

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//////////
	if iv == nil {
		iv = make([]byte, 12)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return "", err
		}

		ciphertext := aesgcm.Seal(nil, iv, []byte(plaintext), additionalData)
		padded := make([]byte, len(ciphertext)+len(iv))
		copy(padded, ciphertext)
		copy(padded[len(ciphertext):], iv)
		return base64.StdEncoding.EncodeToString(padded), nil
	}

	//////////
	ciphertext := aesgcm.Seal(nil, iv, []byte(plaintext), additionalData)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecryptGCM ...
func AESDecryptGCM(k string, ct string, iv []byte, additionalData []byte) (string, error) {
	if ct == "" {
		return "", nil
	}

	key := k // key, _ := base64.StdEncoding.DecodeString(k) // we do not need base64-decoding
	ciphertext, _ := base64.StdEncoding.DecodeString(ct)

	//////////
	if iv == nil {
		iv = ciphertext[len(ciphertext)-12:]
		ciphertext = ciphertext[0 : len(ciphertext)-12]
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := aesgcm.Open(nil, iv, ciphertext, additionalData)
	if err != nil {
		return "", err
	}
	s := string(plaintext[:])
	return s, nil
}

// GetKeyWithIvByRSA ...
func GetKeyWithIvByRSA(rsaPriKey *rsa.PrivateKey, encKey string) ([]string, error) {
	if rsaPriKey == nil {
		return nil, errors.Annotate(errors.New("No Private Key"), "Cipher GetKeyWithIvByRSA")
	}

	keySet := make([]string, 2, 2)

	if encKey == "" {
		return keySet, nil
	}

	aesKey, err := DecryptRSA(rsaPriKey, encKey)
	if err != nil {
		return nil, errors.Annotate(err, "Cipher GetKeyWithIvByRSA")
	}

	keySet[0] = aesKey[0:32]
	keySet[1] = aesKey[32:len(aesKey)]

	return keySet, nil
}
