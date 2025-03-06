package util

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a temporary go.mod file
func createTempGoMod(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "go_mod_*.mod")
	if err != nil {
		return "", err
	}
	_, err = tmpFile.WriteString(content)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name()) // Ensure cleanup
		return "", err
	}
	tmpFile.Close()
	return tmpFile.Name(), nil
}

func TestGetModuleName(t *testing.T) {
	tests := []struct {
		name        string
		modContent  string
		expected    string
		expectError bool
	}{
		{
			name:       "Valid go.mod",
			modContent: "module github.com/example/project\n\ngo 1.20",
			expected:   "github.com/example/project",
		},
		{
			name:       "Valid go.mod with whitespace",
			modContent: "  module   github.com/example/whitespace  \n\n go 1.21",
			expected:   "github.com/example/whitespace",
		},
		{
			name:        "Invalid go.mod format",
			modContent:  "invalid content",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp go.mod file
			tempFile, err := createTempGoMod(tt.modContent)
			if err != nil {
				t.Fatalf("Failed to create temp go.mod: %v", err)
			}
			defer os.Remove(tempFile) // Ensure cleanup

			// Resolve absolute path of the temp file
			absPath, err := filepath.Abs(tempFile)
			if err != nil {
				t.Fatalf("Failed to resolve absolute path: %v", err)
			}

			// Run the function
			modName, err := GetModuleName(absPath)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if modName != tt.expected {
					t.Errorf("Expected module name %q, got %q", tt.expected, modName)
				}
			}
		})
	}
}
