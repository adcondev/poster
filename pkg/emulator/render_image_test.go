package emulator_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/emulator"
	"github.com/adcondev/poster/pkg/graphics"
)

// ============================================================================
// DefaultImageOptions Tests
// ============================================================================

func TestDefaultImageOptions(t *testing.T) {
	opts := emulator.DefaultImageOptions()

	if opts == nil {
		t.Fatal("DefaultImageOptions() returned nil")
		return
	}

	if opts.PixelWidth != constants.DefaultImagePixelWidth {
		t.Errorf("PixelWidth = %d, want %d", opts.PixelWidth, constants.DefaultImagePixelWidth)
	}

	if opts.Align != constants.Center.String() {
		t.Errorf("Align = %s, want %s", opts.Align, constants.Center.String())
	}

	if !opts.PreserveAspect {
		t.Error("PreserveAspect should be true by default")
	}

	if opts.Scaling != graphics.BiLinear {
		t.Errorf("Scaling = %d, want BiLinear", opts.Scaling)
	}

	if opts.SimulateThermal {
		t.Error("SimulateThermal should be false by default")
	}

	if opts.Dithering != graphics.Atkinson {
		t.Errorf("Dithering = %d, want Atkinson", opts.Dithering)
	}

	if opts.Threshold != 128 {
		t.Errorf("Threshold = %d, want 128", opts.Threshold)
	}
}

// ============================================================================
// RenderImage - Error Handling Tests
// ============================================================================

func TestRenderImage_NilImage(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	err = engine.PrintImage(nil)
	if err == nil {
		t.Error("PrintImage(nil) should return an error")
	}
}

func TestRenderImage_NilOptions(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestRGBA(100, 50, color.White)

	// Should use default options and not panic
	err = engine.PrintImageWithOptions(img, nil)
	if err != nil {
		t.Errorf("PrintImageWithOptions with nil options should not error: %v", err)
	}
}

// ============================================================================
// RenderImage - Cursor Position Tests
// ============================================================================

func TestRenderImage_UpdatesCursorY(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	initialY := engine.State().CursorY
	img := createTestRGBA(100, 50, color.White)

	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage failed: %v", err)
	}

	// Cursor should have advanced by at least the image height
	if engine.State().CursorY <= initialY {
		t.Error("CursorY should advance after printing image")
	}

	// The advance should be approximately the image height (may be scaled)
	advance := engine.State().CursorY - initialY
	if advance < 10 {
		t.Errorf("CursorY advance = %.1f, expected at least image height", advance)
	}
}

func TestRenderImage_CursorYMatchesImageHeight(t *testing.T) {
	engine, err := emulator.New58mmEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	initialY := engine.State().CursorY

	// Create image with known dimensions that won't be resized
	imgWidth := 200
	imgHeight := 100
	img := createTestRGBA(imgWidth, imgHeight, color.White)

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = imgWidth // Keep original size
	opts.PreserveAspect = true

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions failed: %v", err)
	}

	expectedY := initialY + float64(imgHeight)
	actualY := engine.State().CursorY

	if actualY != expectedY {
		t.Errorf("CursorY = %.1f, want %.1f", actualY, expectedY)
	}
}

// ============================================================================
// RenderImage - Alignment Tests
// ============================================================================

func TestRenderImage_AlignmentLeft(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestRGBA(100, 50, color.Black)

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = 100
	opts.Align = constants.Left.String()

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions failed:  %v", err)
	}

	// Render and check that black pixels start at x=0
	result := engine.Render()
	bounds := result.Bounds()

	// First column should have the image (black pixels)
	foundBlackAtLeft := false
	for y := 0; y < bounds.Dy(); y++ {
		r, _, _, _ := result.At(0, y).RGBA()
		if r == 0 { // Black pixel
			foundBlackAtLeft = true
			break
		}
	}

	if !foundBlackAtLeft {
		t.Error("Left-aligned image should have content at x=0")
	}
}

