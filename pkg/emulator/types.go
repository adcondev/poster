package emulator

import (
	"image"
)

// RenderResult contains the output of an emulation render
type RenderResult struct {
	Image  image.Image
	Width  int
	Height int
}

// TextStyle represents text formatting options for the emulator
type TextStyle struct {
	Bold      bool
	Underline int     // 0: none, 1: single, 2: double
	Inverse   bool    // White on black
	ScaleW    float64 // Width multiplier (1.0 - 8.0)
	ScaleH    float64 // Height multiplier (1.0 - 8.0)
}

// DefaultTextStyle returns a TextStyle with default values
func DefaultTextStyle() TextStyle {
	return TextStyle{
		Bold:      false,
		Underline: 0,
		Inverse:   false,
		ScaleW:    1.0,
		ScaleH:    1.0,
	}
}

// Config holds configuration for the emulator engine
type Config struct {
	PaperWidthPx int
	DPI          int
	FontAPath    string // Optional: custom Font A path
	FontBPath    string // Optional: custom Font B path
}

// DefaultConfig returns a default configuration for 80mm paper at 203 DPI
func DefaultConfig() Config {
	return Config{
		PaperWidthPx: PaperWidth80mm,
		DPI:          203,
	}
}

// Config58mm returns configuration for 58mm paper at 203 DPI
func Config58mm() Config {
	return Config{
		PaperWidthPx: PaperWidth58mm,
		DPI:          203,
	}
}
