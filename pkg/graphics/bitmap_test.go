package graphics_test

import (
	"testing"

	"github.com/adcondev/poster/pkg/graphics"
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

func TestMonochromeBitmap_ToImage_AllWhite(t *testing.T) {
	width, height := 16, 8
	mb := graphics.NewMonochromeBitmap(width, height)
	// All pixels default to white (false)

	img := mb.ToImage()

	if img == nil {
		t.Fatal("ToImage() returned nil")
	}

	bounds := img.Bounds()
	if bounds.Dx() != width || bounds.Dy() != height {
		t.Errorf("ToImage() dimensions = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), width, height)
	}

	// Verify all pixels are white
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r>>8 != 255 || g>>8 != 255 || b>>8 != 255 || a>>8 != 255 {
				t.Errorf("Pixel (%d,%d) = RGBA(%d,%d,%d,%d), want white (255,255,255,255)",
					x, y, r>>8, g>>8, b>>8, a>>8)
			}
		}
	}
}

func TestMonochromeBitmap_ToImage_AllBlack(t *testing.T) {
	width, height := 16, 8
	mb := graphics.NewMonochromeBitmap(width, height)

	// Set all pixels to black
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			mb.SetPixel(x, y, true)
		}
	}

	img := mb.ToImage()

	// Verify all pixels are black
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r>>8 != 0 || g>>8 != 0 || b>>8 != 0 || a>>8 != 255 {
				t.Errorf("Pixel (%d,%d) = RGBA(%d,%d,%d,%d), want black (0,0,0,255)",
					x, y, r>>8, g>>8, b>>8, a>>8)
			}
		}
	}
}

func TestMonochromeBitmap_ToImage_Checkerboard(t *testing.T) {
	width, height := 8, 8
	mb := graphics.NewMonochromeBitmap(width, height)

	// Create checkerboard pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			mb.SetPixel(x, y, (x+y)%2 == 0)
		}
	}

	img := mb.ToImage()

	// Verify checkerboard pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, _, _, _ := img.At(x, y).RGBA()
			isBlack := r>>8 == 0
			expectBlack := (x+y)%2 == 0

			if isBlack != expectBlack {
				t.Errorf("Pixel (%d,%d) isBlack=%v, want %v", x, y, isBlack, expectBlack)
			}
		}
	}
}

func TestMonochromeBitmap_ToImage_NonByteAlignedWidth(t *testing.T) {
	// Test with width that's not a multiple of 8
	tests := []struct {
		width  int
		height int
	}{
		{1, 1},
		{7, 3},
		{9, 5},
		{13, 7},
		{100, 10},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			mb := graphics.NewMonochromeBitmap(tt.width, tt.height)

			// Set some pixels
			for y := 0; y < tt.height; y++ {
				for x := 0; x < tt.width; x++ {
					if x%2 == 0 {
						mb.SetPixel(x, y, true)
					}
				}
			}

			img := mb.ToImage()

			// Verify dimensions
			bounds := img.Bounds()
			if bounds.Dx() != tt.width || bounds.Dy() != tt.height {
				t.Errorf("Dimensions = %dx%d, want %dx%d",
					bounds.Dx(), bounds.Dy(), tt.width, tt.height)
			}

			// Verify pattern
			for y := 0; y < tt.height; y++ {
				for x := 0; x < tt.width; x++ {
					r, _, _, _ := img.At(x, y).RGBA()
					isBlack := r>>8 == 0
					expectBlack := x%2 == 0

					if isBlack != expectBlack {
						t.Errorf("Pixel (%d,%d) isBlack=%v, want %v", x, y, isBlack, expectBlack)
					}
				}
			}
		})
	}
}

