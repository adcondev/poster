package executor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	posqr "github.com/adcondev/poster/pkg/commands/qrcode"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/service"
)

// QRCommand for QR handler
type QRCommand struct {
	Data        string `json:"data"`
	HumanText   string `json:"human_text,omitempty"`
	PixelWidth  int    `json:"pixel_width,omitempty"`
	Correction  string `json:"correction,omitempty"`
	Align       string `json:"align,omitempty"`
	Logo        string `json:"logo,omitempty"`
	CircleShape bool   `json:"circle_shape,omitempty"`
}

// handleQR manages QR code commands
func (e *Executor) handleQR(printer *service.Printer, data json.RawMessage) error {
	var cmd QRCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse QR command: %w", err)
	}

	// Validación de datos requeridos
	if cmd.Data == "" {
		return fmt.Errorf("QR data cannot be empty")
	}
	if len(cmd.Data) > posqr.MaxDataLength {
		return fmt.Errorf("QR data too long: %d bytes (maximum %d)", len(cmd.Data), posqr.MaxDataLength)
	}

	// Construir opciones con defaults
	opts := graphics.DefaultQROptions()

	// Configurar tamaño del QR
	if cmd.PixelWidth > 0 {
		if cmd.PixelWidth < constants.MinQrPixelWidth { // Mínimo según schema
			log.Printf("Warning: pixel_width %d too small, setting to minimum 87", cmd.PixelWidth)
			cmd.PixelWidth = constants.MinQrPixelWidth
		}
		opts.PixelWidth = cmd.PixelWidth
	} else {
		// Default: 128
		opts.PixelWidth = constants.DefaultQrPixelWidth
	}

	// Mapear corrección de errores (default: Q según schema)
	correction := "Q"
	if cmd.Correction != "" {
		correction = strings.ToUpper(cmd.Correction)
	}

	switch correction {
	case "L":
		opts.ErrorCorrection = posqr.LevelL
	case "M":
		opts.ErrorCorrection = posqr.LevelM
	case "Q":
		opts.ErrorCorrection = posqr.LevelQ
	case "H":
		opts.ErrorCorrection = posqr.LevelH
	default:
		return fmt.Errorf("invalid QR correction level: %s (valid: L, M, Q, H)", cmd.Correction)
	}

	// Configurar logo si se proporciona
	if cmd.Logo != "" {
		opts.LogoData = cmd.Logo
		log.Printf("QR with logo enabled")
	}

	// Configurar forma circular (solo para QR grandes)
	if cmd.CircleShape {
		if opts.PixelWidth > 256 {
			opts.CircleShape = true
			log.Printf("QR with circular blocks enabled")
		} else {
			log.Printf("Warning: circle_shape ignored (requires pixel_width > 256)")
		}
	}

	// Aplicar alineación (default: center según schema)
	align := constants.DefaultQrAlignment.String()
	if cmd.Align != "" {
		align = strings.ToLower(cmd.Align)
	}

	switch align {
	case constants.Center.String():
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	case constants.Right.String():
		if err := printer.AlignRight(); err != nil {
			return err
		}
	case constants.Left.String():
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	default:
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	}

	// Imprimir QR (automáticamente intenta nativo y fallback a imagen)
	if err := printer.PrintQR(cmd.Data, opts); err != nil {
		return fmt.Errorf("failed to print QR: %w", err)
	}

	// Imprimir texto humano si existe
	if cmd.HumanText != "" {
		// Mantener la alineación del QR para el texto
		if err := printer.PrintLine(cmd.HumanText); err != nil {
			return fmt.Errorf("failed to print QR human text: %w", err)
		}
	}

	// Restaurar alineación a la izquierda
	if err := printer.AlignLeft(); err != nil {
		return err
	}

	return nil
}
