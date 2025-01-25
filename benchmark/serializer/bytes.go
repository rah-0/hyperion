package serializer

import (
	"encoding/gob"
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
