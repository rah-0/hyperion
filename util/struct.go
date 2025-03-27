package util

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"unicode"
)

// StructsExtractFromPackage scans a package directory and extracts struct definitions
// includePrivate determines whether to include unexported (private) fields
// limit specifies the maximum number of structs to extract (0 means no limit)
func StructsExtractFromPackage(pkgDir string, includePrivate bool, limit int) ([]StructDef, error) {
	var structs []StructDef

	files, err := os.ReadDir(pkgDir)
	if err != nil {
		return nil, err
	}

	fs := token.NewFileSet()

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".go" || file.IsDir() || filepath.Base(file.Name()) == "test" {
			continue
		}

		srcPath := filepath.Join(pkgDir, file.Name())
		src, err := os.ReadFile(srcPath)
		if err != nil {
			return nil, err
		}

		node, err := parser.ParseFile(fs, srcPath, src, parser.AllErrors)
		if err != nil {
			return nil, err
		}

		done := false
		ast.Inspect(node, func(n ast.Node) bool {
			if done {
				return false
			}

			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return true
			}

			var fields []StructField
			for _, field := range st.Fields.List {
				if len(field.Names) == 0 {
					err = fmt.Errorf("embedded fields are not supported in struct %s", ts.Name.Name)
					return false
				}

				fieldType := exprToString(field.Type)
				tag := ""
				if field.Tag != nil {
					tag = field.Tag.Value
				}

				for _, name := range field.Names {
					if !includePrivate && unicode.IsLower(rune(name.Name[0])) {
						continue
					}

					fields = append(fields, StructField{
						Name: name.Name,
						Type: fieldType,
						Tag:  tag,
					})
				}
			}

			structs = append(structs, StructDef{
				Name:   ts.Name.Name,
				Fields: fields,
			})

			if limit > 0 && len(structs) >= limit {
				// Stop traversal
				done = true
				return false
			}

			return true
		})
		if err != nil {
			return nil, err
		}

		if done || (limit > 0 && len(structs) >= limit) {
			break
		}
	}

	return structs, nil
}

// exprToString converts AST expression to a readable string representation
func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), expr)
	return buf.String()
}
