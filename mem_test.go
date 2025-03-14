package main

import (
	"testing"

	. "github.com/rah-0/hyperion/register"
)

func TestMemoryAdd(t *testing.T) {
	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		instance := e.New()
		instance.SetFieldValue("Name", "Alice")
		instance.SetFieldValue("Surname", "Smith")
		instance.MemoryAdd()

		allInstances := instance.MemoryGetAll()
		if len(allInstances) != 1 {
			t.Fatalf("Expected 1 instance in memory, got %d", len(allInstances))
		}
	}
}

func TestMemoryClear(t *testing.T) {
	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		instance1 := e.New()
		instance1.SetFieldValue("Name", "Charlie")
		instance1.SetFieldValue("Surname", "Brown")

		instance2 := e.New()
		instance2.SetFieldValue("Name", "Daisy")
		instance2.SetFieldValue("Surname", "Williams")

		instance1.MemoryAdd()
		instance2.MemoryAdd()

		instance1.MemoryClear()

		allInstances := instance1.MemoryGetAll()
		if len(allInstances) != 0 {
			t.Fatalf("Expected memory to be empty after MemoryClear, got %d", len(allInstances))
		}
	}
}

func TestMemoryGetAll(t *testing.T) {
	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		instance1 := e.New()
		instance1.SetFieldValue("Name", "Eve")
		instance1.SetFieldValue("Surname", "Clark")

		instance2 := e.New()
		instance2.SetFieldValue("Name", "Frank")
		instance2.SetFieldValue("Surname", "Harris")

		instance1.MemoryAdd()
		instance2.MemoryAdd()

		allInstances := instance1.MemoryGetAll()
		if len(allInstances) != 2 {
			t.Fatalf("Expected 2 instances in memory, got %d", len(allInstances))
		}
	}
}

func TestMemoryContains(t *testing.T) {
	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		instance1 := e.New()
		instance1.SetFieldValue("Name", "Grace")
		instance1.SetFieldValue("Surname", "Lee")

		instance2 := e.New()
		instance2.SetFieldValue("Name", "Hank")
		instance2.SetFieldValue("Surname", "Martinez")

		instance1.MemoryAdd()

		// Verify that instance1 is contained in memory
		if !instance1.MemoryContains(instance1) {
			t.Fatal("MemoryContains failed for instance1")
		}

		// Verify that instance2 is NOT contained in memory
		if instance1.MemoryContains(instance2) {
			t.Fatal("MemoryContains should return false for instance2")
		}
	}
}
