package barcode

import (
	"errors"

	"github.com/adcondev/poster/pkg/commands/common"
)

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for barcode printing functionality.
// ESC/POS is the command system used by thermal receipt printers to control
// barcode generation, HRI (Human Readable Interpretation) text positioning,
// and various barcode symbologies including UPC, EAN, CODE128, and GS1 DataBar

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Type definitions

// HRIPosition represents the print positions for HRI characters
type HRIPosition byte

// HRIFont represents the font types for HRI characters
type HRIFont byte

// Height represents the barcode height (in dots)
type Height byte

// Width represents the horizontal module width of the barcode
type Width byte

// Symbology represents the barcode symbology types
type Symbology byte

// Code128Set represents the code sets for CODE128 (A/B/C)
type Code128Set byte

// Constants

const (
	// HRINotPrinted represents HRI not printed (numeric mode)
	HRINotPrinted HRIPosition = 0x00
	// HRIAbove represents HRI printed above barcode
	HRIAbove HRIPosition = 0x01
	// HRIBelow represents HRI printed below barcode
	HRIBelow HRIPosition = 0x02
	// HRIBoth represents HRI printed above and below barcode
	HRIBoth HRIPosition = 0x03

	// HRINotPrintedASCII represents HRI not printed (ASCII mode)
	HRINotPrintedASCII HRIPosition = '0'
	// HRIAboveASCII represents HRI printed above barcode (ASCII mode)
	HRIAboveASCII HRIPosition = '1'
	// HRIBelowASCII represents HRI printed below barcode (ASCII mode)
	HRIBelowASCII HRIPosition = '2'
	// HRIBothASCII represents HRI printed above and below barcode (ASCII mode)
	HRIBothASCII HRIPosition = '3'
)

const (
	// HRIFontA represents HRI font type A (numeric mode)
	HRIFontA HRIFont = 0x00
	// HRIFontB represents HRI font type B
	HRIFontB HRIFont = 0x01
	// HRIFontC represents HRI font type C
	HRIFontC HRIFont = 0x02
	// HRIFontD represents HRI font type D
	HRIFontD HRIFont = 0x03
	// HRIFontE represents HRI font type E
	HRIFontE HRIFont = 0x04

	// HRIFontAASCII represents HRI font type A (ASCII mode)
	HRIFontAASCII HRIFont = '0'
	// HRIFontBASCII represents HRI font type B (ASCII mode)
	HRIFontBASCII HRIFont = '1'
	// HRIFontCASCII represents HRI font type C (ASCII mode)
	HRIFontCASCII HRIFont = '2'
	// HRIFontDASCII represents HRI font type D (ASCII mode)
	HRIFontDASCII HRIFont = '3'
	// HRIFontEASCII represents HRI font type E (ASCII mode)
	HRIFontEASCII HRIFont = '4'

	// HRISpecialFontA represents special HRI font A
	HRISpecialFontA HRIFont = 97
	// HRISpecialFontB represents special HRI font B
	HRISpecialFontB HRIFont = 98
)

const (
	// MinHeight represents minimum barcode height
	MinHeight Height = 1
	// MaxHeight represents maximum barcode height
	MaxHeight Height = 255
	// DefaultHeight represents default barcode height (model-dependent)
	DefaultHeight Height = 162
)

const (
	// MinWidth represents minimum barcode width
	MinWidth Width = 2
	// MaxWidth represents maximum barcode width
	MaxWidth Width = 6
	// DefaultWidth represents default barcode width
	DefaultWidth Width = 3

	// ExtendedMinWidth represents extended minimum width (model-dependent)
	ExtendedMinWidth Width = 68
	// ExtendedMaxWidth represents extended maximum width (model-dependent)
	ExtendedMaxWidth Width = 76
)

