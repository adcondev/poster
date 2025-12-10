package emulator

import (
	"github.com/adcondev/poster/pkg/constants"
)

// Config holds configuration for the emulator engine
type Config struct {
	PaperPxWidth int
	DPI          int
	FontAPath    string // Optional: custom Font A path
	FontBPath    string // Optional: custom Font B path
	Debug        bool   // Enable debug logging
}

// DefaultConfig returns a default configuration for 80mm paper at 203 DPI
func DefaultConfig() Config {
	return Config{
		PaperPxWidth: constants.PaperPxWidth80mm,
		DPI:          203,
	}
}

// Config58mm returns configuration for 58mm paper at 203 DPI
func Config58mm() Config {
	return Config{
		PaperPxWidth: constants.PaperPxWidth58mm,
		DPI:          203,
	}
}
