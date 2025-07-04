package disk

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/rah-0/testmark/testutil"

	SampleV1 "github.com/rah-0/hyperion/entities/Sample/v1"
	_ "github.com/rah-0/hyperion/template"

	"github.com/rah-0/hyperion/register"
	"github.com/rah-0/hyperion/util"
)

const (
	FieldUuid    = 1
	FieldDeleted = 2
	FieldName    = 3
	FieldSurname = 4
)

func TestMain(m *testing.M) {
	testutil.TestMainWrapper(testutil.TestConfig{
		M: m,
		LoadResources: func() error {
			return SampleV1.Register()
		},
		UnloadResources: func() error {
			return nil
		},
	})
}

// TestDiskClose tests the Close method ensuring proper resource cleanup
func TestDiskClose(t *testing.T) {
	// Create and open a disk
	d := NewDisk()
	d.WithNewRandomPath()
	err := d.OpenFile()
	if err != nil {
		t.Fatalf("Failed to open disk file: %v", err)
	}
	t.Cleanup(func() { util.FileDelete(d.Path) })

	// Ensure background goroutine is running by checking if StopChan and SyncTicker exist
	if d.StopChan == nil {
		t.Fatal("StopChan should be initialized after OpenFile")
	}
	if d.SyncTicker == nil {
		t.Fatal("SyncTicker should be initialized after OpenFile")
	}

	// Close the disk
	err = d.Close()
	if err != nil {
		t.Fatalf("Failed to close disk: %v", err)
	}

	// Verify resources are properly cleaned up
	// After Close(), the fields should be nil as they've been cleaned up
	if d.SyncTicker != nil {
		t.Error("SyncTicker should be nil after Close()")
	}
	
	if d.StopChan != nil {
		t.Error("StopChan should be nil after Close()")
	}

	// The file handle should be closed, attempting to use it should fail
	_, err = d.File.Stat()
	if err == nil {
		t.Error("File handle should be closed after Close()")
	}
}

// TestDiskCloseMultipleTimes ensures that calling Close multiple times is safe
func TestDiskCloseMultipleTimes(t *testing.T) {
	// Create and open a disk
	d := NewDisk()
	d.WithNewRandomPath()
	err := d.OpenFile()
	if err != nil {
		t.Fatalf("Failed to open disk file: %v", err)
	}
	t.Cleanup(func() { util.FileDelete(d.Path) })

	// First close should succeed
	err = d.Close()
	if err != nil {
		t.Fatalf("First close failed: %v", err)
	}

	// Second close should not error
	err = d.Close()
	if err != nil {
		t.Fatalf("Second close should be safe but got error: %v", err)
	}
}

// TestDiskCloseWithoutOpen tests closing a disk that was never opened
func TestDiskCloseWithoutOpen(t *testing.T) {
	// Create disk without opening
	d := NewDisk()
	d.WithNewRandomPath()

	// Close should not panic or error
	err := d.Close()
	if err != nil {
		t.Fatalf("Close on unopened disk should not error: %v", err)
	}
}

func TestDataWrite(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	d.OpenFile()
	t.Cleanup(func() { util.FileDelete(d.Path) })

	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}

		instance := e.EntityExtension.New()
		instance.SetFieldValue(FieldUuid, uuid.New())
		instance.SetFieldValue(FieldName, "John")
		instance.SetFieldValue(FieldSurname, "Doe")
		if err := instance.Encode(); err != nil {
			t.Fatal(err)
		}

		err := d.DataWrite(instance.GetBufferData())
		if err != nil {
			t.Fatalf("DataWrite failed: %v", err)
		}

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

		// Validate length
		expectedLength := uint64(instance.GetBuffer().Len())
		if length != expectedLength {
			t.Fatalf("Incorrect length: expected %d, got %d", expectedLength, length)
		}

		// Read the data back
		data := make([]byte, length)
		if _, err = file.Read(data); err != nil {
			t.Fatalf("Failed to read entity data: %v", err)
		}

		// Decode entity
		readInstance := e.EntityExtension.New()
		readInstance.SetBufferData(data)
		if err = readInstance.Decode(); err != nil {
			t.Fatalf("Failed to decode entity: %v", err)
		}

		// Verify UUID, Name, and Surname
		if readInstance.GetFieldValue(FieldUuid) != instance.GetFieldValue(FieldUuid) {
			t.Fatalf("UUID mismatch: expected %v, got %v",
				instance.GetFieldValue(FieldUuid), readInstance.GetFieldValue(FieldUuid))
		}
		if readInstance.GetFieldValue(FieldName) != "John" {
			t.Fatalf("Incorrect Name: expected 'John', got %s", readInstance.GetFieldValue(FieldName))
		}
		if readInstance.GetFieldValue(FieldSurname) != "Doe" {
			t.Fatalf("Incorrect Surname: expected 'Doe', got %s", readInstance.GetFieldValue(FieldSurname))
		}

		instance.BufferReset()
	}
}

