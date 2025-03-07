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
	serializer := NewSerializer()

	if len(Entities) == 0 {
		t.Fatal("no entities generated")
	}

	for _, e := range Entities {
		if e.EntityName != "Sample" {
			continue
		}

		instance := e.New()
		instance.SetFieldValue("Name", "John")
		instance.SetFieldValue("Surname", "Doe")

		err := serializer.Encode(instance)
		if err != nil {
			t.Fatalf("Encoding failed for %s: %v", e.EntityName, err)
		}

		data := serializer.GetData()
		serializer.Reset()

		newInstance := e.New()
		serializer.SetData(data)
		err = serializer.Decode(newInstance)
		if err != nil {
			t.Fatalf("Decoding failed for %s: %v", e.EntityName, err)
		}

		if newInstance.GetFieldValue("Name") != "John" || newInstance.GetFieldValue("Surname") != "Doe" {
			t.Fatalf("Decoded values mismatch for %s: got Name=%s, Surname=%s",
				e.EntityName, newInstance.GetFieldValue("Name"), newInstance.GetFieldValue("Surname"))
		}
	}
}
