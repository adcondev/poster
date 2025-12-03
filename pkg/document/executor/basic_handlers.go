package executor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adcondev/pos-printer/pkg/constants"
	"github.com/adcondev/pos-printer/pkg/service"
)

// handleSeparator manages separator commands
func (e *Executor) handleSeparator(printer *service.Printer, data json.RawMessage) error {
	var cmd SeparatorCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse separator command: %w", err)
	}

	// Valores por defecto
	if cmd.Char == "" {
		cmd.Char = constants.DefaultSeparatorChar
	}
	if cmd.Length <= 0 {
		// Usar ancho del papel en caracteres (aproximado)
		// TODO: Verify the following line for different fonts, constrained to Font A
		cmd.Length = e.printer.Profile.DotsPerLine / 12 // Aproximación para Font A
	}

	err := printer.SetAlignment(constants.DefaultSeparatorAlignment.String())
	if err != nil {
		return err
	}

	// Construir línea separadora
	line := strings.Repeat(cmd.Char, cmd.Length)
	err = printer.PrintLine(line[:cmd.Length-1]) // Trim to exact length
	if err != nil {
		return err
	}

	err = printer.AlignLeft()
	if err != nil {
		return err
	}

	return nil
}

// handleFeed manages feed commands
func (e *Executor) handleFeed(printer *service.Printer, data json.RawMessage) error {
	var cmd FeedCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse feed command: %w", err)
	}

	if cmd.Lines <= 0 {
		cmd.Lines = constants.DefaultFeedLines
	}

	return printer.FeedLines(byte(cmd.Lines))
}

// handleCut manages cut commands
func (e *Executor) handleCut(printer *service.Printer, data json.RawMessage) error {
	var cmd CutCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse cut command: %w", err)
	}

	// Default feed antes del corte según schema
	if cmd.Feed == 0 {
		cmd.Feed = constants.DefaultCutFeed
	}

	// Avance antes del corte si se especifica
	if cmd.Feed > 0 {
		if err := printer.FeedLines(byte(cmd.Feed)); err != nil {
			return err
		}
	}

	// TODO: Implement field Dots, as this command make the feed based lines only.
	// Instead of 0 for Full and Partial cut, it should be the number of dots to feed before cutting.
	// Maybe calculation from mm to dpl is needed.

	// Ejecutar corte
	switch strings.ToLower(cmd.Mode) {
	case constants.Full.String():
		return printer.FullFeedAndCut(0)
	case constants.Partial.String():
		return printer.PartialFeedAndCut(0)
	default: // partial
		return printer.PartialFeedAndCut(0)
	}
}
