package document

import (
	"encoding/json"
	"log"

	"github.com/adcondev/pos-printer/pkg/tables"
)

// TableBuilder provides a fluent API for building table commands within a document
type TableBuilder struct {
	parent      *Builder
	definition  tables.Definition
	rows        [][]string
	showHeaders bool
	options     *TableOptions
}

// NewTableBuilder creates a new table builder attached to a document builder
func (b *Builder) NewTableBuilder() *TableBuilder {
	return &TableBuilder{
		parent: b,
		definition: tables.Definition{
			Columns: []tables.Column{},
		},
		rows:        [][]string{},
		showHeaders: true,
		options: &TableOptions{
			HeaderBold: true,
			WordWrap:   true,
		},
	}
}

// AddColumn adds a column to the table
func (tb *TableBuilder) AddColumn(header string, width int, align tables.Alignment) *TableBuilder {
	tb.definition.Columns = append(tb.definition.Columns, tables.Column{
		Header: header,
		Width:  width,
		Align:  align,
	})
	return tb
}

// SetPaperWidth sets the paper width in characters
func (tb *TableBuilder) SetPaperWidth(width int) *TableBuilder {
	tb.definition.PaperWidth = width
	return tb
}

// AddRow adds a data row to the table
func (tb *TableBuilder) AddRow(cells ...string) *TableBuilder {
	tb.rows = append(tb.rows, cells)
	return tb
}

// AddRows adds multiple rows at once
func (tb *TableBuilder) AddRows(rows [][]string) *TableBuilder {
	tb.rows = append(tb.rows, rows...)
	return tb
}

// ShowHeaders enables or disables header display
func (tb *TableBuilder) ShowHeaders(show bool) *TableBuilder {
	tb.showHeaders = show
	return tb
}

// HeaderBold enables or disables bold headers
func (tb *TableBuilder) HeaderBold(bold bool) *TableBuilder {
	tb.options.HeaderBold = bold
	return tb
}

// WordWrap enables or disables word wrapping
func (tb *TableBuilder) WordWrap(wrap bool) *TableBuilder {
	tb.options.WordWrap = wrap
	return tb
}

// ColumnSpacing sets the spacing between columns
func (tb *TableBuilder) ColumnSpacing(spacing int) *TableBuilder {
	if spacing < 0 {
		spacing = 0
	}
	tb.options.ColumnSpacing = spacing
	return tb
}

// Align sets the table alignment (left, center, right)
func (tb *TableBuilder) Align(align string) *TableBuilder {
	tb.options.Align = align
	return tb
}

// Done finishes the table and adds it to the parent document builder
func (tb *TableBuilder) Done() *Builder {
	if len(tb.definition.Columns) == 0 {
		log.Printf("Warning: table has no columns defined, skipping")
		return tb.parent
	}

	cmd := TableCommand{
		Definition:  tb.definition,
		ShowHeaders: tb.showHeaders,
		Rows:        tb.rows,
		Options:     tb.options,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("Error marshaling table command: %v", err)
		return tb.parent
	}

	tb.parent.doc.Commands = append(tb.parent.doc.Commands, Command{
		Type: "table",
		Data: data,
	})

	return tb.parent
}

// Build returns the TableCommand without adding to parent (for standalone use)
func (tb *TableBuilder) Build() *TableCommand {
	return &TableCommand{
		Definition:  tb.definition,
		ShowHeaders: tb.showHeaders,
		Rows:        tb.rows,
		Options:     tb.options,
	}
}
