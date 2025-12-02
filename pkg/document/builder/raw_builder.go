package builder

import (
	"encoding/hex"
)

// RawBuilder constructs raw ESC/POS commands
type RawBuilder struct {
	parent   *DocumentBuilder
	hexData  string
	format   string
	comment  string
	safeMode bool
}

type rawCommand struct {
	Hex      string `json:"hex"`
	Format   string `json:"format,omitempty"`
	Comment  string `json:"comment,omitempty"`
	SafeMode bool   `json:"safe_mode,omitempty"`
}

func newRawBuilder(parent *DocumentBuilder, hexData string) *RawBuilder {
	return &RawBuilder{
		parent:  parent,
		hexData: hexData,
		format:  "hex",
	}
}

// FromBytes creates raw command from byte slice
func (b *DocumentBuilder) FromBytes(data []byte) *RawBuilder {
	return &RawBuilder{
		parent:  b,
		hexData: hex.EncodeToString(data),
		format:  "hex",
	}
}

// Base64 sets format to base64
func (rb *RawBuilder) Base64() *RawBuilder {
	rb.format = "base64"
	return rb
}

// Comment adds documentation
func (rb *RawBuilder) Comment(comment string) *RawBuilder {
	rb.comment = comment
	return rb
}

// Safe enables safety checks
func (rb *RawBuilder) Safe() *RawBuilder {
	rb.safeMode = true
	return rb
}

// End finishes the raw command
func (rb *RawBuilder) End() *DocumentBuilder {
	cmd := rawCommand{
		Hex:      rb.hexData,
		Format:   rb.format,
		Comment:  rb.comment,
		SafeMode: rb.safeMode,
	}
	return rb.parent.addCommand("raw", cmd)
}
