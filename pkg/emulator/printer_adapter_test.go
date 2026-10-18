package emulator

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/adcondev/poster/pkg/commands/barcode"
	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/profile"
)

// Ensure *Printer structurally matches service.PrinterActions at compile time
// without depending on pkg/service directly.
type printerActions interface {
	Initialize() error
	Close() error
	Write(data []byte) error

	Print(text string) error
	PrintLine(text string) error
	FeedLines(lines byte) error

	FontA() error
	FontB() error
	EnableBold() error
	DisableBold() error

	AlignLeft() error
	AlignCenter() error
	AlignRight() error
	SetAlignment(align string) error

	SingleSize() error
	DoubleSize() error
	TripleSize() error
	QuadraSize() error
	PentaSize() error
	HexaSize() error
	HeptaSize() error
	OctaSize() error
	CustomSize(width, height byte) error

	NoDot() error
	OneDot() error
	TwoDot() error

	InverseOn() error
	InverseOff() error

	FullFeedAndCut(lines byte) error
	PartialFeedAndCut(lines byte) error

	Pulse(pin, onTime, offTime byte) error
	Beep(times, lapse byte) error
	SetLineSpacing(dots byte) error
	SelectDefaultLineSpacing() error

	PrintBitmap(bitmap *graphics.MonochromeBitmap) error
	PrintQR(data string, opts *graphics.QrOptions) error
	PrintBarcode(cfg graphics.BarcodeConfig, data []byte) error

	SetCodeTable(codeTable character.CodeTable) error
	GetProfile() *profile.Escpos
}

var _ printerActions = (*Printer)(nil)

func TestNewPrinter_DefaultProfile(t *testing.T) {
	p, err := NewPrinter(nil)
	if err != nil {
		t.Fatalf("NewPrinter(nil) failed: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil printer")
	}
	if p.Engine() == nil {
		t.Fatal("expected non-nil engine")
	}
	prof := p.GetProfile()
	if prof == nil {
		t.Fatal("expected non-nil profile")
	}
	if prof.PaperWidth != 80 {
		t.Errorf("expected 80mm paper width, got %v", prof.PaperWidth)
	}
}

func TestNewPrinter_CustomProfile(t *testing.T) {
	customProf := profile.CreateProfile58mm()
	p, err := NewPrinter(customProf)
	if err != nil {
		t.Fatalf("NewPrinter(customProf) failed: %v", err)
	}
	if p.GetProfile().PaperWidth != 58 {
		t.Errorf("expected 58mm paper width, got %v", p.GetProfile().PaperWidth)
	}
}

func TestPrinter_BasicText(t *testing.T) {
	p, err := NewPrinter(nil)
	if err != nil {
		t.Fatalf("NewPrinter failed: %v", err)
	}

	if err := p.Initialize(); err != nil {
		t.Errorf("Initialize() error = %v", err)
	}
	if err := p.Print("Hello "); err != nil {
		t.Errorf("Print() error = %v", err)
	}
	if err := p.PrintLine("World"); err != nil {
		t.Errorf("PrintLine() error = %v", err)
	}
	if err := p.FeedLines(2); err != nil {
		t.Errorf("FeedLines() error = %v", err)
	}
}

