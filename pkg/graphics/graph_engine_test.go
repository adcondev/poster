package graphics_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/adcondev/poster/pkg/graphics"
)

func TestNewPipeline(t *testing.T) {
	// Test with nil options
	p := graphics.NewPipeline(nil)
	if p == nil {
		t.Fatal("NewPipeline(nil) returned nil")
	}

	// Test with valid options
	opts := graphics.DefaultOptions()
	p2 := graphics.NewPipeline(opts)
	if p2 == nil {
		t.Fatal("NewPipeline(opts) returned nil")
	}
}

func TestPipeline_Process(t *testing.T) {
	opts := graphics.DefaultOptions()
	opts.PixelWidth = 10
	opts.Threshold = 128
	p := graphics.NewPipeline(opts)

	// Test nil image
	_, err := p.Process(nil)
	if err == nil {
		t.Error("Process(nil) expected error, got nil")
	}

	// Create a test image
	// 10x10 image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	// Fill with white
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.White)
		}
	}
	// Set a black pixel at 0,0
	img.Set(0, 0, color.Black)

	// Process
	mono, err := p.Process(img)
	if err != nil {
		t.Fatalf("Process(img) error = %v", err)
	}
	if mono == nil {
		t.Fatal("Process(img) returned nil bitmap")
		return
	}

	// Verify dimensions
	if mono.Width != 10 {
		t.Errorf("Process(img) width = %d, want 10", mono.Width)
	}

	// Verify pixels (Threshold mode)
	// Black pixel (0,0) should be true (black)
	if !mono.GetPixel(0, 0) {
		t.Error("Pixel(0,0) should be black (true)")
	}
	// White pixel (1,1) should be false (white)
	if mono.GetPixel(1, 1) {
		t.Error("Pixel(1,1) should be white (false)")
	}
}

func TestPipeline_Process_Resizing(t *testing.T) {
	opts := graphics.DefaultOptions()
	opts.PixelWidth = 5 // Resize down
	p := graphics.NewPipeline(opts)

	img := image.NewRGBA(image.Rect(0, 0, 10, 10)) // 10x10 input
	mono, err := p.Process(img)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if mono.Width != 5 {
		t.Errorf("Output width = %d, want 5", mono.Width)
	}
	// Height should scale proportionally (aspect ratio maintained)
	if mono.Height != 5 {
		t.Errorf("Output height = %d, want 5", mono.Height)
	}
}

func TestPipeline_Process_Atkinson(t *testing.T) {
	opts := graphics.DefaultOptions()
	opts.Dithering = graphics.Atkinson
	opts.Threshold = 128
	p := graphics.NewPipeline(opts)

	// Create a gray image to test dithering pattern
	img := image.NewGray(image.Rect(0, 0, 4, 4))
	for i := 0; i < 16; i++ {
		img.SetGray(i%4, i/4, color.Gray{Y: 128}) // Middle gray
	}

	mono, err := p.Process(img)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Check that we have a mix of black and white pixels
	blackCount := 0
	totalCount := 16
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if mono.GetPixel(x, y) {
				blackCount++
			}
		}
	}

	// With 128 gray and 128 threshold, we expect some dithering.
	// It shouldn't be all black or all white.
	if blackCount == 0 || blackCount == totalCount {
		t.Errorf("Atkinson dithering produced uniform output for gray input. Black count: %d/%d", blackCount, totalCount)
	}
}

func TestPipeline_Resize_Limit(t *testing.T) {
	opts := graphics.DefaultOptions()
	opts.PixelWidth = 1000 // Exceeds 576 limit
	p := graphics.NewPipeline(opts)

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	mono, err := p.Process(img)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Should be capped at 576
	if mono.Width != 576 {
		t.Errorf("Output width = %d, want 576 (capped)", mono.Width)
	}
}