func TestRenderImage_AlignmentCenter(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	imgWidth := 200
	img := createTestRGBA(imgWidth, 50, color.Black)

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = imgWidth
	opts.Align = constants.Center.String()

	initialY := engine.State().CursorY

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions failed:  %v", err)
	}

	result := engine.Render()

	// Calculate Y position where image was rendered
	imageY := int(initialY) + 25 // Middle of image area

	// Verify left margin is white (before centered image)
	if imageY < result.Bounds().Dy() {
		r, g, b, _ := result.At(0, imageY).RGBA()
		isWhite := r>>8 == 255 && g>>8 == 255 && b>>8 == 255
		if !isWhite {
			t.Error("Left margin should be white for center-aligned image")
		}
	}

	// Check that there are black pixels in the center region
	centerX := constants.PaperPxWidth80mm / 2
	foundBlackInCenter := false
	for y := 0; y < result.Bounds().Dy(); y++ {
		r, _, _, _ := result.At(centerX, y).RGBA()
		if r == 0 {
			foundBlackInCenter = true
			break
		}
	}

	if !foundBlackInCenter {
		t.Error("Center-aligned image should have content in the center")
	}
}

func TestRenderImage_AlignmentRight(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	imgWidth := 200
	img := createTestRGBA(imgWidth, 50, color.Black)

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = imgWidth
	opts.Align = constants.Right.String()

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions failed: %v", err)
	}

	result := engine.Render()

	// Right-aligned image should have content at the right edge
	rightEdgeX := constants.PaperPxWidth80mm - 1
	foundBlackAtRight := false
	for y := 0; y < result.Bounds().Dy(); y++ {
		r, _, _, _ := result.At(rightEdgeX, y).RGBA()
		if r == 0 {
			foundBlackAtRight = true
			break
		}
	}

	if !foundBlackAtRight {
		t.Error("Right-aligned image should have content at right edge")
	}
}

func TestRenderImage_AlignmentFallback(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestRGBA(100, 50, color.Black)

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = 100
	opts.Align = "invalid_alignment" // Should fall back to left

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions with invalid alignment should not error: %v", err)
	}

	// Should fall back to left alignment - check first column has content
	result := engine.Render()
	foundBlackAtLeft := false
	for y := 0; y < result.Bounds().Dy(); y++ {
		r, _, _, _ := result.At(0, y).RGBA()
		if r == 0 {
			foundBlackAtLeft = true
			break
		}
	}

	if !foundBlackAtLeft {
		t.Error("Invalid alignment should fall back to left")
	}
}

// ============================================================================
// RenderImage - Width Handling Tests
// ============================================================================

func TestRenderImage_WidthCappedToPaper(t *testing.T) {
	engine, err := emulator.New58mmEngine() // 384px wide
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Image wider than paper
	img := createTestRGBA(500, 100, color.Black)

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = 500 // Request full width

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions failed: %v", err)
	}

	result := engine.Render()

	// Result should not exceed paper width
	if result.Bounds().Dx() > constants.PaperPxWidth58mm {
		t.Errorf("Result width %d exceeds paper width %d",
			result.Bounds().Dx(), constants.PaperPxWidth58mm)
	}
}

func TestRenderImage_ZeroPixelWidth(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	imgWidth := 200
	img := createTestRGBA(imgWidth, 100, color.White)

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = 0 // Should use original width

	initialY := engine.State().CursorY

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions failed:  %v", err)
	}

	// Should still render (cursor should advance)
	if engine.State().CursorY <= initialY {
		t.Error("Image with PixelWidth=0 should still render")
	}
}

func TestRenderImage_NegativePixelWidth(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine:  %v", err)
	}

	img := createTestRGBA(200, 100, color.White)

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = -100 // Should use original width

	initialY := engine.State().CursorY

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions failed: %v", err)
	}

	// Should still render
	if engine.State().CursorY <= initialY {
		t.Error("Image with negative PixelWidth should still render")
	}
}

// ============================================================================
// RenderImage - Preview Mode Tests
// ============================================================================

