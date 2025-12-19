package tables

import (
	"fmt"

	"github.com/adcondev/poster/pkg/constants"
)

// ReduceResult contains the outcome of a reduction attempt
type ReduceResult struct {
	Success       bool
	Columns       []Column
	OriginalWidth int
	ReducedWidth  int
	Reductions    int
}

// ReduceToFit attempts to reduce column widths to fit within maxChars.
// It prioritizes reducing the widest columns first to preserve legibility of smaller columns.
func ReduceToFit(columns []Column, maxChars, columnSpacing, minWidth int) (*ReduceResult, error) {
	if minWidth <= 0 {
		minWidth = constants.MinTableColumnWidth
	}

	// Work on a deep copy to avoid mutating original definition on failure
	reduced := make([]Column, len(columns))
	copy(reduced, columns)

	originalWidth := calculateTotalWidth(reduced, columnSpacing)
	currentWidth := originalWidth
	reductions := 0

	// Safety limit to prevent infinite loops
	const maxIterations = 1000

	for currentWidth > maxChars {
		if reductions > maxIterations {
			return nil, fmt.Errorf("auto-reduction safety limit exceeded")
		}

		// Find the index of the longest column that is still > minWidth
		targetIdx := -1
		maxColWidth := -1

		for i, col := range reduced {
			if col.Width > minWidth {
				if col.Width > maxColWidth {
					maxColWidth = col.Width
					targetIdx = i
				}
			}
		}

		// If no column can be reduced further, we fail
		if targetIdx == -1 {
			return &ReduceResult{
				Success:       false,
				Columns:       reduced,
				OriginalWidth: originalWidth,
				ReducedWidth:  currentWidth,
				Reductions:    reductions,
			}, fmt.Errorf("cannot reduce table further: all columns at minimum width (%d)", minWidth)
		}

		// Reduce the found column
		reduced[targetIdx].Width--
		currentWidth--
		reductions++
	}

	return &ReduceResult{
		Success:       true,
		Columns:       reduced,
		OriginalWidth: originalWidth,
		ReducedWidth:  currentWidth,
		Reductions:    reductions,
	}, nil
}

// calculateTotalWidth calculates total width including gaps
func calculateTotalWidth(columns []Column, spacing int) int {
	total := 0
	for _, col := range columns {
		total += col.Width
	}
	if len(columns) > 1 {
		total += (len(columns) - 1) * spacing
	}
	return total
}
