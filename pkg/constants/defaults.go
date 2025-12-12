// Package constants provides centralized default values for the entire library.
// All packages should reference these constants to ensure consistency.
package constants

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
	// DefaultDotsPerMm is the number of dots per millimeter at 203 DPI
	DefaultDotsPerMm = 8
	// DefaultVersion default
	DefaultVersion = "1.0"
	// DefaultLineSpacing is the default line spacing in pixels
	DefaultLineSpacing = 30
)

// Valid options for profile configuration
var (
	// ValidPaperWidths in millimeters
	ValidPaperWidths = []int{58, 72, 80, 100, 112, 120}

	// ValidDPIs for printers
	ValidDPIs = []int{DefaultDPI, 300, 600}

	// ValidCodeTables for character encoding
	ValidCodeTables = []string{"WPC1252", "PC850", "PC437", "PC858"}
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
	// DefaultSeparatorChar default (repeated to fill DefaultSeparatorLength)
	DefaultSeparatorChar = "-"
	// DefaultSeparatorLength default (58mm paper @ Font A with 203 DPI)
	DefaultSeparatorLength = 32
	// DefaultSeparatorAlignment is the default alignment for separators
	DefaultSeparatorAlignment = Center
)

// Feed defaults
const (
	// DefaultFeedLines default
	DefaultFeedLines = 1
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

// Pulse (Cash Drawer) defaults
const (
	DefaultPulsePin     = 0
	DefaultPulseOnTime  = 50  // ms
	DefaultPulseOffTime = 100 // ms
)

// Beep defaults
const (
	DefaultBeepTimes = 1
	DefaultBeepLapse = 1
)

// Raw defaults
const (
	// DefaultRawFormat is the default format for raw data
	DefaultRawFormat = Hex
)
