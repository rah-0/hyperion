package util

import (
	"crypto/rand"
)

func GenerateRandomMessage(s int) ([]byte, error) {
	bytes := make([]byte, s)
	_, err := rand.Read(bytes)
	if err != nil {
		return []byte{}, err
	}
	for i, b := range bytes {
		// Limit to printable ASCII characters (32 to 126)
		bytes[i] = 32 + (b % 95)
	}

	return bytes, nil
}

func GenerateRandomStringMessage(s int) (string, error) {
	bytes, err := GenerateRandomMessage(s)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
