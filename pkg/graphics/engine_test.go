package graphics_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/adcondev/pos-printer/pkg/graphics"
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
		t.Errorf("Process(img) error = %v", err)
	}
	if mono == nil {
		t.Fatal("Process(img) returned nil bitmap")
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
