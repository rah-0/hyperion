package util

import (
	"testing"
)

func TestNextVersion(t *testing.T) {
	tests := []struct {
		name           string
		currentVersion string
		expectedNext   string
		expectError    bool
	}{
		{"Simple increment", "v1", "v2", false},
		{"Two-digit version", "v42", "v43", false},
		{"Large number", "v999999", "v1000000", false},
		{"Zero version", "v0", "v1", false},
		{"No v prefix", "1", "v2", false},
		{"Empty string", "", "", true},
		{"Invalid format", "version1", "", true},
		{"Negative number", "v-1", "", true},
		{"Non-numeric", "vabc", "", true},
		{"Decimal version", "v1.2", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringNextVersion(tt.currentVersion)
			if (err != nil) != tt.expectError {
				t.Errorf("nextVersion() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if got != tt.expectedNext {
				t.Errorf("nextVersion() = %v, want %v", got, tt.expectedNext)
			}
		})
	}
}
