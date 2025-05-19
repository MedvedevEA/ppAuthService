package secure

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"math/rand"
	"os"
	"strconv"
)

func GetHash(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
func CheckHash(text string, hash string) bool {
	textHash := GetHash(text)
	return textHash == hash
}

func LoadPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	puplicKeyByteArray, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	publicKeyPemBlock, _ := pem.Decode(puplicKeyByteArray)
	if publicKeyPemBlock == nil || publicKeyPemBlock.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("decoding error. The PEM block was not found or the type is not equal to RSA PUBLIC KEY")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyPemBlock.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("error converting the interface to the *rsa.PublicKey type")
	}
	return publicKey, nil
}

func GenerateSecureToken() string {
	return GetHash(strconv.FormatInt(rand.Int63(), 10))
}
