package graphics_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/graphics"
)

// ============================================================================
// ResizeImage Tests
// ============================================================================

func TestResizeImage_NilImage(t *testing.T) {
	result := graphics.ResizeImage(nil, nil)
	if result != nil {
		t.Error("ResizeImage(nil, nil) should return nil")
	}
}

func TestResizeImage_NilOptions(t *testing.T) {
	img := createTestImage(100, 100, color.White)
	result := graphics.ResizeImage(img, nil)

	if result == nil {
		t.Fatal("ResizeImage with nil options should not return nil")
	}

	// Should use default options (80mm paper width = 576px max)
	bounds := result.Bounds()
	if bounds.Dx() > constants.PaperPxWidth80mm {
		t.Errorf("Result width %d exceeds max width %d", bounds.Dx(), constants.PaperPxWidth80mm)
	}
}

func TestResizeImage_NoResizeNeeded(t *testing.T) {
	targetWidth := 200
	img := createTestImage(targetWidth, 100, color.White)

	opts := &graphics.ResizeOptions{
		TargetWidth:    targetWidth,
		MaxWidth:       constants.PaperPxWidth80mm,
		PreserveAspect: true,
		Scaling:        graphics.BiLinear,
	}

	result := graphics.ResizeImage(img, opts)

	// Should return original image when already at target width
	if result.Bounds().Dx() != targetWidth {
		t.Errorf("Expected width %d, got %d", targetWidth, result.Bounds().Dx())
	}
}

func TestResizeImage_ScaleDown(t *testing.T) {
	img := createTestImage(400, 200, color.White)

	opts := &graphics.ResizeOptions{
		TargetWidth:    200,
		MaxWidth:       constants.PaperPxWidth80mm,
		PreserveAspect: true,
		Scaling:        graphics.BiLinear,
	}

	result := graphics.ResizeImage(img, opts)
	bounds := result.Bounds()

	if bounds.Dx() != 200 {
		t.Errorf("Expected width 200, got %d", bounds.Dx())
	}
	// Aspect ratio preserved:  400: 200 = 2:1, so 200:100
	if bounds.Dy() != 100 {
		t.Errorf("Expected height 100 (aspect ratio preserved), got %d", bounds.Dy())
	}
}

func TestResizeImage_ScaleUp(t *testing.T) {
	img := createTestImage(100, 50, color.White)

	opts := &graphics.ResizeOptions{
		TargetWidth:    200,
		MaxWidth:       constants.PaperPxWidth80mm,
		PreserveAspect: true,
		Scaling:        graphics.BiLinear,
	}

	result := graphics.ResizeImage(img, opts)
	bounds := result.Bounds()

	if bounds.Dx() != 200 {
		t.Errorf("Expected width 200, got %d", bounds.Dx())
	}
	// Aspect ratio preserved: 100:50 = 2:1, so 200:100
	if bounds.Dy() != 100 {
		t.Errorf("Expected height 100 (aspect ratio preserved), got %d", bounds.Dy())
	}
}

func TestResizeImage_MaxWidthCap(t *testing.T) {
	img := createTestImage(1000, 500, color.White)

	opts := &graphics.ResizeOptions{
		TargetWidth:    800, // Exceeds max
		MaxWidth:       576,
		PreserveAspect: true,
		Scaling:        graphics.BiLinear,
	}

	result := graphics.ResizeImage(img, opts)
	bounds := result.Bounds()

	if bounds.Dx() != 576 {
		t.Errorf("Expected width capped at 576, got %d", bounds.Dx())
	}
}

func TestResizeImage_NoAspectPreservation(t *testing.T) {
	img := createTestImage(400, 200, color.White)

	opts := &graphics.ResizeOptions{
		TargetWidth:    200,
		MaxWidth:       constants.PaperPxWidth80mm,
		PreserveAspect: false,
		Scaling:        graphics.BiLinear,
	}

	result := graphics.ResizeImage(img, opts)
	bounds := result.Bounds()

	if bounds.Dx() != 200 {
		t.Errorf("Expected width 200, got %d", bounds.Dx())
	}
	// Without aspect preservation, height stays the same
	if bounds.Dy() != 200 {
		t.Errorf("Expected height 200 (no aspect preservation), got %d", bounds.Dy())
	}
}

