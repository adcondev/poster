package service

import (
	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/profile"
)

// PrinterActions defines the interface for printer operations used by handlers.
// This abstraction allows for testing handlers without requiring a real printer.
type PrinterActions interface {
	// Basic operations
	Initialize() error
	Close() error
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

	// Device control
	Pulse(pin, onTime, offTime byte) error
	Beep(times, lapse byte) error

	// Line spacing
	SetLineSpacing(dots byte) error
	SelectDefaultLineSpacing() error

	// Graphics
	PrintBitmap(bitmap *graphics.MonochromeBitmap) error
	PrintQR(data string, opts *graphics.QrOptions) error
	PrintBarcode(cfg graphics.BarcodeConfig, data []byte) error

	// Character encoding
	SetCodeTable(codeTable character.CodeTable) error

	// Profile access - allows handlers to read/modify printer configuration
	GetProfile() *profile.Escpos
}

// Compile-time checks that Printer and MockPrinter implement PrinterActions
var (
	_ PrinterActions = (*Printer)(nil)
	_ PrinterActions = (*MockPrinter)(nil)
)
