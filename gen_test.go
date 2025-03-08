package main

import (
	"testing"

	. "github.com/rah-0/hyperion/register"
)

func TestGenerate(t *testing.T) {
	err := Generate()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGeneratedSampleSerializer(t *testing.T) {
	if len(Entities) == 0 {
		t.Fatal("no entities generated")
	}
	for _, e := range Entities {
		if e.Name != "Sample" {
			continue
		}

		instanceOld := e.New()
		instanceOld.SetFieldValue("Name", "John")
		instanceOld.SetFieldValue("Surname", "Doe")
		err := instanceOld.Encode()
		if err != nil {
			t.Fatalf("Encoding failed for %s: %v", e.Name, err)
		}

		instanceNew := e.New()
		err = instanceNew.Decode()
		if err != nil {
			t.Fatalf("Decoding failed for %s: %v", e.Name, err)
		}

		if instanceNew.GetFieldValue("Name") != "John" || instanceNew.GetFieldValue("Surname") != "Doe" {
			t.Fatalf("Decoded values mismatch for %s: got Name=%s, Surname=%s",
				e.Name, instanceNew.GetFieldValue("Name"), instanceNew.GetFieldValue("Surname"))
		}
	}
}