func TestResizeImage_ZeroTargetWidth(t *testing.T) {
	img := createTestImage(300, 150, color.White)

	opts := &graphics.ResizeOptions{
		TargetWidth:    0, // Invalid
		MaxWidth:       constants.PaperPxWidth80mm,
		PreserveAspect: true,
		Scaling:        graphics.BiLinear,
	}

	result := graphics.ResizeImage(img, opts)
	bounds := result.Bounds()

	// Should use original width when target is 0
	if bounds.Dx() != 300 {
		t.Errorf("Expected original width 300 when target is 0, got %d", bounds.Dx())
	}
}

func TestResizeImage_NegativeTargetWidth(t *testing.T) {
	img := createTestImage(300, 150, color.White)

	opts := &graphics.ResizeOptions{
		TargetWidth:    -100, // Invalid
		MaxWidth:       constants.PaperPxWidth80mm,
		PreserveAspect: true,
		Scaling:        graphics.BiLinear,
	}

	result := graphics.ResizeImage(img, opts)
	bounds := result.Bounds()

	// Should use original width when target is negative
	if bounds.Dx() != 300 {
		t.Errorf("Expected original width 300 when target is negative, got %d", bounds.Dx())
	}
}

func TestResizeImage_NearestNeighborScaling(t *testing.T) {
	img := createTestImage(100, 100, color.Black)

	opts := &graphics.ResizeOptions{
		TargetWidth:    50,
		MaxWidth:       constants.PaperPxWidth80mm,
		PreserveAspect: true,
		Scaling:        graphics.NearestNeighbor,
	}

	result := graphics.ResizeImage(img, opts)
	bounds := result.Bounds()

	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("Expected 50x50, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestResizeImage_VariousAspectRatios(t *testing.T) {
	tests := []struct {
		name       string
		srcW, srcH int
		targetW    int
		expectW    int
		expectH    int
	}{
		{"Square", 100, 100, 50, 50, 50},
		{"Wide 2: 1", 200, 100, 100, 100, 50},
		{"Tall 1:2", 100, 200, 50, 50, 100},
		{"Wide 4:3", 400, 300, 200, 200, 150},
		{"Tall 3:4", 300, 400, 150, 150, 200},
		{"Very wide 10:1", 1000, 100, 500, 500, 50},
		{"Very tall 1:10", 100, 1000, 50, 50, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := createTestImage(tt.srcW, tt.srcH, color.White)
			opts := &graphics.ResizeOptions{
				TargetWidth:    tt.targetW,
				MaxWidth:       constants.PaperPxWidth80mm,
				PreserveAspect: true,
				Scaling:        graphics.BiLinear,
			}

			result := graphics.ResizeImage(img, opts)
			bounds := result.Bounds()

			if bounds.Dx() != tt.expectW {
				t.Errorf("Width:  expected %d, got %d", tt.expectW, bounds.Dx())
			}
			if bounds.Dy() != tt.expectH {
				t.Errorf("Height:  expected %d, got %d", tt.expectH, bounds.Dy())
			}
		})
	}
}

// ============================================================================
// DefaultResizeOptions Tests
// ============================================================================

func TestDefaultResizeOptions(t *testing.T) {
	opts := graphics.DefaultResizeOptions(384)

	if opts.TargetWidth != 384 {
		t.Errorf("TargetWidth = %d, want 384", opts.TargetWidth)
	}
	if opts.MaxWidth != constants.PaperPxWidth80mm {
		t.Errorf("MaxWidth = %d, want %d", opts.MaxWidth, constants.PaperPxWidth80mm)
	}
	if !opts.PreserveAspect {
		t.Error("PreserveAspect should be true by default")
	}
	if opts.Scaling != graphics.BiLinear {
		t.Errorf("Scaling = %d, want BiLinear", opts.Scaling)
	}
}

