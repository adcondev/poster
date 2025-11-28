// Package schema contiene las definiciones de estructuras para documentos de impresión
package schema

import (
	"encoding/json"
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