func TestDataReadAll_InitialWrite(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	d.OpenFile()
	t.Cleanup(func() { util.FileDelete(d.Path) })

	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}
		d.WithEntity(e)

		// Write multiple unique entities
		expectedEntities := make(map[uuid.UUID]register.Model)
		for i := 0; i < 50; i++ {
			instance := e.EntityExtension.New()
			instance.SetFieldValue(FieldUuid, uuid.New())
			instance.SetFieldValue(FieldName, "User"+uuid.NewString())
			instance.SetFieldValue(FieldSurname, "Surname"+uuid.NewString())

			if err := instance.Encode(); err != nil {
				t.Fatal(err)
			}
			if err := d.DataWrite(instance.GetBufferData()); err != nil {
				t.Fatalf("DataWrite failed: %v", err)
			}
			instance.BufferReset()

			expectedEntities[instance.GetUuid()] = instance
		}

		// Read entities from disk
		entities, err := d.DataReadAll()
		if err != nil {
			t.Fatalf("DataReadAll failed: %v", err)
		}

		// Ensure correct entity count
		if len(entities) != len(expectedEntities) {
			t.Fatalf("Mismatch in entity count: expected %d, got %d", len(expectedEntities), len(register.Entities))
		}

		// Ensure UUIDs match
		for _, entity := range entities {
			entityUUID := entity.GetUuid()
			if _, exists := expectedEntities[entityUUID]; !exists {
				t.Fatalf("Unexpected entity found: UUID %v not in expectedEntities", entityUUID)
			}

			// Verify field values
			expected := expectedEntities[entityUUID]
			if entity.GetFieldValue(FieldName) != expected.GetFieldValue(FieldName) {
				t.Fatalf("Mismatch in Name for UUID %v", entityUUID)
			}
			if entity.GetFieldValue(FieldSurname) != expected.GetFieldValue(FieldSurname) {
				t.Fatalf("Mismatch in Surname for UUID %v", entityUUID)
			}
		}
	}
}

func TestDataReadAll_SingleEntity_WithUpdates(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	d.OpenFile()
	t.Cleanup(func() { util.FileDelete(d.Path) })

	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}
		d.WithEntity(e)

		// Write Initial EntityBase
		instance := e.EntityExtension.New()
		instance.SetFieldValue(FieldName, "OriginalUser")
		instance.SetFieldValue(FieldSurname, "OriginalSurname")
		if err := instance.Encode(); err != nil {
			t.Fatal(err)
		}
		err := d.DataWrite(instance.GetBufferData())
		if err != nil {
			t.Fatalf("DataWrite failed: %v", err)
		}
		instance.BufferReset()

		entityUUID := instance.GetUuid()

		// Perform Multiple Updates
		updateNames := []string{"Update1", "Update2", "Update3", "FinalUpdate"}
		updateSurnames := []string{"Surname1", "Surname2", "Surname3", "FinalSurname"}

		for i := range updateNames {
			updatedInstance := e.EntityExtension.New()
			updatedInstance.SetFieldValue(FieldUuid, entityUUID) // Keep same UUID
			updatedInstance.SetFieldValue(FieldName, updateNames[i])
			updatedInstance.SetFieldValue(FieldSurname, updateSurnames[i])
			if err = updatedInstance.Encode(); err != nil {
				t.Fatal(err)
			}
			if err = d.DataWrite(updatedInstance.GetBufferData()); err != nil {
				t.Fatalf("DataWrite failed: %v", err)
			}
			updatedInstance.BufferReset()
		}

		// Read from Disk
		entities, err := d.DataReadAll()
		if err != nil {
			t.Fatalf("DataReadAll failed: %v", err)
		}

		// Validate Data
		if len(register.Entities) != 1 {
			t.Fatalf("Expected only 1 entity after updates, but got %d", len(register.Entities))
		}

		finalEntity := entities[0]

		if finalEntity.GetFieldValue(FieldUuid) != entityUUID {
			t.Fatalf("UUID mismatch: expected %v, got %v",
				entityUUID, finalEntity.GetFieldValue(FieldUuid))
		}

		if finalEntity.GetFieldValue(FieldName) != "FinalUpdate" {
			t.Fatalf("Incorrect Name: expected 'FinalUpdate', got %s",
				finalEntity.GetFieldValue(FieldName))
		}

		if finalEntity.GetFieldValue(FieldSurname) != "FinalSurname" {
			t.Fatalf("Incorrect Surname: expected 'FinalSurname', got %s",
				finalEntity.GetFieldValue(FieldSurname))
		}
	}
}

