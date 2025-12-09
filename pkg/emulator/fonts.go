package emulator

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"unicode"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	emfonts "github.com/adcondev/poster/pkg/emulator/fonts"
)

// FontMetrics contains calculated font dimensions
type FontMetrics struct {
	GlyphWidth  float64
	GlyphHeight float64
	LineHeight  float64
	Ascent      float64
	Descent     float64
}

// ScaledFont holds a font face with its calculated metrics
type ScaledFont struct {
	face    font.Face
	metrics FontMetrics
	ttFont  *opentype.Font // Keep reference for scaling
}

// FontManager handles font loading and scaling for thermal printer emulation
type FontManager struct {
	fonts       map[string]*ScaledFont
	useFallback bool
}

// NewFontManager creates a new FontManager instance
func NewFontManager() *FontManager {
	return &FontManager{
		fonts:       make(map[string]*ScaledFont),
		useFallback: false,
	}
}

// LoadFont loads and calibrates a font to match target pixel dimensions
func (fm *FontManager) LoadFont(name, filename string, targetWidth, targetHeight float64) error {
	log.Printf("[FontManager] Attempting to load font '%s' from file '%s' (target:  %.0fx%.0f)",
		name, filename, targetWidth, targetHeight)

	ttfData, err := emfonts.LoadFontData(filename)
	if err != nil {
		log.Printf("[FontManager] ERROR:  Failed to load font data for '%s': %v", name, err)
		fm.useFallback = true
		return fmt.Errorf("loading font data: %w", err)
	}
	log.Printf("[FontManager] Loaded %d bytes for font '%s'", len(ttfData), name)

	f, err := opentype.Parse(ttfData)
	if err != nil {
		log.Printf("[FontManager] ERROR: Failed to parse TTF for '%s': %v", name, err)
		fm.useFallback = true
		return fmt.Errorf("parsing ttf: %w", err)
	}
	log.Printf("[FontManager] Successfully parsed TTF for '%s'", name)

	// Heuristic search for optimal font size
	var bestFace font.Face
	var bestSize float64
	minDiff := math.MaxFloat64

	// Search range for thermal printer fonts (6pt to 72pt)
	for size := 6.0; size <= 72.0; size += 0.5 {
		opts := &opentype.FaceOptions{
			Size:    size,
			DPI:     72.0,
			Hinting: font.HintingFull,
		}
		face, err := opentype.NewFace(f, opts)
		if err != nil {
			log.Printf("Debug: Failed to create face at size %.1f: %v", size, err)
			continue
		}

		// Use 'M' as reference for width in monospace
		advance, ok := face.GlyphAdvance('M')
		if !ok {
			log.Printf("Debug: Failed to find 'M' glyph at size %.1f", size)
			continue
		}

		currentWidth := float64(advance) / 64.0
		metrics := face.Metrics()
		currentHeight := float64(metrics.Height) / 64.0

		// Calculate difference - prioritize width match
		widthDiff := math.Abs(currentWidth - targetWidth)
		heightDiff := math.Abs(currentHeight - targetHeight)
		totalDiff := widthDiff*2 + heightDiff // Weight width more

		if totalDiff < minDiff {
			minDiff = totalDiff
			bestFace = face
			bestSize = size
		}
	}

	// RELAXED:  Accept if we found ANY face (we'll use target metrics for spacing)
	if bestFace == nil {
		log.Printf("[FontManager] ERROR: Could not create any font face for '%s'", name)
		fm.useFallback = true
		return fmt.Errorf("could not fit font %s to dimensions %.2fx%.2f", filename, targetWidth, targetHeight)
	}

	// Get actual metrics from best face
	actualMetrics := bestFace.Metrics()
	advance, _ := bestFace.GlyphAdvance('M')
	actualWidth := float64(advance) / 64.0
	actualHeight := float64(actualMetrics.Height) / 64.0

	log.Printf("[FontManager] Best match for '%s': size=%.1fpt, actual=%.1fx%.1f, diff=%.2f",
		name, bestSize, actualWidth, actualHeight, minDiff)

	// Store with TARGET metrics for consistent spacing (important!)
	fm.fonts[name] = &ScaledFont{
		face:   bestFace,
		ttFont: f,
		metrics: FontMetrics{
			GlyphWidth:  targetWidth, // Use target for consistent grid
			GlyphHeight: targetHeight,
			LineHeight:  targetHeight + 6,
			Ascent:      targetHeight * 0.8,
			Descent:     targetHeight * 0.2,
		},
	}

	log.Printf("[FontManager] Successfully loaded font '%s' (using target metrics for spacing)", name)
	return nil
}

