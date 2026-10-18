package executor

import (
	"encoding/json"
	"fmt"
	"strings"

	posqr "github.com/adcondev/poster/pkg/commands/qrcode"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/service"
)

// QRCommand for QR handler
type QRCommand = schema.QRCommand

// handleQR manages QR code commands
func (e *Executor) handleQR(printer service.PrinterActions, data json.RawMessage) error {
	var cmd QRCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse QR command: %w", err)
	}

	if cmd.Data == "" {
		return fmt.Errorf("QR data cannot be empty")
	}
	if len(cmd.Data) > posqr.MaxDataLength {
		return fmt.Errorf("QR data too long: %d bytes (maximum %d)", len(cmd.Data), posqr.MaxDataLength)
	}

	opts := graphics.DefaultQROptions()

	if cmd.PixelWidth > 0 {
		if cmd.PixelWidth < constants.MinQrPixelWidth {
			cmd.PixelWidth = constants.MinQrPixelWidth
		}
		opts.PixelWidth = cmd.PixelWidth
	} else {
		opts.PixelWidth = constants.DefaultQrPixelWidth
	}

	correction := constants.DefaultQrErrorCorrection.String()
	if cmd.Correction != "" {
		correction = strings.ToUpper(cmd.Correction)
	}

	switch correction {
	case constants.L.String():
		opts.ErrorCorrection = posqr.LevelL
	case constants.M.String():
		opts.ErrorCorrection = posqr.LevelM
	case constants.Q.String():
		opts.ErrorCorrection = posqr.LevelQ
	case constants.H.String():
		opts.ErrorCorrection = posqr.LevelH
	default:
		return fmt.Errorf("invalid QR correction level: %s (valid: L, M, Q, H)", cmd.Correction)
	}

	if cmd.Logo != "" {
		opts.LogoData = cmd.Logo
	}

	if cmd.CircleShape && opts.PixelWidth > 256 {
		opts.CircleShape = true
	}

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

	if err := printer.PrintQR(cmd.Data, opts); err != nil {
		return fmt.Errorf("failed to print QR: %w", err)
	}

	if cmd.HumanText != "" {
		if err := printer.PrintLine(cmd.HumanText); err != nil {
			return err
		}
	}

	return printer.AlignLeft()
}
