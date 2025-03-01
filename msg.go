package main

type MessageType int

const (
	MessageTypeUndefined MessageType = iota
	MessageTypeTest
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
}
