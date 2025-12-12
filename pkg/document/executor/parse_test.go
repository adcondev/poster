package executor

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

// ============================================================================
// Valid JSON Parsing Tests
// ============================================================================

func TestParseDocument_ValidDocument(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {
			"model": "TestPrinter",
			"paper_width": 80,
			"code_table": "WPC1252"
		},
		"commands": [
			{"type": "text", "data": {"content": {"text": "Hello"}}}
		]
	}`)

	doc, err := ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	if doc == nil {
		t.Fatal("ParseDocument() returned nil document")
	}

	if doc.Version != constants.DefaultVersion {
		t.Errorf("Expected version '1.0', got '%s'", doc.Version)
	}

	if doc.Profile.Model != "TestPrinter" {
		t.Errorf("Expected model 'TestPrinter', got '%s'", doc.Profile.Model)
	}

	if len(doc.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(doc.Commands))
	}
}

func TestParseDocument_MultipleCommands(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"commands": [
			{"type": "text", "data": {"content": {"text": "Line 1"}}},
			{"type": "feed", "data": {"lines": 2}},
			{"type": "cut", "data": {"mode": "partial"}}
		]
	}`)

	doc, err := ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	if len(doc.Commands) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(doc.Commands))
	}

	expectedTypes := []string{"text", "feed", "cut"}
	for i, expected := range expectedTypes {
		if doc.Commands[i].Type != expected {
			t.Errorf("Command %d: expected type '%s', got '%s'", i, expected, doc.Commands[i].Type)
		}
	}
}

func TestParseDocument_AllCommandTypes(t *testing.T) {
	commandTypes := []struct {
		cmdType string
		data    string
	}{
		{"text", `{"content": {"text": "test"}}`},
		{"image", `{"code": "base64data"}`},
		{"separator", `{"char": "-"}`},
		{"feed", `{"lines": 1}`},
		{"cut", `{"mode": "partial"}`},
		{"qr", `{"data": "test"}`},
		{"table", `{"definition": {"columns": [{"name": "A", "width": 10}]}, "rows": [["x"]]}`},
		{"barcode", `{"symbology": "CODE128", "data": "123"}`},
		{"raw", `{"hex": "1B40"}`},
		{"pulse", `{"pin": 0, "on_time":  50, "off_time": 100}`}, // ADD THIS
		{"beep", `{"times": 3, "lapse":  2}`},                    // ADD THIS
	}

	for _, tc := range commandTypes {
		t.Run(tc.cmdType, func(t *testing.T) {
			jsonData := []byte(`{
				"version": "1.0",
				"profile": {"model": "Test"},
				"commands": [{"type": "` + tc.cmdType + `", "data": ` + tc.data + `}]
			}`)

			doc, err := ParseDocument(jsonData)
			if err != nil {
				t.Fatalf("ParseDocument() failed for %s: %v", tc.cmdType, err)
			}

			if doc.Commands[0].Type != tc.cmdType {
				t.Errorf("Expected type '%s', got '%s'", tc.cmdType, doc.Commands[0].Type)
			}
		})
	}
}

// ============================================================================
// Empty Version Default Tests
// ============================================================================

func TestParseDocument_EmptyVersionDefaultsTo1_0(t *testing.T) {
	jsonData := []byte(`{
		"version": "",
		"profile": {"model": "Test"},
		"commands": [{"type": "text", "data": {"content": {"text": "test"}}}]
	}`)

	doc, err := ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	if doc.Version != constants.DefaultVersion {
		t.Errorf("Expected default version '1.0', got '%s'", doc.Version)
	}
}

func TestParseDocument_MissingVersionDefaultsTo1_0(t *testing.T) {
	jsonData := []byte(`{
		"profile": {"model": "Test"},
		"commands": [{"type": "text", "data": {"content": {"text": "test"}}}]
	}`)

	doc, err := ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	if doc.Version != constants.DefaultVersion {
		t.Errorf("Expected default version '1.0', got '%s'", doc.Version)
	}
}