func TestDataReadAll_MultipleEntities_WithUpdates(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	d.OpenFile()
	t.Cleanup(func() { util.FileDelete(d.Path) })

	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}
		d.WithEntity(e)

		// Write initial entities
		expectedEntities := make(map[uuid.UUID]register.Model)
		for i := 0; i < 10; i++ {
			instance := e.EntityExtension.New()
			instance.SetFieldValue(FieldUuid, uuid.New())
			instance.SetFieldValue(FieldName, "User"+uuid.NewString())
			instance.SetFieldValue(FieldSurname, "Surname"+uuid.NewString())

			if err := instance.Encode(); err != nil {
				t.Fatal(err)
			}
			if err := d.DataWrite(instance.GetBufferData()); err != nil {
				t.Fatalf("DataWrite failed: %v", err)
			}
			instance.BufferReset()

			expectedEntities[instance.GetUuid()] = instance
		}

		// Read entities from disk
		entities, err := d.DataReadAll()
		if err != nil {
			t.Fatalf("DataReadAll failed: %v", err)
		}

		// Select 5 entities for update
		for i := 0; i < 5; i++ {
			original := entities[i]
			entityUUID := original.GetUuid()

			// Create an updated version
			updatedInstance := e.EntityExtension.New()
			updatedInstance.SetFieldValue(FieldUuid, entityUUID) // Keep same UUID
			updatedInstance.SetFieldValue(FieldName, "UpdatedUser")
			updatedInstance.SetFieldValue(FieldSurname, "UpdatedSurname")

			if err = updatedInstance.Encode(); err != nil {
				t.Fatal(err)
			}
			if err = d.DataWrite(updatedInstance.GetBufferData()); err != nil {
				t.Fatalf("DataWrite failed: %v", err)
			}
			updatedInstance.BufferReset()

			// Replace in expected map
			expectedEntities[entityUUID] = updatedInstance
		}

		// Read entities after update
		entities, err = d.DataReadAll()
		if err != nil {
			t.Fatalf("DataReadAll failed: %v", err)
		}

		// Ensure count is still correct
		if len(entities) != len(expectedEntities) {
			t.Fatalf("Expected %d entities after updates, got %d", len(expectedEntities), len(register.Entities))
		}

		// Validate updated entities
		for _, entity := range entities {
			entityUUID := entity.GetUuid()
			expected := expectedEntities[entityUUID]

			if entity.GetFieldValue(FieldName) != expected.GetFieldValue(FieldName) {
				t.Fatalf("Mismatch in Name for UUID %v", entityUUID)
			}
			if entity.GetFieldValue(FieldSurname) != expected.GetFieldValue(FieldSurname) {
				t.Fatalf("Mismatch in Surname for UUID %v", entityUUID)
			}
		}
	}
}

