package main

import (
	"os"
	"testing"
)

func TestWriteRead(t *testing.T) {
	d := NewDisk()
	d.WithNewSerializer()
	d.WithNewRandomPath()

	defer os.Remove(d.Path) // Cleanup temp file

	type TestStruct struct {
		ID   int
		Name string
	}

	original := TestStruct{ID: 42, Name: "Hyperion"}
	var decoded TestStruct

	// Write to file
	if err := d.WriteToFile(original); err != nil {
		t.Fatalf("WriteToFile failed: %v", err)
	}

	// Read from file
	if err := d.ReadFromFile(&decoded); err != nil {
		t.Fatalf("ReadFromFile failed: %v", err)
	}

	// Ensure the decoded struct matches the original
	if original != decoded {
		t.Fatalf("Expected %+v, got %+v", original, decoded)
	}
}
