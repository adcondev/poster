package graphics_test

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/adcondev/poster/pkg/graphics"
)

func createSamplePNGBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode sample png: %v", err)
	}
	return buf.Bytes()
}

func TestLoadImageFromFile(t *testing.T) {
	tempDir := t.TempDir()
	pngData := createSamplePNGBytes(t)
	fileName := "sample.png"
	fullPath := filepath.Join(tempDir, fileName)

	if err := os.WriteFile(fullPath, pngData, 0600); err != nil {
		t.Fatalf("failed to write sample png: %v", err)
	}

	t.Run("valid file", func(t *testing.T) {
		img, _, err := graphics.LoadImageFromFile(tempDir, fileName)
		if err != nil {
			t.Fatalf("LoadImageFromFile failed: %v", err)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
		if img.Bounds().Dx() != 10 || img.Bounds().Dy() != 10 {
			t.Errorf("image dimensions = %dx%d, want 10x10", img.Bounds().Dx(), img.Bounds().Dy())
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, _, err := graphics.LoadImageFromFile(tempDir, "missing.png")
		if err == nil {
			t.Error("expected error for missing file, got nil")
		}
	})

	t.Run("directory traversal attempt", func(t *testing.T) {
		_, _, err := graphics.LoadImageFromFile(tempDir, "../../../etc/passwd")
		if err == nil {
			t.Error("expected error for path traversal, got nil")
		}
	})
}

func TestLoadImageFromBase64(t *testing.T) {
	pngData := createSamplePNGBytes(t)
	rawB64 := base64.StdEncoding.EncodeToString(pngData)
	dataURIB64 := "data:image/png;base64," + rawB64

	t.Run("raw base64 string", func(t *testing.T) {
		img, format, err := graphics.LoadImageFromBase64(rawB64)
		if err != nil {
			t.Fatalf("LoadImageFromBase64 failed: %v", err)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
		if format != "png" {
			t.Errorf("format = %q, want png", format)
		}
	})

	t.Run("data URI base64 string rejected", func(t *testing.T) {
		_, _, err := graphics.LoadImageFromBase64(dataURIB64)
		if err == nil {
			t.Error("expected error for data URI prefix in LoadImageFromBase64, got nil")
		}
	})

	t.Run("empty string error", func(t *testing.T) {
		_, _, err := graphics.LoadImageFromBase64("")
		if err == nil {
			t.Error("expected error for empty base64 string, got nil")
		}
	})

	t.Run("invalid base64 content error", func(t *testing.T) {
		_, _, err := graphics.LoadImageFromBase64("this-is-not-valid-base64!@#$")
		if err == nil {
			t.Error("expected error for corrupt base64 string, got nil")
		}
	})

	t.Run("corrupt image payload error", func(t *testing.T) {
		corrupt := base64.StdEncoding.EncodeToString([]byte("not an image"))
		_, _, err := graphics.LoadImageFromBase64(corrupt)
		if err == nil {
			t.Error("expected error for corrupt image payload, got nil")
		}
	})
}
