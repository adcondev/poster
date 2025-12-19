package executor

import (
	"fmt"
	"strings"

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

	tableData := &tables.Data{
		Definition:  cmd.Definition,
		ShowHeaders: cmd.ShowHeaders,
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
