package data_format

import (
	"bytes"
	"encoding/json"
)

func EncodeJson(x any) ([]byte, error) {
	b := &bytes.Buffer{}
	e := json.NewEncoder(b)
	if err := e.Encode(x); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func DecodeJson(data []byte, x any) error {
	return json.Unmarshal(data, x)
}
