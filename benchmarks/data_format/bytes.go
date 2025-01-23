package data_format

import (
	"bytes"
	"encoding/gob"
)

var (
	bytesBuf = &bytes.Buffer{}
	bytesEnc = gob.NewEncoder(bytesBuf)
	bytesDec = gob.NewDecoder(bytesBuf)
)

func EncodeBytes(x any) ([]byte, error) {
	bytesBuf.Reset()
	if err := bytesEnc.Encode(x); err != nil {
		return nil, err
	}
	return bytesBuf.Bytes(), nil
}

func DecodeBytes(data []byte, x any) error {
	bytesBuf.Reset()
	bytesBuf.Write(data)
	return bytesDec.Decode(x)
}