func TestPipeline_Process_Atkinson_DiffusionLogic(t *testing.T) {
	opts := graphics.DefaultOptions()
	opts.Dithering = graphics.Atkinson
	opts.Threshold = 128
	// Use small width to force manual checking and avoid resize
	opts.PixelWidth = 3
	p := graphics.NewPipeline(opts)

	// We test each neighbor independently to avoid cascade interference
	tests := []struct {
		name   string
		source image.Point
		target image.Point
	}{
		{"Right_x+1", image.Point{0, 0}, image.Point{1, 0}},
		{"Right_x+2", image.Point{0, 0}, image.Point{2, 0}},
		{"Below_y+1_x", image.Point{0, 0}, image.Point{0, 1}},
		{"Diag_y+1_x+1", image.Point{0, 0}, image.Point{1, 1}},
		{"Below_y+2_x", image.Point{0, 0}, image.Point{0, 2}},
		{"Diag_y+1_x-1", image.Point{1, 0}, image.Point{0, 1}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			img := image.NewGray(image.Rect(0, 0, 3, 3))
			// Clear
			for i := 0; i < 9; i++ {
				img.SetGray(i%3, i/3, color.Gray{Y: 0})
			}

			// Set source to 100. Error = 100-0 = 100. Diff = 12.
			img.SetGray(tc.source.X, tc.source.Y, color.Gray{Y: 100})

			// Set target to 120. 120+12 = 132 > 128 -> White.
			img.SetGray(tc.target.X, tc.target.Y, color.Gray{Y: 120})

			mono, err := p.Process(img)
			if err != nil {
				t.Fatalf("Process failed: %v", err)
			}

			// Source should be Black
			if !mono.GetPixel(tc.source.X, tc.source.Y) {
				t.Errorf("Pixel(%d,%d) should be Black", tc.source.X, tc.source.Y)
			}

			// Target should be White
			if mono.GetPixel(tc.target.X, tc.target.Y) {
				t.Errorf("Pixel(%d,%d) should be White (received diffusion)", tc.target.X, tc.target.Y)
			}
		})
	}
}

func TestDitherMap_Lookups(t *testing.T) {
	tests := []struct {
		key      string
		expected graphics.DitherMode
	}{
		{"threshold", graphics.Threshold},
		{"atkinson", graphics.Atkinson},
		{"floyd_steinberg", graphics.FloydSteinberg},
		{"ordered", graphics.Ordered},
		{"bayer", graphics.Ordered},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			mode, ok := graphics.DitherMap[tt.key]
			if !ok {
				t.Fatalf("DitherMap missing key %q", tt.key)
			}
			if mode != tt.expected {
				t.Errorf("DitherMap[%q] = %v, want %v", tt.key, mode, tt.expected)
			}
		})
	}
}

func TestPipeline_Process_FloydSteinberg(t *testing.T) {
	opts := graphics.DefaultOptions()
	opts.Dithering = graphics.FloydSteinberg
	opts.Threshold = 128
	p := graphics.NewPipeline(opts)

	// Test with mid-gray image
	img := image.NewGray(image.Rect(0, 0, 8, 8))
	for i := 0; i < 64; i++ {
		img.SetGray(i%8, i/8, color.Gray{Y: 128})
	}

	mono, err := p.Process(img)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	blackCount := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if mono.GetPixel(x, y) {
				blackCount++
			}
		}
	}

	// Mid-gray should yield a mix of black and white pixels
	if blackCount == 0 || blackCount == 64 {
		t.Errorf("Floyd-Steinberg dithering produced uniform output for gray input: %d/64 black", blackCount)
	}

	// Test extreme inputs: all white should produce 0 black pixels
	whiteImg := image.NewGray(image.Rect(0, 0, 4, 4))
	for i := 0; i < 16; i++ {
		whiteImg.SetGray(i%4, i/4, color.Gray{Y: 255})
	}
	whiteMono, err := p.Process(whiteImg)
	if err != nil {
		t.Fatalf("Process whiteImg failed: %v", err)
	}
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if whiteMono.GetPixel(x, y) {
				t.Errorf("Expected white pixel at (%d, %d)", x, y)
			}
		}
	}
}

