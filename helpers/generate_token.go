package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecureRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateSecureRandomToken(length int) (string, error) {
	bytes, err := GenerateSecureRandomBytes(length)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}