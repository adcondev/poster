package load_test

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/adcondev/poster/internal/load"
	"golang.org/x/image/bmp"
)

// createSampleImage returns a simple 8x8 RGBA test image.
func createSampleImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	return img
}

func TestSecureFilepath(t *testing.T) {
	tempDir := t.TempDir()
	realTempDir, err := filepath.EvalSymlinks(tempDir)
	if err != nil {
		realTempDir = tempDir
	}

	// Valid relative path
	validRel := "sub/file.png"
	abs, err := load.SecureFilepath(tempDir, validRel)
	if err != nil {
		t.Fatalf("SecureFilepath() unexpected error = %v", err)
	}
	expected := filepath.Join(realTempDir, "sub", "file.png")
	if abs != expected {
		t.Errorf("SecureFilepath() = %q, want %q", abs, expected)
	}

	// Traversal attempt with ../
	_, err = load.SecureFilepath(tempDir, "../escaped.png")
	if err == nil {
		t.Fatal("SecureFilepath() expected error for traversal path, got nil")
	}

	// Deep traversal attempt
	_, err = load.SecureFilepath(tempDir, "a/b/../../../../escaped.png")
	if err == nil {
		t.Fatal("SecureFilepath() expected error for deep traversal path, got nil")
	}
}

func TestImgFromFile_Formats(t *testing.T) {
	tempDir := t.TempDir()
	sample := createSampleImage()

	formats := []struct {
		ext    string
		encode func(f *os.File, img image.Image) error
	}{
		{
			ext: "sample.png",
			encode: func(f *os.File, img image.Image) error {
				return png.Encode(f, img)
			},
		},
		{
			ext: "sample.jpg",
			encode: func(f *os.File, img image.Image) error {
				return jpeg.Encode(f, img, nil)
			},
		},
		{
			ext: "sample.bmp",
			encode: func(f *os.File, img image.Image) error {
				return bmp.Encode(f, img)
			},
		},
	}

	for _, tt := range formats {
		t.Run(tt.ext, func(t *testing.T) {
			path := filepath.Join(tempDir, tt.ext)
			f, err := os.Create(filepath.Clean(path))
			if err != nil {
				t.Fatalf("os.Create() error = %v", err)
			}
			if err := tt.encode(f, sample); err != nil {
				_ = f.Close()
				t.Fatalf("encode error = %v", err)
			}
			_ = f.Close()

			img, err := load.ImgFromFile(tempDir, tt.ext)
			if err != nil {
				t.Fatalf("ImgFromFile(%s) error = %v", tt.ext, err)
			}
			if img.Bounds().Dx() != 8 || img.Bounds().Dy() != 8 {
				t.Errorf("ImgFromFile(%s) dimensions = %dx%d, want 8x8",
					tt.ext, img.Bounds().Dx(), img.Bounds().Dy())
			}
		})
	}
}

func TestImgFromFile_Errors(t *testing.T) {
	tempDir := t.TempDir()

	// Path traversal error
	_, err := load.ImgFromFile(tempDir, "../forbidden.png")
	if err == nil {
		t.Fatal("ImgFromFile() expected error on path traversal, got nil")
	}

	// Non-existent file
	_, err = load.ImgFromFile(tempDir, "non_existent.png")
	if err == nil {
		t.Fatal("ImgFromFile() expected error for non-existent file, got nil")
	}

	// Corrupt / non-image file
	corruptPath := filepath.Join(tempDir, "corrupt.png")
	if err := os.WriteFile(corruptPath, []byte("not an image file"), 0600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	_, err = load.ImgFromFile(tempDir, "corrupt.png")
	if err == nil {
		t.Fatal("ImgFromFile() expected decode error for corrupted image, got nil")
	}
}

func TestImgFromBase64_Success(t *testing.T) {
	sample := createSampleImage()

	var buf bytes.Buffer
	if err := png.Encode(&buf, sample); err != nil {
		t.Fatalf("png.Encode() error = %v", err)
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	img, format, err := load.ImgFromBase64(encoded)
	if err != nil {
		t.Fatalf("ImgFromBase64() error = %v", err)
	}
	if format != "png" {
		t.Errorf("ImgFromBase64() format = %q, want \"png\"", format)
	}
	if img.Bounds().Dx() != 8 || img.Bounds().Dy() != 8 {
		t.Errorf("ImgFromBase64() dimensions = %dx%d, want 8x8", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestImgFromBase64_Errors(t *testing.T) {
	// Invalid base64 encoding
	_, _, err := load.ImgFromBase64("!!!not-valid-base64###")
	if err == nil {
		t.Fatal("ImgFromBase64() expected error for invalid base64, got nil")
	}

	// Data URI prefix is rejected by raw base64 decoder
	_, _, err = load.ImgFromBase64("data:image/png;base64,invalid")
	if err == nil {
		t.Fatal("ImgFromBase64() expected error for data URI prefix with raw decoder, got nil")
	}

	// Valid base64 but corrupt image payload
	corruptBase64 := base64.StdEncoding.EncodeToString([]byte("hello world"))
	_, _, err = load.ImgFromBase64(corruptBase64)
	if err == nil {
		t.Fatal("ImgFromBase64() expected image decode error for invalid payload, got nil")
	}
}
