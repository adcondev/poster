package emulator

import (
	"image/color"
	"strings"

	"github.com/adcondev/poster/pkg/constants"
)

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

// TextRenderer handles text rendering for the emulator
type TextRenderer struct {
	canvas *DynamicCanvas
	fonts  *FontManager
	state  *PrinterState
	black  color.Color
	white  color.Color
	// TODO: Si true, interpreta secuencias ESC/POS en el texto
	// parseEscPos bool
}

// NewTextRenderer creates a new TextRenderer
func NewTextRenderer(canvas *DynamicCanvas, fonts *FontManager, state *PrinterState) *TextRenderer {
	return &TextRenderer{
		canvas: canvas,
		fonts:  fonts,
		state:  state,
		black:  colorBlack,
		white:  colorWhite,
	}
}

// RenderText renders a string of text at the current cursor position
func (tr *TextRenderer) RenderText(text string) {
	if len(text) == 0 {
		return
	}
	// Filter control characters except common ones
	text = strings.Map(func(r rune) rune {
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			return -1 // Remove character
		}
		return r
	}, text)

	metrics := tr.fonts.GetMetrics(tr.state.FontName)
	charWidth := metrics.GlyphWidth * tr.state.ScaleW
	charHeight := metrics.GlyphHeight * tr.state.ScaleH

	// Calculate text width for alignment
	textWidth := float64(len(text)) * charWidth

	// Determine starting X position based on alignment
	startX := tr.calculateAlignedX(textWidth)

	// Ensure canvas has enough height
	requiredY := tr.state.CursorY + charHeight
	tr.canvas.EnsureHeight(requiredY)

	// Render each character
	x := startX
	for _, char := range text {
		tr.renderChar(char, x, tr.state.CursorY, charWidth, charHeight)
		x += charWidth
	}

	// Update cursor position
	tr.state.CursorX = x
	tr.canvas.UpdateMaxY(tr.state.CursorY + charHeight)
}

// RenderLine renders text and moves to next line
func (tr *TextRenderer) RenderLine(text string) {
	tr.RenderText(text)
	tr.NewLine()
}

// NewLine moves to the beginning of the next line
func (tr *TextRenderer) NewLine() {
	metrics := tr.fonts.GetMetrics(tr.state.FontName)
	tr.state.NewLine(metrics)
}

// Feed advances paper by specified number of lines
func (tr *TextRenderer) Feed(lines int) {
	metrics := tr.fonts.GetMetrics(tr.state.FontName)
	tr.state.Feed(lines, metrics)
	tr.canvas.UpdateMaxY(tr.state.CursorY)
}

// calculateAlignedX calculates the X starting position based on alignment
func (tr *TextRenderer) calculateAlignedX(textWidth float64) float64 {
	paperWidth := float64(tr.state.PaperPxWidth)

	switch tr.state.Align {
	case constants.Center.String():
		return (paperWidth - textWidth) / 2
	case constants.Right.String():
		return paperWidth - textWidth
	default: // AlignLeft
		return 0
	}
}

// renderChar renders a single character with current style
func (tr *TextRenderer) renderChar(char rune, x, y, width, height float64) {
	// Handle inverse mode (white on black)
	if tr.state.IsInverse {
		// Draw black background
		tr.canvas.DrawRect(
			int(x), int(y-height),
			int(width)+1, int(height)+1,
			tr.black,
		)
		// Draw white character
		tr.drawScaledChar(char, x, y, tr.white)
	} else {
		// Normal:  black on white
		tr.drawScaledChar(char, x, y, tr.black)
	}

	// Handle underline
	if tr.state.IsUnderline > 0 {
		underlineY := int(y) + 2
		thickness := tr.state.IsUnderline
		tr.canvas.DrawLine(int(x), underlineY, int(x+width), thickness, tr.black)
	}

	// Handle bold (draw twice with offset)
	if tr.state.IsBold {
		tr.drawScaledChar(char, x+1, y, tr.black)
	}
}

// drawScaledChar draws a character with scaling applied.
// It uses the current state's scale factors and font metrics internally.
func (tr *TextRenderer) drawScaledChar(char rune, x, y float64, col color.Color) {
	baseMetrics := tr.fonts.GetMetrics(tr.state.FontName)
	baseW := int(baseMetrics.GlyphWidth)
	baseH := int(baseMetrics.GlyphHeight)

	if !tr.state.HasScaling() {
		// No scaling - direct draw using TrueType font
		tr.fonts.DrawChar(tr.canvas.Image(), tr.state.FontName, char, int(x), int(y), col)
		return
	}

	// With scaling - use bitmap fallback pattern for pixel-perfect scaling
	scaleX := int(tr.state.ScaleW)
	scaleY := int(tr.state.ScaleH)

	pattern := getFallbackPattern(char)

	pixelW := baseW / 6
	pixelH := baseH / 8
	if pixelW < 1 {
		pixelW = 1
	}
	if pixelH < 1 {
		pixelH = 1
	}

	for py := 0; py < 7; py++ {
		for px := 0; px < 5; px++ {
			if pattern[py]&(1<<(4-px)) != 0 {
				// Draw scaled pixel
				for sy := 0; sy < pixelH*scaleY; sy++ {
					for sx := 0; sx < pixelW*scaleX; sx++ {
						drawX := int(x) + px*pixelW*scaleX + sx
						drawY := int(y) - baseH*scaleY + py*pixelH*scaleY + sy
						tr.canvas.Set(drawX, drawY, col)
					}
				}
			}
		}
	}
}

// WrapText wraps text to fit within paper width and renders each line
func (tr *TextRenderer) WrapText(text string) {
	metrics := tr.fonts.GetMetrics(tr.state.FontName)
	charWidth := metrics.GlyphWidth * tr.state.ScaleW
	charsPerLine := tr.state.CharsPerLine(charWidth)

	if charsPerLine <= 0 {
		charsPerLine = 1
	}

	// Split into words
	words := strings.Fields(text)
	if len(words) == 0 {
		return
	}

	var currentLine strings.Builder
	currentLen := 0

	for _, word := range words {
		wordLen := len(word)

		// If word is longer than line, split it
		if wordLen > charsPerLine {
			// Flush current line first
			if currentLen > 0 {
				tr.RenderLine(currentLine.String())
				currentLine.Reset()
				currentLen = 0
			}
			// Split long word
			for i := 0; i < wordLen; i += charsPerLine {
				end := i + charsPerLine
				if end > wordLen {
					end = wordLen
				}
				tr.RenderLine(word[i:end])
			}
			continue
		}

		// Check if word fits on current line
		spaceNeeded := wordLen
		if currentLen > 0 {
			spaceNeeded++ // space before word
		}

		if currentLen+spaceNeeded > charsPerLine {
			// Word doesn't fit - render current line and start new one
			tr.RenderLine(currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
			currentLen = wordLen
		} else {
			// Word fits
			if currentLen > 0 {
				currentLine.WriteByte(' ')
				currentLen++
			}
			currentLine.WriteString(word)
			currentLen += wordLen
		}
	}

	// Render remaining text
	if currentLen > 0 {
		tr.RenderLine(currentLine.String())
	}
}
