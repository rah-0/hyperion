package data_format

import (
	"bytes"
	"encoding/json"
)

var (
	jsonBuf     = &bytes.Buffer{}
	jsonEncoder = json.NewEncoder(jsonBuf)
)

func EncodeJson(x any) ([]byte, error) {
	jsonBuf.Reset()
	if err := jsonEncoder.Encode(x); err != nil {
		return nil, err
	}
	return jsonBuf.Bytes(), nil
}

func DecodeJson(data []byte, x any) error {
	return json.Unmarshal(data, x)
}
