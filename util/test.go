package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"testing"
)

// TestConfig defines arguments for the TestMain wrapper.
type TestConfig struct {
	M               *testing.M   // The testing.M instance from TestMain.
	LoadResources   func() error // Function to load necessary resources.
	UnloadResources func() error // Function to unload resources.
}

// TestMainWrapper is a wrapper for TestMain to handle resource loading/unloading.
func TestMainWrapper(c TestConfig) {
	// Load resources if a loader function is provided.
	if c.LoadResources != nil {
		if err := c.LoadResources(); err != nil {
			log.Fatalf("Failed to load resources: %v", err)
		}
	}

	// Run the tests.
	exitCode := c.M.Run()

	// Unload resources if an unloader function is provided.
	if c.UnloadResources != nil {
		if err := c.UnloadResources(); err != nil {
			log.Printf("Error unloading resources: %v", err)
		}
	}

	// Exit with the test run's exit code.
	os.Exit(exitCode)
}

// RunTestWithRecover executes a test function and recovers from panics, failing the test if a panic occurs.
func RunTestWithRecover(t *testing.T, testFunc func(*testing.T)) {
	defer RecoverTestHandler(t)
	testFunc(t)
}

// RecoverTestHandler recovers from a panic and marks the test as failed, printing the stack trace.
func RecoverTestHandler(t *testing.T) {
	if r := recover(); r != nil {
		t.Errorf("Test panicked: %v\nStack trace:\n%s", r, debug.Stack())
	}
}

// RunBenchWithRecover executes a bench function and recovers from panics, failing the test if a panic occurs.
func RunBenchWithRecover(b *testing.B, testFunc func(*testing.B)) {
	defer RecoverBenchHandler(b)
	testFunc(b)
}

// RecoverBenchHandler recovers from a panic and marks the bench as failed, printing the stack trace.
func RecoverBenchHandler(b *testing.B) {
	if r := recover(); r != nil {
		b.Errorf("Test panicked: %v\nStack trace:\n%s", r, debug.Stack())
	}
}

func BuildBinary() error {
	binaryPath := filepath.Join(os.TempDir(), "hyperion_test")
	fmt.Println("Building binary at:", binaryPath)

	// Build the Go binary
	cmd := exec.Command("go", "build", "-o", binaryPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build binary: %w", err)
	}

	// Give execution permissions to the binary
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
