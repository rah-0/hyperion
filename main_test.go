package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestCheckPathConfig_WithArgs(t *testing.T) {
	pathConfig = ""

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tempDir := os.TempDir()
	testFilePath := filepath.Join(tempDir, uuid.NewString(), "config.json")

	hostName, _ := os.Hostname()
	err := fileCreate(testFilePath, []byte(`{"Nodes":[{"Host":{"Name":"`+hostName+`"},"Path":{"Data":"/tmp/`+uuid.NewString()+`"}}]}`))
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	os.Args = []string{"cmd", "-pathConfig=" + testFilePath}
	main()

	os.RemoveAll(filepath.Dir(testFilePath))
}

func TestCheckPathConfig_WithEnv(t *testing.T) {
	pathConfig = ""

	tempDir := os.TempDir()
	testFilePath := filepath.Join(tempDir, uuid.NewString(), "config.json")

	hostName, _ := os.Hostname()
	err := fileCreate(testFilePath, []byte(`{"Nodes":[{"Host":{"Name":"`+hostName+`"},"Path":{"Data":"/tmp/`+uuid.NewString()+`"}}]}`))
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
