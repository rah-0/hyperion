package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rah-0/nabu"

	"github.com/rah-0/hyperion/register"
	. "github.com/rah-0/hyperion/util"
)

/*
Disk will structure the bytes of an entity (row) in the following way:
- Length: 8 Bytes
- Data: any length
*/
type Disk struct {
	Mu         sync.Mutex
	Path       string
	Entity     *register.Entity
	File       *os.File
	syncTicker *time.Ticker
	stopChan   chan struct{}
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

func (x *Disk) OpenFile() error {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	file, err := os.OpenFile(x.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	x.File = file

	// Start background sync loop
	x.syncTicker = time.NewTicker(1 * time.Second)
	x.stopChan = make(chan struct{})
	go x.backgroundSync()

	return nil
}

func (x *Disk) backgroundSync() {
	for {
		select {
		case <-x.syncTicker.C:
			x.Mu.Lock()
			x.File.Sync()
			x.Mu.Unlock()
		case <-x.stopChan:
			return
		}
	}
}

func (x *Disk) DataWrite(data []byte) error {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	length := uint64(len(data))
	if err := binary.Write(x.File, binary.LittleEndian, length); err != nil {
		return err
	}

	_, err := x.File.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (x *Disk) DataReadAll() ([]register.Model, error) {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	if _, err := x.File.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	nabu.FromMessage("Starting data read for file: [" + x.Path + "]").Log()

	// Get total file size
	fileInfo, err := x.File.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()
	var bytesRead int64 = 0
	var lastLoggedProgress = -1.0

	latestEntities := make(map[uuid.UUID]register.Model)

	nabu.FromMessage("Reading entities from file...").Log()
	for {
		var length uint64
		if err = binary.Read(x.File, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		bytesRead += int64(binary.Size(length))

		data := make([]byte, length)
		if _, err = io.ReadFull(x.File, data); err != nil {
			return nil, err
		}
		bytesRead += int64(len(data))

		// Decode entity
		instance := x.Entity.New()
		instance.SetBufferData(data)
		if err = instance.Decode(); err != nil {
			return nil, err
		}

		// Extract UUID and keep only the latest version
		entityUUID := instance.GetUuid()
		latestEntities[entityUUID] = instance

		// Log progress every 1% interval, with two decimal places
		progress := (float64(bytesRead) / float64(fileSize)) * 100
		if progress-lastLoggedProgress >= 1.0 {
			nabu.FromMessage(fmt.Sprintf("Reading progress: %.2f%% completed", progress)).Log()
			lastLoggedProgress = progress
		}
	}

	// Convert map to slice
	entities := make([]register.Model, 0, len(latestEntities))
	for _, model := range latestEntities {
		entities = append(entities, model)
	}

	nabu.FromMessage("Finished reading all entities from file. Total entities: " + fmt.Sprintf("%d", len(entities))).Log()
	return entities, nil
}

func (x *Disk) DataCleanup() error {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	if _, err := x.File.Seek(0, io.SeekStart); err != nil {
		return err
	}

	nabu.FromMessage("Starting data cleanup for file: [" + x.Path + "]").Log()
	originalFile, err := os.OpenFile(x.Path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	// First pass: Detect duplicates
	seenUUIDs := make(map[uuid.UUID]bool)
	hasDuplicates := false

	fileInfo, err := originalFile.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()
	var bytesRead int64 = 0
	var lastLoggedProgress float64 = -1.0

	nabu.FromMessage("Scanning file to detect duplicates...").Log()
	for {
		var length uint64
		if err = binary.Read(originalFile, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		bytesRead += int64(binary.Size(length))

		data := make([]byte, length)
		if _, err = originalFile.Read(data); err != nil {
			return err
		}
		bytesRead += int64(len(data))

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
			nabu.FromMessage("Duplicate detected for entity UUID: [" + entityUUID.String() + "]").Log()
			break // Stop scanning early; we now know cleanup is needed
		}
		seenUUIDs[entityUUID] = true

		// Log progress every 1% interval, with two decimal places
		progress := (float64(bytesRead) / float64(fileSize)) * 100
		if progress-lastLoggedProgress >= 1.0 {
			nabu.FromMessage(fmt.Sprintf("Progress: %.2f%% scanned", progress)).Log()
			lastLoggedProgress = progress
		}
	}

	// If no duplicates were found, no need to clean up
	if !hasDuplicates {
		nabu.FromMessage("No duplicates found. Skipping cleanup.").Log()
		return nil
	}

	// Reset file pointer for second pass
	if _, err = originalFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// Second pass: Compact the file by keeping only the latest version per UUID
	tempPath := x.Path + ".tmp"
	exists, err := PathExists(tempPath)
	if err != nil {
		return err
	}
	if exists {
		err = FileDelete(tempPath)
		if err != nil {
			return err
		}
	}

	latestEntities := make(map[uuid.UUID][]byte)

	nabu.FromMessage("Compacting file...").Log()
	bytesRead = 0 // Reset counter
	lastLoggedProgress = -1.0

	for {
		var length uint64
		if err = binary.Read(originalFile, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		bytesRead += int64(binary.Size(length))

		data := make([]byte, length)
		if _, err = originalFile.Read(data); err != nil {
			return err
		}
		bytesRead += int64(len(data))

		// Decode entity to extract UUID
		instance := x.Entity.New()
		instance.SetBufferData(data)
		if err = instance.Decode(); err != nil {
			return err
		}

		entityUUID := instance.GetUuid()
		latestEntities[entityUUID] = data

		// Log progress every 1% interval
		progress := (float64(bytesRead) / float64(fileSize)) * 100
		if progress-lastLoggedProgress >= 1.0 {
			nabu.FromMessage(fmt.Sprintf("Compacting progress: %.2f%% completed", progress)).Log()
			lastLoggedProgress = progress
		}
	}

	// Create new compacted file
	tempFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	nabu.FromMessage("Writing compacted data to new file...").Log()
	bytesWritten := int64(0)
	lastLoggedProgress = -1.0
	totalBytes := int64(0)

	for _, data := range latestEntities {
		totalBytes += int64(len(data) + binary.Size(uint64(len(data))))
	}

	for _, data := range latestEntities {
		length := uint64(len(data))
		if err = binary.Write(tempFile, binary.LittleEndian, length); err != nil {
			return err
		}
		if _, err = tempFile.Write(data); err != nil {
			return err
		}
		bytesWritten += int64(len(data) + binary.Size(length))

		// Log progress every 1% interval
		progress := (float64(bytesWritten) / float64(totalBytes)) * 100
		if progress-lastLoggedProgress >= 1.0 {
			nabu.FromMessage(fmt.Sprintf("Writing progress: %.2f%% completed", progress)).Log()
			lastLoggedProgress = progress
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

	nabu.FromMessage("Data cleanup completed successfully for file: [" + x.Path + "]").Log()
	return nil
}
