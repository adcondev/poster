package emulator

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/graphics"
)

// ImageOptions configures how an image is rendered on the emulated receipt
type ImageOptions struct {
	// PixelWidth is the target width (0 = use original, capped to paper width)
	PixelWidth int

	// Align uses standard alignment constants ("left", "center", "right")
	Align string

	// PreserveAspect maintains aspect ratio when resizing
	PreserveAspect bool

	// Scaling algorithm to use
	Scaling graphics.ScaleMode

	// SimulateThermal applies full thermal print processing for accurate preview.
	// When true, the image goes through the same Pipeline as physical printing,
	// ensuring 100% fidelity between preview and actual output.
	SimulateThermal bool

	// Dithering mode (only used when SimulateThermal is true)
	Dithering graphics.DitherMode

	// Threshold for dithering (only used when SimulateThermal is true)
	Threshold uint8
}

// DefaultImageOptions returns sensible defaults for image embedding
func DefaultImageOptions() *ImageOptions {
	return &ImageOptions{
		PixelWidth:      constants.DefaultImagePixelWidth,
		Align:           constants.Center.String(),
		PreserveAspect:  true,
		Scaling:         graphics.BiLinear,
		SimulateThermal: false,
		Dithering:       graphics.Atkinson,
		Threshold:       128,
	}
}

// ImageRenderer handles image rendering for the emulator
type ImageRenderer struct {
	canvas *DynamicCanvas
	state  *PrinterState
}

// NewImageRenderer creates a new ImageRenderer
func NewImageRenderer(canvas *DynamicCanvas, state *PrinterState) *ImageRenderer {
	return &ImageRenderer{
		canvas: canvas,
		state:  state,
	}
}

// RenderImage embeds an image onto the emulated receipt canvas
func (ir *ImageRenderer) RenderImage(img image.Image, opts *ImageOptions) error {
	if img == nil {
		return fmt.Errorf("image cannot be nil")
	}

	if opts == nil {
		opts = DefaultImageOptions()
	}

	// Determine target width
	targetWidth := opts.PixelWidth
	if targetWidth <= 0 {
		// Use original image width, capped to paper width later
		targetWidth = img.Bounds().Dx()
	}
	if targetWidth > ir.state.PaperPxWidth {
		targetWidth = ir.state.PaperPxWidth
	}

	var processed image.Image

	if opts.SimulateThermal {
		// Use full graphics Pipeline for 100% fidelity thermal preview
		processed = ir.processThermalPreview(img, targetWidth, opts)
	} else {
		// Normal preview: just resize (keep colors)
		processed = ir.processNormalPreview(img, targetWidth, opts)
	}

	// Calculate position based on alignment
	imgBounds := processed.Bounds()
	imgWidth := imgBounds.Dx()
	imgHeight := imgBounds.Dy()

	startX := ir.calculateAlignedX(imgWidth, opts.Align)
	startY := int(ir.state.CursorY)

	// Ensure canvas has enough height
	requiredY := float64(startY + imgHeight)
	ir.canvas.EnsureHeight(requiredY)

	// Draw image onto canvas
	dstRect := image.Rect(startX, startY, startX+imgWidth, startY+imgHeight)
	draw.Draw(ir.canvas.Image(), dstRect, processed, imgBounds.Min, draw.Over)

	// Update cursor position
	ir.state.CursorY = float64(startY + imgHeight)
	ir.canvas.UpdateMaxY(ir.state.CursorY)

	return nil
}

// processNormalPreview resizes the image while preserving colors/grayscale
func (ir *ImageRenderer) processNormalPreview(img image.Image, targetWidth int, opts *ImageOptions) image.Image {
	// Composite over white to handle transparency
	composited := graphics.CompositeOverWhite(img)

	// Resize using extracted utility
	return graphics.ResizeImage(composited, &graphics.ResizeOptions{
		TargetWidth:    targetWidth,
		MaxWidth:       ir.state.PaperPxWidth,
		PreserveAspect: opts.PreserveAspect,
		Scaling:        opts.Scaling,
	})
}

// processThermalPreview uses the full graphics.Pipeline for accurate thermal output preview
func (ir *ImageRenderer) processThermalPreview(img image.Image, targetWidth int, opts *ImageOptions) image.Image {
	// Configure pipeline exactly as the printer service would
	pipelineOpts := &graphics.ImgOptions{
		PixelWidth:     targetWidth,
		Threshold:      opts.Threshold,
		Dithering:      opts.Dithering,
		Scaling:        opts.Scaling,
		PreserveAspect: opts.PreserveAspect,
	}

	pipeline := graphics.NewPipeline(pipelineOpts)

	// Process through the same pipeline the printer uses
	bitmap, err := pipeline.Process(img)
	if err != nil {
		// Fallback to normal preview on error
		return ir.processNormalPreview(img, targetWidth, opts)
	}

	// Convert MonochromeBitmap back to RGBA for canvas drawing
	// This ensures 100% fidelity - what you see is what will print
	return bitmap.ToImage()
}

// calculateAlignedX returns the X position for the given alignment string
func (ir *ImageRenderer) calculateAlignedX(imageWidth int, align string) int {
	switch align {
	case constants.Center.String():
		return (ir.state.PaperPxWidth - imageWidth) / 2
	case constants.Right.String():
		return ir.state.PaperPxWidth - imageWidth
	default: // constants.Left.String() or fallback
		return 0
	}
}
