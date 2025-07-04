package disk

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
	"github.com/rah-0/hyperion/util"
)

/*
Disk will structure the bytes of an entity (row) in the following way:
- Length: 8 Bytes
- Data: any length
*/
type Disk struct {
	Mu         sync.Mutex
	File       *os.File
	Path       string
	Entity     *register.Entity
	SyncTicker *time.Ticker
	StopChan   chan struct{}
	IsClosed   bool
	Wg         sync.WaitGroup
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

// OpenFile sets up the backing file for disk operations and starts a background sync routine.
func (x *Disk) OpenFile() error {
	x.Mu.Lock()

	if x.File != nil {
		x.Mu.Unlock()
		return nil
	}

	file, err := os.OpenFile(x.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		x.Mu.Unlock()
		return err
	}
	x.File = file
	x.IsClosed = false
	x.StopChan = make(chan struct{})
	x.SyncTicker = time.NewTicker(1 * time.Second)

	// Create local copies of resources needed by background goroutine
	// to avoid potential race conditions during shutdown
	syncTicker := x.SyncTicker
	stopChan := x.StopChan

	x.Wg.Add(1)

	// Release lock before starting goroutine to avoid deadlock
	x.Mu.Unlock()

	// Start background sync with copies of needed resources
	go func() {
		defer x.Wg.Done()
		x.runBackgroundSync(syncTicker, stopChan)
	}()

	return nil
}

// runBackgroundSync safely performs periodic sync operations on the file
func (x *Disk) runBackgroundSync(ticker *time.Ticker, stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			x.Mu.Lock()
			file := x.File
			isClosed := x.IsClosed
			x.Mu.Unlock()

			if isClosed || file == nil {
				return
			}
			err := file.Sync()
			if err != nil {
				nabu.FromError(err).WithMessage("failed to sync file").Log()
			}
		}
	}
}

// Close safely closes the Disk, ensuring all data is synced to disk and resources are released.
// Close stops the background sync goroutine, performs a final sync operation, and closes the file handle.
// This method is safe to call multiple times and will only perform cleanup actions once.
func (x *Disk) Close() error {
	x.Mu.Lock()
	if x.IsClosed {
		x.Mu.Unlock()
		return nil
	}
	x.IsClosed = true

	if x.SyncTicker != nil {
		x.SyncTicker.Stop()
		x.SyncTicker = nil
	}
	if x.StopChan != nil {
		close(x.StopChan)
		x.StopChan = nil
	}

	// Store reference to file before unlocking and nullify
	file := x.File
	x.File = nil
	x.Mu.Unlock()

	// Wait for background goroutines to complete
	x.Wg.Wait()

	// Close the file if it was open
	if file != nil {
		if err := file.Sync(); err != nil {
			return err
		}
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
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

// DataReadAll reads all entities from disk, returning only the latest version of each entity.
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
		instance := x.Entity.EntityExtension.New()
		instance.SetBufferData(data)
		if err = instance.Decode(); err != nil {
			return nil, err
		}

		// Extract UUID and keep only the latest version
		entityUUID := instance.GetUuid()
		if instance.IsDeleted() {
			delete(latestEntities, entityUUID)
			continue
		}

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

	fileInfo, err := originalFile.Stat()
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	// First pass: detect if any duplicate UUIDs or deletes exist
	seenUUIDs := make(map[uuid.UUID]struct{})
	hasDuplicates := false

	var bytesRead int64 = 0
	var lastLoggedProgress = -1.0

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
		if _, err = io.ReadFull(originalFile, data); err != nil {
			return err
		}
		bytesRead += int64(len(data))

		instance := x.Entity.EntityExtension.New()
		instance.SetBufferData(data)
		if err = instance.Decode(); err != nil {
			return err
		}

		id := instance.GetUuid()
		if _, exists := seenUUIDs[id]; exists || instance.IsDeleted() {
			hasDuplicates = true
			break
		}
		seenUUIDs[id] = struct{}{}

		progress := (float64(bytesRead) / float64(fileSize)) * 100
		if progress-lastLoggedProgress >= 1.0 {
			nabu.FromMessage(fmt.Sprintf("Progress: %.2f%% scanned", progress)).Log()
			lastLoggedProgress = progress
		}
	}

	if !hasDuplicates {
		nabu.FromMessage("No duplicates or deletions found. Skipping cleanup.").Log()
		return nil
	}

	// Second pass: compact file by keeping latest valid (non-deleted) version per UUID
	if _, err = originalFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	type entityInfo struct {
		data    []byte
		deleted bool
	}
	latest := make(map[uuid.UUID]entityInfo)

	bytesRead = 0
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
		if _, err = io.ReadFull(originalFile, data); err != nil {
			return err
		}
		bytesRead += int64(len(data))

		instance := x.Entity.EntityExtension.New()
		instance.SetBufferData(data)
		if err = instance.Decode(); err != nil {
			return err
		}

		latest[instance.GetUuid()] = entityInfo{
			data:    data,
			deleted: instance.IsDeleted(),
		}

		progress := (float64(bytesRead) / float64(fileSize)) * 100
		if progress-lastLoggedProgress >= 1.0 {
			nabu.FromMessage(fmt.Sprintf("Compacting progress: %.2f%%", progress)).Log()
			lastLoggedProgress = progress
		}
	}

	// Write compacted file
	tempPath := x.Path + ".tmp"
	if exists, err := util.PathExists(tempPath); err != nil {
		return err
	} else if exists {
		if err = util.FileDelete(tempPath); err != nil {
			return err
		}
	}

	tempFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	var totalBytes, writtenBytes int64
	for _, info := range latest {
		if info.deleted {
			continue
		}
		totalBytes += int64(binary.Size(uint64(len(info.data))) + len(info.data))
	}

	lastLoggedProgress = -1.0
	for _, info := range latest {
		if info.deleted {
			continue
		}
		length := uint64(len(info.data))
		if err = binary.Write(tempFile, binary.LittleEndian, length); err != nil {
			return err
		}
		if _, err = tempFile.Write(info.data); err != nil {
			return err
		}
		writtenBytes += int64(binary.Size(length) + len(info.data))

		progress := (float64(writtenBytes) / float64(totalBytes)) * 100
		if progress-lastLoggedProgress >= 1.0 {
			nabu.FromMessage(fmt.Sprintf("Writing progress: %.2f%%", progress)).Log()
			lastLoggedProgress = progress
		}
	}

	if err = tempFile.Sync(); err != nil {
		return err
	}
	if err = os.Rename(tempPath, x.Path); err != nil {
		return err
	}

	nabu.FromMessage("Data cleanup completed successfully for file: [" + x.Path + "]").Log()
	return nil
}
