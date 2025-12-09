package emulator

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

var (
	colorWhite = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	colorBlack = color.RGBA{A: 255}
)

// DynamicCanvas provides an auto-growing canvas for receipt rendering
type DynamicCanvas struct {
	img    *image.RGBA
	width  int
	height int
	maxY   float64
}

// NewDynamicCanvas creates a new canvas with the specified width
func NewDynamicCanvas(width int) *DynamicCanvas {
	dc := &DynamicCanvas{
		width:  width,
		height: InitialHeight,
		maxY:   0,
	}
	dc.img = image.NewRGBA(image.Rect(0, 0, width, InitialHeight))
	dc.Clear()
	return dc
}

// Clear fills the canvas with white (paper color)
func (dc *DynamicCanvas) Clear() {
	white := colorWhite
	draw.Draw(dc.img, dc.img.Bounds(), &image.Uniform{C: white}, image.Point{}, draw.Src)
}

// EnsureHeight grows the canvas if needed to accommodate the required Y position
func (dc *DynamicCanvas) EnsureHeight(requiredY float64) {
	if requiredY < float64(dc.height) {
		return
	}

	// Calculate new height
	newHeight := int(float64(dc.height) * GrowthFactor)
	if newHeight < int(requiredY)+100 {
		newHeight = int(requiredY) + 500
	}
	if newHeight > MaxCanvasHeight {
		newHeight = MaxCanvasHeight
	}

	// Create new larger image
	newImg := image.NewRGBA(image.Rect(0, 0, dc.width, newHeight))

	// Fill with white
	white := colorWhite
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: white}, image.Point{}, draw.Src)

	// Copy old content
	draw.Draw(newImg, dc.img.Bounds(), dc.img, image.Point{}, draw.Src)

	dc.img = newImg
	dc.height = newHeight
}

// UpdateMaxY tracks the maximum Y position used
func (dc *DynamicCanvas) UpdateMaxY(y float64) {
	if y > dc.maxY {
		dc.maxY = y
	}
}

// Image returns the underlying RGBA image for drawing
func (dc *DynamicCanvas) Image() *image.RGBA {
	return dc.img
}

// Width returns the canvas width
func (dc *DynamicCanvas) Width() int {
	return dc.width
}

// Crop returns the final image cropped to the actual content height
func (dc *DynamicCanvas) Crop() image.Image {
	if dc.maxY == 0 {
		// Return minimal image if nothing was printed
		return image.NewRGBA(image.Rect(0, 0, dc.width, 1))
	}

	finalHeight := int(math.Ceil(dc.maxY)) + BottomPadding

	// Create final image of exact size
	rect := image.Rect(0, 0, dc.width, finalHeight)
	dst := image.NewRGBA(rect)

	draw.Draw(dst, rect, dc.img, image.Point{}, draw.Src)
	return dst
}

// DrawRect draws a filled rectangle
func (dc *DynamicCanvas) DrawRect(x, y, w, h int, col color.Color) {
	if w <= 0 || h <= 0 {
		return // Nada que dibujar
	}
	// Clamp to canvas bounds
	if x < 0 {
		w += x
		x = 0
	}
	if x+w > dc.width {
		w = dc.width - x
	}
	dc.EnsureHeight(float64(y + h))
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			dc.img.Set(x+dx, y+dy, col)
		}
	}
	dc.UpdateMaxY(float64(y + h))
}

// DrawLine draws a horizontal line
func (dc *DynamicCanvas) DrawLine(x1, y, x2 int, thickness int, col color.Color) {
	dc.EnsureHeight(float64(y + thickness))
	for t := 0; t < thickness; t++ {
		for x := x1; x <= x2; x++ {
			dc.img.Set(x, y+t, col)
		}
	}
	dc.UpdateMaxY(float64(y + thickness))
}

// DrawDashedLine draws a dashed horizontal line
func (dc *DynamicCanvas) DrawDashedLine(x1, y, x2 int, dashLen, gapLen, thickness int, col color.Color) {
	dc.EnsureHeight(float64(y + thickness))
	x := x1
	drawing := true
	segmentLen := dashLen

	for x <= x2 {
		if drawing {
			endX := x + segmentLen
			if endX > x2 {
				endX = x2
			}
			for t := 0; t < thickness; t++ {
				for px := x; px <= endX; px++ {
					dc.img.Set(px, y+t, col)
				}
			}
		}
		x += segmentLen
		drawing = !drawing
		if drawing {
			segmentLen = dashLen
		} else {
			segmentLen = gapLen
		}
	}
	dc.UpdateMaxY(float64(y + thickness))
}

// Set sets a single pixel
func (dc *DynamicCanvas) Set(x, y int, col color.Color) {
	if x >= 0 && x < dc.width && y >= 0 {
		dc.EnsureHeight(float64(y + 1))
		dc.img.Set(x, y, col)
		dc.UpdateMaxY(float64(y + 1))
	}
}
