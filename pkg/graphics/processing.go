package graphics

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"

	"github.com/adcondev/poster/pkg/constants"
)

// ResizeOptions configures image resizing behavior
type ResizeOptions struct {
	TargetWidth    int
	MaxWidth       int // Safety cap (e.g., 576 for 80mm)
	PreserveAspect bool
	Scaling        ScaleMode
}

// DefaultResizeOptions returns sensible defaults for thermal printer emulation
func DefaultResizeOptions(paperWidth int) *ResizeOptions {
	return &ResizeOptions{
		TargetWidth:    paperWidth,
		MaxWidth:       constants.PaperPxWidth80mm,
		PreserveAspect: true,
		Scaling:        BiLinear,
	}
}

// ResizeImage scales an image to the target width while optionally preserving aspect ratio.
func ResizeImage(img image.Image, opts *ResizeOptions) image.Image {
	if img == nil {
		return nil
	}
	if opts == nil {
		opts = DefaultResizeOptions(constants.PaperPxWidth80mm)
	}

	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()

	// Apply max width cap
	targetW := opts.TargetWidth
	if targetW > opts.MaxWidth {
		targetW = opts.MaxWidth
	}
	if targetW <= 0 {
		targetW = srcW // No resize if invalid
	}

	// Skip if already at target width
	if srcW == targetW {
		return img
	}

	targetH := srcH
	if opts.PreserveAspect && srcW > 0 {
		targetH = (srcH * targetW) / srcW
	}

	// Create destination image
	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))

	var scaler draw.Scaler
	switch opts.Scaling {
	case NearestNeighbor:
		scaler = draw.NearestNeighbor
	case BiLinear:
		scaler = draw.BiLinear
	default:
		scaler = draw.BiLinear
	}

	scaler.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
	return dst
}

// ToGrayscale converts any image to grayscale, compositing transparent pixels
// over white to match thermal paper behavior.
func ToGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()

	// 1. Composite over white first (handles transparency)
	rgba := CompositeOverWhite(img)

	// 2. Convert to Grayscale
	gray := image.NewGray(bounds)
	width, height := bounds.Dx(), bounds.Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := color.GrayModel.Convert(rgba.At(x+bounds.Min.X, y+bounds.Min.Y)).(color.Gray)
			gray.Set(x+bounds.Min.X, y+bounds.Min.Y, c)
		}
	}

	return gray
}

// CompositeOverWhite composites an image over a white background.
// Useful for handling transparent PNGs before any processing.
func CompositeOverWhite(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	// Fill with white (thermal paper color)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	draw.Draw(rgba, bounds, &image.Uniform{C: white}, image.Point{}, draw.Src)

	// Draw source over white
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Over)

	return rgba
}
