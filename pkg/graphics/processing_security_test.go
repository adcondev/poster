package graphics_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/adcondev/poster/pkg/graphics"
)

// HugeImage is a mock image that reports huge dimensions but doesn't store data.
type HugeImage struct {
	W, H int
}

func (m *HugeImage) ColorModel() color.Model { return color.RGBAModel }
func (m *HugeImage) Bounds() image.Rectangle { return image.Rect(0, 0, m.W, m.H) }
func (m *HugeImage) At(x, y int) color.Color { return color.RGBA{0, 0, 0, 255} }

func TestResizeImage_DoS_Vulnerability(t *testing.T) {
	// Setup: srcW=10, srcH=1200, targetW=500.
	// targetH would be (1200 * 500) / 10 = 60,000 without cap.
	// With cap, it should be graphics.MaxImagePixelHeight (32000).

	src := &HugeImage{W: 10, H: 1200}
	opts := &graphics.ResizeOptions{
		TargetWidth:    500,
		MaxWidth:       576,
		PreserveAspect: true,
		Scaling:        graphics.BiLinear,
	}

	res := graphics.ResizeImage(src, opts)

	if res == nil {
		t.Fatal("ResizeImage returned nil")
	}

	bounds := res.Bounds()
	t.Logf("Result bounds: %v", bounds)

	if bounds.Dy() > graphics.MaxImagePixelHeight {
		t.Errorf("Result height %d exceeds MaxImagePixelHeight %d", bounds.Dy(), graphics.MaxImagePixelHeight)
	}

	if bounds.Dy() != graphics.MaxImagePixelHeight {
		t.Errorf("Expected result height to be capped at %d, got %d", graphics.MaxImagePixelHeight, bounds.Dy())
	}
}