// GetFont retrieves a loaded font by name
func (fm *FontManager) GetFont(name string) (*ScaledFont, error) {
	if f, ok := fm.fonts[name]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("font %s not loaded", name)
}

// GetMetrics returns metrics for a font, using fallback if necessary
func (fm *FontManager) GetMetrics(name string) FontMetrics {
	if f, ok := fm.fonts[name]; ok {
		return f.metrics
	}
	// Return default metrics for fallback
	if name == "B" {
		return FontMetrics{
			GlyphWidth:  FontBWidth,
			GlyphHeight: FontBHeight,
			LineHeight:  FontBHeight + 4,
			Ascent:      FontBHeight * 0.8,
			Descent:     FontBHeight * 0.2,
		}
	}
	return FontMetrics{
		GlyphWidth:  FontAWidth,
		GlyphHeight: FontAHeight,
		LineHeight:  FontAHeight + 6,
		Ascent:      FontAHeight * 0.8,
		Descent:     FontAHeight * 0.2,
	}
}

// UseFallback returns true if fallback bitmap rendering should be used
func (fm *FontManager) UseFallback() bool {
	return fm.useFallback
}

// DrawChar draws a single character at the specified position
// Returns the advance width
func (fm *FontManager) DrawChar(dst draw.Image, fontName string, char rune, x, y int, col color.Color) float64 {
	sf, err := fm.GetFont(fontName)
	if err != nil || fm.useFallback {
		// Fallback: draw a simple rectangle placeholder
		metrics := fm.GetMetrics(fontName)
		fm.drawFallbackChar(dst, char, x, y, int(metrics.GlyphWidth), int(metrics.GlyphHeight), col)
		return metrics.GlyphWidth
	}

	// Draw using the font face
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(col),
		Face: sf.face,
		Dot:  point,
	}
	d.DrawString(string(char))

	return sf.metrics.GlyphWidth
}

// drawFallbackChar renders a basic bitmap character when fonts aren't available
func (fm *FontManager) drawFallbackChar(dst draw.Image, char rune, x, y, w, h int, col color.Color) {
	// Simple 5x7 bitmap font patterns for basic ASCII
	// This is a minimal fallback - real implementation would have full charset
	patterns := getFallbackPattern(char)

	scaleX := w / 6 // 5 pixels + 1 spacing
	scaleY := h / 8 // 7 pixels + 1 spacing
	if scaleX < 1 {
		scaleX = 1
	}
	if scaleY < 1 {
		scaleY = 1
	}

	for py := 0; py < 7; py++ {
		for px := 0; px < 5; px++ {
			if patterns[py]&(1<<(4-px)) != 0 {
				// Draw scaled pixel
				for sy := 0; sy < scaleY; sy++ {
					for sx := 0; sx < scaleX; sx++ {
						dst.Set(x+px*scaleX+sx, y-h+py*scaleY+sy, col)
					}
				}
			}
		}
	}
}

