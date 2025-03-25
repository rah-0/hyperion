package SampleV1

import (
	"testing"

	"github.com/rah-0/hyperion/register"
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

func TestGeneratedSampleSerializer(t *testing.T) {
	if len(register.Entities) == 0 {
		t.Fatal("no entities generated")
	}
	for _, e := range register.Entities {
		if e.Name != "Sample" {
			continue
		}

		instanceOld := e.New()
		instanceOld.SetFieldValue(FieldName, "John")
		instanceOld.SetFieldValue(FieldSurname, "Doe")
		err := instanceOld.Encode()
		if err != nil {
			t.Fatalf("Encoding failed for %s: %v", e.Name, err)
		}

		instanceNew := e.New()
		err = instanceNew.Decode()
		if err != nil {
			t.Fatalf("Decoding failed for %s: %v", e.Name, err)
		}
		instanceNew.BufferReset()

		if instanceNew.GetFieldValue(FieldName) != "John" || instanceNew.GetFieldValue(FieldSurname) != "Doe" {
			t.Fatalf("Decoded values mismatch for %s: got Name=%s, Surname=%s",
				e.Name, instanceNew.GetFieldValue(FieldName), instanceNew.GetFieldValue(FieldSurname))
		}
	}
}