func TestPipeline_Process_Ordered(t *testing.T) {
	opts := graphics.DefaultOptions()
	opts.Dithering = graphics.Ordered
	opts.Threshold = 128
	p := graphics.NewPipeline(opts)

	// Mid-gray image (128)
	img := image.NewGray(image.Rect(0, 0, 4, 4))
	for i := 0; i < 16; i++ {
		img.SetGray(i%4, i/4, color.Gray{Y: 128})
	}

	mono, err := p.Process(img)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	blackCount := 0
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if mono.GetPixel(x, y) {
				blackCount++
			}
		}
	}

	// In a 4x4 Bayer matrix with base threshold 128, mid-gray produces a balanced pattern
	if blackCount == 0 || blackCount == 16 {
		t.Errorf("Ordered dithering produced uniform output for gray input: %d/16 black", blackCount)
	}

	// Test solid black input
	blackImg := image.NewGray(image.Rect(0, 0, 4, 4))
	for i := 0; i < 16; i++ {
		blackImg.SetGray(i%4, i/4, color.Gray{Y: 0})
	}
	blackMono, err := p.Process(blackImg)
	if err != nil {
		t.Fatalf("Process blackImg failed: %v", err)
	}
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if !blackMono.GetPixel(x, y) {
				t.Errorf("Expected black pixel at (%d, %d)", x, y)
			}
		}
	}
}

func TestMonochromeBitmap_SubBand(t *testing.T) {
	// Create a 16x10 bitmap
	bm := graphics.NewMonochromeBitmap(16, 10)
	// Set specific pixels
	bm.SetPixel(0, 0, true)
	bm.SetPixel(5, 2, true)
	bm.SetPixel(10, 4, true)
	bm.SetPixel(15, 7, true)
	bm.SetPixel(7, 9, true)

	t.Run("first band", func(t *testing.T) {
		band := bm.SubBand(0, 4)
		if band.Width != 16 || band.Height != 4 {
			t.Fatalf("unexpected band dimensions: %dx%d, want 16x4", band.Width, band.Height)
		}
		if !band.GetPixel(0, 0) {
			t.Error("expected pixel at (0, 0) to be black")
		}
		if !band.GetPixel(5, 2) {
			t.Error("expected pixel at (5, 2) to be black")
		}
		if band.GetPixel(10, 3) {
			t.Error("expected pixel at (10, 3) to be white")
		}
	})

	t.Run("middle band", func(t *testing.T) {
		band := bm.SubBand(4, 4)
		if band.Width != 16 || band.Height != 4 {
			t.Fatalf("unexpected band dimensions: %dx%d, want 16x4", band.Width, band.Height)
		}
		// Pixel at (10, 4) in original is at y=0 in band
		if !band.GetPixel(10, 0) {
			t.Error("expected pixel at (10, 0) in band to be black")
		}
		// Pixel at (15, 7) in original is at y=3 in band
		if !band.GetPixel(15, 3) {
			t.Error("expected pixel at (15, 3) in band to be black")
		}
	})

	t.Run("tail band clamped to height", func(t *testing.T) {
		band := bm.SubBand(8, 4) // 8 + 4 = 12 > 10, clamped to 2 rows
		if band.Width != 16 || band.Height != 2 {
			t.Fatalf("unexpected band dimensions: %dx%d, want 16x2", band.Width, band.Height)
		}
		// Pixel at (7, 9) in original is at y=1 in band
		if !band.GetPixel(7, 1) {
			t.Error("expected pixel at (7, 1) in band to be black")
		}
	})

	t.Run("invalid bounds return empty bitmap", func(t *testing.T) {
		invalidCases := []struct {
			name   string
			bm     *graphics.MonochromeBitmap
			startY int
			height int
		}{
			{"negative startY", bm, -1, 4},
			{"startY at height", bm, 10, 2},
			{"startY beyond height", bm, 15, 2},
			{"zero height", bm, 0, 0},
			{"negative height", bm, 0, -2},
			{"nil bitmap", nil, 0, 4},
		}

		for _, ic := range invalidCases {
			t.Run(ic.name, func(t *testing.T) {
				sub := ic.bm.SubBand(ic.startY, ic.height)
				if sub == nil {
					t.Fatal("SubBand returned nil, expected empty bitmap")
				}
				if sub.Width != 0 || sub.Height != 0 {
					t.Errorf("SubBand returned %dx%d, want 0x0", sub.Width, sub.Height)
				}
			})
		}
	})
}
