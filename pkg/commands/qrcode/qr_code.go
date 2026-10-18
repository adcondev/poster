package qrcode

import (
	"errors"
	"fmt"
)

// ============================================================================
// Context
// ============================================================================
// This package implements ESC/POS commands for QR Code generation and printing.
// ESC/POS is the command system used by thermal receipt printers to control
// QR Code symbol encoding, storage, and printing operations.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Model is QR Code model selection
type Model byte

// Model constants - QR Code model types
const (
	// Model1 is original QR specification
	Model1 Model = 49
	// Model2 is the most common QR Code model
	Model2 Model = 50
	// MicroQR is the smaller version
	MicroQR Model = 51
)

// ModuleSize is the size of QR Code modules (dots)
type ModuleSize byte

const (
	// MinModuleSize represents the minimum module size for QR Codes (1 dot)
	MinModuleSize ModuleSize = 1
	// DefaultModuleSize represents the default module size for QR Codes (3 dots)
	DefaultModuleSize ModuleSize = 3
	// MaxModuleSize represents the maximum module size for QR Codes (16 dots)
	MaxModuleSize ModuleSize = 16
)

// ErrorCorrection is the error correction level for QR Codes
type ErrorCorrection byte

const (
	// LevelL (~7% recovery)
	LevelL ErrorCorrection = 48
	// LevelM (~15% recovery)
	LevelM ErrorCorrection = 49
	// LevelQ (~25% recovery)
	LevelQ ErrorCorrection = 50
	// LevelH (~30% recovery)
	LevelH ErrorCorrection = 51
)

// Data limits
const (
	MinDataLength = 1    // Minimum data length
	MaxDataLength = 7089 // Maximum data length (7092 - 3 header bytes)
)

// ============================================================================
// Error Definitions
// ============================================================================

var (
	// ErrQRModel indicates an invalid QR Code model
	ErrQRModel = errors.New("invalid QR model (try 49-51)")
	// ErrParameter indicates an invalid parameter
	ErrParameter = errors.New("invalid parameter (must be 0)")
	// ErrModuleSize indicates an invalid module size
	ErrModuleSize = errors.New("invalid module size (try 1-16)")
	// ErrErrorCorrection indicates an invalid error correction level
	ErrErrorCorrection = errors.New("invalid error correction level (try 48-51)")
	// ErrDataTooShort indicates data is too short
	ErrDataTooShort = errors.New("data too short (minimum 1 byte)")
	// ErrDataTooLong indicates data is too long
	ErrDataTooLong = errors.New("data too long (maximum 7089 bytes)")
)

// ============================================================================
// Interface Definitions
// ============================================================================

// Interface compliance check
var _ Capability = (*Commands)(nil)

// Capability defines the QR Code printing interface
type Capability interface {
	SelectQRCodeModel(n1 Model, n2 byte) ([]byte, error)
	SetQRCodeModuleSize(n ModuleSize) ([]byte, error)
	SetQRCodeErrorCorrectionLevel(n ErrorCorrection) ([]byte, error)
	StoreQRCodeData(data []byte) ([]byte, error)
	PrintQRCode() []byte
	GetQRCodeSize() []byte
}

// ============================================================================
// Main Implementation
// ============================================================================

// Commands implements QR Code ESC/POS commands
type Commands struct {
	// No sub-modules needed for QR Code
}

// NewCommands creates a new QR Code commands instance
func NewCommands() *Commands {
	return &Commands{}
}

// ============================================================================
// Helper Functions
// ============================================================================

// IsNumericData checks if data can be encoded in Numeric mode
func IsNumericData(data []byte) bool {
	for _, b := range data {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// IsAlphanumericData checks if data can be encoded in Alphanumeric mode
func IsAlphanumericData(data []byte) bool {
	// Alphanumeric mode supports: 0-9, A-Z, space, $, %, *, +, -, ., /, :
	for _, b := range data {
		switch {
		case b >= '0' && b <= '9':
			continue
		case b >= 'A' && b <= 'Z':
			continue
		case b == ' ' || b == '$' || b == '%' || b == '*' ||
			b == '+' || b == '-' || b == '.' || b == '/' || b == ':':
			continue
		default:
			return false
		}
	}
	return true
}

// ============================================================================
// Validation Functions
// ============================================================================

// ValidateQRModel validates if the QR Code model is valid
func ValidateQRModel(model Model) error {
	if model < Model1 || model > MicroQR {
		return fmt.Errorf("%w: %d", ErrQRModel, model)
	}
	return nil
}

// ValidateModuleSize validates if the module size is valid
func ValidateModuleSize(size ModuleSize) error {
	if size < MinModuleSize || size > MaxModuleSize {
		return fmt.Errorf("%w: %d", ErrModuleSize, size)
	}
	return nil
}

// ValidateErrorCorrection validates if the error correction level is valid
func ValidateErrorCorrection(level ErrorCorrection) error {
	if level < LevelL || level > LevelH {
		return fmt.Errorf("%w: %d", ErrErrorCorrection, level)
	}
	return nil
}

// ValidateDataLength validates if the data length is within bounds
func ValidateDataLength(data []byte) error {
	if len(data) < MinDataLength {
		return fmt.Errorf("%w: %d bytes", ErrDataTooShort, len(data))
	}
	if len(data) > MaxDataLength {
		return fmt.Errorf("%w: %d bytes", ErrDataTooLong, len(data))
	}
	return nil
}
