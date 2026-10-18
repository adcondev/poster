package bitimage

import (
	"errors"
)

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for bit image and graphics functionality.
// ESC/POS is the command system used by thermal receipt printers to control
// bit image printing in various formats, manage non-volatile and downloaded images,
// and handle graphics data for logos and pictures.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Type and Constant definition

// Mode represents the bit image mode for SelectBitImageMode command
type Mode byte

const (
	// SingleDensity8 represents 8-dot single-density mode
	SingleDensity8 Mode = 0
	// DoubleDensity8 represents 8-dot double-density mode
	DoubleDensity8 Mode = 1
	// SingleDensity24 represents 24-dot single-density mode
	SingleDensity24 Mode = 32
	// DoubleDensity24 represents 24-dot double-density mode
	DoubleDensity24 Mode = 33
)

// PrintMode represents the print mode for NV/downloaded bit images
type PrintMode byte

const (
	// Normal represents normal print mode (x1 horizontal, x1 vertical)
	Normal PrintMode = 0
	// DoubleWidth represents double-width print mode (x2 horizontal, x1 vertical)
	DoubleWidth PrintMode = 1
	// DoubleHeight represents double-height print mode (x1 horizontal, x2 vertical)
	DoubleHeight PrintMode = 2
	// Quadruple represents quadruple print mode (x2 horizontal, x2 vertical)
	Quadruple PrintMode = 3
	// NormalASCII represents normal print mode in ASCII format
	NormalASCII PrintMode = 48
	// DoubleWidthASCII represents double-width print mode in ASCII format
	DoubleWidthASCII PrintMode = 49
	// DoubleHeightASCII represents double-height print mode in ASCII format
	DoubleHeightASCII PrintMode = 50
	// QuadrupleASCII represents quadruple print mode in ASCII format
	QuadrupleASCII PrintMode = 51
)

// NVBitImageData represents data for a single NV bit image
type NVBitImageData struct {
	Width  uint16 // Horizontal size in bytes (1-1023)
	Height uint16 // Vertical size in bytes (1-288)
	Data   []byte // Bit image data in column format
}

// Constants for image size limits
const (
	// MaxHorizontalDots represents maximum horizontal dots for bit images
	MaxHorizontalDots = 2400
	// MaxVerticalBytes8Dot represents maximum vertical bytes for 8-dot modes
	MaxVerticalBytes8Dot = 1
	// MaxVerticalBytes24Dot represents maximum vertical bytes for 24-dot modes
	MaxVerticalBytes24Dot = 3
	// MaxNVImageWidth represents maximum NV bit image width in bytes
	MaxNVImageWidth = 1023
	// MaxNVImageHeight represents maximum NV bit image height in bytes
	MaxNVImageHeight = 288
	// MaxDownloadedWidth represents maximum downloaded bit image width in bytes
	MaxDownloadedWidth = 255
	// MaxDownloadedHeight represents maximum downloaded bit image height in bytes
	MaxDownloadedHeight = 48
	// MaxDownloadedSize represents maximum downloaded bit image total size
	MaxDownloadedSize = 1536
	// MaxVariableWidth represents maximum variable vertical size bit image width
	MaxVariableWidth = 4256
	// MaxVariableHeight represents maximum variable vertical size bit image height
	MaxVariableHeight = 16
	// MaxRasterWidth represents maximum raster bit image width in bytes
	MaxRasterWidth = 65535
	// MaxRasterHeight represents maximum raster bit image height in dots
	MaxRasterHeight = 2303
)

// ============================================================================
// Error Definitions
// ============================================================================