// ============================================================================
// ToGrayscale Tests
// ============================================================================

func TestToGrayscale_WhiteImage(t *testing.T) {
	img := createTestImage(10, 10, color.White)
	gray := graphics.ToGrayscale(img)

	if gray == nil {
		t.Fatal("ToGrayscale returned nil")
	}

	// All pixels should be white (255)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			g := gray.GrayAt(x, y).Y
			if g != 255 {
				t.Errorf("Pixel (%d,%d) = %d, want 255", x, y, g)
			}
		}
	}
}

func TestToGrayscale_BlackImage(t *testing.T) {
	img := createTestImage(10, 10, color.Black)
	gray := graphics.ToGrayscale(img)

	if gray == nil {
		t.Fatal("ToGrayscale returned nil")
	}

	// All pixels should be black (0)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			g := gray.GrayAt(x, y).Y
			if g != 0 {
				t.Errorf("Pixel (%d,%d) = %d, want 0", x, y, g)
			}
		}
	}
}

func TestToGrayscale_ColorImage(t *testing.T) {
	// Create image with different colors
	img := image.NewRGBA(image.Rect(0, 0, 3, 1))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255}) // Red
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Green
	img.Set(2, 0, color.RGBA{R: 0, G: 0, B: 255, A: 255}) // Blue

	gray := graphics.ToGrayscale(img)

	// Each color should convert to different gray values
	// Standard luminance:  Y = 0.299*R + 0.587*G + 0.114*B
	redGray := gray.GrayAt(0, 0).Y
	greenGray := gray.GrayAt(1, 0).Y
	blueGray := gray.GrayAt(2, 0).Y

	// Green should be brightest due to luminance weights
	if greenGray <= redGray || greenGray <= blueGray {
		t.Errorf("Green (%d) should be brightest.  Red: %d, Blue: %d", greenGray, redGray, blueGray)
	}
	// Blue should be darkest
	if blueGray >= redGray || blueGray >= greenGray {
		t.Errorf("Blue (%d) should be darkest. Red: %d, Green: %d", blueGray, redGray, greenGray)
	}
}

func TestToGrayscale_TransparentPixels(t *testing.T) {
	// Create image with transparent pixels
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.Set(0, 0, color.RGBA{R: 0, G: 0, B: 0, A: 0})   // Fully transparent
	img.Set(1, 0, color.RGBA{R: 0, G: 0, B: 0, A: 255}) // Fully opaque black

	gray := graphics.ToGrayscale(img)

	// Transparent pixel should become white (composited over white background)
	transparentGray := gray.GrayAt(0, 0).Y
	opaqueGray := gray.GrayAt(1, 0).Y

	if transparentGray != 255 {
		t.Errorf("Transparent pixel should be white (255), got %d", transparentGray)
	}
	if opaqueGray != 0 {
		t.Errorf("Opaque black pixel should be black (0), got %d", opaqueGray)
	}
}

func TestToGrayscale_SemiTransparent(t *testing.T) {
	// Create semi-transparent black pixel
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 0, G: 0, B: 0, A: 128}) // 50% transparent black

	gray := graphics.ToGrayscale(img)

	// Should be approximately middle gray (composited over white)
	g := gray.GrayAt(0, 0).Y
	// 50% black over white ≈ 128 gray (may vary slightly due to alpha blending)
	if g < 100 || g > 160 {
		t.Errorf("Semi-transparent black should be ~128 gray, got %d", g)
	}
}

