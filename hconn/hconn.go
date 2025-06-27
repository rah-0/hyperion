package hconn

import (
	"encoding/binary"
	"net"
	"sync"
	"time"

	"github.com/rah-0/nabu"

	"github.com/rah-0/hyperion/model"
)

var Timeout = 120 * time.Second

type HConn struct {
	C  net.Conn
	S  *Serializer
	mu sync.Mutex
}

func NewHConn(conn net.Conn) *HConn {
	return &HConn{
		C: conn,
		S: NewSerializer(),
	}
}

func (hc *HConn) Close() error {
	return hc.C.Close()
}

// Send sends a message with a length-prefixed format
func (hc *HConn) Send(a any) error {
	if err := hc.S.Encode(a); err != nil {
		return nabu.FromError(err).Log()
	}

	data := hc.S.GetData()
	hc.S.Reset()

	dataLen := uint64(len(data))
	lengthPrefix := make([]byte, 8)
	binary.BigEndian.PutUint64(lengthPrefix, dataLen)

	if err := hc.write(lengthPrefix); err != nil {
		return err
	}
	if err := hc.write(data); err != nil {
		return err
	}

	return nil
}

// Receive reads a message using the length-prefixed format
func (hc *HConn) Receive() (msg model.Message, err error) {
	if err = hc.C.SetReadDeadline(time.Now().Add(Timeout)); err != nil {
		return
	}

	lengthPrefix, err := hc.read(8)
	if err != nil {
		return
	}

	messageLength := binary.BigEndian.Uint64(lengthPrefix)
	if messageLength == 0 {
		err = model.ErrMessageEmpty
		return
	}

	message, err := hc.read(int(messageLength))
	if err != nil {
		return
	}

	hc.S.SetData(message)
	err = hc.S.Decode(&msg)
	hc.S.Reset()
	return msg, err
}

func (hc *HConn) SendReceive(msg model.Message) (model.Message, error) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	if err := hc.Send(msg); err != nil {
		return model.Message{}, err
	}
	return hc.Receive()
}

// Ensures all bytes are sent
func (hc *HConn) write(data []byte) error {
	totalSent := 0
	for totalSent < len(data) {
		n, err := hc.C.Write(data[totalSent:])
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
		n, err := hc.C.Read(buffer[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += n
	}

	return buffer, nil
}
