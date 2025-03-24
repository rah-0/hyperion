package Sample

import (
	"testing"

	"github.com/google/uuid"

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

func TestMemoryRemove(t *testing.T) {
	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		instance := e.New()
		u := uuid.New()
		instance.SetFieldValue("Uuid", u)
		instance.SetFieldValue("Name", "Ian")
		instance.SetFieldValue("Surname", "Miller")
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
	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		u := uuid.New()
		instance := e.New()
		instance.SetFieldValue("Uuid", u)
		instance.SetFieldValue("Name", "Jane")
		instance.SetFieldValue("Surname", "Doe")
		instance.MemoryAdd()

		instance.SetFieldValue("Name", "Janet")
		instance.MemoryUpdate()

		all := instance.MemoryGetAll()
		found := false
		for _, m := range all {
			if m.GetUuid() == u {
				found = true
				if m.GetFieldValue("Name") != "Janet" {
					t.Fatalf("Expected name to be 'Janet', got %v", m.GetFieldValue("Name"))
				}
			}
		}
		if !found {
			t.Fatalf("Instance with UUID %v not found in memory", u)
		}
	}
}

func TestMemorySet(t *testing.T) {
	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		instance1 := e.New()
		u1 := uuid.New()
		instance1.SetFieldValue("Uuid", u1)
		instance1.SetFieldValue("Name", "Kyle")

		instance2 := e.New()
		u2 := uuid.New()
		instance2.SetFieldValue("Uuid", u2)
		instance2.SetFieldValue("Name", "Laura")

		models := []Model{instance1, instance2}
		instance1.MemorySet(models)

		all := instance1.MemoryGetAll()
		found1 := false
		found2 := false
		for _, m := range all {
			if m.GetUuid() == u1 {
				found1 = true
			}
			if m.GetUuid() == u2 {
				found2 = true
			}
		}
		if !found1 || !found2 {
			t.Fatalf("Expected both UUIDs in memory. Found1: %v, Found2: %v", found1, found2)
		}
	}
}
