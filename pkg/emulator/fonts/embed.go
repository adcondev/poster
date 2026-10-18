// Package fonts provides embedded font resources for the emulator.
package fonts

import (
	"embed"
	"fmt"
)

// Note: Place your .ttf files in this directory.
// Recommended fonts: monospace fonts like "DejaVu Sans Mono", "Liberation Mono",
// or "JetBrains Mono" work well for thermal printer emulation.
//
// If no fonts are embedded, the system will use basic bitmap rendering.

//go:embed *.ttf
var fs embed.FS

// LoadFontData loads font data from embedded resources
func LoadFontData(filename string) ([]byte, error) {
	data, err := fs.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("embedded font not found: %s: %w", filename, err)
	}
	return data, nil
}

// HasEmbeddedFonts checks if any fonts are embedded
func HasEmbeddedFonts() bool {
	entries, err := fs.ReadDir(".")
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			return true
		}
	}
	return false
}

// ListFonts returns a list of embedded font filenames
func ListFonts() []string {
	entries, err := fs.ReadDir(".")
	if err != nil {
		return nil
	}
	var fonts []string
	for _, entry := range entries {
		if !entry.IsDir() {
			fonts = append(fonts, entry.Name())
		}
	}
	return fonts
}

// DiagnoseEmbedding checks the embedded fonts directory.
func DiagnoseEmbedding() {
	_, _ = fs.ReadDir(".")
}
