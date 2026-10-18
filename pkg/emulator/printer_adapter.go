package emulator

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/profile"
)

// Printer adapts the emulator engine to the service.PrinterActions interface.
type Printer struct {
	engine  *Engine
	profile *profile.Escpos
}

// NewPrinter creates a new Printer adapter wrapping an emulator Engine.
func NewPrinter(prof *profile.Escpos) (*Printer, error) {
	if prof == nil {
		prof = profile.CreateProfile80mm()
	}

	dotsPerLine := prof.DotsPerLine
	if dotsPerLine <= 0 {
		dotsPerLine = constants.PaperPxWidth80mm
	}

	dpi := prof.DPI
	if dpi <= 0 {
		dpi = constants.DefaultDPI
	}

	cfg := Config{
		PaperPxWidth:            dotsPerLine,
		DPI:                     dpi,
		AutoAdjustCursorOnScale: true,
	}

	eng, err := NewEngine(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating emulator engine: %w", err)
	}

	return &Printer{
		engine:  eng,
		profile: prof,
	}, nil
}

// Engine returns the underlying emulator Engine.
func (p *Printer) Engine() *Engine {
	return p.engine
}

// RenderPNG renders the receipt canvas to PNG encoded bytes.
func (p *Printer) RenderPNG() ([]byte, error) {
	var buf bytes.Buffer
	if err := p.engine.WritePNG(&buf); err != nil {
		return nil, fmt.Errorf("writing png: %w", err)
	}
	return buf.Bytes(), nil
}

// SavePNG renders the receipt canvas and saves it to a file.
func (p *Printer) SavePNG(filename string) error {
	data, err := p.RenderPNG()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("saving png file %s: %w", filename, err)
	}
	return nil
}

// Initialize resets the printer state and canvas.
func (p *Printer) Initialize() error {
	p.engine.Reset()
	return nil
}

// Close closes any resources associated with the printer adapter.
func (p *Printer) Close() error {
	return nil
}

// Write is a no-op implementation for raw bytes on the visual emulator canvas.
func (p *Printer) Write(_ []byte) error {
	return nil
}

// Print renders text without a trailing line feed.
func (p *Printer) Print(text string) error {
	p.engine.Print(text)
	return nil
}

// PrintLine renders text with a trailing line feed.
func (p *Printer) PrintLine(text string) error {
	p.engine.PrintLine(text)
	return nil
}

// FeedLines advances the paper by the specified number of lines.
func (p *Printer) FeedLines(lines byte) error {
	p.engine.Feed(int(lines))
	return nil
}

// FontA selects font A.
func (p *Printer) FontA() error {
	p.engine.SetFont("A")
	return nil
}

// FontB selects font B.
func (p *Printer) FontB() error {
	p.engine.SetFont("B")
	return nil
}

// EnableBold enables bold text rendering.
func (p *Printer) EnableBold() error {
	p.engine.SetBold(true)
	return nil
}

// DisableBold disables bold text rendering.
func (p *Printer) DisableBold() error {
	p.engine.SetBold(false)
	return nil
}

// AlignLeft sets left alignment.
func (p *Printer) AlignLeft() error {
	p.engine.AlignLeft()
	return nil
}

// AlignCenter sets center alignment.
func (p *Printer) AlignCenter() error {
	p.engine.AlignCenter()
	return nil
}

// AlignRight sets right alignment.
func (p *Printer) AlignRight() error {
	p.engine.AlignRight()
	return nil
}

// SetAlignment sets text alignment to "left", "center", or "right".
func (p *Printer) SetAlignment(align string) error {
	switch strings.ToLower(strings.TrimSpace(align)) {
	case "left":
		p.engine.AlignLeft()
	case "center":
		p.engine.AlignCenter()
	case "right":
		p.engine.AlignRight()
	default:
		return fmt.Errorf("invalid alignment %q", align)
	}
	return nil
}

// SingleSize sets character size to 1x1.
func (p *Printer) SingleSize() error {
	p.engine.SetSize(1, 1)
	return nil
}

// DoubleSize sets character size to 2x2.
func (p *Printer) DoubleSize() error {
	p.engine.SetSize(2, 2)
	return nil
}

// TripleSize sets character size to 3x3.
func (p *Printer) TripleSize() error {
	p.engine.SetSize(3, 3)
	return nil
}

// QuadraSize sets character size to 4x4.
func (p *Printer) QuadraSize() error {
	p.engine.SetSize(4, 4)
	return nil
}

// PentaSize sets character size to 5x5.
func (p *Printer) PentaSize() error {
	p.engine.SetSize(5, 5)
	return nil
}

// HexaSize sets character size to 6x6.
func (p *Printer) HexaSize() error {
	p.engine.SetSize(6, 6)
	return nil
}

// HeptaSize sets character size to 7x7.
func (p *Printer) HeptaSize() error {
	p.engine.SetSize(7, 7)
	return nil
}

// OctaSize sets character size to 8x8.
func (p *Printer) OctaSize() error {
	p.engine.SetSize(8, 8)
	return nil
}

// CustomSize sets character size with width and height multipliers clamped to [1, 8].
func (p *Printer) CustomSize(width, height byte) error {
	w := int(width)
	if w < 1 {
		w = 1
	} else if w > 8 {
		w = 8
	}
	h := int(height)
	if h < 1 {
		h = 1
	} else if h > 8 {
		h = 8
	}
	p.engine.SetSize(w, h)
	return nil
}

