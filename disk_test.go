package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"

	. "github.com/rah-0/hyperion/register"
	. "github.com/rah-0/hyperion/util"
)

type TestDataSample struct {
	Name    string
	Surname string
}

func TestDataWriteSingleRecord(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	t.Cleanup(func() { FileDelete(d.Path) })

	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		// Create and write a single entity
		instance := e.New()
		instance.SetFieldValue("Name", "John")
		instance.SetFieldValue("Surname", "Doe")
		if err := instance.Encode(); err != nil {
			t.Fatal(err)
		}
		if _, err := d.DataWrite(instance.GetBufferData()); err != nil {
			t.Fatalf("DataWrite failed: %v", err)
		}
		bufferLen := instance.GetBuffer().Len()
		instance.BufferReset()

		// Verify file contents
		file, err := os.Open(d.Path)
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		// Read length field
		var length uint64
		if err = binary.Read(file, binary.LittleEndian, &length); err != nil {
			t.Fatalf("Failed to read length field: %v", err)
		}

		// Read status flag
		var status byte
		if err = binary.Read(file, binary.LittleEndian, &status); err != nil {
			t.Fatalf("Failed to read status flag: %v", err)
		}

		// Validate length and status flag
		expectedLength := uint64(bufferLen + LEN_BYTE_STATUS)
		if length != expectedLength {
			t.Fatalf("Incorrect length: expected %d, got %d", expectedLength, length)
		}

		if status != STATUS_BYTE_ACTIVE {
			t.Fatalf("Incorrect status flag: expected %d, got %d", STATUS_BYTE_ACTIVE, status)
		}
	}
}

func TestDataWriteReadMultipleRecords(t *testing.T) {
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
			if err := instance.Encode(); err != nil {
				t.Fatal(err)
			}
			if _, err := d.DataWrite(instance.GetBufferData()); err != nil {
				t.Fatalf("DataWrite failed for %s %s: %v", name, surname, err)
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

		entities, err := d.DataReadAll(e)
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

func TestDeleteRecord(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	t.Cleanup(func() {
		FileDelete(d.Path)
	})

	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		instance1 := e.New()
		instance1.SetFieldValue("Name", "John")
		instance1.SetFieldValue("Surname", "Doe")
		if err := instance1.Encode(); err != nil {
			t.Fatal(err)
		}
		offset, err := d.DataWrite(instance1.GetBufferData())
		if err != nil {
			t.Fatal(err)
		}
		instance1.SetOffset(offset)
		instance1.BufferReset()

		instance2 := e.New()
		instance2.SetFieldValue("Name", "Jane")
		instance2.SetFieldValue("Surname", "Doe")
		if err = instance2.Encode(); err != nil {
			t.Fatal(err)
		}
		offset, err = d.DataWrite(instance2.GetBufferData())
		if err != nil {
			t.Fatal(err)
		}
		instance2.SetOffset(offset)
		instance2.BufferReset()

		// Delete the first record
		if err := d.DataChangeStatus(instance1.GetOffset(), STATUS_BYTE_DELETED); err != nil {
			t.Fatalf("Failed to delete first record: %v", err)
		}

		// Verify deletion by reading file bytes directly
		file, err := os.Open(d.Path)
		if err != nil {
			t.Fatalf("Failed to open file for verification: %v", err)
		}
		defer file.Close()

		// Read length field of the first record
		var length uint64
		if err = binary.Read(file, binary.LittleEndian, &length); err != nil {
			t.Fatalf("Failed to read first record length field: %v", err)
		}

		// Read status flag of the first record
		var status byte
		if err = binary.Read(file, binary.LittleEndian, &status); err != nil {
			t.Fatalf("Failed to read status flag: %v", err)
		}

		// Validate that the status flag is now set to 1
		if status != 1 {
			t.Fatalf("First record was not properly status, expected status flag = 1, got %d", status)
		}

		// Perform cleanup
		if err = d.DataCleanup(); err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}

		// Read remaining records after cleanup
		remainingEntities, err := d.DataReadAll(e)
		if err != nil {
			t.Fatalf("Failed to read remaining entities after cleanup: %v", err)
		}

		// Ensure only one entity remains in the file after cleanup
		if len(remainingEntities) != 1 {
			t.Fatalf("Expected 1 remaining entity after cleanup, got %d", len(remainingEntities))
		}

		// Ensure the remaining record is the second one
		if remainingEntities[0].GetFieldValue("Name") != "Jane" {
			t.Fatalf("Unexpected remaining record: expected Jane, got %v", remainingEntities[0].GetFieldValue("Name"))
		}
	}
}

