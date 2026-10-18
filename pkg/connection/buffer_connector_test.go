package connection

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func TestBufferConnector_Basic(t *testing.T) {
	b := NewBufferConnector()
	if b.Closed() {
		t.Error("expected new buffer connector to be open")
	}
	if b.Len() != 0 {
		t.Errorf("expected len 0, got %d", b.Len())
	}
	if b.Bytes() != nil {
		t.Errorf("expected nil bytes for empty buffer, got %v", b.Bytes())
	}
	if b.String() != "" {
		t.Errorf("expected empty string, got %q", b.String())
	}

	// Write empty slice
	n, err := b.Write([]byte{})
	if err != nil || n != 0 {
		t.Errorf("empty Write() = (%d, %v), want (0, nil)", n, err)
	}

	// Write data
	msg := []byte("hello thermal buffer")
	n, err = b.Write(msg)
	if err != nil || n != len(msg) {
		t.Fatalf("Write() = (%d, %v), want (%d, nil)", n, err, len(msg))
	}

	if b.Len() != len(msg) {
		t.Errorf("expected len %d, got %d", len(msg), b.Len())
	}
	if !bytes.Equal(b.Bytes(), msg) {
		t.Errorf("Bytes() = %v, want %v", b.Bytes(), msg)
	}
	if b.String() != string(msg) {
		t.Errorf("String() = %q, want %q", b.String(), string(msg))
	}

	// Verify Bytes() returns an independent copy
	cp := b.Bytes()
	cp[0] = 'X'
	if b.Bytes()[0] == 'X' {
		t.Error("Bytes() did not return an independent copy")
	}

	// Close
	if err := b.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
	if !b.Closed() {
		t.Error("expected Closed() to be true")
	}

	// Idempotent second Close
	if err := b.Close(); err != nil {
		t.Errorf("second Close() error = %v", err)
	}

	// Write after close
	_, err = b.Write([]byte("more data"))
	if err == nil {
		t.Error("expected error writing to closed buffer connector, got nil")
	}

	// Reset re-opens the connector and clears buffer
	b.Reset()
	if b.Closed() {
		t.Error("expected connector to be open after Reset()")
	}
	if b.Len() != 0 {
		t.Errorf("expected len 0 after Reset(), got %d", b.Len())
	}

	// Can write again after Reset
	_, err = b.Write([]byte("reopened"))
	if err != nil {
		t.Fatalf("Write after Reset failed: %v", err)
	}
	if b.String() != "reopened" {
		t.Errorf("String() = %q, want 'reopened'", b.String())
	}
}

func TestBufferConnector_ConcurrentWrites(t *testing.T) {
	b := NewBufferConnector()
	const numGoroutines = 20
	const writesPerGoroutine = 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < writesPerGoroutine; j++ {
				data := []byte(fmt.Sprintf("[%d:%d]", id, j))
				_, err := b.Write(data)
				if err != nil {
					t.Errorf("concurrent write error: %v", err)
				}
				_ = b.Len()
				_ = b.Bytes()
				_ = b.String()
			}
		}(i)
	}

	wg.Wait()

	if b.Len() == 0 {
		t.Error("expected non-zero length after concurrent writes")
	}
}
