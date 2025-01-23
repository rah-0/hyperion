package data_format

import (
	"bytes"
	"encoding/gob"
)

func EncodeBytes(x any) ([]byte, error) {
	b := &bytes.Buffer{}
	e := gob.NewEncoder(b)
	if err := e.Encode(x); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func DecodeBytes(data []byte, x any) error {
	b := &bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(b)
	return d.Decode(x)
}
