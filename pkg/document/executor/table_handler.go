package executor

import (
	"encoding/json"
	"fmt"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/service"
	"github.com/adcondev/poster/pkg/tables"
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

	prof := printer.GetProfile()
	var dotsPerLine int
	var paperWidth float64
	var dpi int
	if prof != nil {
		dotsPerLine = prof.DotsPerLine
		paperWidth = prof.PaperWidth
		dpi = prof.DPI
	}

	maxChars := constants.MaxCharsForPaperFontA(dotsPerLine)
	if maxChars == 0 {
		if paperWidth >= 80 {
			maxChars = tables.Width80mm203dpi
		} else {
			maxChars = tables.Width58mm203dpi
		}
	}

	spacing := constants.DefaultTableColumnSpacing
	if cmd.Options != nil && cmd.Options.ColumnSpacing != nil {
		spacing = *cmd.Options.ColumnSpacing
	}

	numberOfGaps := len(cmd.Definition.Columns) - 1
	if numberOfGaps < 0 {
		numberOfGaps = 0
	}

	totalGapWidth := numberOfGaps * spacing
	totalRequiredWidth := totalColumnWidth + totalGapWidth

	if totalRequiredWidth > maxChars {
		shouldReduce := constants.DefaultTableAutoReduce
		if cmd.Options != nil && cmd.Options.AutoReduce != nil {
			shouldReduce = *cmd.Options.AutoReduce
		}

		if shouldReduce {
			result, err := tables.ReduceToFit(
				cmd.Definition.Columns,
				maxChars,
				spacing,
				constants.MinTableColumnWidth,
			)

			if err == nil && result.Success {
				cmd.Definition.Columns = result.Columns
			} else {
				return fmt.Errorf(
					"table overflow: required %d chars > max %d; auto-reduction failed: %w",
					totalRequiredWidth, maxChars, err,
				)
			}
		} else {
			return fmt.Errorf(
				"table overflow: columns (%d) + gaps (%d) = %d chars, exceeds max %d chars "+
					"(%.0fmm paper @ %d DPI, Font A)",
				totalColumnWidth,
				totalGapWidth,
				totalRequiredWidth,
				maxChars,
				paperWidth,
				dpi,
			)
		}
	}

	showHeaders := constants.DefaultTableShowHeaders
	if cmd.ShowHeaders != nil {
		showHeaders = *cmd.ShowHeaders
	}

	wordWrap := constants.DefaultTableWordWrap
	headerBold := constants.DefaultTableHeaderBold
	if cmd.Options != nil {
		if cmd.Options.WordWrap != nil {
			wordWrap = *cmd.Options.WordWrap
		}
		if cmd.Options.HeaderBold != nil {
			headerBold = *cmd.Options.HeaderBold
		}
	}

	opts := &tables.Options{
		PaperWidth:    maxChars,
		ShowHeaders:   showHeaders,
		WordWrap:      wordWrap,
		ColumnSpacing: spacing,
		HeaderStyle:   tables.Style{Bold: headerBold},
	}

	if err := printer.FontA(); err != nil {
		return fmt.Errorf("failed to set Font A for table: %w", err)
	}

	if cmd.Definition.PaperWidth > 0 {
		opts.PaperWidth = cmd.Definition.PaperWidth
	}

	finalTable, err := RenderTable(&cmd, opts)
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

	if err := printer.Print(finalTable); err != nil {
		return err
	}

	return printer.AlignLeft()
}
