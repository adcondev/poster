// Package constants provides centralized default values for the entire library.
// All packages should reference these constants to ensure consistency.
package constants

// ============================================================================
// Profile Defaults (from schema)
// ============================================================================

const (
	// PaperWidth default in millimeters
	PaperWidth = 80
	// CodeTable default
	CodeTable = "WPC1252"
	// DPI default
	DPI = 203
)

// ValidPaperWidths in millimeters
var ValidPaperWidths = []int{58, 72, 80, 100, 112, 120}

// ValidDPIs for printers
var ValidDPIs = []int{203, 300, 600}

// ============================================================================
// Image Defaults
// ============================================================================

const (
	// ImagePixelWidth default in pixels
	ImagePixelWidth = 128
	// ImageThreshold default (0-255)
	ImageThreshold = 128
	// ImageAlign default
	ImageAlign = "center"
)

// ============================================================================
// Dithering
// ============================================================================

// Dithering algorithms for image processing
type Dithering string

const (
	// DitheringThreshold uses simple thresholding
	DitheringThreshold Dithering = "threshold"
	// DitheringAtkinson uses Atkinson dithering algorithm
	DitheringAtkinson Dithering = "atkinson"
	// DitheringDefault is the default dithering algorithm
	DitheringDefault = DitheringAtkinson
)

func (d Dithering) String() string {
	return string(d)
}

// ============================================================================
// Scaling
// ============================================================================

// Scaling algorithms for image resizing
type Scaling string

func (s Scaling) String() string {
	return string(s)
}

const (
	// ScalingNNS uses Nearest Neighbor Scaling
	ScalingNNS Scaling = "nns"
	// ScalingBilinear uses Bilinear Scaling
	ScalingBilinear Scaling = "bilinear"
	// ScalingDefault is the default scaling algorithm
	ScalingDefault = ScalingBilinear
)

// ============================================================================
// Alignment
// ============================================================================

// Alignment options for content positioning
type Alignment string

const (
	// AlignLeft aligns content to the left
	AlignLeft Alignment = "left"
	// AlignCenter centers the content
	AlignCenter Alignment = "center"
	// AlignRight aligns content to the right
	AlignRight Alignment = "right"
	// AlignDefault is the default alignment
	AlignDefault = AlignLeft
)

func (a Alignment) String() string {
	return string(a)
}

// ============================================================================
// Text Defaults
// ============================================================================

const (
	// TextAlign default
	TextAlign = "left"
	// TextSize default
	TextSize = "1x1"
	// TextFont default
	TextFont = "A"
	// TextNewLine default
	TextNewLine = true
)

// ============================================================================
// Separator Defaults
// ============================================================================

const (
	// SeparatorChar default
	SeparatorChar = "- "
	// SeparatorLength default (80mms Font A)
	SeparatorLength = 48
)

// ============================================================================
// Cut Defaults
// ============================================================================

const (
	// CutMode default
	CutMode = "partial"
	// CutFeed default in lines
	CutFeed = 2
)

// ============================================================================
// QR Defaults
// ============================================================================

const (
	// QRPixelWidth default in pixels
	QRPixelWidth = 128
	// QRMinPixelWidth minimum in pixels
	QRMinPixelWidth = 87
	// QRMaxPixelWidth maximum in pixels
	QRMaxPixelWidth = 576 // 80mm @ 203 DPI
	// QRCorrection default error correction level
	QRCorrection = "Q"
	// QRAlign default
	QRAlign = "center"
)

// ============================================================================
// Barcode Defaults
// ============================================================================

const (
	// BarcodeHRIPosition default
	BarcodeHRIPosition = "below"
	// BarcodeHRIFont default
	BarcodeHRIFont = "A"
	// BarcodeAlign default
	BarcodeAlign = "center"
	// BarcodeMinWidth limit
	BarcodeMinWidth = 2
	// BarcodeMaxWidth limit
	BarcodeMaxWidth = 6
	// BarcodeMinHeight limit
	BarcodeMinHeight = 1
	// BarcodeMaxHeight limit
	BarcodeMaxHeight = 255
)

// ============================================================================
// Table Defaults
// ============================================================================

const (
	// TableShowHeaders indicates if table headers are shown by default
	TableShowHeaders = true
	// TableHeaderBold indicates if table headers are bold by default
	TableHeaderBold = true
	// TableWordWrap indicates if table cells should wrap text by default
	TableWordWrap = true
	// TableColumnSpacing default spacing between columns
	TableColumnSpacing = 1
	// TableAlign default table alignment
	TableAlign = "center"
	// TableColumnAlign default font for table columns
	TableColumnAlign = "center"
)

// Paper width in characters
const (
	// TableWidth58mm paper width in characters
	TableWidth58mm = 32
	// TableWidth80mm paper width in characters
	TableWidth80mm = 48
)

// ============================================================================
// Raw Command Limits
// ============================================================================

const (
	// RawMaxBytes maximum bytes for raw command
	RawMaxBytes = 4096
	// RawCommentMaxLen maximum length for raw command comments
	RawCommentMaxLen = 500
)
