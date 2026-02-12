package constants

import (
	"fmt"
)

// ============================================================================
// Paper Constants
// ============================================================================

// Paper width constants in pixels at 203 DPI
const (
	// PaperPxWidth58mm = 48mm printable area * 8 dots/mm = 384 pixels
	PaperPxWidth58mm = 384
	// PaperPxWidth80mm = 72mm printable area * 8 dots/mm = 576 pixels
	PaperPxWidth80mm = 576
	// Paper58mm is 58mm paper width
	Paper58mm = 58
	// Paper80mm is 80mm paper width
	Paper80mm = 80
)

// ============================================================================
// Alignment Constants
// ============================================================================

// Ensure Alignment implements fmt.Stringer
var _ fmt.Stringer = Alignment("")

// Alignment options for content positioning
type Alignment string

func (a Alignment) String() string {
	return string(a)
}

const (
	// Left aligns content to the left
	Left Alignment = "left"
	// Center centers the content
	Center Alignment = "center"
	// Right aligns content to the right
	Right Alignment = "right"
)

// ============================================================================
// Dithering Constants
// ============================================================================

// Ensure Dithering implements fmt.Stringer
var _ fmt.Stringer = Dithering("")

// Dithering algorithms for image processing
type Dithering string

func (d Dithering) String() string {
	return string(d)
}

const (
	// Threshold uses simple thresholding
	Threshold Dithering = "threshold"
	// Atkinson uses Atkinson dithering algorithm
	Atkinson Dithering = "atkinson"
	// FloydSteinberg uses Floyd-Steinberg dithering algorithm
	FloydSteinberg Dithering = "floyd_steinberg"
	// TODO: Add more dithering algorithms if needed
)

// ============================================================================
// Scaling Constants
// ============================================================================

// Ensure Scaling implements fmt.Stringer
var _ fmt.Stringer = Scaling("")

// Scaling algorithms for image resizing
type Scaling string

func (s Scaling) String() string {
	return string(s)
}

const (
	// NearestNeighbor uses Nearest Neighbor Scaling
	NearestNeighbor Scaling = "nearest_neighbor"
	// Bilinear uses Bilinear Scaling
	Bilinear Scaling = "bilinear"
	// Bicubic uses Bicubic Scaling
	Bicubic Scaling = "bicubic"
	// TODO: Add more scaling algorithms if needed
)

// ============================================================================
// Size Constants
// ============================================================================

// Ensure Size implements fmt.Stringer
var _ fmt.Stringer = Size("")

// Size options for text sizing
type Size string

func (s Size) String() string {
	return string(s)
}

const (
	// Normal is normal text size
	Normal Size = "1x1"
	// Double is double text size
	Double Size = "2x2"
	// Triple is triple text size
	Triple Size = "3x3"
	// Quad is quadruple text size
	Quad Size = "4x4"
	// Penta is quintuple text size
	Penta Size = "5x5"
	// Hexa is hexa text size
	Hexa Size = "6x6"
	// Hepta is hepta text size
	Hepta Size = "7x7"
	// Octa is octa text size
	Octa Size = "8x8"
)

const (
	// MinScale is the minimum scale multiplier for text
	MinScale = 1.0
	// MaxScale is the maximum scale multiplier for text (ESC/POS limit)
	MaxScale = 8.0
)

// ============================================================================
// Underline Constants
// ============================================================================

// Ensure Underline implements fmt.Stringer
var _ fmt.Stringer = Underline("")

// Underline options for text underlining
type Underline string

func (u Underline) String() string {
	return string(u)
}

const (
	// NoDot indicates no underline
	NoDot Underline = "0pt"
	// OneDot indicates single underline
	OneDot Underline = "1pt"
	// TwoDot indicates double underline
	TwoDot Underline = "2pt"
)

const (
	// MinUnderline is the minimum underline mode
	MinUnderline = 0
	// MaxUnderline is the maximum underline mode (0=off, 1=single, 2=double)
	MaxUnderline = 2
)

// ============================================================================
// Font Constants
// ============================================================================

// Ensure Font implements fmt.Stringer
var _ fmt.Stringer = Font("")

// Font options for text font selection
type Font string

func (t Font) String() string {
	return string(t)
}

const (
	// A is 12x24 dots font
	A Font = "a"
	// B is 9x17 dots font
	B Font = "b"
	// C is Font C
	C Font = "c"
	// D is Font D
	D Font = "d"
)

// Font metrics constants (based on typical ESC/POS thermal printer fonts)
const (
	// FontAWidth is the character width in pixels for Font A (12x24)
	FontAWidth = 12
	// FontAHeight is the character height in pixels for Font A
	FontAHeight = 24

	// FontBWidth is the character width in pixels for Font B (9x17)
	FontBWidth = 9
	// FontBHeight is the character height in pixels for Font B
	FontBHeight = 17
)

// ============================================================================
// Cut Constants
// ============================================================================

// Ensure Mode implements fmt.Stringer
var _ fmt.Stringer = Mode("")

// Mode options for paper cutting
type Mode string

func (c Mode) String() string {
	return string(c)
}

const (
	// Partial indicates a partial cut
	Partial Mode = "partial"
	// Full indicates a full cut
	Full Mode = "full"
)

const (
	// CutDashLength is the length of dashes in partial cut lines
	CutDashLength = 8
	// CutGapLength is the gap between dashes in partial cut lines
	CutGapLength = 4
	// CutTeethSpacing is the spacing between "teeth" marks on partial cuts
	CutTeethSpacing = 20
	// CutSpaceAfter is the vertical space added after a cut line
	CutSpaceAfter = 10
	// FullCutThickness is the line thickness for full cuts
	FullCutThickness = 2
	// PartialCutThickness is the line thickness for partial cuts
	PartialCutThickness = 1
)