// ============================================================================
// Empty Commands Array Error Tests
// ============================================================================

func TestParseDocument_EmptyCommandsArray(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"commands": []
	}`)

	_, err := ParseDocument(jsonData)
	if err == nil {
		t.Fatal("Expected error for empty commands array, got nil")
	}

	if !strings.Contains(err.Error(), "at least one command") {
		t.Errorf("Expected 'at least one command' error, got: %v", err)
	}
}

func TestParseDocument_MissingCommandsField(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"}
	}`)

	_, err := ParseDocument(jsonData)
	if err == nil {
		t.Fatal("Expected error for missing commands, got nil")
	}

	if !strings.Contains(err.Error(), "at least one command") {
		t.Errorf("Expected 'at least one command' error, got: %v", err)
	}
}

func TestParseDocument_NullCommands(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"commands": null
	}`)

	_, err := ParseDocument(jsonData)
	if err == nil {
		t.Fatal("Expected error for null commands, got nil")
	}
}

// ============================================================================
// Malformed JSON Error Handling Tests
// ============================================================================

func TestParseDocument_MalformedJSON(t *testing.T) {
	testCases := []struct {
		name string
		json string
	}{
		{"missing closing brace", `{"version": "1.0"`},
		{"invalid syntax", `{version: "1.0"}`},
		{"trailing comma", `{"version": "1.0",}`},
		{"single quotes", `{'version': '1.0'}`},
		{"unclosed string", `{"version": "1.0}`},
		{"empty input", ``},
		{"just whitespace", `   `},
		{"random text", `not json at all`},
		{"array instead of object", `[1, 2, 3]`},
		{"number", `42`},
		{"boolean", `true`},
		{"null", `null`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseDocument([]byte(tc.json))
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tc.name)
			}
		})
	}
}

func TestParseDocument_InvalidFieldTypes(t *testing.T) {
	testCases := []struct {
		name string
		json string
	}{
		{
			"version as number",
			`{"version": 1.0, "profile": {"model": "Test"}, "commands": [{"type": "text", "data": {}}]}`,
		},
		{
			"commands as object",
			`{"version": "1.0", "profile": {"model": "Test"}, "commands": {"type": "text"}}`,
		},
		{
			"profile as string",
			`{"version": "1.0", "profile": "Test", "commands": [{"type": "text", "data": {}}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseDocument([]byte(tc.json))
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tc.name)
			}
		})
	}
}

// ============================================================================
// Profile Parsing Tests
// ============================================================================

func TestParseDocument_ProfileFields(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {
			"model": "EPSON TM-T20II",
			"paper_width": 80,
			"code_table": "WPC1252",
			"dpi": 203,
			"has_qr": true
		},
		"commands": [{"type": "text", "data": {"content": {"text": "test"}}}]
	}`)

	doc, err := ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	if doc.Profile.Model != "EPSON TM-T20II" {
		t.Errorf("Expected model 'EPSON TM-T20II', got '%s'", doc.Profile.Model)
	}

	if doc.Profile.PaperWidth != 80 {
		t.Errorf("Expected paper_width 80, got %d", doc.Profile.PaperWidth)
	}

	if doc.Profile.CodeTable != "WPC1252" {
		t.Errorf("Expected code_table 'WPC1252', got '%s'", doc.Profile.CodeTable)
	}

	if doc.Profile.DPI != 203 {
		t.Errorf("Expected DPI 203, got %d", doc.Profile.DPI)
	}

	if !doc.Profile.HasQR {
		t.Error("Expected has_qr to be true")
	}
}

func TestParseDocument_MinimalProfile(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"commands": [{"type": "text", "data": {"content": {"text": "test"}}}]
	}`)

	doc, err := ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	if doc.Profile.Model != "Test" {
		t.Errorf("Expected model 'Test', got '%s'", doc.Profile.Model)
	}

	// Optional fields should be zero values
	if doc.Profile.PaperWidth != 0 {
		t.Errorf("Expected default paper_width 0, got %d", doc.Profile.PaperWidth)
	}
}

// ============================================================================
// Command Data Parsing Tests
// ============================================================================

func TestParseDocument_CommandDataPreserved(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"commands": [
			{
				"type": "text",
				"data": {
					"content": {
						"text": "Hello World",
						"content_style": {"bold": true}
					},
					"new_line": true
				}
			}
		]
	}`)

	doc, err := ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	// Verify raw JSON data is preserved
	if doc.Commands[0].Data == nil {
		t.Fatal("Command data is nil")
	}

	// Try to unmarshal the preserved data
	var textCmd TextCommand
	if err := json.Unmarshal(doc.Commands[0].Data, &textCmd); err != nil {
		t.Fatalf("Failed to unmarshal command data: %v", err)
	}

	if textCmd.Content.Text != "Hello World" {
		t.Errorf("Expected text 'Hello World', got '%s'", textCmd.Content.Text)
	}
}

