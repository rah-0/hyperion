package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestExtractStructsFromPackage(t *testing.T) {
	// Generate a unique temporary directory
	tempDir := filepath.Join(os.TempDir(), "testpkg_"+uuid.New().String())
	err := DirectoryCreate(tempDir)
	if err != nil {
		t.Fatalf("Failed to create temp test directory: %v", err)
	}
	defer DirectoryRemove(tempDir) // Cleanup after test

	// Create test Go files
	goFileContent := `package testpkg

type User struct {
	ID    int
	Name  string
	Email string
}

type Product struct {
	Name  string
	Price float64
}`

	err = FileCreate(filepath.Join(tempDir, "models.go"), []byte(goFileContent))
	if err != nil {
		t.Fatalf("Failed to create Go file: %v", err)
	}

	// Run the function to extract structs
	structs, err := StructsExtractFromPackage(tempDir, false, 0)
	if err != nil {
		t.Fatalf("Failed to extract structs: %v", err)
	}

	// Expected struct definitions
	expected := map[string][]StructField{
		"User": {
			{"ID", "int", "", false},
			{"Name", "string", "", false},
			{"Email", "string", "", false},
		},
		"Product": {
			{"Name", "string", "", false},
			{"Price", "float64", "", false},
		},
	}

	// Validate extracted structs
	if len(structs) != len(expected) {
		t.Fatalf("Expected %d structs, got %d", len(expected), len(structs))
	}

	for _, structDef := range structs {
		expFields, exists := expected[structDef.Name]
		if !exists {
			t.Errorf("Unexpected struct found: %s", structDef.Name)
			continue
		}
		if len(structDef.Fields) != len(expFields) {
			t.Errorf("Struct %s: expected %d fields, got %d", structDef.Name, len(expFields), len(structDef.Fields))
			continue
		}
		for i, field := range structDef.Fields {
			if field.Name != expFields[i].Name || field.Type != expFields[i].Type {
				t.Errorf("Struct %s: expected field %s %s, got %s %s",
					structDef.Name, expFields[i].Name, expFields[i].Type, field.Name, field.Type)
			}
		}
	}
}

func TestExtractStructsWithLimit(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "testpkg_"+uuid.New().String())
	err := DirectoryCreate(tempDir)
	if err != nil {
		t.Fatalf("Failed to create temp test directory: %v", err)
	}
	defer DirectoryRemove(tempDir)

	goFileContent := `package testpkg

type A struct {}
type B struct {}
type C struct {}
type D struct {}`

	err = FileCreate(filepath.Join(tempDir, "models.go"), []byte(goFileContent))
	if err != nil {
		t.Fatalf("Failed to create Go file: %v", err)
	}

	structs, err := StructsExtractFromPackage(tempDir, false, 2)
	if err != nil {
		t.Fatalf("Failed to extract structs: %v", err)
	}

	if len(structs) != 2 {
		t.Errorf("Expected 2 structs due to limit, got %d", len(structs))
	}
}

func TestExtractStructsIncludePrivate(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "testpkg_"+uuid.New().String())
	err := DirectoryCreate(tempDir)
	if err != nil {
		t.Fatalf("Failed to create temp test directory: %v", err)
	}
	defer DirectoryRemove(tempDir)

	goFileContent := `package testpkg

type Demo struct {
	PublicField  string
	privateField int
}`

	err = FileCreate(filepath.Join(tempDir, "model.go"), []byte(goFileContent))
	if err != nil {
		t.Fatalf("Failed to create Go file: %v", err)
	}

	structs, err := StructsExtractFromPackage(tempDir, true, 0)
	if err != nil {
		t.Fatalf("Failed to extract structs: %v", err)
	}

	if len(structs) != 1 {
		t.Fatalf("Expected 1 struct, got %d", len(structs))
	}

	fieldNames := map[string]bool{}
	for _, f := range structs[0].Fields {
		fieldNames[f.Name] = true
	}

	if !fieldNames["PublicField"] || !fieldNames["privateField"] {
		t.Errorf("Expected both exported and unexported fields, got: %v", fieldNames)
	}
}
