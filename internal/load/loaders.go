// Package load provides utility functions for image and file loading and decoding.
package load

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"os"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/bmp" // Register BMP decoder
)

// ImgFromFile loads an image from a file path within baseDir.
//
// It validates the path to prevent directory traversal attacks.
func ImgFromFile(baseDir, relPath string) (image.Image, error) {
	// gosec: Path traversal check
	targetAbs, err := SecureFilepath(baseDir, relPath)
	if err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	// Open and decode the image file
	file, err := os.Open(targetAbs) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}

// ImgFromBase64 converts a base64-encoded string to an image.Image.
func ImgFromBase64(data string) (image.Image, string, error) {
	// Decode base64 to byte slice
	imgBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base64 string: %w", err)
	}

	// Convert to image.Image
	img, format, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	return img, format, nil
}

// SecureFilepath constructs an absolute file path from baseDir and relPath,
// ensuring that the resulting path is within baseDir to prevent directory traversal.
func SecureFilepath(baseDir, relPath string) (string, error) {
	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("invalid base directory: %w", err)
	}

	// Resolve symlinks on base directory first to canonicalize path
	if eval, err := filepath.EvalSymlinks(baseAbs); err == nil {
		baseAbs = eval
	}

	target := filepath.Join(baseAbs, relPath)
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return "", fmt.Errorf("invalid target path: %w", err)
	}

	// Resolve symlinks on target to prevent escaping via symlinks
	if eval, err := filepath.EvalSymlinks(targetAbs); err == nil {
		targetAbs = eval
	}

	// Ensure target is inside base directory
	rel, err := filepath.Rel(baseAbs, targetAbs)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("file %s is outside allowed directory", relPath)
	}

	return targetAbs, nil
}
