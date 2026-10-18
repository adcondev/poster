package connection

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// ErrClosedConnector is returned when attempting an operation on a closed connector.
	ErrClosedConnector = errors.New("connector is closed")

	// ErrEmptyPath is returned when the provided file path is empty.
	ErrEmptyPath = errors.New("file path cannot be empty")
)

var _ Connector = (*FileConnector)(nil)

// FileConnector writes raw ESC/POS output to a file or device node.
type FileConnector struct {
	file   *os.File
	path   string
	closed bool
	mu     sync.Mutex
}

// NewFileConnector creates a new FileConnector at the specified path.
func NewFileConnector(path string) (*FileConnector, error) {
	if strings.TrimSpace(path) == "" {
		return nil, ErrEmptyPath
	}

	cleanPath := filepath.Clean(path)
	file, err := os.OpenFile(cleanPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", cleanPath, err)
	}

	return &FileConnector{
		file: file,
		path: cleanPath,
	}, nil
}

// Write writes data bytes to the file connector.
func (c *FileConnector) Write(data []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return 0, ErrClosedConnector
	}
	if len(data) == 0 {
		return 0, nil
	}

	n, err := c.file.Write(data)
	if err != nil {
		return n, fmt.Errorf("failed to write to file: %w", err)
	}
	return n, nil
}

// Close closes the file connector idempotently.
func (c *FileConnector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}
	c.closed = true

	if err := c.file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

// Path returns the cleaned file path of the connector.
func (c *FileConnector) Path() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.path
}
