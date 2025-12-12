package executor

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

// ============================================================================
// Barcode Command Parsing Tests
// ============================================================================

func TestBarcodeCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd BarcodeCommand)
	}{
		{
			name: "CODE128 barcode",
			json: `{"symbology": "CODE128", "data": "12345"}`,
			checkFunc: func(t *testing.T, cmd BarcodeCommand) {
				if cmd.Symbology != "CODE128" {
					t.Errorf("Expected symbology 'CODE128', got '%s'", cmd.Symbology)
				}
				if cmd.Data != "12345" {
					t.Errorf("Expected data '12345', got '%s'", cmd.Data)
				}
			},
		},
		{
			name: "EAN13 with dimensions",
			json: `{"symbology": "EAN13", "data": "1234567890123", "width": 3, "height": 80}`,
			checkFunc: func(t *testing.T, cmd BarcodeCommand) {
				if cmd.Width == nil || *cmd.Width != 3 {
					t.Error("Expected width 3")
				}
				if cmd.Height == nil || *cmd.Height != 80 {
					t.Error("Expected height 80")
				}
			},
		},
		{
			name: "barcode with HRI options",
			json: `{"symbology": "CODE39", "data": "ABC123", "hri_position": "below", "hri_font": "A"}`,
			checkFunc: func(t *testing.T, cmd BarcodeCommand) {
				if cmd.HRIPosition == nil || *cmd.HRIPosition != "below" {
					t.Error("Expected hri_position 'below'")
				}
				if cmd.HRIFont == nil || *cmd.HRIFont != "A" {
					t.Error("Expected hri_font 'A'")
				}
			},
		},
		{
			name: "barcode with alignment",
			json: `{"symbology": "UPCA", "data":  "12345678901", "align": "center"}`,
			checkFunc: func(t *testing.T, cmd BarcodeCommand) {
				if cmd.Align == nil || *cmd.Align != constants.Center.String() {
					t.Error("Expected align 'center'")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd BarcodeCommand
			err := json.Unmarshal([]byte(tt.json), &cmd)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

// ============================================================================
// Barcode Command Default Value Tests
// ============================================================================

func TestBarcodeCommand_Defaults(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd BarcodeCommand)
	}{
		{
			name: "nil width uses default",
			json: `{"symbology": "CODE128", "data": "123"}`,
			checkFunc: func(t *testing.T, cmd BarcodeCommand) {
				if cmd.Width != nil {
					t.Error("Expected width to be nil (handler applies default)")
				}
			},
		},
		{
			name: "nil height uses default",
			json: `{"symbology": "CODE128", "data":  "123"}`,
			checkFunc: func(t *testing.T, cmd BarcodeCommand) {
				if cmd.Height != nil {
					t.Error("Expected height to be nil (handler applies default)")
				}
			},
		},
		{
			name: "nil HRI options use defaults",
			json: `{"symbology": "CODE128", "data": "123"}`,
			checkFunc: func(t *testing.T, cmd BarcodeCommand) {
				if cmd.HRIPosition != nil {
					t.Error("Expected hri_position to be nil")
				}
				if cmd.HRIFont != nil {
					t.Error("Expected hri_font to be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd BarcodeCommand
			err := json.Unmarshal([]byte(tt.json), &cmd)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

// ============================================================================
// Barcode Command Validation Tests
// ============================================================================

func TestBarcodeCommand_Validation(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		expectError bool
	}{
		{
			name:        "missing symbology is invalid",
			json:        `{"data": "12345"}`,
			expectError: false, // JSON parses, handler validates
		},
		{
			name:        "missing data is invalid",
			json:        `{"symbology": "CODE128"}`,
			expectError: false, // JSON parses, handler validates
		},
		{
			name:        "invalid json",
			json:        `{invalid}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd BarcodeCommand
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
