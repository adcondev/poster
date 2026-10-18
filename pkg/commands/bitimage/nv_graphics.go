package bitimage

import (
	"errors"
)

// ============================================================================
// Context
// ============================================================================
// This sub-module implements ESC/POS commands for NV (non-volatile) graphics functionality.
// ESC/POS is the command system used by thermal receipt printers to control
// graphics data storage in non-volatile memory, including definition, deletion,
// printing, and capacity management of persistent graphics like logos and images.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Type definitions

// NVFunctionCode represents the function code for NV graphics commands
type NVFunctionCode byte

const (
	// NVFuncGetCapacity represents function code for getting NV capacity
	NVFuncGetCapacity NVFunctionCode = 0
	// NVFuncGetCapacityASCII represents function code for getting NV capacity (ASCII)
	NVFuncGetCapacityASCII NVFunctionCode = 48
	// NVFuncGetRemaining represents function code for getting remaining capacity
	NVFuncGetRemaining NVFunctionCode = 3
	// NVFuncGetRemainingASCII represents function code for getting remaining capacity (ASCII)
	NVFuncGetRemainingASCII NVFunctionCode = 51
)

// NVGraphicsColorData represents color data for NV graphics
type NVGraphicsColorData struct {
	Color GraphicsColor // Color value (49-52)
	Data  []byte        // Graphics data
}

// Constants for NV graphics limits
const (
	// MaxNVGraphicsWidth represents maximum NV graphics width in dots
	MaxNVGraphicsWidth = 8192
	// MaxNVGraphicsHeight represents maximum NV graphics height in dots
	MaxNVGraphicsHeight = 2304
	// MaxNVColorGroups represents maximum number of color groups for multiple tone
	MaxNVColorGroups = 4
	// MinKeyCode represents minimum key code value (ASCII space)
	MinKeyCode = 32
	// MaxKeyCode represents maximum key code value (ASCII ~)
	MaxKeyCode = 126
	// NVGraphicsAreaSize represents total NV graphics area size in KB
)

// ============================================================================
// Error Definitions
// ============================================================================

// Error variables
var (
	// ErrInvalidNVFunctionCode indicates an invalid NV function code
	ErrInvalidNVFunctionCode = errors.New("invalid NV function code")
	// ErrInvalidKeyCode indicates an invalid key code
	ErrInvalidKeyCode = errors.New("invalid key code (try ASCII 32-126)")
	// ErrInvalidColorCount indicates invalid number of color data groups
	ErrInvalidColorCount = errors.New("invalid color count for specified tone")
	// ErrInvalidBMPFormat indicates invalid Windows BMP format
	ErrInvalidBMPFormat = errors.New("invalid Windows BMP format")
	// ErrDuplicateColor indicates duplicate color in color data
	ErrDuplicateColor = errors.New("duplicate color in color data")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ NVGraphicsCapability = (*NvGraphicsCommands)(nil)

// NVGraphicsCapability defines NV graphics-related operations
type NVGraphicsCapability interface {
	GetNVGraphicsCapacity(fn NVFunctionCode) ([]byte, error)
	GetNVGraphicsRemainingCapacity(fn NVFunctionCode) ([]byte, error)
	GetNVGraphicsKeyCodeList() []byte
	DeleteAllNVGraphics() []byte
	DeleteNVGraphicsByKeyCode(kc1, kc2 byte) ([]byte, error)
	DefineNVRasterGraphics(tone GraphicsTone, kc1, kc2 byte, width, height uint16, colorData []NVGraphicsColorData) ([]byte, error)
	DefineNVRasterGraphicsLarge(tone GraphicsTone, kc1, kc2 byte, width, height uint16, colorData []NVGraphicsColorData) ([]byte, error)
	DefineNVColumnGraphics(kc1, kc2 byte, width, height uint16, colorData []NVGraphicsColorData) ([]byte, error)
	DefineNVColumnGraphicsLarge(kc1, kc2 byte, width, height uint16, colorData []NVGraphicsColorData) ([]byte, error)
	PrintNVGraphics(kc1, kc2 byte, horizontalScale, verticalScale GraphicsScale) ([]byte, error)
	DefineWindowsBMPNVGraphics(kc1, kc2 byte, tone GraphicsTone, bmpData []byte) ([]byte, error)
}

// ============================================================================
// Main Implementation
// ============================================================================

// NvGraphicsCommands implements NVGraphicsCapability
type NvGraphicsCommands struct{}

// NewNVGraphicsCommands creates a new instance of NvGraphicsCommands
func NewNVGraphicsCommands() *NvGraphicsCommands {
	return &NvGraphicsCommands{}
}

// ============================================================================
// Helper Functions
// ============================================================================

// calculateNVRasterDataSize calculates the data size for NV raster format graphics
func calculateNVRasterDataSize(width, height uint16) int {
	widthBytes := (int(width) + 7) / 8
	return widthBytes * int(height)
}

// calculateNVColumnDataSize calculates the data size for NV column format graphics
func calculateNVColumnDataSize(width, height uint16) int {
	heightBytes := (int(height) + 7) / 8
	return int(width) * heightBytes
}

// ============================================================================
// Validation Helper Functions
// ============================================================================

// ValidateNVCapacityFunctionCode validates if function code is valid for capacity query
func ValidateNVCapacityFunctionCode(fn NVFunctionCode) error {
	if fn != NVFuncGetCapacity && fn != NVFuncGetCapacityASCII {
		return ErrInvalidNVFunctionCode
	}
	return nil
}

// ValidateNVRemainingFunctionCode validates if function code is valid for remaining capacity query
func ValidateNVRemainingFunctionCode(fn NVFunctionCode) error {
	if fn != NVFuncGetRemaining && fn != NVFuncGetRemainingASCII {
		return ErrInvalidNVFunctionCode
	}
	return nil
}

// ValidateKeyCode validates if key code is valid
func ValidateKeyCode(kc byte) error {
	if kc < MinKeyCode || kc > MaxKeyCode {
		return ErrInvalidKeyCode
	}
	return nil
}

// ValidateNVGraphicsDimensions validates NV graphics dimensions
func ValidateNVGraphicsDimensions(width, height uint16) error {
	if width < 1 || width > MaxNVGraphicsWidth {
		return ErrInvalidWidth
	}
	if height < 1 || height > MaxNVGraphicsHeight {
		return ErrInvalidHeight
	}
	return nil
}

// ValidateBMPData validates Windows BMP format data
func ValidateBMPData(data []byte) error {
	// Check minimum BMP file size (file header + DIB header)
	if len(data) < 54 {
		return ErrInvalidBMPFormat
	}

	// Check BMP signature "BM"
	if data[0] != 'B' || data[1] != 'M' {
		return ErrInvalidBMPFormat
	}

	return nil
}

// validateColorDataForTone validates color data based on tone
func validateColorDataForTone(tone GraphicsTone, colorData []NVGraphicsColorData) error {
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
		// Monochrome only allows 1 color data group with color 1 or 2
		if len(colorData) != 1 {
			return ErrInvalidColorCount
		}
		if colorData[0].Color != Color1 && colorData[0].Color != Color2 {
			return ErrInvalidColor
		}
	} else { // MultipleTone
		// Multiple tone allows 1-4 color data groups with colors 1-4
		if len(colorData) > MaxNVColorGroups {
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

// validateColumnColorData validates color data for column format
func validateColumnColorData(colorData []NVGraphicsColorData) error {
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
