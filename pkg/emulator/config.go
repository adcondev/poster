package emulator

import (
	"github.com/adcondev/poster/pkg/constants"
)

// Config holds configuration for the emulator engine
type Config struct {
	PaperPxWidth int
	DPI          int
	FontAPath    string // Optional: custom Font A path
	FontBPath    string // Optional:  custom Font B path
	Debug        bool   // Enable debug logging

	// AutoAdjustCursorOnScale controls whether SetSize automatically adjusts
	// the cursor position when scaling up text to prevent overlap with
	// previous content.  Default:  true for ESC/POS-like behavior.
	// Set to false for explicit cursor control.
	AutoAdjustCursorOnScale bool
}

// DefaultConfig returns a default configuration for 80mm paper at 203 DPI
func DefaultConfig() Config {
	return Config{
		PaperPxWidth:            constants.PaperPxWidth80mm,
		DPI:                     constants.DefaultDPI,
		FontAPath:               "",
		FontBPath:               "",
		Debug:                   true,
		AutoAdjustCursorOnScale: true, // ESC/POS-like behavior by default
	}
}

// ConfigTest58mm returns configuration for 58mm paper at 203 DPI
func ConfigTest58mm() Config {
	return Config{
		PaperPxWidth:            constants.PaperPxWidth58mm,
		DPI:                     constants.DefaultDPI,
		FontAPath:               "",
		FontBPath:               "",
		Debug:                   true,
		AutoAdjustCursorOnScale: true,
	}
}
