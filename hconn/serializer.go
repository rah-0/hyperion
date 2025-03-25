package hconn

import (
	"bytes"
	"encoding/gob"
)

type Serializer struct {
	Buffer *bytes.Buffer
	E      *gob.Encoder
	D      *gob.Decoder
}

func NewSerializer() *Serializer {
	b := new(bytes.Buffer)
	return &Serializer{
		Buffer: b,
		E:      gob.NewEncoder(b),
		D:      gob.NewDecoder(b),
	}
}

func (x *Serializer) Encode(a any) error {
	if err := x.E.Encode(a); err != nil {
		return err
	}
	return nil
}

func (x *Serializer) Decode(a any) error {
	return x.D.Decode(a)
}

func (x *Serializer) Reset() {
	x.Buffer.Reset()
}

func (x *Serializer) SetData(d []byte) {
	x.Buffer.Write(d)
}

func (x *Serializer) GetData() []byte {
	return x.Buffer.Bytes()
}

func (x *Serializer) GetBuffer() *bytes.Buffer {
	return x.Buffer
}
