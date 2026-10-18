package executor

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/adcondev/poster/pkg/service"
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
				if cmd.Options.HeaderBold == nil || !*cmd.Options.HeaderBold {
					t.Error("Expected header_bold to be true")
				}
				if cmd.Options.WordWrap == nil || *cmd.Options.WordWrap {
					t.Error("Expected word_wrap to be false")
				}
				if cmd.Options.ColumnSpacing == nil || *cmd.Options.ColumnSpacing != 2 {
					t.Errorf("Expected column_spacing 2, got %v", cmd.Options.ColumnSpacing)
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
				if cmd.ShowHeaders == nil || !*cmd.ShowHeaders {
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

func TestTableCommand_Parsing_AutoReduceTrue(t *testing.T) {
	raw := `{
		"definition": {"columns": [{"name": "Item", "width": 20}, {"name": "Price", "width": 10}]},
		"rows": [["Coffee", "$4.50"]],
		"options": {"auto_reduce": true}
	}`

	var cmd TableCommand
	if err := json.Unmarshal([]byte(raw), &cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cmd.Options == nil || cmd.Options.AutoReduce == nil {
		t.Fatalf("expected options.auto_reduce to be parsed")
	}
	if *cmd.Options.AutoReduce != true {
		t.Fatalf("expected auto_reduce=true, got %v", *cmd.Options.AutoReduce)
	}
}

func TestTableCommand_Parsing_AutoReduceFalse(t *testing.T) {
	raw := `{
		"definition": {"columns": [{"name": "Item", "width": 20}, {"name": "Price", "width": 10}]},
		"rows": [["Coffee", "$4.50"]],
		"options": {"auto_reduce": false}
	}`

	var cmd TableCommand
	if err := json.Unmarshal([]byte(raw), &cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cmd.Options == nil || cmd.Options.AutoReduce == nil {
		t.Fatalf("expected options.auto_reduce to be parsed")
	}
	if *cmd.Options.AutoReduce != false {
		t.Fatalf("expected auto_reduce=false, got %v", *cmd.Options.AutoReduce)
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
			name: "show_headers defaults to nil (handler applies true)",
			json: `{
				"definition": {"columns":  [{"name": "Col", "width":  10}]},
				"rows":  [["Data"]]
			}`,
			checkFunc: func(t *testing.T, cmd TableCommand) {
				if cmd.ShowHeaders != nil {
					t.Error("Expected show_headers to be nil when omitted")
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

func TestTableCommand_Parsing_AutoReduceOmitted_IsNil(t *testing.T) {
	raw := `{
		"definition": {"columns": [{"name": "Item", "width": 20}, {"name": "Price", "width": 10}]},
		"rows": [["Coffee", "$4.50"]],
		"options": {"header_bold": true}
	}`

	var cmd TableCommand
	if err := json.Unmarshal([]byte(raw), &cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cmd.Options == nil {
		t.Fatalf("expected options to be non-nil")
	}
	if cmd.Options.AutoReduce != nil {
		t.Fatalf("expected auto_reduce to be nil when omitted, got %v", *cmd.Options.AutoReduce)
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

// ============================================================================
// Table Command Execution Tests
// ============================================================================

func TestHandleTable_DefaultFontAAndBoldHeaders(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	tableJSON := []byte(`{
		"definition": {
			"columns": [
				{"name": "Item", "width": 10},
				{"name": "Price", "width": 8}
			]
		},
		"rows": [
			["Burger", "$5.99"],
			["Fries", "$2.49"]
		]
	}`)

	err := exec.handleTable(mock, tableJSON)
	if err != nil {
		t.Fatalf("handleTable failed: %v", err)
	}

	if mock.CallCount("FontA") < 1 {
		t.Errorf("Expected FontA to be called, got %d calls", mock.CallCount("FontA"))
	}
	if mock.CallCount("EnableBold") != 1 {
		t.Errorf("Expected EnableBold to be called 1 time, got %d", mock.CallCount("EnableBold"))
	}
	if mock.CallCount("DisableBold") != 1 {
		t.Errorf("Expected DisableBold to be called 1 time, got %d", mock.CallCount("DisableBold"))
	}

	allOutput := strings.Join(mock.PrintedText, "")
	if !strings.Contains(allOutput, "Item") || !strings.Contains(allOutput, "Burger") {
		t.Errorf("Expected table output to contain headers and body, got: %q", allOutput)
	}
}

func TestHandleTable_FontB(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	tableJSON := []byte(`{
		"definition": {
			"columns": [
				{"name": "Item", "width": 15},
				{"name": "Qty", "width": 5}
			]
		},
		"options": {
			"font": "B"
		},
		"rows": [
			["Coffee", "2"]
		]
	}`)

	err := exec.handleTable(mock, tableJSON)
	if err != nil {
		t.Fatalf("handleTable failed: %v", err)
	}

	if mock.CallCount("FontB") != 1 {
		t.Errorf("Expected FontB to be called 1 time, got %d", mock.CallCount("FontB"))
	}
	// FontA should be restored deferred
	if mock.CallCount("FontA") < 1 {
		t.Errorf("Expected FontA to be restored, got %d calls", mock.CallCount("FontA"))
	}
}

func TestHandleTable_NoHeaders(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	tableJSON := []byte(`{
		"definition": {
			"columns": [
				{"name": "Col1", "width": 10}
			]
		},
		"show_headers": false,
		"rows": [
			["Value1"]
		]
	}`)

	err := exec.handleTable(mock, tableJSON)
	if err != nil {
		t.Fatalf("handleTable failed: %v", err)
	}

	if mock.CallCount("EnableBold") != 0 {
		t.Errorf("Expected 0 calls to EnableBold when ShowHeaders is false, got %d", mock.CallCount("EnableBold"))
	}
	allOutput := strings.Join(mock.PrintedText, "")
	if strings.Contains(allOutput, "Col1") {
		t.Errorf("Expected output not to contain header 'Col1', got: %q", allOutput)
	}
	if !strings.Contains(allOutput, "Value1") {
		t.Errorf("Expected output to contain 'Value1', got: %q", allOutput)
	}
}

func TestHandleTable_Align(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	tableJSON := []byte(`{
		"definition": {
			"columns": [
				{"name": "Name", "width": 8}
			]
		},
		"options": {
			"align": "center"
		},
		"rows": [
			["Alice"]
		]
	}`)

	err := exec.handleTable(mock, tableJSON)
	if err != nil {
		t.Fatalf("handleTable failed: %v", err)
	}

	if mock.CallCount("AlignCenter") != 1 {
		t.Errorf("Expected AlignCenter to be called 1 time, got %d", mock.CallCount("AlignCenter"))
	}
	if mock.CallCount("AlignLeft") < 1 {
		t.Errorf("Expected AlignLeft to be restored at end, got %d", mock.CallCount("AlignLeft"))
	}
}

func TestHandleTable_ValidationErrors(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("empty columns", func(t *testing.T) {
		err := exec.handleTable(mock, []byte(`{"definition": {"columns": []}, "rows": []}`))
		if err == nil {
			t.Fatal("Expected error for empty columns, got nil")
		}
	})

	t.Run("invalid column width", func(t *testing.T) {
		err := exec.handleTable(mock, []byte(`{"definition": {"columns": [{"name": "A", "width": 0}]}, "rows": []}`))
		if err == nil {
			t.Fatal("Expected error for zero width column, got nil")
		}
	})

	t.Run("table overflow auto-reduce false", func(t *testing.T) {
		// Mock profile paper width is 80mm, dotsPerLine 576. Font A maxChars is 48.
		// Two columns of 30 chars each + 1 gap = 61 > 48.
		overflowJSON := []byte(`{
			"definition": {
				"columns": [
					{"name": "Col1", "width": 30},
					{"name": "Col2", "width": 30}
				]
			},
			"options": {
				"auto_reduce": false
			},
			"rows": [["a", "b"]]
		}`)
		err := exec.handleTable(mock, overflowJSON)
		if err == nil {
			t.Fatal("Expected error for table overflow without auto-reduce, got nil")
		}
	})
}

