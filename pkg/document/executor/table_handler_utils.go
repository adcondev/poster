package executor

import (
	"fmt"
	"strings"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/profile"
	"github.com/adcondev/poster/pkg/service"
	"github.com/adcondev/poster/pkg/tables"
)

// ============================================================================
// Table Handler Utilities
// ============================================================================

// ValidateColumns checks that all columns have valid widths and returns total width.
func ValidateColumns(columns []tables.Column) (int, error) {
	total := 0
	for _, col := range columns {
		if col.Width <= 0 {
			return 0, fmt.Errorf("column '%s' has invalid width: %d (must be > 0)",
				col.Name, col.Width)
		}
		total += col.Width
	}
	return total, nil
}

// CalculateTotalWidth computes total required width including gaps.
func CalculateTotalWidth(columnWidth, columnCount, spacing int) int {
	gaps := columnCount - 1
	if gaps < 0 {
		gaps = 0
	}
	return columnWidth + (gaps * spacing)
}

// RenderTable renders the table to a string.
func RenderTable(cmd *TableCommand, opts *tables.Options) (string, error) {
	engine := tables.NewEngine(&cmd.Definition, opts)

	showHeaders := constants.DefaultTableShowHeaders
	if cmd.ShowHeaders != nil {
		showHeaders = *cmd.ShowHeaders
	}
	if opts != nil {
		showHeaders = opts.ShowHeaders
	}

	tableData := &tables.Data{
		Definition:  cmd.Definition,
		ShowHeaders: showHeaders,
		Rows:        make([]tables.Row, len(cmd.Rows)),
	}

	for i, row := range cmd.Rows {
		tableData.Rows[i] = row
	}

	var buf strings.Builder
	if err := engine.Render(&buf, tableData); err != nil {
		return "", fmt.Errorf("failed to render table: %w", err)
	}

	return buf.String(), nil
}

// RenderTableStructured renders the table to structured header and body lines.
func RenderTableStructured(cmd *TableCommand, opts *tables.Options) (*tables.RenderedTable, error) {
	engine := tables.NewEngine(&cmd.Definition, opts)

	showHeaders := constants.DefaultTableShowHeaders
	if cmd.ShowHeaders != nil {
		showHeaders = *cmd.ShowHeaders
	}
	if opts != nil {
		showHeaders = opts.ShowHeaders
	}

	tableData := &tables.Data{
		Definition:  cmd.Definition,
		ShowHeaders: showHeaders,
		Rows:        make([]tables.Row, len(cmd.Rows)),
	}

	for i, row := range cmd.Rows {
		tableData.Rows[i] = row
	}

	rendered, err := engine.RenderStructured(tableData)
	if err != nil {
		return nil, fmt.Errorf("failed to render structured table: %w", err)
	}

	return rendered, nil
}

// calculateMaxChars resolves maximum characters per line, paper width and DPI based on printer profile and font.
func calculateMaxChars(prof *profile.Escpos, isFontB bool) (int, float64, int) {
	var dotsPerLine int
	var paperWidth float64
	var dpi int
	if prof != nil {
		dotsPerLine = prof.DotsPerLine
		paperWidth = prof.PaperWidth
		dpi = prof.DPI
	}

	if isFontB {
		maxChars := constants.MaxCharsForPaperFontB(dotsPerLine)
		if maxChars == 0 {
			if paperWidth >= 80 {
				return 64, paperWidth, dpi
			}
			return 42, paperWidth, dpi
		}
		return maxChars, paperWidth, dpi
	}

	maxChars := constants.MaxCharsForPaperFontA(dotsPerLine)
	if maxChars == 0 {
		if paperWidth >= 80 {
			return tables.Width80mm203dpi, paperWidth, dpi
		}
		return tables.Width58mm203dpi, paperWidth, dpi
	}
	return maxChars, paperWidth, dpi
}

// fitTableColumns validates table width against maxChars and applies auto-reduction if needed.
func fitTableColumns(cmd *TableCommand, totalColumnWidth, maxChars, spacing int, isFontB bool, paperWidth float64, dpi int) error {
	numberOfGaps := len(cmd.Definition.Columns) - 1
	if numberOfGaps < 0 {
		numberOfGaps = 0
	}

	totalGapWidth := numberOfGaps * spacing
	totalRequiredWidth := totalColumnWidth + totalGapWidth

	if totalRequiredWidth <= maxChars {
		return nil
	}

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
			return nil
		}
		return fmt.Errorf(
			"table overflow: required %d chars > max %d; auto-reduction failed: %w",
			totalRequiredWidth, maxChars, err,
		)
	}

	fontName := "Font A"
	if isFontB {
		fontName = "Font B"
	}
	return fmt.Errorf(
		"table overflow: columns (%d) + gaps (%d) = %d chars, exceeds max %d chars "+
			"(%.0fmm paper @ %d DPI, %s)",
		totalColumnWidth,
		totalGapWidth,
		totalRequiredWidth,
		maxChars,
		paperWidth,
		dpi,
		fontName,
	)
}

// resolveTableOptions constructs the tables.Options from TableCommand.
func resolveTableOptions(cmd *TableCommand, maxChars, spacing int) *tables.Options {
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

	if cmd.Definition.PaperWidth > 0 {
		opts.PaperWidth = cmd.Definition.PaperWidth
	}
	return opts
}

// printRenderedTable outputs the rendered table lines with bold header control.
func printRenderedTable(printer service.PrinterActions, rendered *tables.RenderedTable) error {
	if len(rendered.HeaderLines) > 0 {
		if rendered.HeaderBold {
			if err := printer.EnableBold(); err != nil {
				return err
			}
		}
		if len(rendered.BodyLines) > 0 {
			if err := printer.PrintLine(strings.Join(rendered.HeaderLines, "\n")); err != nil {
				return err
			}
		} else {
			if err := printer.Print(strings.Join(rendered.HeaderLines, "\n")); err != nil {
				return err
			}
		}
		if rendered.HeaderBold {
			if err := printer.DisableBold(); err != nil {
				return err
			}
		}
	}

	if len(rendered.BodyLines) > 0 {
		if err := printer.Print(strings.Join(rendered.BodyLines, "\n")); err != nil {
			return err
		}
	}

	return nil
}


