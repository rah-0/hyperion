package model

import (
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