const (
	// UPCA represents UPC-A barcode symbology (11-12 digits)
	UPCA Symbology = 0
	// UPCE represents UPC-E barcode symbology (6-8, 11-12 digits)
	UPCE Symbology = 1
	// JAN13 represents JAN13/EAN13 barcode symbology (12-13 digits)
	JAN13 Symbology = 2
	// JAN8 represents JAN8/EAN8 barcode symbology (7-8 digits)
	JAN8 Symbology = 3
	// CODE39 represents CODE39 barcode symbology (variable length)
	CODE39 Symbology = 4
	// ITF represents Interleaved 2 of 5 barcode symbology (pairs of digits)
	ITF Symbology = 5
	// CODABAR represents CODABAR/NW-7 barcode symbology (variable length)
	CODABAR Symbology = 6
)

const (
	// UPCAB represents UPC-A barcode symbology (11-12 digits, length-prefixed)
	UPCAB Symbology = 65
	// UPCEB represents UPC-E barcode symbology (6-8, 11-12 digits)
	UPCEB Symbology = 66
	// EAN13 represents EAN13 barcode symbology (12-13 digits)
	EAN13 Symbology = 67
	// EAN8 represents EAN8 barcode symbology (7-8 digits)
	EAN8 Symbology = 68
	// CODE39B represents CODE39 barcode symbology (1-255 chars)
	CODE39B Symbology = 69
	// ITFB represents ITF barcode symbology (2-254 pairs)
	ITFB Symbology = 70
	// CODABARB represents CODABAR barcode symbology (2-255 chars)
	CODABARB Symbology = 71
	// CODE93 represents CODE93 barcode symbology (1-255 chars)
	CODE93 Symbology = 72
	// CODE128 represents CODE128 barcode symbology (2-255 bytes)
	CODE128 Symbology = 73
	// GS1128 represents GS1-128 barcode symbology (2-255 bytes)
	GS1128 Symbology = 74
	// GS1DataBarOmni represents GS1 DataBar Omnidirectional (13 digits)
	GS1DataBarOmni Symbology = 75
	// GS1DataBarTrunc represents GS1 DataBar Truncated (13 digits)
	GS1DataBarTrunc Symbology = 76
	// GS1DataBarLim represents GS1 DataBar Limited (13 digits)
	GS1DataBarLim Symbology = 77
	// GS1DataBarExp represents GS1 DataBar Expanded (2-255 chars)
	GS1DataBarExp Symbology = 78
	// CODE128Auto represents CODE128 Auto barcode symbology (1-255 bytes)
	CODE128Auto Symbology = 79
)

const (
	// Code128SetA represents CODE128 Set A (ASCII 0-95)
	Code128SetA Code128Set = 65
	// Code128SetB represents CODE128 Set B (ASCII 32-127)
	Code128SetB Code128Set = 66
	// Code128SetC represents CODE128 Set C (numeric pairs 00-99)
	Code128SetC Code128Set = 67
)

// ============================================================================
// Error Variables
// ============================================================================

// ErrHRIPosition represents an invalid HRI position error
var (
	ErrHRIPosition      = errors.New("invalid HRI position (try 0-3 or '0'..'3')")
	ErrHRIFont          = errors.New("invalid HRI font (try 0-4, '0'..'4', 97, or 98)")
	ErrHeight           = errors.New("invalid barcode height (try 1-255)")
	ErrWidth            = errors.New("invalid barcode width (try 2-6 or 68-76 for extended)")
	ErrSymbology        = errors.New("invalid barcode symbology")
	ErrDataTooShort     = errors.New("barcode data too short")
	ErrDataTooLong      = errors.New("barcode data too long (max 255 bytes)")
	ErrOddITFLength     = errors.New("ITF barcode requires even number of digits")
	ErrCode128Set       = errors.New("invalid CODE128 code set (try 65-67)")
	ErrCode128NoCodeSet = errors.New("CODE128 requires code set specification")
)

// ============================================================================
// Main Implementation
// ============================================================================

// Commands implements the Capability interface for barcode commands
type Commands struct{}

// NewCommands creates a new instance of Commands
func NewCommands() *Commands {
	return &Commands{}
}

// ============================================================================
// Interface Definitions
// ============================================================================

// Compile-time check that Commands implements Capability
var _ Capability = (*Commands)(nil)

