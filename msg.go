package main

type MessageType int

const (
	MessageTypeUndefined MessageType = iota
	MessageTypeTest
	MessageTypeInsert = 2
)

type Mode int

const (
	ModeSync Mode = iota
	ModeAsync
)

type Message struct {
	Type   MessageType
	Mode   Mode
	String string
	Bytes  []byte
	Entity []byte
}