func TestPrinter_FontsAndStyles(t *testing.T) {
	p, err := NewPrinter(nil)
	if err != nil {
		t.Fatalf("NewPrinter failed: %v", err)
	}

	// Fonts
	if err := p.FontA(); err != nil {
		t.Errorf("FontA() error = %v", err)
	}
	if p.Engine().State().FontName != "A" {
		t.Errorf("expected Font A, got %s", p.Engine().State().FontName)
	}
	if err := p.FontB(); err != nil {
		t.Errorf("FontB() error = %v", err)
	}
	if p.Engine().State().FontName != "B" {
		t.Errorf("expected Font B, got %s", p.Engine().State().FontName)
	}

	// Bold
	if err := p.EnableBold(); err != nil {
		t.Errorf("EnableBold() error = %v", err)
	}
	if !p.Engine().State().IsBold {
		t.Error("expected bold enabled")
	}
	if err := p.DisableBold(); err != nil {
		t.Errorf("DisableBold() error = %v", err)
	}
	if p.Engine().State().IsBold {
		t.Error("expected bold disabled")
	}

	// Underline
	if err := p.NoDot(); err != nil {
		t.Errorf("NoDot() error = %v", err)
	}
	if err := p.OneDot(); err != nil {
		t.Errorf("OneDot() error = %v", err)
	}
	if p.Engine().State().IsUnderline != 1 {
		t.Errorf("expected underline 1, got %d", p.Engine().State().IsUnderline)
	}
	if err := p.TwoDot(); err != nil {
		t.Errorf("TwoDot() error = %v", err)
	}
	if p.Engine().State().IsUnderline != 2 {
		t.Errorf("expected underline 2, got %d", p.Engine().State().IsUnderline)
	}

	// Inverse
	if err := p.InverseOn(); err != nil {
		t.Errorf("InverseOn() error = %v", err)
	}
	if !p.Engine().State().IsInverse {
		t.Error("expected inverse enabled")
	}
	if err := p.InverseOff(); err != nil {
		t.Errorf("InverseOff() error = %v", err)
	}
	if p.Engine().State().IsInverse {
		t.Error("expected inverse disabled")
	}
}

func TestPrinter_Alignment(t *testing.T) {
	p, err := NewPrinter(nil)
	if err != nil {
		t.Fatalf("NewPrinter failed: %v", err)
	}

	if err := p.AlignLeft(); err != nil {
		t.Errorf("AlignLeft() error = %v", err)
	}
	if err := p.AlignCenter(); err != nil {
		t.Errorf("AlignCenter() error = %v", err)
	}
	if err := p.AlignRight(); err != nil {
		t.Errorf("AlignRight() error = %v", err)
	}
	if err := p.SetAlignment("left"); err != nil {
		t.Errorf("SetAlignment('left') error = %v", err)
	}
	if err := p.SetAlignment("center"); err != nil {
		t.Errorf("SetAlignment('center') error = %v", err)
	}
	if err := p.SetAlignment("right"); err != nil {
		t.Errorf("SetAlignment('right') error = %v", err)
	}
	if err := p.SetAlignment("invalid"); err == nil {
		t.Error("expected error for invalid alignment")
	}
}

func TestPrinter_Sizes(t *testing.T) {
	p, err := NewPrinter(nil)
	if err != nil {
		t.Fatalf("NewPrinter failed: %v", err)
	}

	sizes := []func() error{
		p.SingleSize, p.DoubleSize, p.TripleSize, p.QuadraSize,
		p.PentaSize, p.HexaSize, p.HeptaSize, p.OctaSize,
	}
	for i, sizeFn := range sizes {
		if err := sizeFn(); err != nil {
			t.Errorf("sizeFn %d error = %v", i+1, err)
		}
		expected := float64(i + 1)
		if p.Engine().State().ScaleW != expected || p.Engine().State().ScaleH != expected {
			t.Errorf("size %d scale mismatch: %v x %v", i+1, p.Engine().State().ScaleW, p.Engine().State().ScaleH)
		}
	}

	// Custom size with clamping
	if err := p.CustomSize(0, 10); err != nil {
		t.Errorf("CustomSize(0, 10) error = %v", err)
	}
	if p.Engine().State().ScaleW != 1.0 || p.Engine().State().ScaleH != 8.0 {
		t.Errorf("clamping failed: %v x %v", p.Engine().State().ScaleW, p.Engine().State().ScaleH)
	}
}

func TestPrinter_HardwareControls(t *testing.T) {
	p, err := NewPrinter(nil)
	if err != nil {
		t.Fatalf("NewPrinter failed: %v", err)
	}

	if err := p.FullFeedAndCut(2); err != nil {
		t.Errorf("FullFeedAndCut() error = %v", err)
	}
	if err := p.PartialFeedAndCut(1); err != nil {
		t.Errorf("PartialFeedAndCut() error = %v", err)
	}
	if err := p.Write([]byte{0x1b, 0x40}); err != nil {
		t.Errorf("Write() error = %v", err)
	}
	if err := p.Pulse(0, 10, 10); err != nil {
		t.Errorf("Pulse() error = %v", err)
	}
	if err := p.Beep(1, 1); err != nil {
		t.Errorf("Beep() error = %v", err)
	}
	if err := p.SetLineSpacing(30); err != nil {
		t.Errorf("SetLineSpacing() error = %v", err)
	}
	if err := p.SelectDefaultLineSpacing(); err != nil {
		t.Errorf("SelectDefaultLineSpacing() error = %v", err)
	}
	if err := p.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	if err := p.SetCodeTable(character.PC437); err != nil {
		t.Errorf("SetCodeTable() error = %v", err)
	}
	if p.GetProfile().CodeTable != character.PC437 {
		t.Errorf("expected code table PC437, got %v", p.GetProfile().CodeTable)
	}
}

