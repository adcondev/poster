package emulator

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"sync"

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
	scaledFaces map[string]font.Face // Cache:  "fontName_scaleW_scaleH" -> face
	mu          sync.RWMutex
	useFallback bool
}

// NewFontManager creates a new FontManager instance
func NewFontManager() *FontManager {
	return &FontManager{
		fonts:       make(map[string]*ScaledFont),
		scaledFaces: make(map[string]font.Face),
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
		log.Printf("[FontManager] ERROR:  Failed to parse TTF for '%s': %v", name, err)
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
			GlyphWidth:  targetWidth,
			GlyphHeight: targetHeight,
			LineHeight:  targetHeight + 6,
			Ascent:      targetHeight * 0.8,
			Descent:     targetHeight * 0.2,
		},
	}

	// Pre-cache the 1x1 scale for this font
	cacheKey := fmt.Sprintf("%s_1.0_1.0", name)
	fm.mu.Lock()
	fm.scaledFaces[cacheKey] = bestFace
	fm.mu.Unlock()

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

// GetScaledMetrics returns font metrics adjusted for the given scale factors
func (fm *FontManager) GetScaledMetrics(name string, scaleW, scaleH float64) FontMetrics {
	base := fm.GetMetrics(name)
	return FontMetrics{
		GlyphWidth:  base.GlyphWidth * scaleW,
		GlyphHeight: base.GlyphHeight * scaleH,
		LineHeight:  base.LineHeight * scaleH,
		Ascent:      base.Ascent * scaleH,
		Descent:     base.Descent * scaleH,
	}
}

// GetOrCreateScaledFace returns a cached scaled font face, creating it if necessary
func (fm *FontManager) GetOrCreateScaledFace(fontName string, scaleW, scaleH float64) (font.Face, error) {
	// Generate cache key
	key := fmt.Sprintf("%s_%.1f_%.1f", fontName, scaleW, scaleH)

	// Check cache first
	fm.mu.RLock()
	if face, ok := fm.scaledFaces[key]; ok {
		fm.mu.RUnlock()
		return face, nil
	}
	fm.mu.RUnlock()

	// Get the base font
	sf, err := fm.GetFont(fontName)
	if err != nil {
		return nil, err
	}

	// Calculate target height for the scaled font
	// Use the larger scale factor to determine the font size
	scaleFactor := scaleH
	if scaleW > scaleH {
		scaleFactor = scaleW
	}
	targetHeight := sf.metrics.GlyphHeight * scaleFactor

	// Search for the best matching font size
	var bestFace font.Face
	minDiff := math.MaxFloat64

	// Extended range for larger scales (up to 8x means we need larger sizes)
	maxSize := 72.0 + (scaleFactor * 20.0)
	if maxSize > 200.0 {
		maxSize = 200.0
	}

	for size := 6.0; size <= maxSize; size += 0.5 {
		opts := &opentype.FaceOptions{
			Size:    size,
			DPI:     72.0,
			Hinting: font.HintingFull,
		}
		face, err := opentype.NewFace(sf.ttFont, opts)
		if err != nil {
			continue
		}

		metrics := face.Metrics()
		currentHeight := float64(metrics.Height) / 64.0

		diff := math.Abs(currentHeight - targetHeight)
		if diff < minDiff {
			minDiff = diff
			bestFace = face
		}

		// Early exit if we found a very close match
		if diff < 0.5 {
			break
		}
	}

	if bestFace == nil {
		return nil, fmt.Errorf("could not create scaled face for %s at %.1fx%.1f", fontName, scaleW, scaleH)
	}

	// Cache the result
	fm.mu.Lock()
	// Double check cache in case another goroutine populated it
	if face, ok := fm.scaledFaces[key]; ok {
		fm.mu.Unlock()
		if bestFace != nil {
			_ = bestFace.Close()
		}
		return face, nil
	}
	fm.scaledFaces[key] = bestFace
	fm.mu.Unlock()

	log.Printf("[FontManager] Created and cached scaled face:  %s at %.1fx%.1f", fontName, scaleW, scaleH)

	return bestFace, nil
}

// DrawChar draws a single character at the specified position (1x1 scale)
// Returns the advance width
func (fm *FontManager) DrawChar(dst draw.Image, fontName string, char rune, x, y int, col color.Color) float64 {
	sf, err := fm.GetFont(fontName)
	if err != nil || fm.useFallback {
		// Fallback:  draw using bitmap
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

// DrawCharScaled draws a character with the specified scaling factors
// Returns the advance width (scaled)
func (fm *FontManager) DrawCharScaled(dst draw.Image, fontName string, char rune, x, y int, scaleW, scaleH float64, col color.Color) float64 {
	metrics := fm.GetScaledMetrics(fontName, scaleW, scaleH)

	// If fallback mode, use bitmap fallback
	if fm.useFallback {
		fm.drawFallbackCharScaled(dst, char, x, y, int(metrics.GlyphWidth), int(metrics.GlyphHeight), col)
		return metrics.GlyphWidth
	}

	// Get or create the scaled font face
	face, err := fm.GetOrCreateScaledFace(fontName, scaleW, scaleH)
	if err != nil {
		// Fallback to bitmap rendering on error
		fm.drawFallbackCharScaled(dst, char, x, y, int(metrics.GlyphWidth), int(metrics.GlyphHeight), col)
		return metrics.GlyphWidth
	}

	// Draw using the scaled font face
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(string(char))

	return metrics.GlyphWidth
}

// drawFallbackCharScaled renders a bitmap character scaled to fit the given dimensions
func (fm *FontManager) drawFallbackCharScaled(dst draw.Image, char rune, x, y, w, h int, col color.Color) {
	pattern := getFallbackPattern(char)

	// The 5x7 bitmap pattern needs to be scaled to fit w x h
	pixelW := w / 5 // 5 columns in the pattern
	pixelH := h / 7 // 7 rows in the pattern
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
				for sy := 0; sy < pixelH; sy++ {
					for sx := 0; sx < pixelW; sx++ {
						drawX := x + px*pixelW + sx
						drawY := y - h + py*pixelH + sy
						dst.Set(drawX, drawY, col)
					}
				}
			}
		}
	}
}

// FIXME: ClearScaledFaceCache recreates the map but doesn't close the existing font.Face resources.

// ClearScaledFaceCache clears the cached scaled font faces
// Useful when resetting the engine or changing fonts
func (fm *FontManager) ClearScaledFaceCache() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	for _, face := range fm.scaledFaces {
		_ = face.Close()
	}
	fm.scaledFaces = make(map[string]font.Face)
}
