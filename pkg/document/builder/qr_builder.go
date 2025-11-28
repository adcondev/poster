package builder

import (
	"log"
)

// QRBuilder constructs QR code commands
type QRBuilder struct {
	parent      *DocumentBuilder
	data        string
	humanText   string
	pixelWidth  int
	correction  string
	align       string
	logo        string
	circleShape bool
}

type qrCommand struct {
	Data        string `json:"data"`
	HumanText   string `json:"human_text,omitempty"`
	PixelWidth  int    `json:"pixel_width,omitempty"`
	Correction  string `json:"correction,omitempty"`
	Align       string `json:"align,omitempty"`
	Logo        string `json:"logo,omitempty"`
	CircleShape bool   `json:"circle_shape,omitempty"`
}

func newQRBuilder(parent *DocumentBuilder, data string) *QRBuilder {
	return &QRBuilder{
		parent:     parent,
		data:       data,
		pixelWidth: 128,
		correction: "Q",
		align:      Center.String(),
	}
}

// Size sets the QR code pixel width
func (qb *QRBuilder) Size(pixels int) *QRBuilder {
	qb.pixelWidth = pixels
	return qb
}

// Correction sets error correction level (L, M, Q, H)
func (qb *QRBuilder) Correction(level string) *QRBuilder {
	qb.correction = level
	return qb
}

// Left aligns QR to the left
func (qb *QRBuilder) Left() *QRBuilder {
	qb.align = Left.String()
	return qb
}

// Center centers the QR (default)
func (qb *QRBuilder) Center() *QRBuilder {
	qb.align = Center.String()
	return qb
}

// Right aligns QR to the right
func (qb *QRBuilder) Right() *QRBuilder {
	qb.align = Right.String()
	return qb
}

// WithText adds human-readable text below QR
func (qb *QRBuilder) WithText(text string) *QRBuilder {
	qb.humanText = text
	return qb
}

// TODO: Add WithLogoFile
// Load logo from image(jpg, png) or txt file with base64 content

// WithLogo adds a logo in the center (base64)
func (qb *QRBuilder) WithLogo(base64Logo string) *QRBuilder {
	qb.logo = base64Logo
	return qb
}

// CircleShape uses circular modules (for large QRs)
func (qb *QRBuilder) CircleShape() *QRBuilder {
	if qb.pixelWidth < 256 {
		log.Printf("Warning: circle_shape can lead to unreadable QR Code (requires pixel_width > 256)")
	}
	qb.circleShape = true
	return qb
}

// End finishes the QR command
func (qb *QRBuilder) End() *DocumentBuilder {
	cmd := qrCommand{
		Data:        qb.data,
		HumanText:   qb.humanText,
		PixelWidth:  qb.pixelWidth,
		Correction:  qb.correction,
		Align:       qb.align,
		Logo:        qb.logo,
		CircleShape: qb.circleShape,
	}
	return qb.parent.addCommand("qr", cmd)
}
