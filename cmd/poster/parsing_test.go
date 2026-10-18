package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
)

const validDocContent = `{
  "version": "1.0",
  "profile": {
    "model": "POS-80"
  },
  "commands": [
    {
      "type": "text",
      "data": {"text": "data"}
    }
  ]
}`

func TestLoadJSON_SizeLimit(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "poster_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	// Case 1: File within limit
	smallFile := filepath.Join(tempDir, "small.json")
	if err := os.WriteFile(smallFile, []byte(validDocContent), 0600); err != nil {
		t.Fatalf("Failed to write small file: %v", err)
	}

	_, err = loadJSON(smallFile)
	if err != nil {
		t.Errorf("Expected no error for small file, got: %v", err)
	}

	// Case 2: File exceeding limit
	largeFile := filepath.Join(tempDir, "large.json")
	f, err := os.Create(filepath.Clean(largeFile))
	if err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	// Write slightly more than MaxJSONSize (11MB)
	chunk := make([]byte, 1024*1024)
	for i := range chunk {
		chunk[i] = ' '
	}

	for i := 0; i < 11; i++ {
		if _, err := f.Write(chunk); err != nil {
			_ = f.Close()
			t.Fatalf("Failed to write to large file: %v", err)
		}
	}
	if _, err := f.Write([]byte(validDocContent)); err != nil {
		_ = f.Close()
		t.Fatalf("Failed to write end of large file: %v", err)
	}
	_ = f.Close()

	_, err = loadJSON(largeFile)
	if err == nil {
		t.Error("Expected error for large file, got nil")
	} else if !strings.Contains(err.Error(), "file too large") {
		t.Errorf("Expected 'file too large' error, got: %v", err)
	}
}

func TestLoadJSON_ExactLimit(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "poster_test_exact")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	exactFile := filepath.Join(tempDir, "exact.json")
	f, err := os.Create(filepath.Clean(exactFile))
	if err != nil {
		t.Fatalf("Failed to create exact file: %v", err)
	}

	targetSize := constants.MaxJSONSize
	validBytes := []byte(validDocContent)
	if _, err := f.Write(validBytes); err != nil {
		_ = f.Close()
		t.Fatalf("Failed to write valid doc prefix: %v", err)
	}

	chunkSize := 1024 * 1024
	chunk := make([]byte, chunkSize)
	for i := range chunk {
		chunk[i] = ' '
	}

	remaining := targetSize - len(validBytes)
	for remaining > 0 {
		toWrite := remaining
		if toWrite > chunkSize {
			toWrite = chunkSize
		}
		if _, err := f.Write(chunk[:toWrite]); err != nil {
			_ = f.Close()
			t.Fatalf("Failed to write exact file padding: %v", err)
		}
		remaining -= toWrite
	}
	_ = f.Close()

	info, err := os.Stat(exactFile)
	if err != nil {
		t.Fatalf("Failed to stat exact file: %v", err)
	}
	if info.Size() != int64(constants.MaxJSONSize) {
		t.Fatalf("Created file size %d, expected %d", info.Size(), constants.MaxJSONSize)
	}

	data, err := loadJSON(exactFile)
	if err != nil {
		t.Fatalf("Expected no error for exact size file, got: %v", err)
	}

	var doc schema.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("Failed to unmarshal exact size file JSON: %v", err)
	}
	if err := doc.Validate(); err != nil {
		t.Errorf("Document in exact size file failed doc.Validate(): %v", err)
	}
}

func TestValidateJSON(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "poster_validate_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	t.Run("Valid document passes validation", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "valid.json")
		if err := os.WriteFile(filePath, []byte(validDocContent), 0600); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		if err := validateJSON(filePath); err != nil {
			t.Errorf("Expected valid document to pass validateJSON, got: %v", err)
		}
	})

	t.Run("Invalid document schema fails validation", func(t *testing.T) {
		invalidDoc := `{"commands": [{"type": "text", "data": {"text": "hello"}}]}`
		filePath := filepath.Join(tempDir, "invalid_schema.json")
		if err := os.WriteFile(filePath, []byte(invalidDoc), 0600); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		err := validateJSON(filePath)
		if err == nil {
			t.Error("Expected error for invalid document schema, got nil")
		} else if !strings.Contains(err.Error(), "invalid document") {
			t.Errorf("Expected 'invalid document' error, got: %v", err)
		}
	})

	t.Run("Invalid JSON syntax fails validation", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "invalid_syntax.json")
		if err := os.WriteFile(filePath, []byte(`{not valid json}`), 0600); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		err := validateJSON(filePath)
		if err == nil {
			t.Error("Expected error for invalid syntax, got nil")
		} else if !strings.Contains(err.Error(), "invalid JSON") {
			t.Errorf("Expected 'invalid JSON' error, got: %v", err)
		}
	})
}

