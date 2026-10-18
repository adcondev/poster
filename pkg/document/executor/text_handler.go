package executor

import (
	"encoding/json"
	"fmt"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/service"
)

type (
	// TextCommand for text handler
	TextCommand = schema.TextCommand
	// Content for text
	Content = schema.Content
	// Label for text
	Label = schema.Label
	// TextStyle for text formatting
	TextStyle = schema.TextStyle
)

func strPtr(s string) *string {
	return &s
}

// handleText manages text commands
func (e *Executor) handleText(printer service.PrinterActions, data json.RawMessage) error {
	var cmd TextCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse text command: %w", err)
	}

	if cmd.Label != nil && cmd.Label.Text != "" {
		if cmd.Label.Style != nil {
			if err := e.applyTextStyle(printer, cmd.Label.Style); err != nil {
				return fmt.Errorf("failed to apply label style: %w", err)
			}
		}

		labelText := cmd.Label.Text
		if cmd.Label.Separator != nil && *cmd.Label.Separator != "" {
			labelText += *cmd.Label.Separator
		} else {
			labelText += ": "
		}

		if err := e.applyAlign(printer, cmd.Label.Align); err != nil {
			return err
		}

		if !sameOrNil(cmd.Label.Align, cmd.Content.Align) {
			if err := printer.PrintLine(labelText); err != nil {
				return err
			}
			if cmd.Label.Style != nil {
				if err := e.resetTextStyle(printer, cmd.Label.Style); err != nil {
					return err
				}
			}
		} else {
			if err := printer.Print(labelText); err != nil {
				return err
			}
			if cmd.Label.Style != nil && cmd.Content.Style != nil {
				if err := e.resetDifferingStyles(printer, cmd.Label.Style, cmd.Content.Style); err != nil {
					return err
				}
			} else if cmd.Label.Style != nil {
				if err := e.resetTextStyle(printer, cmd.Label.Style); err != nil {
					return err
				}
			}
		}
	}

	if cmd.Content.Text != "" {
		if err := e.applyAlign(printer, cmd.Content.Align); err != nil {
			return err
		}

		if cmd.Content.Style != nil {
			if err := e.applyTextStyle(printer, cmd.Content.Style); err != nil {
				return fmt.Errorf("failed to apply content style: %w", err)
			}
		}

		newLine := true
		if cmd.NewLine != nil {
			newLine = *cmd.NewLine
		}

		if newLine {
			if err := printer.PrintLine(cmd.Content.Text); err != nil {
				return err
			}
		} else {
			if err := printer.Print(cmd.Content.Text); err != nil {
				return err
			}
		}

		if cmd.Content.Style != nil {
			if err := e.resetTextStyle(printer, cmd.Content.Style); err != nil {
				return fmt.Errorf("failed to reset content style: %w", err)
			}
		}
	}

	return e.applyAlign(printer, strPtr(constants.Left.String()))
}
