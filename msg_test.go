package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/rah-0/hyperion/disk"
	SampleV1 "github.com/rah-0/hyperion/entities/Sample/v1"
	"github.com/rah-0/hyperion/model"
	"github.com/rah-0/hyperion/register"
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

	msg, err := entity.DbInsert(c)
	if err != nil {
		t.Fatal(err)
	}

	if msg.Status != model.StatusSuccess {
		t.Fatalf("Unexpected status: got %v, want %v", msg.Status, model.StatusSuccess)
	}
}

func TestMessageInsertAndDelete(t *testing.T) {
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		t.Fatal(err)
	}

	entity := SampleV1.Sample{
		Name:    "Something",
		Surname: "Else",
	}

	msg, err := entity.DbInsert(c)
	if err != nil {
		t.Fatal(err)
	}

	if msg.Status != model.StatusSuccess {
		t.Fatalf("Unexpected status: got %v, want %v", msg.Status, model.StatusSuccess)
	}

	msg, err = entity.DbDelete(c)
	if err != nil {
		t.Fatal(err)
	}

	if msg.Status != model.StatusSuccess {
		t.Fatalf("Unexpected status: got %v, want %v", msg.Status, model.StatusSuccess)
	}
}

func TestMessageInsert1000(t *testing.T) {
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		t.Fatal(err)
	}

	var expected []SampleV1.Sample
	for i := 0; i < 1000; i++ {
		entity := SampleV1.Sample{
			Name:    fmt.Sprintf("Something%d", i),
			Surname: fmt.Sprintf("Else%d", i),
		}

		msg, err := entity.DbInsert(c)
		if err != nil {
			t.Fatal(err)
		}
		if msg.Status != model.StatusSuccess {
			t.Fatalf("Unexpected status: got %v, want %v", msg.Status, model.StatusSuccess)
		}

		expected = append(expected, entity)
	}

	d := disk.NewDisk()
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

		msg, err := entity.DbInsert(c)
		if err != nil {
			b.Fatal(err)
		}

		if msg.Status != model.StatusSuccess {
			b.Fatalf("Unexpected status: got %v, want %v", msg.Status, model.StatusSuccess)
		}
	}
}

func TestMessageUpdate(t *testing.T) {
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		t.Fatal(err)
	}

	entity := &SampleV1.Sample{
		Name:    "Initial",
		Surname: "User",
	}
	if _, err := entity.DbInsert(c); err != nil {
		t.Fatal(err)
	}

	// Modify values
	entity.Name = "Updated"
	entity.Surname = "User"

	resp, err := entity.DbUpdate(c)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != model.StatusSuccess {
		t.Fatalf("Unexpected update response status: %v", resp.Status)
	}

	all, err := SampleV1.DbGetAll(c)
	if err != nil {
		t.Fatal(err)
	}

	var match *SampleV1.Sample
	for _, e := range all {
		if e.GetUuid() == entity.GetUuid() {
			match = e
			break
		}
	}

	if match == nil {
		t.Fatalf("Updated entity not found")
	}
	if match.Name != "Updated" || match.Surname != "User" {
		t.Fatalf("Update not applied correctly. Got Name: %s, Surname: %s", match.Name, match.Surname)
	}
}

func TestMessageGetAll(t *testing.T) {
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		t.Fatal(err)
	}

	var inserted []*SampleV1.Sample

	// Insert 3 entities and track them
	for i := 0; i < 3; i++ {
		entity := &SampleV1.Sample{
			Name:    fmt.Sprintf("Name%d", i),
			Surname: fmt.Sprintf("Surname%d", i),
		}
		if _, err := entity.DbInsert(c); err != nil {
			t.Fatal(err)
		}
		inserted = append(inserted, entity)
	}

	// Get all from remote
	entities, err := SampleV1.DbGetAll(c)
	if err != nil {
		t.Fatal(err)
	}

	if len(entities) < len(inserted) {
		t.Fatalf("Expected at least %d entities, got %d", len(inserted), len(entities))
	}

	// Check UUIDs exist in received list
	for _, ins := range inserted {
		found := false
		for _, got := range entities {
			if ins.Uuid == got.Uuid {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Expected UUID %v not found in response", ins.Uuid)
		}
	}
}
