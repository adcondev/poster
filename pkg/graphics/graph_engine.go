package graphics

import (
	"fmt"
	"image"

	"github.com/adcondev/poster/pkg/constants"
)

// DitherMode defines how images are converted to monochrome
type DitherMode int

const (
	// Threshold applies simple threshold conversion
	Threshold DitherMode = iota
	// Atkinson applies Atkinson dithering algorithm
	Atkinson
	// FloydSteinberg applies Floyd-Steinberg dithering algorithm
	FloydSteinberg
	// Ordered applies ordered dithering with a 4x4 Bayer matrix
	Ordered
)

// DitherMap maps constants to DitherMode
var DitherMap = map[string]DitherMode{
	constants.Threshold.String():      Threshold,
	constants.Atkinson.String():       Atkinson,
	constants.FloydSteinberg.String(): FloydSteinberg,
	constants.Ordered.String():        Ordered,
	"bayer":                           Ordered,
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
	MaxWidth       int        // Maximum width in pixels
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
		MaxWidth:       constants.PaperPxWidth80mm,
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
	case FloydSteinberg:
		mono = p.applyFloydSteinberg(gray)
	case Ordered:
		mono = p.applyOrdered(gray)
	case Threshold:
		fallthrough
	default:
		mono = p.applyThreshold(gray)
	}

	return mono, nil
}

// resize scales (up or down) the image to target width maintaining aspect ratio
func (p *Pipeline) resize(img image.Image) image.Image {
	maxWidth := p.opts.MaxWidth
	if maxWidth <= 0 {
		maxWidth = constants.PaperPxWidth80mm
	}
	if p.opts.PixelWidth > maxWidth {
		p.opts.PixelWidth = maxWidth
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
	if width <= 0 || height <= 0 {
		return mono
	}

	threshold := p.opts.Threshold
	bytesPerRow := mono.bytesPerRow
	fullBytes := width / 8
	remBits := width % 8

	bitMasks := [8]byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01}

	for y := 0; y < height; y++ {
		srcRowOffset := gray.PixOffset(bounds.Min.X, bounds.Min.Y+y)
		dstRowOffset := y * bytesPerRow

		for byteIdx := 0; byteIdx < fullBytes; byteIdx++ {
			var b byte
			srcX := srcRowOffset + byteIdx*8
			if gray.Pix[srcX] < threshold {
				b |= 0x80
			}
			if gray.Pix[srcX+1] < threshold {
				b |= 0x40
			}
			if gray.Pix[srcX+2] < threshold {
				b |= 0x20
			}
			if gray.Pix[srcX+3] < threshold {
				b |= 0x10
			}
			if gray.Pix[srcX+4] < threshold {
				b |= 0x08
			}
			if gray.Pix[srcX+5] < threshold {
				b |= 0x04
			}
			if gray.Pix[srcX+6] < threshold {
				b |= 0x02
			}
			if gray.Pix[srcX+7] < threshold {
				b |= 0x01
			}
			mono.data[dstRowOffset+byteIdx] = b
		}

		if remBits > 0 {
			var b byte
			srcX := srcRowOffset + fullBytes*8
			for bit := 0; bit < remBits; bit++ {
				if gray.Pix[srcX+bit] < threshold {
					b |= bitMasks[bit]
				}
			}
			mono.data[dstRowOffset+fullBytes] = b
		}
	}

	return mono
}

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

// applyFloydSteinberg implements Floyd-Steinberg error diffusion dithering.
func (p *Pipeline) applyFloydSteinberg(gray *image.Gray) *MonochromeBitmap {
	bounds := gray.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	mono := NewMonochromeBitmap(width, height)
	if width <= 0 || height <= 0 {
		return mono
	}

	work := make([]int, width*height)
	for y := 0; y < height; y++ {
		rowOffset := y * width
		for x := 0; x < width; x++ {
			work[rowOffset+x] = int(gray.GrayAt(bounds.Min.X+x, bounds.Min.Y+y).Y)
		}
	}

	threshold := int(p.opts.Threshold)

	for y := 0; y < height; y++ {
		rowOffset := y * width
		nextRowOffset := (y + 1) * width

		for x := 0; x < width; x++ {
			idx := rowOffset + x
			oldPixel := work[idx]
			newPixel := 0
			if oldPixel > threshold {
				newPixel = 255
			}

			mono.SetPixel(x, y, newPixel == 0)

			err := oldPixel - newPixel

			if x+1 < width {
				work[idx+1] += err * 7 / 16
			}
			if y+1 < height {
				if x-1 >= 0 {
					work[nextRowOffset+(x-1)] += err * 3 / 16
				}
				work[nextRowOffset+x] += err * 5 / 16
				if x+1 < width {
					work[nextRowOffset+(x+1)] += err * 1 / 16
				}
			}
		}
	}

	return mono
}

// applyOrdered implements ordered dithering using a 4x4 Bayer matrix.
func (p *Pipeline) applyOrdered(gray *image.Gray) *MonochromeBitmap {
	bounds := gray.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	mono := NewMonochromeBitmap(width, height)
	if width <= 0 || height <= 0 {
		return mono
	}

	bayer4x4 := [4][4]int{
		{0, 8, 2, 10},
		{12, 4, 14, 6},
		{3, 11, 1, 9},
		{15, 7, 13, 5},
	}

	baseThreshold := int(p.opts.Threshold)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			bayerVal := (bayer4x4[y%4][x%4]*16 + 8) - 128
			localThreshold := baseThreshold + bayerVal/2
			val := int(gray.GrayAt(bounds.Min.X+x, bounds.Min.Y+y).Y)
			if val < localThreshold {
				mono.SetPixel(x, y, true)
			}
		}
	}

	return mono
}

// SubBand extracts a horizontal band of the monochrome bitmap from startY with bandHeight.
func (m *MonochromeBitmap) SubBand(startY, bandHeight int) *MonochromeBitmap {
	if m == nil || startY < 0 || startY >= m.Height || bandHeight <= 0 {
		return NewMonochromeBitmap(0, 0)
	}
	if startY+bandHeight > m.Height {
		bandHeight = m.Height - startY
	}
	sub := NewMonochromeBitmap(m.Width, bandHeight)
	copy(sub.data, m.data[startY*m.bytesPerRow:(startY+bandHeight)*m.bytesPerRow])
	return sub
}