func TestDataCleanup_NoDuplicates(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	d.OpenFile()
	t.Cleanup(func() { util.FileDelete(d.Path) })

	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}
		d.WithEntity(e)

		// Write multiple unique entities
		expectedEntities := make(map[uuid.UUID]register.Model)
		for i := 0; i < 1000; i++ {
			instance := e.EntityExtension.New()
			instance.SetFieldValue(FieldUuid, uuid.New())
			instance.SetFieldValue(FieldName, "User"+uuid.NewString())
			instance.SetFieldValue(FieldSurname, "Surname"+uuid.NewString())

			if err := instance.Encode(); err != nil {
				t.Fatal(err)
			}
			if err := d.DataWrite(instance.GetBufferData()); err != nil {
				t.Fatalf("DataWrite failed: %v", err)
			}
			instance.BufferReset()

			expectedEntities[instance.GetUuid()] = instance
		}

		// Perform cleanup (should do nothing)
		err := d.DataCleanup()
		if err != nil {
			t.Fatalf("DataCleanup failed: %v", err)
		}

		// Read entities after cleanup
		entities, err := d.DataReadAll()
		if err != nil {
			t.Fatalf("DataReadAll failed: %v", err)
		}

		// Ensure the number of entities remains unchanged
		if len(entities) != len(expectedEntities) {
			t.Fatalf("Mismatch in entity count after cleanup: expected %d, got %d", len(expectedEntities), len(register.Entities))
		}
	}
}

func TestDataCleanup_WithDuplicates(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	d.OpenFile()
	t.Cleanup(func() { util.FileDelete(d.Path) })

	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}
		d.WithEntity(e)

		// Write initial entities
		expectedEntities := make(map[uuid.UUID]register.Model)
		for i := 0; i < 1000; i++ {
			instance := e.EntityExtension.New()
			instance.SetFieldValue(FieldUuid, uuid.New())
			instance.SetFieldValue(FieldName, "User"+uuid.NewString())
			instance.SetFieldValue(FieldSurname, "Surname"+uuid.NewString())

			if err := instance.Encode(); err != nil {
				t.Fatal(err)
			}
			if err := d.DataWrite(instance.GetBufferData()); err != nil {
				t.Fatalf("DataWrite failed: %v", err)
			}
			instance.BufferReset()

			expectedEntities[instance.GetUuid()] = instance
		}

		for i := 0; i < 500; i++ {
			var entityUUID uuid.UUID
			for key := range expectedEntities {
				entityUUID = key
				break
			}

			// Create a newer version
			updatedInstance := e.EntityExtension.New()
			updatedInstance.SetFieldValue(FieldUuid, entityUUID) // Same UUID
			updatedInstance.SetFieldValue(FieldName, "UpdatedUser")
			updatedInstance.SetFieldValue(FieldSurname, "UpdatedSurname")

			if err := updatedInstance.Encode(); err != nil {
				t.Fatal(err)
			}
			if err := d.DataWrite(updatedInstance.GetBufferData()); err != nil {
				t.Fatalf("DataWrite failed: %v", err)
			}
			updatedInstance.BufferReset()

			// Replace in expected map
			expectedEntities[entityUUID] = updatedInstance
		}

		// Perform cleanup (should remove old versions)
		err := d.DataCleanup()
		if err != nil {
			t.Fatalf("DataCleanup failed: %v", err)
		}

		// Read entities after cleanup
		entities, err := d.DataReadAll()
		if err != nil {
			t.Fatalf("DataReadAll failed: %v", err)
		}

		// Ensure only 10 entities remain (latest versions)
		if len(entities) != len(expectedEntities) {
			t.Fatalf("Mismatch in entity count after cleanup: expected %d, got %d", len(expectedEntities), len(register.Entities))
		}

		// Validate updated entities
		for _, entity := range entities {
			entityUUID := entity.GetUuid()
			expected := expectedEntities[entityUUID]

			if entity.GetFieldValue(FieldName) != expected.GetFieldValue(FieldName) {
				t.Fatalf("Mismatch in Name for UUID %v", entityUUID)
			}
			if entity.GetFieldValue(FieldSurname) != expected.GetFieldValue(FieldSurname) {
				t.Fatalf("Mismatch in Surname for UUID %v", entityUUID)
			}
		}
	}
}

