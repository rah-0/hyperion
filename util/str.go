package util

import (
	"fmt"
	"strconv"
	"strings"
)

func StringNextVersion(currentVersion string) (string, error) {
	// Remove the 'v' prefix
	numStr := strings.TrimPrefix(currentVersion, "v")

	// Convert the remaining string to an integer
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return "", fmt.Errorf("invalid version format: %s", currentVersion)
	}

	// Validate that the number is non-negative
	if num < 0 {
		return "", fmt.Errorf("version number cannot be negative: %s", currentVersion)
	}

	// Increment the version number
	nextNum := num + 1

	// Return the new version string
	return fmt.Sprintf("v%d", nextNum), nil
}