// Capability groups barcode-related capabilities
type Capability interface {
	// HRI settings
	SelectHRICharacterPosition(position HRIPosition) ([]byte, error)
	SelectFontForHRI(font HRIFont) ([]byte, error)

	// Dimensions
	SetBarcodeHeight(height Height) ([]byte, error)
	SetBarcodeWidth(width Width) ([]byte, error)

	// Printing
	PrintBarcode(symbology Symbology, data []byte) ([]byte, error)
	PrintBarcodeWithCodeSet(symbology Symbology, codeSet Code128Set, data []byte) ([]byte, error)
}

// ============================================================================
// Helper Functions
// ============================================================================

// buildFunctionA builds Function A command (NUL-terminated)
func (c *Commands) buildFunctionA(symbology Symbology, data []byte) ([]byte, error) {
	// Basic validations for Function A symbologies
	if symbology == ITF {
		if len(data)%2 != 0 {
			return nil, ErrOddITFLength
		}
	}

	// Build command: GS k m data... NUL
	cmd := []byte{common.GS, 'k', byte(symbology)}
	cmd = append(cmd, data...)
	cmd = append(cmd, common.NUL)
	return cmd, nil
}

// buildFunctionB builds Function B command (length-prefixed)
func (c *Commands) buildFunctionB(symbology Symbology, data []byte) ([]byte, error) {
	// Validate data length (maximum 255 for one-byte length)
	if len(data) > 255 {
		return nil, ErrDataTooLong
	}

	// Special validation for certain symbologies
	switch symbology {
	case CODE128, GS1128:
		// Check if data has the required code set prefix
		if len(data) < 2 || data[0] != '{' ||
			data[1] < byte(Code128SetA) || data[1] > byte(Code128SetC) {
			return nil, ErrCode128NoCodeSet
		}
	case ITFB:
		// ITF requires even number of digits
		if len(data)%2 != 0 {
			return nil, ErrOddITFLength
		}
	}

	// Build command: GS k m n data...
	cmd := []byte{common.GS, 'k', byte(symbology), byte(len(data))}
	cmd = append(cmd, data...)
	return cmd, nil
}

// ============================================================================
// Utility Functions for Validation
// ============================================================================

// ValidateHRIPosition validates if HRI position is valid
func ValidateHRIPosition(position HRIPosition) error {
	switch position {
	case HRINotPrinted, HRIAbove, HRIBelow, HRIBoth,
		HRINotPrintedASCII, HRIAboveASCII, HRIBelowASCII, HRIBothASCII:
		return nil
	default:
		return ErrHRIPosition
	}
}

// ValidateHRIFont validates if HRI font is valid
func ValidateHRIFont(font HRIFont) error {
	switch font {
	case HRIFontA, HRIFontB, HRIFontC, HRIFontD, HRIFontE,
		HRIFontAASCII, HRIFontBASCII, HRIFontCASCII, HRIFontDASCII, HRIFontEASCII,
		HRISpecialFontA, HRISpecialFontB:
		return nil
	default:
		return ErrHRIFont
	}
}

// ValidateNumericData checks if all bytes are numeric digits
func ValidateNumericData(data []byte) bool {
	for _, b := range data {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// ValidateCode39Data checks if all bytes are valid for CODE39
func ValidateCode39Data(data []byte) bool {
	for _, b := range data {
		switch {
		case b >= '0' && b <= '9':
		case b >= 'A' && b <= 'Z':
		case b == ' ' || b == '$' || b == '%' || b == '*' ||
			b == '+' || b == '-' || b == '.' || b == '/':
		default:
			return false
		}
	}
	return true
}

// ValidateCodabarData checks if CODABAR data is valid
func ValidateCodabarData(data []byte) bool {
	if len(data) < 2 {
		return false
	}
	// Check start character
	start := data[0]
	if (start < 'A' || start > 'D') && (start < 'a' || start > 'd') {
		return false
	}
	// Check stop character
	stop := data[len(data)-1]
	if (stop < 'A' || stop > 'D') && (stop < 'a' || stop > 'd') {
		return false
	}
	return true
}
