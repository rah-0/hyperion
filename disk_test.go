package main

import (
	"testing"

	. "github.com/rah-0/hyperion/register"
	. "github.com/rah-0/hyperion/util"
)

func TestWriteReadSingle(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	defer FileDelete(d.Path) // Cleanup temp file

	if len(Entities) == 0 {
		t.Fatal("no entities generated")
	}
	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		instanceToSave := e.New()
		instanceToSave.SetFieldValue("Name", "John")
		instanceToSave.SetFieldValue("Surname", "Doe")
		if err := d.WriteToFile(instanceToSave); err != nil {
			t.Fatalf("WriteToFile failed: %v", err)
		}

		instanceToSave.BufferReset()

		instanceToLoad := e.New()
		if err := d.ReadFromFile(instanceToLoad); err != nil {
			t.Fatalf("ReadFromFile failed: %v", err)
		}
		if instanceToLoad.GetFieldValue("Name") != "John" {
			t.Fatalf("ReadFromFile failed: %v != John", instanceToLoad.GetFieldValue("Name"))
		}
		if instanceToLoad.GetFieldValue("Surname") != "Doe" {
			t.Fatalf("ReadFromFile failed: %v != Doe", instanceToLoad.GetFieldValue("Surname"))
		}
	}
}
