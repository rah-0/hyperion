package main

import (
	"bytes"
	"encoding/gob"
)

type Serializer struct {
	b *bytes.Buffer
	e *gob.Encoder
	d *gob.Decoder
}

func NewSerializer() *Serializer {
	b := new(bytes.Buffer)
	return &Serializer{
		b: b,
		e: gob.NewEncoder(b),
		d: gob.NewDecoder(b),
	}
}

func (x *Serializer) Encode(a any) error {
	if err := x.e.Encode(a); err != nil {
		return err
	}
	return nil
}

func (x *Serializer) Decode(a any) error {
	return x.d.Decode(a)
}

func (x *Serializer) Reset() {
	x.b.Reset()
}

func (x *Serializer) SetData(d []byte) {
	x.b.Write(d)
}

func (x *Serializer) GetData() []byte {
	return x.b.Bytes()
}

func (x *Serializer) GetBuffer() *bytes.Buffer {
	return x.b
}
