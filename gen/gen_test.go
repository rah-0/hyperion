package gen

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	err := Generate()
	if err != nil {
		t.Fatal(err)
	}
}
