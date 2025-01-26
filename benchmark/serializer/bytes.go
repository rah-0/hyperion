package serializer

import (
	"bytes"
	"encoding/gob"
)

var (
	bytesBuf = &bytes.Buffer{}
	bytesEnc = gob.NewEncoder(bytesBuf)
	bytesDec = gob.NewDecoder(bytesBuf)
)

func EncodeBytes(enc *gob.Encoder, x any) error {
	if err := enc.Encode(x); err != nil {
		return err
	}
	return nil
}

func DecodeBytes(dec *gob.Decoder, x any) error {
	return dec.Decode(x)
}
