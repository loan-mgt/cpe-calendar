package decrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func DecryptMessage(encryptedBase64 string, privateKey *rsa.PrivateKey) (string, error) {
	// Decode the Base64-encoded message
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 string: %v", err)
	}

	// Create a new SHA-256 hash for OAEP
	hash := sha256.New()

	// Decrypt the message using the private key
	decryptedBytes, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, encryptedBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt message: %v", err)
	}

	// Return the decrypted message as a string
	return string(decryptedBytes), nil
}

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	pemFile := "secret/private.pem"

	// Read the private key file
	keyData, err := os.ReadFile(pemFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(keyData)
	if block == nil {
		log.Fatalf("Failed to decode PEM block containing private key")
	}

	var privateKey interface{}
	// Parse the private key
	if block.Type == "PRIVATE KEY" {
		// PKCS#8 format
		privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			log.Fatalf("Failed to parse PKCS#8 private key: %v", err)
		}
	} else if block.Type == "RSA PRIVATE KEY" {
		// PKCS#1 format
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Fatalf("Failed to parse PKCS#1 private key: %v", err)
		}
	} else {
		log.Fatalf("Unknown PEM block type")
	}

	return privateKey.(*rsa.PrivateKey), nil
}
