package printposition

import (
	"errors"
	"fmt"
)

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for print position control.
// ESC/POS is the command system used by thermal receipt printers to control
// print positioning, justification, margins, tab positions, and print area
// configuration in both Standard and Page modes.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Justification represents text justification modes
type Justification byte

// Justification modes
const (
	Left        Justification = 0x00
	Center      Justification = 0x01
	Right       Justification = 0x02
	LeftASCII   Justification = '0'
	CenterASCII Justification = '1'
	RightASCII  Justification = '2'
)

// PrintDirection represents print direction modes in Page mode
type PrintDirection byte

const (
	// LeftToRight means upper left start
	LeftToRight PrintDirection = 0x00
	// BottomToTop means lower left start
	BottomToTop PrintDirection = 0x01
	// RightToLeft means lower right start
	RightToLeft PrintDirection = 0x02
	// TopToBottom means upper right start
	TopToBottom PrintDirection = 0x03

	// LeftToRightASCII ASCII for upper left start
	LeftToRightASCII PrintDirection = '0'
	// BottomToTopASCII ASCII for lower left start
	BottomToTopASCII PrintDirection = '1'
	// RightToLeftASCII ASCII for lower right start
	RightToLeftASCII PrintDirection = '2'
	// TopToBottomASCII ASCII for upper right start
	TopToBottomASCII PrintDirection = '3'
)

// BeginLine of line operations
type BeginLine byte

const (
	// Erase clears the print buffer from the current position to the end of the line
	Erase BeginLine = 0x00
	// Print prints the buffer from the current position to the end of the line
	Print BeginLine = 0x01

	// EraseASCII clears the print buffer from the current position to the end of the line
	EraseASCII BeginLine = '0'
	// PrintASCII prints the buffer from the current position to the end of the line
	PrintASCII BeginLine = '1'
)

// Tab position limits
const (
	MaxTabPositions = 32
	MaxTabValue     = 255
)

// ============================================================================
// Error Definitions
// ============================================================================

// ErrJustification represents an invalid justification mode error
var (
	ErrJustification       = errors.New("invalid justification mode (try 0-2 or '0'..'2')")
	ErrPrintDirection      = errors.New("invalid print direction (try 0-3 or '0'..'3')")
	ErrBeginLineMode       = errors.New("invalid begin line mode (try 0-1 or '0'..'1')")
	ErrTooManyTabPositions = fmt.Errorf("too many tab positions (max %d)", MaxTabPositions)
	ErrTabPosition         = errors.New("invalid tab position (must be 1-255 in ascending order)")
	ErrPrintAreaWidthSize  = errors.New("invalid print area size (width must be >= 1)")
	ErrPrintAreaHeightSize = errors.New("invalid print area size (height must be >= 1)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ Capability = (*Commands)(nil)

// Capability defines the interface for print position commands
type Capability interface {
	SetAbsolutePrintPosition(position uint16) []byte
	SetRelativePrintPosition(distance int16) []byte
	HorizontalTab() []byte
	SetHorizontalTabPositions(positions []byte) ([]byte, error)
	SelectJustification(mode Justification) ([]byte, error)
	SetLeftMargin(margin uint16) []byte
	SetPrintAreaWidth(width uint16) []byte
	SetPrintPositionBeginningLine(mode BeginLine) ([]byte, error)
	SelectPrintDirectionPageMode(direction PrintDirection) ([]byte, error)
	SetPrintAreaPageMode(x, y, width, height uint16) ([]byte, error)
	SetAbsoluteVerticalPrintPosition(position uint16) []byte
	SetRelativeVerticalPrintPosition(distance int16) []byte
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands implements the Capability interface for print position commands
type Commands struct{}

// NewCommands creates a new instance of print position Commands
func NewCommands() *Commands {
	return &Commands{}
}

// ============================================================================
// Helper Functions
// ============================================================================

// (Add any private helper functions here if needed)

// ============================================================================
// Validation Helper Functions
// ============================================================================

// ValidateJustification validates if justification mode is valid.
func ValidateJustification(mode Justification) error {
	switch mode {
	case Left, Center, Right,
		LeftASCII, CenterASCII, RightASCII:
		return nil
	default:
		return ErrJustification
	}
}

// ValidatePrintDirection validates if print direction is valid.
func ValidatePrintDirection(direction PrintDirection) error {
	switch direction {
	case LeftToRight, BottomToTop, RightToLeft, TopToBottom,
		LeftToRightASCII, BottomToTopASCII, RightToLeftASCII, TopToBottomASCII:
		return nil
	default:
		return ErrPrintDirection
	}
}

// ValidateBeginLineMode validates if begin line mode is valid.
func ValidateBeginLineMode(mode BeginLine) error {
	switch mode {
	case Erase, Print, EraseASCII, PrintASCII:
		return nil
	default:
		return ErrBeginLineMode
	}
}

// ValidatePrintArea validates print area dimensions.
func ValidatePrintArea(width, height uint16) error {
	if width == 0 {
		return ErrPrintAreaWidthSize
	}
	if height == 0 {
		return ErrPrintAreaHeightSize
	}
	return nil
}
