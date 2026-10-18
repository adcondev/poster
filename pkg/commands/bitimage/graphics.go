package bitimage

import (
	"errors"
)

// ============================================================================
// Context
// ============================================================================
// This sub-module implements ESC/POS commands for graphics data handling.
// ESC/POS is the command system used by thermal receipt printers to control
// graphics data storage, printing, and configuration including dot density,
// color selection, and various graphics formats (raster and column).

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Type definitions

// DotDensity represents the dot density setting for graphics
type DotDensity byte

const (
	// Density180x180 represents 180 dpi × 180 dpi
	Density180x180 DotDensity = 50
	// Density360x360 represents 360 dpi × 360 dpi
	Density360x360 DotDensity = 51
)

// FunctionCode represents the function code for graphics commands
type FunctionCode byte

const (
	// FunctionCodeDensity1 represents function code 1 for density setting
	FunctionCodeDensity1 FunctionCode = 1
	// FunctionCodeDensity49 represents function code 49 for density setting
	FunctionCodeDensity49 FunctionCode = 49
	// FunctionCodePrint2 represents function code 2 for printing
	FunctionCodePrint2 FunctionCode = 2
	// FunctionCodePrint50 represents function code 50 for printing
	FunctionCodePrint50 FunctionCode = 50
)

// GraphicsTone represents the graphics tone mode
type GraphicsTone byte

const (
	// Monochrome represents monochrome (digital) graphics
	Monochrome GraphicsTone = 48
	// MultipleTone represents multiple tone graphics
	MultipleTone GraphicsTone = 52
)

// GraphicsScale represents the scaling factor for graphics
type GraphicsScale byte

const (
	// NormalScale represents normal scale (x1)
	NormalScale GraphicsScale = 1
	// DoubleScale represents double scale (x2)
	DoubleScale GraphicsScale = 2
)

// GraphicsColor represents the color selection for graphics
type GraphicsColor byte

const (
	// Color1 represents color 1
	Color1 GraphicsColor = 49
	// Color2 represents color 2
	Color2 GraphicsColor = 50
	// Color3 represents color 3
	Color3 GraphicsColor = 51
	// Color4 represents color 4
	Color4 GraphicsColor = 52
)

// Constants for graphics limits
const (
	// MaxGraphicsWidth represents maximum graphics width in dots
	MaxGraphicsWidth = 2400
	// MaxMonochromeHeightNormal represents maximum monochrome height with normal scale
	MaxMonochromeHeightNormal = 2400
	// MaxMonochromeHeightDouble represents maximum monochrome height with double scale
	MaxMonochromeHeightDouble = 1200
	// MaxMultiToneHeightNormal represents maximum multi-tone height with normal scale
	MaxMultiToneHeightNormal = 600
	// MaxMultiToneHeightDouble represents maximum multi-tone height with double scale
	MaxMultiToneHeightDouble = 300
	// MaxColumnWidth represents maximum column format width in dots
	MaxColumnWidth = 2048
	// MaxColumnHeight represents maximum column format height in dots
	MaxColumnHeight = 128
	// MaxStandardCommandSize represents maximum size for standard command format
	MaxStandardCommandSize = 65535
	// MaxExtendedCommandSize represents maximum size for extended command format
	MaxExtendedCommandSize uint32 = 4294967295
)

// ============================================================================
// Error Definitions
// ============================================================================

