// Package schema contiene las definiciones de estructuras para documentos de impresión
package schema

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Default values according to document.schema.json
const (
	// Profile defaults
	DefaultPaperWidth = 80
	DefaultCodeTable  = "WPC1252"
	DefaultDPI        = 203

	// Image defaults
	DefaultImagePixelWidth = 128
	DefaultImageThreshold  = 128
	DefaultImageDithering  = "atkinson"
	DefaultImageAlign      = "center"

	// Separator defaults
	DefaultSeparatorChar   = "- "
	DefaultSeparatorLength = 48

	// Cut defaults
	DefaultCutMode = "partial"
	DefaultCutFeed = 2

	// QR defaults
	DefaultQRPixelWidth = 128
	DefaultQRCorrection = "Q"
	DefaultQRAlign      = "center"
	MinQRPixelWidth     = 87

	// Barcode defaults
	DefaultBarcodeHRIPosition = "below"
	DefaultBarcodeHRIFont     = "A"
	DefaultBarcodeAlign       = "center"

	// Text defaults
	DefaultTextAlign = "left"
	DefaultTextSize  = "1x1"
	DefaultTextFont  = "A"
	DefaultNewLine   = true
)

var versionPattern = regexp.MustCompile(`^\d+\.\d+$`)

// Validate valida un documento según las reglas del schema
func (d *Document) Validate() error {
	if d.Version != "" && !versionPattern.MatchString(d.Version) {
		return fmt.Errorf("invalid version format: %s (expected: X.Y)", d.Version)
	}

	if d.Profile.Model == "" {
		return fmt.Errorf("profile.model is required")
	}

	if len(d.Commands) == 0 {
		return fmt.Errorf("document must contain at least one command")
	}

	// Validar paper_width
	validWidths := map[int]bool{58: true, 72: true, 80: true, 100: true, 112: true, 120: true}
	if d.Profile.PaperWidth != 0 && !validWidths[d.Profile.PaperWidth] {
		return fmt.Errorf("invalid paper_width: %d (valid: 58, 72, 80, 100, 112, 120)", d.Profile.PaperWidth)
	}

	// Validar DPI
	validDPI := map[int]bool{203: true, 300: true, 600: true}
	if d.Profile.DPI != 0 && !validDPI[d.Profile.DPI] {
		return fmt.Errorf("invalid DPI: %d (valid: 203, 300, 600)", d.Profile.DPI)
	}

	return nil
}

// TODO: Define all_mayus y all_bold options for commands

// Document representa un documento de impresión completo
type Document struct {
	Version  string        `json:"version"`             // Requerido: >1.0
	Profile  ProfileConfig `json:"profile"`             // Requerido: profile.model
	DebugLog bool          `json:"debug_log,omitempty"` // Default: false
	Commands []Command     `json:"commands"`            // Requerido: len > 0
}

// ToJSON convierte el documento a JSON
func (d *Document) ToJSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

// ProfileConfig configuración del perfil de impresora
type ProfileConfig struct {
	Model      string `json:"model"`                 // Requerido
	PaperWidth int    `json:"paper_width,omitempty"` // Default: 80
	CodeTable  string `json:"code_table,omitempty"`  // Default: WPC1252
	DPI        int    `json:"dpi,omitempty"`         // Default: 203
	HasQR      bool   `json:"has_qr,omitempty"`      // Default: false
}

// TODO: Define an order field for reordering or grouping commands. Check if it's worth it.

// Command represents a single command in the document
type Command struct {
	Type string          `json:"type"` // Tipo de comando
	Data json.RawMessage `json:"data"` // Datos específicos del comando
}
