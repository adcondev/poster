// Package composer provides a high-level interface for composing ESC/POS commands
// for POS printers by combining various capabilities.
package composer

import (
	"bytes"
	"fmt"

	"github.com/adcondev/poster/pkg/commands/barcode"
	"github.com/adcondev/poster/pkg/commands/bitimage"
	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/commands/linespacing"
	"github.com/adcondev/poster/pkg/commands/mechanismcontrol"
	"github.com/adcondev/poster/pkg/commands/print"
	"github.com/adcondev/poster/pkg/commands/printposition"
	"github.com/adcondev/poster/pkg/commands/qrcode"
	"github.com/adcondev/poster/pkg/commands/shared"
	"github.com/adcondev/poster/pkg/graphics"
)

// EscposProtocol implements the ESCPOS Commands
type EscposProtocol struct {
	Barcode          barcode.Capability
	BitImage         bitimage.Capability
	Character        character.Capability
	LineSpacing      linespacing.Capability
	MechanismControl mechanismcontrol.Capability
	Print            print.Capability
	PrintPosition    printposition.Capability
	QRCode           qrcode.Capability
	// TODO: Implement other capabilities
	// PrintingPaper    printingpaper.Capability
	// PaperSensor      papersensor.Capability
	// PanelButton      panelbutton.Capability
	// Status           status.Capability
	// MacroFunctions   macrofunctions.Capability
	// Kanji 		    kanji.Capability
	// Miscellaneous 	miscellaneous.Capability
	// Customize 	    customize.Capability
	// CounterPrinting  counterprinting.Capability
	// PDF417           pdf417.Capability
	// MaxiCode         maxicode.Capability
	// DataBar          databar.Capability
	// CompositeSym     compositesym.Capability
	// AztecCode        azteccode.Capability
	// DataMatrix       datamatrix.Capability
}

// NewEscpos creates a new instance of the ESC/POS protocol
func NewEscpos() *EscposProtocol {
	return &EscposProtocol{
		Barcode:          barcode.NewCommands(),
		BitImage:         bitimage.NewCommands(),
		Character:        character.NewCommands(),
		LineSpacing:      linespacing.NewCommands(),
		MechanismControl: mechanismcontrol.NewCommands(),
		Print:            print.NewCommands(),
		PrintPosition:    printposition.NewCommands(),
		QRCode:           qrcode.NewCommands(),
	}
}

// TODO: Implement other methods to access capabilities related to initialization and state management (miscellaneous)

// InitializePrinter provides a reset of the printer to its power-on state for RAM settings.
//
// Format:
//
//	ASCII:   ESC @
//	Hex:     0x1B 0x40
//	Decimal: 27 64
//
// Range:
//
//	None
//
// Default:
//
//	None
//
// Parameters:
//
//	None
//
// Notes:
//   - Clears the data in the Print buffer and resets printer modes to those at power-on.
//   - Macro definitions are NOT cleared.
//   - Offline response selection is NOT cleared.
//   - Contents of user NV memory are NOT cleared.
//   - NV graphics (NV bit image) and NV user memory are NOT cleared.
//   - The maintenance counter value is NOT affected by this command.
//   - Software setting values are NOT cleared.
//   - DIP switch settings are NOT re-read.
//   - The data in the receiver buffer is NOT cleared.
//   - In Page mode: deletes data in Print areas, initializes all settings, and selects Standard mode.
//   - Cancels many active settings (Print mode, line feed, etc.) and moves the Print position to the left side
//     of the printable area; printer status becomes "Beginning of the line".
//   - Certain ESC = behavior is preserved/adjusted as described by the printer (see model notes).
//   - Use with care when expecting persistent RAM/NV behavior — only RAM settings are reset to power-on defaults.
//
// Errors:
//
//	This function is safe and does not return errors
func (c *EscposProtocol) InitializePrinter() []byte {
	return []byte{shared.ESC, '@'}
}

// ============================================================================
// Minimal Print Methods
// ============================================================================

