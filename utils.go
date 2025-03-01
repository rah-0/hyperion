package main

import (
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	Size2KB   = 2 * 1024
	Size4KB   = 4 * 1024
	Size8KB   = 8 * 1024
	Size16KB  = 16 * 1024
	Size32KB  = 32 * 1024
	Size64KB  = 64 * 1024
	Size128KB = 128 * 1024
	Size256KB = 256 * 1024
	Size512KB = 512 * 1024
	Size1MB   = 1 * 1024 * 1024
	Size10MB  = 10 * 1024 * 1024
	Size100MB = 100 * 1024 * 1024
	Size1GB   = 1 * 1024 * 1024 * 1024
)

func generateRandomMessage(s int) ([]byte, error) {
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

func generateRandomStringMessage(s int) (string, error) {
	bytes, err := generateRandomMessage(s)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func getEnvKeyValue(key string) string {
	value, _ := os.LookupEnv(key)
	return value
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func pathCreateDirs(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func fileCreate(path string, content []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

func fileRead(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func fileIsEditable(filename string) bool {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

func portIsInUse(port int) bool {
	address := fmt.Sprintf(":%d", port)
	conn, err := net.DialTimeout("tcp", address, time.Millisecond*100)
	if err != nil {
		return false // Port is likely not in use
	}
	if conn != nil {
		conn.Close()
		return true // Port is in use
	}
	return false
}
