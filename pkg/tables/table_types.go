package tables

import (
	"fmt"

	"github.com/adcondev/poster/pkg/constants"
)

// Column defines a table column configuration
type Column struct {
	Header string              `json:"header"`
	Width  int                 `json:"width"`
	Align  constants.Alignment `json:"align"`
}

// Definition defines the structure of a table
type Definition struct {
	Columns    []Column `json:"columns"`
	PaperWidth int      `json:"paper_width,omitempty"`
}

// Row represents a single row of data
type Row []string

// Data holds the complete table data
type Data struct {
	Definition  Definition `json:"definition"`
	ShowHeaders bool       `json:"show_headers,omitempty"`
	Rows        []Row      `json:"rows"`
}

// ValidateWidths checks if the total column widths fit within the specified max characters.
//
// Parameters:
//   - maxChars: The maximum characters per line (calculated from DotsPerLine / FontWidth)
//   - columnSpacing:  Spaces between columns
//
// This method should be called with maxChars calculated from the printer profile,
// not from user-provided values.
func (d *Definition) ValidateWidths(maxChars, columnSpacing int) error {
	if maxChars <= 0 {
		return fmt.Errorf("maxChars must be positive, got %d", maxChars)
	}
	if columnSpacing < 0 {
		columnSpacing = 0
	}

	total := 0
	for _, col := range d.Columns {
		if col.Width <= 0 {
			return fmt.Errorf("column '%s' has invalid width: %d", col.Header, col.Width)
		}
		total += col.Width
	}

	gaps := 0
	if len(d.Columns) > 1 {
		gaps = (len(d.Columns) - 1) * columnSpacing
	}

	totalRequired := total + gaps
	if totalRequired > maxChars {
		return fmt.Errorf(
			"columns too wide: %d chars (columns) + %d chars (gaps) = %d, max allowed: %d",
			total, gaps, totalRequired, maxChars,
		)
	}
	return nil
}

// Validate checks if the table data is valid
func (dt *Data) Validate() error {
	if len(dt.Definition.Columns) == 0 {
		return fmt.Errorf("table must have at least one column")
	}

	// Validate each row has correct number of cells
	expectedCells := len(dt.Definition.Columns)
	for i, row := range dt.Rows {
		if len(row) != expectedCells {
			return fmt.Errorf("row %d has %d cells, expected %d", i, len(row), expectedCells)
		}
	}

	return nil
}