// Error variables
var (
	// ErrBitImageMode indicates an invalid bit image mode
	ErrBitImageMode = errors.New("invalid bit image mode (try 0, 1, 32, or 33)")
	// ErrHorizontalDotsRange indicates horizontal dots out of range
	ErrHorizontalDotsRange = errors.New("invalid horizontal dots range (try 1-2400)")
	// ErrDataLength indicates data length mismatch
	ErrDataLength = errors.New("data length does not match expected size")
	// ErrInvalidMode indicates an invalid print mode
	ErrInvalidMode = errors.New("invalid print mode (try 0-3 or 48-51)")
	// ErrInvalidBitImageNumber indicates an invalid bit image number
	ErrInvalidBitImageNumber = errors.New("invalid bit image number (try 1-255)")
	// ErrInvalidImageCount indicates invalid number of images
	ErrInvalidImageCount = errors.New("invalid image count or mismatch with data")
	// ErrInvalidImageDimensions indicates invalid image dimensions
	ErrInvalidImageDimensions = errors.New("invalid image dimensions")
	// ErrInvalidHorizontalSize indicates invalid horizontal size
	ErrInvalidHorizontalSize = errors.New("invalid horizontal size (try 1-255)")
	// ErrInvalidVerticalSize indicates invalid vertical size
	ErrInvalidVerticalSize = errors.New("invalid vertical size (try 1-48)")
	// ErrInvalidDimensions indicates total dimensions exceed limit
	ErrInvalidDimensions = errors.New("invalid dimensions (x × y must be 1-1536)")
	// ErrInvalidDataLength indicates data length mismatch
	ErrInvalidDataLength = errors.New("data length does not match dimensions")
	// ErrInvalidImageWidth indicates invalid image width
	ErrInvalidImageWidth = errors.New("invalid image width")
	// ErrInvalidImageHeight indicates invalid image height
	ErrInvalidImageHeight = errors.New("invalid image height")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ Capability = (*Commands)(nil)

// Capability defines the interface for bit image commands
type Capability interface {
	SelectBitImageMode(mode Mode, width uint16, data []byte) ([]byte, error) // Legacy bit image commands
	PrintNVBitImage(n byte, mode PrintMode) ([]byte, error)                  // NV bit image commands (deprecated)
	DefineNVBitImage(n byte, images []NVBitImageData) ([]byte, error)
	DefineDownloadedBitImage(x, y byte, data []byte) ([]byte, error) // Downloaded bit image commands (deprecated)
	PrintDownloadedBitImage(mode PrintMode) ([]byte, error)
	PrintVariableVerticalSizeBitImage(mode PrintMode, width, height uint16, data []byte) ([]byte, error) // Variable size bit image commands (deprecated)
	PrintRasterBitImage(mode PrintMode, width, height uint16, data []byte) ([]byte, error)
	GraphicsCmds() GraphicsCapability
	NVGraphicsCmds() NVGraphicsCapability
	DownloadGraphicsCmds() DownloadGraphicsCapability
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands groups all bit image-related capabilities
type Commands struct {
	Graphics         GraphicsCapability
	NvGraphics       NVGraphicsCapability
	DownloadGraphics DownloadGraphicsCapability
}

// NewCommands creates a new Commands instance with initialized sub-commands
func NewCommands() *Commands {
	return &Commands{
		Graphics:         &GraphicsCommands{},
		NvGraphics:       &NvGraphicsCommands{},
		DownloadGraphics: &DownloadGraphicsCommands{},
	}
}

// GraphicsCmds returns the graphics capability.
func (c *Commands) GraphicsCmds() GraphicsCapability {
	return c.Graphics
}

// NVGraphicsCmds returns the NV graphics capability.
func (c *Commands) NVGraphicsCmds() NVGraphicsCapability {
	return c.NvGraphics
}

// DownloadGraphicsCmds returns the download graphics capability.
func (c *Commands) DownloadGraphicsCmds() DownloadGraphicsCapability {
	return c.DownloadGraphics
}

// ============================================================================
// Helper Functions
// ============================================================================

// CalculateDataLength calculates the expected data length for bit image data
func CalculateDataLength(mode Mode, width uint16) int {
	switch mode {
	case SingleDensity8, DoubleDensity8:
		return int(width)
	case SingleDensity24, DoubleDensity24:
		return int(width) * 3
	default:
		return 0
	}
}

// CalculateRasterDataLength calculates expected data length for raster format
func CalculateRasterDataLength(widthBytes, heightDots uint16) int {
	return int(widthBytes) * int(heightDots)
}

// CalculateColumnDataLength calculates expected data length for column format
func CalculateColumnDataLength(widthDots, heightBytes uint16) int {
	return int(widthDots) * int(heightBytes)
}

// ============================================================================
// Validation Helper Functions
// ============================================================================

// ValidateBitImageMode validates if bit image mode is valid
func ValidateBitImageMode(mode Mode) error {
	switch mode {
	case SingleDensity8, DoubleDensity8, SingleDensity24, DoubleDensity24:
		return nil
	default:
		return ErrBitImageMode
	}
}

// ValidatePrintMode validates if print mode is valid
func ValidatePrintMode(mode PrintMode) error {
	switch mode {
	case Normal, DoubleWidth, DoubleHeight, Quadruple,
		NormalASCII, DoubleWidthASCII, DoubleHeightASCII, QuadrupleASCII:
		return nil
	default:
		return ErrInvalidMode
	}
}

// ValidateHorizontalDots validates horizontal dot count
func ValidateHorizontalDots(dots uint16) error {
	if dots < 1 || dots > MaxHorizontalDots {
		return ErrHorizontalDotsRange
	}
	return nil
}

// ValidateNVBitImageNumber validates NV bit image number
func ValidateNVBitImageNumber(n byte) error {
	if n < 1 {
		return ErrInvalidBitImageNumber
	}
	return nil
}

// ValidateNVImageDimensions validates NV bit image dimensions
func ValidateNVImageDimensions(width, height uint16) error {
	if width < 1 || width > MaxNVImageWidth {
		return ErrInvalidImageDimensions
	}
	if height < 1 || height > MaxNVImageHeight {
		return ErrInvalidImageDimensions
	}
	return nil
}

// ValidateDownloadedImageDimensions validates downloaded bit image dimensions
func ValidateDownloadedImageDimensions(x, y byte) error {
	if x < 1 {
		return ErrInvalidHorizontalSize
	}
	if y < 1 || y > MaxDownloadedHeight {
		return ErrInvalidVerticalSize
	}
	if uint16(x)*uint16(y) > MaxDownloadedSize {
		return ErrInvalidDimensions
	}
	return nil
}

// ValidateVariableImageDimensions validates variable vertical size bit image dimensions
func ValidateVariableImageDimensions(width, height uint16, _ PrintMode) error {
	if width < 1 || width > MaxVariableWidth {
		return ErrInvalidImageWidth
	}
	if height < 1 || height > MaxVariableHeight {
		return ErrInvalidImageHeight
	}
	return nil
}

// ValidateRasterImageDimensions validates raster bit image dimensions
func ValidateRasterImageDimensions(width, height uint16, _ PrintMode) error {
	if width < 1 || width > MaxRasterWidth {
		return ErrInvalidImageWidth
	}
	if height < 1 || height > MaxRasterHeight {
		return ErrInvalidImageHeight
	}
	return nil
}
