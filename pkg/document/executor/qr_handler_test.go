package executor

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/service"
)

// ============================================================================
// QR Command Parsing Tests
// ============================================================================

func TestQRCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd QRCommand)
	}{
		{
			name: "minimal QR",
			json: `{"data": "https://example.com"}`,
			checkFunc: func(t *testing.T, cmd QRCommand) {
				if cmd.Data != "https://example.com" {
					t.Errorf("Expected data 'https://example.com', got '%s'", cmd.Data)
				}
			},
		},
		{
			name: "QR with all options",
			json: `{"data": "test", "pixel_width": 200, "correction": "H", "align": "center", "human_text": "Scan me"}`,
			checkFunc: func(t *testing.T, cmd QRCommand) {
				if cmd.PixelWidth != 200 {
					t.Errorf("Expected pixel_width 200, got %d", cmd.PixelWidth)
				}
				if cmd.Correction != "H" {
					t.Errorf("Expected correction 'H', got '%s'", cmd.Correction)
				}
				if cmd.Align != constants.Center.String() {
					t.Errorf("Expected align 'center', got '%s'", cmd.Align)
				}
				if cmd.HumanText != "Scan me" {
					t.Errorf("Expected human_text 'Scan me', got '%s'", cmd.HumanText)
				}
			},
		},
		{
			name: "QR with circle shape",
			json: `{"data": "test", "circle_shape": true}`,
			checkFunc: func(t *testing.T, cmd QRCommand) {
				if !cmd.CircleShape {
					t.Error("Expected circle_shape to be true")
				}
			},
		},
		{
			name: "QR with logo",
			json: `{"data": "test", "logo":  "base64logodata"}`,
			checkFunc: func(t *testing.T, cmd QRCommand) {
				if cmd.Logo != "base64logodata" {
					t.Errorf("Expected logo 'base64logodata', got '%s'", cmd.Logo)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd QRCommand
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
// QR Command Default Value Tests
// ============================================================================

func TestQRCommand_Defaults(t *testing.T) {
	var cmd QRCommand
	if err := json.Unmarshal([]byte(`{"data": "test"}`), &cmd); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.PixelWidth != 0 {
		t.Error("Expected pixel_width 0 (handler applies default)")
	}
	if cmd.Correction != "" {
		t.Error("Expected empty correction (handler applies default)")
	}
	if cmd.CircleShape {
		t.Error("Expected circle_shape to be false")
	}
	if cmd.Logo != "" {
		t.Error("Expected empty logo")
	}
	if cmd.HumanText != "" {
		t.Error("Expected empty human_text")
	}
}

// ============================================================================
// QR Command Validation Tests
// ============================================================================

func TestQRCommand_Validation(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		expectError bool
	}{
		{"missing data parses", `{"pixel_width": 200}`, false},
		{"invalid json", `{invalid}`, true},
		{"empty object parses", `{}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd QRCommand
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
// Direct Handler Execution Tests
// ============================================================================

func TestHandleQR_Execution(t *testing.T) {
	mock := service.NewMockPrinter()
	exec := NewExecutor(mock)

	t.Run("empty data returns error", func(t *testing.T) {
		err := exec.handleQR(mock, json.RawMessage(`{"data": ""}`))
		if err == nil {
			t.Error("expected error for empty QR data")
		}
	})

	t.Run("invalid json returns error", func(t *testing.T) {
		err := exec.handleQR(mock, json.RawMessage(`{bad-json`))
		if err == nil {
			t.Error("expected error for invalid json")
		}
	})

	t.Run("invalid correction level returns error", func(t *testing.T) {
		err := exec.handleQR(mock, json.RawMessage(`{"data": "test", "correction": "X"}`))
		if err == nil {
			t.Error("expected error for invalid correction level")
		}
	})

	t.Run("valid QR with all options, human text, circle shape, logo", func(t *testing.T) {
		mock.Reset()
		jsonStr := `{
			"data": "https://example.com",
			"human_text": "Scan me",
			"pixel_width": 300,
			"correction": "H",
			"align": "right",
			"logo": "dGVzdA==",
			"circle_shape": true
		}`
		err := exec.handleQR(mock, json.RawMessage(jsonStr))
		if err != nil {
			t.Fatalf("handleQR failed: %v", err)
		}
		if mock.CallCount("PrintQR") != 1 {
			t.Errorf("expected 1 PrintQR call, got %d", mock.CallCount("PrintQR"))
		}
		if mock.CallCount("PrintLine") != 1 {
			t.Errorf("expected 1 PrintLine call for human text, got %d", mock.CallCount("PrintLine"))
		}
		if mock.CallCount("AlignRight") != 1 {
			t.Errorf("expected 1 AlignRight call, got %d", mock.CallCount("AlignRight"))
		}
	})

	t.Run("various correction levels and alignments", func(t *testing.T) {
		corrections := []string{"L", "M", "Q", "H"}
		alignments := []string{"left", "center", "right", "default"}
		for _, corr := range corrections {
			for _, align := range alignments {
				mock.Reset()
				payload := `{"data": "test", "correction": "` + corr + `", "align": "` + align + `", "pixel_width": 50}`
				if err := exec.handleQR(mock, json.RawMessage(payload)); err != nil {
					t.Errorf("handleQR failed for corr=%s align=%s: %v", corr, align, err)
				}
			}
		}
	})
}
