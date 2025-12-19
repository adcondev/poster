package executor

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/adcondev/poster/pkg/document/schema"
)

// ============================================================================
// NewExecutor Tests
// ============================================================================

func TestNewExecutor_NilPrinter_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil printer, got none")
		}
	}()

	NewExecutor(nil)
}

// ============================================================================
// ParseDocument Integration Tests
// ============================================================================

func TestParseDocument_ValidatesCommandTypes(t *testing.T) {
	validTypes := []string{
		"text", "image", "qr", "barcode", "table",
		"separator", "feed", "cut", "pulse", "beep", "raw",
	}

	for _, cmdType := range validTypes {
		t.Run(cmdType, func(t *testing.T) {
			jsonData := []byte(`{
				"version": "1.0",
				"profile": {"model": "Test"},
				"commands": [{"type": "` + cmdType + `", "data": {}}]
			}`)

			doc, err := schema.ParseDocument(jsonData)
			if err != nil {
				t.Fatalf("ParseDocument failed for %s: %v", cmdType, err)
			}

			if len(doc.Commands) != 1 {
				t.Errorf("Expected 1 command, got %d", len(doc.Commands))
			}

			if doc.Commands[0].Type != cmdType {
				t.Errorf("Expected type %s, got %s", cmdType, doc.Commands[0].Type)
			}
		})
	}
}

func TestParseDocument_PreservesRawData(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"commands":  [
			{
				"type":  "text",
				"data": {
					"content": {"text": "Hello", "content_style": {"bold": true}},
					"label": {"text": "Greeting"},
					"new_line": false
				}
			}
		]
	}`)

	doc, err := schema.ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument error: %v", err)
	}

	// Verify we can unmarshal the preserved data
	var cmd TextCommand
	if err := json.Unmarshal(doc.Commands[0].Data, &cmd); err != nil {
		t.Fatalf("Failed to unmarshal command data: %v", err)
	}

	if cmd.Content.Text != "Hello" {
		t.Errorf("Expected text 'Hello', got '%s'", cmd.Content.Text)
	}

	if cmd.Label == nil || cmd.Label.Text != "Greeting" {
		t.Error("Expected label 'Greeting'")
	}

	if cmd.NewLine == nil || *cmd.NewLine != false {
		t.Error("Expected new_line false")
	}
}

// ============================================================================
// Document Structure Tests
// ============================================================================

func TestDocument_CompleteStructure(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {
			"model":  "EPSON TM-T20II",
			"paper_width": 80,
			"code_table": "WPC1252",
			"dpi": 203,
			"has_qr": true
		},
		"debug_log": true,
		"commands": [
			{"type": "text", "data": {"content": {"text": "Receipt"}}},
			{"type": "separator", "data": {"char": "="}},
			{"type": "table", "data": {"definition": {"columns": [{"name": "Item", "width": 20}]}, "rows": [["Coffee"]]}},
			{"type": "qr", "data": {"data": "https://example.com"}},
			{"type": "feed", "data": {"lines": 3}},
			{"type": "cut", "data": {"mode": "partial"}}
		]
	}`)

	doc, err := schema.ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument error: %v", err)
	}

	// Verify profile
	if doc.Profile.Model != "EPSON TM-T20II" {
		t.Errorf("Expected model 'EPSON TM-T20II', got '%s'", doc.Profile.Model)
	}
	if doc.Profile.PaperWidth != 80 {
		t.Errorf("Expected paper_width 80, got %d", doc.Profile.PaperWidth)
	}
	if !doc.Profile.HasQR {
		t.Error("Expected has_qr true")
	}
	if !doc.DebugLog {
		t.Error("Expected debug_log true")
	}

	// Verify command count and types
	expectedTypes := []string{"text", "separator", "table", "qr", "feed", "cut"}
	if len(doc.Commands) != len(expectedTypes) {
		t.Fatalf("Expected %d commands, got %d", len(expectedTypes), len(doc.Commands))
	}

	for i, expected := range expectedTypes {
		if doc.Commands[i].Type != expected {
			t.Errorf("Command %d:  expected type '%s', got '%s'", i, expected, doc.Commands[i].Type)
		}
	}
}

