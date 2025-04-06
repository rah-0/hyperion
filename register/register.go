package register

import (
	"bytes"

	"github.com/google/uuid"
)

type IndexAccessor struct {
	GetByValue func(value any) []Model
}

var Entities []*Entity

type Entity struct {
	EntityBase      *EntityBase
	EntityExtension *EntityExtension
}

type EntityBase struct {
	// These fields are used for dynamic loading
	Version    string
	Name       string
	DbFileName string
	Data       []byte
}

type EntityExtension struct {
	New            func() Model
	FieldTypes     map[int]string
	Indexes        map[int]any
	IndexAccessors map[int]IndexAccessor
}

func RegisterEntity(entity *EntityBase, entityExtension *EntityExtension) {
	Entities = append(Entities, &Entity{
		EntityBase:      entity,
		EntityExtension: entityExtension,
	})
}

type Model interface {
	GetUuid() uuid.UUID
	SetUuid(uuid uuid.UUID)
	WithNewUuid()
	IsDeleted() bool
	SetFieldValue(int, any)
	GetFieldValue(int) any

	Encode() error
	Decode() error

	BufferReset()
	GetBuffer() *bytes.Buffer
	GetBufferData() []byte
	SetBufferData([]byte)

	MemoryAdd()
	MemoryRemove()
	MemoryUpdate()
	MemoryClear()
	MemoryGetAll() []Model
	MemoryContains(Model) bool
}
