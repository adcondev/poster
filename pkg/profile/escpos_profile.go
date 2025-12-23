package profile

import (
	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/commands/common"
	"github.com/adcondev/poster/pkg/graphics"
)

// Escpos defines all the physical characteristics and capabilities of a printer.
type Escpos struct {
	// Basic Information
	Model string // Same name used for printer connection

	// Physical Characteristics
	PaperWidth  float64 // in mm (e.g., 58mm, 80mm)
	PaperHeight float64 // in mm (0 for continuous roll)
	DPI         int     // Dots Per Inch (e.g., 203, 300)
	DotsPerLine int     // Dots per line (e.g., 384, 576)
	PrintWidth  int     // Print width in mm (e.g., 48, 42 in Font A)

	// Capabilities
	SupportsGraphics bool // Supports graphics (images)
	SupportsBarcode  bool // Supports native barcodes
	HasQR            bool // Supports native QR codes
	SupportsCutter   bool // Has automatic cutter
	SupportsDrawer   bool // Supports cash drawer

	// Maximum supported version
	QRMaxSize byte

	// Code table and encoding configuration
	CodeTable character.CodeTable

	// Advanced Configuration (optional)
	ImageThreshold int                 // Threshold for B/W conversion (0-255)
	Dithering      graphics.DitherMode // Default dithering type

	// Debug info
	DebugLog bool // Enable debug logs
}

// CreatePt210 creates a profile for the 58mm PT-58N thermal printer.
// TODO: Hardcoded Models: Factory functions for specific models might become hard to maintain. Consider configuration-based approach.
func CreatePt210() *Escpos {
	p := CreateProfile58mm()
	p.Model = "58mm PT-210"

	p.CodeTable = character.PC850
	p.QRMaxSize = 19 // Max supported QR version
	p.HasQR = true   // Supports native QR
	return p
}

// CreateGP58N creates a profile for the 58mm GP-58N thermal printer.
func CreateGP58N() *Escpos {
	p := CreateProfile58mm()
	p.Model = "58mm GP-58N"

	p.CodeTable = character.PC850
	return p
}

// CreateProfile58mm creates a profile for a common 58mm thermal printer.
func CreateProfile58mm() *Escpos {
	return &Escpos{
		Model: "Generic 58mm",

		PaperWidth:  58,
		DPI:         203,
		DotsPerLine: 384, // Typical for 58mm at 203 DPI
		PrintWidth:  48,  // Effective print width

		SupportsGraphics: true,
		SupportsBarcode:  true,
		HasQR:            false, // Many cheap printers do not support native QR
		SupportsCutter:   false,
		SupportsDrawer:   false,

		CodeTable: character.PC850,
	}
}

// CreateECPM80250 creates a profile for the 80mm EC-PM-80250 thermal printer.
func CreateECPM80250() *Escpos {
	p := CreateProfile80mm()
	p.Model = "80mm EC-PM-80250"
	return p
}

// CreateProfile80mm creates a profile for a common 80mm thermal printer.
func CreateProfile80mm() *Escpos {
	return &Escpos{
		Model: "Generic 80mm",

		PaperWidth:  80,
		DPI:         203,
		DotsPerLine: common.Dpl80mm203dpi, // Typical for 80mm (72mm) at 203 DPI

		SupportsGraphics: true,
		SupportsBarcode:  true,
		HasQR:            true, // 80mm printers usually have more features
		SupportsCutter:   true,
		SupportsDrawer:   true,

		// More character sets
		CodeTable: character.PC850, // CP850

		ImageThreshold: 128,
	}
}
