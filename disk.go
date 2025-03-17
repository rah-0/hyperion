package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"

	"github.com/rah-0/hyperion/register"
)

/*
Disk will structure the bytes of an entity (row) in the following way:
- Length: 8 Bytes
- Status: 1 Byte
- Data: any length

Notes:
- the length includes the Status 1 byte. So if Data is 10 bytes,
its final length will be 1 (Status byte) + 10 (Data bytes) = 11 Bytes
*/
type Disk struct {
	Mu     sync.Mutex
	Path   string
	Entity *register.Entity
}

func NewDisk() *Disk {
	return &Disk{}
}

func (x *Disk) WithNewRandomPath() *Disk {
	x.Path = filepath.Join(os.TempDir(), uuid.NewString())
	return x
}

func (x *Disk) WithEntity(e *register.Entity) *Disk {
	x.Entity = e
	return x
}

func (x *Disk) WithPath(path string) *Disk {
	x.Path = path
	return x
}

func (x *Disk) DataWrite(data []byte) error {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	file, err := os.OpenFile(x.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	length := uint64(len(data))
	if err = binary.Write(file, binary.LittleEndian, length); err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return file.Sync()
}

func (x *Disk) DataReadAll() ([]register.Model, error) {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	file, err := os.OpenFile(x.Path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	latestEntities := make(map[uuid.UUID]register.Model)

	for {
		var length uint64
		if err = binary.Read(file, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		data := make([]byte, length)
		if _, err = io.ReadFull(file, data); err != nil {
			return nil, fmt.Errorf("failed to read full record data: %w", err)
		}

		// Decode entity
		instance := x.Entity.New()
		instance.SetBufferData(data)
		if err = instance.Decode(); err != nil {
			return nil, fmt.Errorf("failed to decode entity: %w", err)
		}

		// Extract UUID and keep only the latest version
		entityUUID := instance.GetUuid()
		latestEntities[entityUUID] = instance
	}

	// Convert map to slice
	entities := make([]register.Model, 0, len(latestEntities))
	for _, model := range latestEntities {
		entities = append(entities, model)
	}

	return entities, nil
}

func (x *Disk) DataCleanup() error {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	originalFile, err := os.OpenFile(x.Path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	// First pass: Detect duplicates
	seenUUIDs := make(map[uuid.UUID]bool)
	hasDuplicates := false

	for {
		var length uint64
		if err = binary.Read(originalFile, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		data := make([]byte, length)
		if _, err = originalFile.Read(data); err != nil {
			return err
		}

		// Decode entity to extract UUID
		instance := x.Entity.New()
		instance.SetBufferData(data)
		if err = instance.Decode(); err != nil {
			return err
		}

		entityUUID := instance.GetUuid()

		// Check if UUID has already been seen
		if seenUUIDs[entityUUID] {
			hasDuplicates = true
			break // Stop scanning early; we now know cleanup is needed
		}
		seenUUIDs[entityUUID] = true
	}

	// If no duplicates were found, no need to clean up
	if !hasDuplicates {
		return nil
	}

	// Reset file pointer for second pass
	if _, err = originalFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// Second pass: Compact the file by keeping only the latest version per UUID
	tempPath := x.Path + ".tmp"
	latestEntities := make(map[uuid.UUID][]byte)

	for {
		var length uint64
		if err = binary.Read(originalFile, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		data := make([]byte, length)
		if _, err = originalFile.Read(data); err != nil {
			return err
		}

		// Decode entity to extract UUID
		instance := x.Entity.New()
		instance.SetBufferData(data)
		if err = instance.Decode(); err != nil {
			return err
		}

		entityUUID := instance.GetUuid()

		// Store only the latest version
		latestEntities[entityUUID] = data
	}

	// Create new compacted file
	tempFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// Write the latest versions to the compacted file
	for _, data := range latestEntities {
		length := uint64(len(data))
		if err = binary.Write(tempFile, binary.LittleEndian, length); err != nil {
			return err
		}
		if _, err = tempFile.Write(data); err != nil {
			return err
		}
	}

	// Ensure data is fully written
	if err = tempFile.Sync(); err != nil {
		return err
	}

	// Replace old file with compacted file
	if err = os.Rename(tempPath, x.Path); err != nil {
		return err
	}

	return nil
}
