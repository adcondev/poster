package service

import (
	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/graphics"
)

// PrinterActions defines the interface for printer operations used by handlers.
// This abstraction allows for testing handlers without requiring a real printer.
type PrinterActions interface {
	// Basic operations
	Initialize() error
	Write(data []byte) error

	// Text printing
	Print(text string) error
	PrintLine(text string) error
	FeedLines(lines byte) error

	// Text formatting
	FontA() error
	FontB() error
	EnableBold() error
	DisableBold() error

	// Alignment
	AlignLeft() error
	AlignCenter() error
	AlignRight() error
	SetAlignment(align string) error

	// Size
	SingleSize() error
	DoubleSize() error
	TripleSize() error
	QuadraSize() error
	PentaSize() error
	HexaSize() error
	HeptaSize() error
	OctaSize() error
	CustomSize(width, height byte) error

	// Underline
	NoDot() error
	OneDot() error
	TwoDot() error

	// Inverse
	InverseOn() error
	InverseOff() error

	// Paper control
	FullFeedAndCut(lines byte) error
	PartialFeedAndCut(lines byte) error

	// Graphics
	PrintBitmap(bitmap *graphics.MonochromeBitmap) error
	PrintQR(data string, opts *graphics.QROptions) error
	PrintBarcode(cfg graphics.BarcodeConfig, data []byte) error

	// Character encoding
	SetCodeTable(codeTable character.CodeTable) error
}

// Ensure service. Printer implements PrinterActions
// This is a compile-time check that will fail if the interface doesn't match
// Note:  Uncomment this after adding the import for service package
// var _ PrinterActions = (*service.Printer)(nil)
