package executor

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

// ============================================================================
// Image Command Parsing Tests
// ============================================================================

func TestImageCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd ImageCommand)
	}{
		{
			name: "minimal image",
			json: `{"code":  "base64imagedata"}`,
			checkFunc: func(t *testing.T, cmd ImageCommand) {
				if cmd.Code != "base64imagedata" {
					t.Errorf("Expected code 'base64imagedata', got '%s'", cmd.Code)
				}
			},
		},
		{
			name: "image with all options",
			json: `{"code": "data", "pixel_width": 256, "align": "center", "threshold": 128, "dithering": "atkinson", "scaling": "bilinear"}`,
			checkFunc: func(t *testing.T, cmd ImageCommand) {
				if cmd.PixelWidth != 256 {
					t.Errorf("Expected pixel_width 256, got %d", cmd.PixelWidth)
				}
				if cmd.Align != constants.Center.String() {
					t.Errorf("Expected align 'center', got '%s'", cmd.Align)
				}
				if cmd.Threshold != 128 {
					t.Errorf("Expected threshold 128, got %d", cmd.Threshold)
				}
				if cmd.Dithering != "atkinson" {
					t.Errorf("Expected dithering 'atkinson', got '%s'", cmd.Dithering)
				}
				if cmd.Scaling != "bilinear" {
					t.Errorf("Expected scaling 'bilinear', got '%s'", cmd.Scaling)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd ImageCommand
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
// Image Command Default Value Tests
// ============================================================================

func TestImageCommand_Defaults(t *testing.T) {
	var cmd ImageCommand
	if err := json.Unmarshal([]byte(`{"code": "data"}`), &cmd); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cmd.PixelWidth != 0 {
		t.Error("Expected pixel_width 0 (handler applies default)")
	}
	if cmd.Threshold != 0 {
		t.Error("Expected threshold 0 (handler applies default)")
	}
	if cmd.Dithering != "" {
		t.Error("Expected empty dithering (handler applies default)")
	}
	if cmd.Scaling != "" {
		t.Error("Expected empty scaling (handler applies default)")
	}
}

// ============================================================================
// Image Command Validation Tests
// ============================================================================

func TestImageCommand_Validation(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		expectError bool
	}{
		{"missing code parses", `{"pixel_width": 256}`, false},
		{"invalid json", `{invalid}`, true},
		{"empty object parses", `{}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd ImageCommand
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
