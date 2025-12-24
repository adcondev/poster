package executor

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/service"
	"github.com/adcondev/poster/pkg/tables"
)

// TODO: Consider upper_separator y lower_separator for tables

// TODO: Add TextStyle support for column names without alignmentón

// TODO: Manage text_under and text_above options instead of human_text

// TODO: Consider a title fields for tables as field in TableCommand

// TableCommand for table handler
type TableCommand struct {
	Definition  tables.Definition `json:"definition"`
	ShowHeaders bool              `json:"show_headers,omitempty"`
	Rows        [][]string        `json:"rows"`
	Options     *TableOptions     `json:"options,omitempty"`
}

// TableOptions for table configuration
type TableOptions struct {
	HeaderBold    bool   `json:"header_bold,omitempty"`
	WordWrap      bool   `json:"word_wrap,omitempty"`
	ColumnSpacing int    `json:"column_spacing,omitempty"`
	Align         string `json:"align,omitempty"`
	// TODO: Setup Font usage
	Font       string `json:"font,omitempty"`
	AutoReduce *bool  `json:"auto_reduce,omitempty"`
}

// TODO: Reducir código duplicado entre handlers

// handleTable manages table commands
func (e *Executor) handleTable(printer *service.Printer, data json.RawMessage) error {
	var cmd TableCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse table command: %w", err)
	}

	// Normalize options
	if cmd.Options != nil && cmd.Options.ColumnSpacing < 0 {
		cmd.Options.ColumnSpacing = 0
		log.Printf("ColumnSpacing cannot be negative, using 0")
	}

	// Validate columns exist
	if len(cmd.Definition.Columns) == 0 {
		return fmt.Errorf("table must have at least one column defined")
	}

	// Sum all column widths
	totalColumnWidth, err := ValidateColumns(cmd.Definition.Columns)
	if err != nil {
		return fmt.Errorf("invalid table definition: %w", err)
	}

	// Calculate max chars based on printer profile and Font A
	maxChars := constants.MaxCharsForPaperFontA(printer.Profile.DotsPerLine)

	// Fallback for incomplete profiles (e.g., mock printers in tests)
	if maxChars == 0 {
		if printer.Profile.PaperWidth >= 80 {
			maxChars = tables.Width80mm203dpi // 48 chars
		} else {
			maxChars = tables.Width58mm203dpi // 32 chars
		}
		log.Printf("DotsPerLine not set, falling back to %d chars based on %.0fmm paper",
			maxChars, printer.Profile.PaperWidth)
	}

	spacing := constants.DefaultTableColumnSpacing
	if cmd.Options != nil && cmd.Options.ColumnSpacing > 0 {
		spacing = cmd.Options.ColumnSpacing
	}

	// Calculate gap width (spaces between columns)
	numberOfGaps := len(cmd.Definition.Columns) - 1
	if numberOfGaps < 0 {
		numberOfGaps = 0
	}

	// Gaps are on both sides of each column
	totalGapWidth := numberOfGaps * spacing
	totalRequiredWidth := totalColumnWidth + totalGapWidth

	if totalRequiredWidth > maxChars {
		// Check if AutoReduce is enabled (default is true)
		shouldReduce := constants.DefaultTableAutoReduce
		if cmd.Options != nil && cmd.Options.AutoReduce != nil {
			shouldReduce = *cmd.Options.AutoReduce
		}

		if shouldReduce {
			log.Printf("Table overflow detected (%d > %d). Attempting auto-reduction...",
				totalRequiredWidth, maxChars)

			result, err := tables.ReduceToFit(
				cmd.Definition.Columns,
				maxChars,
				spacing,
				constants.MinTableColumnWidth,
			)

			if err == nil && result.Success {
				// Update the definition with reduced columns
				cmd.Definition.Columns = result.Columns
				log.Printf("Table auto-reduced:  %d → %d chars (%d reductions applied)",
					result.OriginalWidth, result.ReducedWidth, result.Reductions)
			} else {
				// Auto-reduction failed
				return fmt.Errorf(
					"table overflow: required %d chars > max %d; auto-reduction failed:  %w",
					totalRequiredWidth, maxChars, err,
				)
			}
		} else {
			// AutoReduce disabled by user
			return fmt.Errorf(
				"table overflow: columns (%d) + gaps (%d) = %d chars, exceeds max %d chars "+
					"(%.0fmm paper @ %d DPI, Font A)",
				totalColumnWidth,
				totalGapWidth,
				totalRequiredWidth,
				maxChars,
				printer.Profile.PaperWidth,
				printer.Profile.DPI,
			)
		}
	}

	// Create options with validated paper width
	opts := &tables.Options{
		PaperWidth:    maxChars,
		ShowHeaders:   constants.DefaultTableShowHeaders,
		WordWrap:      constants.DefaultTableWordWrap,
		ColumnSpacing: spacing,
		HeaderStyle:   tables.Style{Bold: constants.DefaultTableHeaderBold},
	}

	// Apply custom options if provided
	if cmd.Options != nil {
		opts.WordWrap = cmd.Options.WordWrap
		if cmd.Options.HeaderBold {
			opts.HeaderStyle.Bold = true
		}
	}

	// Enforce Font A for consistent table rendering
	if err := printer.FontA(); err != nil {
		return fmt.Errorf("failed to set Font A for table:  %w", err)
	}

	// Set paper width
	switch {
	case cmd.Definition.PaperWidth > 0:
		opts.PaperWidth = cmd.Definition.PaperWidth
	case printer.Profile.PrintWidth > 0:
		opts.PaperWidth = printer.Profile.PrintWidth
	default:
		if printer.Profile.PaperWidth >= 80 {
			opts.PaperWidth = tables.Width80mm203dpi
		} else {
			opts.PaperWidth = tables.Width58mm203dpi
		}
	}

	finalTable, err := RenderTable(&cmd, opts)
	if err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	// Apply alignment
	align := constants.DefaultTableAlignment.String()
	if cmd.Options != nil && cmd.Options.Align != "" {
		align = cmd.Options.Align
	}

	err = ApplyAlignment(printer, align)
	if err != nil {
		return err
	}

	err = printer.Print(finalTable)
	if err != nil {
		return err
	}

	// Restaurar alineación
	err = printer.AlignLeft()
	if err != nil {
		return err
	}

	// Send the raw output
	return nil
}
