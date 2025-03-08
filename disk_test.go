package main

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	. "github.com/rah-0/hyperion/register"
	. "github.com/rah-0/hyperion/util"
)

type TestDataSample struct {
	Name    string
	Surname string
}

func TestWriteRead(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()

	t.Cleanup(func() {
		size, _ := FileSizeHuman(d.Path)
		fmt.Println("Size used: " + size)
		fmt.Println("Memory used: " + GetMemoryUsage())
		FileDelete(d.Path)
	})

	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		writeEntity := func(name, surname string) Model {
			instance := e.New()
			instance.SetFieldValue("Name", name)
			instance.SetFieldValue("Surname", surname)
			if err := d.WriteToFile(instance); err != nil {
				t.Fatalf("WriteToFile failed for %s %s: %v", name, surname, err)
			}
			instance.BufferReset()
			return instance
		}

		expectedEntities := make([]TestDataSample, 100)
		for i := range expectedEntities {
			expectedEntities[i] = TestDataSample{
				Name:    uuid.NewString(),
				Surname: uuid.NewString(),
			}
		}

		for _, exp := range expectedEntities {
			writeEntity(exp.Name, exp.Surname)
		}

		entities, err := d.LoadAllEntities(e)
		if err != nil {
			t.Fatalf("Failed to load entities: %v", err)
		}

		if len(entities) != len(expectedEntities) {
			t.Fatalf("Expected %d entities, got %d", len(expectedEntities), len(entities))
		}

		for i, exp := range expectedEntities {
			if entities[i].GetFieldValue("Name") != exp.Name ||
				entities[i].GetFieldValue("Surname") != exp.Surname {
				t.Fatalf("Entity %d mismatch: expected (%s, %s), got (%s, %s)",
					i, exp.Name, exp.Surname,
					entities[i].GetFieldValue("Name"),
					entities[i].GetFieldValue("Surname"))
			}
		}
	}
}
