package util

import (
	"os"
	"path/filepath"
)

func FileCreate(path string, content []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

func FileRead(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func FileIsEditable(filename string) bool {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return false
	}
	file.Close()
	return true
}
