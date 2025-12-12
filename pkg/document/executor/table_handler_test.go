package executor

import (
	"encoding/json"
	"testing"
)

// ============================================================================
// Table Command Parsing Tests
// ============================================================================

func TestTableCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd TableCommand)
	}{
		{
			name: "simple table",
			json: `{
				"definition": {"columns": [{"name": "Item", "width": 20}, {"name": "Price", "width": 10}]},
				"rows": [["Coffee", "$4.50"], ["Muffin", "$3.00"]]
			}`,
			checkFunc: func(t *testing.T, cmd TableCommand) {
				if len(cmd.Definition.Columns) != 2 {
					t.Errorf("Expected 2 columns, got %d", len(cmd.Definition.Columns))
				}
				if len(cmd.Rows) != 2 {
					t.Errorf("Expected 2 rows, got %d", len(cmd.Rows))
				}
			},
		},
		{
			name: "table with options",
			json: `{
				"definition": {"columns":  [{"name": "Test", "width": 10}]},
				"rows": [["Data"]],
				"options": {"header_bold": true, "word_wrap": false, "column_spacing": 2}
			}`,
			checkFunc: func(t *testing.T, cmd TableCommand) {
				if cmd.Options == nil {
					t.Fatal("Expected options to be set")
				}
				if !cmd.Options.HeaderBold {
					t.Error("Expected header_bold to be true")
				}
				if cmd.Options.WordWrap {
					t.Error("Expected word_wrap to be false")
				}
				if cmd.Options.ColumnSpacing != 2 {
					t.Errorf("Expected column_spacing 2, got %d", cmd.Options.ColumnSpacing)
				}
			},
		},
		{
			name: "table with show_headers",
			json: `{
				"definition": {"columns":  [{"name": "Col", "width": 10}]},
				"rows": [["Data"]],
				"show_headers": true
			}`,
			checkFunc: func(t *testing.T, cmd TableCommand) {
				if !cmd.ShowHeaders {
					t.Error("Expected show_headers to be true")
				}
			},
		},
		{
			name: "table with paper width",
			json: `{
				"definition": {"columns":  [{"name": "Col", "width":  10}], "paper_width":  48},
				"rows": [["Data"]]
			}`,
			checkFunc: func(t *testing.T, cmd TableCommand) {
				if cmd.Definition.PaperWidth != 48 {
					t.Errorf("Expected paper_width 48, got %d", cmd.Definition.PaperWidth)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd TableCommand
			err := json.Unmarshal([]byte(tt.json), &cmd)
			if err != nil {
				t.Fatalf("Unexpected error:  %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

// ============================================================================
// Table Command Default Value Tests
// ============================================================================

func TestTableCommand_Defaults(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd TableCommand)
	}{
		{
			name: "nil options uses defaults",
			json: `{
				"definition": {"columns":  [{"name": "Col", "width":  10}]},
				"rows":  [["Data"]]
			}`,
			checkFunc: func(t *testing.T, cmd TableCommand) {
				if cmd.Options != nil {
					t.Error("Expected options to be nil (handler applies defaults)")
				}
			},
		},
		{
			name: "show_headers defaults to false",
			json: `{
				"definition": {"columns":  [{"name": "Col", "width":  10}]},
				"rows":  [["Data"]]
			}`,
			checkFunc: func(t *testing.T, cmd TableCommand) {
				if cmd.ShowHeaders {
					t.Error("Expected show_headers to default to false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd TableCommand
			err := json.Unmarshal([]byte(tt.json), &cmd)
			if err != nil {
				t.Fatalf("Unexpected error:  %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

// ============================================================================
// Table Command Validation Tests
// ============================================================================

func TestTableCommand_Validation(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		expectError bool
	}{
		{
			name:        "empty columns parses but invalid for handler",
			json:        `{"definition": {"columns": []}, "rows": [["Data"]]}`,
			expectError: false,
		},
		{
			name:        "empty rows is valid",
			json:        `{"definition": {"columns": [{"name": "Col", "width": 10}]}, "rows": []}`,
			expectError: false,
		},
		{
			name:        "invalid json",
			json:        `{invalid}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd TableCommand
			err := json.Unmarshal([]byte(tt.json), &cmd)
			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
