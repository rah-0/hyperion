package model

import (
	"encoding/gob"
	"time"

	"github.com/google/uuid"

	"github.com/rah-0/hyperion/util"
)

var (
	// GlobalStructFields are the fields that are added to all Entities
	// Uuid and Deleted cannot be removed, only [Tag] can be modified
	GlobalStructFields = []util.StructField{
		{
			Name: "Uuid",
			Type: "uuid.UUID",
			Tag:  "`json:\",omitzero\"`",
		}, {
			Name: "Deleted",
			Type: "bool",
			Tag:  "`json:\",omitzero\"`",
		},
	}

	// GlobalIndexesTypesSortable represents the list of types where an index can be maintained in a sorted way
	// If you have a custom type, it can be defined here, but you will have to implement your own sorting mechanism
	// This list affects Entity template generation, so if you add a custom type make sure you modify the template as well
	GlobalIndexesTypesSortable = []string{
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
		"string", "time.Time",
	}
)

// Custom types should be registered for GOB here:
func init() {
	gob.Register(time.Time{})
	gob.Register(uuid.UUID{})
}
