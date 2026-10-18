package bitimage

import (
	"errors"
)

// ============================================================================
// Context
// ============================================================================
// This sub-module implements ESC/POS commands for download graphics functionality.
// ESC/POS is the command system used by thermal receipt printers to control
// graphics data storage in volatile memory (RAM), including definition, deletion,
// printing, and capacity management of temporary graphics like logos and images.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Type definitions

// DLFunctionCode represents the function code for download graphics commands
type DLFunctionCode byte

const (
	// DLFuncGetRemaining represents function code for getting remaining capacity
	DLFuncGetRemaining DLFunctionCode = 4
	// DLFuncGetRemainingASCII represents function code for getting remaining capacity (ASCII)
	DLFuncGetRemainingASCII DLFunctionCode = 52
)

// DLGraphicsColorData represents color data for download graphics
type DLGraphicsColorData struct {
	Color GraphicsColor // Color value (49-52)
	Data  []byte        // Graphics data
}

// Constants for download graphics limits
const (
	// MaxDLGraphicsWidth represents maximum download graphics width in dots
	MaxDLGraphicsWidth = 8192
	// MaxDLGraphicsHeight represents maximum download graphics height in dots
	MaxDLGraphicsHeight = 2304
	// MaxDLColorGroups represents maximum number of color groups for multiple tone
	MaxDLColorGroups = 4
)

// ============================================================================
// Error Definitions
// ============================================================================

// Error variables
var (
	// ErrInvalidDLFunctionCode indicates an invalid download function code
	ErrInvalidDLFunctionCode = errors.New("invalid download function code")
	// ErrInvalidNumColors indicates invalid number of colors
	ErrInvalidNumColors = errors.New("invalid number of colors for specified tone")
	// ErrInvalidScaleFactor indicates invalid scale factor
	ErrInvalidScaleFactor = errors.New("invalid scale factor (try 1 or 2)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ DownloadGraphicsCapability = (*DownloadGraphicsCommands)(nil)

// DownloadGraphicsCapability defines download graphics-related operations
type DownloadGraphicsCapability interface {
	GetDownloadGraphicsRemainingCapacity(fn DLFunctionCode) ([]byte, error)
	GetDownloadGraphicsKeyCodeList() []byte
	DeleteAllDownloadGraphics() []byte
	DeleteDownloadGraphicsByKeyCode(kc1, kc2 byte) ([]byte, error)
	DefineDownloadGraphics(tone GraphicsTone, kc1, kc2 byte, width, height uint16,
		colorData []DLGraphicsColorData) ([]byte, error)
	DefineDownloadGraphicsLarge(tone GraphicsTone, kc1, kc2 byte, width, height uint16,
		colorData []DLGraphicsColorData) ([]byte, error)
	DefineDownloadGraphicsColumn(kc1, kc2 byte, width, height uint16,
		colorData []DLGraphicsColorData) ([]byte, error)
	DefineDownloadGraphicsColumnLarge(kc1, kc2 byte, width, height uint16,
		colorData []DLGraphicsColorData) ([]byte, error)
	PrintDownloadGraphics(kc1, kc2 byte, horizontalScale, verticalScale GraphicsScale) ([]byte, error)
	DefineBMPDownloadGraphics(kc1, kc2 byte, tone GraphicsTone, bmpData []byte) ([]byte, error)
}

// ============================================================================
// Main Implementation
// ============================================================================

// DownloadGraphicsCommands implements DownloadGraphicsCapability
type DownloadGraphicsCommands struct{}

// NewDownloadGraphicsCommands creates a new instance of DownloadGraphicsCommands
func NewDownloadGraphicsCommands() *DownloadGraphicsCommands {
	return &DownloadGraphicsCommands{}
}

// ============================================================================
// Helper Functions
// ============================================================================

// calculateDLRasterDataSize calculates the data size for download raster format graphics
func calculateDLRasterDataSize(width, height uint16) int {
	widthBytes := (int(width) + 7) / 8
	return widthBytes * int(height)
}

// calculateDLColumnDataSize calculates the data size for download column format graphics
func calculateDLColumnDataSize(width, height uint16) int {
	heightBytes := (int(height) + 7) / 8
	return int(width) * heightBytes
}

// ============================================================================
// Validation Helper Functions
// ============================================================================

// ValidateDLRemainingFunctionCode validates if function code is valid for remaining capacity query
func ValidateDLRemainingFunctionCode(fn DLFunctionCode) error {
	if fn != DLFuncGetRemaining && fn != DLFuncGetRemainingASCII {
		return ErrInvalidDLFunctionCode
	}
	return nil
}

// ValidateDLGraphicsDimensions validates download graphics dimensions
func ValidateDLGraphicsDimensions(width, height uint16) error {
	if width < 1 || width > MaxDLGraphicsWidth {
		return ErrInvalidWidth
	}
	if height < 1 || height > MaxDLGraphicsHeight {
		return ErrInvalidHeight
	}
	return nil
}

// ValidateDLColorDataForTone validates color data based on tone for download graphics
func ValidateDLColorDataForTone(tone GraphicsTone, colorData []DLGraphicsColorData) error {
	if len(colorData) == 0 {
		return ErrInvalidColorCount
	}

	// Check for duplicates
	colorMap := make(map[GraphicsColor]bool)
	for _, cd := range colorData {
		if colorMap[cd.Color] {
			return ErrDuplicateColor
		}
		colorMap[cd.Color] = true
	}

	if tone == Monochrome {
		// Monochrome only allows 1 color data group with color 1
		if len(colorData) != 1 {
			return ErrInvalidColorCount
		}
		if colorData[0].Color != Color1 {
			return ErrInvalidColor
		}
	} else { // MultipleTone
		// Multiple tone allows 1-4 color data groups with colors 1-4
		if len(colorData) > MaxDLColorGroups {
			return ErrInvalidColorCount
		}
		for _, cd := range colorData {
			if cd.Color < Color1 || cd.Color > Color4 {
				return ErrInvalidColor
			}
		}
	}

	return nil
}

// ValidateDLColumnColorData validates color data for column format download graphics
func ValidateDLColumnColorData(colorData []DLGraphicsColorData) error {
	if len(colorData) == 0 {
		return ErrInvalidColorCount
	}

	// Check for duplicates
	colorMap := make(map[GraphicsColor]bool)
	for _, cd := range colorData {
		if colorMap[cd.Color] {
			return ErrDuplicateColor
		}
		colorMap[cd.Color] = true
	}

	// Column format restrictions
	hasColor3 := false
	for _, cd := range colorData {
		if cd.Color == Color3 {
			hasColor3 = true
		} else if cd.Color != Color1 && cd.Color != Color2 {
			return ErrInvalidColor
		}
	}

	// If color 3 is used, it must be the only color
	if hasColor3 && len(colorData) != 1 {
		return ErrInvalidColorCount
	}

	// Colors 1 and 2 can be used alone or together (max 2)
	if !hasColor3 && len(colorData) > 2 {
		return ErrInvalidColorCount
	}

	return nil
}
