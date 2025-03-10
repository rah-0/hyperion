package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"

	"github.com/rah-0/hyperion/register"
	. "github.com/rah-0/hyperion/util"
)

const (
	LEN_BYTE_STATUS  = 1
	LEN_BYTES_LENGTH = 8

	STATUS_BYTE_ACTIVE  = byte(0)
	STATUS_BYTE_DELETED = byte(1)
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
	Offset uint64
}

func NewDisk() *Disk {
	return &Disk{}
}

func (x *Disk) WithNewRandomPath() {
	x.Path = filepath.Join(os.TempDir(), uuid.NewString())
}

func (x *Disk) DataWrite(data []byte) (uint64, error) {
	x.Mu.Lock()
	defer x.Mu.Unlock()
	lastOffset := x.Offset

	file, err := os.OpenFile(x.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	length := uint64(len(data) + LEN_BYTE_STATUS)
	if err = binary.Write(file, binary.LittleEndian, length); err != nil {
		return 0, err
	}
	if err = binary.Write(file, binary.LittleEndian, STATUS_BYTE_ACTIVE); err != nil {
		return 0, err
	}
	_, err = file.Write(data)
	if err != nil {
		return 0, err
	}

	x.Offset += LEN_BYTES_LENGTH + length
	return lastOffset, nil
}

func (x *Disk) DataReadAll(e *register.Entity) ([]register.Model, error) {
	x.Mu.Lock()
	defer x.Mu.Unlock()
	x.Offset = 0

	file, err := os.OpenFile(x.Path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entities []register.Model
	for {
		var length uint64
		if err = binary.Read(file, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		var status byte // Read status flag (1 byte)
		if err = binary.Read(file, binary.LittleEndian, &status); err != nil {
			return nil, err
		}
		if status == STATUS_BYTE_DELETED {
			if length > math.MaxInt64 {
				return nil, fmt.Errorf("invalid record length: %d", length)
			}
			_, err = file.Seek(int64(length-LEN_BYTE_STATUS), io.SeekCurrent) // Skip the remaining data
			if err != nil {
				return nil, err
			}
			continue
		}

		// Read the data
		data := make([]byte, length-LEN_BYTE_STATUS) // -1 since deleted flag was read already
		if _, err = file.Read(data); err != nil {
			return nil, err
		}

		instance := e.New()
		instance.SetBufferData(data)
		if err = instance.Decode(); err != nil {
			return nil, err
		}
		instance.SetOffset(x.Offset)

		x.Offset += LEN_BYTES_LENGTH + length
		entities = append(entities, instance)
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

	tempPath := x.Path + ".tmp"
	exists, err := PathExists(tempPath)
	if exists {
		// Delete old tmp file if it exists.
		// Partial writes could happen on a power loss.
		err = FileDelete(tempPath)
		if err != nil {
			return err
		}
	}

	// Create a temporary file for cleanup data
	tempFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	compactionNeeded := false
	for {
		var length uint64
		if err = binary.Read(originalFile, binary.LittleEndian, &length); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		var status byte
		if err = binary.Read(originalFile, binary.LittleEndian, &status); err != nil {
			return err
		}

		if status == STATUS_BYTE_DELETED {
			if length > math.MaxInt64 {
				return fmt.Errorf("invalid record length: %d", length)
			}
			_, err = originalFile.Seek(int64(length-LEN_BYTE_STATUS), io.SeekCurrent) // Skip remaining data
			if err != nil {
				return err
			}
			compactionNeeded = true
			continue
		}

		data := make([]byte, length-LEN_BYTE_STATUS)
		if _, err = originalFile.Read(data); err != nil {
			return err
		}

		newLength := uint64(len(data) + LEN_BYTE_STATUS)
		if err = binary.Write(tempFile, binary.LittleEndian, newLength); err != nil {
			return err
		}
		if err = binary.Write(tempFile, binary.LittleEndian, STATUS_BYTE_ACTIVE); err != nil {
			return err
		}
		if _, err = tempFile.Write(data); err != nil {
			return err
		}
	}

	if compactionNeeded {
		// Remove the old file before renaming
		if err = os.Remove(x.Path); err != nil {
			return err
		}
		// Rename the new compacted file to the original filename
		return os.Rename(tempPath, x.Path)
	} else {
		// No compaction needed, remove temporary file
		err = os.Remove(tempPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (x *Disk) DataChangeStatus(offset uint64, status byte) error {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	file, err := os.OpenFile(x.Path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Seek to status flag position
	_, err = file.Seek(int64(offset+LEN_BYTES_LENGTH), io.SeekStart)
	if err != nil {
		return err
	}

	// Write status
	if err = binary.Write(file, binary.LittleEndian, status); err != nil {
		return err
	}

	return file.Sync()
}

func (x *Disk) DataUpdate(offset uint64, data []byte) (uint64, error) {
	// Mark the existing record as deleted
	if err := x.DataChangeStatus(offset, STATUS_BYTE_DELETED); err != nil {
		return 0, err
	}
	// Write the new data as a separate entry
	return x.DataWrite(data)
}
