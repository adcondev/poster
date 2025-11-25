// Package document provides structures and functions to build print documents.
package document

import (
	"encoding/json"
	"fmt"
)

// PrintJob representa los datos de un documento de impresión
type PrintJob struct {
	Data Document `json:"data"`
}

// TODO: Define all_mayus y all_bold options for commands

// Document representa un documento de impresión completo
type Document struct {
	Version  string        `json:"version"`             // Requerido: >1.0
	Profile  ProfileConfig `json:"profile"`             // Requerido: profile.model
	DebugLog bool          `json:"debug_log,omitempty"` // Default: false
	Commands []Command     `json:"commands"`            // Requerido: len > 0
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

// ImageCommand represents an image command
type ImageCommand struct {
	Code       string `json:"code"`                  // Base64 - Requerido
	Format     string `json:"format,omitempty"`      // png, jpg, bmp
	PixelWidth int    `json:"pixel_width,omitempty"` // Default: 128
	Align      string `json:"align,omitempty"`       // Default: center
	Threshold  byte   `json:"threshold,omitempty"`   // Default: 128
	Dithering  string `json:"dithering,omitempty"`   // Default: atkinson
	Scaling    string `json:"scaling,omitempty"`     // Default: bilinear
}

// SeparatorCommand represents a separator command
type SeparatorCommand struct {
	Char   string `json:"char,omitempty"`   // Default: "-"
	Length int    `json:"length,omitempty"` // Default: 48
}

// FeedCommand represents a feed command
type FeedCommand struct {
	Lines int `json:"lines"` // Requerido (1-255)
}

// CutCommand represents a cut command
type CutCommand struct {
	Mode string `json:"mode,omitempty"` // full, partial
	Feed int    `json:"feed,omitempty"` // Líneas antes del corte
}

// ParseDocument parsea un documento JSON
func ParseDocument(data []byte) (*Document, error) {
	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	// Validación básica
	if doc.Version == "" {
		// TODO: Review an smart way to handle versioning
		doc.Version = "1.0"
	}

	if len(doc.Commands) == 0 {
		return nil, fmt.Errorf("document must contain at least one command")
	}

	return &doc, nil
}

// ToJSON convierte el documento a JSON
func (d *Document) ToJSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}
