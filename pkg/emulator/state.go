package emulator

import (
	"log"
)

// PrinterState maintains the virtual printer state during emulation
type PrinterState struct {
	// Cursor position
	CursorX float64
	CursorY float64

	// Paper configuration
	PaperWidthPx int
	DPI          int

	// Text properties
	FontName    string // "A" or "B"
	IsBold      bool
	IsUnderline int     // 0: none, 1: single, 2: double
	IsInverse   bool    // White on black
	ScaleW      float64 // Width multiplier (1.0 - 8.0)
	ScaleH      float64 // Height multiplier (1.0 - 8.0)

	// Alignment
	Align string // "left", "center", "right"

	// Line spacing
	LineSpacing float64
}

// NewPrinterState creates a new PrinterState with default values
func NewPrinterState(paperWidthPx int) *PrinterState {
	return &PrinterState{
		CursorX:      0,
		CursorY:      float64(DefaultLineSpacing), // Start with some margin
		PaperWidthPx: paperWidthPx,
		DPI:          203,
		FontName:     "A",
		IsBold:       false,
		IsUnderline:  0,
		IsInverse:    false,
		ScaleW:       1.0,
		ScaleH:       1.0,
		Align:        AlignLeft,
		LineSpacing:  float64(DefaultLineSpacing),
	}
}

// Reset resets the printer state to defaults (like ESC @)
func (s *PrinterState) Reset() {
	s.CursorX = 0
	s.CursorY = float64(DefaultLineSpacing)
	s.FontName = "A"
	s.IsBold = false
	s.IsUnderline = 0
	s.IsInverse = false
	s.ScaleW = 1.0
	s.ScaleH = 1.0
	s.Align = AlignLeft
	s.LineSpacing = float64(DefaultLineSpacing)
}

// HasScaling checks if current state has scaling applied
func (s *PrinterState) HasScaling() bool {
	return s.ScaleW != 1.0 || s.ScaleH != 1.0
}

// NewLine moves cursor to beginning of next line
func (s *PrinterState) NewLine(fontMetrics FontMetrics) {
	lineHeight := fontMetrics.LineHeight * s.ScaleH
	if lineHeight < s.LineSpacing {
		lineHeight = s.LineSpacing
	}
	s.CursorY += lineHeight
	s.CursorX = 0
}

// Feed advances paper by the specified number of lines
func (s *PrinterState) Feed(lines int, fontMetrics FontMetrics) {
	lineHeight := fontMetrics.LineHeight * s.ScaleH
	if lineHeight < s.LineSpacing {
		lineHeight = s.LineSpacing
	}
	s.CursorY += lineHeight * float64(lines)
	s.CursorX = 0
}

// SetSize sets the character size multipliers
func (s *PrinterState) SetSize(w, h float64) {
	if w < 1.0 {
		log.Printf("Warning: Character width scale %f too small, setting to 1.0", w)
		w = 1.0
	}
	if w > 8.0 {
		log.Printf("Warning: Character width scale %f too large, setting to 8.0", w)
		w = 8.0
	}
	if h < 1.0 {
		log.Printf("Warning: Character height scale %f too small, setting to 1.0", h)
		h = 1.0
	}
	if h > 8.0 {
		log.Printf("Warning: Character height scale %f too large, setting to 8.0", h)
		h = 8.0
	}
	s.ScaleW = w
	s.ScaleH = h
}

// GetCharWidth returns the current character width considering scale
func (s *PrinterState) GetCharWidth(baseWidth float64) float64 {
	return baseWidth * s.ScaleW
}

// GetCharHeight returns the current character height considering scale
func (s *PrinterState) GetCharHeight(baseHeight float64) float64 {
	return baseHeight * s.ScaleH
}

// CharsPerLine calculates how many characters fit on one line
func (s *PrinterState) CharsPerLine(charWidth float64) int {
	return int(float64(s.PaperWidthPx) / (charWidth * s.ScaleW))
}
