package hconn

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/rah-0/hyperion/model"
)

// TestSendAndReceive tests sending and receiving a Message using HConn
func TestSendAndReceive(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := NewHConn(server)
	clientConn := NewHConn(client)

	testMessage := model.Message{
		Type:   model.MessageTypeTest,
		String: "Test",
	}

	go func() {
		if err := clientConn.Send(testMessage); err != nil {
			t.Errorf("Send failed: %v", err)
		}
		client.Close()
	}()

	receivedMessage, err := serverConn.Receive()
	if err != nil {
		t.Fatalf("Receive failed: %v", err)
	}

	if !reflect.DeepEqual(testMessage, receivedMessage) {
		t.Errorf("Expected %+v but got %+v", testMessage, receivedMessage)
	}
}

// TestMultipleMessages ensures multiple sequential messages are handled correctly
func TestMultipleMessages(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := NewHConn(server)
	clientConn := NewHConn(client)

	messages := []model.Message{
		{Type: model.MessageTypeTest, String: "First Message"},
		{Type: model.MessageTypeTest, String: "Second Message"},
		{Type: model.MessageTypeTest, String: "Third Message"},
	}

	go func() {
		for _, msg := range messages {
			if err := clientConn.Send(msg); err != nil {
				t.Errorf("Send failed: %v", err)
			}
		}
		client.Close()
	}()

	for i, expected := range messages {
		received, err := serverConn.Receive()
		if err != nil {
			t.Fatalf("Receive failed on message %d: %v", i, err)
		}
		if !reflect.DeepEqual(received, expected) {
			t.Errorf("Message %d mismatch: expected %+v, got %+v", i, expected, received)
		}
	}
}

// TestLargeMessage ensures large messages are handled correctly
func TestLargeMessage(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := NewHConn(server)
	clientConn := NewHConn(client)

	largeStr := make([]byte, 10*1024) // 10KB message
	for i := range largeStr {
		largeStr[i] = byte(i % 256)
	}
	largeMessage := model.Message{Type: model.MessageTypeTest, String: string(largeStr)}

	go func() {
		if err := clientConn.Send(largeMessage); err != nil {
			t.Errorf("Send failed: %v", err)
		}
		client.Close()
	}()

	received, err := serverConn.Receive()
	if err != nil {
		t.Fatalf("Receive failed: %v", err)
	}

	if received.String != largeMessage.String {
		t.Errorf("Large message mismatch: expected %d bytes, got %d bytes", len(largeMessage.String), len(received.String))
	}
}

// TestEmptyMessage ensures an empty message is handled correctly
func TestEmptyMessage(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := NewHConn(server)
	clientConn := NewHConn(client)

	go func() {
		err := clientConn.Send("")
		if err != nil {
			t.Errorf("Send failed: %v", err)
		}
		client.Close()
	}()

	received, err := serverConn.Receive()
	if err == nil && received.Type != 0 {
		t.Fatalf("Expected error for empty message, but got none")
	}
	if received.String != "" {
		t.Errorf("Expected empty message, but got %+v", received)
	}
}

// TestPartialRead simulates a slow network where data arrives in chunks
func TestPartialRead(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	clientConn := NewHConn(client)

	testMessage := model.Message{Type: model.MessageTypeTest, String: "Hello, Chunked Read!"}

	go func() {
		if err := clientConn.Send(testMessage); err != nil {
			t.Errorf("Send failed: %v", err)
		}
		client.Close()
	}()

	// Simulating slow/chunked read
	lengthBuf := make([]byte, 8)
	_, err := io.ReadFull(server, lengthBuf)
	if err != nil {
		t.Fatalf("Failed to read length prefix: %v", err)
	}

	messageLength := binary.BigEndian.Uint64(lengthBuf)
	receivedBuf := make([]byte, messageLength)

	totalRead := 0
	for totalRead < int(messageLength) {
		chunkSize := 3 // Read in chunks of 3 bytes
		if totalRead+chunkSize > int(messageLength) {
			chunkSize = int(messageLength) - totalRead
		}
		n, err := server.Read(receivedBuf[totalRead : totalRead+chunkSize])
		if err != nil {
			t.Fatalf("Error reading message: %v", err)
		}
		totalRead += n
		time.Sleep(10 * time.Millisecond) // Simulating network delay
	}

	hc := NewHConn(server)
	hc.S.SetData(receivedBuf)
	var decoded model.Message
	err = hc.S.Decode(&decoded)
	if err != nil {
		t.Fatalf("Decoding failed: %v", err)
	}

	if !reflect.DeepEqual(decoded, testMessage) {
		t.Errorf("Partial read mismatch: expected %+v, got %+v", testMessage, decoded)
	}
}

// TestConnectionClose ensures an error is returned when trying to receive on a closed connection
func TestConnectionClose(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()

	serverConn := NewHConn(server)

	client.Close() // Close client before receiving

	_, err := serverConn.Receive()
	if err == nil {
		t.Errorf("Expected error when reading from closed connection, but got none")
	}
}

// TestCorruptedData ensures that incorrect data format (length mismatch) is handled
func TestCorruptedData(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := NewHConn(server)

	go func() {
		// Send incorrect length prefix (indicating more bytes than actually sent)
		lengthPrefix := make([]byte, 8)
		binary.BigEndian.PutUint64(lengthPrefix, 100) // Expecting 100 bytes

		client.Write(lengthPrefix) // Only send the length, no actual message
		client.Close()
	}()

	_, err := serverConn.Receive()
	if err == nil {
		t.Errorf("Expected error due to corrupted data, but got none")
	}
}

func TestReceiveTimeout(t *testing.T) {
	originalTimeout := Timeout
	defer func() { Timeout = originalTimeout }()
	Timeout = 500 * time.Millisecond

	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	hc := NewHConn(server)

	done := make(chan struct{})
	go func() {
		_, err := hc.Receive()
		if err == nil {
			t.Error("Expected timeout error, but got nil")
		} else {
			var netErr net.Error
			if !errors.As(err, &netErr) || !netErr.Timeout() {
				t.Errorf("Expected timeout error, got: %v", err)
			}
		}
		close(done)
	}()

	select {
	case <-done:
		// OK, finished in time
	case <-time.After(2 * Timeout):
		t.Error("Receive did not time out as expected")
	}
}

func BenchmarkSendReceive(b *testing.B) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	sender := NewHConn(client)
	receiver := NewHConn(server)

	msg := model.Message{
		Type:   model.MessageTypeTest,
		String: "BenchmarkSendReceive",
	}

	done := make(chan error, 1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			_, err := receiver.Receive()
			done <- err
		}()

		if err := sender.Send(msg); err != nil {
			b.Fatalf("Send failed: %v", err)
		}

		if err := <-done; err != nil {
			b.Fatalf("Receive failed: %v", err)
		}
	}
}
