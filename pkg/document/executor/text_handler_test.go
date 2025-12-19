package executor

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
)

// ============================================================================
// Text Command Parsing Tests
// ============================================================================

func TestTextCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd TextCommand)
	}{
		{
			name: "simple text",
			json: `{"content": {"text": "Hello World"}}`,
			checkFunc: func(t *testing.T, cmd TextCommand) {
				if cmd.Content.Text != "Hello World" {
					t.Errorf("Expected text 'Hello World', got '%s'", cmd.Content.Text)
				}
			},
		},
		{
			name: "text with bold style",
			json: `{"content": {"text": "Bold", "content_style": {"bold": true}}}`,
			checkFunc: func(t *testing.T, cmd TextCommand) {
				if cmd.Content.Style == nil || cmd.Content.Style.Bold == nil || !*cmd.Content.Style.Bold {
					t.Error("Expected bold to be true")
				}
			},
		},
		{
			name: "text with label",
			json: `{"content": {"text": "$100.00"}, "label": {"text": "Total", "separator": ": "}}`,
			checkFunc: func(t *testing.T, cmd TextCommand) {
				if cmd.Label == nil || cmd.Label.Text != "Total" {
					t.Error("Expected label 'Total'")
				}
				if cmd.Label.Separator == nil || *cmd.Label.Separator != ": " {
					t.Error("Expected separator ': '")
				}
			},
		},
		{
			name: "text with alignment",
			json: `{"content": {"text": "Centered", "align": "center"}}`,
			checkFunc: func(t *testing.T, cmd TextCommand) {
				if cmd.Content.Align == nil || *cmd.Content.Align != constants.Center.String() {
					t.Error("Expected align 'center'")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd TextCommand
			if err := json.Unmarshal([]byte(tt.json), &cmd); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

// ============================================================================
// Text Command Default Value Tests
// ============================================================================

func TestTextCommand_Defaults(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd TextCommand)
	}{
		{
			name: "nil style uses defaults",
			json: `{"content": {"text": "test"}}`,
			checkFunc: func(t *testing.T, cmd TextCommand) {
				if cmd.Content.Style != nil {
					t.Error("Expected style to be nil (handler applies defaults)")
				}
			},
		},
		{
			name: "nil new_line uses default true",
			json: `{"content": {"text": "test"}}`,
			checkFunc: func(t *testing.T, cmd TextCommand) {
				if cmd.NewLine != nil {
					t.Error("Expected new_line to be nil (handler applies default true)")
				}
			},
		},
		{
			name: "nil label is valid",
			json: `{"content": {"text": "test"}}`,
			checkFunc: func(t *testing.T, cmd TextCommand) {
				if cmd.Label != nil {
					t.Error("Expected label to be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd TextCommand
			if err := json.Unmarshal([]byte(tt.json), &cmd); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

// ============================================================================
// Text Command Validation Tests
// ============================================================================

func TestTextCommand_Validation(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		expectError bool
	}{
		{
			name:        "missing content parses but invalid for handler",
			json:        `{"label": {"text": "Label"}}`,
			expectError: false,
		},
		{
			name:        "invalid json",
			json:        `{invalid}`,
			expectError: true,
		},
		{
			name:        "empty object parses",
			json:        `{}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd TextCommand
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

	doc, err := schema.ParseDocument(jsonData)
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

func TestParseDocument_UnicodeContent(t *testing.T) {
	jsonData := []byte(`{
		"version": "1.0",
		"profile": {"model": "Test"},
		"commands": [
			{"type": "text", "data": {"content": {"text": "Héllo Wörld 中文 🎉"}}}
		]
	}`)

	doc, err := schema.ParseDocument(jsonData)
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
