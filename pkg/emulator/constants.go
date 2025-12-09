package emulator

// Paper width constants in pixels at 203 DPI
const (
	// PaperWidth58mm = 48mm printable area * 8 dots/mm = 384 pixels
	PaperWidth58mm = 384
	// PaperWidth80mm = 72mm printable area * 8 dots/mm = 576 pixels
	PaperWidth80mm = 576
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

	// DefaultLineSpacing is the default line spacing in pixels
	DefaultLineSpacing = 30
)

// Canvas constants
const (
	// InitialHeight is the initial canvas height in pixels
	InitialHeight = 1000
	// MaxCanvasHeight is the maximum canvas height (safety limit ~2. 5m of paper)
	MaxCanvasHeight = 20000
	// GrowthFactor is the multiplier when canvas needs to grow
	GrowthFactor = 1.5
	// BottomPadding is the padding added after the last printed element
	BottomPadding = 20
)

// Alignment constants
const (
	AlignLeft   = "left"
	AlignCenter = "center"
	AlignRight  = "right"
)
