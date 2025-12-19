// Package tables provides table generation and rendering for ESC/POS printers
package tables

import (
	"strings"
	"unicode/utf8"

	"github.com/adcondev/poster/pkg/constants"
)

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
func PadString(s string, width int, align constants.Alignment) string {
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
	case constants.Center:
		padLeft := padTotal / 2
		padRight := padTotal - padLeft
		return strings.Repeat(" ", padLeft) + s + strings.Repeat(" ", padRight)
	case constants.Right:
		return strings.Repeat(" ", padTotal) + s
	case constants.Left:
		return s + strings.Repeat(" ", padTotal)
	default:
		// FIXME: it has to depend on constant.DefaultTableColumnAlignment
		padLeft := padTotal / 2
		padRight := padTotal - padLeft
		return strings.Repeat(" ", padLeft) + s + strings.Repeat(" ", padRight)
	}
}