// ============================================================================
// Debug Log Field Tests
// ============================================================================

func TestParseDocument_DebugLogField(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"debug_log": true,
		"commands": [{"type": "text", "data": {"content": {"text": "test"}}}]
	}`)

	doc, err := ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	if !doc.DebugLog {
		t.Error("Expected debug_log to be true")
	}
}

// ============================================================================
// Edge Cases
// ============================================================================

func TestParseDocument_LargeDocument(t *testing.T) {
	// Build a document with many commands
	var commands []string
	for i := 0; i < 100; i++ {
		commands = append(commands, `{"type": "text", "data": {"content": {"text": "Line"}}}`)
	}

	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"commands": [` + strings.Join(commands, ",") + `]
	}`)

	doc, err := ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	if len(doc.Commands) != 100 {
		t.Errorf("Expected 100 commands, got %d", len(doc.Commands))
	}
}

func TestParseDocument_UnicodeContent(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"commands": [
			{"type": "text", "data": {"content": {"text": "Héllo Wörld 中文 🎉"}}}
		]
	}`)

	doc, err := ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	var textCmd TextCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &textCmd)

	expected := "Héllo Wörld 中文 🎉"
	if textCmd.Content.Text != expected {
		t.Errorf("Expected text '%s', got '%s'", expected, textCmd.Content.Text)
	}
}

func TestParseDocument_SpecialCharactersInStrings(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test \"Quoted\""},
		"commands": [
			{"type": "text", "data": {"content": {"text": "Line with\nnewline and\ttab"}}}
		]
	}`)

	doc, err := ParseDocument(jsonData)
	if err != nil {
		t.Fatalf("ParseDocument() error: %v", err)
	}

	if doc.Profile.Model != `Test "Quoted"` {
		t.Errorf("Expected model with quotes, got '%s'", doc.Profile.Model)
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkParseDocument_Simple(b *testing.B) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"commands": [{"type": "text", "data": {"content": {"text": "Hello"}}}]
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseDocument(jsonData)
	}
}

func BenchmarkParseDocument_Complex(b *testing.B) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {
			"model": "EPSON TM-T20II",
			"paper_width": 80,
			"code_table": "WPC1252",
			"dpi": 203,
			"has_qr": true
		},
		"debug_log": true,
		"commands": [
			{"type": "text", "data": {"content": {"text": "Header", "content_style": {"bold": true}}}},
			{"type": "separator", "data": {"char": "="}},
			{"type": "table", "data": {"definition": {"columns": [{"name": "A", "width": 10}]}, "rows": [["x"]]}},
			{"type": "qr", "data": {"data": "https://example.com"}},
			{"type": "cut", "data": {"mode": "partial"}}
		]
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseDocument(jsonData)
	}
}