func TestDataCleanupRemovesDeletedRecords(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	t.Cleanup(func() { FileDelete(d.Path) })

	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		// Create and write first record
		instance1 := e.New()
		instance1.SetFieldValue("Name", "Alice")
		instance1.SetFieldValue("Surname", "Smith")
		if err := instance1.Encode(); err != nil {
			t.Fatal(err)
		}
		offset, err := d.DataWrite(instance1.GetBufferData())
		if err != nil {
			t.Fatalf("Failed to write first record: %v", err)
		}
		instance1.SetOffset(offset)
		instance1.BufferReset()

		// Create and write second record
		instance2 := e.New()
		instance2.SetFieldValue("Name", "Bob")
		instance2.SetFieldValue("Surname", "Johnson")
		if err = instance2.Encode(); err != nil {
			t.Fatal(err)
		}
		offset, err = d.DataWrite(instance2.GetBufferData())
		if err != nil {
			t.Fatalf("Failed to write second record: %v", err)
		}
		instance2.SetOffset(offset)
		instance2.BufferReset()

		// Create and write third record
		instance3 := e.New()
		instance3.SetFieldValue("Name", "Charlie")
		instance3.SetFieldValue("Surname", "Brown")
		if err = instance3.Encode(); err != nil {
			t.Fatal(err)
		}
		offset, err = d.DataWrite(instance3.GetBufferData())
		if err != nil {
			t.Fatalf("Failed to write third record: %v", err)
		}
		instance3.SetOffset(offset)
		instance3.BufferReset()

		// Delete the second record
		if err = d.DataChangeStatus(instance2.GetOffset(), STATUS_BYTE_DELETED); err != nil {
			t.Fatalf("Failed to delete second record: %v", err)
		}

		// Get initial file size before cleanup
		initialSize, _ := readFileSize(d.Path)

		// Perform cleanup
		if err = d.DataCleanup(); err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}

		// Get new file size after cleanup
		finalSize, _ := readFileSize(d.Path)

		// Ensure the file size is reduced
		if finalSize >= initialSize {
			t.Fatalf("Cleanup did not reduce file size: before=%d, after=%d", initialSize, finalSize)
		}

		// Read remaining records
		remainingEntities, err := d.DataReadAll(e)
		if err != nil {
			t.Fatalf("Failed to read remaining entities after cleanup: %v", err)
		}

		// Ensure only two entities remain
		if len(remainingEntities) != 2 {
			t.Fatalf("Expected 2 remaining entities after cleanup, got %d", len(remainingEntities))
		}

		// Verify that the remaining entities are Alice and Charlie
		if remainingEntities[0].GetFieldValue("Name") != "Alice" {
			t.Fatalf("Unexpected first entity: expected Alice, got %v", remainingEntities[0].GetFieldValue("Name"))
		}
		if remainingEntities[1].GetFieldValue("Name") != "Charlie" {
			t.Fatalf("Unexpected second entity: expected Charlie, got %v", remainingEntities[1].GetFieldValue("Name"))
		}
	}
}

