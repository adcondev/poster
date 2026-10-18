package connection_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/adcondev/poster/pkg/connection"
)

func TestFileConnector_WriteAndVerify(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_output.bin")

	conn, err := connection.NewFileConnector(filePath)
	if err != nil {
		t.Fatalf("NewFileConnector failed: %v", err)
	}

	if conn.Path() != filePath {
		t.Errorf("Path() = %q, want %q", conn.Path(), filePath)
	}

	testData := []byte{0x1B, 0x40, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x0A} // ESC @ + "Hello\n"

	// Write empty slice
	n, err := conn.Write([]byte{})
	if err != nil {
		t.Errorf("Write([]byte{}) error: %v", err)
	}
	if n != 0 {
		t.Errorf("Write([]byte{}) = %d, want 0", n)
	}

	// Write test data
	n, err = conn.Write(testData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(testData) {
		t.Fatalf("Write returned %d, want %d", n, len(testData))
	}

	if err := conn.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Verify file content on disk
	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(content) != string(testData) {
		t.Errorf("file content = %v, want %v", content, testData)
	}
}

func TestFileConnector_EmptyPath(t *testing.T) {
	testCases := []string{
		"",
		"   ",
		"\t\n",
	}

	for _, tc := range testCases {
		conn, err := connection.NewFileConnector(tc)
		if err == nil {
			if conn != nil {
				_ = conn.Close()
			}
			t.Errorf("NewFileConnector(%q) expected error, got nil", tc)
		}
		if !errors.Is(err, connection.ErrEmptyPath) {
			t.Errorf("NewFileConnector(%q) error = %v, want %v", tc, err, connection.ErrEmptyPath)
		}
	}
}

func TestFileConnector_IdempotentCloseAndWriteAfterClose(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "idempotent.bin")

	conn, err := connection.NewFileConnector(filePath)
	if err != nil {
		t.Fatalf("NewFileConnector failed: %v", err)
	}

	// First close
	if err := conn.Close(); err != nil {
		t.Fatalf("First Close failed: %v", err)
	}

	// Second close should be idempotent
	if err := conn.Close(); err != nil {
		t.Errorf("Second Close should be idempotent, got: %v", err)
	}

	// Third close
	if err := conn.Close(); err != nil {
		t.Errorf("Third Close should be idempotent, got: %v", err)
	}

	// Write after close should fail
	_, err = conn.Write([]byte{0x01, 0x02})
	if err == nil {
		t.Fatal("Write after Close expected error, got nil")
	}
	if !errors.Is(err, connection.ErrClosedConnector) {
		t.Errorf("Write after Close error = %v, want %v", err, connection.ErrClosedConnector)
	}
}
