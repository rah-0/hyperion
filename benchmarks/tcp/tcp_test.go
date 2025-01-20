package main

import (
	"testing"

	"github.com/rah-0/hyperion/utils/testutil"
)

func TestMain(m *testing.M) {
	testutil.TestMainWrapper(testutil.TestConfig{
		M: m,
		LoadResources: func() error {
			if err := tcpStartServer(); err != nil {
				return err
			}

			client = &TCPClient{}
			if err := client.connect(); err != nil {
				return err
			}

			return nil
		},
		UnloadResources: func() error {
			if err := client.Close(); err != nil {
				return err
			}

			if err := server.Close(); err != nil {
				return err
			}

			return nil
		},
	})
}

func BenchmarkTCPClientSendSingleMessage(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)
	msg := []byte("something")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(msg)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage2KB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(1) // 2KB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage4KB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(2) // 4KB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage8KB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(3) // 8KB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage16KB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(4) // 16KB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage32KB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(5) // 32KB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage64KB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(6) // 64KB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage128KB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(7) // 128KB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage256KB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(8) // 256KB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage512KB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(9) // 512KB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage1MB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(10) // 1MB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage10MB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(11) // 10MB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage100MB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(12) // 100MB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkTCPClientSendSingleMessage1GB(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	message, err := testutil.GenerateMessage(13) // 1GB
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := client.send(message)
		if err != nil {
			b.Fatalf("Sending failed on iteration %d: %v", i, err)
		}
	}
}
