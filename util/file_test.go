package util

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestFileCreateAndRead(t *testing.T) {
	tempDir := os.TempDir()
	testFilePath := filepath.Join(tempDir, uuid.NewString(), "testfile.txt")

	expectedContent := []byte("Hello, World!")
	err := FileCreate(testFilePath, expectedContent)
	if err != nil {
		t.Fatalf("fileCreate failed: %v", err)
	}

	actualContent, err := FileRead(testFilePath)
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

	err := FileCreate(testFilePath, []byte{})
	if err != nil {
		t.Fatalf("fileCreate failed for empty content: %v", err)
	}

	content, err := FileRead(testFilePath)
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

	err := FileCreate(invalidPath, []byte("test data"))
	if err == nil {
		t.Fatalf("Expected error when creating file in an invalid path, but got nil")
	}
}

func TestFileExpanderWithOneExpander(t *testing.T) {
	fileContent := []byte(`
aeaerhhae
fwaefaerwe
gerrett
---Start Tag---
---End Tag---
fewr546543
f3453457587634
	`)

	filePath := "testfile1.txt"
	err := FileCreate(filePath, fileContent)
	if err != nil {
		t.Error(err)
	}
	defer FileDelete(filePath)

	FileExpand(filePath, []FileExpanderTags{
		{
			StartTag:   []byte("---Start Tag---"),
			EndTag:     []byte("---End Tag---"),
			ExpandWith: []byte("test-now"),
		},
	})

	fileContentRetrieved, err := FileRead(filePath)
	if err != nil {
		t.Error(err)
	}

	expectedFileContent := []byte(`
aeaerhhae
fwaefaerwe
gerrett
---Start Tag---
test-now
---End Tag---
fewr546543
f3453457587634
	`)

	if !bytes.Equal(fileContentRetrieved, expectedFileContent) {
		t.Error()
	}
}

func TestFileExpanderWithTwoExpanders(t *testing.T) {
	fileContent := []byte(`
aeaerhhae
fwaefaerwe
gerrett
---Start Tag---
---End Tag---
fewr546543
f3453457587634
---A Start---
---A End---
	`)

	filePath := "testfile2.txt"
	err := FileCreate(filePath, fileContent)
	if err != nil {
		t.Error(err)
	}
	defer FileDelete(filePath)

	FileExpand(filePath, []FileExpanderTags{
		{
			StartTag:   []byte("---Start Tag---"),
			EndTag:     []byte("---End Tag---"),
			ExpandWith: []byte("test-now"),
		},
		{
			StartTag:   []byte("---A Start---"),
			EndTag:     []byte("---A End---"),
			ExpandWith: []byte("Multi expansion"),
		},
	})

	fileContentRetrieved, err := FileRead(filePath)
	if err != nil {
		t.Error(err)
	}

	expectedFileContent := []byte(`
aeaerhhae
fwaefaerwe
gerrett
---Start Tag---
test-now
---End Tag---
fewr546543
f3453457587634
---A Start---
Multi expansion
---A End---
	`)

	if !bytes.Equal(fileContentRetrieved, expectedFileContent) {
		t.Error()
	}
}

func TestFileExpanderWithCount(t *testing.T) {
	fileContent := []byte(`
---Start Tag---
---End Tag---
---Start Tag---
---End Tag---
	`)

	filePath := "testfile_count.txt"
	err := FileCreate(filePath, fileContent)
	if err != nil {
		t.Error(err)
	}
	defer FileDelete(filePath)

	FileExpand(filePath, []FileExpanderTags{
		{
			StartTag:   []byte("---Start Tag---"),
			EndTag:     []byte("---End Tag---"),
			ExpandWith: []byte("Inserted Once"),
			Count:      1,
		},
	})

	fileContentRetrieved, err := FileRead(filePath)
	if err != nil {
		t.Error(err)
	}

	expectedFileContent := []byte(`
---Start Tag---
Inserted Once
---End Tag---
---Start Tag---
---End Tag---
	`)

	if !bytes.Equal(fileContentRetrieved, expectedFileContent) {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", expectedFileContent, fileContentRetrieved)
	}
}
