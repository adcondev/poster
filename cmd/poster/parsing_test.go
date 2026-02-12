package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

func TestLoadJSON_SizeLimit(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "poster_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Case 1: File within limit
	smallFile := filepath.Join(tempDir, "small.json")
	smallContent := `{"test": "data"}`
	if err := os.WriteFile(smallFile, []byte(smallContent), 0644); err != nil {
		t.Fatalf("Failed to write small file: %v", err)
	}

	_, err = loadJSON(smallFile)
	if err != nil {
		t.Errorf("Expected no error for small file, got: %v", err)
	}

	// Case 2: File exceeding limit
	largeFile := filepath.Join(tempDir, "large.json")
	f, err := os.Create(largeFile)
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
			f.Close()
			t.Fatalf("Failed to write to large file: %v", err)
		}
	}
	if _, err := f.Write([]byte("{}")); err != nil {
		f.Close()
		t.Fatalf("Failed to write end of large file: %v", err)
	}
	f.Close()

	_, err = loadJSON(largeFile)
	if err == nil {
		t.Error("Expected error for large file, got nil")
	} else if !strings.Contains(err.Error(), "file too large") {
		t.Errorf("Expected 'file too large' error, got: %v", err)
	}
}

func TestLoadJSON_ExactLimit(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "poster_test_exact")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Case: File exactly at limit
	exactFile := filepath.Join(tempDir, "exact.json")
	f, err := os.Create(exactFile)
	if err != nil {
		t.Fatalf("Failed to create exact file: %v", err)
	}

	targetSize := constants.MaxJSONSize
	chunkSize := 1024 * 1024
	chunk := make([]byte, chunkSize)
	for i := range chunk {
		chunk[i] = ' '
	}

	// Reserve 2 bytes for "{}"
	remaining := targetSize - 2
	for remaining > 0 {
		toWrite := remaining
		if toWrite > chunkSize {
			toWrite = chunkSize
		}
		if _, err := f.Write(chunk[:toWrite]); err != nil {
			f.Close()
			t.Fatalf("Failed to write exact file: %v", err)
		}
		remaining -= toWrite
	}

	if _, err := f.Write([]byte("{}")); err != nil {
		f.Close()
		t.Fatalf("Failed to write end of exact file: %v", err)
	}
	f.Close()

	info, err := os.Stat(exactFile)
	if err != nil {
		t.Fatalf("Failed to stat exact file: %v", err)
	}
	if info.Size() != int64(constants.MaxJSONSize) {
		t.Fatalf("Created file size %d, expected %d", info.Size(), constants.MaxJSONSize)
	}

	_, err = loadJSON(exactFile)
	if err != nil {
		t.Errorf("Expected no error for exact size file, got: %v", err)
	}
}
