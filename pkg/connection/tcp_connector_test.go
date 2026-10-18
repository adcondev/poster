package connection

import (
	"context"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestNewTCPConnector_Success(t *testing.T) {
	var lc net.ListenConfig
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen failed: %v", err)
	}
	defer func() { _ = ln.Close() }()

	serverReceived := make(chan []byte, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer func() { _ = conn.Close() }()
		data, _ := io.ReadAll(conn)
		serverReceived <- data
	}()

	addr := ln.Addr().String()
	connector, err := NewTCPConnector(addr,
		WithDialTimeout(2*time.Second),
		WithWriteTimeout(2*time.Second),
	)
	if err != nil {
		t.Fatalf("NewTCPConnector failed: %v", err)
	}
	if connector.Address() != addr {
		t.Errorf("expected address %s, got %s", addr, connector.Address())
	}

	// Empty write
	n, err := connector.Write([]byte{})
	if err != nil || n != 0 {
		t.Errorf("empty Write() = (%d, %v), want (0, nil)", n, err)
	}

	// Data write
	msg := []byte("Hello Thermal Printer\n")
	n, err = connector.Write(msg)
	if err != nil || n != len(msg) {
		t.Errorf("Write() = (%d, %v), want (%d, nil)", n, err, len(msg))
	}

	// Close
	if err := connector.Close(); err != nil {
		t.Errorf("Close() failed: %v", err)
	}

	// Idempotent second Close
	if err := connector.Close(); err != nil {
		t.Errorf("second Close() failed: %v", err)
	}

	// Write after close
	_, err = connector.Write([]byte("after close"))
	if err == nil {
		t.Error("expected error writing after close, got nil")
	}

	select {
	case received := <-serverReceived:
		if string(received) != string(msg) {
			t.Errorf("server received %q, want %q", string(received), string(msg))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for server to receive data")
	}
}

func TestNewTCPConnector_DefaultPortNormalization(t *testing.T) {
	_, err := NewTCPConnector("192.0.2.1", WithDialTimeout(20*time.Millisecond))
	if err == nil {
		t.Fatal("expected dial error to test address, got nil")
	}
	if !strings.Contains(err.Error(), ":9100") {
		t.Errorf("expected error message to contain ':9100', got %v", err)
	}
}

func TestNewTCPConnector_EmptyAddress(t *testing.T) {
	_, err := NewTCPConnector("   ")
	if err == nil {
		t.Fatal("expected error for empty address, got nil")
	}
}

func TestNewTCPConnector_DialFailure(t *testing.T) {
	_, err := NewTCPConnector("127.0.0.1:1", WithDialTimeout(50*time.Millisecond))
	if err == nil {
		t.Fatal("expected dial error for unreachable port, got nil")
	}
}
