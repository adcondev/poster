/*
Package emulator provides a high-fidelity thermal printer emulation engine for Go.

It allows generating images (PNG) that accurately represent how a receipt would look
when printed on standard ESC/POS thermal printers (58mm and 80mm).

# Core Features

  - Dynamic Canvas: Auto-growing canvas that adapts to content length.
  - Font Calibration: Heuristic engine that scales TrueType fonts to match exact physical dot pitches (e.g., 12x24 dots for Font A).
  - ESC/POS Styling: Supports Bold, Underline, Inverse, Double Width/Height, and Justification.
  - Fidelity: Includes a bitmap fallback renderer to mimic the dot-matrix look of thermal heads when scaling or when fonts are missing.

# Basic Usage

To use the emulator, create an engine instance (usually with default 80mm config),
issue print commands, and finally render the output to an image.

	package main

	import (
		"os"
		"github.com/adcondev/poster/pkg/emulator"
	)

	func main() {
		// 1. Create a default 80mm engine
		eng, _ := emulator.NewDefaultEngine()

		// 2. Issue commands (simulating printer instructions)
		eng.SetAlign("center")
		eng.SetBold(true)
		eng.PrintLine("STORE NAME")
		eng.SetBold(false)
		eng.PrintLine("--------------------------------")
		eng.SetAlign("left")
		eng.PrintLine("Item 1 ................. $10.00")
		eng.PrintLine("Item 2 ................. $25.50")
		eng.Feed(2)
		eng.Cut(true)

		// 3. Render to PNG
		f, _ := os.Create("receipt.png")
		defer f.Close()
		eng.WritePNG(f)
	}

# Configuration

The engine can be customized for different paper widths (58mm/80mm) and DPI settings via the Config struct.
Embedded fonts (JetBrains Mono recommended) are used to simulate Font A and Font B, ensuring consistent character columns (48 cols for 80mm Font A).

# Architecture

The package is built around an Engine struct that maintains the virtual PrinterState (cursor position, active styles).
Rendering is delegated to specialized sub-renderers that draw onto a DynamicCanvas. This canvas manages memory efficiently, growing only as needed and cropping the final image to the exact content height.
*/
package emulator
