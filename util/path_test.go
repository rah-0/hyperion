package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestPathExists(t *testing.T) {
	tempDir := os.TempDir()

	exists, _ := PathExists(tempDir)
	if !exists {
		t.Fatalf("Temp dir should exist")
	}

	exists, err := PathExists(filepath.Join(tempDir, uuid.NewString()))
	if exists {
		t.Fatalf("Random path should not exist")
	}
	if err != nil {
		t.Fatalf("Error should be nil")
	}
}
