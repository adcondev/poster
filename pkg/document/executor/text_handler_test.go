package executor

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/service"
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

// ============================================================================
// Direct Handler Execution Tests
// ============================================================================

func TestHandleText_BasicAndStyles(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("simple text print and printline", func(t *testing.T) {
		mock.Reset()
		payload := `{"content": {"text": "Line 1"}}`
		if err := exec.handleText(mock, json.RawMessage(payload)); err != nil {
			t.Fatalf("handleText failed: %v", err)
		}
		if mock.CallCount("PrintLine") != 1 {
			t.Errorf("expected 1 PrintLine call, got %d", mock.CallCount("PrintLine"))
		}

		mock.Reset()
		payloadNoNL := `{"content": {"text": "No newline"}, "new_line": false}`
		if err := exec.handleText(mock, json.RawMessage(payloadNoNL)); err != nil {
			t.Fatalf("handleText failed: %v", err)
		}
		if mock.CallCount("Print") != 1 {
			t.Errorf("expected 1 Print call, got %d", mock.CallCount("Print"))
		}
	})

	t.Run("bold, inverse, underline styles", func(t *testing.T) {
		mock.Reset()
		payload := `{
			"content": {
				"text": "Styled",
				"content_style": {
					"bold": true,
					"inverse": true,
					"underline": "1pt"
				},
				"align": "center"
			}
		}`
		if err := exec.handleText(mock, json.RawMessage(payload)); err != nil {
			t.Fatalf("handleText failed: %v", err)
		}
		if mock.CallCount("EnableBold") != 1 {
			t.Errorf("expected 1 EnableBold call, got %d", mock.CallCount("EnableBold"))
		}
		if mock.CallCount("DisableBold") != 1 {
			t.Errorf("expected 1 DisableBold call on reset, got %d", mock.CallCount("DisableBold"))
		}
		if mock.CallCount("InverseOn") != 1 {
			t.Errorf("expected 1 InverseOn call, got %d", mock.CallCount("InverseOn"))
		}
		if mock.CallCount("InverseOff") != 1 {
			t.Errorf("expected 1 InverseOff call on reset, got %d", mock.CallCount("InverseOff"))
		}
		if mock.CallCount("OneDot") != 1 {
			t.Errorf("expected 1 OneDot call, got %d", mock.CallCount("OneDot"))
		}
		if mock.CallCount("NoDot") != 1 {
			t.Errorf("expected 1 NoDot call on reset, got %d", mock.CallCount("NoDot"))
		}
	})
}

func TestHandleText_SizesAndFonts(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	sizes := []string{
		"normal", "double", "triple", "quad", "penta", "hexa", "hepta", "octa",
		"2x2", "3x4", "8x8", "invalid_size",
	}
	for _, size := range sizes {
		mock.Reset()
		payload := `{"content": {"text": "SizeTest", "content_style": {"size": "` + size + `"}}}`
		if err := exec.handleText(mock, json.RawMessage(payload)); err != nil {
			t.Errorf("handleText failed for size=%s: %v", size, err)
		}
	}

	underlines := []string{"0", "none", "1", "2"}
	for _, ul := range underlines {
		mock.Reset()
		payload := `{"content": {"text": "ULTest", "content_style": {"underline": "` + ul + `"}}}`
		if err := exec.handleText(mock, json.RawMessage(payload)); err != nil {
			t.Errorf("handleText failed for underline=%s: %v", ul, err)
		}
	}

	fonts := []string{"a", "b", "A", "unknown"}
	for _, font := range fonts {
		mock.Reset()
		payload := `{"content": {"text": "FontTest", "content_style": {"font": "` + font + `"}}}`
		if err := exec.handleText(mock, json.RawMessage(payload)); err != nil {
			t.Errorf("handleText failed for font=%s: %v", font, err)
		}
	}
}

