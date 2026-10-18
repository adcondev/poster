// Package service provides implementations for various POS printer service.
package service

import (
	"fmt"

	"github.com/adcondev/poster/internal/calculate"
	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/commands/mechanismcontrol"
	"github.com/adcondev/poster/pkg/composer"
	"github.com/adcondev/poster/pkg/connection"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/profile"
)

var _ PrinterActions = (*Printer)(nil)

// Printer represents a POS printer device
type Printer struct {
	Profile    profile.Escpos
	Connection connection.Connector
	Protocol   composer.EscposProtocol
}

// NewPrinter creates a new Printer instance
func NewPrinter(proto *composer.EscposProtocol, prof *profile.Escpos, conn connection.Connector) (*Printer, error) {
	if proto == nil {
		return nil, fmt.Errorf("protocol cannot be nil")
	}
	if conn == nil {
		return nil, fmt.Errorf("connection cannot be nil")
	}
	if prof == nil {
		return nil, fmt.Errorf("profile cannot be nil")
	}

	return &Printer{
		Profile:    *prof,
		Connection: conn,
		Protocol:   *proto,
	}, nil
}

// PrinterOption configures a Printer instance.
type PrinterOption func(*Printer) error

// WithProfile sets a custom printer profile.
func WithProfile(prof *profile.Escpos) PrinterOption {
	return func(p *Printer) error {
		if prof == nil {
			return fmt.Errorf("profile cannot be nil")
		}
		p.Profile = *prof
		return nil
	}
}

// WithProtocol sets a custom ESC/POS protocol.
func WithProtocol(proto *composer.EscposProtocol) PrinterOption {
	return func(p *Printer) error {
		if proto == nil {
			return fmt.Errorf("protocol cannot be nil")
		}
		p.Protocol = *proto
		return nil
	}
}

// NewPrinterWithOptions creates a new Printer instance with optional configuration.
func NewPrinterWithOptions(conn connection.Connector, opts ...PrinterOption) (*Printer, error) {
	if conn == nil {
		return nil, fmt.Errorf("connection cannot be nil")
	}

	p := &Printer{
		Profile:    *profile.CreateProfile80mm(),
		Connection: conn,
		Protocol:   *composer.NewEscpos(),
	}

	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}

	return p, nil
}

// ============================================================================
// Basic Control Methods
// ============================================================================

// Initialize resets the printer to default settings
func (p *Printer) Initialize() error {
	var ct []byte
	var err error
	if p.Profile.IsSupported(p.Profile.CodeTable) {
		ct, err = p.Protocol.Character.SelectCharacterCodeTable(p.Profile.CodeTable)
	} else {
		ct, err = p.Protocol.Character.SelectCharacterCodeTable(character.WPC1252)
	}
	if err != nil {
		return fmt.Errorf("select code table: %w", err)
	}

	init := append(p.Protocol.InitializePrinter(), ct...)
	return p.Write(init)
}

// Close closes the connection to the printer
func (p *Printer) Close() error {
	return p.Connection.Close()
}

// Write sends raw bytes directly to the printer
func (p *Printer) Write(data []byte) error {
	_, err := p.Connection.Write(data)
	return err
}

// ============================================================================
// Text Printing Methods
// ============================================================================

