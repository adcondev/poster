package tables_test

import (
	"strings"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/tables"
)

// ============================================================================
// Tests para WrapText - Función con lógica compleja de word wrapping
// ============================================================================

func TestWrapText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		expected []string
	}{
		{
			name:     "no wrap needed - text fits",
			text:     "Hello",
			width:    10,
			expected: []string{"Hello"},
		},
		{
			name:     "basic word wrap",
			text:     "Hello World",
			width:    7,
			expected: []string{"Hello", "World"},
		},
		{
			name:     "multiple words wrap to multiple lines",
			text:     "The quick brown fox jumps",
			width:    10,
			expected: []string{"The quick", "brown fox", "jumps"},
		},
		{
			name:     "empty text returns empty string slice",
			text:     "",
			width:    10,
			expected: []string{""},
		},
		{
			name:     "single word longer than width forces break",
			text:     "Supercalifragilisticexpialidocious",
			width:    10,
			expected: []string{"Supercalif", "ragilistic", "expialidoc", "ious"},
		},
		{
			name:     "zero width returns original text",
			text:     "Test",
			width:    0,
			expected: []string{"Test"},
		},
		{
			name:     "negative width returns original text",
			text:     "Test",
			width:    -5,
			expected: []string{"Test"},
		},
		{
			name:     "exact width fit",
			text:     "12345",
			width:    5,
			expected: []string{"12345"},
		},
		{
			name:     "multiple spaces collapsed by Fields",
			text:     "Hello    World",
			width:    10,
			expected: []string{"Hello", "World"},
		},
		{
			name:     "unicode characters counted correctly",
			text:     "Hëllo Wörld",
			width:    7,
			expected: []string{"Hëllo", "Wörld"},
		},
		{
			name:     "mixed short and long words",
			text:     "A verylongwordhere B",
			width:    8,
			expected: []string{"A", "verylong", "wordhere", "B"},
		},
		{
			name:     "word exactly at width boundary",
			text:     "12345 67890",
			width:    5,
			expected: []string{"12345", "67890"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tables.WrapText(tt.text, tt.width)
			if len(result) != len(tt.expected) {
				t.Errorf("WrapText(%q, %d) returned %d lines, want %d\nGot: %v\nWant: %v",
					tt.text, tt.width, len(result), len(tt.expected), result, tt.expected)
				return
			}
			for i, line := range result {
				if line != tt.expected[i] {
					t.Errorf("WrapText(%q, %d) line %d = %q, want %q",
						tt.text, tt.width, i, line, tt.expected[i])
				}
			}
		})
	}
}

// ============================================================================
// Tests para PadString - Lógica de alineación y padding
// ============================================================================

// TestPadString tests the PadString function for various alignments and edge cases
func TestPadString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		align    constants.Alignment
		expected string
	}{
		{
			name:     "center align default",
			input:    "Hi",
			width:    5,
			align:    "",
			expected: " Hi  ",
		},
		{
			name:     "Right align",
			input:    "Hi",
			width:    5,
			align:    constants.Right,
			expected: "   Hi",
		},
		{
			name:     "Center align even padding",
			input:    "Hi",
			width:    6,
			align:    constants.Center,
			expected: "  Hi  ",
		},
		{
			name:     "Center align odd padding - left gets less",
			input:    "Hi",
			width:    5,
			align:    constants.Center,
			expected: " Hi  ",
		},
		{
			name:     "string exactly matches width",
			input:    "Hello",
			width:    5,
			align:    constants.Center,
			expected: "Hello",
		},
		{
			name:     "string longer than width gets truncated",
			input:    "Hello World",
			width:    5,
			align:    constants.Center,
			expected: "Hello",
		},
		{
			name:     "empty string left aligned",
			input:    "",
			width:    3,
			align:    "",
			expected: "   ",
		},
		{
			name:     "empty string Center aligned",
			input:    "",
			width:    4,
			align:    constants.Center,
			expected: "    ",
		},
		{
			name:     "unicode string Right aligned",
			input:    "Hëllo",
			width:    8,
			align:    constants.Right,
			expected: "   Hëllo",
		},
		{
			name:     "unicode truncation",
			input:    "Hëllo Wörld",
			width:    5,
			align:    "",
			expected: "Hëllo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tables.PadString(tt.input, tt.width, tt.align)
			if result != tt.expected {
				t.Errorf("PadString(%q, %d, %q) = %q, want %q",
					tt.input, tt.width, tt.align, result, tt.expected)
			}
			// Verificar que el resultado tiene exactamente el ancho esperado
			if len([]rune(result)) != tt.width {
				t.Errorf("PadString result rune length = %d, want %d",
					len([]rune(result)), tt.width)
			}
		})
	}
}

