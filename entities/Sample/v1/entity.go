package Sample

// The code in this file is autogenerated, do not modify manually!

import (
	"encoding/gob"
)

func init() {
	gob.Register(Sample{})
}

const (
	Version = "v1"
)

type Sample struct {
	Name    string
	Surname string
}
