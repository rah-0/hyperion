package main

import (
	"testing"

	"github.com/google/uuid"

	SampleV1 "github.com/rah-0/hyperion/entities/Sample/v1"

	. "github.com/rah-0/hyperion/util"
)

func TestWriteRead(t *testing.T) {
	d := NewDisk()
	d.WithNewSerializer()
	d.WithNewRandomPath()
	defer FileDelete(d.Path) // Cleanup temp file

	entity := SampleV1.New()
	entity.SetFieldValue("Name", uuid.NewString())

	decoded := SampleV1.Sample{}

	if err := d.WriteToFile(entity); err != nil {
		t.Fatalf("WriteToFile failed: %v", err)
	}
	if err := d.ReadFromFile(&decoded); err != nil {
		t.Fatalf("ReadFromFile failed: %v", err)
	}

	original := entity.GetFieldValue("Name")
	fromStorage := decoded.GetFieldValue("Name")
	if original != fromStorage {
		t.Fatalf("Expected %+v, got %+v", original, fromStorage)
	}
}
