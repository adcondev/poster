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
	// TableCommand for table handler
	TableCommand = schema.TableCommand
	// TableOptions for table configuration
	TableOptions = schema.TableOptions
)

// handleTable manages table commands
func (e *Executor) handleTable(printer service.PrinterActions, data json.RawMessage) error {
	var cmd TableCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse table command: %w", err)
	}

	if cmd.Options != nil && cmd.Options.ColumnSpacing != nil && *cmd.Options.ColumnSpacing < 0 {
		zero := 0
		cmd.Options.ColumnSpacing = &zero
	}

	if len(cmd.Definition.Columns) == 0 {
		return fmt.Errorf("table must have at least one column defined")
	}

	totalColumnWidth, err := ValidateColumns(cmd.Definition.Columns)
	if err != nil {
		return fmt.Errorf("invalid table definition: %w", err)
	}

	isFontB := cmd.Options != nil && strings.EqualFold(cmd.Options.Font, "B")
	maxChars, paperWidth, dpi := calculateMaxChars(printer.GetProfile(), isFontB)

	spacing := constants.DefaultTableColumnSpacing
	if cmd.Options != nil && cmd.Options.ColumnSpacing != nil {
		spacing = *cmd.Options.ColumnSpacing
	}

	if err := fitTableColumns(&cmd, totalColumnWidth, maxChars, spacing, isFontB, paperWidth, dpi); err != nil {
		return err
	}

	opts := resolveTableOptions(&cmd, maxChars, spacing)

	if isFontB {
		if err := printer.FontB(); err != nil {
			return fmt.Errorf("failed to set Font B for table: %w", err)
		}
		defer func() {
			_ = printer.FontA()
		}()
	} else {
		if err := printer.FontA(); err != nil {
			return fmt.Errorf("failed to set Font A for table: %w", err)
		}
	}

	rendered, err := RenderTableStructured(&cmd, opts)
	if err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	align := constants.DefaultTableAlignment.String()
	if cmd.Options != nil && cmd.Options.Align != "" {
		align = cmd.Options.Align
	}

	if err := ApplyAlignment(printer, align); err != nil {
		return err
	}

	if err := printRenderedTable(printer, rendered); err != nil {
		return err
	}

	return printer.AlignLeft()
}
