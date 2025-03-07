package main

import (
	"encoding/binary"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Disk struct {
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
	if err := x.Serializer.Encode(data); err != nil {
		return err
	}
	serializedData := x.Serializer.GetData()

	// Write length prefix
	length := uint32(len(serializedData))
	if err := binary.Write(file, binary.LittleEndian, length); err != nil {
		return err
	}

	// Write actual data
	_, err = file.Write(serializedData)
	return err
}

func (x *Disk) ReadFromFile(data any) error {
	file, err := os.Open(x.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read length prefix
	var length uint32
	if err := binary.Read(file, binary.LittleEndian, &length); err != nil {
		return err
	}

	// Read actual serialized data
	x.Serializer.Reset()
	buf := make([]byte, length)
	if _, err := file.Read(buf); err != nil {
		return err
	}

	// Deserialize
	x.Serializer.Reset()
	x.Serializer.SetData(buf)
	return x.Serializer.Decode(data)
}
