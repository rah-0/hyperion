package gen

import (
	"fmt"

	. "github.com/rah-0/hyperion/util"
)

func Sample() {
	pkgDir := "../entities" // Change this to your package path

	structs, err := ExtractStructsFromPackage(pkgDir)
	if err != nil {
		fmt.Println("Error extracting structs:", err)
		return
	}

	// Print extracted struct definitions
	for _, s := range structs {
		fmt.Printf("Struct: %s\n", s.Name)
		for _, f := range s.Fields {
			fmt.Printf("  - %s: %s\n", f.Name, f.Type)
		}
	}
}
