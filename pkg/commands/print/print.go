package print

import (
	"errors"
	"fmt"

	"github.com/adcondev/poster/pkg/commands/shared"
)

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for basic printing operations.
// ESC/POS is the command system used by thermal receipt printers to control
// text printing, paper feeding, line feeds, and print buffer management
// in both Standard and Page modes.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Control characters used in printing commands
const (
	// LF (Line FeedLines)
	LF byte = 0x0A
	// CR (Carriage Return)
	CR byte = 0x0D
	// FF (Form FeedLines)
	FF byte = 0x0C
	// CAN (Cancel)
	CAN byte = 0x18
)

// Reverse motion units and lines
var (
	// MaxReverseMotionUnits is the maximum number of motion units for reverse printing
	MaxReverseMotionUnits byte = 48
	// MaxReverseFeedLines is the maximum number of lines for reverse printing
	MaxReverseFeedLines byte = 2
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	// ErrEmptyText indicates that the provided text is empty
	ErrEmptyText = shared.ErrEmptyBuffer
	// ErrTextTooLarge indicates that the provided text exceeds buffer limits
	ErrTextTooLarge = shared.ErrBufferOverflow
	// ErrReverseUnits invalid number of motion units for reverse print
	ErrReverseUnits = fmt.Errorf("invalid reverse feed units (try 0-%d)", MaxReverseMotionUnits)
	// ErrReverseLines invalid number of lines for reverse print
	ErrReverseLines = fmt.Errorf("invalid reverse feed lines (try 0-%d)", MaxReverseFeedLines)
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ Capability = (*Commands)(nil)

// Capability defines the interface for print commands
type Capability interface {
	// Text operations
	Text(text string) ([]byte, error)
	PrintAndLineFeed() []byte
	PrintAndCarriageReturn() []byte
	FormFeed() []byte
	PrintAndFeedPaper(units byte) []byte
	PrintAndFeedLines(lines byte) []byte
	PrintAndReverseFeed(units byte) ([]byte, error)
	PrintAndReverseFeedLines(lines byte) ([]byte, error)
	PrintDataInPageMode() []byte
	CancelData() []byte
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands implements the Capability interface for print commands
type Commands struct{}

// NewCommands creates a new instance of print Commands
func NewCommands() *Commands {
	return &Commands{}
}

// Text formats and sends a string for printing.
//
// Description:
//
//	Converts a string to bytes and applies ESC/POS formatting.
//
// Notes:
//   - Replaces '\n' with LF (0x0A)
//   - Replaces '\r' with CR (0x0D)
//   - Replaces '\t' with HT (0x09)
//   - Validates buffer size according to printer limitations
func (c *Commands) Text(n string) ([]byte, error) {
	if err := shared.IsBufLenOk([]byte(n)); err != nil {
		switch {
		case errors.Is(err, shared.ErrEmptyBuffer):
			return nil, ErrEmptyText
		case errors.Is(err, shared.ErrBufferOverflow):
			return nil, ErrTextTooLarge
		default:
			return nil, err
		}
	}

	return Formatting([]byte(n)), nil
}

// ============================================================================
// Helper Functions
// ============================================================================

// Formatting replaces specific characters in the byte slice with their ESC/POS equivalents.
func Formatting(data []byte) []byte {
	formatted := make([]byte, len(data))
	copy(formatted, data)

	for i := range formatted {
		switch formatted[i] {
		case '\n':
			formatted[i] = LF
		case '\r':
			formatted[i] = CR
		case '\t':
			formatted[i] = shared.HT
		}
	}
	return formatted
}

// ============================================================================
// Validation Helper Functions
// ============================================================================

// ValidateReverseFeedUnits validates reverse feed units are within limits.
func ValidateReverseFeedUnits(units byte) error {
	if units > MaxReverseMotionUnits {
		return ErrReverseUnits
	}
	return nil
}

// ValidateReverseFeedLines validates reverse feed lines are within limits.
func ValidateReverseFeedLines(lines byte) error {
	if lines > MaxReverseFeedLines {
		return ErrReverseLines
	}
	return nil
}