func TestToGrayscale_PreservesDimensions(t *testing.T) {
	tests := []struct {
		width, height int
	}{
		{1, 1},
		{10, 10},
		{100, 50},
		{50, 100},
		{576, 1000},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			img := createTestImage(tt.width, tt.height, color.White)
			gray := graphics.ToGrayscale(img)

			bounds := gray.Bounds()
			if bounds.Dx() != tt.width || bounds.Dy() != tt.height {
				t.Errorf("Expected %dx%d, got %dx%d", tt.width, tt.height, bounds.Dx(), bounds.Dy())
			}
		})
	}
}

// ============================================================================
// CompositeOverWhite Tests
// ============================================================================

func TestCompositeOverWhite_OpaqueImage(t *testing.T) {
	img := createTestImage(10, 10, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	result := graphics.CompositeOverWhite(img)

	if result == nil {
		t.Fatal("CompositeOverWhite returned nil")
	}

	// All pixels should still be red
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			r, g, b, a := result.At(x, y).RGBA()
			// RGBA returns 16-bit values, convert to 8-bit
			if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 || a>>8 != 255 {
				t.Errorf("Pixel (%d,%d) color changed unexpectedly", x, y)
			}
		}
	}
}

func TestCompositeOverWhite_FullyTransparent(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	// All pixels are fully transparent (default zero value)

	result := graphics.CompositeOverWhite(img)

	// All pixels should be white
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			r, g, b, a := result.At(x, y).RGBA()
			if r>>8 != 255 || g>>8 != 255 || b>>8 != 255 || a>>8 != 255 {
				t.Errorf("Transparent pixel (%d,%d) should be white, got RGBA(%d,%d,%d,%d)",
					x, y, r>>8, g>>8, b>>8, a>>8)
			}
		}
	}
}

func TestCompositeOverWhite_SemiTransparent(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	// 50% transparent red
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 128})

	result := graphics.CompositeOverWhite(img)

	r, g, b, a := result.At(0, 0).RGBA()
	// Convert from 16-bit to 8-bit
	r8, g8, b8, a8 := r>>8, g>>8, b>>8, a>>8

	// Red channel:  should be higher than G and B (it's red over white = pink)
	if r8 < g8 || r8 < b8 {
		t.Errorf("Red channel (%d) should be >= green (%d) and blue (%d) for pink", r8, g8, b8)
	}

	// Green and Blue should be similar (both were 0, blended with white)
	diff := int(g8) - int(b8)
	if diff < -10 || diff > 10 {
		t.Errorf("Green (%d) and Blue (%d) should be similar", g8, b8)
	}

	// Result should be fully opaque
	if a8 != 255 {
		t.Errorf("Alpha should be 255 after compositing, got %d", a8)
	}

	// The result should not be pure white or pure red
	if r8 == 255 && g8 == 255 && b8 == 255 {
		t.Error("Result should not be pure white")
	}
	if g8 == 0 && b8 == 0 {
		t.Error("Result should not be pure red (G and B should have some white blended in)")
	}
}

func TestCompositeOverWhite_PreservesDimensions(t *testing.T) {
	tests := []struct {
		width, height int
	}{
		{1, 1},
		{10, 10},
		{100, 50},
		{576, 1000},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			img := createTestImage(tt.width, tt.height, color.White)
			result := graphics.CompositeOverWhite(img)

			bounds := result.Bounds()
			if bounds.Dx() != tt.width || bounds.Dy() != tt.height {
				t.Errorf("Expected %dx%d, got %dx%d", tt.width, tt.height, bounds.Dx(), bounds.Dy())
			}
		})
	}
}

