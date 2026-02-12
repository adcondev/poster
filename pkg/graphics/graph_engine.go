package graphics

import (
	"fmt"
	"image"
	"log"

	"github.com/adcondev/poster/pkg/constants"
)

// DitherMode defines how images are converted to monochrome
type DitherMode int

// TODO: Implement Floyd-Steinberg and Ordered dithering methods in future

const (
	// Threshold applies simple threshold conversion
	Threshold DitherMode = iota
	// Atkinson applies Atkinson dithering algorithm
	Atkinson
	// FloydSteinberg applies Floyd-Steinberg dithering (future)
	// FloydSteinberg
	// Ordered applies ordered dithering with Bayer matrix (future)
	// Ordered
)

// DitherMap maps constants to DitherMode
var DitherMap = map[string]DitherMode{
	constants.Threshold.String(): Threshold,
	constants.Atkinson.String():  Atkinson,
	// constants.DitheringFloydSteinberg: FloydSteinberg,
	// constants.DitheringOrdered:        Ordered,
}

// ScaleMode defines the scaling algorithm
type ScaleMode int

const (
	// NearestNeighbor uses nearest-neighbor scaling
	NearestNeighbor = iota
	// BiLinear uses bi-linear scaling
	BiLinear
	// BiCubic
)

// ScaleMap maps constants to ScaleMode
var ScaleMap = map[constants.Scaling]ScaleMode{
	constants.NearestNeighbor: NearestNeighbor,
	constants.Bilinear:        BiLinear,
	// constants.ScalingBicubic:  BiCubic,
}

// ImgOptions configures the graphics processing pipeline
type ImgOptions struct {
	PixelWidth     int        // Target width in pixels
	Threshold      uint8      // Threshold for black/white (0-255)
	Dithering      DitherMode // Processing algorithm
	Scaling        ScaleMode  // Up/Down Scale algorithm
	AutoRotate     bool       // Auto-rotate for best fit
	PreserveAspect bool       // Maintain aspect ratio
}

// DefaultOptions returns sensible defaults for 80mm printers
func DefaultOptions() *ImgOptions {
	return &ImgOptions{
		PixelWidth:     constants.DefaultImagePixelWidth,
		Threshold:      constants.DefaultImageThreshold,
		Dithering:      DitherMap[constants.Atkinson.String()],
		Scaling:        ScaleMap[constants.Bilinear],
		AutoRotate:     false,
		PreserveAspect: true,
	}
}

// Pipeline represents the image processing pipeline
type Pipeline struct {
	opts *ImgOptions
}

// NewPipeline creates a new processing pipeline with given options
func NewPipeline(opts *ImgOptions) *Pipeline {
	if opts == nil {
		opts = DefaultOptions()
	}
	return &Pipeline{opts: opts}
}

// Process transforms an image through the complete pipeline
func (p *Pipeline) Process(img image.Image) (*MonochromeBitmap, error) {
	if img == nil {
		return nil, fmt.Errorf("input image cannot be nil")
	}

	// Step 1: Resize if needed
	if p.opts.PixelWidth > 0 && img.Bounds().Dx() != p.opts.PixelWidth {
		img = p.resize(img)
	}

	// Step 2: Convert to grayscale
	gray := p.toGrayscale(img)

	// Step 3: Apply processing mode
	var mono *MonochromeBitmap
	switch p.opts.Dithering {
	case Atkinson:
		mono = p.applyAtkinson(gray)
	case Threshold:
		fallthrough
	default:
		mono = p.applyThreshold(gray)
	}

	return mono, nil
}

// TODO: Consider supporting other scaling algorithms (e.g., NN, Lanczos, Catmull-Rom) for even better quality or performance tuning.

// resize scales (up or down) the image to target width maintaining aspect ratio
func (p *Pipeline) resize(img image.Image) image.Image {
	// TODO: ImgOption should have a MaxWidth to cap resizing based on ticket size.
	maxWidth := constants.PaperPxWidth80mm
	if p.opts.PixelWidth > maxWidth {
		p.opts.PixelWidth = maxWidth
		log.Printf("resize: limiting target width to %d pixels", p.opts.PixelWidth)
	}

	return ResizeImage(img, &ResizeOptions{
		TargetWidth:    p.opts.PixelWidth,
		MaxWidth:       maxWidth,
		PreserveAspect: p.opts.PreserveAspect,
		Scaling:        p.opts.Scaling,
	})
}

// toGrayscale converts any image to grayscale
func (p *Pipeline) toGrayscale(img image.Image) *image.Gray {
	return ToGrayscale(img)
}

// applyThreshold applies simple threshold conversion
func (p *Pipeline) applyThreshold(gray *image.Gray) *MonochromeBitmap {
	bounds := gray.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	mono := NewMonochromeBitmap(width, height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := gray.GrayAt(x, y).Y
			// Set pixel to black (true) if below threshold
			if pixel < p.opts.Threshold {
				mono.SetPixel(x, y, true)
			}
		}
	}

	return mono
}

// TODO: Implement Floyd-Steinberg and Ordered dithering methods in future

// applyAtkinson implements Atkinson dithering algorithm
func (p *Pipeline) applyAtkinson(gray *image.Gray) *MonochromeBitmap {
	bounds := gray.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	mono := NewMonochromeBitmap(width, height)

	// Create a working copy for error diffusion
	work := make([]int, width*height)
	for y := 0; y < height; y++ {
		rowOffset := y * width
		for x := 0; x < width; x++ {
			work[rowOffset+x] = int(gray.GrayAt(x, y).Y)
		}
	}

	// Atkinson dithering pattern:
	//     *  1  1
	//  1  1  1
	//     1
	// Error is distributed as 1/8 to each neighbor (total 6/8 = 3/4)

	for y := 0; y < height; y++ {
		rowOffset := y * width
		nextRowOffset := (y + 1) * width
		nextNextRowOffset := (y + 2) * width

		for x := 0; x < width; x++ {
			idx := rowOffset + x
			oldPixel := work[idx]
			newPixel := 0
			if oldPixel > int(p.opts.Threshold) {
				newPixel = 255
			}

			// Set the monochrome pixel
			mono.SetPixel(x, y, newPixel == 0)

			// Calculate error
			err := oldPixel - newPixel

			// Atkinson only diffuses 3/4 (6/8) of the error
			// Each of the 6 neighbors gets 1/8 of the original error
			diffusedError := err / 8

			// Distribute to neighbors
			if x+1 < width {
				work[idx+1] += diffusedError
			}
			if x+2 < width {
				work[idx+2] += diffusedError
			}
			if y+1 < height {
				if x-1 >= 0 {
					work[nextRowOffset+(x-1)] += diffusedError
				}
				work[nextRowOffset+x] += diffusedError
				if x+1 < width {
					work[nextRowOffset+(x+1)] += diffusedError
				}
			}
			if y+2 < height {
				work[nextNextRowOffset+x] += diffusedError
			}
		}
	}

	return mono
}
