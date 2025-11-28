package builder

// ImageBuilder constructs image commands
type ImageBuilder struct {
	parent     *DocumentBuilder
	code       string
	pixelWidth int
	align      string
	threshold  int
	dithering  string
}

type imageCommand struct {
	Code       string `json:"code"`
	PixelWidth int    `json:"pixel_width,omitempty"`
	Align      string `json:"align,omitempty"`
	Threshold  int    `json:"threshold,omitempty"`
	Dithering  string `json:"dithering,omitempty"`
}

func newImageBuilder(parent *DocumentBuilder, base64Data string) *ImageBuilder {
	return &ImageBuilder{
		parent:     parent,
		code:       base64Data,
		pixelWidth: 128,
		align:      "center",
		threshold:  128,
		// TODO: make dithering enum
		dithering: "atkinson",
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
func (ib *ImageBuilder) Dithering(mode string) *ImageBuilder {
	ib.dithering = mode
	return ib
}

// Left aligns image to the left
func (ib *ImageBuilder) Left() *ImageBuilder {
	ib.align = "left"
	return ib
}

// Center centers the image (default)
func (ib *ImageBuilder) Center() *ImageBuilder {
	ib.align = "center"
	return ib
}

// Right aligns image to the right
func (ib *ImageBuilder) Right() *ImageBuilder {
	ib.align = "right"
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
	}
	return ib.parent.addCommand("image", cmd)
}
