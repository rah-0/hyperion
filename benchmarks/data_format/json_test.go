package data_format

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/hyperion/utils/testutil"
)

func TestEncodeDecodeJson(t *testing.T) {
	original := &Person{
		Name:    "John",
		Age:     30,
		Surname: "Doe",
	}

	data, err := EncodeJson(original)
	assert.NoError(t, err, "JSON serialization should not fail")
	assert.NotEmpty(t, data, "Serialized data should not be empty")

	decoded := &Person{}
	err = DecodeJson(data, decoded)
	assert.NoError(t, err, "JSON deserialization should not fail")

	assert.Equal(t, original.Name, decoded.Name, "Names should match")
	assert.Equal(t, original.Age, decoded.Age, "Ages should match")
	assert.Equal(t, original.Surname, decoded.Surname, "Surnames should match")
}

func BenchmarkEncodeJson1Small(b *testing.B) {
	runEncodeJsonBenchmark(b, 1)
}

func BenchmarkEncodeJson100Small(b *testing.B) {
	runEncodeJsonBenchmark(b, 100)
}

func BenchmarkEncodeJson10000Small(b *testing.B) {
	runEncodeJsonBenchmark(b, 10000)
}

func BenchmarkEncodeJson1000000Small(b *testing.B) {
	runEncodeJsonBenchmark(b, 1000000)
}

func runEncodeJsonBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := generateRandomPersons(count)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EncodeJson(data)
		if err != nil {
			b.Fatalf("Failed to Encode JSON: %v", err)
		}
	}
}

func BenchmarkEncodeJson1Unreal(b *testing.B) {
	runEncodeJsonUnrealBenchmark(b, 1)
}

func BenchmarkEncodeJson10Unreals(b *testing.B) {
	runEncodeJsonUnrealBenchmark(b, 10)
}

func BenchmarkEncodeJson100Unreals(b *testing.B) {
	runEncodeJsonUnrealBenchmark(b, 100)
}

func BenchmarkEncodeJson1000Unreals(b *testing.B) {
	runEncodeJsonUnrealBenchmark(b, 1000)
}

func runEncodeJsonUnrealBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := generateRandomUnreals(count)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EncodeJson(data)
		if err != nil {
			b.Fatalf("Failed to Encode JSON: %v", err)
		}
	}
}
