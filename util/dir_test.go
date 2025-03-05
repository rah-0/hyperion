package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirectoryGetHighestVersion(t *testing.T) {
	basePath := filepath.Join(os.TempDir(), "test_version_dirs")
	defer os.RemoveAll(basePath)

	testDirs := []string{"v1", "v2", "v10", "v5", "v100", "v200"}
	for _, dir := range testDirs {
		err := DirectoryCreate(filepath.Join(basePath, dir))
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	highestDir, err := DirectoryGetHighestVersion(basePath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if highestDir != "v200" {
		t.Errorf("Expected highest directory 'v200', got '%s'", highestDir)
	}
}

func TestDirectoryGetHighestVersion_EmptyDir(t *testing.T) {
	basePath := filepath.Join(os.TempDir(), "test_empty_dir")
	defer os.RemoveAll(basePath)

	err := DirectoryCreate(basePath)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	highestDir, err := DirectoryGetHighestVersion(basePath)
	if err == nil {
		t.Errorf("Expected an error for empty directory, got nil")
	}

	if highestDir != "" {
		t.Errorf("Expected empty string for highest directory, got '%s'", highestDir)
	}
}

func TestDirectoryGetHighestVersion_NoVersionDirs(t *testing.T) {
	basePath := filepath.Join(os.TempDir(), "test_non_version_dirs")
	defer os.RemoveAll(basePath)

	// Create non-versioned directories
	testDirs := []string{"alpha", "beta", "gamma"}
	for _, dir := range testDirs {
		err := DirectoryCreate(filepath.Join(basePath, dir))
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	highestDir, err := DirectoryGetHighestVersion(basePath)
	if err == nil {
		t.Errorf("Expected an error for non-versioned directories, got nil")
	}

	if highestDir != "" {
		t.Errorf("Expected empty string for highest directory, got '%s'", highestDir)
	}
}
