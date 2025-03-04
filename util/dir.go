package util

import (
	"os"
)

// DirectoryCreate creates a new directory with the given path
func DirectoryCreate(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// DirectoryRemove deletes the directory and its contents
func DirectoryRemove(path string) error {
	return os.RemoveAll(path)
}
