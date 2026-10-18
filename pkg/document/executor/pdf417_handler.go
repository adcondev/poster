package executor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adcondev/poster/pkg/commands/qrcode"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/service"
)

// handlePDF417 manages PDF417 2D barcode commands.
func (e *Executor) handlePDF417(printer service.PrinterActions, data json.RawMessage) error {
	var cmd schema.PDF417Command
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse PDF417 command: %w", err)
	}

	if cmd.Data == "" {
		return fmt.Errorf("PDF417 data cannot be empty")
	}

	columns := byte(0)
	if cmd.Columns >= 0 && cmd.Columns <= 30 {
		columns = byte(cmd.Columns)
	}

	rows := byte(0)
	if cmd.Rows >= 3 && cmd.Rows <= 90 {
		rows = byte(cmd.Rows)
	}

	moduleWidth := byte(3)
	if cmd.ModuleWidth >= 2 && cmd.ModuleWidth <= 8 {
		moduleWidth = byte(cmd.ModuleWidth)
	}

	rowHeight := byte(3)
	if cmd.RowHeight >= 2 && cmd.RowHeight <= 8 {
		rowHeight = byte(cmd.RowHeight)
	}

	ecLevel := byte(2)
	if cmd.ErrorCorrection > 0 && cmd.ErrorCorrection <= 8 {
		ecLevel = byte(cmd.ErrorCorrection)
	}

	align := constants.DefaultQrAlignment.String()
	if cmd.Align != "" {
		align = strings.ToLower(cmd.Align)
	}

	switch align {
	case constants.Left.String():
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	case constants.Right.String():
		if err := printer.AlignRight(); err != nil {
			return err
		}
	case constants.Center.String():
		fallthrough
	default:
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	}

	rawCmd, err := qrcode.ComposePDF417(columns, rows, moduleWidth, rowHeight, ecLevel, []byte(cmd.Data))
	if err != nil {
		return fmt.Errorf("failed to compose PDF417 command: %w", err)
	}

	if err := printer.Write(rawCmd); err != nil {
		return fmt.Errorf("failed to write PDF417 command: %w", err)
	}

	return printer.AlignLeft()
}
