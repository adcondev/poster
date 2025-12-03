// Package constants provides centralized default values for the entire library.
// All packages should reference these constants to ensure consistency.
package constants

import (
	"fmt"
)

// TODO: Add MustCompile validation for some constants (Text Size, etc.)

// ============================================================================
// Defaults
// ============================================================================

// Default profile values
const (
	// DefaultPaperWidth default in millimeters
	DefaultPaperWidth = 80
	// DefaultCodeTable default
	DefaultCodeTable = "WPC1252"
	// DefaultDPI default
	DefaultDPI = 203
)

// Valid options for profile configuration
var (
	// ValidPaperWidths in millimeters
	ValidPaperWidths = []int{58, 72, 80, 100, 112, 120}

	// ValidDPIs for printers
	ValidDPIs = []int{203, 300, 600}

	// ValidCodeTables for character encoding
	ValidCodeTables = []string{"WPC1252", "PC850", "PC437", "PC858"}
)

// Default alignment constants for various elements
const (
	// DefaultAlignTable is the default alignment for tables
	DefaultAlignTable = Center
	// DefaultAlignTableColumn is the default alignment for table columns
	DefaultAlignTableColumn = Center
)

// Text defaults
const (
	// DefaultTextAlignment is the default text alignment
	DefaultTextAlignment = Left
	// DefaultTextFont is the default font
	DefaultTextFont = A
	// DefaultTextSize is the default text size
	DefaultTextSize = Normal
	// DefaultTextNewLine indicates if a new line is added after text by default
	DefaultTextNewLine = true
	// DefaultUnderline is the default underline style
	DefaultUnderline = NoDot
)

// Image defaults for processing
const (
	// DefaultImagePixelWidth default in pixels (0 - 576 @ 203 DPI)
	DefaultImagePixelWidth = 128
	// DefaultImageThreshold default (0-255)
	DefaultImageThreshold = 128
	// DefaultImageDithering is the default dithering algorithm
	DefaultImageDithering = Atkinson
	// DefaultImageScaling is the default scaling algorithm
	DefaultImageScaling = Bilinear
	// DefaultImageAlignment is the default alignment for images
	DefaultImageAlignment = Center
)

// QR Code defaults
const (
	// DefaultQrThreshold is the default threshold for QR codes
	DefaultQrThreshold = 128
	// DefaultQrPixelWidth is the default pixel width for QR codes (87 - 576 @ 203 DPI)
	DefaultQrPixelWidth = 256
	// DefaultQrDithering is the default dithering algorithm for QR codes
	DefaultQrDithering = Threshold
	// DefaultQrScaling is the default scaling algorithm for QR codes
	DefaultQrScaling = NearestNeighbor
	// DefaultQrAlignment is the default alignment for QR codes
	DefaultQrAlignment = Center
	// DefaultQrErrorCorrection is the default error correction level for QR codes
	DefaultQrErrorCorrection = Q
)

// Barcode defaults
const (
	// DefaultBarcodeDithering is the default dithering algorithm for barcodes
	DefaultBarcodeDithering = Threshold
	// DefaultBarcodeModuleWidth is the default pixel width for barcodes (2 - 6 dots per module)
	DefaultBarcodeModuleWidth = 3
	// DefaultBarcodeHeight is the default pixel height for barcodes
	DefaultBarcodeHeight = 64
	// DefaultBarcodeHriFont is the default font for Human Readable Interpretation
	DefaultBarcodeHriFont = A
	// DefaultBarcodeAlignment default
	DefaultBarcodeAlignment = Center
	// DefaultBarcodeHriPosition default
	DefaultBarcodeHriPosition = Below
)

// Separator defaults
const (
	// DefaultSeparatorChar default (trimmed to DefaultSeparatorLength)
	DefaultSeparatorChar = "- "
	// DefaultSeparatorLength default (58mm paper @ Font A with 203 DPI)
	DefaultSeparatorLength = 32
	// DefaultSeparatorAlignment is the default alignment for separators
	DefaultSeparatorAlignment = Center
)

// Feed defaults
const (
	// DefaultFeedLines default
	DefaultFeedLines = 2
)

// Cut defaults
const (
	// DefaultCutMode default
	DefaultCutMode = Partial
	// DefaultCutFeed default
	DefaultCutFeed = 2
)

// Table Defaults
const (
	// DefaultTableShowHeaders indicates if table headers are shown by default
	DefaultTableShowHeaders = true
	// DefaultTableHeaderBold indicates if table headers are bold by default
	DefaultTableHeaderBold = true
	// DefaultTableWordWrap indicates if table cells should wrap text by default
	DefaultTableWordWrap = true
	// DefaultTableColumnSpacing default spacing between columns
	DefaultTableColumnSpacing = 1
	// DefaultTableAlignment default table alignment
	DefaultTableAlignment = Center
	// DefaultTableColumnAlignment default font for table columns
	DefaultTableColumnAlignment = Center
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

// ============================================================================
// Image Pixel Width Constants
// ============================================================================

const (
	// MaxImagePixelWidth maximum in pixels (0 - 576 @ 203 DPI)
	MaxImagePixelWidth = 576
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
	MaxQrPixelWidth = 576
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
// Raw Command Limits
// ============================================================================

const (
	// RawMaxBytes maximum bytes for raw command
	RawMaxBytes = 4096
	// RawCommentMaxLen maximum length for raw command comments
	RawCommentMaxLen = 512
)
