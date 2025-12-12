package emulator

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/adcondev/poster/pkg/constants"
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
	ttfData, err := emfonts.LoadFontData(filename)
	if err != nil {
		log.Printf("[FontManager] ERROR:  Failed to load font data for '%s': %v", name, err)
		fm.useFallback = true
		return fmt.Errorf("loading font data: %w", err)
	}

	f, err := opentype.Parse(ttfData)
	if err != nil {
		log.Printf("[FontManager] ERROR: Failed to parse TTF for '%s': %v", name, err)
		fm.useFallback = true
		return fmt.Errorf("parsing ttf: %w", err)
	}

	// Heuristic search for optimal font size
	var bestFace font.Face
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
			continue
		}

		// Use 'M' as reference for width in monospace
		advance, ok := face.GlyphAdvance('M')
		if !ok {
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
		}
	}

	// Accept if we found ANY face (we'll use target metrics for spacing)
	if bestFace == nil {
		log.Printf("[FontManager] ERROR: Could not create any font face for '%s'", name)
		fm.useFallback = true
		return fmt.Errorf("could not fit font %s to dimensions %.2fx%.2f", filename, targetWidth, targetHeight)
	}

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
			GlyphWidth:  constants.FontBWidth,
			GlyphHeight: constants.FontBHeight,
			LineHeight:  constants.FontBHeight + 4,
			Ascent:      constants.FontBHeight * 0.8,
			Descent:     constants.FontBHeight * 0.2,
		}
	}
	return FontMetrics{
		GlyphWidth:  constants.FontAWidth,
		GlyphHeight: constants.FontAHeight,
		LineHeight:  constants.FontAHeight + 6,
		Ascent:      constants.FontAHeight * 0.8,
		Descent:     constants.FontAHeight * 0.2,
	}
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
