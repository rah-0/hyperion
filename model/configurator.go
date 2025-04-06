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
)

// Custom types should be registered for GOB here:
func init() {
	gob.Register(time.Time{})
	gob.Register(uuid.UUID{})
}
