package entities

import (
	"time"
)

/*
	Define your entities (structs) in this file.
*/

type Sample struct {
	Name    string    `json:"-"`
	Surname string    `json:"-"`
	Birth   time.Time `json:"-"`
	//FullName string
}
