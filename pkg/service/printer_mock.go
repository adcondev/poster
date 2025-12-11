package service

import (
	"fmt"

	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/graphics"
)

var _ PrinterActions = (*MockPrinter)(nil)

// MockPrinter implements PrinterActions for testing purposes.
// It records all calls made to it for verification in tests.
type MockPrinter struct {
	// Call tracking
	Calls []MockCall

	// Configuration for behavior
	FailOn       map[string]error // Map of method names to errors to return
	InitializeOK bool

	// State tracking
	CurrentAlignment string
	CurrentFont      string
	BoldEnabled      bool
	InverseEnabled   bool
	CurrentSize      string
	UnderlineMode    string

	// Output capture
	PrintedText    []string
	WrittenBytes   [][]byte
	PrintedQRs     []string
	PrintedBitmaps []*graphics.MonochromeBitmap
}

// MockCall represents a single call to the mock printer
type MockCall struct {
	Method string
	Args   []interface{}
}

// NewMockPrinter creates a new MockPrinter with default settings
func NewMockPrinter() *MockPrinter {
	return &MockPrinter{
		Calls:            []MockCall{},
		FailOn:           make(map[string]error),
		InitializeOK:     true,
		CurrentAlignment: "left",
		CurrentFont:      "A",
		CurrentSize:      "1x1",
		UnderlineMode:    "none",
	}
}

// record adds a call to the call log
func (m *MockPrinter) record(method string, args ...interface{}) {
	m.Calls = append(m.Calls, MockCall{Method: method, Args: args})
}

// checkError returns an error if the method is configured to fail
func (m *MockPrinter) checkError(method string) error {
	if err, ok := m.FailOn[method]; ok {
		return err
	}
	return nil
}

// LastCall returns the last call made, or nil if no calls
func (m *MockPrinter) LastCall() *MockCall {
	if len(m.Calls) == 0 {
		return nil
	}
	return &m.Calls[len(m.Calls)-1]
}

// CallCount returns the number of calls to a specific method
func (m *MockPrinter) CallCount(method string) int {
	count := 0
	for _, call := range m.Calls {
		if call.Method == method {
			count++
		}
	}
	return count
}

// Reset clears all recorded calls and state
func (m *MockPrinter) Reset() {
	m.Calls = []MockCall{}
	m.PrintedText = []string{}
	m.WrittenBytes = [][]byte{}
	m.PrintedQRs = []string{}
	m.PrintedBitmaps = nil
	m.CurrentAlignment = "left"
	m.CurrentFont = "A"
	m.BoldEnabled = false
	m.InverseEnabled = false
	m.CurrentSize = "1x1"
	m.UnderlineMode = "none"
}

// ============================================================================
// PrinterActions Implementation
// ============================================================================

// Initialize simulates printer initialization
func (m *MockPrinter) Initialize() error {
	m.record("Initialize")
	if !m.InitializeOK {
		return fmt.Errorf("initialize failed")
	}
	return m.checkError("Initialize")
}

// Write simulates writing raw bytes to the printer
func (m *MockPrinter) Write(data []byte) error {
	m.record("Write", data)
	m.WrittenBytes = append(m.WrittenBytes, data)
	return m.checkError("Write")
}

// Print simulates printing text
func (m *MockPrinter) Print(text string) error {
	m.record("Print", text)
	m.PrintedText = append(m.PrintedText, text)
	return m.checkError("Print")
}

// PrintLine simulates printing a line of text
func (m *MockPrinter) PrintLine(text string) error {
	m.record("PrintLine", text)
	m.PrintedText = append(m.PrintedText, text+"\n")
	return m.checkError("PrintLine")
}

// FeedLines simulates feeding paper by a number of lines
func (m *MockPrinter) FeedLines(lines byte) error {
	m.record("FeedLines", lines)
	return m.checkError("FeedLines")
}

// FontA sets the font to A
func (m *MockPrinter) FontA() error {
	m.record("FontA")
	m.CurrentFont = "A"
	return m.checkError("FontA")
}

// FontB sets the font to B
func (m *MockPrinter) FontB() error {
	m.record("FontB")
	m.CurrentFont = "B"
	return m.checkError("FontB")
}

// EnableBold enables bold text
func (m *MockPrinter) EnableBold() error {
	m.record("EnableBold")
	m.BoldEnabled = true
	return m.checkError("EnableBold")
}

// DisableBold disables bold text
func (m *MockPrinter) DisableBold() error {
	m.record("DisableBold")
	m.BoldEnabled = false
	return m.checkError("DisableBold")
}

// AlignLeft sets left alignment
func (m *MockPrinter) AlignLeft() error {
	m.record("AlignLeft")
	m.CurrentAlignment = "left"
	return m.checkError("AlignLeft")
}

// AlignCenter sets center alignment
func (m *MockPrinter) AlignCenter() error {
	m.record("AlignCenter")
	m.CurrentAlignment = "center"
	return m.checkError("AlignCenter")
}

