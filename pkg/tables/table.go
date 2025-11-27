// Package tables provides table generation and rendering for ESC/POS printers
package tables

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Alignment represents text alignment within a cell
type Alignment string

const (
	// Center alignment
	Center Alignment = "Center"
	// Right alignment
	Right Alignment = "Right"
)

// Column defines a table column configuration
type Column struct {
	Header string    `json:"header"`
	Width  int       `json:"width"`
	Align  Alignment `json:"align"`
}

// Definition defines the structure of a table
type Definition struct {
	Columns    []Column `json:"columns"`
	PaperWidth int      `json:"paper_width,omitempty"`
}

// ValidateWidths checks if the total column widths fit within the paper width
func (d *Definition) ValidateWidths() error {
	total := 0
	for _, col := range d.Columns {
		total += col.Width
	}
	gaps := (len(d.Columns) - 1) * 1 // ColumnSpacing default
	if total+gaps > d.PaperWidth {
		return fmt.Errorf("columns too wide: %d chars (max %d)",
			total+gaps, d.PaperWidth)
	}
	return nil
}

// Row represents a single row of data
type Row []string

// Data holds the complete table data
type Data struct {
	Definition  Definition `json:"definition"`
	ShowHeaders bool       `json:"show_headers,omitempty"`
	Rows        []Row      `json:"rows"`
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

// Validate checks if the table definition is valid
func (d *Definition) Validate() error {
	totalWidth := 0
	for _, col := range d.Columns {
		totalWidth += col.Width
	}

	// Add space for column spacing
	totalWidth += (len(d.Columns) - 1) * 1 // Default spacing

	if totalWidth > d.PaperWidth {
		return fmt.Errorf("total column width (%d) exceeds paper width (%d)",
			totalWidth, d.PaperWidth)
	}
	return nil
}

// WrapText wraps text to fit within the specified width
func WrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)

	if len(words) == 0 {
		return []string{""}
	}

	currentLine := words[0]

	for _, word := range words[1:] {
		// Check if adding this word would exceed the width
		testLine := currentLine + " " + word
		if utf8.RuneCountInString(testLine) > width {
			// Current line is full, start a new one
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	// Add the last line
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	// Handle case where a single word is longer than width
	var finalLines []string
	for _, line := range lines {
		if utf8.RuneCountInString(line) > width {
			// Force break long words
			runes := []rune(line)
			for i := 0; i < len(runes); i += width {
				end := i + width
				if end > len(runes) {
					end = len(runes)
				}
				finalLines = append(finalLines, string(runes[i:end]))
			}
		} else {
			finalLines = append(finalLines, line)
		}
	}

	return finalLines
}

// PadString pads a string according to alignment
func PadString(s string, width int, align Alignment) string {
	length := utf8.RuneCountInString(s)
	if length >= width {
		// Truncate if necessary
		runes := []rune(s)
		if len(runes) > width {
			return string(runes[:width])
		}
		return s
	}

	padTotal := width - length

	switch align {
	case Center:
		padLeft := padTotal / 2
		padRight := padTotal - padLeft
		return strings.Repeat(" ", padLeft) + s + strings.Repeat(" ", padRight)
	case Right:
		return strings.Repeat(" ", padTotal) + s
	default: // AlignLeft
		return s + strings.Repeat(" ", padTotal)
	}
}
