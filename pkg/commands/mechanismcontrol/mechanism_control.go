package mechanismcontrol

import (
	"errors"
)

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for mechanism control functionality.
// ESC/POS is the command system used by thermal receipt printers to control
// physical printer mechanisms such as paper cutting, print head movement,
// paper feeding, and unidirectional print mode settings.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Type definitions and Constants

// UnidirectionalMode represents the unidirectional print mode setting
type UnidirectionalMode byte

const (
	// UnidirOff turns off unidirectional mode (bidirectional on)
	UnidirOff UnidirectionalMode = 0x00
	// UnidirOn turns on unidirectional mode
	UnidirOn UnidirectionalMode = 0x01
)

// CutMode represents the paper cut mode for simple cut commands
type CutMode byte

const (
	// CutModeFull performs a full cut
	CutModeFull CutMode = 0x00
	// CutModePartial performs a partial cut
	CutModePartial CutMode = 0x01
	// CutModeFullASCII performs a full cut (ASCII mode)
	CutModeFullASCII CutMode = '0'
	// CutModePartialASCII performs a partial cut (ASCII mode)
	CutModePartialASCII CutMode = '1'
)

// CutType represents the paper cut type for extended cut commands
type CutType byte

const (
	// CutTypeFull performs a full cut
	CutTypeFull CutType = 0x00
	// CutTypePartial performs a partial cut
	CutTypePartial CutType = 0x01
)

// FeedCut represents the feed and cut mode
type FeedCut byte

const (
	// FeedCutFull performs a full cut after feeding
	FeedCutFull FeedCut = 65
	// FeedCutPartial performs a partial cut after feeding
	FeedCutPartial FeedCut = 66
)

// PositionCut represents the position cut mode
type PositionCut byte

const (
	// PositionCutFull performs a full cut at position
	PositionCutFull PositionCut = 97
	// PositionCutPartial performs a partial cut at position
	PositionCutPartial PositionCut = 98
)

// FeedCutReturn represents the feed, cut and return mode
type FeedCutReturn byte

const (
	// FeedCutReturnFull performs a full cut with return
	FeedCutReturnFull FeedCutReturn = 103
	// FeedCutReturnPartial performs a partial cut with return
	FeedCutReturnPartial FeedCutReturn = 104
)

const (
	// MaxFeedAmount represents the maximum feed amount before cutting
	MaxFeedAmount byte = 255
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	// ErrCutMode indicates an invalid cut mode
	ErrCutMode = errors.New("invalid cut mode (try 0, 1, 48, or 49)")
	// ErrCutType indicates an invalid cut type
	ErrCutType = errors.New("invalid cut type (try 0 or 1)")
	// ErrFeedCutMode indicates an invalid feed and cut mode
	ErrFeedCutMode = errors.New("invalid feed and cut mode (try 65 or 66)")
	// ErrPositionCutMode indicates an invalid position cut mode
	ErrPositionCutMode = errors.New("invalid position cut mode (try 97 or 98)")
	// ErrFeedCutReturnMode indicates an invalid feed, cut and return mode
	ErrFeedCutReturnMode = errors.New("invalid feed, cut and return mode (try 103 or 104)")
	// ErrInvalidDrawerPin indicates an invalid drawer pin
	ErrInvalidDrawerPin = errors.New("invalid drawer pin (must be 0 or 1)")
	// ErrInvalidBeepParams indicates invalid beep parameters
	ErrInvalidBeepParams = errors.New("invalid beep parameters (times and lapse must be >= 1)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Compile-time check that Commands implements Capability
var _ Capability = (*Commands)(nil)

// Capability defines the interface for mechanism control commands
type Capability interface {
	ReturnHome() []byte
	SetUnidirectionalPrintMode(mode UnidirectionalMode) []byte
	PartialCut() []byte
	PartialCutThreePoints() []byte
	PaperCut(cutType CutType) ([]byte, error)
	PaperFeedAndCut(cutType CutType, feedAmount byte) ([]byte, error)
	ReservePaperCut(cutType CutType, feedAmount byte) ([]byte, error)
	CutPaper(mode CutMode) ([]byte, error)
	FeedAndCutPaper(mode FeedCut, feedAmount byte) ([]byte, error)
	SetCutPosition(mode PositionCut, position byte) ([]byte, error)
	FeedCutAndReturnPaper(mode FeedCutReturn, feedAmount byte) ([]byte, error)
	Pulse(pin byte, onTime byte, offTime byte) ([]byte, error)
	Beep(times byte, lapse byte) ([]byte, error)
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands implements the Capability interface for mechanism control
type Commands struct{}

// NewCommands creates a new Commands instance
func NewCommands() *Commands {
	return &Commands{}
}

// ============================================================================
// Validation Helper Functions
// ============================================================================

// ValidateCutMode validates if cut mode is valid
func ValidateCutMode(mode CutMode) error {
	switch mode {
	case CutModeFull, CutModePartial, CutModeFullASCII, CutModePartialASCII:
		return nil
	default:
		return ErrCutMode
	}
}

// ValidateCutType validates if cut type is valid
func ValidateCutType(cutType CutType) error {
	switch cutType {
	case CutTypeFull, CutTypePartial:
		return nil
	default:
		return ErrCutType
	}
}

// ValidateFeedCutMode validates if feed cut mode is valid
func ValidateFeedCutMode(mode FeedCut) error {
	switch mode {
	case FeedCutFull, FeedCutPartial:
		return nil
	default:
		return ErrFeedCutMode
	}
}

// ValidatePositionCutMode validates if position cut mode is valid
func ValidatePositionCutMode(mode PositionCut) error {
	switch mode {
	case PositionCutFull, PositionCutPartial:
		return nil
	default:
		return ErrPositionCutMode
	}
}

// ValidateFeedCutReturnMode validates if feed cut return mode is valid
func ValidateFeedCutReturnMode(mode FeedCutReturn) error {
	switch mode {
	case FeedCutReturnFull, FeedCutReturnPartial:
		return nil
	default:
		return ErrFeedCutReturnMode
	}
}
