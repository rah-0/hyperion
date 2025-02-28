package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

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
