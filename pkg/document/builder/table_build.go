package builder

import (
	"github.com/adcondev/poster/pkg/constants"
)

// TableBuilder constructs table commands
type TableBuilder struct {
	parent      *DocumentBuilder
	columns     []tableColumn
	rows        [][]string
	paperWidth  int
	showHeaders bool
	options     *tableOptions
}

type tableColumn struct {
	Name  string `json:"name"`
	Width int    `json:"width"`
	Align string `json:"align,omitempty"`
}

type tableDefinition struct {
	Columns    []tableColumn `json:"columns"`
	PaperWidth int           `json:"paper_width,omitempty"`
}

type tableOptions struct {
	HeaderBold    bool   `json:"header_bold,omitempty"`
	WordWrap      bool   `json:"word_wrap,omitempty"`
	ColumnSpacing int    `json:"column_spacing,omitempty"`
	Align         string `json:"align,omitempty"`
}

type tableCommand struct {
	Definition  tableDefinition `json:"definition"`
	ShowHeaders bool            `json:"show_headers,omitempty"`
	Rows        [][]string      `json:"rows"`
	Options     *tableOptions   `json:"options,omitempty"`
}

func newTableBuilder(parent *DocumentBuilder) *TableBuilder {
	return &TableBuilder{
		parent:      parent,
		columns:     []tableColumn{},
		rows:        [][]string{},
		showHeaders: true,
		options: &tableOptions{
			HeaderBold: true,
			WordWrap:   true,
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
		col.Align = string(align[0])
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
	tb.showHeaders = false
	return tb
}

// NoHeaderBold disables bold headers
func (tb *TableBuilder) NoHeaderBold() *TableBuilder {
	tb.options.HeaderBold = false
	return tb
}

// NoWordWrap disables word wrapping
func (tb *TableBuilder) NoWordWrap() *TableBuilder {
	tb.options.WordWrap = false
	return tb
}

// ColumnSpacing sets spacing between columns
func (tb *TableBuilder) ColumnSpacing(spacing int) *TableBuilder {
	if spacing < 0 {
		spacing = 0
	}
	tb.options.ColumnSpacing = spacing
	return tb
}

// Align sets table alignment
func (tb *TableBuilder) Align(align constants.Alignment) *TableBuilder {
	tb.options.Align = string(align)
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