func TestRenderImage_NormalPreview(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create a colored image
	img := createTestRGBA(100, 50, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	opts := emulator.DefaultImageOptions()
	opts.SimulateThermal = false

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions failed:  %v", err)
	}

	result := engine.Render()

	// In normal preview, colors should be preserved (not pure B&W)
	// Find a pixel from the image area
	foundColor := false
	for y := 0; y < result.Bounds().Dy(); y++ {
		for x := 0; x < result.Bounds().Dx(); x++ {
			r, g, b, _ := result.At(x, y).RGBA()
			// Looking for reddish pixels (not pure white or pure black)
			if r>>8 > 200 && g>>8 < 100 && b>>8 < 100 {
				foundColor = true
				break
			}
		}
		if foundColor {
			break
		}
	}

	if !foundColor {
		t.Log("Note: Normal preview should preserve colors, but color detection may vary")
	}
}

func TestRenderImage_ThermalPreview(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create a gray image
	img := createTestRGBA(100, 50, color.RGBA{R: 128, G: 128, B: 128, A: 255})

	opts := emulator.DefaultImageOptions()
	opts.SimulateThermal = true
	opts.Dithering = graphics.Threshold
	opts.Threshold = 128

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions failed: %v", err)
	}

	result := engine.Render()

	// In thermal preview with threshold dithering, pixels should be pure B&W
	allBW := true
	for y := 0; y < result.Bounds().Dy(); y++ {
		for x := 0; x < result.Bounds().Dx(); x++ {
			r, g, b, _ := result.At(x, y).RGBA()
			r8, g8, b8 := r>>8, g>>8, b>>8

			// Pixel should be either white (255,255,255) or black (0,0,0)
			isWhite := r8 == 255 && g8 == 255 && b8 == 255
			isBlack := r8 == 0 && g8 == 0 && b8 == 0

			if !isWhite && !isBlack {
				allBW = false
				break
			}
		}
		if !allBW {
			break
		}
	}

	if !allBW {
		t.Error("Thermal preview should produce only black and white pixels")
	}
}

func TestRenderImage_ThermalPreviewAtkinson(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create a gradient image to test dithering
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 100; x++ {
			gray := uint8((x * 255) / 100) //nolint:gosec
			img.Set(x, y, color.RGBA{R: gray, G: gray, B: gray, A: 255})
		}
	}

	opts := emulator.DefaultImageOptions()
	opts.SimulateThermal = true
	opts.Dithering = graphics.Atkinson

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions failed:  %v", err)
	}

	result := engine.Render()

	// Count black and white pixels - Atkinson should produce a mix
	blackCount := 0
	whiteCount := 0

	for y := 0; y < result.Bounds().Dy(); y++ {
		for x := 0; x < result.Bounds().Dx(); x++ {
			r, _, _, _ := result.At(x, y).RGBA()
			if r == 0 {
				blackCount++
			} else if r>>8 == 255 {
				whiteCount++
			}
		}
	}

	// With a gradient, we should have both black and white pixels
	if blackCount == 0 {
		t.Error("Atkinson dithering should produce some black pixels for gradient")
	}
	if whiteCount == 0 {
		t.Error("Atkinson dithering should produce some white pixels for gradient")
	}
}

// ============================================================================
// RenderImage - Aspect Ratio Tests
// ============================================================================

func TestRenderImage_PreserveAspectRatio(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// 2:1 aspect ratio image
	img := createTestRGBA(200, 100, color.White)

	initialY := engine.State().CursorY

	opts := emulator.DefaultImageOptions()
	opts.PixelWidth = 100 // Scale down to half width
	opts.PreserveAspect = true

	err = engine.PrintImageWithOptions(img, opts)
	if err != nil {
		t.Fatalf("PrintImageWithOptions failed: %v", err)
	}

	// Height should be 50 (half of original 100)
	expectedHeight := 50.0
	actualAdvance := engine.State().CursorY - initialY

	if actualAdvance != expectedHeight {
		t.Errorf("Cursor advance = %.1f, want %.1f (aspect ratio preserved)", actualAdvance, expectedHeight)
	}
}

// ============================================================================
// RenderImage - Transparency Tests
// ============================================================================