func BenchmarkDataWrite(b *testing.B) {
	rowCounts := []int{1000, 10000, 100000, 1000000}

	for _, numRows := range rowCounts {
		b.Run(fmt.Sprintf("%d_Rows", numRows), func(b *testing.B) {
			d := NewDisk()
			d.WithNewRandomPath()
			defer util.FileDelete(d.Path)

			if len(register.Entities) == 0 {
				b.Fatal("No entities generated")
			}

			var testEntity *register.EntityExtension
			for _, e := range register.Entities {
				if e.EntityBase.Name == "Sample" {
					testEntity = e.EntityExtension
					break
				}
			}

			if testEntity == nil {
				b.Fatal("Test entity 'Sample' not found")
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				for j := 0; j < numRows; j++ {
					instance := testEntity.New()
					instance.SetFieldValue(FieldUuid, uuid.New())
					instance.SetFieldValue(FieldName, "User"+uuid.NewString())
					instance.SetFieldValue(FieldSurname, "Surname"+uuid.NewString())

					if err := instance.Encode(); err != nil {
						b.Fatal(err)
					}
					if err := d.DataWrite(instance.GetBufferData()); err != nil {
						b.Fatalf("DataWrite failed: %v", err)
					}
					instance.BufferReset()
				}
			}
		})
	}
}

func BenchmarkDataReadAll(b *testing.B) {
	rowCounts := []int{1000, 10000, 100000, 1000000}

	for _, numRows := range rowCounts {
		b.Run(fmt.Sprintf("%d_Rows", numRows), func(b *testing.B) {
			d := NewDisk()
			d.WithNewRandomPath()
			defer util.FileDelete(d.Path)

			if len(register.Entities) == 0 {
				b.Fatal("No entities generated")
			}

			var testEntity *register.EntityExtension
			for _, e := range register.Entities {
				if e.EntityBase.Name == "Sample" {
					testEntity = e.EntityExtension
					break
				}
			}

			// Populate file with numRows entities
			for i := 0; i < numRows; i++ {
				instance := testEntity.New()
				instance.SetFieldValue(FieldUuid, uuid.New())
				instance.SetFieldValue(FieldName, "User"+uuid.NewString())
				instance.SetFieldValue(FieldSurname, "Surname"+uuid.NewString())

				if err := instance.Encode(); err != nil {
					b.Fatal(err)
				}
				if err := d.DataWrite(instance.GetBufferData()); err != nil {
					b.Fatalf("DataWrite failed: %v", err)
				}
				instance.BufferReset()
			}

			// Run benchmark
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := d.DataReadAll()
				if err != nil {
					b.Fatalf("DataReadAll failed: %v", err)
				}
			}
		})
	}
}

func TestGenerateLargeDataset(t *testing.T) {
	t.Skip()

	d := NewDisk()
	d.WithNewRandomPath()
	// t.Cleanup(func() { util.FileDelete(d.Path) })

	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	var testEntity *register.EntityExtension
	for _, e := range register.Entities {
		if e.EntityBase.Name == "Sample" {
			testEntity = e.EntityExtension
			break
		}
	}

	if testEntity == nil {
		t.Fatal("Test entity 'Sample' not found")
	}

	const numRows = 1000000

	t.Logf("Starting to write %d entities to disk...", numRows)

	// Write entities
	ticks := 0
	for i := 0; i < numRows; i++ {
		instance := testEntity.New()
		instance.SetFieldValue(FieldUuid, uuid.New())
		instance.SetFieldValue(FieldName, fmt.Sprintf("User_%d", i))
		instance.SetFieldValue(FieldSurname, fmt.Sprintf("Surname_%d", i))

		if err := instance.Encode(); err != nil {
			t.Fatal(err)
		}
		if err := d.DataWrite(instance.GetBufferData()); err != nil {
			t.Fatalf("DataWrite failed at row %d: %v", i, err)
		}
		instance.BufferReset()

		ticks++
		if ticks == 1000 {
			ticks = 0
			percentage := float64(i+1) / float64(numRows) * 100
			fmt.Printf("Progress: %.2f%%\n", percentage)
		}
	}

	fmt.Println("\nData writing completed.")
	fmt.Println("File Path:", d.Path)
}

