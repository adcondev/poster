package executor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/service"
)

type (
	// SeparatorCommand for separator handler
	SeparatorCommand = schema.SeparatorCommand
	// FeedCommand for feed handler
	FeedCommand = schema.FeedCommand
	// CutCommand for cut handler
	CutCommand = schema.CutCommand
	// PulseCommand for cash drawer pulse handler
	PulseCommand = schema.PulseCommand
	// BeepCommand for beep sound handler
	BeepCommand = schema.BeepCommand
)

// handleSeparator manages separator commands
func (e *Executor) handleSeparator(printer service.PrinterActions, data json.RawMessage) error {
	var cmd SeparatorCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse separator command: %w", err)
	}

	if cmd.Char == "" {
		cmd.Char = constants.DefaultSeparatorChar
	}
	if cmd.Length <= 0 {
		prof := printer.GetProfile()
		if prof != nil && prof.DotsPerLine > 0 {
			cmd.Length = prof.DotsPerLine / constants.FontAWidth
		} else {
			cmd.Length = constants.DefaultSeparatorLength
		}
	}

	if err := printer.SetAlignment(constants.DefaultSeparatorAlignment.String()); err != nil {
		return err
	}

	line := strings.Repeat(cmd.Char, cmd.Length)
	runes := []rune(line)
	if len(runes) > cmd.Length {
		line = string(runes[:cmd.Length])
	}
	if err := printer.PrintLine(line); err != nil {
		return err
	}

	return printer.AlignLeft()
}

// handleFeed manages feed commands
func (e *Executor) handleFeed(printer service.PrinterActions, data json.RawMessage) error {
	var cmd FeedCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse feed command: %w", err)
	}

	if cmd.Lines <= 0 {
		cmd.Lines = constants.DefaultFeedLines
	} else if cmd.Lines > 255 {
		cmd.Lines = 255
	}

	return printer.FeedLines(byte(cmd.Lines)) //nolint:gosec // bounded to [1, 255]
}

// handleCut manages cut commands
func (e *Executor) handleCut(printer service.PrinterActions, data json.RawMessage) error {
	var cmd CutCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse cut command: %w", err)
	}

	feed := constants.DefaultCutFeed
	if cmd.Feed != nil {
		feed = *cmd.Feed
	}
	if feed < 0 {
		feed = 0
	} else if feed > 255 {
		feed = 255
	}

	if feed > 0 {
		if err := printer.FeedLines(byte(feed)); err != nil {
			return err
		}
	}

	switch strings.ToLower(cmd.Mode) {
	case constants.Full.String():
		return printer.FullFeedAndCut(0)
	case constants.Partial.String():
		return printer.PartialFeedAndCut(0)
	default:
		return printer.PartialFeedAndCut(0)
	}
}

// handlePulse manages cash drawer pulse commands
func (e *Executor) handlePulse(printer service.PrinterActions, data json.RawMessage) error {
	var cmd PulseCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse pulse command: %w", err)
	}

	if cmd.Pin < 0 || cmd.Pin > 1 {
		cmd.Pin = constants.DefaultPulsePin
	}

	if cmd.OnTime <= 0 {
		cmd.OnTime = constants.DefaultPulseOnTime
	} else if cmd.OnTime > 510 {
		cmd.OnTime = 510
	}
	if cmd.OffTime <= 0 {
		cmd.OffTime = constants.DefaultPulseOffTime
	} else if cmd.OffTime > 510 {
		cmd.OffTime = 510
	}

	return printer.Pulse(byte(cmd.Pin), byte(cmd.OnTime/2), byte(cmd.OffTime/2)) //nolint:gosec // bounded above
}

// handleBeep manages beep sound commands
func (e *Executor) handleBeep(printer service.PrinterActions, data json.RawMessage) error {
	var cmd BeepCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse beep command: %w", err)
	}

	if cmd.Times <= 0 {
		cmd.Times = constants.DefaultBeepTimes
	} else if cmd.Times > 255 {
		cmd.Times = 255
	}
	if cmd.Lapse <= 0 {
		cmd.Lapse = constants.DefaultBeepLapse
	} else if cmd.Lapse > 255 {
		cmd.Lapse = 255
	}

	return printer.Beep(byte(cmd.Times), byte(cmd.Lapse)) //nolint:gosec // bounded to [1, 255]
}