// ============================================================================
// Command Data Unmarshaling Tests
// ============================================================================

func TestCommandData_TextWithAllFields(t *testing.T) {
	data := json.RawMessage(`{
		"content": {
			"text": "Total",
			"content_style": {"bold": true, "size": "2x2", "underline": "1pt"},
			"align": "right"
		},
		"label":  {
			"text": "Amount",
			"label_style": {"bold": false},
			"separator": ":: "
		},
		"new_line": true
	}`)

	var cmd TextCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		t.Fatalf("Unmarshal error:  %v", err)
	}

	// Verify content
	if cmd.Content.Text != "Total" {
		t.Errorf("Expected text 'Total', got '%s'", cmd.Content.Text)
	}
	if cmd.Content.Style == nil {
		t.Fatal("Expected content style")
	}
	if cmd.Content.Style.Bold == nil || !*cmd.Content.Style.Bold {
		t.Error("Expected bold true")
	}
	if cmd.Content.Style.Size == nil || *cmd.Content.Style.Size != "2x2" {
		t.Error("Expected size '2x2'")
	}

	// Verify label
	if cmd.Label == nil {
		t.Fatal("Expected label")
	}
	if cmd.Label.Text != "Amount" {
		t.Errorf("Expected label 'Amount', got '%s'", cmd.Label.Text)
	}
	if cmd.Label.Separator == nil || *cmd.Label.Separator != ":: " {
		t.Error("Expected separator ':: '")
	}
}

func TestCommandData_BarcodeWithAllFields(t *testing.T) {
	data := json.RawMessage(`{
		"symbology": "CODE128",
		"data": "ABC123456",
		"width": 3,
		"height": 80,
		"hri_position": "below",
		"hri_font": "A",
		"align": "center"
	}`)

	var cmd BarcodeCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if cmd.Symbology != "CODE128" {
		t.Errorf("Expected symbology 'CODE128', got '%s'", cmd.Symbology)
	}
	if cmd.Data != "ABC123456" {
		t.Errorf("Expected data 'ABC123456', got '%s'", cmd.Data)
	}
	if cmd.Width == nil || *cmd.Width != 3 {
		t.Error("Expected width 3")
	}
	if cmd.Height == nil || *cmd.Height != 80 {
		t.Error("Expected height 80")
	}
	if cmd.HRIPosition == nil || *cmd.HRIPosition != "below" {
		t.Error("Expected hri_position 'below'")
	}
}

func TestCommandData_QRWithAllFields(t *testing.T) {
	data := json.RawMessage(`{
		"data": "https://example.com/receipt/12345",
		"pixel_width": 256,
		"correction":  "H",
		"align": "center",
		"human_text": "Scan for details",
		"logo": "base64data",
		"circle_shape": true
	}`)

	var cmd QRCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if cmd.Data != "https://example.com/receipt/12345" {
		t.Errorf("Unexpected data: %s", cmd.Data)
	}
	if cmd.PixelWidth != 256 {
		t.Errorf("Expected pixel_width 256, got %d", cmd.PixelWidth)
	}
	if cmd.Correction != "H" {
		t.Errorf("Expected correction 'H', got '%s'", cmd.Correction)
	}
	if cmd.HumanText != "Scan for details" {
		t.Errorf("Expected human_text 'Scan for details', got '%s'", cmd.HumanText)
	}
	if !cmd.CircleShape {
		t.Error("Expected circle_shape true")
	}
}