// getFallbackPattern returns a 5x7 bitmap pattern for a character
func getFallbackPattern(char rune) [7]byte {
	// Convert to uppercase first
	char = unicode.ToUpper(char)

	// Basic ASCII patterns (5 bits per row, 7 rows)
	patterns := map[rune][7]byte{
		' ': {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		'A': {0x0E, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11},
		'B': {0x1E, 0x11, 0x11, 0x1E, 0x11, 0x11, 0x1E},
		'C': {0x0E, 0x11, 0x10, 0x10, 0x10, 0x11, 0x0E},
		'D': {0x1E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x1E},
		'E': {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x1F},
		'F': {0x1F, 0x10, 0x10, 0x1E, 0x10, 0x10, 0x10},
		'G': {0x0E, 0x11, 0x10, 0x17, 0x11, 0x11, 0x0E},
		'H': {0x11, 0x11, 0x11, 0x1F, 0x11, 0x11, 0x11},
		'I': {0x0E, 0x04, 0x04, 0x04, 0x04, 0x04, 0x0E},
		'J': {0x1F, 0x04, 0x04, 0x04, 0x04, 0x14, 0x08},
		'K': {0x11, 0x12, 0x14, 0x18, 0x14, 0x12, 0x11},
		'L': {0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x1F},
		'M': {0x11, 0x1B, 0x15, 0x11, 0x11, 0x11, 0x11},
		'N': {0x11, 0x19, 0x15, 0x13, 0x11, 0x11, 0x11},
		'O': {0x0E, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E},
		'P': {0x1E, 0x11, 0x11, 0x1E, 0x10, 0x10, 0x10},
		'Q': {0x0E, 0x11, 0x11, 0x11, 0x15, 0x12, 0x0D},
		'R': {0x1E, 0x11, 0x11, 0x1E, 0x14, 0x12, 0x11},
		'S': {0x0E, 0x11, 0x10, 0x0E, 0x01, 0x11, 0x0E},
		'T': {0x1F, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04},
		'U': {0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x0E},
		'V': {0x11, 0x11, 0x11, 0x11, 0x0A, 0x0A, 0x04},
		'W': {0x11, 0x11, 0x11, 0x15, 0x15, 0x0A, 0x0A},
		'X': {0x11, 0x11, 0x0A, 0x04, 0x0A, 0x11, 0x11},
		'Y': {0x11, 0x11, 0x0A, 0x04, 0x04, 0x04, 0x04},
		'Z': {0x1F, 0x01, 0x02, 0x04, 0x08, 0x10, 0x1F},
		'0': {0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E},
		'1': {0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E},
		'2': {0x0E, 0x11, 0x01, 0x06, 0x08, 0x10, 0x1F},
		'3': {0x0E, 0x11, 0x01, 0x06, 0x01, 0x11, 0x0E},
		'4': {0x02, 0x06, 0x0A, 0x12, 0x1F, 0x02, 0x02},
		'5': {0x1F, 0x10, 0x1E, 0x01, 0x01, 0x11, 0x0E},
		'6': {0x06, 0x08, 0x10, 0x1E, 0x11, 0x11, 0x0E},
		'7': {0x1F, 0x01, 0x02, 0x04, 0x08, 0x08, 0x08},
		'8': {0x0E, 0x11, 0x11, 0x0E, 0x11, 0x11, 0x0E},
		'9': {0x0E, 0x11, 0x11, 0x1F, 0x01, 0x02, 0x0C},
		'-': {0x00, 0x00, 0x00, 0x0E, 0x00, 0x00, 0x00},
		'=': {0x00, 0x00, 0x0E, 0x00, 0x0E, 0x00, 0x00},
		'.': {0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06},
		',': {0x00, 0x00, 0x00, 0x00, 0x06, 0x02, 0x04},
		':': {0x00, 0x0C, 0x0C, 0x00, 0x0C, 0x0C, 0x00},
		'$': {0x04, 0x0F, 0x14, 0x0E, 0x05, 0x1E, 0x04},
		'#': {0x0A, 0x0A, 0x1F, 0x0A, 0x1F, 0x0A, 0x0A},
		'!': {0x04, 0x04, 0x04, 0x04, 0x04, 0x00, 0x04},
		'?': {0x0E, 0x11, 0x01, 0x02, 0x04, 0x00, 0x04},
		'/': {0x01, 0x01, 0x02, 0x04, 0x08, 0x10, 0x10},
		'*': {0x00, 0x04, 0x15, 0x0E, 0x15, 0x04, 0x00},
		'+': {0x00, 0x04, 0x04, 0x1F, 0x04, 0x04, 0x00},
		'(': {0x02, 0x04, 0x08, 0x08, 0x08, 0x04, 0x02},
		')': {0x08, 0x04, 0x02, 0x02, 0x02, 0x04, 0x08},
		// Extended Latin characters for Spanish/Portuguese support
		// N with tilde: Tilde on Row 0, Compressed N on Rows 1-6
		'Ñ': {0x0E, 0x11, 0x19, 0x15, 0x13, 0x11, 0x11},
		// Lowercase n with tilde
		'Á': {0x02, 0x04, 0x0E, 0x11, 0x1F, 0x11, 0x11},
		'É': {0x02, 0x04, 0x1F, 0x10, 0x1E, 0x10, 0x1F},
		'Í': {0x02, 0x04, 0x0E, 0x04, 0x04, 0x04, 0x0E},
		'Ó': {0x02, 0x04, 0x0E, 0x11, 0x11, 0x11, 0x0E},
		'Ú': {0x02, 0x04, 0x11, 0x11, 0x11, 0x11, 0x0E},
		'Ü': {0x0A, 0x00, 0x11, 0x11, 0x11, 0x11, 0x0E},
		'¿': {0x04, 0x00, 0x04, 0x08, 0x10, 0x11, 0x0E},
		'¡': {0x04, 0x00, 0x04, 0x04, 0x04, 0x04, 0x04},
	}

	// Check for match (already uppercase)
	if p, ok := patterns[char]; ok {
		return p
	}

	// Default: filled rectangle for unknown chars
	return [7]byte{0x1F, 0x11, 0x11, 0x11, 0x11, 0x11, 0x1F}
}
