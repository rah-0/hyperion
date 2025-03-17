package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
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
			err = checkConfigs()
			if err != nil {
				return err
			}

			var wg sync.WaitGroup
			for _, n := range GlobalConfig.Nodes {
				wg.Add(1)

				go func(node *Node) {
					defer wg.Done()
					_ = BuildBinary(node.Host.Name)

					pathConfigForNode := filepath.Join(os.TempDir(), "hyperion_test_"+node.Host.Name+".config")
					_ = FileCOpy(p, pathConfigForNode)

					logFilePath := filepath.Join(os.TempDir(), "hyperion_test_"+node.Host.Name+".log")
					_ = FileDelete(logFilePath)
					logFile, _ := os.Create(logFilePath)
					defer logFile.Close()

					cmd := exec.Command(filepath.Join(os.TempDir(), "hyperion_test_"+node.Host.Name),
						"-pathConfig", pathConfigForNode,
						"-forceHost", node.Host.Name)
					cmd.Stdout = logFile
					cmd.Stderr = logFile

					if err := cmd.Start(); err != nil {
						fmt.Printf("Error running instance for host %s: %v\n", node.Host.Name, err)
					}
				}(n)
			}
			wg.Wait()
			return nil
		},
		UnloadResources: func() error {
			return Pkill("hyperion_test_*")
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