func TestCompositeOverWhite_MixedTransparency(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4, 1))
	img.Set(0, 0, color.RGBA{R: 0, G: 0, B: 0, A: 0})   // Fully transparent
	img.Set(1, 0, color.RGBA{R: 0, G: 0, B: 0, A: 64})  // 25% opaque
	img.Set(2, 0, color.RGBA{R: 0, G: 0, B: 0, A: 128}) // 50% opaque
	img.Set(3, 0, color.RGBA{R: 0, G: 0, B: 0, A: 255}) // Fully opaque

	result := graphics.CompositeOverWhite(img)

	// Get gray values (all channels should be equal for black over white)
	g0, _, _, _ := result.At(0, 0).RGBA()
	g1, _, _, _ := result.At(1, 0).RGBA()
	g2, _, _, _ := result.At(2, 0).RGBA()
	g3, _, _, _ := result.At(3, 0).RGBA()

	// Should be in descending order of brightness
	// More transparent = more white = brighter
	if g0>>8 <= g1>>8 || g1>>8 <= g2>>8 || g2>>8 <= g3>>8 {
		t.Errorf("Brightness should decrease with opacity:  %d > %d > %d > %d",
			g0>>8, g1>>8, g2>>8, g3>>8)
	}

	// Fully transparent should be white (255)
	if g0>>8 != 255 {
		t.Errorf("Fully transparent should be white (255), got %d", g0>>8)
	}
	// Fully opaque black should be black (0)
	if g3>>8 != 0 {
		t.Errorf("Fully opaque black should be black (0), got %d", g3>>8)
	}
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestProcessingPipeline_ResizeAndGrayscale(t *testing.T) {
	// Create a color image
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 200; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 128, A: 255}) //nolint:gosec
		}
	}

	// Resize
	resized := graphics.ResizeImage(img, &graphics.ResizeOptions{
		TargetWidth:    100,
		MaxWidth:       576,
		PreserveAspect: true,
		Scaling:        graphics.BiLinear,
	})

	// Convert to grayscale
	gray := graphics.ToGrayscale(resized)

	// Verify dimensions
	bounds := gray.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 50 {
		t.Errorf("Expected 100x50, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestProcessingPipeline_TransparentPNGHandling(t *testing.T) {
	// Simulate a transparent PNG logo
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Create a circle-like pattern with transparency
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			dx, dy := x-50, y-50
			if dx*dx+dy*dy < 40*40 {
				// Inside circle:  opaque black
				img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			} else {
				// Outside circle: transparent
				img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0})
			}
		}
	}

	// Process through pipeline
	composited := graphics.CompositeOverWhite(img)
	resized := graphics.ResizeImage(composited, &graphics.ResizeOptions{
		TargetWidth:    50,
		MaxWidth:       576,
		PreserveAspect: true,
		Scaling:        graphics.BiLinear,
	})
	gray := graphics.ToGrayscale(resized)

	// Verify:  corners should be white (was transparent), center should be dark
	cornerGray := gray.GrayAt(0, 0).Y
	centerGray := gray.GrayAt(25, 25).Y

	if cornerGray != 255 {
		t.Errorf("Corner (transparent area) should be white (255), got %d", cornerGray)
	}
	if centerGray > 50 {
		t.Errorf("Center (opaque black) should be dark, got %d", centerGray)
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkResizeImage_BiLinear(b *testing.B) {
	img := createTestImage(576, 800, color.White)
	opts := &graphics.ResizeOptions{
		TargetWidth:    384,
		MaxWidth:       576,
		PreserveAspect: true,
		Scaling:        graphics.BiLinear,
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = graphics.ResizeImage(img, opts)
	}
}

func BenchmarkResizeImage_NearestNeighbor(b *testing.B) {
	img := createTestImage(576, 800, color.White)
	opts := &graphics.ResizeOptions{
		TargetWidth:    384,
		MaxWidth:       576,
		PreserveAspect: true,
		Scaling:        graphics.NearestNeighbor,
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = graphics.ResizeImage(img, opts)
	}
}

func BenchmarkToGrayscale(b *testing.B) {
	img := createTestImage(384, 500, color.White)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = graphics.ToGrayscale(img)
	}
}

func BenchmarkCompositeOverWhite(b *testing.B) {
	img := createTestImage(384, 500, color.RGBA{R: 128, G: 128, B: 128, A: 128})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = graphics.CompositeOverWhite(img)
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

// createTestImage creates a solid color test image
func createTestImage(width, height int, col color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, col)
		}
	}
	return img
}