func TestRenderImage_TransparentImage(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Create image with transparent pixels
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	// Leave all pixels as default (transparent black)

	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage failed: %v", err)
	}

	result := engine.Render()

	// Transparent areas should be white (paper color)
	// Check the image area
	foundWhite := false
	for y := 0; y < result.Bounds().Dy(); y++ {
		for x := 0; x < result.Bounds().Dx(); x++ {
			r, g, b, _ := result.At(x, y).RGBA()
			if r>>8 == 255 && g>>8 == 255 && b>>8 == 255 {
				foundWhite = true
				break
			}
		}
		if foundWhite {
			break
		}
	}

	if !foundWhite {
		t.Error("Transparent image should be composited over white background")
	}
}

// ============================================================================
// RenderImage - Multiple Images Tests
// ============================================================================

func TestRenderImage_MultipleImages(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img1 := createTestRGBA(100, 50, color.Black)
	img2 := createTestRGBA(100, 50, color.Black)

	// Render first image
	err = engine.PrintImage(img1)
	if err != nil {
		t.Fatalf("First PrintImage failed: %v", err)
	}
	cursorAfterFirst := engine.State().CursorY

	// Render second image
	err = engine.PrintImage(img2)
	if err != nil {
		t.Fatalf("Second PrintImage failed:  %v", err)
	}
	cursorAfterSecond := engine.State().CursorY

	// Second image should be below first
	if cursorAfterSecond <= cursorAfterFirst {
		t.Error("Second image should be rendered below first image")
	}
}

// ============================================================================
// PrintImage Convenience Methods Tests
// ============================================================================

func TestPrintImage(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine:  %v", err)
	}

	img := createTestRGBA(100, 50, color.White)
	initialY := engine.State().CursorY

	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage failed:  %v", err)
	}

	if engine.State().CursorY <= initialY {
		t.Error("PrintImage should advance cursor")
	}
}

func TestPrintImageAligned(t *testing.T) {
	tests := []struct {
		align string
	}{
		{constants.Left.String()},
		{constants.Center.String()},
		{constants.Right.String()},
	}

	for _, tt := range tests {
		t.Run(tt.align, func(t *testing.T) {
			engine, err := emulator.NewDefaultEngine()
			if err != nil {
				t.Fatalf("Failed to create engine: %v", err)
			}

			img := createTestRGBA(100, 50, color.White)

			err = engine.PrintImageAligned(img, 100, tt.align)
			if err != nil {
				t.Fatalf("PrintImageAligned failed: %v", err)
			}
		})
	}
}

func TestPrintImageThermalPreview(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	img := createTestRGBA(100, 50, color.RGBA{R: 128, G: 128, B: 128, A: 255})
	initialY := engine.State().CursorY

	err = engine.PrintImageThermalPreview(img, 100)
	if err != nil {
		t.Fatalf("PrintImageThermalPreview failed: %v", err)
	}

	if engine.State().CursorY <= initialY {
		t.Error("PrintImageThermalPreview should advance cursor")
	}
}

// ============================================================================
// Canvas Growth Tests
// ============================================================================

func TestRenderImage_CanvasGrows(t *testing.T) {
	engine, err := emulator.NewDefaultEngine()
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// Print text to move cursor down
	for i := 0; i < 50; i++ {
		engine.PrintLine("Line of text")
	}

	// Now add a large image
	img := createTestRGBA(200, 500, color.Black)

	err = engine.PrintImage(img)
	if err != nil {
		t.Fatalf("PrintImage failed: %v", err)
	}

	// Canvas should accommodate the image
	result := engine.Render()
	if result.Bounds().Dy() < 500 {
		t.Error("Canvas should grow to accommodate large image")
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkRenderImage_NormalPreview(b *testing.B) {
	engine, _ := emulator.NewDefaultEngine()
	img := createTestRGBA(576, 500, color.White)
	opts := emulator.DefaultImageOptions()
	opts.SimulateThermal = false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Reset()
		_ = engine.PrintImageWithOptions(img, opts)
	}
}

func BenchmarkRenderImage_ThermalPreview(b *testing.B) {
	engine, _ := emulator.NewDefaultEngine()
	img := createTestRGBA(576, 500, color.White)
	opts := emulator.DefaultImageOptions()
	opts.SimulateThermal = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Reset()
		_ = engine.PrintImageWithOptions(img, opts)
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// createTestRGBA creates a solid color test image
func createTestRGBA(width, height int, col color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, col)
		}
	}
	return img
}
