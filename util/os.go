package util

import (
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func GetEnvKeyValue(key string) string {
	value, _ := os.LookupEnv(key)
	return value
}

// GetModuleName reads a go.mod file from the specified path and returns the module name.
func GetModuleName(modPath string) (string, error) {
	// Resolve absolute path
	absPath, err := filepath.Abs(modPath)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", err
	}

	modFile, err := modfile.Parse(absPath, data, nil)
	if err != nil {
		return "", err
	}

	return modFile.Module.Mod.Path, nil
}
