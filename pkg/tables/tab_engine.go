// Package tables provides table generation and rendering for ESC/POS printers
package tables

import (
	"fmt"
	"io"
	"strings"

	"github.com/adcondev/pos-printer/pkg/commands/print"
	"github.com/adcondev/pos-printer/pkg/composer"
)

// Paper width constants (in characters)
const (
	Width58mm203dpi = 32 // Conservative for 58mm
	Width80mm203dpi = 48 // Conservative for 80mm
)

// Style represents text styling options
type Style struct {
	Bold       bool
	DoubleSize bool
	Underline  bool
}

// Options configures the table engine
type Options struct {
	PaperWidth    int   // Total width in characters
	ShowHeaders   bool  // Whether to show column headers
	HeaderStyle   Style // Style for headers
	WordWrap      bool  // Enable automatic word wrapping
	ColumnSpacing int   // Spaces between columns (default: 1)
}

// DefaultOptions returns sensible defaults for 80mm printers
func DefaultOptions() *Options {
	return &Options{
		PaperWidth:    Width80mm203dpi,
		ShowHeaders:   true,
		HeaderStyle:   Style{Bold: true},
		WordWrap:      true,
		ColumnSpacing: 1,
	}
}

// TabEngine handles table rendering
type TabEngine struct {
	definition *Definition
	options    *Options
}

// NewEngine creates a new table engine
func NewEngine(def *Definition, opts *Options) *TabEngine {
	if opts == nil {
		opts = DefaultOptions()
	}
	if opts.ColumnSpacing <= 0 {
		opts.ColumnSpacing = 1
	}
	if def.PaperWidth > 0 {
		opts.PaperWidth = def.PaperWidth
	}

	return &TabEngine{
		definition: def,
		options:    opts,
	}
}

// Render renders the table data to the writer
func (te *TabEngine) Render(w io.Writer, data *Data) error {
	if data == nil {
		return fmt.Errorf("table data cannot be nil")
	}
	if err := data.validate(); err != nil {
		return fmt.Errorf("invalid table data: %w", err)
	}

	def := te.definition
	if len(data.Definition.Columns) > 0 {
		def = &data.Definition
	}

	// Headers
	if te.options.ShowHeaders || data.ShowHeaders {
		headerLine := te.formatHeaderRow(te.makeHeaderRow(def), def)
		if _, err := w.Write([]byte(headerLine + string(print.LF))); err != nil {
			return err
		}

	}

	// Data rows (without blank lines between them)
	for _, row := range data.Rows {
		if te.options.WordWrap {
			wrapped := te.wrapRow(row, def)
			for _, wr := range wrapped {
				line := te.formatRow(wr, def)
				if _, err := w.Write([]byte(line + string(print.LF))); err != nil {
					return err
				}
			}
		} else {
			line := te.formatRow(row, def)
			if _, err := w.Write([]byte(line + string(print.LF))); err != nil {
				return err
			}
		}
	}

	return nil
}

// formatHeaderRow formats a header row with bold styling
func (te *TabEngine) formatHeaderRow(cells []string, def *Definition) string {
	var result strings.Builder

	cmds := composer.NewEscpos()

	// Apply bold command at the beginning if enabled
	if te.options.HeaderStyle.Bold {
		result.WriteString(string(cmds.EnableBold())) // ESC E 1 (Bold ON)
	}

	// Format header cells
	result.WriteString(te.formatRow(cells, def))

	// Reset bold at the end if it was enabled
	if te.options.HeaderStyle.Bold {
		result.WriteString(string(cmds.DisableBold())) // ESC E 0 (Bold OFF)
	}

	return result.String()
}

// TODO: Consider row styles in the future

// formatRow formats a regular data row without styling
func (te *TabEngine) formatRow(cells []string, def *Definition) string {
	var result strings.Builder

	for i, cell := range cells {
		if i < len(def.Columns) {
			padded := padString(cell, def.Columns[i].Width, def.Columns[i].Align)
			result.WriteString(padded)

			// Add spacing between columns
			if i < len(cells)-1 {
				result.WriteString(strings.Repeat(" ", te.options.ColumnSpacing))
			}
		}
	}

	return result.String()
}

// wrapRow handles word wrapping for a single row
func (te *TabEngine) wrapRow(row Row, def *Definition) [][]string {
	wrappedCells := make([][]string, len(row))
	maxLines := 0

	for i, cell := range row {
		if i < len(def.Columns) {
			wrapped := wrapText(cell, def.Columns[i].Width)
			wrappedCells[i] = wrapped
			if len(wrapped) > maxLines {
				maxLines = len(wrapped)
			}
		}
	}

	result := make([][]string, maxLines)
	for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
		result[lineIdx] = make([]string, len(row))
		for colIdx := range row {
			if colIdx < len(def.Columns) && lineIdx < len(wrappedCells[colIdx]) {
				result[lineIdx][colIdx] = wrappedCells[colIdx][lineIdx]
			} else {
				// Empty string for missing cells
				result[lineIdx][colIdx] = ""
			}
		}
	}

	return result
}

// makeHeaderRow creates header row from column definitions
func (te *TabEngine) makeHeaderRow(def *Definition) []string {
	headers := make([]string, len(def.Columns))
	for i, col := range def.Columns {
		headers[i] = col.Header
	}
	return headers
}
