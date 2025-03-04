package util

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// ExtractStructsFromPackage scans a package directory and extracts struct definitions
func ExtractStructsFromPackage(pkgDir string) ([]StructDef, error) {
	var structs []StructDef

	// Read all files in the directory
	files, err := os.ReadDir(pkgDir)
	if err != nil {
		return nil, err
	}

	// Create a new token file set
	fs := token.NewFileSet()

	// Loop through files in the package directory
	for _, file := range files {
		// Only parse .go files and ignore test files
		if filepath.Ext(file.Name()) != ".go" || file.IsDir() || filepath.Base(file.Name()) == "test" {
			continue
		}

		// Open and parse the Go file
		srcPath := filepath.Join(pkgDir, file.Name())
		src, err := os.ReadFile(srcPath)
		if err != nil {
			return nil, err
		}

		node, err := parser.ParseFile(fs, srcPath, src, parser.AllErrors)
		if err != nil {
			return nil, err
		}

		// Walk through AST to find struct definitions
		ast.Inspect(node, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			// Ensure it's a struct type
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return true
			}

			// Extract struct fields
			var fields []StructField
			for _, field := range st.Fields.List {
				// Skip embedded fields (anonymous fields)
				if len(field.Names) == 0 {
					continue
				}
				for _, name := range field.Names {
					fieldType := fmt.Sprintf("%s", field.Type)
					fields = append(fields, StructField{
						Name: name.Name,
						Type: fieldType,
					})
				}
			}

			// Store the struct definition
			structs = append(structs, StructDef{
				Name:   ts.Name.Name,
				Fields: fields,
			})

			return true
		})
	}

	return structs, nil
}
