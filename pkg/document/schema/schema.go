// Package schema contiene las definiciones de estructuras para documentos de impresión
package schema

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/adcondev/pos-printer/pkg/constants"
)

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

// Version pattern: X.Y where X and Y are digits
var versionPattern = regexp.MustCompile(`^\d+\.\d+$`)

// Validate checks the document structure and returns an error if invalid
func (d *Document) Validate() error {
	if d.Version == "" {
		return fmt.Errorf("version is required")
	}
	if !versionPattern.MatchString(d.Version) {
		return fmt.Errorf("invalid version format: %s (expected X.Y pattern, e.g., '1.0')", d.Version)
	}

	if d.Profile.Model == "" {
		return fmt.Errorf("profile.model is required")
	}

	// Use constants for validation
	if d.Profile.PaperWidth != 0 && !isValidPaperWidth(d.Profile.PaperWidth) {
		return fmt.Errorf("invalid paper_width: %d (valid values: %v)",
			d.Profile.PaperWidth, constants.ValidPaperWidths)
	}

	if d.Profile.DPI != 0 && !isValidDPI(d.Profile.DPI) {
		return fmt.Errorf("invalid dpi: %d (valid values: %v)",
			d.Profile.DPI, constants.ValidDPIs)
	}

	if len(d.Commands) == 0 {
		return fmt.Errorf("document must contain at least one command")
	}

	return nil
}

func isValidPaperWidth(width int) bool {
	for _, valid := range constants.ValidPaperWidths {
		if width == valid {
			return true
		}
	}
	return false
}

func isValidDPI(dpi int) bool {
	for _, valid := range constants.ValidDPIs {
		if dpi == valid {
			return true
		}
	}
	return false
}
