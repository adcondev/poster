package builder

import (
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/tables"
)

type (
	tableColumn     = tables.Column
	tableDefinition = tables.Definition
	tableOptions    = schema.TableOptions
	tableCommand    = schema.TableCommand
)

func boolPtr(v bool) *bool {
	return &v
}

// TableBuilder constructs table commands
type TableBuilder struct {
	parent      *DocumentBuilder
	columns     []tableColumn
	rows        [][]string
	paperWidth  int
	showHeaders *bool
	options     *tableOptions
}

func newTableBuilder(parent *DocumentBuilder) *TableBuilder {
	return &TableBuilder{
		parent:      parent,
		columns:     []tableColumn{},
		rows:        [][]string{},
		showHeaders: boolPtr(constants.DefaultTableShowHeaders),
		options: &tableOptions{
			HeaderBold: boolPtr(constants.DefaultTableHeaderBold),
			WordWrap:   boolPtr(constants.DefaultTableWordWrap),
		},
	}
}

// Column adds a column definition
func (tb *TableBuilder) Column(header string, width int, align ...constants.Alignment) *TableBuilder {
	col := tableColumn{
		Name:  header,
		Width: width,
	}
	if len(align) > 0 {
		col.Align = align[0]
	}
	tb.columns = append(tb.columns, col)
	return tb
}

// Row adds a data row
func (tb *TableBuilder) Row(cells ...string) *TableBuilder {
	tb.rows = append(tb.rows, cells)
	return tb
}

// Rows adds multiple rows at once
func (tb *TableBuilder) Rows(rows [][]string) *TableBuilder {
	tb.rows = append(tb.rows, rows...)
	return tb
}

// PaperWidth sets the paper width in characters
func (tb *TableBuilder) PaperWidth(width int) *TableBuilder {
	tb.paperWidth = width
	return tb
}

// HideHeaders hides the table headers
func (tb *TableBuilder) HideHeaders() *TableBuilder {
	tb.showHeaders = boolPtr(false)
	return tb
}

// ShowHeaders shows the table headers
func (tb *TableBuilder) ShowHeaders() *TableBuilder {
	tb.showHeaders = boolPtr(true)
	return tb
}

// HeaderBold enables bold headers
func (tb *TableBuilder) HeaderBold() *TableBuilder {
	tb.options.HeaderBold = boolPtr(true)
	return tb
}

// NoHeaderBold disables bold headers
func (tb *TableBuilder) NoHeaderBold() *TableBuilder {
	tb.options.HeaderBold = boolPtr(false)
	return tb
}

// NoWordWrap disables word wrapping
func (tb *TableBuilder) NoWordWrap() *TableBuilder {
	tb.options.WordWrap = boolPtr(false)
	return tb
}

// WordWrap enables word wrapping
func (tb *TableBuilder) WordWrap() *TableBuilder {
	tb.options.WordWrap = boolPtr(true)
	return tb
}

// ColumnSpacing sets spacing between columns
func (tb *TableBuilder) ColumnSpacing(spacing int) *TableBuilder {
	if spacing < 0 {
		spacing = 0
	}
	tb.options.ColumnSpacing = intPtr(spacing)
	return tb
}

// Align sets table alignment
func (tb *TableBuilder) Align(align constants.Alignment) *TableBuilder {
	tb.options.Align = string(align)
	return tb
}

// AutoReduce enables automatic column width reduction when table exceeds paper width.
func (tb *TableBuilder) AutoReduce() *TableBuilder {
	tb.options.AutoReduce = boolPtr(true)
	return tb
}

// NoAutoReduce disables automatic column width reduction.
func (tb *TableBuilder) NoAutoReduce() *TableBuilder {
	tb.options.AutoReduce = boolPtr(false)
	return tb
}

// End finishes the table and returns to document builder
func (tb *TableBuilder) End() *DocumentBuilder {
	cmd := tableCommand{
		Definition: tableDefinition{
			Columns:    tb.columns,
			PaperWidth: tb.paperWidth,
		},
		ShowHeaders: tb.showHeaders,
		Rows:        tb.rows,
		Options:     tb.options,
	}
	return tb.parent.addCommand("table", cmd)
}
