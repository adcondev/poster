package emulator

import (
	"image/color"
	"strings"
	"unicode/utf8"

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

	// Get scaled metrics for current font and size
	metrics := tr.fonts.GetScaledMetrics(tr.state.FontName, tr.state.ScaleW, tr.state.ScaleH)
	charWidth := metrics.GlyphWidth
	charHeight := metrics.GlyphHeight

	// Count actual characters (runes), not bytes - important for UTF-8 text
	runeCount := utf8.RuneCountInString(text)

	// Calculate text width for alignment using rune count
	textWidth := float64(runeCount) * charWidth

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

// FIXME: The NewLine method now directly manipulates cursor position instead of delegating to PrinterState.NewLine(metrics).
// This duplicates line height calculation logic.
// Consider whether the PrinterState.NewLine method should also be updated to use GetScaledMetrics,
// or if this logic should remain centralized in one place.

// NewLine moves to the beginning of the next line
func (tr *TextRenderer) NewLine() {
	metrics := tr.fonts.GetScaledMetrics(tr.state.FontName, tr.state.ScaleW, tr.state.ScaleH)
	lineHeight := metrics.LineHeight
	if lineHeight < tr.state.LineSpacing {
		lineHeight = tr.state.LineSpacing
	}
	tr.state.CursorY += lineHeight
	tr.state.CursorX = 0
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
		// Draw character in white
		tr.drawChar(char, x, y, tr.white)
	} else {
		// Normal:  black on white
		tr.drawChar(char, x, y, tr.black)
	}

	// Handle underline
	if tr.state.IsUnderline > 0 {
		underlineY := int(y) + 2
		thickness := tr.state.IsUnderline
		tr.canvas.DrawLine(int(x), underlineY, int(x+width), thickness, tr.black)
	}

	// Handle bold (draw twice with offset for extra weight)
	if tr.state.IsBold {
		if tr.state.IsInverse {
			tr.drawChar(char, x+1, y, tr.white)
		} else {
			tr.drawChar(char, x+1, y, tr.black)
		}
	}
}

// drawChar draws a character using TrueType fonts (scaled or unscaled)
func (tr *TextRenderer) drawChar(char rune, x, y float64, col color.Color) {
	if !tr.state.HasScaling() {
		// No scaling - direct draw using base TrueType font
		tr.fonts.DrawChar(tr.canvas.Image(), tr.state.FontName, char, int(x), int(y), col)
		return
	}

	// With scaling - use cached scaled TrueType font
	tr.fonts.DrawCharScaled(
		tr.canvas.Image(),
		tr.state.FontName,
		char,
		int(x),
		int(y),
		tr.state.ScaleW,
		tr.state.ScaleH,
		col,
	)
}

// FIXME: The WrapText method now calculates charsPerLine directly instead of using state.CharsPerLine(charWidth).

// WrapText wraps text to fit within paper width and renders each line
func (tr *TextRenderer) WrapText(text string) {
	metrics := tr.fonts.GetScaledMetrics(tr.state.FontName, tr.state.ScaleW, tr.state.ScaleH)
	charWidth := metrics.GlyphWidth
	charsPerLine := int(float64(tr.state.PaperPxWidth) / charWidth)

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
		// Count runes, not bytes - important for UTF-8 text
		wordLen := utf8.RuneCountInString(word)

		// If word is longer than line, split it
		if wordLen > charsPerLine {
			// Flush current line first
			if currentLen > 0 {
				tr.RenderLine(currentLine.String())
				currentLine.Reset()
				currentLen = 0
			}
			// Split long word by runes (not bytes)
			runes := []rune(word)
			for i := 0; i < len(runes); i += charsPerLine {
				end := i + charsPerLine
				if end > len(runes) {
					end = len(runes)
				}
				tr.RenderLine(string(runes[i:end]))
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
