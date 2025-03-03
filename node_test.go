package main

import (
	"testing"
)

func TestNodesDirectConnection(t *testing.T) {
	t.Skip()

	pathConfig = ""
	err := run()
	if err != nil {
		t.Fatal(err)
	}

	for _, node := range nodes {
		c, err := ConnectToNode(node)
		if err != nil {
			t.Fatal(err)
		}

		msg := Message{
			Type:   MessageTypeTest,
			Mode:   ModeSync,
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
}

func TestNodeToNodeConnection(t *testing.T) {
	t.Skip()

	pathConfig = ""
	err := run()
	if err != nil {
		t.Fatal(err)
	}

	var totalExpectedMessages int
	var totalSuccessfulMessages int

	for _, node := range nodes {
		expectedMessages := len(node.Peers) // Each node should send/receive a message for every peer
		totalExpectedMessages += expectedMessages

		successfulMessages := 0

		for _, peer := range node.Peers {
			msg := Message{
				Type:   MessageTypeTest,
				Mode:   ModeSync,
				String: "Test",
			}

			err = peer.HConn.Send(msg)
			if err != nil {
				t.Errorf("Failed to send message to peer [%s:%d]: %v", peer.Host.Name, peer.Host.Port, err)
				continue
			}

			msg, err = peer.HConn.Receive()
			if err != nil {
				t.Errorf("Failed to receive message from peer [%s:%d]: %v", peer.Host.Name, peer.Host.Port, err)
				continue
			}

			if msg.String != "TestReceived" {
				t.Errorf("Unexpected response from peer [%s:%d]: got %q, want %q", peer.Host.Name, peer.Host.Port, msg.String, "TestReceived")
				continue
			}

			successfulMessages++
		}

		if successfulMessages != expectedMessages {
			t.Errorf("Node [%s:%d] had %d/%d successful message exchanges", node.Host.Name, node.Host.Port, successfulMessages, expectedMessages)
		}

		totalSuccessfulMessages += successfulMessages
	}

	if totalSuccessfulMessages != totalExpectedMessages {
		t.Fatalf("Mismatch in total successful messages: got %d, expected %d", totalSuccessfulMessages, totalExpectedMessages)
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
				c, err := ConnectToNode(node)
				if err != nil {
					b.Fatal(err)
				}

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