func TestDataCleanupPreservesValidRecords(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	t.Cleanup(func() { FileDelete(d.Path) })

	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		// Write valid records
		instance1 := e.New()
		instance1.SetFieldValue("Name", "David")
		instance1.SetFieldValue("Surname", "Miller")
		if err := instance1.Encode(); err != nil {
			t.Fatal(err)
		}
		if _, err := d.DataWrite(instance1.GetBufferData()); err != nil {
			t.Fatalf("Failed to write first record: %v", err)
		}
		instance1.BufferReset()

		instance2 := e.New()
		instance2.SetFieldValue("Name", "Emma")
		instance2.SetFieldValue("Surname", "Davis")
		if err := instance2.Encode(); err != nil {
			t.Fatal(err)
		}
		if _, err := d.DataWrite(instance2.GetBufferData()); err != nil {
			t.Fatalf("Failed to write second record: %v", err)
		}
		instance2.BufferReset()

		// Get initial file size before cleanup
		initialSize, _ := readFileSize(d.Path)

		// Perform cleanup (nothing should change)
		if err := d.DataCleanup(); err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}

		// Get new file size after cleanup
		finalSize, _ := readFileSize(d.Path)

		// Ensure the file size remains unchanged
		if finalSize != initialSize {
			t.Fatalf("Cleanup modified file unexpectedly: before=%d, after=%d", initialSize, finalSize)
		}

		// Read back records to confirm they remain the same
		entities, err := d.DataReadAll(e)
		if err != nil {
			t.Fatalf("Failed to read entities after cleanup: %v", err)
		}
		if len(entities) != 2 {
			t.Fatalf("Expected 2 entities after cleanup, got %d", len(entities))
		}

		if entities[0].GetFieldValue("Name") != "David" {
			t.Fatalf("Unexpected first entity: expected David, got %v", entities[0].GetFieldValue("Name"))
		}
		if entities[1].GetFieldValue("Name") != "Emma" {
			t.Fatalf("Unexpected second entity: expected Emma, got %v", entities[1].GetFieldValue("Name"))
		}
	}
}

func TestDataCleanupHandlesMultipleDeletedRecords(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	t.Cleanup(func() { FileDelete(d.Path) })

	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		// Write multiple records
		instances := []Model{
			e.New(),
			e.New(),
			e.New(),
			e.New(),
		}
		names := []string{"Keep1", "Delete1", "Keep2", "Delete2"}

		for i, name := range names {
			instances[i].SetFieldValue("Name", name)
			instances[i].SetFieldValue("Surname", "Doe")
			if err := instances[i].Encode(); err != nil {
				t.Fatal(err)
			}
			offset, err := d.DataWrite(instances[i].GetBufferData())
			if err != nil {
				t.Fatalf("Failed to write record %s: %v", name, err)
			}
			instances[i].SetOffset(offset)
			instances[i].BufferReset()
		}

		// Delete second and fourth records
		if err := d.DataChangeStatus(instances[1].GetOffset(), STATUS_BYTE_DELETED); err != nil {
			t.Fatalf("Failed to delete second record: %v", err)
		}
		if err := d.DataChangeStatus(instances[3].GetOffset(), STATUS_BYTE_DELETED); err != nil {
			t.Fatalf("Failed to delete fourth record: %v", err)
		}

		// Get initial file size before cleanup
		initialSize, _ := readFileSize(d.Path)

		// Perform cleanup
		if err := d.DataCleanup(); err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}

		// Get new file size after cleanup
		finalSize, _ := readFileSize(d.Path)

		// Ensure the file size is reduced
		if finalSize >= initialSize {
			t.Fatalf("Cleanup did not reduce file size: before=%d, after=%d", initialSize, finalSize)
		}

		// Read remaining records
		remainingEntities, err := d.DataReadAll(e)
		if err != nil {
			t.Fatalf("Failed to read remaining entities after cleanup: %v", err)
		}

		// Ensure only "Keep1" and "Keep2" remain
		if len(remainingEntities) != 2 {
			t.Fatalf("Unexpected number of remaining records: expected 2, got %d", len(remainingEntities))
		}
		if remainingEntities[0].GetFieldValue("Name") != "Keep1" {
			t.Fatalf("Unexpected first remaining entity: expected Keep1, got %v", remainingEntities[0].GetFieldValue("Name"))
		}
		if remainingEntities[1].GetFieldValue("Name") != "Keep2" {
			t.Fatalf("Unexpected second remaining entity: expected Keep2, got %v", remainingEntities[1].GetFieldValue("Name"))
		}
	}
}