// NoDot disables underline.
func (p *Printer) NoDot() error {
	p.engine.SetUnderline(0)
	return nil
}

// OneDot enables single-dot underline.
func (p *Printer) OneDot() error {
	p.engine.SetUnderline(1)
	return nil
}

// TwoDot enables double-dot underline.
func (p *Printer) TwoDot() error {
	p.engine.SetUnderline(2)
	return nil
}

// InverseOn enables white-on-black inverse mode.
func (p *Printer) InverseOn() error {
	p.engine.SetInverse(true)
	return nil
}

// InverseOff disables inverse mode.
func (p *Printer) InverseOff() error {
	p.engine.SetInverse(false)
	return nil
}

// FullFeedAndCut feeds specified lines and performs a full paper cut.
func (p *Printer) FullFeedAndCut(lines byte) error {
	if lines > 0 {
		p.engine.Feed(int(lines))
	}
	p.engine.Cut(false)
	return nil
}

// PartialFeedAndCut feeds specified lines and performs a partial paper cut.
func (p *Printer) PartialFeedAndCut(lines byte) error {
	if lines > 0 {
		p.engine.Feed(int(lines))
	}
	p.engine.Cut(true)
	return nil
}

// Pulse is a no-op on the visual emulator canvas.
func (p *Printer) Pulse(_, _, _ byte) error {
	return nil
}

// Beep is a no-op on the visual emulator canvas.
func (p *Printer) Beep(_, _ byte) error {
	return nil
}

// SetLineSpacing is a no-op on the visual emulator canvas.
func (p *Printer) SetLineSpacing(_ byte) error {
	return nil
}

// SelectDefaultLineSpacing is a no-op on the visual emulator canvas.
func (p *Printer) SelectDefaultLineSpacing() error {
	return nil
}

// PrintBitmap renders a monochrome bitmap onto the receipt canvas.
func (p *Printer) PrintBitmap(bitmap *graphics.MonochromeBitmap) error {
	if bitmap == nil {
		return fmt.Errorf("bitmap cannot be nil")
	}
	return p.engine.PrintImage(bitmap.ToImage())
}

// PrintQR generates and renders a QR code image onto the receipt canvas.
func (p *Printer) PrintQR(data string, opts *graphics.QrOptions) error {
	img, err := graphics.ProcessQRImage(data, opts)
	if err != nil {
		return fmt.Errorf("processing QR image: %w", err)
	}

	pipeline := graphics.NewPipeline(&graphics.ImgOptions{
		PixelWidth:     img.Bounds().Dx(),
		MaxWidth:       p.profile.DotsPerLine,
		Threshold:      128,
		Dithering:      graphics.Threshold,
		Scaling:        graphics.NearestNeighbor,
		PreserveAspect: true,
	})
	bitmap, err := pipeline.Process(img)
	if err != nil {
		return fmt.Errorf("processing QR bitmap pipeline: %w", err)
	}

	if err := p.engine.PrintImage(bitmap.ToImage()); err != nil {
		return fmt.Errorf("printing QR image: %w", err)
	}
	return nil
}

// PrintBarcode renders a clean visual barcode representation on the canvas.
func (p *Printer) PrintBarcode(cfg graphics.BarcodeConfig, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("barcode data cannot be empty")
	}

	height := int(cfg.Height)
	if height <= 0 {
		height = constants.DefaultBarcodeHeight
	}

	moduleWidth := int(cfg.Width)
	if moduleWidth <= 0 {
		moduleWidth = 2
	}

	barcodeWidth := len(data) * moduleWidth * 6
	if p.profile.DotsPerLine > 0 && barcodeWidth > p.profile.DotsPerLine {
		barcodeWidth = p.profile.DotsPerLine
	}
	if barcodeWidth <= 0 {
		barcodeWidth = 60
	}

	bitmap := graphics.NewMonochromeBitmap(barcodeWidth, height)
	for x := 0; x < barcodeWidth; x++ {
		charIdx := x / (moduleWidth * 6)
		if charIdx >= len(data) {
			charIdx = len(data) - 1
		}
		bitIdx := (x / moduleWidth) % 6
		b := data[charIdx]
		isBlack := ((b >> bitIdx) & 1) == 1
		if ((x / moduleWidth) % 2) == 0 {
			isBlack = !isBlack
		}
		if isBlack {
			for y := 0; y < height; y++ {
				bitmap.SetPixel(x, y, true)
			}
		}
	}

	hriAbove := cfg.HRIPosition == 0x01 || cfg.HRIPosition == 0x03 || cfg.HRIPosition == '1' || cfg.HRIPosition == '3'
	hriBelow := cfg.HRIPosition == 0x02 || cfg.HRIPosition == 0x03 || cfg.HRIPosition == '2' || cfg.HRIPosition == '3'

	if hriAbove {
		p.engine.PrintLine(string(data))
	}
	if err := p.engine.PrintImage(bitmap.ToImage()); err != nil {
		return fmt.Errorf("printing barcode image: %w", err)
	}
	if hriBelow {
		p.engine.PrintLine(string(data))
	}

	return nil
}

// SetCodeTable updates the active code table in the profile.
func (p *Printer) SetCodeTable(codeTable character.CodeTable) error {
	p.profile.CodeTable = codeTable
	return nil
}

// GetProfile returns the active printer profile.
func (p *Printer) GetProfile() *profile.Escpos {
	return p.profile
}
