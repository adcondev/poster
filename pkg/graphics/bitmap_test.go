package graphics_test

import (
	"testing"

	"github.com/adcondev/pos-printer/pkg/graphics"
)

func TestMonochromeBitmap_New(t *testing.T) {
	width, height := 16, 10
	mb := graphics.NewMonochromeBitmap(width, height)

	if mb.Width != width {
		t.Errorf("NewMonochromeBitmap(%d, %d).Width = %d, want %d", width, height, mb.Width, width)
	}
	if mb.Height != height {
		t.Errorf("NewMonochromeBitmap(%d, %d).Height = %d, want %d", width, height, mb.Height, height)
	}

	expectedBytes := ((width + 7) / 8) * height
	if len(mb.GetRasterData()) != expectedBytes {
		t.Errorf("NewMonochromeBitmap(%d, %d) data length = %d, want %d", width, height, len(mb.GetRasterData()), expectedBytes)
	}
}

func TestMonochromeBitmap_SetPixel_GetPixel(t *testing.T) {
	width, height := 8, 8
	mb := graphics.NewMonochromeBitmap(width, height)

	// Test valid pixels
	tests := []struct {
		x, y  int
		black bool
	}{
		{0, 0, true},
		{7, 7, true},
		{4, 4, true},
		{0, 7, true},
		{7, 0, true},
		{2, 2, false}, // Should remain false (default) or be set to false
	}

	for _, tt := range tests {
		mb.SetPixel(tt.x, tt.y, tt.black)
		got := mb.GetPixel(tt.x, tt.y)
		if got != tt.black {
			t.Errorf("GetPixel(%d, %d) = %v, want %v", tt.x, tt.y, got, tt.black)
		}
	}

	// Test out of bounds
	oobTests := []struct {
		x, y int
	}{
		{-1, 0},
		{0, -1},
		{width, 0},
		{0, height},
	}

	for _, tt := range oobTests {
		// SetPixel should not panic
		mb.SetPixel(tt.x, tt.y, true)
		// GetPixel should return false
		if mb.GetPixel(tt.x, tt.y) {
			t.Errorf("GetPixel(%d, %d) = true, want false (out of bounds)", tt.x, tt.y)
		}
	}
}

func TestMonochromeBitmap_GetWidthBytes(t *testing.T) {
	tests := []struct {
		width int
		want  int
	}{
		{8, 1},
		{1, 1},
		{9, 2},
		{16, 2},
		{0, 0},
	}

	for _, tt := range tests {
		mb := graphics.NewMonochromeBitmap(tt.width, 1)
		if got := mb.GetWidthBytes(); got != tt.want {
			t.Errorf("GetWidthBytes() for width %d = %d, want %d", tt.width, got, tt.want)
		}
	}
}

func TestMonochromeBitmap_ManualBitCheck(t *testing.T) {
	// Verify specific bit manipulation
	mb := graphics.NewMonochromeBitmap(8, 1)
	mb.SetPixel(0, 0, true) // Should set MSB (bit 7) of first byte: 10000000 -> 0x80
	mb.SetPixel(7, 0, true) // Should set LSB (bit 0) of first byte: 00000001 -> combined 0x81

	data := mb.GetRasterData()
	if data[0] != 0x81 {
		t.Errorf("GetRasterData() = 0x%X, want 0x81", data[0])
	}
}
