package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestPathExists(t *testing.T) {
	tempDir := os.TempDir()

	exists, _ := pathExists(tempDir)
	if !exists {
		t.Fatalf("Temp dir should exist")
	}

	exists, err := pathExists(filepath.Join(tempDir, uuid.NewString()))
	if exists {
		t.Fatalf("Random path should not exist")
	}
	if err != nil {
		t.Fatalf("Error should be nil")
	}
}

func TestFileCreateAndRead(t *testing.T) {
	tempDir := os.TempDir()
	testFilePath := filepath.Join(tempDir, uuid.NewString(), "testfile.txt")

	expectedContent := []byte("Hello, World!")
	err := fileCreate(testFilePath, expectedContent)
	if err != nil {
		t.Fatalf("fileCreate failed: %v", err)
	}

	actualContent, err := fileRead(testFilePath)
	if err != nil {
		t.Fatalf("fileRead failed: %v", err)
	}

	if string(actualContent) != string(expectedContent) {
		t.Fatalf("File content mismatch: got %s, expected %s", actualContent, expectedContent)
	}

	os.RemoveAll(filepath.Dir(testFilePath))
}

func TestFileCreateWithEmptyContent(t *testing.T) {
	tempDir := os.TempDir()
	testFilePath := filepath.Join(tempDir, uuid.NewString(), "emptyfile.txt")

	err := fileCreate(testFilePath, []byte{})
	if err != nil {
		t.Fatalf("fileCreate failed for empty content: %v", err)
	}

	content, err := fileRead(testFilePath)
	if err != nil {
		t.Fatalf("fileRead failed for empty content: %v", err)
	}

	if len(content) != 0 {
		t.Fatalf("File should be empty but got content: %s", content)
	}

	os.RemoveAll(filepath.Dir(testFilePath))
}

func TestFileCreateInvalidPath(t *testing.T) {
	invalidPath := "/invalid_path/testfile.txt"

	err := fileCreate(invalidPath, []byte("test data"))
	if err == nil {
		t.Fatalf("Expected error when creating file in an invalid path, but got nil")
	}
}

func TestPortAvailability(t *testing.T) {
	// Ports that are typically open on many systems
	commonOpenPorts := []int{22}

	// Ports over 60000 (likely to be unused)
	highPorts := []int{60001, 62000, 65000}

	// Test common open ports
	for _, port := range commonOpenPorts {
		if !portIsInUse(port) {
			t.Errorf("Unexpected: Port %d is not in use while it should be.", port)
		}
	}

	// Test high ports (likely to be unused)
	for _, port := range highPorts {
		if portIsInUse(port) {
			t.Errorf("Unexpected: High port %d is in use. This port was expected to be free.", port)
		}
	}
}
