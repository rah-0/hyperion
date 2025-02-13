package main

import (
	"os"
	"path/filepath"
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
