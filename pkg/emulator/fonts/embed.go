// Package fonts provides embedded font resources for the emulator.
package fonts

import (
	"embed"
	"fmt"
	"log"
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

// DiagnoseEmbedding prints diagnostic info about embedded fonts
func DiagnoseEmbedding() {
	entries, err := fs.ReadDir(".")
	if err != nil {
		log.Printf("[fonts] ERROR reading embed directory: %v", err)
		return
	}
	log.Printf("[fonts] Found %d entries in embedded filesystem:", len(entries))
	for _, entry := range entries {
		info, _ := entry.Info()
		if info != nil {
			log.Printf("[fonts]   - %s (%d bytes)", entry.Name(), info.Size())
		} else {
			log.Printf("[fonts]   - %s (size unknown)", entry.Name())
		}
	}
}
