package executor

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

// ============================================================================
// Raw Command Parsing Tests
// ============================================================================

func TestRawCommand_Parsing(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd RawCommand)
	}{
		{
			name: "hex format",
			json: `{"hex": "1B 40", "format": "hex"}`,
			checkFunc: func(t *testing.T, cmd RawCommand) {
				if cmd.Hex != "1B 40" {
					t.Errorf("Expected hex '1B 40', got '%s'", cmd.Hex)
				}
				if cmd.Format != "hex" {
					t.Errorf("Expected format 'hex', got '%s'", cmd.Format)
				}
			},
		},
		{
			name: "with comment and safe mode",
			json: `{"hex": "1B40", "comment": "Initialize", "safe_mode": true}`,
			checkFunc: func(t *testing.T, cmd RawCommand) {
				if cmd.Comment != "Initialize" {
					t.Errorf("Expected comment 'Initialize', got '%s'", cmd.Comment)
				}
				if !cmd.SafeMode {
					t.Error("Expected safe_mode to be true")
				}
			},
		},
		{
			name: "base64 format",
			json: `{"hex": "G0A=", "format": "base64"}`,
			checkFunc: func(t *testing.T, cmd RawCommand) {
				if cmd.Format != "base64" {
					t.Errorf("Expected format 'base64', got '%s'", cmd.Format)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd RawCommand
			if err := json.Unmarshal([]byte(tt.json), &cmd); err != nil {
				t.Fatalf("Unexpected error:  %v", err)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, cmd)
			}
		})
	}
}

// ============================================================================
// Raw Command Default Value Tests
// ============================================================================

func TestRawCommand_Defaults(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		checkFunc func(t *testing.T, cmd RawCommand)
	}{
		{
			name: "empty format uses default hex",
			json: `{"hex": "1B40"}`,
			checkFunc: func(t *testing.T, cmd RawCommand) {
				if cmd.Format != "" {
					t.Errorf("Expected empty format (handler applies default %s)", constants.DefaultRawFormat)
				}
			},
		},
		{
			name: "safe_mode defaults to false",
			json: `{"hex": "1B40"}`,
			checkFunc: func(t *testing.T, cmd RawCommand) {
				if cmd.SafeMode {
					t.Error("Expected safe_mode to default to false")
				}
			},
		},
		{
			name: "comment defaults to empty",
			json: `{"hex": "1B40"}`,
			checkFunc: func(t *testing.T, cmd RawCommand) {
				if cmd.Comment != "" {
					t.Error("Expected comment to be empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd RawCommand
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
// Raw Command Validation Tests
// ============================================================================

func TestRawCommand_Validation(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		expectError bool
	}{
		{
			name:        "missing hex parses but invalid for handler",
			json:        `{"format": "hex"}`,
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
			var cmd RawCommand
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
// Hex String Parsing Tests
// ============================================================================

func TestParseHexString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []byte
		expectError bool
	}{
		{"simple hex", "1B40", []byte{0x1B, 0x40}, false},
		{"hex with spaces", "1B 40 0A", []byte{0x1B, 0x40, 0x0A}, false},
		{"hex with 0x prefix", "0x1B0x40", []byte{0x1B, 0x40}, false},
		{"lowercase hex", "1b40", []byte{0x1B, 0x40}, false},
		{"empty string", "", nil, true},
		{"odd characters", "1B4", nil, true},
		{"invalid hex", "GGHI", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseHexString(tt.input)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if len(result) != len(tt.expected) {
				t.Fatalf("Expected %d bytes, got %d", len(tt.expected), len(result))
			}
			for i, b := range tt.expected {
				if result[i] != b {
					t.Errorf("Byte %d:  expected 0x%02X, got 0x%02X", i, b, result[i])
				}
			}
		})
	}
}