// LeftMargin sets the left margin
func (c *EscposProtocol) LeftMargin(margin uint16) []byte {
	return c.PrintPosition.SetLeftMargin(margin)
}

// PrintWidth sets the Print area width
func (c *EscposProtocol) PrintWidth(width uint16) []byte {
	return c.PrintPosition.SetPrintAreaWidth(width)
}

// PrintLn sends text to the printer followed by a line feed.
func (c *EscposProtocol) PrintLn(text string) ([]byte, error) {
	cmd, err := c.Print.Text(text)
	if err != nil {
		return nil, fmt.Errorf("println: text: %w", err)
	}
	cmd = append(cmd, c.Print.PrintAndLineFeed()...)
	return cmd, nil
}

// DisableBold disables bold mode.
func (c *EscposProtocol) DisableBold() []byte {
	return c.Character.SetEmphasizedMode(character.OffEm)
}

// SetAlign sets the text alignment.
func (c *EscposProtocol) SetAlign(mode printposition.Justification) ([]byte, error) {
	cmd, err := c.PrintPosition.SelectJustification(mode)
	if err != nil {
		return nil, fmt.Errorf("set align: select justification: %w", err)
	}
	return cmd, nil
}

// CenterAlign centers the text.
func (c *EscposProtocol) CenterAlign() []byte {
	cmd, _ := c.PrintPosition.SelectJustification(printposition.Center)
	return cmd
}

// LeftAlign left-aligns the text.
func (c *EscposProtocol) LeftAlign() []byte {
	cmd, _ := c.PrintPosition.SelectJustification(printposition.Left)
	return cmd
}

// RightAlign right-aligns the text.
func (c *EscposProtocol) RightAlign() []byte {
	cmd, _ := c.PrintPosition.SelectJustification(printposition.Right)
	return cmd
}

// SingleSizeText sets the smallest(regular) text size.
func (c *EscposProtocol) SingleSizeText() []byte {
	size, _ := character.NewSize(1, 1)
	return c.Character.SelectCharacterSize(size)
}

// DoubleSizeText sets double size text.
func (c *EscposProtocol) DoubleSizeText() []byte {
	size, _ := character.NewSize(2, 2)
	return c.Character.SelectCharacterSize(size)
}

// TripleSizeText sets triple size text.
func (c *EscposProtocol) TripleSizeText() []byte {
	size, _ := character.NewSize(3, 3)
	return c.Character.SelectCharacterSize(size)
}

// QuadraSizeText sets quadruple size text.
func (c *EscposProtocol) QuadraSizeText() []byte {
	size, _ := character.NewSize(4, 4)
	return c.Character.SelectCharacterSize(size)
}

// PentaSizeText sets pentuple size text.
func (c *EscposProtocol) PentaSizeText() []byte {
	size, _ := character.NewSize(5, 5)
	return c.Character.SelectCharacterSize(size)
}

// HexaSizeText sets hexuple size text.
func (c *EscposProtocol) HexaSizeText() []byte {
	size, _ := character.NewSize(6, 6)
	return c.Character.SelectCharacterSize(size)
}

// HeptaSizeText sets heptuple size text.
func (c *EscposProtocol) HeptaSizeText() []byte {
	size, _ := character.NewSize(7, 7)
	return c.Character.SelectCharacterSize(size)
}

// OctaSizeText sets octuple size text.
func (c *EscposProtocol) OctaSizeText() []byte {
	size, _ := character.NewSize(8, 8)
	return c.Character.SelectCharacterSize(size)
}

// CustomSizeText sets a custom text size.
func (c *EscposProtocol) CustomSizeText(widthMultiplier, heightMultiplier byte) []byte {
	if widthMultiplier < 1 {
		widthMultiplier = 1
	}
	if widthMultiplier > 8 {
		widthMultiplier = 8
	}
	if heightMultiplier < 1 {
		heightMultiplier = 1
	}
	if heightMultiplier > 8 {
		heightMultiplier = 8
	}
	size, _ := character.NewSize(widthMultiplier, heightMultiplier)
	return c.Character.SelectCharacterSize(size)
}

