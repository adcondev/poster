package executor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adcondev/poster/pkg/commands/barcode"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/graphics"
	"github.com/adcondev/poster/pkg/service"
)

// BarcodeCommand for barcode handler
type BarcodeCommand = schema.BarcodeCommand

// handleBarcode manages barcode commands
func (e *Executor) handleBarcode(printer service.PrinterActions, data json.RawMessage) error {
	var cmd BarcodeCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse barcode command: %w", err)
	}

	if cmd.Symbology == "" {
		return fmt.Errorf("barcode symbology is required")
	}
	if cmd.Data == "" {
		return fmt.Errorf("barcode data is required")
	}

	cfg := graphics.DefaultBarcodeConfig()

	symbology, err := graphics.MapSymbology(cmd.Symbology)
	if err != nil {
		symbology = barcode.CODE128
	}
	cfg.Symbology = symbology

	if cmd.Width != nil {
		switch {
		case *cmd.Width < constants.MinBarcodeModuleWidth:
			cfg.Width = constants.MinBarcodeModuleWidth
		case *cmd.Width > constants.MaxBarcodeModuleWidth:
			cfg.Width = constants.MaxBarcodeModuleWidth
		default:
			cfg.Width = barcode.Width(*cmd.Width)
		}
	} else {
		cfg.Width = constants.DefaultBarcodeModuleWidth
	}

	if cmd.Height != nil {
		switch {
		case *cmd.Height < constants.MinBarcodeHeight:
			cfg.Height = constants.MinBarcodeHeight
		case *cmd.Height > constants.MaxBarcodeHeight:
			cfg.Height = constants.MaxBarcodeHeight
		default:
			cfg.Height = barcode.Height(*cmd.Height)
		}
	} else {
		cfg.Height = constants.DefaultBarcodeHeight
	}

	hriPos := constants.DefaultBarcodeHriPosition
	if cmd.HRIPosition != nil {
		hriPos = constants.HriPosition(strings.ToLower(*cmd.HRIPosition))
	}
	cfg.HRIPosition = mapHRIPosition(hriPos)

	hriFont := constants.DefaultBarcodeHriFont
	if cmd.HRIFont != nil {
		hriFont = constants.Font(strings.ToLower(*cmd.HRIFont))
	}
	cfg.HRIFont = mapHRIFont(hriFont)

	if symbology == barcode.CODE128 || symbology == barcode.GS1128 {
		cfg.CodeSet = barcode.Code128SetB
		if barcode.ValidateNumericData([]byte(cmd.Data)) && len(cmd.Data)%2 == 0 {
			cfg.CodeSet = barcode.Code128SetC
		}
	}

	align := constants.DefaultBarcodeAlignment
	if cmd.Align != nil {
		align = constants.Alignment(strings.ToLower(*cmd.Align))
	}

	switch align {
	case constants.Center:
		if err := printer.AlignCenter(); err != nil {
			return err
		}
	case constants.Right:
		if err := printer.AlignRight(); err != nil {
			return err
		}
	case constants.Left:
		if err := printer.AlignLeft(); err != nil {
			return err
		}
	default:
		if err := printer.SetAlignment(constants.DefaultTextAlignment.String()); err != nil {
			return err
		}
	}

	if err := printer.PrintBarcode(cfg, []byte(cmd.Data)); err != nil {
		return fmt.Errorf("failed to print barcode: %w", err)
	}

	return printer.SetAlignment(constants.DefaultTextAlignment.String())
}

// mapHRIPosition maps schema HRI position constants to barcode package constants
func mapHRIPosition(pos constants.HriPosition) barcode.HRIPosition {
	posMap := map[constants.HriPosition]barcode.HRIPosition{
		constants.None:  barcode.HRINotPrinted,
		constants.Above: barcode.HRIAbove,
		constants.Below: barcode.HRIBelow,
		constants.Both:  barcode.HRIBoth,
	}

	hriPos, ok := posMap[pos]
	if !ok {
		return posMap[constants.DefaultBarcodeHriPosition]
	}

	return hriPos
}

// mapHRIFont maps schema font constants to barcode HRI font constants
func mapHRIFont(font constants.Font) barcode.HRIFont {
	fontMap := map[constants.Font]barcode.HRIFont{
		constants.A: barcode.HRIFontA,
		constants.B: barcode.HRIFontB,
	}

	hriFont, ok := fontMap[font]
	if !ok {
		return fontMap[constants.DefaultBarcodeHriFont]
	}

	return hriFont
}
