package main

import (
	"testing"
)

// TestSerializerEncodeDecode tests the Encode and Decode methods.
func TestSerializerEncodeDecode(t *testing.T) {
	s := NewSerializer()

	// Define a test structure
	type TestStruct struct {
		Name  string
		Value int
	}

	original := TestStruct{Name: "Test", Value: 42}
	var decoded TestStruct

	// Encode the original struct
	if err := s.Encode(original); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode into a new struct
	if err := s.Decode(&decoded); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Ensure the decoded struct matches the original
	if original != decoded {
		t.Fatalf("Expected %+v, got %+v", original, decoded)
	}
}

// TestSerializerReuseBuffer ensures that the buffer is reusable.
func TestSerializerReuseBuffer(t *testing.T) {
	s := NewSerializer()

	type Data struct {
		ID   int
		Text string
	}

	data1 := Data{ID: 1, Text: "First"}
	data2 := Data{ID: 2, Text: "Second"}

	var decoded Data

	// Encode and Decode first object
	if err := s.Encode(data1); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if err := s.Decode(&decoded); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded != data1 {
		t.Fatalf("Expected %+v, got %+v", data1, decoded)
	}

	if err := s.Encode(data2); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if err := s.Decode(&decoded); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded != data2 {
		t.Fatalf("Expected %+v, got %+v", data2, decoded)
	}
}