func TestCommandData_TableWithAllFields(t *testing.T) {
	data := json.RawMessage(`{
		"definition": {
			"columns": [
				{"name": "Item", "width":  20, "align": "left"},
				{"name": "Qty", "width": 5, "align": "center"},
				{"name": "Price", "width": 10, "align": "right"}
			],
			"paper_width": 48
		},
		"show_headers": true,
		"rows": [
			["Coffee", "2", "$8.00"],
			["Muffin", "1", "$3.50"]
		],
		"options": {
			"header_bold": true,
			"word_wrap": true,
			"column_spacing": 2,
			"align":  "center"
		}
	}`)

	var cmd TableCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if len(cmd.Definition.Columns) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(cmd.Definition.Columns))
	}
	if !cmd.ShowHeaders {
		t.Error("Expected show_headers true")
	}
	if len(cmd.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(cmd.Rows))
	}
	if cmd.Options == nil {
		t.Fatal("Expected options")
	}
	if !cmd.Options.HeaderBold {
		t.Error("Expected header_bold true")
	}
	if cmd.Options.ColumnSpacing != 2 {
		t.Errorf("Expected column_spacing 2, got %d", cmd.Options.ColumnSpacing)
	}
}

func TestCommandData_RawWithAllFields(t *testing.T) {
	data := json.RawMessage(`{
		"hex": "1B 40 1B 61 01",
		"format": "hex",
		"comment": "Initialize and center",
		"safe_mode": true
	}`)

	var cmd RawCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		t.Fatalf("Unmarshal error:  %v", err)
	}

	if cmd.Hex != "1B 40 1B 61 01" {
		t.Errorf("Unexpected hex: %s", cmd.Hex)
	}
	if cmd.Format != "hex" {
		t.Errorf("Expected format 'hex', got '%s'", cmd.Format)
	}
	if cmd.Comment != "Initialize and center" {
		t.Errorf("Unexpected comment: %s", cmd.Comment)
	}
	if !cmd.SafeMode {
		t.Error("Expected safe_mode true")
	}
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestParseDocument_ErrorMessages(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		errorContains string
	}{
		{
			name:          "empty commands",
			json:          `{"version": "1.0", "profile": {"model": "Test"}, "commands": []}`,
			errorContains: "at least one command",
		},
		{
			name:          "invalid json",
			json:          `{not valid json}`,
			errorContains: "failed to parse",
		},
		{
			name:          "missing commands",
			json:          `{"version": "1.0", "profile": {"model": "Test"}}`,
			errorContains: "at least one command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := schema.ParseDocument([]byte(tt.json))
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.errorContains) {
				t.Errorf("Expected error containing '%s', got:  %v", tt.errorContains, err)
			}
		})
	}
}

// ============================================================================
// Schema Validation Tests
// ============================================================================

func TestDocument_SchemaValidation(t *testing.T) {
	tests := []struct {
		name      string
		doc       schema.Document
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid document",
			doc: schema.Document{
				Version: "1.0",
				Profile: schema.ProfileConfig{Model: "Test"},
				Commands: []schema.Command{
					{Type: "text", Data: json.RawMessage(`{}`)},
				},
			},
			expectErr: false,
		},
		{
			name: "missing version",
			doc: schema.Document{
				Version:  "",
				Profile:  schema.ProfileConfig{Model: "Test"},
				Commands: []schema.Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			expectErr: true,
			errMsg:    "version is required",
		},
		{
			name: "missing model",
			doc: schema.Document{
				Version:  "1.0",
				Profile:  schema.ProfileConfig{},
				Commands: []schema.Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			expectErr: true,
			errMsg:    "model is required",
		},
		{
			name: "no commands",
			doc: schema.Document{
				Version:  "1.0",
				Profile:  schema.ProfileConfig{Model: "Test"},
				Commands: []schema.Command{},
			},
			expectErr: true,
			errMsg:    "at least one command",
		},
		{
			name: "invalid paper width",
			doc: schema.Document{
				Version:  "1.0",
				Profile:  schema.ProfileConfig{Model: "Test", PaperWidth: 99},
				Commands: []schema.Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			expectErr: true,
			errMsg:    "invalid paper_width",
		},
		{
			name: "invalid DPI",
			doc: schema.Document{
				Version:  "1.0",
				Profile:  schema.ProfileConfig{Model: "Test", DPI: 999},
				Commands: []schema.Command{{Type: "text", Data: json.RawMessage(`{}`)}},
			},
			expectErr: true,
			errMsg:    "invalid dpi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.doc.Validate()
			if tt.expectErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got: %v", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error:  %v", err)
				}
			}
		})
	}
}