// Error variables
var (
	// ErrInvalidFunctionCode indicates an invalid function code
	ErrInvalidFunctionCode = errors.New("invalid function code")
	// ErrInvalidDensityValue indicates an invalid density value
	ErrInvalidDensityValue = errors.New("invalid density value (try 50 or 51)")
	// ErrInvalidDensityCombination indicates an invalid density combination
	ErrInvalidDensityCombination = errors.New("invalid density combination (x and y must match)")
	// ErrInvalidTone indicates an invalid graphics tone
	ErrInvalidTone = errors.New("invalid tone (try 48 for monochrome or 52 for multiple tone)")
	// ErrInvalidScale indicates an invalid scale factor
	ErrInvalidScale = errors.New("invalid scale (try 1 or 2)")
	// ErrInvalidColor indicates an invalid color selection
	ErrInvalidColor = errors.New("invalid color (try 49-52 for colors 1-4)")
	// ErrInvalidWidth indicates invalid graphics width
	ErrInvalidWidth = errors.New("invalid width (1-2400 for raster, 1-2048 for column)")
	// ErrInvalidHeight indicates invalid graphics height
	ErrInvalidHeight = errors.New("invalid height (check limits based on tone and scale)")
	// ErrDataTooLarge indicates data size exceeds command limits
	ErrDataTooLarge = errors.New("data size exceeds command limits")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ GraphicsCapability = (*GraphicsCommands)(nil)

// GraphicsCapability defines graphics-related operations
type GraphicsCapability interface {
	SetGraphicsDotDensity(fn FunctionCode, x, y DotDensity) ([]byte, error)
	PrintBufferedGraphics(fn FunctionCode) ([]byte, error)
	StoreRasterGraphicsInBuffer(tone GraphicsTone, horizontalScale, verticalScale GraphicsScale,
		color GraphicsColor, width, height uint16, data []byte) ([]byte, error)
	StoreRasterGraphicsInBufferLarge(tone GraphicsTone, horizontalScale, verticalScale GraphicsScale,
		color GraphicsColor, width, height uint16, data []byte) ([]byte, error)
	StoreColumnGraphicsInBuffer(horizontalScale, verticalScale GraphicsScale,
		color GraphicsColor, width, height uint16, data []byte) ([]byte, error)
	StoreColumnGraphicsInBufferLarge(horizontalScale, verticalScale GraphicsScale,
		color GraphicsColor, width, height uint16, data []byte) ([]byte, error)
}

// ============================================================================
// Main Implementation
// ============================================================================

// GraphicsCommands implements GraphicsCapability
type GraphicsCommands struct{}

// NewGraphicsCommands creates a new instance of GraphicsCommands
func NewGraphicsCommands() *GraphicsCommands {
	return &GraphicsCommands{}
}

// ============================================================================
// Helper Functions
// ============================================================================

// calculateRasterDataSize calculates the data size for raster format graphics
func calculateRasterDataSize(width, height uint16) int {
	widthBytes := (int(width) + 7) / 8
	return widthBytes * int(height)
}

// calculateColumnDataSize calculates the data size for column format graphics
func calculateColumnDataSize(width, height uint16) int {
	heightBytes := (int(height) + 7) / 8
	return int(width) * heightBytes
}

// getMaxHeight returns the maximum height based on tone and vertical scale
func getMaxHeight(tone GraphicsTone, verticalScale GraphicsScale) uint16 {
	if tone == Monochrome {
		if verticalScale == NormalScale {
			return MaxMonochromeHeightNormal
		}
		return MaxMonochromeHeightDouble
	}
	// Multiple tone
	if verticalScale == NormalScale {
		return MaxMultiToneHeightNormal
	}
	return MaxMultiToneHeightDouble
}

// ============================================================================
// Validation Helper Functions
// ============================================================================

// ValidateDensityFunctionCode validates if function code is valid for density setting
func ValidateDensityFunctionCode(fn FunctionCode) error {
	if fn != FunctionCodeDensity1 && fn != FunctionCodeDensity49 {
		return ErrInvalidFunctionCode
	}
	return nil
}

// ValidatePrintFunctionCode validates if function code is valid for printing
func ValidatePrintFunctionCode(fn FunctionCode) error {
	if fn != FunctionCodePrint2 && fn != FunctionCodePrint50 {
		return ErrInvalidFunctionCode
	}
	return nil
}

// ValidateDotDensity validates if dot density values are valid
func ValidateDotDensity(x, y DotDensity) error {
	if x != Density180x180 && x != Density360x360 {
		return ErrInvalidDensityValue
	}
	if x != y {
		return ErrInvalidDensityCombination
	}
	return nil
}

// ValidateGraphicsTone validates if graphics tone is valid
func ValidateGraphicsTone(tone GraphicsTone) error {
	if tone != Monochrome && tone != MultipleTone {
		return ErrInvalidTone
	}
	return nil
}

// ValidateGraphicsScale validates if scale factor is valid
func ValidateGraphicsScale(scale GraphicsScale) error {
	if scale != NormalScale && scale != DoubleScale {
		return ErrInvalidScale
	}
	return nil
}

// ValidateGraphicsColor validates if color selection is valid
func ValidateGraphicsColor(color GraphicsColor) error {
	if color < Color1 || color > Color4 {
		return ErrInvalidColor
	}
	return nil
}

// ValidateRasterDimensions validates raster format graphics dimensions
func ValidateRasterDimensions(width, height uint16, tone GraphicsTone, verticalScale GraphicsScale) error {
	if width < 1 || width > MaxGraphicsWidth {
		return ErrInvalidWidth
	}

	maxHeight := getMaxHeight(tone, verticalScale)
	if height < 1 || height > maxHeight {
		return ErrInvalidHeight
	}

	return nil
}

// ValidateColumnDimensions validates column format graphics dimensions
func ValidateColumnDimensions(width, height uint16) error {
	if width < 1 || width > MaxColumnWidth {
		return ErrInvalidWidth
	}
	if height < 1 || height > MaxColumnHeight {
		return ErrInvalidHeight
	}
	return nil
}
