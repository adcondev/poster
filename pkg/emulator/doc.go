/*
Package emulator provides a high-fidelity thermal printer emulation engine for Go.

It generates images (PNG) that accurately represent how a receipt would look
when printed on standard ESC/POS thermal printers (58mm and 80mm).

# Core Features

  - Dynamic Canvas:  Auto-growing canvas that adapts to content length
  - TrueType Font Rendering: Scales fonts to match thermal printer dot pitches (12x24 for Font A, 9x17 for Font B)
  - ESC/POS Styling: Bold, Underline, Inverse, Double Width/Height (1x-8x), and Justification
  - Image Embedding: Embed images with optional thermal preview simulation
  - Bitmap Fallback: Basic 5x7 bitmap font when TrueType fonts are unavailable

# Basic Usage

	package main

	import (
		"os"
		"github.com/adcondev/poster/pkg/emulator"
	)

	func main() {
		// Create engine (80mm default or 58mm)
		eng, _ := emulator.NewDefaultEngine()

		// Text with styling
		eng.AlignCenter()
		eng.SetBold(true)
		eng.SetSize(2, 2)
		eng.PrintLine("STORE NAME")
		eng.SetSize(1, 1)
		eng.SetBold(false)

		// Content
		eng.AlignLeft()
		eng.PrintLine("Item 1 ................ . $10.00")
		eng. Separator("-", 48)

		// Images (normal or thermal preview)
		img, _ := loadImage("logo.png")
		eng.PrintImage(img)
		eng.PrintImageThermalPreview(img, 200) // B&W dithered

		// Output
		eng.Cut(true)
		f, _ := os. Create("receipt.png")
		defer f.Close()
		eng.WritePNG(f)
	}

# Configuration

	// 80mm paper (576px width)
	eng, _ := emulator.NewDefaultEngine()

	// 58mm paper (384px width)
	eng, _ := emulator.New58mmEngine()

	// Custom configuration
	config := emulator.Config{
		PaperPxWidth:            576,
		DPI:                     203,
		AutoAdjustCursorOnScale: true, // Auto-adjust cursor when scaling up
	}
	eng, _ := emulator.NewEngine(config)

# Image Embedding

The emulator supports two image preview modes:

Normal Preview (default): Resizes and composites over white, preserving colors.
Ideal for digital receipts.

	eng.PrintImage(img)

Thermal Preview: Processes through the same pipeline as physical printing
(grayscale → dithering → monochrome). Shows exactly how it will print.

	eng. PrintImageThermalPreview(img, 256)

	// Or with full control:
	opts := emulator.DefaultImageOptions()
	opts.SimulateThermal = true
	opts. Dithering = graphics.Atkinson
	opts. Threshold = 128
	eng.PrintImageWithOptions(img, opts)

# Text Scaling

Text can be scaled 1x-8x in both dimensions.  By default, the cursor auto-adjusts
when scaling up to prevent overlap with previous content.

	eng. SetSize(2, 2) // Double size, cursor adjusts automatically
	eng.PrintLine("BIG TEXT")
	eng.SetSize(1, 1) // Back to normal

Disable auto-adjustment for manual cursor control:

	config := emulator.DefaultConfig()
	config.AutoAdjustCursorOnScale = false
	eng, _ := emulator.NewEngine(config)
*/
package emulator
