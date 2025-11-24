// Package service provides implementations for various POS printer service.
package service

import (
	"fmt"
	"log"

	"github.com/adcondev/pos-printer/pkg/commands/character"
	"github.com/adcondev/pos-printer/pkg/commands/mechanismcontrol"
	"github.com/adcondev/pos-printer/pkg/composer"
	"github.com/adcondev/pos-printer/pkg/connection"
	"github.com/adcondev/pos-printer/pkg/graphics"
	"github.com/adcondev/pos-printer/pkg/profile"
	"github.com/adcondev/pos-printer/pkg/twodimensional"
)

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

// ============================================================================
// Basic Control Methods
// ============================================================================

// Initialize resets the printer to default settings
func (p *Printer) Initialize() error {
	// TODO: Add profile-specific initialization if needed
	var ct []byte
	if p.Profile.IsSupported(p.Profile.CodeTable) {
		ct, _ = p.Protocol.Character.SelectCharacterCodeTable(p.Profile.CodeTable)
	} else {
		ct, _ = p.Protocol.Character.SelectCharacterCodeTable(character.WPC1252)
		log.Printf("warning: unsupported code table %v, falling back to Windows-1252", p.Profile.CodeTable)
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
	encText, err := p.Profile.EncodeString(text)
	if err != nil {
		return err
	}
	cmd, err := p.Protocol.Print.Text(encText)
	if err != nil {
		return err
	}
	return p.Write(cmd)
}

// PrintLine sends text with line feed
func (p *Printer) PrintLine(text string) error {
	encText, err := p.Profile.EncodeString(text)
	if err != nil {
		return err
	}
	cmd, err := p.Protocol.PrintLn(encText)
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
	cmd, _ := p.Protocol.MechanismControl.FeedAndCutPaper(mechanismcontrol.FeedCutFull, lines)
	return p.Write(cmd)
}

// PartialFeedAndCut performs a partial paper cut
func (p *Printer) PartialFeedAndCut(lines byte) error {
	cmd, _ := p.Protocol.MechanismControl.FeedAndCutPaper(mechanismcontrol.FeedCutPartial, lines)
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
		log.Printf("warning: unsupported code table %v, falling back to Windows-1252", codeTable)
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

// PrintQR imprime un QR con detección automática y fallback
func (p *Printer) PrintQR(data string, opts *graphics.QROptions) error {
	if opts == nil {
		// TODO: Automatic options based on profile (DPI and Paper PixelWidth config to calculate Dots Per Line)
		opts = graphics.DefaultQROptions()
	}

	// Intentar QR nativo si está soportado
	opts.Qr = graphics.QrInfo{}
	opts.Logo = graphics.LogoInfo{}

	if p.Profile.HasQR {
		err := p.printQRNative(data, opts)
		if err == nil {
			return nil
		}
		log.Printf("Native QR failed, falling back to image: %v", err)
	}

	// Fallback a imagen
	return p.printQRAsImage(data, opts)
}

// printQRNative imprime usando protocolo ESC/POS nativo
func (p *Printer) printQRNative(data string, opts *graphics.QROptions) error {
	// Configurar modelo
	if cmd, err := p.Protocol.QRCode.SelectQRCodeModel(opts.Model, 0); err != nil {
		return err
	} else if err := p.Write(cmd); err != nil {
		return err
	}

	// Configurar tamaño de módulo
	_, err := opts.GenerateQR(data)
	if err != nil {
		return err
	}
	if cmd, err := p.Protocol.QRCode.SetQRCodeModuleSize(opts.GetModuleSize()); err != nil {
		return err
	} else if err := p.Write(cmd); err != nil {
		return err
	}

	// Configurar corrección de errores
	if cmd, err := p.Protocol.QRCode.SetQRCodeErrorCorrectionLevel(opts.ErrorCorrection); err != nil {
		return err
	} else if err := p.Write(cmd); err != nil {
		return err
	}

	// Almacenar datos
	if cmd, err := p.Protocol.QRCode.StoreQRCodeData([]byte(data)); err != nil {
		return err
	} else if err := p.Write(cmd); err != nil {
		return err
	}

	// Imprimir
	return p.Write(p.Protocol.QRCode.PrintQRCode())
}

// printQRAsImage genera y imprime QR como imagen
func (p *Printer) printQRAsImage(data string, opts *graphics.QROptions) error {
	// Generar imagen QR
	img, err := graphics.ProcessQRImage(data, opts)
	if err != nil {
		return fmt.Errorf("generate QR image: %w", err)
	}

	imgOpts := &graphics.ImgOptions{
		PixelWidth:     opts.PixelWidth,
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
func (p *Printer) PrintBarcode(cfg twodimensional.BarcodeConfig, data []byte) error {
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
