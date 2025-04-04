package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func BuildBinary(suffix string) error {
	binaryPath := filepath.Join(os.TempDir(), "hyperion_test_"+suffix)
	fmt.Println("Building binary at:", binaryPath)

	projectRoot, err := filepath.Abs(filepath.Join(".."))
	if err != nil {
		return fmt.Errorf("failed to resolve project root: %w", err)
	}

	cmd := exec.Command("go", "build", "-o", binaryPath, projectRoot)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build binary: %w", err)
	}

	if err := os.Chmod(binaryPath, 0755); err != nil {
		return fmt.Errorf("failed to set execute permissions: %w", err)
	}

	fmt.Println("Binary built and execution permissions granted:", binaryPath)
	return nil
}

func Pkill(processName string) error {
	cmd := exec.Command("pkill", processName)
	output, err := cmd.CombinedOutput()
	s := string(output)
	if err != nil {
		return fmt.Errorf("pkill failed: %v, output: %s", err, s)
	}
	if len(s) > 0 {
		fmt.Printf("pkill succeeded, output: %s\n", s)
	}
	return nil
}
