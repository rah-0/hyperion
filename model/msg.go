package model

import (
	"github.com/rah-0/hyperion/query"
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
	MessageTypeQuery
)

type Status int

const (
	StatusSuccess Status = iota
	StatusError
)

type Message struct {
	Type   MessageType
	Status Status
	String string
	Bytes  []byte
	Entity register.Entity
	Models []register.Model
	Query  *query.Query
}

func (x *Message) Error(errMsg string) *Message {
	x.Status = StatusError
	x.String = errMsg
	return x
}