// ============================================================================
// Tests para Definition. ValidateWidths - Validación de anchos de columna
// ============================================================================

func TestDefinition_ValidateWidths_WithMaxChars(t *testing.T) {
	tests := []struct {
		name          string
		columns       []tables.Column
		maxChars      int
		columnSpacing int
		wantErr       bool
		errContains   string
	}{
		{
			name: "valid table for 58mm printer",
			columns: []tables.Column{
				{Header: "Item", Width: 15},
				{Header: "Price", Width: 15},
			},
			maxChars:      32,
			columnSpacing: 1,
			wantErr:       false,
			// 15 + 15 + 1 gap = 31 chars <= 32 max ✓
		},
		{
			name: "overflow on 58mm printer",
			columns: []tables.Column{
				{Header: "Product", Width: 20},
				{Header: "Qty", Width: 10},
				{Header: "Price", Width: 10},
			},
			maxChars:      32,
			columnSpacing: 1,
			wantErr:       true,
			errContains:   "columns too wide",
			// 20 + 10 + 10 + 2 gaps = 42 chars > 32 max ✗
		},
		{
			name: "exact fit on 80mm printer",
			columns: []tables.Column{
				{Header: "A", Width: 15},
				{Header: "B", Width: 15},
				{Header: "C", Width: 16},
			},
			maxChars:      48,
			columnSpacing: 1,
			wantErr:       false,
			// 15 + 15 + 16 + 2 gaps = 48 chars = 48 max ✓
		},
		{
			name: "single column uses full width",
			columns: []tables.Column{
				{Header: "Full Width Column", Width: 32},
			},
			maxChars:      32,
			columnSpacing: 1,
			wantErr:       false,
			// No gaps for single column, 32 = 32 ✓
		},
		{
			name: "rejects zero width column",
			columns: []tables.Column{
				{Header: "Bad Column", Width: 0},
			},
			maxChars:      32,
			columnSpacing: 1,
			wantErr:       true,
			errContains:   "invalid width",
		},
		{
			name: "rejects negative width column",
			columns: []tables.Column{
				{Header: "Negative", Width: -5},
			},
			maxChars:      32,
			columnSpacing: 1,
			wantErr:       true,
			errContains:   "invalid width",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def := tables.Definition{Columns: tt.columns}
			err := def.ValidateWidths(tt.maxChars, tt.columnSpacing)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWidths() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ValidateWidths() error = %q, should contain %q",
						err.Error(), tt.errContains)
				}
			}
		})
	}
}

// ============================================================================
// Tests para Data.Validate - Validación de datos de tabla
// ============================================================================

func TestData_Validate(t *testing.T) {
	tests := []struct {
		name    string
		data    tables.Data
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid data with matching columns and cells",
			data: tables.Data{
				Definition: tables.Definition{
					Columns: []tables.Column{{Width: 10}, {Width: 10}},
				},
				Rows: []tables.Row{{"a", "b"}, {"c", "d"}},
			},
			wantErr: false,
		},
		{
			name: "no columns defined",
			data: tables.Data{
				Definition: tables.Definition{Columns: []tables.Column{}},
				Rows:       []tables.Row{},
			},
			wantErr: true,
			errMsg:  "at least one column",
		},
		{
			name: "row with fewer cells than columns",
			data: tables.Data{
				Definition: tables.Definition{
					Columns: []tables.Column{{Width: 10}, {Width: 10}},
				},
				Rows: []tables.Row{{"only one"}},
			},
			wantErr: true,
			errMsg:  "has 1 cells, expected 2",
		},
		{
			name: "row with more cells than columns",
			data: tables.Data{
				Definition: tables.Definition{
					Columns: []tables.Column{{Width: 10}},
				},
				Rows: []tables.Row{{"one", "two", "three"}},
			},
			wantErr: true,
			errMsg:  "has 3 cells, expected 1",
		},
		{
			name: "empty rows is valid",
			data: tables.Data{
				Definition: tables.Definition{
					Columns: []tables.Column{{Width: 10}},
				},
				Rows: []tables.Row{},
			},
			wantErr: false,
		},
		{
			name: "mixed valid and invalid rows - fails on first invalid",
			data: tables.Data{
				Definition: tables.Definition{
					Columns: []tables.Column{{Width: 10}, {Width: 10}},
				},
				Rows: []tables.Row{{"a", "b"}, {"only one"}},
			},
			wantErr: true,
			errMsg:  "row 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %q, should contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TODO: Use string.Contains from strings package when available
func Contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && SearchSubstring(s, substr)))
}

func SearchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
