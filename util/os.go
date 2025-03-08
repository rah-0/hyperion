package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

func GetMemoryUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	var unitIndex int
	memory := float64(m.Sys)

	for memory >= 1024 && unitIndex < len(units)-1 {
		memory /= 1024
		unitIndex++
	}

	return fmt.Sprintf("%.2f %s", memory, units[unitIndex])
}
