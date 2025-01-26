package compression

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rah-0/hyperion/util/testutil"
)

func TestGzip_2KB(t *testing.T) {
	runGzipTest(t, 1)
}

func TestGzip_4KB(t *testing.T) {
	runGzipTest(t, 2)
}

func TestGzip_8KB(t *testing.T) {
	runGzipTest(t, 3)
}

func TestGzip_16KB(t *testing.T) {
	runGzipTest(t, 4)
}

func TestGzip_32KB(t *testing.T) {
	runGzipTest(t, 5)
}

func TestGzip_64KB(t *testing.T) {
	runGzipTest(t, 6)
}

func TestGzip_128KB(t *testing.T) {
	runGzipTest(t, 7)
}

func TestGzip_256KB(t *testing.T) {
	runGzipTest(t, 8)
}

func TestGzip_512KB(t *testing.T) {
	runGzipTest(t, 9)
}

func TestGzip_1MB(t *testing.T) {
	runGzipTest(t, 10)
}

func TestGzip_10MB(t *testing.T) {
	runGzipTest(t, 11)
}

func runGzipTest(t *testing.T, sizeType int) {
	originalData, err := testutil.GenerateMessage(sizeType)
	if err != nil {
		t.Fatalf("Failed to generate message of size type %d: %v", sizeType, err)
	}

	compressedData, err := CompressGzip(originalData)
	if err != nil {
		t.Fatalf("Failed to compress data for size type %d: %v", sizeType, err)
	}

	decompressedData, err := DecompressGzip(compressedData)
	if err != nil {
		t.Fatalf("Failed to decompress data for size type %d: %v", sizeType, err)
	}

	assert.Equal(t, originalData, decompressedData, "Decompressed data should match the original for size type %d", sizeType)

	dataSize := len(originalData)
	dataSizeCompressed := len(compressedData)
	fmt.Printf("Encoded size: %d bytes, Compressed size: %d bytes, Gain %.2f%%\n",
		dataSize, dataSizeCompressed, testutil.PercentDifference(dataSize, dataSizeCompressed))
}

func BenchmarkCompressDecompressGzip(b *testing.B) {
	data, err := testutil.GenerateMessage(10) // 10 corresponds to 1MB
	if err != nil {
		b.Fatalf("Failed to generate message of size type %d: %v", 10, err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		compressedData, err := CompressGzip(data)
		if err != nil {
			b.Fatalf("Failed to Compress with Gzip: %v", err)
		}

		_, err = DecompressGzip(compressedData)
		if err != nil {
			b.Fatalf("Failed to Decompress with Gzip: %v", err)
		}
	}
}