// Print sends text without line feed
func (p *Printer) Print(text string) error {
	encoded, err := p.Profile.EncodeBytesSafe(text)
	if err != nil {
		return err
	}
	cmd, err := p.Protocol.Print.Text(string(encoded))
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// PrintLine sends text with line feed
func (p *Printer) PrintLine(text string) error {
	encoded, err := p.Profile.EncodeBytesSafe(text)
	if err != nil {
		return fmt.Errorf("encode text: %w", err)
	}
	cmd, err := p.Protocol.PrintLn(string(encoded))
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// FeedLines advances paper by n lines
func (p *Printer) FeedLines(lines byte) error {
	return p.Write(p.Protocol.Print.PrintAndFeedLines(lines))
}

// ============================================================================
// Text Formatting Methods
// ============================================================================

// FontA sets the font to Font A
func (p *Printer) FontA() error {
	return p.Write(p.Protocol.SetFontA())
}

// FontB sets the font to Font B
func (p *Printer) FontB() error {
	return p.Write(p.Protocol.SetFontB())
}

// EnableBold enables or disables bold text
func (p *Printer) EnableBold() error {
	return p.Write(p.Protocol.Character.SetEmphasizedMode(character.OnEm))
}

// DisableBold disables bold text
func (p *Printer) DisableBold() error {
	return p.Write(p.Protocol.Character.SetEmphasizedMode(character.OffEm))
}

// AlignLeft sets left alignment
func (p *Printer) AlignLeft() error {
	cmd := p.Protocol.LeftAlign()
	return p.Write(cmd)
}

// AlignCenter sets center alignment
func (p *Printer) AlignCenter() error {
	cmd := p.Protocol.CenterAlign()
	return p.Write(cmd)
}

// AlignRight sets right alignment
func (p *Printer) AlignRight() error {
	cmd := p.Protocol.RightAlign()
	return p.Write(cmd)
}

// SetAlignment sets text alignment
func (p *Printer) SetAlignment(align string) error {
	switch align {
	case constants.Left.String():
		return p.AlignLeft()
	case constants.Center.String():
		return p.AlignCenter()
	case constants.Right.String():
		return p.AlignRight()
	default:
		return nil
	}
}

// SingleSize resets text to normal size
func (p *Printer) SingleSize() error {
	return p.Write(p.Protocol.SingleSizeText())
}

// DoubleSize enables or disables double width
func (p *Printer) DoubleSize() error {
	return p.Write(p.Protocol.DoubleSizeText())
}

// TripleSize enables or disables triple width
func (p *Printer) TripleSize() error {
	return p.Write(p.Protocol.TripleSizeText())
}

// QuadraSize enables or disables quadruple width
func (p *Printer) QuadraSize() error {
	return p.Write(p.Protocol.QuadraSizeText())
}

// PentaSize enables or disables pentuple width
func (p *Printer) PentaSize() error {
	return p.Write(p.Protocol.PentaSizeText())
}

// HexaSize enables or disables hexaple width
func (p *Printer) HexaSize() error {
	return p.Write(p.Protocol.HexaSizeText())
}

// HeptaSize enables or disables heptuple width
func (p *Printer) HeptaSize() error {
	return p.Write(p.Protocol.HeptaSizeText())
}

// OctaSize enables or disables octuple width
func (p *Printer) OctaSize() error {
	return p.Write(p.Protocol.OctaSizeText())
}

// CustomSize sets a custom text size
func (p *Printer) CustomSize(width, height byte) error {
	return p.Write(p.Protocol.CustomSizeText(width, height))
}

// OneDot enables single underline
func (p *Printer) OneDot() error {
	return p.Write(p.Protocol.OneDotUnderline())
}

// TwoDot enables double underline
func (p *Printer) TwoDot() error {
	return p.Write(p.Protocol.TwoDotUnderline())
}

// NoDot disables underline
func (p *Printer) NoDot() error {
	return p.Write(p.Protocol.DisableUnderline())
}

// InverseOn enables reverse mode
func (p *Printer) InverseOn() error {
	return p.Write(p.Protocol.EnableReverseMode())
}

// InverseOff disables reverse mode
func (p *Printer) InverseOff() error {
	return p.Write(p.Protocol.DisableReverseMode())
}

// ============================================================================
// Paper Control Methods
// ============================================================================

// FullFeedAndCut performs a full paper cut
func (p *Printer) FullFeedAndCut(lines byte) error {
	cmd, err := p.Protocol.MechanismControl.FeedAndCutPaper(mechanismcontrol.FeedCutFull, lines)
	if err != nil {
		return fmt.Errorf("feed and cut paper: %w", err)
	}
	return p.Write(cmd)
}

// PartialFeedAndCut performs a partial paper cut
func (p *Printer) PartialFeedAndCut(lines byte) error {
	cmd, err := p.Protocol.MechanismControl.FeedAndCutPaper(mechanismcontrol.FeedCutPartial, lines)
	if err != nil {
		return fmt.Errorf("feed and cut paper: %w", err)
	}
	return p.Write(cmd)
}

// Pulse generates a cash drawer pulse on pin 2 or pin 5
func (p *Printer) Pulse(pin, onTime, offTime byte) error {
	cmd, err := p.Protocol.Pulse(pin, onTime, offTime)
	if err != nil {
		return fmt.Errorf("pulse command: %w", err)
	}
	return p.Write(cmd)
}

// Beep triggers the printer buzzer
func (p *Printer) Beep(times, lapse byte) error {
	cmd, err := p.Protocol.Beep(times, lapse)
	if err != nil {
		return fmt.Errorf("beep command: %w", err)
	}
	return p.Write(cmd)
}

// SetLineSpacing sets the line spacing to the specified number of dots
func (p *Printer) SetLineSpacing(dots byte) error {
	cmd := p.Protocol.SetLineSpacing(dots)
	return p.Write(cmd)
}

// SelectDefaultLineSpacing resets line spacing to the printer default
func (p *Printer) SelectDefaultLineSpacing() error {
	cmd := p.Protocol.SelectDefaultLineSpacing()
	return p.Write(cmd)
}

// ============================================================================
// Image Printing Methods
// ============================================================================

// PrintBitmap prints a monochrome bitmap using raster graphics
func (p *Printer) PrintBitmap(bitmap *graphics.MonochromeBitmap) error {
	if bitmap == nil {
		return fmt.Errorf("bitmap cannot be nil")
	}

	width := bitmap.GetWidthBytes()
	height := bitmap.Height

	if width < 0 {
		return fmt.Errorf("invalid bitmap width: %d", width)
	}
	if height < 0 {
		return fmt.Errorf("invalid bitmap height: %d", height)
	}

	const maxUint16 = 1<<16 - 1
	if width > maxUint16 {
		return fmt.Errorf("bitmap width in bytes %d exceeds uint16 max %d", width, maxUint16)
	}
	if height > maxUint16 {
		return fmt.Errorf("bitmap height %d exceeds uint16 max %d", height, maxUint16)
	}

	if p.Profile.MaxBandHeight > 0 && bitmap.Height > p.Profile.MaxBandHeight {
		for y := 0; y < bitmap.Height; y += p.Profile.MaxBandHeight {
			band := bitmap.SubBand(y, p.Profile.MaxBandHeight)
			bandWidth := band.GetWidthBytes()
			bandHeight := band.Height
			if bandWidth < 0 || bandWidth > maxUint16 || bandHeight < 0 || bandHeight > maxUint16 {
				return fmt.Errorf("band dimensions out of range at y=%d: %dx%d", y, bandWidth, bandHeight)
			}
			cmd, err := p.Protocol.BitImage.PrintRasterBitImage(
				0, // normal mode
				uint16(bandWidth),
				uint16(bandHeight),
				band.GetRasterData(),
			)
			if err != nil {
				return fmt.Errorf("generate raster band at y=%d: %w", y, err)
			}
			if err := p.Write(cmd); err != nil {
				return fmt.Errorf("write raster band at y=%d: %w", y, err)
			}
		}
		return nil
	}

	cmd, err := p.Protocol.BitImage.PrintRasterBitImage(
		0, // normal mode
		uint16(width),
		uint16(height),
		bitmap.GetRasterData(),
	)
	if err != nil {
		return fmt.Errorf("generate raster command: %w", err)
	}

	return p.Write(cmd)
}

// ============================================================================
// Character Code Table Methods
// ============================================================================

// SetCodeTable changes the character code table
func (p *Printer) SetCodeTable(codeTable character.CodeTable) error {
	var cmd []byte
	var err error
	if !p.Profile.IsSupported(codeTable) {
		cmd, err = p.Protocol.Character.SelectCharacterCodeTable(character.WPC1252)
	} else {
		cmd, err = p.Protocol.Character.SelectCharacterCodeTable(codeTable)
	}
	if err != nil {
		return fmt.Errorf("set code table: %w", err)
	}
	if err := p.Write(cmd); err != nil {
		return fmt.Errorf("write code table command: %w", err)
	}
	p.Profile.CodeTable = codeTable
	return nil
}

// ============================================================================
// QR Code Printing Methods
// ============================================================================

// calculateMaxPixelWidth calculates the maximum printable pixel width based on printer profile.
func (p *Printer) calculateMaxPixelWidth() int {
	if p.Profile.DotsPerLine > 0 {
		return p.Profile.DotsPerLine
	}

	if p.Profile.PrintWidth > 0 && p.Profile.DPI > 0 {
		return calculate.MmToDots(float64(p.Profile.PrintWidth), p.Profile.DPI)
	}

	if p.Profile.PaperWidth > 0 && p.Profile.DPI > 0 {
		var printWidthMm float64
		if p.Profile.PaperWidth >= 78 { // ~80mm
			printWidthMm = 72.0
		} else { // ~58mm or smaller
			printWidthMm = 48.0
		}
		return calculate.MmToDots(printWidthMm, p.Profile.DPI)
	}

	return constants.PaperPxWidth80mm
}

// PrintQR imprime un QR con detección automática y fallback
func (p *Printer) PrintQR(data string, qropts *graphics.QrOptions) error {
	maxWidth := p.calculateMaxPixelWidth()
	if qropts == nil {
		qropts = graphics.DefaultQROptions()
		qropts.MaxPixelWidth = maxWidth
	}

	qropts.MaxPixelWidth = maxWidth

	// Use native ESC/POS QR only when supported AND neither logo nor circle_shape is requested
	qropts.Qr = graphics.QrInfo{}
	if qropts.Logo.Image == nil {
		qropts.Logo = graphics.LogoInfo{}
	}

	useNativeQR := p.Profile.HasQR && qropts.LogoData == "" && qropts.Logo.Image == nil && !qropts.CircleShape
	if useNativeQR {
		err := p.printQRNative(data, qropts)
		if err == nil {
			return nil
		}
	}

	// Fallback / software bitmap rendering (required for logo, circle_shape, or HasQR == false)
	return p.printQRAsImage(data, qropts)
}

// printQRNative imprime usando protocolo ESC/POS nativo
func (p *Printer) printQRNative(data string, qropts *graphics.QrOptions) error {
	// 1. Select Model
	if cmd, err := p.Protocol.QRCode.SelectQRCodeModel(qropts.Model, 0); err != nil {
		return err
	} else if err := p.Write(cmd); err != nil {
		return err
	}

	// 2. Calculate Module Size (Pure Go calculation)
	// Note: GenerateQR modifies opts.Qr.moduleSize based on opts.PixelWidth and Data
	_, err := qropts.GenerateQR(data)
	if err != nil {
		return err
	}

	// 3. Set Module Size
	if cmd, err := p.Protocol.QRCode.SetQRCodeModuleSize(qropts.GetModuleSize()); err != nil {
		return err
	} else if err := p.Write(cmd); err != nil {
		return err
	}

	// 4. Set Error Correction Level
	if cmd, err := p.Protocol.QRCode.SetQRCodeErrorCorrectionLevel(qropts.ErrorCorrection); err != nil {
		return err
	} else if err := p.Write(cmd); err != nil {
		return err
	}

	// 5. Store Data
	if cmd, err := p.Protocol.QRCode.StoreQRCodeData([]byte(data)); err != nil {
		return err
	} else if err := p.Write(cmd); err != nil {
		return err
	}

	// 6. Print
	return p.Write(p.Protocol.QRCode.PrintQRCode())
}

// printQRAsImage genera y imprime QR como imagen
func (p *Printer) printQRAsImage(data string, opts *graphics.QrOptions) error {
	// Generar imagen QR
	img, err := graphics.ProcessQRImage(data, opts)
	if err != nil {
		return fmt.Errorf("generate QR image: %w", err)
	}

	targetWidth := opts.PixelWidth
	generatedWidth := img.Bounds().Dx()

	if generatedWidth > targetWidth {
		targetWidth = generatedWidth
	}

	imgOpts := &graphics.ImgOptions{
		PixelWidth:     targetWidth,
		Threshold:      128,
		Scaling:        graphics.NearestNeighbor,
		Dithering:      graphics.Threshold,
		PreserveAspect: true,
	}

	if opts.Logo.Image != nil {
		imgOpts.Dithering = graphics.Atkinson
	}

	// Procesar imagen para impresora térmica
	pipeline := graphics.NewPipeline(imgOpts)

	bitmap, err := pipeline.Process(img)
	if err != nil {
		return fmt.Errorf("process QR image: %w", err)
	}

	return p.PrintBitmap(bitmap)
}

// ============================================================================
// Barcode Printing Methods
// ============================================================================

// PrintBarcode imprime un código de barras configurando todos sus parámetros
// en una sola transmisión para asegurar consistencia (stateless).
func (p *Printer) PrintBarcode(cfg graphics.BarcodeConfig, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("barcode data cannot be empty")
	}

	// Delegamos la construcción completa al Composer
	fullCommand, err := p.Protocol.GenerateBarcode(cfg, data)
	if err != nil {
		return fmt.Errorf("compose barcode: %w", err)
	}

	// Enviamos el bloque atómico a la impresora
	return p.Write(fullCommand)
}

// GetProfile returns a pointer to the printer's profile configuration
func (p *Printer) GetProfile() *profile.Escpos {
	return &p.Profile
}
