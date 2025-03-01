package main

import (
	"encoding/binary"
	"net"

	"github.com/rah-0/nabu"
)

type HConn struct {
	c net.Conn
	s *Serializer
}

func NewHConn(conn net.Conn) *HConn {
	return &HConn{
		c: conn,
		s: NewSerializer(),
	}
}

// Send sends a message with a length-prefixed format over the connection
func (hc *HConn) Send(a any) error {
	if err := hc.s.Encode(a); err != nil {
		return nabu.FromError(err).Log()
	}

	data := hc.s.GetData()
	hc.s.Reset()

	dataLen := uint32(len(data))
	lengthPrefix := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthPrefix, dataLen)

	if err := hc.write(lengthPrefix); err != nil {
		return err
	}
	if err := hc.write(data); err != nil {
		return err
	}

	return nil
}

// Receive reads a message using the length-prefixed format
func (hc *HConn) Receive() (msg Message, err error) {
	lengthPrefix, err := hc.read(4)
	if err != nil {
		return
	}

	messageLength := binary.BigEndian.Uint32(lengthPrefix)
	if messageLength == 0 {
		err = ErrMessageEmpty
		return
	}

	message, err := hc.read(int(messageLength))
	if err != nil {
		return
	}

	hc.s.SetData(message)
	err = hc.s.Decode(&msg)
	hc.s.Reset()
	return msg, err
}

// Ensures all bytes are sent
func (hc *HConn) write(data []byte) error {
	totalSent := 0
	for totalSent < len(data) {
		n, err := hc.c.Write(data[totalSent:])
		if err != nil {
			return err
		}
		totalSent += n
	}
	return nil
}

// Reads exactly `size` bytes from the connection
func (hc *HConn) read(size int) ([]byte, error) {
	buffer := make([]byte, size)
	totalRead := 0

	for totalRead < size {
		n, err := hc.c.Read(buffer[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += n
	}

	return buffer, nil
}