// ============================================================================
// Image Pixel Width Constants
// ============================================================================

const (
	// MaxImagePixelWidth maximum in pixels (0 - 576 @ 203 DPI)
	MaxImagePixelWidth = PaperPxWidth80mm
	// MinImagePixelWidth minimum in pixels (0 - 576 @ 203 DPI)
	MinImagePixelWidth = 32
)

// ============================================================================
// QR Code Pixel Width Constants
// ============================================================================

const (
	// MinQrPixelWidth minimum in pixels (4 + 4 silence zone + 21 modules grid @ 3 dots/module)
	MinQrPixelWidth = 87
	// MaxQrPixelWidth maximum in pixels (80mm @ 203 DPI)
	MaxQrPixelWidth = PaperPxWidth80mm
)

// ============================================================================
// Error Correction Constants
// ============================================================================

// Ensure ErrorCorrection implements fmt.Stringer
var _ fmt.Stringer = ErrorCorrection("")

// ErrorCorrection levels for QR codes
type ErrorCorrection string

func (e ErrorCorrection) String() string {
	return string(e)
}

const (
	// L is Low error correction level (7%)
	L ErrorCorrection = "L"
	// M is Medium error correction level (15%)
	M ErrorCorrection = "M"
	// Q is Quartile error correction level (25%)
	Q ErrorCorrection = "Q"
	// H is the Highest error correction level (30%)
	H ErrorCorrection = "H"
)

// ============================================================================
// Image Threshold Constants
// ============================================================================

const (
	// MaxImageThreshold maximum
	MaxImageThreshold = 255
	// MinImageThreshold minimum
	MinImageThreshold = 0
)

// ============================================================================
// QR Code Threshold Constants
// ============================================================================

const (
	// MaxQrThreshold maximum
	MaxQrThreshold = 255
	// MinQrThreshold minimum
	MinQrThreshold = 0
)

// ============================================================================
// Barcode Constants
// ============================================================================

// Ensure Symbology implements fmt.Stringer
var _ fmt.Stringer = Symbology("")

// Symbology options for barcode types
type Symbology string

func (b Symbology) String() string {
	return string(b)
}

const (
	// CODE39 symbology
	CODE39 Symbology = "code39"
	// CODE128 symbology
	CODE128 Symbology = "code128"
	// EAN13 symbology
	EAN13 Symbology = "ean13"
	// UPCA symbology
	UPCA Symbology = "upca"
	// ITF symbology
	ITF Symbology = "itf"
	// CODABAR symbology
	CODABAR Symbology = "codabar"
	// EAN8 symbology
	EAN8 Symbology = "ean8"
	// UPCE symbology
	UPCE Symbology = "upce"
	// JAN13 symbology (alias for EAN13)
	JAN13 Symbology = "jan13"
	// JAN8 symbology (alias for EAN8)
	JAN8 Symbology = "jan8"
)

// ============================================================================
// Human Readable Interpretation (HRI) Position Constants
// ============================================================================

// Ensure HRIPosition implements fmt.Stringer
var _ fmt.Stringer = HriPosition("")

// HriPosition options for Human Readable Interpretation placement
type HriPosition string

func (b HriPosition) String() string {
	return string(b)
}

const (
	// Below indicates HRI is printed below the barcode
	Below HriPosition = "below"
	// Above indicates HRI is printed above the barcode
	Above HriPosition = "above"
	// Both indicates HRI is printed both above and below the barcode
	Both HriPosition = "both"
	// None indicates no HRI is printed
	None HriPosition = "none"
)

// ============================================================================
// Width and Height Barcode Constants
// ============================================================================

const (
	// MinBarcodeModuleWidth is the minimum dot per module
	MinBarcodeModuleWidth = 2
	// MaxBarcodeModuleWidth is the maximum dot per module
	MaxBarcodeModuleWidth = 6
	// MinBarcodeHeight is the minimum height in dots
	MinBarcodeHeight = 1
	// MaxBarcodeHeight is the maximum height in dots
	MaxBarcodeHeight = 255
)

// ============================================================================
// Raw Command Constants
// ============================================================================

// Ensure Format implements fmt.Stringer
var _ fmt.Stringer = Format("")

// Format options for Human Readable Interpretation placement
type Format string

func (r Format) String() string {
	return string(r)
}

const (
	// Hex indicates raw command is in hexadecimal format
	Hex Format = "hex"
	// Base64 indicates raw command is in base64 format
	Base64 Format = "base64"
)

const (
	// RawMaxBytes maximum bytes for raw command
	RawMaxBytes = 4096
	// RawCommentMaxLen maximum length for raw command comments
	RawCommentMaxLen = 512
)

// ============================================================================
// Canvas Constants
// ============================================================================

const (
	// InitialHeight is the initial canvas height in pixels
	InitialHeight = 1000
	// MaxCanvasHeight is the maximum canvas height (safety limit ~1.25m of paper)
	MaxCanvasHeight = 10000
	// GrowthFactor is the multiplier when canvas needs to grow
	GrowthFactor = 1.5
	// BottomPadding is the padding added after the last printed element
	BottomPadding = 20
	// TopMarginMultiplier is used to calculate initial top margin
	TopMarginMultiplier = 2
)

// ============================================================================
// File Size Constants
// ============================================================================

const (
	// MaxJSONSize is the maximum size in bytes for a JSON document (10MB)
	MaxJSONSize = 10 * 1024 * 1024
)