func TestPrinter_Graphics(t *testing.T) {
	p, err := NewPrinter(nil)
	if err != nil {
		t.Fatalf("NewPrinter failed: %v", err)
	}

	// Bitmap
	if err := p.PrintBitmap(nil); err == nil {
		t.Error("expected error for nil bitmap")
	}
	bmp := graphics.NewMonochromeBitmap(32, 32)
	bmp.SetPixel(10, 10, true)
	if err := p.PrintBitmap(bmp); err != nil {
		t.Errorf("PrintBitmap error = %v", err)
	}

	// QR Code
	if err := p.PrintQR("", nil); err == nil {
		t.Error("expected error for empty QR data")
	}
	qrOpts := graphics.DefaultQROptions()
	if err := p.PrintQR("https://example.com", qrOpts); err != nil {
		t.Errorf("PrintQR error = %v", err)
	}

	// Barcode
	bCfg := graphics.DefaultBarcodeConfig()
	if err := p.PrintBarcode(bCfg, nil); err == nil {
		t.Error("expected error for empty barcode data")
	}
	bCfg.HRIPosition = barcode.HRIAbove
	if err := p.PrintBarcode(bCfg, []byte("12345678")); err != nil {
		t.Errorf("PrintBarcode (above) error = %v", err)
	}
	bCfg.HRIPosition = barcode.HRIBelow
	if err := p.PrintBarcode(bCfg, []byte("87654321")); err != nil {
		t.Errorf("PrintBarcode (below) error = %v", err)
	}
	bCfg.HRIPosition = barcode.HRINotPrinted
	if err := p.PrintBarcode(bCfg, []byte("99999999")); err != nil {
		t.Errorf("PrintBarcode (none) error = %v", err)
	}
}

func TestPrinter_RenderPNGAndSavePNG(t *testing.T) {
	p, err := NewPrinter(nil)
	if err != nil {
		t.Fatalf("NewPrinter failed: %v", err)
	}

	if err := p.PrintLine("Receipt Title"); err != nil {
		t.Fatalf("PrintLine failed: %v", err)
	}
	if err := p.PrintLine("Item 1: $10.00"); err != nil {
		t.Fatalf("PrintLine failed: %v", err)
	}
	if err := p.FullFeedAndCut(1); err != nil {
		t.Fatalf("FullFeedAndCut failed: %v", err)
	}

	pngBytes, err := p.RenderPNG()
	if err != nil {
		t.Fatalf("RenderPNG() failed: %v", err)
	}
	if len(pngBytes) == 0 {
		t.Fatal("expected non-empty PNG bytes")
	}

	// Verify valid PNG image
	img, err := png.Decode(bytes.NewReader(pngBytes))
	if err != nil {
		t.Fatalf("png.Decode failed: %v", err)
	}
	if img.Bounds().Dx() != 576 {
		t.Errorf("expected 576px width for 80mm printer, got %d", img.Bounds().Dx())
	}
	if img.Bounds().Dy() <= 0 {
		t.Errorf("expected positive height, got %d", img.Bounds().Dy())
	}

	// Test SavePNG
	tmpDir := t.TempDir()
	cleanPath := filepath.Clean(filepath.Join(tmpDir, "test_receipt.png"))
	if err := p.SavePNG(cleanPath); err != nil {
		t.Fatalf("SavePNG failed: %v", err)
	}

	readBytes, err := os.ReadFile(cleanPath)
	if err != nil {
		t.Fatalf("reading saved PNG failed: %v", err)
	}
	if !bytes.Equal(readBytes, pngBytes) {
		t.Error("saved file bytes do not match RenderPNG bytes")
	}
}