// EnableBold enables bold text.
func (c *EscposProtocol) EnableBold() []byte {
	return c.Character.SetEmphasizedMode(character.OnEm)
}

// OneDotUnderline enables one-dot underline.
func (c *EscposProtocol) OneDotUnderline() []byte {
	cmd, _ := c.Character.SetUnderlineMode(character.OneDot)
	return cmd
}

// TwoDotUnderline enables two-dot underline.
func (c *EscposProtocol) TwoDotUnderline() []byte {
	cmd, _ := c.Character.SetUnderlineMode(character.TwoDot)
	return cmd
}

// DisableUnderline disables underline.
func (c *EscposProtocol) DisableUnderline() []byte {
	cmd, _ := c.Character.SetUnderlineMode(character.NoDot)
	return cmd
}

// EnableReverseMode enables reverse mode (white on black).
func (c *EscposProtocol) EnableReverseMode() []byte {
	return c.Character.SetWhiteBlackReverseMode(character.OnRm)
}

// DisableReverseMode disables reverse mode (black on white).
func (c *EscposProtocol) DisableReverseMode() []byte {
	return c.Character.SetWhiteBlackReverseMode(character.OffRm)
}

// SetFontA selects Font A.
func (c *EscposProtocol) SetFontA() []byte {
	cmd, _ := c.Character.SelectCharacterFont(character.FontA)
	return cmd
}

// SetFontB selects Font B.
func (c *EscposProtocol) SetFontB() []byte {
	cmd, _ := c.Character.SelectCharacterFont(character.FontB)
	return cmd
}

// GenerateBarcode crea una secuencia completa y atómica para imprimir un código de barras.
// Incluye comandos de configuración (ancho, alto, fuentes) seguidos inmediatamente por los datos.
func (c *EscposProtocol) GenerateBarcode(cfg graphics.BarcodeConfig, data []byte) ([]byte, error) {
	var buffer bytes.Buffer

	// 1. Configurar Ancho del Módulo (GS w)
	// Si el valor es 0, usamos default o ignoramos (según tu lógica de negocio).
	// Aquí asumo que queremos asegurar el estado, así que enviamos el comando.
	cmd, err := c.Barcode.SetBarcodeWidth(cfg.Width)
	if err != nil {
		return nil, fmt.Errorf("config width: %w", err)
	}
	buffer.Write(cmd)

	// 2. Configurar Altura (GS h)
	cmd, err = c.Barcode.SetBarcodeHeight(cfg.Height)
	if err != nil {
		return nil, fmt.Errorf("config height: %w", err)
	}
	buffer.Write(cmd)

	// 3. Configurar Posición HRI (GS H)
	cmd, err = c.Barcode.SelectHRICharacterPosition(cfg.HRIPosition)
	if err != nil {
		return nil, fmt.Errorf("config HRI pos: %w", err)
	}
	buffer.Write(cmd)

	// 4. Configurar Fuente HRI (GS f)
	cmd, err = c.Barcode.SelectFontForHRI(cfg.HRIFont)
	if err != nil {
		return nil, fmt.Errorf("config HRI font: %w", err)
	}
	buffer.Write(cmd)

	// 5. Generar Comando de Impresión (GS k)
	var printCmd []byte

	// Lógica especial para CODE128 segura vs Estándar
	if cfg.Symbology == barcode.CODE128 {
		// Usamos la función segura que escapa caracteres e inyecta el CodeSet
		printCmd, err = c.Barcode.PrintBarcodeWithCodeSet(cfg.Symbology, cfg.CodeSet, data)
	} else {
		// Impresión estándar
		printCmd, err = c.Barcode.PrintBarcode(cfg.Symbology, data)
	}

	if err != nil {
		return nil, fmt.Errorf("generate payload: %w", err)
	}
	buffer.Write(printCmd)

	return buffer.Bytes(), nil
}