func TestMonochromeBitmap_ToImage_SinglePixel(t *testing.T) {
	// Edge case: single pixel bitmap
	mb := graphics.NewMonochromeBitmap(1, 1)

	// Test white
	img := mb.ToImage()
	r, g, b, a := img.At(0, 0).RGBA()
	if r>>8 != 255 || g>>8 != 255 || b>>8 != 255 {
		t.Errorf("Single white pixel = RGBA(%d,%d,%d,%d), want white",
			r>>8, g>>8, b>>8, a>>8)
	}

	// Test black
	mb.SetPixel(0, 0, true)
	img = mb.ToImage()
	r, g, b, a = img.At(0, 0).RGBA()
	if r>>8 != 0 || g>>8 != 0 || b>>8 != 0 {
		t.Errorf("Single black pixel = RGBA(%d,%d,%d,%d), want black",
			r>>8, g>>8, b>>8, a>>8)
	}
}

func TestMonochromeBitmap_ToImage_SpecificBits(t *testing.T) {
	// Test specific bit positions to verify bit ordering
	mb := graphics.NewMonochromeBitmap(8, 1)

	// Set only the first pixel (MSB of first byte)
	mb.SetPixel(0, 0, true)

	img := mb.ToImage()

	// Only pixel 0 should be black
	for x := 0; x < 8; x++ {
		r, _, _, _ := img.At(x, 0).RGBA()
		isBlack := r>>8 == 0
		expectBlack := x == 0

		if isBlack != expectBlack {
			t.Errorf("Pixel (%d,0) isBlack=%v, want %v", x, isBlack, expectBlack)
		}
	}

	// Now test last pixel (LSB of first byte)
	mb2 := graphics.NewMonochromeBitmap(8, 1)
	mb2.SetPixel(7, 0, true)

	img2 := mb2.ToImage()

	// Only pixel 7 should be black
	for x := 0; x < 8; x++ {
		r, _, _, _ := img2.At(x, 0).RGBA()
		isBlack := r>>8 == 0
		expectBlack := x == 7

		if isBlack != expectBlack {
			t.Errorf("Pixel (%d,0) isBlack=%v, want %v", x, isBlack, expectBlack)
		}
	}
}

func TestMonochromeBitmap_ToImage_RoundTrip(t *testing.T) {
	// Verify that SetPixel -> ToImage produces correct results
	width, height := 24, 16
	mb := graphics.NewMonochromeBitmap(width, height)

	// Create a specific pattern
	pattern := make([][]bool, height)
	for y := 0; y < height; y++ {
		pattern[y] = make([]bool, width)
		for x := 0; x < width; x++ {
			// Diagonal stripes
			isBlack := (x+y)%3 == 0
			pattern[y][x] = isBlack
			mb.SetPixel(x, y, isBlack)
		}
	}

	img := mb.ToImage()

	// Verify pattern matches
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, _, _, _ := img.At(x, y).RGBA()
			isBlack := r>>8 == 0

			if isBlack != pattern[y][x] {
				t.Errorf("Pixel (%d,%d) isBlack=%v, want %v", x, y, isBlack, pattern[y][x])
			}
		}
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkMonochromeBitmap_ToImage_Small(b *testing.B) {
	mb := graphics.NewMonochromeBitmap(100, 100)
	// Set some pixels
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if (x+y)%2 == 0 {
				mb.SetPixel(x, y, true)
			}
		}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = mb.ToImage()
	}
}

func BenchmarkMonochromeBitmap_ToImage_Medium(b *testing.B) {
	mb := graphics.NewMonochromeBitmap(384, 500)
	// Set some pixels
	for y := 0; y < 500; y++ {
		for x := 0; x < 384; x++ {
			if (x+y)%2 == 0 {
				mb.SetPixel(x, y, true)
			}
		}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = mb.ToImage()
	}
}

func BenchmarkMonochromeBitmap_ToImage_Large(b *testing.B) {
	mb := graphics.NewMonochromeBitmap(576, 1000)
	// Set some pixels
	for y := 0; y < 1000; y++ {
		for x := 0; x < 576; x++ {
			if (x+y)%2 == 0 {
				mb.SetPixel(x, y, true)
			}
		}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = mb.ToImage()
	}
}
