package connection

import (
	"bytes"
	"errors"
	"sync"
)

// BufferConnector is a thread-safe in-memory Connector for testing and previewing.
type BufferConnector struct {
	buf    bytes.Buffer
	mu     sync.RWMutex
	closed bool
}

// NewBufferConnector creates a new in-memory buffer connector.
func NewBufferConnector() *BufferConnector {
	return &BufferConnector{}
}

// Write appends data to the internal buffer. Returns an error if the connector is closed.
func (b *BufferConnector) Write(data []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return 0, errors.New("buffer connector is closed")
	}

	if len(data) == 0 {
		return 0, nil
	}

	return b.buf.Write(data)
}

// Close marks the buffer connector as closed. It is idempotent.
func (b *BufferConnector) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.closed = true
	return nil
}

// Bytes returns a safe copy of the buffered bytes.
func (b *BufferConnector) Bytes() []byte {
	b.mu.RLock()
	defer b.mu.RUnlock()

	raw := b.buf.Bytes()
	if len(raw) == 0 {
		return nil
	}
	cp := make([]byte, len(raw))
	copy(cp, raw)
	return cp
}

// String returns the contents of the buffer as a string.
func (b *BufferConnector) String() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.buf.String()
}

// Len returns the number of bytes written to the buffer.
func (b *BufferConnector) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.buf.Len()
}

// Reset clears the buffer and re-opens the connector if it was closed.
func (b *BufferConnector) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buf.Reset()
	b.closed = false
}

// Closed returns true if the buffer connector has been closed.
func (b *BufferConnector) Closed() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.closed
}

// Ensure BufferConnector implements Connector.
var _ Connector = (*BufferConnector)(nil)
