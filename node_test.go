package main

import (
	"net"
	"testing"
)

func testConnectNode(t *testing.T, x *Node) (*HConn, error) {
	if t != nil {
		t.Helper()
	}
	conn, err := net.Dial("tcp", x.getListenAddress())
	if err != nil {
		return &HConn{}, err
	}
	return NewHConn(conn), nil
}

func TestListenPortForStatus(t *testing.T) {
	pathConfig = ""
	err := run()
	if err != nil {
		t.Fatal(err)
	}

	for _, node := range nodes {
		conn, err := testConnectNode(t, node)
		if err != nil {
			t.Fatal(err)
		}

		msg := Message{
			Type:   MessageTypeTest,
			Mode:   ModeSync,
			String: "Test",
		}

		err = conn.Send(msg)
		if err != nil {
			t.Fatal(err)
		}
		msg, err = conn.Receive()
		if err != nil {
			t.Fatal(err)
		}

		if msg.String != "TestReceived" {
			t.Fatalf("Unexpected message: %s", msg.String)
		}
	}
}

var messageSizes = map[string]int{
	"2KB":   Size2KB,
	"4KB":   Size4KB,
	"8KB":   Size8KB,
	"16KB":  Size16KB,
	"32KB":  Size32KB,
	"64KB":  Size64KB,
	"128KB": Size128KB,
	"256KB": Size256KB,
	"512KB": Size512KB,
	"1MB":   Size1MB,
	"10MB":  Size10MB,
	"100MB": Size100MB,
}

func BenchmarkListenPortForStatus(b *testing.B) {
	pathConfig = ""
	err := run()
	if err != nil {
		b.Fatal(err)
	}

	for sizeLabel, size := range messageSizes {
		b.Run("Size: "+sizeLabel, func(b *testing.B) {
			// Generate a random message of the given size
			testStr, err := generateRandomStringMessage(size)
			if err != nil {
				b.Fatalf("Failed to generate message of size %d: %v", size, err)
			}

			msg := Message{
				Mode:   ModeSync,
				String: testStr,
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < 1; i++ {
				node := nodes[0]

				conn, err := testConnectNode(nil, node)
				if err != nil {
					b.Fatal(err)
				}

				// Send message
				if err := conn.Send(msg); err != nil {
					b.Fatal(err)
				}

				receivedMsg, err := conn.Receive()
				if err != nil {
					b.Fatal(err)
				}

				// Validate response
				if receivedMsg.String != testStr {
					b.Fatalf("Unexpected response: got %q, want %q", receivedMsg.String, testStr)
				}
			}
		})
	}
}
