package serializer

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/hyperion/util/testutil"
)

var (
	bytesBuf = &bytes.Buffer{}
	bytesEnc = gob.NewEncoder(bytesBuf)
	bytesDec = gob.NewDecoder(bytesBuf)
)

func TestEncodeDecodeOptBytes(t *testing.T) {
	original := generateRandomPersons(1)[0]

	err := EncodeBytes(bytesEnc, original)
	assert.NoError(t, err, "GOB serialization should not fail")

	decoded := &Person{}
	err = DecodeBytes(bytesDec, decoded)
	assert.NoError(t, err, "GOB deserialization should not fail")

	assert.Equal(t, original.Name, decoded.Name, "Names should match")
	assert.Equal(t, original.Age, decoded.Age, "Ages should match")
	assert.Equal(t, original.Surname, decoded.Surname, "Surnames should match")

	bytesBuf.Reset()
}

func BenchmarkEncodeDecodeBytesOpt1Small(b *testing.B) {
	runEncodeDecodeBytesOptBenchmark(b, 1)
}

func BenchmarkEncodeDecodeBytesOpt100Small(b *testing.B) {
	runEncodeDecodeBytesOptBenchmark(b, 100)
}

func BenchmarkEncodeDecodeBytesOpt10000Small(b *testing.B) {
	runEncodeDecodeBytesOptBenchmark(b, 10000)
}

func BenchmarkEncodeDecodeBytesOpt1000000Small(b *testing.B) {
	runEncodeDecodeBytesOptBenchmark(b, 1000000)
}

func BenchmarkEncodeDecodeBytesOpt1Unreal(b *testing.B) {
	runEncodeDecodeBytesOptUnrealBenchmark(b, 1)
}

func BenchmarkEncodeDecodeBytesOpt10Unreals(b *testing.B) {
	runEncodeDecodeBytesOptUnrealBenchmark(b, 10)
}

func BenchmarkEncodeDecodeBytesOpt100Unreals(b *testing.B) {
	runEncodeDecodeBytesOptUnrealBenchmark(b, 100)
}

func BenchmarkEncodeDecodeBytesOpt1000Unreals(b *testing.B) {
	runEncodeDecodeBytesOptUnrealBenchmark(b, 1000)
}

func runEncodeDecodeBytesOptBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := generateRandomPersons(count)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := EncodeBytes(bytesEnc, data)
		if err != nil {
			b.Fatalf("Failed to Encode Gob: %v", err)
		}
		var decoded []Person
		err = DecodeBytes(bytesDec, &decoded)
		if err != nil {
			b.Fatalf("Failed to Decode Gob: %v", err)
		}
	}

	bytesBuf.Reset()
}

func runEncodeDecodeBytesOptUnrealBenchmark(b *testing.B, count int) {
	defer testutil.RecoverBenchHandler(b)

	data := generateRandomUnreals(count)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := EncodeBytes(bytesEnc, data)
		if err != nil {
			b.Fatalf("Failed to Encode Gob: %v", err)
		}
		var decoded []Unreal
		err = DecodeBytes(bytesDec, &decoded)
		if err != nil {
			b.Fatalf("Failed to Decode Gob: %v", err)
		}
	}

	bytesBuf.Reset()
}
