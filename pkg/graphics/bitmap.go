// Package graphics provides utilities for handling monochrome bitmaps for ESC/POS printers
package graphics

// MonochromeBitmap represents a black and white bitmap optimized for ESC/POS
type MonochromeBitmap struct {
	Width  int
	Height int
	data   []byte // Packed bit data
}

// NewMonochromeBitmap creates a new monochrome bitmap
func NewMonochromeBitmap(width, height int) *MonochromeBitmap {
	bytesPerRow := (width + 7) / 8
	return &MonochromeBitmap{
		Width:  width,
		Height: height,
		data:   make([]byte, bytesPerRow*height),
	}
}

// SetPixel sets a pixel to black (true) or white (false)
func (m *MonochromeBitmap) SetPixel(x, y int, black bool) {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return
	}

	// TODO: Make sure linter suppression is safe here

	bytesPerRow := (m.Width + 7) / 8
	byteIndex := y*bytesPerRow + x/8
	bitIndex := uint(7 - (x % 8)) //nolint:gosec

	if black {
		m.data[byteIndex] |= 1 << bitIndex
	} else {
		m.data[byteIndex] &^= 1 << bitIndex
	}
}

// GetPixel returns whether a pixel is black (true) or white (false)
func (m *MonochromeBitmap) GetPixel(x, y int) bool {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return false
	}

	bytesPerRow := (m.Width + 7) / 8
	byteIndex := y*bytesPerRow + x/8
	bitIndex := uint(7 - (x % 8)) //nolint:gosec

	return (m.data[byteIndex] & (1 << bitIndex)) != 0
}

// GetRasterData returns the raw raster data for ESC/POS printing
func (m *MonochromeBitmap) GetRasterData() []byte {
	return m.data
}

// GetWidthBytes returns the width in bytes (for ESC/POS commands)
func (m *MonochromeBitmap) GetWidthBytes() int {
	return (m.Width + 7) / 8
}
