package main

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type Disk struct {
	Mu         sync.Mutex
	Serializer *Serializer
	Path       string
}

func NewDisk() *Disk {
	return &Disk{}
}

func (x *Disk) WithNewSerializer() {
	x.Serializer = NewSerializer()
}

func (x *Disk) WithNewRandomPath() {
	x.Path = filepath.Join(os.TempDir(), uuid.NewString())
}

func (x *Disk) WriteToFile(data any) error {
	file, err := os.OpenFile(x.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	x.Serializer.Reset()
	if err = x.Serializer.Encode(data); err != nil {
		return err
	}

	length := uint32(x.Serializer.Buffer.Len())
	if err = binary.Write(file, binary.LittleEndian, length); err != nil {
		return err
	}

	_, err = file.Write(x.Serializer.GetData())
	return err
}

func (x *Disk) ReadFromFile(data any) error {
	file, err := os.OpenFile(x.Path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var length uint32
	if err = binary.Read(file, binary.LittleEndian, &length); err != nil {
		return err
	}

	x.Serializer.Reset()
	if _, err = x.Serializer.Buffer.ReadFrom(file); err != nil {
		return err
	}
	return x.Serializer.Decode(data)
}
