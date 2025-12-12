package executor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/service"
)

// SeparatorCommand for separator handler
type SeparatorCommand struct {
	Char   string `json:"char,omitempty"`
	Length int    `json:"length,omitempty"`
	// TODO: Add TextStyle if needed
	// TODO: Add Align if needed
	// TODO: Add Pre and Post feed line if needed
}

// FeedCommand for feed handler
type FeedCommand struct {
	Lines int `json:"lines"`
}

// CutCommand for cut handler
type CutCommand struct {
	Mode string `json:"mode,omitempty"`
	Feed int    `json:"feed,omitempty"`
}

// PulseCommand for cash drawer pulse handler
type PulseCommand struct {
	Pin     int `json:"pin,omitempty"`      // Drawer pin (0 or 1), default 0
	OnTime  int `json:"on_time,omitempty"`  // On time in ms, default 50
	OffTime int `json:"off_time,omitempty"` // Off time in ms, default 100
}

// BeepCommand for beep sound handler
type BeepCommand struct {
	Times int `json:"times,omitempty"` // Number of beeps, default 1
	Lapse int `json:"lapse,omitempty"` // Duration/interval factor, default 1
}

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

// handlePulse manages cash drawer pulse commands
func (e *Executor) handlePulse(printer *service.Printer, data json.RawMessage) error {
	var cmd PulseCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse pulse command: %w", err)
	}

	// Validate pin (must be 0 or 1)
	if cmd.Pin < 0 || cmd.Pin > 1 {
		log.Printf("Warning: invalid pulse pin %d, using 0", cmd.Pin)
		cmd.Pin = constants.DefaultPulsePin
	}

	// Set defaults
	if cmd.OnTime <= 0 {
		cmd.OnTime = constants.DefaultPulseOnTime // 50ms default
	}
	if cmd.OffTime <= 0 {
		cmd.OffTime = constants.DefaultPulseOffTime // 100ms default
	}

	// ESC p m t1 t2 - Generate pulse on connector pin
	// m = 0 or 1 (pin selection)
	// t1 = on time (t1 * 2 ms)
	// t2 = off time (t2 * 2 ms)
	bytes := []byte{0x1B, 0x70, byte(cmd.Pin), byte(cmd.OnTime / 2), byte(cmd.OffTime / 2)}
	return printer.Write(bytes)
}

// handleBeep manages beep sound commands
func (e *Executor) handleBeep(printer *service.Printer, data json.RawMessage) error {
	var cmd BeepCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse beep command: %w", err)
	}

	// Set defaults
	if cmd.Times <= 0 {
		cmd.Times = constants.DefaultBeepTimes
	}
	if cmd.Lapse <= 0 {
		cmd.Lapse = constants.DefaultBeepLapse
	}

	// ESC B n t - Beep command (non-standard, hardware-specific)
	bytes := []byte{0x1B, 0x42, byte(cmd.Times), byte(cmd.Lapse)}
	return printer.Write(bytes)
}
