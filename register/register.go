package register

import (
	"bytes"

	"github.com/google/uuid"
)

var Entities []*Entity

type Entity struct {
	// These fields are used for dynamic loading
	Version    string
	Name       string
	DbFileName string
	Fields     map[string]int
	New        func() Model
	Data       []byte
}

func RegisterEntity(entity *Entity) {
	Entities = append(Entities, entity)
}

type Model interface {
	GetUuid() uuid.UUID
	SetUuid(uuid uuid.UUID)
	WithNewUuid()
	SetFieldValue(string, any)
	GetFieldValue(string) any

	Encode() error
	Decode() error

	BufferReset()
	GetBuffer() *bytes.Buffer
	GetBufferData() []byte
	SetBufferData([]byte)

	MemoryAdd()
	MemoryRemove() bool
	MemoryClear()
	MemoryGetAll() []Model
	MemoryContains(Model) bool
	MemorySet([]Model)
}