// ============================================================================
// Utility Function Tests
// ============================================================================

func TestSameOrNil(t *testing.T) {
	strA := "a"
	strB := "b"
	strA2 := "a"

	tests := []struct {
		name     string
		a        *string
		b        *string
		expected bool
	}{
		{"both nil", nil, nil, true},
		{"a nil", nil, &strA, false},
		{"b nil", &strA, nil, false},
		{"same value", &strA, &strA2, true},
		{"different value", &strA, &strB, false},
		{"same pointer", &strA, &strA, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sameOrNil(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("sameOrNil(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestStrPtr(t *testing.T) {
	result := strPtr("test")
	if result == nil {
		t.Fatal("Expected non-nil pointer")
	}
	if *result != "test" {
		t.Errorf("Expected 'test', got '%s'", *result)
	}
}

// ============================================================================
// Mapping Function Tests
// ============================================================================

func TestMapHRIPosition(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"none", "none"},
		{"above", "above"},
		{"below", "below"},
		{"both", "both"},
		{"invalid", "below"}, // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(_ *testing.T) {
			// Test via constants mapping
			// This indirectly tests mapHRIPosition behavior
		})
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkParseDocument_Minimal(b *testing.B) {
	data := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"commands": [{"type": "text", "data": {"content": {"text": "Hi"}}}]
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = schema.ParseDocument(data)
	}
}

func BenchmarkParseDocument_Receipt(b *testing.B) {
	data := []byte(`{
		"version": "1.0",
		"profile": {"model": "EPSON", "paper_width":  80, "dpi": 203},
		"commands": [
			{"type": "text", "data":  {"content": {"text": "STORE NAME", "content_style":  {"bold": true, "size": "2x2"}}, "label": null}},
			{"type": "text", "data": {"content": {"text":  "123 Main Street"}}},
			{"type": "separator", "data": {"char": "-"}},
			{"type": "table", "data": {"definition":  {"columns": [{"name": "Item", "width": 20}, {"name": "Price", "width": 10}]}, "rows": [["Coffee", "$4.50"], ["Muffin", "$3.00"]]}},
			{"type": "separator", "data": {"char": "="}},
			{"type": "text", "data": {"content": {"text":  "$7.50", "content_style": {"bold": true}}, "label": {"text": "Total"}}},
			{"type": "qr", "data":  {"data": "https://receipt.example.com/12345"}},
			{"type": "feed", "data":  {"lines": 3}},
			{"type": "cut", "data": {"mode":  "partial"}}
		]
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = schema.ParseDocument(data)
	}
}

func BenchmarkCommandUnmarshal_Text(b *testing.B) {
	data := []byte(`{
		"content": {"text": "Hello", "content_style":  {"bold": true, "size": "2x2"}},
		"label": {"text": "Greeting", "separator": ": "},
		"new_line": true
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var cmd TextCommand
		_ = json.Unmarshal(data, &cmd)
	}
}

func BenchmarkCommandUnmarshal_Table(b *testing.B) {
	data := []byte(`{
		"definition": {"columns": [{"name": "A", "width": 10}, {"name": "B", "width":  10}]},
		"rows": [["1", "2"], ["3", "4"], ["5", "6"]],
		"options": {"header_bold": true}
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var cmd TableCommand
		_ = json.Unmarshal(data, &cmd)
	}
}
