package data_format

import (
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/hyperion/utils/testutil"
)

func TestEncodeDecodeBytes(t *testing.T) {
	gob.Register(Person{})

	original := &Person{
		Name:    "John",
		Age:     30,
		Surname: "Doe",
	}

	data, err := EncodeBytes(original)
	assert.NoError(t, err, "GOB serialization should not fail")
	assert.NotEmpty(t, data, "Serialized data should not be empty")

	decoded := &Person{}
	err = DecodeBytes(data, decoded)
	assert.NoError(t, err, "GOB deserialization should not fail")

	// Assert that the original and decoded structs match
	assert.Equal(t, original.Name, decoded.Name, "Names should match")
	assert.Equal(t, original.Age, decoded.Age, "Ages should match")
	assert.Equal(t, original.Surname, decoded.Surname, "Surnames should match")
}

func BenchmarkEncodeBytes1Small(b *testing.B) {
	runEncodeBytesBenchmark(b, 1)
}

func BenchmarkEncodeBytes100Small(b *testing.B) {
	runEncodeBytesBenchmark(b, 100)
}

func BenchmarkEncodeBytes10000Small(b *testing.B) {
	runEncodeBytesBenchmark(b, 10000)
}

func BenchmarkEncodeBytes1000000Small(b *testing.B) {
	runEncodeBytesBenchmark(b, 1000000)
}

func runEncodeBytesBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := generateRandomPersons(count)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EncodeBytes(data)
		if err != nil {
			b.Fatalf("Failed to Encode Gob: %v", err)
		}
	}
}

func BenchmarkEncodeBytes1Unreal(b *testing.B) {
	runEncodeBytesUnrealBenchmark(b, 1)
}

func BenchmarkEncodeBytes10Unreals(b *testing.B) {
	runEncodeBytesUnrealBenchmark(b, 10)
}

func BenchmarkEncodeBytes100Unreals(b *testing.B) {
	runEncodeBytesUnrealBenchmark(b, 100)
}

func BenchmarkEncodeBytes1000Unreals(b *testing.B) {
	runEncodeBytesUnrealBenchmark(b, 1000)
}

func runEncodeBytesUnrealBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := generateRandomUnreals(count)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EncodeBytes(data)
		if err != nil {
			b.Fatalf("Failed to Encode Gob: %v", err)
		}
	}
}
