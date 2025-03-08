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

	length := uint32(m.GetBuffer().Len())
	if err = binary.Write(file, binary.LittleEndian, length); err != nil {
		return err
	}

	_, err = file.Write(m.GetBufferData())
	return err
}

func (x *Disk) ReadFromFile(m register.Model) error {
	x.Mu.Lock()
	defer x.Mu.Unlock()

	file, err := os.OpenFile(x.Path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var length uint32
	if err = binary.Read(file, binary.LittleEndian, &length); err != nil {
		return err
	}

	if _, err = m.GetBuffer().ReadFrom(file); err != nil {
		return err
	}
	return m.Decode()
}
