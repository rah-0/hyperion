package serializer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/hyperion/util/testutil"
)

func TestEncodeDecodeJsonSinglePerson(t *testing.T) {
	original := generateRandomPersons(1)[0]

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

func TestEncodeDecodeJsonTwoPersons(t *testing.T) {
	persons := generateRandomPersons(2)
	original1 := &persons[0]
	original2 := &persons[1]

	data1, err := EncodeJson(original1)
	assert.NoError(t, err, "JSON serialization for first person should not fail")
	assert.NotEmpty(t, data1, "Serialized data for first person should not be empty")

	data2, err := EncodeJson(original2)
	assert.NoError(t, err, "JSON serialization for second person should not fail")
	assert.NotEmpty(t, data2, "Serialized data for second person should not be empty")

	decoded1 := &Person{}
	err = DecodeJson(data1, decoded1)
	assert.NoError(t, err, "JSON deserialization for first person should not fail")

	decoded2 := &Person{}
	err = DecodeJson(data2, decoded2)
	assert.NoError(t, err, "JSON deserialization for second person should not fail")

	assert.Equal(t, original1.Name, decoded1.Name, "Names of first person should match")
	assert.Equal(t, original1.Age, decoded1.Age, "Ages of first person should match")
	assert.Equal(t, original1.Surname, decoded1.Surname, "Surnames of first person should match")

	assert.Equal(t, original2.Name, decoded2.Name, "Names of second person should match")
	assert.Equal(t, original2.Age, decoded2.Age, "Ages of second person should match")
	assert.Equal(t, original2.Surname, decoded2.Surname, "Surnames of second person should match")
}

func BenchmarkEncodeDecodeJson1Small(b *testing.B) {
	runEncodeDecodeJsonBenchmark(b, 1)
}

func BenchmarkEncodeDecodeJson100Small(b *testing.B) {
	runEncodeDecodeJsonBenchmark(b, 100)
}

func BenchmarkEncodeDecodeJson10000Small(b *testing.B) {
	runEncodeDecodeJsonBenchmark(b, 10000)
}

func BenchmarkEncodeDecodeJson1000000Small(b *testing.B) {
	runEncodeDecodeJsonBenchmark(b, 1000000)
}

func runEncodeDecodeJsonBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := generateRandomPersons(count)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encodedData, err := EncodeJson(data)
		if err != nil {
			b.Fatalf("Failed to Encode JSON: %v", err)
		}

		var decoded []Person
		err = DecodeJson(encodedData, &decoded)
		if err != nil {
			b.Fatalf("Failed to Decode JSON: %v", err)
		}
	}
}

func BenchmarkEncodeDecodeJson1Unreal(b *testing.B) {
	runEncodeDecodeJsonUnrealBenchmark(b, 1)
}

func BenchmarkEncodeDecodeJson10Unreals(b *testing.B) {
	runEncodeDecodeJsonUnrealBenchmark(b, 10)
}

func BenchmarkEncodeDecodeJson100Unreals(b *testing.B) {
	runEncodeDecodeJsonUnrealBenchmark(b, 100)
}

func BenchmarkEncodeDecodeJson1000Unreals(b *testing.B) {
	runEncodeDecodeJsonUnrealBenchmark(b, 1000)
}

func runEncodeDecodeJsonUnrealBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := generateRandomUnreals(count)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encodedData, err := EncodeJson(data)
		if err != nil {
			b.Fatalf("Failed to Encode JSON: %v", err)
		}

		var decoded []Unreal
		err = DecodeJson(encodedData, &decoded)
		if err != nil {
			b.Fatalf("Failed to Decode JSON: %v", err)
		}
	}
}
