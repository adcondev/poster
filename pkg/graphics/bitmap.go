// Package graphics provides utilities for handling monochrome bitmaps for ESC/POS printers
package graphics

import (
	"image"
	"image/color"
)

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

	bytesPerRow := (m.Width + 7) / 8
	byteIndex := y*bytesPerRow + x/8
	// Safe: x is non-negative (checked above) so x%8 is [0,7], result is positive.
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
	// Safe: x is non-negative (checked above) so x%8 is [0,7], result is positive.
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

// ToImage converts the packed monochrome bitmap back to a drawable RGBA image.
// This is used by the Emulator to preview exactly what will be printed.
// The implementation is optimized by processing bytes directly instead of
// calling GetPixel for each pixel.
func (m *MonochromeBitmap) ToImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))

	black := color.RGBA{A: 255}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	bytesPerRow := (m.Width + 7) / 8

	for y := 0; y < m.Height; y++ {
		rowStart := y * bytesPerRow
		for byteIdx := 0; byteIdx < bytesPerRow; byteIdx++ {
			b := m.data[rowStart+byteIdx]
			for bit := 0; bit < 8; bit++ {
				x := byteIdx*8 + bit
				if x >= m.Width {
					break
				}
				if (b & (1 << (7 - bit))) != 0 {
					img.Set(x, y, black)
				} else {
					img.Set(x, y, white)
				}
			}
		}
	}

	return img
}
