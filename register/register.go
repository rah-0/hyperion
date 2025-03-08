package register

import (
	"bytes"
)

var Entities []*Entity

type Entity struct {
	// These fields are used for dynamic loading
	Version    string
	Name       string
	DbFileName string
	Fields     map[string]int
	New        func() Model
}

func RegisterEntity(entity *Entity) {
	Entities = append(Entities, entity)
}

type Model interface {
	SetFieldValue(fieldName string, value any)
	GetFieldValue(fieldName string) any
	Encode() error
	Decode() error
	BufferReset()
	GetBuffer() *bytes.Buffer
	GetBufferData() []byte
	SetBufferData([]byte)
}
