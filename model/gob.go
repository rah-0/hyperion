package model

import (
	"encoding/gob"
	"time"

	"github.com/google/uuid"
)

func init() {
	gob.Register(time.Time{})
	gob.Register(uuid.UUID{})
}