func TestDataCleanupRemovesAllEntities(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	t.Cleanup(func() { FileDelete(d.Path) })

	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		instances := []Model{
			e.New(),
			e.New(),
			e.New(),
		}
		names := []string{"Entity1", "Entity2", "Entity3"}

		for i, name := range names {
			instances[i].SetFieldValue("Name", name)
			instances[i].SetFieldValue("Surname", "Doe")
			if err := instances[i].Encode(); err != nil {
				t.Fatal(err)
			}
			offset, err := d.DataWrite(instances[i].GetBufferData())
			if err != nil {
				t.Fatal(err)
			}
			instances[i].SetOffset(offset)
			instances[i].BufferReset()
		}

		// Mark all entities as deleted
		for _, entity := range instances {
			if err := d.DataChangeStatus(entity.GetOffset(), STATUS_BYTE_DELETED); err != nil {
				t.Fatalf("Failed to delete entity: %v", err)
			}
		}

		// Get initial file size before cleanup
		initialSize, _ := readFileSize(d.Path)

		// Perform cleanup
		if err := d.DataCleanup(); err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}

		// Get new file size after cleanup
		finalSize, _ := readFileSize(d.Path)

		// Ensure the file size is reduced to **almost zero** (only metadata might remain)
		if finalSize >= initialSize || finalSize > 0 {
			t.Fatalf("Cleanup did not remove all records: before=%d, after=%d", initialSize, finalSize)
		}

		// Read remaining records (should be **empty**)
		remainingEntities, err := d.DataReadAll(e)
		if err != nil {
			t.Fatalf("Failed to read remaining entities after cleanup: %v", err)
		}

		// Ensure no entities remain
		if len(remainingEntities) != 0 {
			t.Fatalf("Expected 0 remaining entities after cleanup, got %d", len(remainingEntities))
		}
	}
}

func TestDataUpdateSingleRecord(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	t.Cleanup(func() { FileDelete(d.Path) })

	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		// Write initial record
		instance := e.New()
		instance.SetFieldValue("Name", "Alice")
		instance.SetFieldValue("Surname", "Smith")
		if err := instance.Encode(); err != nil {
			t.Fatal(err)
		}
		offset, err := d.DataWrite(instance.GetBufferData())
		if err != nil {
			t.Fatalf("Failed to write initial record: %v", err)
		}
		instance.SetOffset(offset)
		instance.BufferReset()

		// Read and get offset
		entities, err := d.DataReadAll(e)
		if err != nil {
			t.Fatalf("Failed to read entities: %v", err)
		}
		if len(entities) != 1 {
			t.Fatalf("Expected 1 entity, got %d", len(entities))
		}
		oldOffset := entities[0].GetOffset()

		// Update the record with new data
		updatedInstance := e.New()
		updatedInstance.SetFieldValue("Name", "AliceUpdated")
		updatedInstance.SetFieldValue("Surname", "SmithUpdated")
		if err = updatedInstance.Encode(); err != nil {
			t.Fatal(err)
		}
		newOffset, err := d.DataUpdate(oldOffset, updatedInstance.GetBufferData())
		if err != nil {
			t.Fatalf("Failed to update record: %v", err)
		}
		updatedInstance.SetOffset(newOffset)
		updatedInstance.BufferReset()

		// Read entities after update
		entities, err = d.DataReadAll(e)
		if err != nil {
			t.Fatalf("Failed to read entities after update: %v", err)
		}

		// Ensure only one entity exists after update
		if len(entities) != 1 {
			t.Fatalf("Expected 1 entity after update, got %d", len(entities))
		}

		// Ensure updated data is present
		if entities[0].GetFieldValue("Name") != "AliceUpdated" {
			t.Fatalf("Update failed: expected Name = AliceUpdated, got %v", entities[0].GetFieldValue("Name"))
		}
		if entities[0].GetFieldValue("Surname") != "SmithUpdated" {
			t.Fatalf("Update failed: expected Surname = SmithUpdated, got %v", entities[0].GetFieldValue("Surname"))
		}

		// Ensure offset changed
		if oldOffset == newOffset {
			t.Fatalf("Offset did not update: expected new offset different from old offset")
		}
	}
}

