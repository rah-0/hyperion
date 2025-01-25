package testutil

import (
	"crypto/rand"
	"errors"
)

var (
	sizes = map[int]int{
		1:  2 * 1024,               // 2KB
		2:  4 * 1024,               // 4KB
		3:  8 * 1024,               // 8KB
		4:  16 * 1024,              // 16KB
		5:  32 * 1024,              // 32KB
		6:  64 * 1024,              // 64KB
		7:  128 * 1024,             // 128KB
		8:  256 * 1024,             // 256KB
		9:  512 * 1024,             // 512KB
		10: 1 * 1024 * 1024,        // 1MB
		11: 10 * 1024 * 1024,       // 10MB
		12: 100 * 1024 * 1024,      // 100MB
		13: 1 * 1024 * 1024 * 1024, // 1GB
	}
)

func GenerateMessage(sizeType int) ([]byte, error) {
	// Get the target size
	targetSize, ok := sizes[sizeType]
	if !ok {
		return []byte{}, errors.New("invalid size type: must be between 1 and 6")
	}

	// Generate random content
	bytes := make([]byte, targetSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return []byte{}, err
	}

	// Convert random bytes to a string
	for i, b := range bytes {
		// Limit to printable ASCII characters (32 to 126)
		bytes[i] = 32 + (b % 95)
	}

	return bytes, nil
}
