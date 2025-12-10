package emulator

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"strings"

	"github.com/adcondev/poster/pkg/constants"
)

// RenderResult contains the output of an emulation render
type RenderResult struct {
	Image  image.Image
	Width  int
	Height int
}

// Engine is the main entry point for ESC/POS emulation
type Engine struct {
	config Config
	canvas *DynamicCanvas
	fonts  *FontManager
	state  *PrinterState

	// Renderers
	textRenderer  *TextRenderer
	basicRenderer *BasicRenderer

	// Logs and debug info
	debug bool
}

// NewEngine creates a new emulator engine with the given configuration
func NewEngine(config Config) (*Engine, error) {
	if config.PaperPxWidth <= 0 {
		return nil, fmt.Errorf("invalid paper width: %d (must be > 0)", config.PaperPxWidth)
	}
	if config.DPI <= 0 {
		config.DPI = 203 // Default sensible
	}

	// Create canvas
	canvas := NewDynamicCanvas(config.PaperPxWidth)

	// TODO: Remove later - for debugging font embedding issues
	// emfonts.DiagnoseEmbedding()

	// Create font manager and load fonts
	fonts := NewFontManager()

	// Try to load embedded fonts, fall back to bitmap if not available
	errA := fonts.LoadFont("A", "JetBrainsMono-ExtraBold.ttf", constants.FontAWidth, constants.FontAHeight)
	errB := fonts.LoadFont("B", "JetBrainsMono-Bold.ttf", constants.FontBWidth, constants.FontBHeight)
	if errA != nil {
		log.Printf("[Emulator] Warning: Font A failed to load:  %v (using fallback)", errA)
	}
	if errB != nil {
		log.Printf("[Emulator] Warning: Font B failed to load: %v (using fallback)", errB)
	}

	if fonts.UseFallback() {
		log.Printf("[Emulator] Using bitmap fallback renderer")
	} else {
		log.Printf("[Emulator] Using TrueType font renderer")
	}

	// Create printer state
	state := NewPrinterState(config.PaperPxWidth)
	state.DPI = config.DPI

	engine := &Engine{
		config: config,
		canvas: canvas,
		fonts:  fonts,
		state:  state,
	}

	// Create renderers
	engine.textRenderer = NewTextRenderer(canvas, fonts, state)
	engine.basicRenderer = NewBasicRenderer(canvas, fonts, state)

	// Aplicar margen superior inicial para evitar recorte
	engine.applyTopMargin()

	return engine, nil
}

// NewDefaultEngine creates an engine with default 80mm paper configuration
func NewDefaultEngine() (*Engine, error) {
	return NewEngine(DefaultConfig())
}

// New58mmEngine creates an engine for 58mm paper
func New58mmEngine() (*Engine, error) {
	return NewEngine(Config58mm())
}

// Reset resets the engine to initial state (like ESC @)
func (e *Engine) Reset() {
	e.applyTopMargin()
	e.canvas = NewDynamicCanvas(e.config.PaperPxWidth)
	e.state.Reset()
	e.textRenderer = NewTextRenderer(e.canvas, e.fonts, e.state)
	e.basicRenderer = NewBasicRenderer(e.canvas, e.fonts, e.state)
	if e.debug {
		log.Printf("[Emulator] Engine reset to initial state")
	}
}

// ============================================================================
// Text Methods
// ============================================================================

// Print renders text without line feed
func (e *Engine) Print(text string) {
	e.textRenderer.RenderText(text)
}

// PrintLine renders text with line feed
func (e *Engine) PrintLine(text string) {
	e.textRenderer.RenderLine(text)
}

// PrintWrapped renders text with automatic word wrapping
func (e *Engine) PrintWrapped(text string) {
	e.textRenderer.WrapText(text)
}

// NewLine advances to the next line
func (e *Engine) NewLine() {
	e.textRenderer.NewLine()
}

// ============================================================================
// Style Methods
// ============================================================================

// SetAlign sets text alignment
func (e *Engine) SetAlign(align string) {
	// Normalize alignment input
	align = strings.TrimSpace(align)
	align = strings.ToLower(align)

	// Set alignment based on input
	switch align {
	case "center":
		e.state.Align = constants.Center.String()
	case "right":
		e.state.Align = constants.Right.String()
	default:
		e.state.Align = constants.Left.String()
	}
	if e.debug {
		log.Printf("[Emulator] Alignment set to: %s", e.state.Align)
	}
}

// AlignLeft sets left alignment
func (e *Engine) AlignLeft() {
	e.state.Align = constants.Left.String()
}