func TestHandleText_LabelsAndDifferingStyles(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("label with different align prints on separate lines", func(t *testing.T) {
		mock.Reset()
		payload := `{
			"label": {
				"text": "Item",
				"separator": ": ",
				"align": "left",
				"label_style": {"bold": true}
			},
			"content": {
				"text": "$10.00",
				"align": "right",
				"content_style": {"bold": false}
			}
		}`
		if err := exec.handleText(mock, json.RawMessage(payload)); err != nil {
			t.Fatalf("handleText failed: %v", err)
		}
		if mock.CallCount("PrintLine") != 2 {
			t.Errorf("expected 2 PrintLine calls for separate alignments, got %d", mock.CallCount("PrintLine"))
		}
	})

	t.Run("label with same align inline and differing styles", func(t *testing.T) {
		mock.Reset()
		payload := `{
			"label": {
				"text": "Status",
				"align": "left",
				"label_style": {"bold": true, "size": "double"}
			},
			"content": {
				"text": "OK",
				"align": "left",
				"content_style": {"bold": false, "size": "normal"}
			}
		}`
		if err := exec.handleText(mock, json.RawMessage(payload)); err != nil {
			t.Fatalf("handleText failed: %v", err)
		}
		if mock.CallCount("Print") != 1 {
			t.Errorf("expected 1 Print call for inline label, got %d", mock.CallCount("Print"))
		}
		if mock.CallCount("PrintLine") != 1 {
			t.Errorf("expected 1 PrintLine call for content, got %d", mock.CallCount("PrintLine"))
		}
	})
}

func TestHandleText_Errors(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("invalid json returns error", func(t *testing.T) {
		err := exec.handleText(mock, json.RawMessage(`{not-json`))
		if err == nil {
			t.Error("expected error for invalid json, got nil")
		}
	})

	t.Run("nil printer in applyTextStyle returns error", func(t *testing.T) {
		b := true
		style := &TextStyle{Bold: &b}
		if err := exec.applyTextStyle(nil, style); err == nil {
			t.Error("expected error for nil printer in applyTextStyle")
		}
		if err := exec.resetTextStyle(nil, style); err == nil {
			t.Error("expected error for nil printer in resetTextStyle")
		}
	})
}

func TestHandleText_LabelCustomSeparatorAndEmptyContent(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("label with custom separator and no content", func(t *testing.T) {
		mock.Reset()
		payload := `{
			"label": {
				"text": "Item",
				"separator": " -> "
			},
			"content": {
				"text": ""
			}
		}`
		if err := exec.handleText(mock, json.RawMessage(payload)); err != nil {
			t.Fatalf("handleText failed: %v", err)
		}
		if mock.CallCount("Print") != 1 {
			t.Errorf("expected 1 Print call for label without content, got %d", mock.CallCount("Print"))
		}
	})

	t.Run("content without newline", func(t *testing.T) {
		mock.Reset()
		payload := `{
			"content": {
				"text": "Inline content"
			},
			"new_line": false
		}`
		if err := exec.handleText(mock, json.RawMessage(payload)); err != nil {
			t.Fatalf("handleText failed: %v", err)
		}
		if mock.CallCount("Print") != 1 {
			t.Errorf("expected 1 Print call for new_line false, got %d", mock.CallCount("Print"))
		}
	})

	t.Run("content with style", func(t *testing.T) {
		mock.Reset()
		payload := `{
			"content": {
				"text": "Styled text",
				"content_style": {"bold": true, "underline": "1pt"}
			}
		}`
		if err := exec.handleText(mock, json.RawMessage(payload)); err != nil {
			t.Fatalf("handleText failed: %v", err)
		}
		if mock.CallCount("EnableBold") != 1 {
			t.Errorf("expected 1 EnableBold, got %d", mock.CallCount("EnableBold"))
		}
		if mock.CallCount("DisableBold") != 1 {
			t.Errorf("expected 1 DisableBold, got %d", mock.CallCount("DisableBold"))
		}
	})
}

func TestApplyTextStyle_Sizes(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	sizes := []string{
		constants.Normal.String(),
		constants.Double.String(),
		constants.Triple.String(),
		constants.Quad.String(),
		constants.Penta.String(),
		constants.Hexa.String(),
		constants.Hepta.String(),
		constants.Octa.String(),
		"2x3",
		"9x9",
		"invalid",
	}

	for _, sz := range sizes {
		t.Run(sz, func(t *testing.T) {
			mock.Reset()
			style := &TextStyle{Size: &sz}
			if err := exec.applyTextStyle(mock, style); err != nil {
				t.Fatalf("applyTextStyle failed for size %s: %v", sz, err)
			}
		})
	}
}

