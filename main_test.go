package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"

	. "github.com/rah-0/hyperion/util"
)


func TestMain(m *testing.M) {
	TestMainWrapper(TestConfig{
		M: m,
		LoadResources: func() error {
			p, err := filepath.Abs("./config.json")
			if err != nil {
				return err
			}

			pathConfig = p
			err = run()
			if err != nil {
				return err
			}

			return nil
		},
		UnloadResources: func() error {
			return nil
		},
	})
}

func TestCheckPathConfig_WithArgs(t *testing.T) {
	t.Skip()

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tempDir := os.TempDir()
	testFilePath := filepath.Join(tempDir, uuid.NewString(), "config.json")

	hostName, _ := os.Hostname()
	err := FileCreate(testFilePath, []byte(`{"Nodes":[{"Host":{"Name":"`+hostName+`"},"Path":{"Data":"/tmp/`+uuid.NewString()+`"}}]}`))
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	os.Args = []string{"cmd", "-pathConfig=" + testFilePath}
	main()

	os.RemoveAll(filepath.Dir(testFilePath))
}

func TestCheckPathConfig_WithEnv(t *testing.T) {
	t.Skip()

	originalValue := os.Getenv("HyperionPathConfig")
	defer func() {
		if originalValue != "" {
			os.Setenv("HyperionPathConfig", originalValue)
		} else {
			os.Unsetenv("HyperionPathConfig")
		}
	}()

	tempDir := os.TempDir()
	testFilePath := filepath.Join(tempDir, uuid.NewString(), "config.json")

	hostName, _ := os.Hostname()
	err := FileCreate(testFilePath, []byte(`{"Nodes":[{"Host":{"Name":"`+hostName+`"},"Path":{"Data":"/tmp/`+uuid.NewString()+`"}}]}`))
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	os.Setenv("HyperionPathConfig", testFilePath)

	err = checkPathConfig()
	if err != nil {
		t.Fatalf("Expected checkPathConfig to succeed, got error: %v", err)
	}

	os.RemoveAll(filepath.Dir(testFilePath))
}
