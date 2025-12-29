//nolint:revive,nolintlint
package shared

import "errors"

// ============================================================================
// Context
// ============================================================================
// This package provides shared utilities and constants for ESC/POS command packages.
// It includes buffer validation functions, little-endian byte conversion utilities,
// and shared error definitions used across multiple ESC/POS command implementations.

// ============================================================================
// Constant and Var Definitions
// ============================================================================

// Buffer limits
var (
	// MinBuf es el tamaño mínimo del buffer
	MinBuf = 1
	// MaxBuf es el tamaño máximo del buffer
	MaxBuf = 65535
)

var (
	// ErrBufferOverflow buffer is too large
	ErrBufferOverflow = errors.New("can't print overflowed buffer (protocol max 64KB; model may be lower)")
	// ErrEmptyBuffer buffer is empty
	ErrEmptyBuffer = errors.New("can't print an empty buffer")
)

// ============================================================================
// Validation and Helper Functions
// ============================================================================

// IsBufLenOk validates if the buffer size is within acceptable limits.
func IsBufLenOk(buf []byte) error {
	if len(buf) < MinBuf {
		return ErrEmptyBuffer
	}
	if len(buf) > MaxBuf {
		return ErrBufferOverflow
	}
	return nil
}

// ToLittleEndian convierte una longitud en dos bytes little-endian (dL,dH) para usar en ESCPOS.
func ToLittleEndian(number uint16) (nL, nH byte) {
	nL = byte(number & 0xFF)        // byte de menor peso
	nH = byte((number >> 8) & 0xFF) // byte de mayor peso
	return nL, nH
}

// ToLittleEndian32 convierte una longitud en cuatro bytes little-endian (nL, nH, nHH, nHHH) para usar en ESCPOS.
func ToLittleEndian32(number uint32) (nL, nH, nHH, nHHH byte) {
	nL = byte(number & 0xFF)           // byte de menor peso
	nH = byte((number >> 8) & 0xFF)    // segundo byte
	nHH = byte((number >> 16) & 0xFF)  // tercer byte
	nHHH = byte((number >> 24) & 0xFF) // byte de mayor peso
	return nL, nH, nHH, nHHH
}