// AlignRight sets right alignment
func (m *MockPrinter) AlignRight() error {
	m.record("AlignRight")
	m.CurrentAlignment = "right"
	return m.checkError("AlignRight")
}

// SetAlignment sets alignment to the specified value
func (m *MockPrinter) SetAlignment(align string) error {
	m.record("SetAlignment", align)
	m.CurrentAlignment = align
	return m.checkError("SetAlignment")
}

// SingleSize sets text size to 1x1
func (m *MockPrinter) SingleSize() error {
	m.record("SingleSize")
	m.CurrentSize = "1x1"
	return m.checkError("SingleSize")
}

// DoubleSize sets text size to 2x2
func (m *MockPrinter) DoubleSize() error {
	m.record("DoubleSize")
	m.CurrentSize = "2x2"
	return m.checkError("DoubleSize")
}

// TripleSize sets text size to 3x3
func (m *MockPrinter) TripleSize() error {
	m.record("TripleSize")
	m.CurrentSize = "3x3"
	return m.checkError("TripleSize")
}

// QuadraSize sets text size to 4x4
func (m *MockPrinter) QuadraSize() error {
	m.record("QuadraSize")
	m.CurrentSize = "4x4"
	return m.checkError("QuadraSize")
}

// PentaSize sets text size to 5x5
func (m *MockPrinter) PentaSize() error {
	m.record("PentaSize")
	m.CurrentSize = "5x5"
	return m.checkError("PentaSize")
}

// HexaSize sets text size to 6x6
func (m *MockPrinter) HexaSize() error {
	m.record("HexaSize")
	m.CurrentSize = "6x6"
	return m.checkError("HexaSize")
}

// HeptaSize sets text size to 7x7
func (m *MockPrinter) HeptaSize() error {
	m.record("HeptaSize")
	m.CurrentSize = "7x7"
	return m.checkError("HeptaSize")
}

// OctaSize sets text size to 8x8
func (m *MockPrinter) OctaSize() error {
	m.record("OctaSize")
	m.CurrentSize = "8x8"
	return m.checkError("OctaSize")
}

// CustomSize sets text size to specified width and height
func (m *MockPrinter) CustomSize(width, height byte) error {
	m.record("CustomSize", width, height)
	m.CurrentSize = fmt.Sprintf("%dx%d", width, height)
	return m.checkError("CustomSize")
}

// NoDot sets underline mode to none
func (m *MockPrinter) NoDot() error {
	m.record("NoDot")
	m.UnderlineMode = "none"
	return m.checkError("NoDot")
}

// OneDot sets underline mode to single
func (m *MockPrinter) OneDot() error {
	m.record("OneDot")
	m.UnderlineMode = "1pt"
	return m.checkError("OneDot")
}

// TwoDot sets underline mode to double
func (m *MockPrinter) TwoDot() error {
	m.record("TwoDot")
	m.UnderlineMode = "2pt"
	return m.checkError("TwoDot")
}

// InverseOn enables inverse mode
func (m *MockPrinter) InverseOn() error {
	m.record("InverseOn")
	m.InverseEnabled = true
	return m.checkError("InverseOn")
}

// InverseOff disables inverse mode
func (m *MockPrinter) InverseOff() error {
	m.record("InverseOff")
	m.InverseEnabled = false
	return m.checkError("InverseOff")
}

// FullFeedAndCut simulates full paper feed and cut
func (m *MockPrinter) FullFeedAndCut(lines byte) error {
	m.record("FullFeedAndCut", lines)
	return m.checkError("FullFeedAndCut")
}

// PartialFeedAndCut simulates partial paper feed and cut
func (m *MockPrinter) PartialFeedAndCut(lines byte) error {
	m.record("PartialFeedAndCut", lines)
	return m.checkError("PartialFeedAndCut")
}

// PrintBitmap simulates printing a monochrome bitmap
func (m *MockPrinter) PrintBitmap(bitmap *graphics.MonochromeBitmap) error {
	m.record("PrintBitmap", bitmap)
	m.PrintedBitmaps = append(m.PrintedBitmaps, bitmap)
	return m.checkError("PrintBitmap")
}

// PrintQR simulates printing a QR code
func (m *MockPrinter) PrintQR(data string, opts *graphics.QROptions) error {
	m.record("PrintQR", data, opts)
	m.PrintedQRs = append(m.PrintedQRs, data)
	return m.checkError("PrintQR")
}

// PrintBarcode simulates printing a barcode
func (m *MockPrinter) PrintBarcode(cfg graphics.BarcodeConfig, data []byte) error {
	m.record("PrintBarcode", cfg, data)
	return m.checkError("PrintBarcode")
}

// SetCodeTable sets the character code table
func (m *MockPrinter) SetCodeTable(codeTable character.CodeTable) error {
	m.record("SetCodeTable", codeTable)
	return m.checkError("SetCodeTable")
}
