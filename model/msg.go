package model

import (
	"github.com/rah-0/hyperion/register"
)

type MessageType int

const (
	MessageTypeUndefined MessageType = iota
	MessageTypeTest
	MessageTypeInsert
	MessageTypeDelete
	MessageTypeUpdate
	MessageTypeGetAll
)

type Mode int

const (
	ModeSync Mode = iota
	ModeAsync
)

type Status int

const (
	StatusSuccess Status = iota
)

type Message struct {
	Type   MessageType
	Mode   Mode
	Status Status
	String string
	Bytes  []byte
	Entity register.Entity
	Models []register.Model
	Error  error
}
