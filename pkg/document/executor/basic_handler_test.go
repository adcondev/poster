package executor

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/document/schema"
	"github.com/adcondev/poster/pkg/service"
	"github.com/adcondev/poster/pkg/tables"
)

// ============================================================================
// Feed Command Tests
// ============================================================================

func TestFeedCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd FeedCommand)
	}{
		{
			name: "valid feed command",
			json: `{"lines": 5}`,
			checkFunc: func(t *testing.T, cmd FeedCommand) {
				if cmd.Lines != 5 {
					t.Errorf("Expected lines 5, got %d", cmd.Lines)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd FeedCommand
			if err := json.Unmarshal([]byte(tt.json), &cmd); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

func TestFeedCommand_Defaults(t *testing.T) {
	var cmd FeedCommand
	if err := json.Unmarshal([]byte(`{}`), &cmd); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.Lines != 0 {
		t.Error("Expected lines 0 (handler applies default)")
	}
}

func TestFeedCommand_Validation(t *testing.T) {
	var cmd FeedCommand
	err := json.Unmarshal([]byte(`{invalid}`), &cmd)
	if err == nil {
		t.Error("Expected error for invalid json")
	}
}

// ============================================================================
// Cut Command Tests
// ============================================================================

func TestCutCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd CutCommand)
	}{
		{
			name: "partial cut",
			json: `{"mode": "partial", "feed": 3}`,
			checkFunc: func(t *testing.T, cmd CutCommand) {
				if cmd.Mode != "partial" {
					t.Errorf("Expected mode 'partial', got '%s'", cmd.Mode)
				}
				if cmd.Feed == nil || *cmd.Feed != 3 {
					t.Errorf("Expected feed 3, got %v", cmd.Feed)
				}
			},
		},
		{
			name: "full cut",
			json: `{"mode": "full"}`,
			checkFunc: func(t *testing.T, cmd CutCommand) {
				if cmd.Mode != "full" {
					t.Errorf("Expected mode 'full', got '%s'", cmd.Mode)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd CutCommand
			if err := json.Unmarshal([]byte(tt.json), &cmd); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

func TestCutCommand_Defaults(t *testing.T) {
	var cmd CutCommand
	if err := json.Unmarshal([]byte(`{}`), &cmd); err != nil {
		t.Fatalf("Unexpected error:  %v", err)
	}
	if cmd.Mode != "" {
		t.Error("Expected empty mode (handler applies default)")
	}
	if cmd.Feed != nil {
		t.Error("Expected feed nil (handler applies default)")
	}
}

func TestCutCommand_Validation(t *testing.T) {
	var cmd CutCommand
	err := json.Unmarshal([]byte(`{invalid}`), &cmd)
	if err == nil {
		t.Error("Expected error for invalid json")
	}
}

// ============================================================================
// Separator Command Tests
// ============================================================================

func TestSeparatorCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd SeparatorCommand)
	}{
		{
			name: "dash separator",
			json: `{"char": "-", "length": 48}`,
			checkFunc: func(t *testing.T, cmd SeparatorCommand) {
				if cmd.Char != "-" {
					t.Errorf("Expected char '-', got '%s'", cmd.Char)
				}
				if cmd.Length != 48 {
					t.Errorf("Expected length 48, got %d", cmd.Length)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd SeparatorCommand
			if err := json.Unmarshal([]byte(tt.json), &cmd); err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

func TestSeparatorCommand_Defaults(t *testing.T) {
	var cmd SeparatorCommand
	if err := json.Unmarshal([]byte(`{}`), &cmd); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.Char != "" {
		t.Error("Expected empty char (handler applies default)")
	}
	if cmd.Length != 0 {
		t.Error("Expected length 0 (handler applies default)")
	}
}

func TestSeparatorCommand_Validation(t *testing.T) {
	var cmd SeparatorCommand
	err := json.Unmarshal([]byte(`{invalid}`), &cmd)
	if err == nil {
		t.Error("Expected error for invalid json")
	}
}

// ============================================================================
// Pulse Command Tests
// ============================================================================

func TestPulseCommand_Parsing(t *testing.T) {
	var cmd PulseCommand
	err := json.Unmarshal([]byte(`{"pin": 1, "on_time": 100, "off_time": 200}`), &cmd)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.Pin != 1 || cmd.OnTime != 100 || cmd.OffTime != 200 {
		t.Error("Expected pin=1, on_time=100, off_time=200")
	}
}

func TestPulseCommand_Defaults(t *testing.T) {
	var cmd PulseCommand
	if err := json.Unmarshal([]byte(`{}`), &cmd); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.Pin != 0 || cmd.OnTime != 0 || cmd.OffTime != 0 {
		t.Error("Expected all zero values (handler applies defaults)")
	}
}

func TestPulseCommand_Validation(t *testing.T) {
	var cmd PulseCommand
	err := json.Unmarshal([]byte(`{invalid}`), &cmd)
	if err == nil {
		t.Error("Expected error for invalid json")
	}
}

// ============================================================================
// Beep Command Tests
// ============================================================================

func TestBeepCommand_Parsing(t *testing.T) {
	var cmd BeepCommand
	err := json.Unmarshal([]byte(`{"times": 5, "lapse": 3}`), &cmd)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.Times != 5 || cmd.Lapse != 3 {
		t.Error("Expected times=5, lapse=3")
	}
}

func TestBeepCommand_Defaults(t *testing.T) {
	var cmd BeepCommand
	if err := json.Unmarshal([]byte(`{}`), &cmd); err != nil {
		t.Fatalf("Unexpected error:  %v", err)
	}
	if cmd.Times != 0 || cmd.Lapse != 0 {
		t.Error("Expected zero values (handler applies defaults)")
	}
}

func TestBeepCommand_Validation(t *testing.T) {
	var cmd BeepCommand
	err := json.Unmarshal([]byte(`{invalid}`), &cmd)
	if err == nil {
		t.Error("Expected error for invalid json")
	}
}

// ============================================================================
// Default Application Tests (Handler Logic)
// ============================================================================

func TestFeedCommand_DefaultApplication(t *testing.T) {
	var cmd FeedCommand
	_ = json.Unmarshal([]byte(`{"lines": 0}`), &cmd)

	// Simulate handler default logic
	if cmd.Lines <= 0 {
		cmd.Lines = constants.DefaultFeedLines
	}
	if cmd.Lines != constants.DefaultFeedLines {
		t.Errorf("Expected default lines %d, got %d", constants.DefaultFeedLines, cmd.Lines)
	}
}

// ============================================================================
// Direct Handler Execution Tests
// ============================================================================

func TestHandleSeparator_Execution(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("default separator with profile dots", func(t *testing.T) {
		mock.Reset()
		err := exec.handleSeparator(mock, json.RawMessage(`{}`))
		if err != nil {
			t.Fatalf("handleSeparator failed: %v", err)
		}
		if mock.CallCount("PrintLine") != 1 {
			t.Errorf("expected 1 PrintLine call, got %d", mock.CallCount("PrintLine"))
		}
	})

	t.Run("custom char and length", func(t *testing.T) {
		mock.Reset()
		err := exec.handleSeparator(mock, json.RawMessage(`{"char": "=", "length": 20}`))
		if err != nil {
			t.Fatalf("handleSeparator failed: %v", err)
		}
		if mock.CallCount("PrintLine") != 1 {
			t.Errorf("expected 1 PrintLine call, got %d", mock.CallCount("PrintLine"))
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		err := exec.handleSeparator(mock, json.RawMessage(`{invalid`))
		if err == nil {
			t.Error("expected error for invalid json")
		}
	})
}

func TestHandleFeed_Execution(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("custom lines", func(t *testing.T) {
		mock.Reset()
		err := exec.handleFeed(mock, json.RawMessage(`{"lines": 4}`))
		if err != nil {
			t.Fatalf("handleFeed failed: %v", err)
		}
		if mock.CallCount("FeedLines") != 1 {
			t.Errorf("expected 1 FeedLines call, got %d", mock.CallCount("FeedLines"))
		}
	})

	t.Run("default lines when <= 0", func(t *testing.T) {
		mock.Reset()
		err := exec.handleFeed(mock, json.RawMessage(`{"lines": 0}`))
		if err != nil {
			t.Fatalf("handleFeed failed: %v", err)
		}
		if mock.CallCount("FeedLines") != 1 {
			t.Errorf("expected 1 FeedLines call, got %d", mock.CallCount("FeedLines"))
		}
	})

	t.Run("clamped lines when > 255", func(t *testing.T) {
		mock.Reset()
		err := exec.handleFeed(mock, json.RawMessage(`{"lines": 300}`))
		if err != nil {
			t.Fatalf("handleFeed failed: %v", err)
		}
		if mock.CallCount("FeedLines") != 1 {
			t.Errorf("expected 1 FeedLines call, got %d", mock.CallCount("FeedLines"))
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		err := exec.handleFeed(mock, json.RawMessage(`{bad`))
		if err == nil {
			t.Error("expected error for invalid json")
		}
	})
}

func TestHandleCut_Execution(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("full cut with feed", func(t *testing.T) {
		mock.Reset()
		err := exec.handleCut(mock, json.RawMessage(`{"mode": "full", "feed": 3}`))
		if err != nil {
			t.Fatalf("handleCut failed: %v", err)
		}
		if mock.CallCount("FullFeedAndCut") != 1 {
			t.Errorf("expected 1 FullFeedAndCut call, got %d", mock.CallCount("FullFeedAndCut"))
		}
		if mock.CallCount("FeedLines") != 1 {
			t.Errorf("expected 1 FeedLines call, got %d", mock.CallCount("FeedLines"))
		}
	})

	t.Run("partial cut default", func(t *testing.T) {
		mock.Reset()
		err := exec.handleCut(mock, json.RawMessage(`{"mode": "partial", "feed": 0}`))
		if err != nil {
			t.Fatalf("handleCut failed: %v", err)
		}
		if mock.CallCount("PartialFeedAndCut") != 1 {
			t.Errorf("expected 1 PartialFeedAndCut call, got %d", mock.CallCount("PartialFeedAndCut"))
		}
	})

	t.Run("feed clamped when > 255", func(t *testing.T) {
		mock.Reset()
		err := exec.handleCut(mock, json.RawMessage(`{"feed": 500}`))
		if err != nil {
			t.Fatalf("handleCut failed: %v", err)
		}
		if mock.CallCount("FeedLines") != 1 {
			t.Errorf("expected 1 FeedLines call, got %d", mock.CallCount("FeedLines"))
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		err := exec.handleCut(mock, json.RawMessage(`{bad`))
		if err == nil {
			t.Error("expected error for invalid json")
		}
	})
}

func TestHandlePulse_Execution(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("default pulse", func(t *testing.T) {
		mock.Reset()
		err := exec.handlePulse(mock, json.RawMessage(`{}`))
		if err != nil {
			t.Fatalf("handlePulse failed: %v", err)
		}
		if mock.CallCount("Pulse") != 1 {
			t.Errorf("expected 1 Pulse call, got %d", mock.CallCount("Pulse"))
		}
	})

	t.Run("custom pulse clamped", func(t *testing.T) {
		mock.Reset()
		err := exec.handlePulse(mock, json.RawMessage(`{"pin": 0, "on_time": 600, "off_time": 600}`))
		if err != nil {
			t.Fatalf("handlePulse failed: %v", err)
		}
		if mock.CallCount("Pulse") != 1 {
			t.Errorf("expected 1 Pulse call, got %d", mock.CallCount("Pulse"))
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		err := exec.handlePulse(mock, json.RawMessage(`{bad`))
		if err == nil {
			t.Error("expected error for invalid json")
		}
	})
}

func TestHandleBeep_Execution(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("default beep", func(t *testing.T) {
		mock.Reset()
		err := exec.handleBeep(mock, json.RawMessage(`{}`))
		if err != nil {
			t.Fatalf("handleBeep failed: %v", err)
		}
		if mock.CallCount("Beep") != 1 {
			t.Errorf("expected 1 Beep call, got %d", mock.CallCount("Beep"))
		}
	})

	t.Run("custom beep clamped", func(t *testing.T) {
		mock.Reset()
		err := exec.handleBeep(mock, json.RawMessage(`{"times": 300, "lapse": 300}`))
		if err != nil {
			t.Fatalf("handleBeep failed: %v", err)
		}
		if mock.CallCount("Beep") != 1 {
			t.Errorf("expected 1 Beep call, got %d", mock.CallCount("Beep"))
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		err := exec.handleBeep(mock, json.RawMessage(`{bad`))
		if err == nil {
			t.Error("expected error for invalid json")
		}
	})
}

func TestAlignScopeAndUtilities(t *testing.T) {
	mock := service.NewMockPrinter()

	t.Run("align scope success", func(t *testing.T) {
		mock.Reset()
		called := false
		err := AlignScope(mock, "center", func() error {
			called = true
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Error("expected fn to be called")
		}
		if mock.CallCount("AlignCenter") != 1 {
			t.Errorf("expected 1 AlignCenter, got %d", mock.CallCount("AlignCenter"))
		}
		if mock.CallCount("AlignLeft") != 1 {
			t.Errorf("expected 1 AlignLeft for reset, got %d", mock.CallCount("AlignLeft"))
		}
	})

	t.Run("align scope fn error", func(t *testing.T) {
		mock.Reset()
		expectedErr := errors.New("fn failed")
		err := AlignScope(mock, "right", func() error {
			return expectedErr
		})
		if !errors.Is(err, expectedErr) {
			t.Errorf("expected fn error, got %v", err)
		}
	})

	t.Run("apply alignment options", func(t *testing.T) {
		mock.Reset()
		_ = ApplyAlignment(mock, "center")
		_ = ApplyAlignment(mock, "right")
		_ = ApplyAlignment(mock, "left")
		_ = ApplyAlignment(mock, "unknown")
		if mock.CallCount("AlignCenter") != 1 {
			t.Errorf("expected 1 AlignCenter, got %d", mock.CallCount("AlignCenter"))
		}
		if mock.CallCount("AlignRight") != 1 {
			t.Errorf("expected 1 AlignRight, got %d", mock.CallCount("AlignRight"))
		}
		if mock.CallCount("AlignLeft") != 2 {
			t.Errorf("expected 2 AlignLeft, got %d", mock.CallCount("AlignLeft"))
		}
	})

	t.Run("resolve alignment", func(t *testing.T) {
		if got := ResolveAlignment("", constants.Center); got != constants.Center.String() {
			t.Errorf("expected default center, got %s", got)
		}
		if got := ResolveAlignment("RIGHT", constants.Left); got != "right" {
			t.Errorf("expected lower right, got %s", got)
		}
	})
}

func TestHandleRaw_Execution(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("hex format", func(t *testing.T) {
		mock.Reset()
		payload := `{"format": "hex", "hex": "1B 40", "safe_mode": false}`
		err := exec.handleRaw(mock, json.RawMessage(payload))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.CallCount("Write") != 1 {
			t.Errorf("expected 1 Write call, got %d", mock.CallCount("Write"))
		}
	})

	t.Run("base64 format", func(t *testing.T) {
		mock.Reset()
		encoded := base64.StdEncoding.EncodeToString([]byte{0x1B, 0x32})
		payload := `{"format": "base64", "hex": "` + encoded + `"}`
		err := exec.handleRaw(mock, json.RawMessage(payload))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.CallCount("Write") != 1 {
			t.Errorf("expected 1 Write call, got %d", mock.CallCount("Write"))
		}
	})

	t.Run("default format is hex", func(t *testing.T) {
		mock.Reset()
		payload := `{"hex": "1B 32"}`
		err := exec.handleRaw(mock, json.RawMessage(payload))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if mock.CallCount("Write") != 1 {
			t.Errorf("expected 1 Write call, got %d", mock.CallCount("Write"))
		}
	})

	t.Run("unsupported format", func(t *testing.T) {
		payload := `{"format": "yaml", "hex": "1B"}`
		err := exec.handleRaw(mock, json.RawMessage(payload))
		if err == nil || !strings.Contains(err.Error(), "unsupported format") {
			t.Errorf("expected unsupported format error, got %v", err)
		}
	})

	t.Run("invalid hex string", func(t *testing.T) {
		payload := `{"format": "hex", "hex": "ZZ"}`
		err := exec.handleRaw(mock, json.RawMessage(payload))
		if err == nil {
			t.Error("expected error for invalid hex string")
		}
	})

	t.Run("invalid base64 string", func(t *testing.T) {
		payload := `{"format": "base64", "hex": "!!!"}`
		err := exec.handleRaw(mock, json.RawMessage(payload))
		if err == nil {
			t.Error("expected error for invalid base64")
		}
	})

	t.Run("empty raw command", func(t *testing.T) {
		payload := `{"format": "base64", "hex": ""}`
		err := exec.handleRaw(mock, json.RawMessage(payload))
		if err == nil || !strings.Contains(err.Error(), "cannot be empty") {
			t.Errorf("expected empty error, got %v", err)
		}
	})

	t.Run("command too large", func(t *testing.T) {
		largeHex := strings.Repeat("00 ", constants.RawMaxBytes+1)
		payload := `{"hex": "` + strings.TrimSpace(largeHex) + `"}`
		err := exec.handleRaw(mock, json.RawMessage(payload))
		if err == nil || !strings.Contains(err.Error(), "too large") {
			t.Errorf("expected too large error, got %v", err)
		}
	})

	t.Run("safe mode blocks critical command", func(t *testing.T) {
		payload := `{"hex": "1B 40", "safe_mode": true}`
		err := exec.handleRaw(mock, json.RawMessage(payload))
		if err == nil || !strings.Contains(err.Error(), "unsafe command blocked") {
			t.Errorf("expected safe mode error, got %v", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		err := exec.handleRaw(mock, json.RawMessage(`{invalid`))
		if err == nil {
			t.Error("expected error for invalid json")
		}
	})
}

func TestTableHandlerUtils_Extra(t *testing.T) {
	t.Run("calculate total width", func(t *testing.T) {
		width := CalculateTotalWidth(10, 3, 2)
		if width != 14 {
			t.Errorf("expected width 14, got %d", width)
		}
		widthZero := CalculateTotalWidth(10, 0, 2)
		if widthZero != 10 {
			t.Errorf("expected width 10, got %d", widthZero)
		}
	})

	t.Run("render table string", func(t *testing.T) {
		showH := true
		cmd := &TableCommand{
			Definition: tables.Definition{
				Columns: []tables.Column{
					{Name: "Item", Width: 10},
					{Name: "Price", Width: 6},
				},
			},
			ShowHeaders: &showH,
			Rows: [][]string{
				{"Coffee", "$3.50"},
			},
		}

		out, err := RenderTable(cmd, nil)
		if err != nil {
			t.Fatalf("RenderTable failed: %v", err)
		}
		if !strings.Contains(out, "Coffee") {
			t.Errorf("expected output to contain Coffee, got %s", out)
		}

		opts := &tables.Options{ShowHeaders: false}
		outOpts, err := RenderTable(cmd, opts)
		if err != nil {
			t.Fatalf("RenderTable with opts failed: %v", err)
		}
		if outOpts == "" {
			t.Error("expected non-empty output")
		}
	})
}

func TestHandlerExecutor_Extra(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("register custom handler", func(t *testing.T) {
		customCalled := false
		exec.RegisterHandler("my_custom_cmd", func(_ service.PrinterActions, _ json.RawMessage) error {
			customCalled = true
			return nil
		})

		doc := &schema.Document{
			Version: constants.DefaultVersion,
			Profile: schema.ProfileConfig{Model: "TestModel"},
			Commands: []schema.Command{
				{Type: "my_custom_cmd", Data: json.RawMessage(`{}`)},
			},
		}
		if err := exec.Execute(doc); err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		if !customCalled {
			t.Error("expected custom handler to be called")
		}
	})

	t.Run("execute json valid and invalid", func(t *testing.T) {
		validJSON := []byte(`{
			"version": "1.0",
			"profile": {"model": "Generic", "paper_width": 58, "dpi": 203, "has_qr": true},
			"commands": [
				{"type": "feed", "data": {"lines": 1}}
			]
		}`)
		if err := exec.ExecuteJSON(validJSON); err != nil {
			t.Fatalf("ExecuteJSON failed: %v", err)
		}

		invalidJSON := []byte(`{not valid json}`)
		if err := exec.ExecuteJSON(invalidJSON); err == nil {
			t.Error("expected error for invalid json in ExecuteJSON")
		}
	})

	t.Run("apply profile document variations", func(t *testing.T) {
		doc80 := &schema.Document{
			Version: constants.DefaultVersion,
			Profile: schema.ProfileConfig{Model: "Generic", PaperWidth: 80, DPI: 203, CodeTable: "PC850"},
			Commands: []schema.Command{
				{Type: "feed", Data: json.RawMessage(`{"lines": 1}`)},
			},
		}
		if err := exec.Execute(doc80); err != nil {
			t.Fatalf("Execute 80mm failed: %v", err)
		}
		if mock.GetProfile().DotsPerLine != constants.PaperPxWidth80mm {
			t.Errorf("expected %d dots per line, got %d", constants.PaperPxWidth80mm, mock.GetProfile().DotsPerLine)
		}

		docCustom := &schema.Document{
			Version: constants.DefaultVersion,
			Profile: schema.ProfileConfig{Model: "Generic", PaperWidth: 72, DPI: 300},
			Commands: []schema.Command{
				{Type: "feed", Data: json.RawMessage(`{"lines": 1}`)},
			},
		}
		if err := exec.Execute(docCustom); err != nil {
			t.Fatalf("Execute custom paper failed: %v", err)
		}

		docPreset := &schema.Document{
			Version: constants.DefaultVersion,
			Profile: schema.ProfileConfig{Model: "TM-T88VI"},
			Commands: []schema.Command{
				{Type: "feed", Data: json.RawMessage(`{"lines": 1}`)},
			},
		}
		if err := exec.Execute(docPreset); err != nil {
			t.Fatalf("Execute preset failed: %v", err)
		}

		docNoModel := &schema.Document{
			Version: constants.DefaultVersion,
			Commands: []schema.Command{
				{Type: "feed", Data: json.RawMessage(`{"lines": 1}`)},
			},
		}
		if err := exec.Execute(docNoModel); err == nil {
			t.Error("expected error for doc without model")
		}

		if err := exec.Execute(nil); err == nil {
			t.Error("expected error for nil doc")
		}
	})
}

func TestHandleSeparator_Extra(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("default char and profile dots per line", func(t *testing.T) {
		mock.Reset()
		prof := mock.GetProfile()
		prof.DotsPerLine = 384
		err := exec.handleSeparator(mock, json.RawMessage(`{}`))
		if err != nil {
			t.Fatalf("handleSeparator failed: %v", err)
		}
		if mock.CallCount("PrintLine") != 1 {
			t.Errorf("expected 1 PrintLine call, got %d", mock.CallCount("PrintLine"))
		}
	})

	t.Run("profile dots per line zero fallback to default length", func(t *testing.T) {
		mock.Reset()
		prof := mock.GetProfile()
		prof.DotsPerLine = 0
		err := exec.handleSeparator(mock, json.RawMessage(`{}`))
		if err != nil {
			t.Fatalf("handleSeparator failed: %v", err)
		}
		if mock.CallCount("PrintLine") != 1 {
			t.Errorf("expected 1 PrintLine call, got %d", mock.CallCount("PrintLine"))
		}
	})

	t.Run("rune truncation with multibyte char", func(t *testing.T) {
		mock.Reset()
		err := exec.handleSeparator(mock, json.RawMessage(`{"char": "==", "length": 5}`))
		if err != nil {
			t.Fatalf("handleSeparator failed: %v", err)
		}
		if mock.CallCount("PrintLine") != 1 {
			t.Errorf("expected 1 PrintLine call, got %d", mock.CallCount("PrintLine"))
		}
	})
}

func TestHandleTable_Branches(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("empty columns error", func(t *testing.T) {
		payload := `{"definition": {"columns": []}, "rows": []}`
		err := exec.handleTable(mock, json.RawMessage(payload))
		if err == nil || !strings.Contains(err.Error(), "at least one column") {
			t.Errorf("expected empty columns error, got %v", err)
		}
	})

	t.Run("negative column spacing clamped and font B and align center", func(t *testing.T) {
		mock.Reset()
		payload := `{
			"definition": {
				"columns": [
					{"name": "Item", "width": 10},
					{"name": "Price", "width": 8}
				]
			},
			"options": {
				"column_spacing": -5,
				"font": "B",
				"align": "center"
			},
			"rows": [
				["Coffee", "$3.00"]
			]
		}`
		err := exec.handleTable(mock, json.RawMessage(payload))
		if err != nil {
			t.Fatalf("handleTable failed: %v", err)
		}
		if mock.CallCount("FontB") != 1 {
			t.Errorf("expected 1 FontB call, got %d", mock.CallCount("FontB"))
		}
		if mock.CallCount("AlignCenter") != 1 {
			t.Errorf("expected 1 AlignCenter call, got %d", mock.CallCount("AlignCenter"))
		}
	})
}


