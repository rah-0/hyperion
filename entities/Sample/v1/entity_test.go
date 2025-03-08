package Sample

import (
	"testing"
)

func TestSample_EncodeDecode(t *testing.T) {
	Buffer.Reset()

	original := Sample{Name: "John", Surname: "Doe"}
	err := original.Encode()
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	var decoded Sample
	err = decoded.Decode()
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	if original.Name != decoded.Name || original.Surname != decoded.Surname {
		t.Fatalf("Decoded struct does not match original. Got: %+v, Expected: %+v", decoded, original)
	}
}