func TestApplyTextStyle_UnderlineFontAlignInverse(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("underlines", func(t *testing.T) {
		for _, u := range []string{constants.NoDot.String(), constants.OneDot.String(), constants.TwoDot.String(), "unknown"} {
			mock.Reset()
			style := &TextStyle{Underline: &u}
			if err := exec.applyTextStyle(mock, style); err != nil {
				t.Fatalf("applyTextStyle underline %s failed: %v", u, err)
			}
		}
	})

	t.Run("fonts", func(t *testing.T) {
		for _, f := range []string{"", "a", "b", "c"} {
			mock.Reset()
			style := &TextStyle{Font: &f}
			if err := exec.applyTextStyle(mock, style); err != nil {
				t.Fatalf("applyTextStyle font %s failed: %v", f, err)
			}
		}
	})

	t.Run("inverses", func(t *testing.T) {
		for _, inv := range []bool{true, false} {
			mock.Reset()
			style := &TextStyle{Inverse: &inv}
			if err := exec.applyTextStyle(mock, style); err != nil {
				t.Fatalf("applyTextStyle inverse %v failed: %v", inv, err)
			}
		}
	})

	t.Run("aligns", func(t *testing.T) {
		for _, a := range []*string{nil, strPtr("left"), strPtr("center"), strPtr("right"), strPtr("other")} {
			mock.Reset()
			if err := exec.applyAlign(mock, a); err != nil {
				t.Fatalf("applyAlign failed: %v", err)
			}
		}
	})
}

func TestResetDifferingStyles_AllDifferences(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("nil label style", func(t *testing.T) {
		if err := exec.resetDifferingStyles(mock, nil, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("nil content style delegates to resetTextStyle", func(t *testing.T) {
		mock.Reset()
		b := true
		labelStyle := &TextStyle{Bold: &b}
		if err := exec.resetDifferingStyles(mock, labelStyle, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if mock.CallCount("DisableBold") != 1 {
			t.Errorf("expected 1 DisableBold, got %d", mock.CallCount("DisableBold"))
		}
	})

	t.Run("all differing fields", func(t *testing.T) {
		mock.Reset()
		bTrue := true
		bFalse := false
		szDouble := "double"
		szNormal := "normal"
		u1 := "1pt"
		uNone := "none"
		invTrue := true
		invFalse := false
		fontB := "B"
		fontA := "A"

		label := &TextStyle{
			Bold:      &bTrue,
			Size:      &szDouble,
			Underline: &u1,
			Inverse:   &invTrue,
			Font:      &fontB,
		}
		content := &TextStyle{
			Bold:      &bFalse,
			Size:      &szNormal,
			Underline: &uNone,
			Inverse:   &invFalse,
			Font:      &fontA,
		}

		if err := exec.resetDifferingStyles(mock, label, content); err != nil {
			t.Fatalf("resetDifferingStyles failed: %v", err)
		}
		if mock.CallCount("DisableBold") != 1 {
			t.Errorf("expected 1 DisableBold, got %d", mock.CallCount("DisableBold"))
		}
		if mock.CallCount("SingleSize") != 1 {
			t.Errorf("expected 1 SingleSize, got %d", mock.CallCount("SingleSize"))
		}
		if mock.CallCount("NoDot") != 1 {
			t.Errorf("expected 1 NoDot, got %d", mock.CallCount("NoDot"))
		}
		if mock.CallCount("InverseOff") != 1 {
			t.Errorf("expected 1 InverseOff, got %d", mock.CallCount("InverseOff"))
		}
		if mock.CallCount("FontA") != 1 {
			t.Errorf("expected 1 FontA, got %d", mock.CallCount("FontA"))
		}
	})
}

func TestResetTextStyle_FullReset(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("nil style", func(t *testing.T) {
		if err := exec.resetTextStyle(mock, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("all fields populated", func(t *testing.T) {
		mock.Reset()
		b := true
		sz := "double"
		u := "1pt"
		inv := true
		f := "B"
		style := &TextStyle{
			Bold:      &b,
			Size:      &sz,
			Underline: &u,
			Inverse:   &inv,
			Font:      &f,
		}

		if err := exec.resetTextStyle(mock, style); err != nil {
			t.Fatalf("resetTextStyle failed: %v", err)
		}
		if mock.CallCount("DisableBold") != 1 {
			t.Errorf("expected 1 DisableBold, got %d", mock.CallCount("DisableBold"))
		}
		if mock.CallCount("SingleSize") != 1 {
			t.Errorf("expected 1 SingleSize, got %d", mock.CallCount("SingleSize"))
		}
		if mock.CallCount("NoDot") != 1 {
			t.Errorf("expected 1 NoDot, got %d", mock.CallCount("NoDot"))
		}
		if mock.CallCount("InverseOff") != 1 {
			t.Errorf("expected 1 InverseOff, got %d", mock.CallCount("InverseOff"))
		}
		if mock.CallCount("FontA") != 1 {
			t.Errorf("expected 1 FontA, got %d", mock.CallCount("FontA"))
		}
	})
}
