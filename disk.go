package main

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"

	"github.com/rah-0/hyperion/register"
)

type Disk struct {
	Mu   sync.Mutex
	Path string
}

func NewDisk() *Disk {
	return &Disk{}
}

func (x *Disk) WithNewRandomPath() {
	x.Path = filepath.Join(os.TempDir(), uuid.NewString())
}

func (x *Disk) WriteToFile(m register.Model) error {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	file, err := os.OpenFile(x.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = m.Encode(); err != nil {
		return err
	}

	length := uint64(m.GetBuffer().Len())
	if err = binary.Write(file, binary.LittleEndian, length); err != nil {
		return err
	}

	_, err = file.Write(m.GetBufferData())
	return err
}

func (x *Disk) LoadAllEntities(e *register.Entity) ([]register.Model, error) {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	file, err := os.OpenFile(x.Path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entities []register.Model
	var offset uint64
	for {
		var length uint64
		if err = binary.Read(file, binary.LittleEndian, &length); err != nil {
			break // EOF or error
		}

		data := make([]byte, length)
		if _, err = file.Read(data); err != nil {
			return nil, err
		}

		instance := e.New()
		instance.SetBufferData(data)
		if err = instance.Decode(); err != nil {
			return nil, err
		}

		instance.SetOffset(offset)
		offset += 8 + length

		entities = append(entities, instance)
	}
	return entities, nil
}