const validTicketJSON = `{
  "version": "1.0",
  "profile": {
    "model": "POS-80",
    "paper_width": 80,
    "dpi": 203
  },
  "commands": [
    {
      "type": "text",
      "data": {
        "content": {"text": "Receipt Header"}
      }
    },
    {
      "type": "cut",
      "data": {
        "mode": "partial",
        "feed": 0
      }
    }
  ]
}`

func TestExecutePrint_Preview(t *testing.T) {
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "ticket.json")
	if err := os.WriteFile(jsonPath, []byte(validTicketJSON), 0600); err != nil {
		t.Fatalf("Failed to write ticket json: %v", err)
	}

	previewPath := filepath.Join(tempDir, "preview.png")
	cfg := &Config{
		JSONFile:    jsonPath,
		PreviewPath: previewPath,
	}

	if err := executePrint(cfg); err != nil {
		t.Fatalf("executePrint preview failed: %v", err)
	}

	fi, err := os.Stat(previewPath)
	if err != nil {
		t.Fatalf("failed to stat preview png: %v", err)
	}
	if fi.Size() == 0 {
		t.Errorf("expected non-empty preview file, got 0 bytes")
	}

	data, err := os.ReadFile(filepath.Clean(previewPath))
	if err != nil {
		t.Fatalf("failed to read preview png: %v", err)
	}
	if len(data) < 8 || string(data[:8]) != "\x89PNG\r\n\x1a\n" {
		t.Errorf("preview file does not have valid PNG header")
	}
}

func TestExecutePrint_FileOutput(t *testing.T) {
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "ticket.json")
	if err := os.WriteFile(jsonPath, []byte(validTicketJSON), 0600); err != nil {
		t.Fatalf("Failed to write ticket json: %v", err)
	}

	binPath := filepath.Join(tempDir, "receipt.bin")
	cfg := &Config{
		JSONFile:       jsonPath,
		ConnectionType: "file",
		OutputFile:     binPath,
	}

	if err := executePrint(cfg); err != nil {
		t.Fatalf("executePrint file output failed: %v", err)
	}

	fi, err := os.Stat(binPath)
	if err != nil {
		t.Fatalf("failed to stat output bin file: %v", err)
	}
	if fi.Size() == 0 {
		t.Errorf("expected non-empty bin output file, got 0 bytes")
	}

	data, err := os.ReadFile(filepath.Clean(binPath))
	if err != nil {
		t.Fatalf("failed to read bin output: %v", err)
	}
	if !strings.Contains(string(data), "Receipt Header") {
		t.Errorf("expected bin output to contain 'Receipt Header', got %q", string(data))
	}
}

func TestExecutePrint_TCP(t *testing.T) {
	tempDir := t.TempDir()
	jsonPath := filepath.Join(tempDir, "ticket.json")
	if err := os.WriteFile(jsonPath, []byte(validTicketJSON), 0600); err != nil {
		t.Fatalf("Failed to write ticket json: %v", err)
	}

	var lc net.ListenConfig
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on local tcp: %v", err)
	}
	defer func() {
		_ = ln.Close()
	}()

	addr := ln.Addr().String()

	receivedCh := make(chan []byte, 1)
	errCh := make(chan error, 1)

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			errCh <- err
			return
		}
		defer func() {
			_ = conn.Close()
		}()

		var buf bytes.Buffer
		tmp := make([]byte, 1024)
		for {
			_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			n, err := conn.Read(tmp)
			if n > 0 {
				buf.Write(tmp[:n])
			}
			if err != nil {
				break
			}
		}
		receivedCh <- buf.Bytes()
	}()

	cfg := &Config{
		JSONFile:       jsonPath,
		ConnectionType: "tcp",
		PrinterName:    addr,
	}

	if err := executePrint(cfg); err != nil {
		t.Fatalf("executePrint tcp failed: %v", err)
	}

	select {
	case err := <-errCh:
		t.Fatalf("tcp server error: %v", err)
	case received := <-receivedCh:
		if len(received) == 0 {
			t.Errorf("expected non-empty bytes over tcp")
		}
		if !strings.Contains(string(received), "Receipt Header") {
			t.Errorf("expected tcp data to contain 'Receipt Header', got %q", string(received))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for tcp transmission")
	}
}