func TestDataUpdateMultipleTimes(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	t.Cleanup(func() { FileDelete(d.Path) })

	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		// Write initial record
		instance := e.New()
		instance.SetFieldValue("Name", "Bob")
		instance.SetFieldValue("Surname", "Johnson")
		if err := instance.Encode(); err != nil {
			t.Fatal(err)
		}
		offset, err := d.DataWrite(instance.GetBufferData())
		if err != nil {
			t.Fatalf("Failed to write initial record: %v", err)
		}
		instance.SetOffset(offset)
		instance.BufferReset()

		// First update
		instance.SetFieldValue("Name", "Bob1")
		instance.SetFieldValue("Surname", "Johnson1")
		if err = instance.Encode(); err != nil {
			t.Fatal(err)
		}
		offset, err = d.DataUpdate(instance.GetOffset(), instance.GetBufferData())
		if err != nil {
			t.Fatalf("Failed to write initial record: %v", err)
		}
		instance.SetOffset(offset)
		instance.BufferReset()

		// Second update
		instance.SetFieldValue("Name", "Bob2")
		instance.SetFieldValue("Surname", "Johnson2")
		if err = instance.Encode(); err != nil {
			t.Fatal(err)
		}
		offset, err = d.DataUpdate(instance.GetOffset(), instance.GetBufferData())
		if err != nil {
			t.Fatalf("Failed to write initial record: %v", err)
		}
		instance.SetOffset(offset)
		instance.BufferReset()

		// Read entities after multiple updates
		entities, err := d.DataReadAll(e)
		if err != nil {
			t.Fatalf("Failed to read entities after multiple updates: %v", err)
		}

		// Ensure only one entity exists after multiple updates
		if len(entities) != 1 {
			t.Fatalf("Expected 1 entity after multiple updates, got %d", len(entities))
		}

		// Ensure last update is present
		if entities[0].GetFieldValue("Name") != "Bob2" {
			t.Fatalf("Update failed: expected Name = Bob2, got %v", entities[0].GetFieldValue("Name"))
		}

		// Ensure offset changed after multiple updates
		if offset == entities[0].GetOffset() {
			t.Fatalf("Offset did not update after multiple updates")
		}
	}
}

func TestDataUpdateAllRecords(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	t.Cleanup(func() { FileDelete(d.Path) })

	if len(Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		// Write multiple records
		instances := []Model{
			e.New(),
			e.New(),
			e.New(),
		}
		names := []string{"Original1", "Original2", "Original3"}
		var offsets []uint64

		for i, name := range names {
			instances[i].SetFieldValue("Name", name)
			instances[i].SetFieldValue("Surname", "Doe")
			if err := instances[i].Encode(); err != nil {
				t.Fatal(err)
			}
			offset, err := d.DataWrite(instances[i].GetBufferData())
			if err != nil {
				t.Fatalf("Failed to write record %s: %v", name, err)
			}
			instances[i].SetOffset(offset)
			offsets = append(offsets, offset)
			instances[i].BufferReset()
		}

		// Read all to get offsets
		entities, err := d.DataReadAll(e)
		if err != nil {
			t.Fatalf("Failed to read entities: %v", err)
		}

		// Update all records
		newNames := []string{"Updated1", "Updated2", "Updated3"}
		newOffsets := make([]uint64, len(entities))

		for i, entity := range entities {
			updatedInstance := e.New()
			updatedInstance.SetFieldValue("Name", newNames[i])
			updatedInstance.SetFieldValue("Surname", "DoeUpdated")
			if err = updatedInstance.Encode(); err != nil {
				t.Fatal(err)
			}
			newOffset, err := d.DataUpdate(entity.GetOffset(), updatedInstance.GetBufferData())
			if err != nil {
				t.Fatalf("Failed to update entity %d: %v", i, err)
			}
			newOffsets[i] = newOffset
			updatedInstance.SetOffset(newOffset)
			updatedInstance.BufferReset()
		}

		// Read entities after updating all
		entities, err = d.DataReadAll(e)
		if err != nil {
			t.Fatalf("Failed to read entities after updating all: %v", err)
		}

		// Ensure the same number of entities exist
		if len(entities) != 3 {
			t.Fatalf("Expected 3 entities after updates, got %d", len(entities))
		}

		// Verify updated names
		for i, entity := range entities {
			if entity.GetFieldValue("Name") != newNames[i] {
				t.Fatalf("Update failed for entity %d: expected %s, got %v", i, newNames[i], entity.GetFieldValue("Name"))
			}
			// Ensure offset has been updated
			if newOffsets[i] == offsets[i] {
				t.Fatalf("Offset did not update for entity %d", i)
			}
		}
	}
}

func readFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
