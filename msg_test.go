package main

import (
	"fmt"
	"path/filepath"
	"testing"

	SampleV1 "github.com/rah-0/hyperion/entities/Sample/v1"
	"github.com/rah-0/hyperion/register"
	"github.com/rah-0/hyperion/util"
)

func TestMessageInsert(t *testing.T) {
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		t.Fatal(err)
	}

	entity := SampleV1.Sample{
		Name:    "Something",
		Surname: "Else",
	}

	_, err = entity.DbInsert(c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMessageInsert1000(t *testing.T) {
	c, err := ConnectToNode(GlobalNode)
	defer util.FileDelete(filepath.Join(GlobalNode.Path.Data, "SampleV1.bin"))
	if err != nil {
		t.Fatal(err)
	}

	var expected []SampleV1.Sample
	for i := 0; i < 1000; i++ {
		entity := SampleV1.Sample{
			Name:    fmt.Sprintf("Something%d", i),
			Surname: fmt.Sprintf("Else%d", i),
		}
		if _, err := entity.DbInsert(c); err != nil {
			t.Fatal(err)
		}
		expected = append(expected, entity)
	}

	d := NewDisk()
	d.WithPath(filepath.Join(GlobalNode.Path.Data, "SampleV1.bin"))
	if err = d.OpenFile(); err != nil {
		t.Fatal(err)
	}
	for _, e := range register.Entities {
		if e.Name != "Sample" {
			continue
		}
		d.WithEntity(e)
	}

	entities, err := d.DataReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(entities) != len(expected) {
		t.Fatalf("Expected %d entities, got %d", len(expected), len(entities))
	}

	for _, expectedEntity := range expected {
		found := false
		for _, readEntity := range entities {
			if readEntity.GetUuid() == expectedEntity.GetUuid() {
				found = true
				if readEntity.GetFieldValue("Name") != expectedEntity.GetFieldValue("Name") {
					t.Fatalf("Name mismatch for UUID %v", readEntity.GetUuid())
				}
				if readEntity.GetFieldValue("Surname") != expectedEntity.GetFieldValue("Surname") {
					t.Fatalf("Surname mismatch for UUID %v", readEntity.GetUuid())
				}
				break
			}
		}
		if !found {
			t.Fatalf("Expected entity UUID %v not found", expectedEntity.GetUuid())
		}
	}
}

func BenchmarkMessageInsert(b *testing.B) {
	c, err := ConnectToNode(GlobalNode)
	defer util.FileDelete(filepath.Join(GlobalNode.Path.Data, "SampleV1.bin"))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		entity := SampleV1.Sample{
			Name:    fmt.Sprintf("Something%d", i),
			Surname: fmt.Sprintf("Else%d", i),
		}
		if _, err := entity.DbInsert(c); err != nil {
			b.Fatal(err)
		}
	}
}
