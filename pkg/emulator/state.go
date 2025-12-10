package emulator

import (
	"log"

	"github.com/adcondev/poster/pkg/constants"
)

// PrinterState maintains the virtual printer state during emulation
type PrinterState struct {
	// Cursor position
	CursorX float64
	CursorY float64

	// Paper configuration
	PaperPxWidth int
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
		CursorY:      float64(constants.DefaultLineSpacing),
		PaperPxWidth: paperWidthPx,
		DPI:          constants.DefaultDPI,
		FontName:     "A",
		IsBold:       false,
		IsUnderline:  constants.MinUnderline,
		IsInverse:    false,
		ScaleW:       constants.MinScale,
		ScaleH:       constants.MinScale,
		Align:        constants.Left.String(),
		LineSpacing:  float64(constants.DefaultLineSpacing),
	}
}

// Reset resets the printer state to defaults (like ESC @)
func (s *PrinterState) Reset() {
	s.CursorX = 0
	s.CursorY = float64(constants.DefaultLineSpacing)
	s.FontName = "A"
	s.IsBold = false
	s.IsUnderline = constants.MinUnderline
	s.IsInverse = false
	s.ScaleW = constants.MinScale
	s.ScaleH = constants.MinScale
	s.Align = constants.Left.String()
	s.LineSpacing = float64(constants.DefaultLineSpacing)
}

// HasScaling checks if current state has scaling applied
func (s *PrinterState) HasScaling() bool {
	return s.ScaleW != constants.MinScale || s.ScaleH != constants.MinScale
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
	if lines <= 0 {
		lines = 1
	}
	lineHeight := fontMetrics.LineHeight * s.ScaleH
	if lineHeight < s.LineSpacing {
		lineHeight = s.LineSpacing
	}
	s.CursorY += lineHeight * float64(lines)
	s.CursorX = 0
}

// SetSize sets the character size multipliers with clamping
func (s *PrinterState) SetSize(w, h float64) {
	s.ScaleW = clampScale(w, "width")
	s.ScaleH = clampScale(h, "height")
}

// clampScale clamps a scale value to valid range and logs warnings
func clampScale(value float64, name string) float64 {
	if value < constants.MinScale {
		log.Printf("Warning: Character %s scale %.2f too small, clamping to %.1f", name, value, constants.MinScale)
		return constants.MinScale
	}
	if value > constants.MaxScale {
		log.Printf("Warning: Character %s scale %.2f too large, clamping to %.1f", name, value, constants.MaxScale)
		return constants.MaxScale
	}
	return value
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
	if charWidth <= 0 || s.ScaleW <= 0 {
		return 0
	}
	return int(float64(s.PaperPxWidth) / (charWidth * s.ScaleW))
}