func TestDataReadAll_WithDeletedEntity(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	d.OpenFile()
	t.Cleanup(func() { util.FileDelete(d.Path) })

	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}
		d.WithEntity(e)

		// Create and write entity
		instance := e.EntityExtension.New()
		instance.SetFieldValue(FieldUuid, uuid.New())
		instance.SetFieldValue(FieldName, "ToBeDeleted")
		instance.SetFieldValue(FieldSurname, "ShouldNotAppear")
		if err := instance.Encode(); err != nil {
			t.Fatal(err)
		}
		if err := d.DataWrite(instance.GetBufferData()); err != nil {
			t.Fatalf("Write failed: %v", err)
		}
		entityUUID := instance.GetUuid()
		instance.BufferReset()

		// Mark entity as deleted
		tombstone := e.EntityExtension.New()
		tombstone.SetFieldValue(FieldUuid, entityUUID)
		tombstone.SetFieldValue(FieldDeleted, true)
		if err := tombstone.Encode(); err != nil {
			t.Fatal(err)
		}
		if err := d.DataWrite(tombstone.GetBufferData()); err != nil {
			t.Fatalf("Delete marker write failed: %v", err)
		}
		tombstone.BufferReset()

		// Read all
		entities, err := d.DataReadAll()
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}
		if len(entities) != 0 {
			t.Fatalf("Expected 0 entities after delete, got %d", len(register.Entities))
		}
	}
}

func TestDataReadAll_DeletedThenInserted_Survives(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	d.OpenFile()
	t.Cleanup(func() { util.FileDelete(d.Path) })

	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}
		d.WithEntity(e)

		u := uuid.New()

		// Write deleted tombstone first
		del := e.EntityExtension.New()
		del.SetFieldValue(FieldUuid, u)
		del.SetFieldValue(FieldDeleted, true)
		if err := del.Encode(); err != nil {
			t.Fatal(err)
		}
		if err := d.DataWrite(del.GetBufferData()); err != nil {
			t.Fatal(err)
		}
		del.BufferReset()

		// Write a new version after deletion
		ins := e.EntityExtension.New()
		ins.SetFieldValue(FieldUuid, u)
		ins.SetFieldValue(FieldName, "Alive")
		ins.SetFieldValue(FieldSurname, "User")
		if err := ins.Encode(); err != nil {
			t.Fatal(err)
		}
		if err := d.DataWrite(ins.GetBufferData()); err != nil {
			t.Fatal(err)
		}
		ins.BufferReset()

		entities, err := d.DataReadAll()
		if err != nil {
			t.Fatal(err)
		}
		if len(register.Entities) != 1 {
			t.Fatalf("Expected 1 entity after insert following delete, got %d", len(register.Entities))
		}
		if entities[0].GetFieldValue(FieldName) != "Alive" {
			t.Fatalf("Expected Name 'Alive', got %v", entities[0].GetFieldValue(FieldName))
		}
	}
}

func TestDataCleanup_DeletesArePurged(t *testing.T) {
	d := NewDisk()
	d.WithNewRandomPath()
	d.OpenFile()
	t.Cleanup(func() { util.FileDelete(d.Path) })

	if len(register.Entities) == 0 {
		t.Fatal("No entities generated")
	}

	for _, e := range register.Entities {
		if e.EntityBase.Name != "Sample" {
			continue
		}
		d.WithEntity(e)

		u := uuid.New()

		// Insert then delete
		ins := e.EntityExtension.New()
		ins.SetFieldValue(FieldUuid, u)
		ins.SetFieldValue(FieldName, FieldDeleted)
		ins.SetFieldValue(FieldSurname, "EntityBase")
		if err := ins.Encode(); err != nil {
			t.Fatal(err)
		}
		if err := d.DataWrite(ins.GetBufferData()); err != nil {
			t.Fatal(err)
		}
		ins.BufferReset()

		del := e.EntityExtension.New()
		del.SetFieldValue(FieldUuid, u)
		del.SetFieldValue(FieldDeleted, true)
		if err := del.Encode(); err != nil {
			t.Fatal(err)
		}
		if err := d.DataWrite(del.GetBufferData()); err != nil {
			t.Fatal(err)
		}
		del.BufferReset()

		// Cleanup should purge both
		if err := d.DataCleanup(); err != nil {
			t.Fatal(err)
		}

		entities, err := d.DataReadAll()
		if err != nil {
			t.Fatal(err)
		}
		if len(entities) != 0 {
			t.Fatalf("Expected 0 entities after cleanup, got %d", len(register.Entities))
		}
	}
}
