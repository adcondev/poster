package graphics

import (
	"image"

	"github.com/adcondev/poster/internal/load"
)

// LoadImageFromFile loads an image from a file path within baseDir.
// It validates the path to prevent directory traversal attacks.
func LoadImageFromFile(baseDir, relPath string) (image.Image, string, error) {
	img, err := load.ImgFromFile(baseDir, relPath)
	if err != nil {
		return nil, "", err
	}
	return img, "", nil
}

// LoadImageFromBase64 converts a base64-encoded string to an image.Image.
func LoadImageFromBase64(b64Data string) (image.Image, string, error) {
	return load.ImgFromBase64(b64Data)
}