// AlignCenter sets center alignment
func (e *Engine) AlignCenter() {
	e.state.Align = constants.Center.String()
}

// AlignRight sets right alignment
func (e *Engine) AlignRight() {
	e.state.Align = constants.Right.String()
}

// SetBold enables or disables bold text
func (e *Engine) SetBold(enabled bool) {
	e.state.IsBold = enabled
	if e.debug {
		log.Printf("[Emulator] Bold mode set to: %v", enabled)
	}
}

// SetUnderline sets underline mode (0=off, 1=single, 2=double)
func (e *Engine) SetUnderline(mode int) {
	if mode < 0 {
		mode = 0
	}
	if mode > 2 {
		mode = 2
	}
	e.state.IsUnderline = mode
	if e.debug {
		log.Printf("[Emulator] Underline mode set to: %d", mode)
	}
}

// SetInverse enables or disables inverse mode (white on black)
func (e *Engine) SetInverse(enabled bool) {
	e.state.IsInverse = enabled
	if e.debug {
		log.Printf("[Emulator] Inverse mode set to: %v", enabled)
	}
}

// SetFont sets the current font ("A" or "B")
func (e *Engine) SetFont(name string) {
	// Normalize alignment input
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)

	// Set font based on input
	if name == "b" {
		e.state.FontName = "B"
	} else {
		e.state.FontName = "A"
	}
	if e.debug {
		log.Printf("[Emulator] Font set to: %s", e.state.FontName)
	}
}

// SetSize sets character size multipliers (1-8 for both width and height)
func (e *Engine) SetSize(width, height int) {
	e.state.SetSize(float64(width), float64(height))
	if e.debug {
		log.Printf("[Emulator] Character size set to: width=%d, height=%d", width, height)
	}
}

// applyTopMargin ensures there's enough space at the top for scaled text
func (e *Engine) applyTopMargin() {
	// Añadir espacio equivalente a una línea con la escala máxima típica (2x)
	e.state.CursorY = float64(constants.DefaultLineSpacing) + float64(constants.FontAHeight)
}

// ============================================================================
// Basic Operations
// ============================================================================

// Feed advances paper by specified number of lines
func (e *Engine) Feed(lines int) {
	e.basicRenderer.RenderFeed(lines)
}

// Separator renders a separator line
func (e *Engine) Separator(char string, length int) {
	e.basicRenderer.RenderSeparator(char, length)
}

// Cut renders a cut line (visual indicator)
func (e *Engine) Cut(partial bool) {
	e.basicRenderer.RenderCut(partial)
	if e.debug {
		log.Printf("[Emulator] Cut rendered: partial=%v", partial)
	}
}

// HorizontalLine renders a solid horizontal line
func (e *Engine) HorizontalLine(thickness int) {
	e.basicRenderer.RenderHorizontalLine(thickness)
	if e.debug {
		log.Printf("[Emulator] Horizontal line rendered: thickness=%d", thickness)
	}
}

// ============================================================================
// Output Methods
// ============================================================================

// Render returns the final cropped image
func (e *Engine) Render() image.Image {
	return e.canvas.Crop()
}

// RenderWithInfo returns the render result with dimensions
func (e *Engine) RenderWithInfo() RenderResult {
	img := e.canvas.Crop()
	bounds := img.Bounds()
	if e.debug {
		log.Printf("[Emulator] Rendered image size: %dx%d", bounds.Dx(), bounds.Dy())
	}
	return RenderResult{
		Image:  img,
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
	}
}

// WritePNG writes the rendered image as PNG to the writer
func (e *Engine) WritePNG(w io.Writer) error {
	img := e.Render()
	if e.debug {
		log.Printf("[Emulator] Writing PNG output (%dx%d)", img.Bounds().Dx(), img.Bounds().Dy())
	}
	return png.Encode(w, img)
}

// ============================================================================
// State Access
// ============================================================================

// State returns the current printer state (read-only access)
func (e *Engine) State() PrinterState {
	return *e.state
}

// CharsPerLine returns the number of characters that fit on a line with current font
func (e *Engine) CharsPerLine() int {
	metrics := e.fonts.GetMetrics(e.state.FontName)
	chars := e.state.CharsPerLine(metrics.GlyphWidth)
	if e.debug {
		log.Printf("[Emulator] CharsPerLine: %d (Font: %s, CharWidth: %.2f)", chars, e.state.FontName, metrics.GlyphWidth)
	}
	return chars
}

// ===========================================================================
// Debugging
// ===========================================================================

// EnableDebug enables debug logging
func (e *Engine) EnableDebug() {
	e.debug = true
}
