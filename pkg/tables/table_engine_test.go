package tables

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/adcondev/poster/pkg/constants"
)

// TestRender tests the TabEngine.Render method
func TestRender(t *testing.T) {
	// Common test data
	def := Definition{
		Columns: []Column{
			{Name: "Item", Width: 10, Align: constants.Left},
			{Name: "Qty", Width: 5, Align: constants.Center},
			{Name: "Price", Width: 8, Align: constants.Right},
		},
	}

	tests := []struct {
		name        string
		data        *Data
		opts        *Options
		expected    []string // Expected substrings
		unexpected  []string // Unexpected substrings
		expectError bool
		errorMsg    string
	}{
		{
			name: "Render Basic Table with Headers",
			data: &Data{
				Definition:  def,
				ShowHeaders: true,
				Rows: []Row{
					{"Apple", "10", "1.50"},
				},
			},
			opts: DefaultOptions(),
			expected: []string{
				"\x1bE\x01",                 // Enable Bold
				"Item        Qty     Price", // Header row
				"\x1bE\x00",                 // Disable Bold
				"Apple       10       1.50", // Data row
			},
			expectError: false,
		},
		{
			name: "Render Table without Headers",
			data: &Data{
				Definition:  def,
				ShowHeaders: false,
				Rows: []Row{
					{"Banana", "5", "0.99"},
				},
			},
			opts: &Options{
				PaperWidth:    80,
				ShowHeaders:   false,
				HeaderStyle:   Style{Bold: true},
				WordWrap:      true,
				ColumnSpacing: 1,
			},
			expected: []string{
				"Banana       5       0.99", // Data row
			},
			unexpected: []string{
				"Item", "Qty", "Price", // Headers should not be present
				"\x1bE\x01", // No bold command
			},
			expectError: false,
		},
		{
			name: "Render Table with Alignment",
			data: &Data{
				Definition: def,
				Rows: []Row{
					{"Left", "Cnt", "Right"},
				},
			},
			opts: &Options{
				PaperWidth:    80,
				ShowHeaders:   false,
				HeaderStyle:   Style{Bold: true},
				WordWrap:      true,
				ColumnSpacing: 1,
			},
			expected: []string{
				"Left        Cnt     Right",
			},
			expectError: false,
		},
		{
			name: "Render Table with Word Wrap",
			data: &Data{
				Definition: Definition{
					Columns: []Column{
						{Name: "Desc", Width: 5, Align: constants.Left},
					},
				},
				Rows: []Row{
					{"Long Text"}, // Should wrap to "Long " and "Text "
				},
			},
			opts: DefaultOptions(),
			expected: []string{
				"Long ",
				"Text ",
			},
			expectError: false,
		},
		{
			name: "Error: Nil Data",
			data: nil,
			opts: DefaultOptions(),
			expectError: true,
			errorMsg:    "table data cannot be nil",
		},
		{
			name: "Error: Invalid Data (Row Length Mismatch)",
			data: &Data{
				Definition: def,
				Rows: []Row{
					{"Apple", "10"}, // Missing one column
				},
			},
			opts: DefaultOptions(),
			expectError: true,
			errorMsg:    "invalid table data",
		},
		{
			name: "Error: Invalid Data (Empty definition)",
			data: &Data{
				Definition: Definition{},
				Rows: []Row{
					{"Apple"},
				},
			},
			opts: DefaultOptions(),
			expectError: true,
			errorMsg:    "table must have at least one column",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewEngine(&def, tt.opts)
			var buf bytes.Buffer

			err := engine.Render(&buf, tt.data)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				output := buf.String()

				for _, exp := range tt.expected {
					assert.Contains(t, output, exp, "Output should contain expected string")
				}

				for _, unexp := range tt.unexpected {
					assert.NotContains(t, output, unexp, "Output should NOT contain unexpected string")
				}

				// Verify Line Feeds are present
				if len(tt.expected) > 0 {
					assert.True(t, strings.HasSuffix(output, "\n"), "Output should end with a newline")
				}
			}
		})
	}
}
