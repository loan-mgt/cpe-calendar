package decrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"cpe/calendar/logger"
)

func DecryptMessage(encryptedBase64 string, privateKey *rsa.PrivateKey) (string, error) {
	// Log the decryption attempt with context
	logger.Log.Info().
		Str("encryptedBase64", encryptedBase64).
		Msg("Attempting to decrypt message")

	// Decode the Base64-encoded message
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		logger.Log.Error().
			Str("encryptedBase64", encryptedBase64).
			Err(err).
			Msg("Failed to decode base64 string")
		return "", fmt.Errorf("failed to decode base64 string: %v", err)
	}

	// Create a new SHA-256 hash for OAEP
	hash := sha256.New()

	// Decrypt the message using the private key
	decryptedBytes, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, encryptedBytes, nil)
	if err != nil {
		logger.Log.Error().
			Str("encryptedBase64", encryptedBase64).
			Err(err).
			Msg("Failed to decrypt message")
		return "", fmt.Errorf("failed to decrypt message: %v", err)
	}

	logger.Log.Info().
		Str("decryptedMessage", string(decryptedBytes)[:5]).
		Msg("Message decrypted successfully")

	// Return the decrypted message as a string
	return string(decryptedBytes), nil
}

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	pemFile := "secret/private.pem"

	// Log the private key loading attempt
	logger.Log.Info().
		Str("pemFile", pemFile).
		Msg("Loading private key from PEM file")

	// Read the private key file
	keyData, err := os.ReadFile(pemFile)
	if err != nil {
		logger.Log.Error().
			Str("pemFile", pemFile).
			Err(err).
			Msg("Failed to read private key file")
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(keyData)
	if block == nil {
		logger.Log.Fatal().
			Str("pemFile", pemFile).
			Msg("Failed to decode PEM block containing private key")
	}

	var privateKey interface{}
	// Parse the private key
	if block.Type == "PRIVATE KEY" {
		// PKCS#8 format
		privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			logger.Log.Fatal().
				Str("pemFile", pemFile).
				Err(err).
				Msg("Failed to parse PKCS#8 private key")
		}
	} else if block.Type == "RSA PRIVATE KEY" {
		// PKCS#1 format
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			logger.Log.Fatal().
				Str("pemFile", pemFile).
				Err(err).
				Msg("Failed to parse PKCS#1 private key")
		}
	} else {
		logger.Log.Fatal().
			Str("pemFile", pemFile).
			Msg("Unknown PEM block type")
	}

	logger.Log.Info().
		Str("pemFile", pemFile).
		Msg("Private key loaded successfully")

	return privateKey.(*rsa.PrivateKey), nil
}
