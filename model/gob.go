package model

import (
	"encoding/gob"

	"github.com/google/uuid"
)

func init() {
	gob.Register(uuid.UUID{})
}
