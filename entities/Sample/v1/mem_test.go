package SampleV1

import (
	"testing"

	"github.com/google/uuid"

	"github.com/rah-0/hyperion/register"
)

func TestMemoryAdd(t *testing.T) {
	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}

		instance := e.EntityExtension.New()
		instance.SetFieldValue(FieldName, "Alice")
		instance.SetFieldValue(FieldSurname, "Smith")
		instance.MemoryAdd()

		allInstances := instance.MemoryGetAll()
		if len(allInstances) != 1 {
			t.Fatalf("Expected 1 instance in memory, got %d", len(allInstances))
		}
	}
}

func TestMemoryClear(t *testing.T) {
	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}

		instance1 := e.EntityExtension.New()
		instance1.SetFieldValue(FieldName, "Charlie")
		instance1.SetFieldValue(FieldSurname, "Brown")

		instance2 := e.EntityExtension.New()
		instance2.SetFieldValue(FieldName, "Daisy")
		instance2.SetFieldValue(FieldSurname, "Williams")

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
	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}

		instance1 := e.EntityExtension.New()
		instance1.SetFieldValue(FieldName, "Eve")
		instance1.SetFieldValue(FieldSurname, "Clark")

		instance2 := e.EntityExtension.New()
		instance2.SetFieldValue(FieldName, "Frank")
		instance2.SetFieldValue(FieldSurname, "Harris")

		instance1.MemoryAdd()
		instance2.MemoryAdd()

		allInstances := instance1.MemoryGetAll()
		if len(allInstances) != 2 {
			t.Fatalf("Expected 2 instances in memory, got %d", len(allInstances))
		}
	}
}

func TestMemoryContains(t *testing.T) {
	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}

		instance1 := e.EntityExtension.New()
		instance1.SetFieldValue(FieldName, "Grace")
		instance1.SetFieldValue(FieldSurname, "Lee")

		instance2 := e.EntityExtension.New()
		instance2.SetFieldValue(FieldName, "Hank")
		instance2.SetFieldValue(FieldSurname, "Martinez")

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

func TestMemoryRemove(t *testing.T) {
	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}

		instance := e.EntityExtension.New()
		u := uuid.New()
		instance.SetFieldValue(FieldUuid, u)
		instance.SetFieldValue(FieldName, "Ian")
		instance.SetFieldValue(FieldSurname, "Miller")
		instance.MemoryAdd()

		instance.MemoryRemove()
		all := instance.MemoryGetAll()

		for _, m := range all {
			if m.GetUuid() == u {
				t.Fatalf("Instance with UUID %v should have been removed from memory", u)
			}
		}
	}
}

func TestMemoryUpdate(t *testing.T) {
	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}

		u := uuid.New()
		instance := e.EntityExtension.New()
		instance.SetFieldValue(FieldUuid, u)
		instance.SetFieldValue(FieldName, "Jane")
		instance.SetFieldValue(FieldSurname, "Doe")
		instance.MemoryAdd()

		instance.SetFieldValue(FieldName, "Janet")
		instance.MemoryUpdate()

		all := instance.MemoryGetAll()
		found := false
		for _, m := range all {
			if m.GetUuid() == u {
				found = true
				if m.GetFieldValue(FieldName) != "Janet" {
					t.Fatalf("Expected name to be 'Janet', got %v", m.GetFieldValue(FieldName))
				}
			}
		}
		if !found {
			t.Fatalf("Instance with UUID %v not found in memory", u)
		}
	}
}
