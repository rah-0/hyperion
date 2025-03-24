package main

import (
	"testing"
	"time"

	. "github.com/rah-0/hyperion/model"
	. "github.com/rah-0/hyperion/util"
)

func TestNodeDirectConnectionIpAndPort(t *testing.T) {
	c, err := ConnectToNodeWithHostAndPort("127.0.0.1", "5000")
	if err != nil {
		t.Fatal(err)
	}

	msg := Message{
		Type:   MessageTypeTest,
		String: "Test",
	}

	err = c.Send(msg)
	if err != nil {
		t.Fatal(err)
	}
	msg, err = c.Receive()
	if err != nil {
		t.Fatal(err)
	}

	if msg.String != "TestReceived" {
		t.Fatalf("Unexpected message: %s", msg.String)
	}
}

func TestNodeDirectConnection(t *testing.T) {
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		t.Fatal(err)
	}

	msg := Message{
		Type:   MessageTypeTest,
		String: "Test",
	}

	err = c.Send(msg)
	if err != nil {
		t.Fatal(err)
	}
	msg, err = c.Receive()
	if err != nil {
		t.Fatal(err)
	}

	if msg.String != "TestReceived" {
		t.Fatalf("Unexpected message: %s", msg.String)
	}
}

func TestNodesDirectConnection(t *testing.T) {
	var totalExpectedMessages = len(GlobalConfig.Nodes)
	var totalSuccessfulMessages int

	for _, node := range GlobalConfig.Nodes {
		c, err := ConnectToNode(node)
		if err != nil {
			t.Errorf("Failed to connect to node [%s:%d]: %v", node.Host.Name, node.Host.Port, err)
			continue
		}

		msg := Message{
			Type:   MessageTypeTest,
			String: "Test",
		}

		err = c.Send(msg)
		if err != nil {
			t.Errorf("Failed to send message to node [%s:%d]: %v", node.Host.Name, node.Host.Port, err)
			continue
		}

		msg, err = c.Receive()
		if err != nil {
			t.Errorf("Failed to receive message from node [%s:%d]: %v", node.Host.Name, node.Host.Port, err)
			continue
		}

		if msg.String != "TestReceived" {
			t.Errorf("Unexpected response from node [%s:%d]: got %q, want %q",
				node.Host.Name, node.Host.Port, msg.String, "TestReceived")
			continue
		}

		totalSuccessfulMessages++
	}

	if totalSuccessfulMessages != totalExpectedMessages {
		t.Fatalf("Mismatch in total successful messages: got %d, expected %d",
			totalSuccessfulMessages, totalExpectedMessages)
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
	c, err := ConnectToNode(GlobalNode)
	if err != nil {
		b.Fatal(err)
	}

	for sizeLabel, size := range messageSizes {
		b.Run("Size: "+sizeLabel, func(b *testing.B) {
			// Generate a random message of the given size
			testStr, err := GenerateRandomStringMessage(size)
			if err != nil {
				b.Fatalf("Failed to generate message of size %d: %v", size, err)
			}

			msg := Message{
				String: testStr,
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Send message
				if err := c.Send(msg); err != nil {
					b.Fatal(err)
				}

				receivedMsg, err := c.Receive()
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

func BenchmarkConnectionEstablishment(b *testing.B) {
	// Ensure the node is reachable once before benchmarking
	for {
		warmupConn, err := ConnectToNode(GlobalNode)
		if err == nil {
			warmupConn.C.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn, err := ConnectToNode(GlobalNode)
		if err != nil {
			b.Fatalf("Connection failed on iteration %d: %v", i, err)
		}
		conn.C.Close()
	}
}
