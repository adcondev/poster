// Package tables provides table generation and rendering for ESC/POS printers
package tables

import (
	"fmt"
	"io"
	"strings"

	"github.com/adcondev/poster/pkg/commands/print"
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

// RenderedTable holds structured header and body lines without embedded ESC/POS control bytes.
type RenderedTable struct {
	HeaderLines []string
	BodyLines   []string
	HeaderBold  bool
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
	if opts.ColumnSpacing < 0 {
		opts.ColumnSpacing = 1
	}

	return &TabEngine{
		definition: def,
		options:    opts,
	}
}

// RenderStructured produces structured pure text lines for headers and body without ESC/POS control bytes.
func (te *TabEngine) RenderStructured(data *Data) (*RenderedTable, error) {
	if data == nil {
		return nil, fmt.Errorf("table data cannot be nil")
	}
	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("invalid table data: %w", err)
	}

	def := te.definition
	if len(data.Definition.Columns) > 0 {
		def = &data.Definition
	}
	if def == nil || len(def.Columns) == 0 {
		return nil, fmt.Errorf("table must have at least one column")
	}

	var headerLines []string
	if te.options.ShowHeaders || data.ShowHeaders {
		headerRow := te.makeHeaderRow(def)
		if te.options.WordWrap {
			wrapped := te.wrapRow(headerRow, def)
			for _, wr := range wrapped {
				headerLines = append(headerLines, te.formatRow(wr, def))
			}
		} else {
			headerLines = append(headerLines, te.formatRow(headerRow, def))
		}
	}

	var bodyLines []string
	for _, row := range data.Rows {
		if te.options.WordWrap {
			wrapped := te.wrapRow(row, def)
			for _, wr := range wrapped {
				bodyLines = append(bodyLines, te.formatRow(wr, def))
			}
		} else {
			bodyLines = append(bodyLines, te.formatRow(row, def))
		}
	}

	return &RenderedTable{
		HeaderLines: headerLines,
		BodyLines:   bodyLines,
		HeaderBold:  te.options.HeaderStyle.Bold,
	}, nil
}

// Render renders the table data to the writer as pure text lines.
func (te *TabEngine) Render(w io.Writer, data *Data) error {
	rendered, err := te.RenderStructured(data)
	if err != nil {
		return err
	}

	for _, line := range rendered.HeaderLines {
		if _, err := w.Write([]byte(line + string(print.LF))); err != nil {
			return err
		}
	}

	for _, line := range rendered.BodyLines {
		if _, err := w.Write([]byte(line + string(print.LF))); err != nil {
			return err
		}
	}

	return nil
}

// formatRow formats a row into a single aligned line
func (te *TabEngine) formatRow(cells []string, def *Definition) string {
	var result strings.Builder

	for i, cell := range cells {
		if i < len(def.Columns) {
			padded := PadString(cell, def.Columns[i].Width, def.Columns[i].Align)
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
			wrapped := WrapText(cell, def.Columns[i].Width)
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
		headers[i] = col.Name
	}
	return headers
}
