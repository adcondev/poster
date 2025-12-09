package emulator

import (
	"image/color"
	"strings"
)

// BasicRenderer handles non-text rendering operations
type BasicRenderer struct {
	canvas *DynamicCanvas
	fonts  *FontManager
	state  *PrinterState
	black  color.Color
}

// NewBasicRenderer creates a new BasicRenderer
func NewBasicRenderer(canvas *DynamicCanvas, fonts *FontManager, state *PrinterState) *BasicRenderer {
	return &BasicRenderer{
		canvas: canvas,
		fonts:  fonts,
		state:  state,
		black:  color.RGBA{A: 255},
	}
}

// RenderSeparator renders a separator line using the specified character pattern
func (br *BasicRenderer) RenderSeparator(char string, length int) {
	if len(char) == 0 {
		char = "-"
	}

	// Build separator string
	var sep strings.Builder
	for sep.Len() < length {
		sep.WriteString(char)
	}
	separator := sep.String()
	if len(separator) > length {
		separator = separator[:length]
	}

	// Save current alignment and set to center for separator
	oldAlign := br.state.Align
	br.state.Align = AlignCenter

	// Get text renderer to render the separator
	textRenderer := NewTextRenderer(br.canvas, br.fonts, br.state)
	textRenderer.RenderLine(separator)

	// Restore alignment
	br.state.Align = oldAlign
}

// RenderFeed advances the paper by the specified number of lines
func (br *BasicRenderer) RenderFeed(lines int) {
	if lines <= 0 {
		lines = 1
	}

	metrics := br.fonts.GetMetrics(br.state.FontName)
	br.state.Feed(lines, metrics)
	br.canvas.UpdateMaxY(br.state.CursorY)
}

// RenderCut renders a visual cut line (dashed line indicating paper cut position)
func (br *BasicRenderer) RenderCut(isPartial bool) {
	// Add some space before the cut line
	br.RenderFeed(1)

	y := int(br.state.CursorY)
	width := br.state.PaperWidthPx

	if isPartial {
		// Partial cut: dashed line with small gap in center
		dashLen := 8
		gapLen := 4
		thickness := 1

		// Draw dashed line across the page
		br.canvas.DrawDashedLine(0, y, width, dashLen, gapLen, thickness, br.black)

		// Draw small "teeth" marks on the edge to indicate partial cut
		for x := 0; x < width; x += 20 {
			br.canvas.DrawRect(x, y-2, 2, 5, br.black)
		}
	} else {
		// Full cut: solid line
		br.canvas.DrawLine(0, y, width-1, 2, br.black)
	}

	br.state.CursorY += 10 // Space after cut
	br.canvas.UpdateMaxY(br.state.CursorY)
}

// RenderHorizontalLine renders a solid horizontal line
func (br *BasicRenderer) RenderHorizontalLine(thickness int) {
	if thickness <= 0 {
		thickness = 1
	}

	y := int(br.state.CursorY)
	br.canvas.DrawLine(0, y, br.state.PaperWidthPx-1, thickness, br.black)

	br.state.CursorY += float64(thickness) + 4
	br.canvas.UpdateMaxY(br.state.CursorY)
}
