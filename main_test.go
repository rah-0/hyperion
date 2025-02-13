package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestCheckPathConfig_WithEnv(t *testing.T) {
	tempDir := os.TempDir()
	testFilePath := filepath.Join(tempDir, uuid.NewString(), "config.json")

	err := fileCreate(testFilePath, []byte(`{"config": "test"}`))
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	os.Setenv("HyperionPathConfig", testFilePath)
	defer os.Unsetenv("HyperionPathConfig")

	err = checkPathConfig()
	if err != nil {
		t.Fatalf("Expected checkPathConfig to succeed, got error: %v", err)
	}

	os.RemoveAll(filepath.Dir(testFilePath))
}
