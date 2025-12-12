package builder

import (
	"github.com/adcondev/poster/pkg/constants"
)

// ImageBuilder constructs image commands
type ImageBuilder struct {
	parent     *DocumentBuilder
	code       string
	pixelWidth int
	align      string
	threshold  int
	dithering  string
	scaling    string
}

type imageCommand struct {
	Code       string `json:"code"`
	PixelWidth int    `json:"pixel_width,omitempty"`
	Align      string `json:"align,omitempty"`
	Threshold  int    `json:"threshold,omitempty"`
	Dithering  string `json:"dithering,omitempty"`
	Scaling    string `json:"scaling,omitempty"`
}

func newImageBuilder(parent *DocumentBuilder, base64Data string) *ImageBuilder {
	return &ImageBuilder{
		parent:     parent,
		code:       base64Data,
		pixelWidth: constants.DefaultImagePixelWidth,
		threshold:  constants.DefaultImageThreshold,
		align:      constants.Center.String(),
		dithering:  constants.Atkinson.String(),
		scaling:    constants.Bilinear.String(),
	}
}

// Width sets image width in pixels
func (ib *ImageBuilder) Width(pixels int) *ImageBuilder {
	ib.pixelWidth = pixels
	return ib
}

// Threshold sets B/W threshold (0-255)
func (ib *ImageBuilder) Threshold(t int) *ImageBuilder {
	ib.threshold = t
	return ib
}

// Dithering sets dithering algorithm ("threshold" or "atkinson")
func (ib *ImageBuilder) Dithering(mode constants.Dithering) *ImageBuilder {
	ib.dithering = mode.String()
	return ib
}

// Scaling sets scaling algorithm ("nns" or "bilinear")
func (ib *ImageBuilder) Scaling(mode constants.Scaling) *ImageBuilder {
	ib.scaling = mode.String()
	return ib
}

// Left aligns image to the left
func (ib *ImageBuilder) Left() *ImageBuilder {
	ib.align = constants.Left.String()
	return ib
}

// Center centers the image (default)
func (ib *ImageBuilder) Center() *ImageBuilder {
	ib.align = constants.Center.String()
	return ib
}

// Right aligns image to the right
func (ib *ImageBuilder) Right() *ImageBuilder {
	ib.align = constants.Right.String()
	return ib
}

// Align sets image alignment (left, center, right)
func (ib *ImageBuilder) Align(align constants.Alignment) *ImageBuilder {
	ib.align = align.String()
	return ib
}

// End finishes the image command
func (ib *ImageBuilder) End() *DocumentBuilder {
	cmd := imageCommand{
		Code:       ib.code,
		PixelWidth: ib.pixelWidth,
		Align:      ib.align,
		Threshold:  ib.threshold,
		Dithering:  ib.dithering,
		Scaling:    ib.scaling,
	}
	return ib.parent.addCommand("image", cmd)
}
