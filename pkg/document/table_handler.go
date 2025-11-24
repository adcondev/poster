package document

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/adcondev/pos-printer/pkg/service"
	"github.com/adcondev/pos-printer/pkg/tables"
)

// TODO: Consider upper_separator y lower_separator for tables

// TableCommand represents a table command in the document
type TableCommand struct {
	Definition  tables.Definition `json:"definition"`
	ShowHeaders bool              `json:"show_headers,omitempty"`
	Rows        [][]string        `json:"rows"`
	Options     *TableOptions     `json:"options,omitempty"`
}

// TODO: Implementar Header con TextStyle sin alineación

// TableOptions configures table rendering options
type TableOptions struct {
	// HeaderBold enables bold styling for table headers
	HeaderBold bool `json:"header_bold,omitempty"`
	// WordWrap enables automatic text wrapping in cells
	WordWrap bool `json:"word_wrap,omitempty"`
	// ColumnSpacing sets the number of spaces between columns (default: 1)
	ColumnSpacing int `json:"column_spacing,omitempty"`
	// Align sets the default alignment for table content (left, center, right)
	Align string `json:"align,omitempty"`
}

// TODO: Manage text_under and text_above options instead of human_text

// TODO: Consider a title fields for tables as field in TableCommand

// handleTable manages table commands
func (e *Executor) handleTable(printer *service.Printer, data json.RawMessage) error {
	var cmd TableCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return fmt.Errorf("failed to parse table command: %w", err)
	}

	if cmd.Options != nil && cmd.Options.ColumnSpacing < 0 {
		// Validar que ColumnSpacing no sea negativo
		cmd.Options.ColumnSpacing = 0
		log.Printf("ColumnSpacing cannot be negative, using 0")
	}
	// Validate table command
	if len(cmd.Definition.Columns) == 0 {
		return fmt.Errorf("table must have at least one column defined")
	}

	// Create table options with defaults
	opts := &tables.Options{
		ShowHeaders:   cmd.ShowHeaders,
		WordWrap:      true,
		ColumnSpacing: 1,
		HeaderStyle:   tables.Style{Bold: true},
	}

	// Apply custom options if provided
	if cmd.Options != nil {
		opts.WordWrap = cmd.Options.WordWrap
		if cmd.Options.HeaderBold {
			opts.HeaderStyle.Bold = true
		}
		if cmd.Options.ColumnSpacing > 0 {
			opts.ColumnSpacing = cmd.Options.ColumnSpacing
		}
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

	// Create table engine
	engine := tables.NewEngine(&cmd.Definition, opts)

	// Prepare table data
	tableData := &tables.Data{
		Definition:  cmd.Definition,
		ShowHeaders: cmd.ShowHeaders,
		Rows:        make([]tables.Row, len(cmd.Rows)),
	}

	// Convert rows
	for i, row := range cmd.Rows {
		tableData.Rows[i] = row
	}

	// Render table to string
	var buf strings.Builder
	if err := engine.Render(&buf, tableData); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	// Aplicar alineación
	align := "left"
	if cmd.Options != nil && cmd.Options.Align != "" {
		align = cmd.Options.Align
	}
	switch strings.ToLower(align) {
	case center:
		err := printer.AlignCenter()
		if err != nil {
			return err
		}
	case right:
		err := printer.AlignRight()
		if err != nil {
			return err
		}
	default:
		err := printer.AlignLeft()
		if err != nil {
			return err
		}
	}

	err := printer.Print(buf.String())
	if err != nil {
		return err
	}

	// Restaurar alineación
	err = printer.AlignLeft()
	if err != nil {
		return err
	}

	// Send the raw output (includes ESC/POS commands for bold)
	return nil
}
