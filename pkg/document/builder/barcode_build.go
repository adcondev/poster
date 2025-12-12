package builder

import (
	"github.com/adcondev/poster/pkg/constants"
)

// BarcodeBuilder constructs barcode commands
type BarcodeBuilder struct {
	parent      *DocumentBuilder
	symbology   string
	data        string
	width       *int
	height      *int
	hriPosition *string
	hriFont     *string
	align       *string
}

type barcodeCommand struct {
	Symbology   string  `json:"symbology"`
	Data        string  `json:"data"`
	Width       *int    `json:"width,omitempty"`
	Height      *int    `json:"height,omitempty"`
	HRIPosition *string `json:"hri_position,omitempty"`
	HRIFont     *string `json:"hri_font,omitempty"`
	Align       *string `json:"align,omitempty"`
}

func newBarcodeBuilder(parent *DocumentBuilder, symbology, data string) *BarcodeBuilder {
	return &BarcodeBuilder{
		parent:    parent,
		symbology: symbology,
		data:      data,
	}
}

// Width sets module width (2-6)
func (bb *BarcodeBuilder) Width(w int) *BarcodeBuilder {
	bb.width = &w
	return bb
}

// Height sets barcode height in dots (1-255)
func (bb *BarcodeBuilder) Height(h int) *BarcodeBuilder {
	bb.height = &h
	return bb
}

// HRIPosition sets HRI position (none, above, below, both)
func (bb *BarcodeBuilder) HRIPosition(pos string) *BarcodeBuilder {
	bb.hriPosition = &pos
	return bb
}

// HRIFont sets HRI font (A or B)
func (bb *BarcodeBuilder) HRIFont(font string) *BarcodeBuilder {
	bb.hriFont = &font
	return bb
}

// Left aligns barcode to the left
func (bb *BarcodeBuilder) Left() *BarcodeBuilder {
	align := constants.Left.String()
	bb.align = &align
	return bb
}

// Center centers the barcode (default)
func (bb *BarcodeBuilder) Center() *BarcodeBuilder {
	align := constants.Center.String()
	bb.align = &align
	return bb
}

// Right aligns barcode to the right
func (bb *BarcodeBuilder) Right() *BarcodeBuilder {
	align := constants.Right.String()
	bb.align = &align
	return bb
}

// Align sets barcode alignment (left, center, right)
func (bb *BarcodeBuilder) Align(align constants.Alignment) *BarcodeBuilder {
	alignStr := align.String()
	bb.align = &alignStr
	return bb
}

// End finishes the barcode command
func (bb *BarcodeBuilder) End() *DocumentBuilder {
	cmd := barcodeCommand{
		Symbology:   bb.symbology,
		Data:        bb.data,
		Width:       bb.width,
		Height:      bb.height,
		HRIPosition: bb.hriPosition,
		HRIFont:     bb.hriFont,
		Align:       bb.align,
	}
	return bb.parent.addCommand("barcode", cmd)
}
