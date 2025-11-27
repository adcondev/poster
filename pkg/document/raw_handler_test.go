package document_test

import (
	"encoding/json"
	"strings"
	"testing"

	tu "github.com/adcondev/pos-printer/internal/testutils"
	"github.com/adcondev/pos-printer/pkg/document"
	"github.com/adcondev/pos-printer/pkg/service"
)

// ============================================================================
// Test Helpers
// ============================================================================

// newMockPrinterForRaw crea un printer mock para testing
func newMockPrinterForRaw(t *testing.T) *service.Printer {
	t.Helper()

	conn := &tu.MockConnector{
		WriteFunc: func(data []byte) (int, error) {
			return len(data), nil
		},
	}

	proto := tu.NewTestProtocol()
	profile := tu.NewTestProfile()

	printer, err := service.NewPrinter(proto, profile, conn)
	if err != nil {
		t.Fatalf("Failed to create mock printer: %v", err)
	}

	return printer
}

// ============================================================================
// Tests para CleanHexString
// ============================================================================

func TestCleanHexString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple hex with spaces",
			input:    "1B 40",
			expected: "1B40",
		},
		{
			name:     "Hex with 0x prefix at start",
			input:    "0x1B40",
			expected: "1B40",
		},
		{
			name:     "Hex with multiple 0x prefixes",
			input:    "0x1B 0x40",
			expected: "1B40",
		},
		{
			name:     "Hex with commas",
			input:    "1B,40,0A",
			expected: "1B400A",
		},
		{
			name:     "Hex with colons",
			input:    "1B:40:0A",
			expected: "1B400A",
		},
		{
			name:     "Hex with dashes",
			input:    "1B-40-0A",
			expected: "1B400A",
		},
		{
			name:     "Mixed separators",
			input:    "1B 40,0A:FF-00",
			expected: "1B400AFF00",
		},
		{
			name:     "Lowercase to uppercase",
			input:    "1b 4a ab cd ef",
			expected: "1B4AABCDEF",
		},
		{
			name:     "With invalid characters",
			input:    "1B@#$40!  ~0A",
			expected: "1B400A",
		},
		{
			name:     "Multiple spaces",
			input:    "  1B    40   0A  ",
			expected: "1B400A",
		},
		{
			name:     "Only hex digits",
			input:    "1B400A",
			expected: "1B400A",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only separators",
			input:    "   , : - ",
			expected: "",
		},
		{
			name:     "Only invalid characters",
			input:    "!  @#$%^&*()",
			expected: "",
		},
		{
			name:     "Real ESC/POS commands",
			input:    "1B 70 00 19 FA",
			expected: "1B700019FA",
		},
		{
			name:     "Complex 0x format",
			input:    "0x1B,0x70,0x00",
			expected: "1B7000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := document.CleanHexString(tt.input)
			if result != tt.expected {
				t.Errorf("CleanHexString(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Tests para ParseHexString
// ============================================================================

func TestParseHexString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []byte
		shouldErr bool
	}{
		{
			name:      "Valid hex string",
			input:     "1B 40",
			expected:  []byte{0x1B, 0x40},
			shouldErr: false,
		},
		{
			name:      "Valid hex without spaces",
			input:     "1B400A",
			expected:  []byte{0x1B, 0x40, 0x0A},
			shouldErr: false,
		},
		{
			name:      "Valid hex with 0x prefix",
			input:     "0x1B 0x40",
			expected:  []byte{0x1B, 0x40},
			shouldErr: false,
		},
		{
			name:      "Valid hex with mixed separators",
			input:     "1B,40:0A",
			expected:  []byte{0x1B, 0x40, 0x0A},
			shouldErr: false,
		},
		{
			name:      "Odd number of characters should fail",
			input:     "1B4",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "Empty string should fail",
			input:     "",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "Only spaces should fail",
			input:     "   ",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "Only invalid characters should fail",
			input:     "! @#$%",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "Cash drawer command",
			input:     "1B 70 00 19 FA",
			expected:  []byte{0x1B, 0x70, 0x00, 0x19, 0xFA},
			shouldErr: false,
		},
		{
			name:      "Full reset command",
			input:     "1B 40",
			expected:  []byte{0x1B, 0x40},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := document.ParseHexString(tt.input)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("ParseHexString(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseHexString(%q) unexpected error: %v", tt.input, err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("ParseHexString(%q) length = %d, want %d",
					tt.input, len(result), len(tt.expected))
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("ParseHexString(%q)[%d] = 0x%02X, want 0x%02X",
						tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// ============================================================================
// Tests para CheckDangerousCommands
// ============================================================================

func TestCheckDangerousCommands(t *testing.T) {
	tests := []struct {
		name        string
		bytes       []byte
		shouldWarn  bool
		warnPattern string
	}{
		{
			name:        "ESC @ (Full reset)",
			bytes:       []byte{0x1B, 0x40},
			shouldWarn:  true,
			warnPattern: "ESC @",
		},
		{
			name:        "ESC = 0 (Disable printer)",
			bytes:       []byte{0x1B, 0x3D, 0x00},
			shouldWarn:  true,
			warnPattern: "ESC =",
		},
		{
			name:        "DLE ENQ (Real-time status)",
			bytes:       []byte{0x10, 0x05, 0x04},
			shouldWarn:  true,
			warnPattern: "DLE ENQ",
		},
		{
			name:        "ESC p (Cash drawer)",
			bytes:       []byte{0x1B, 0x70},
			shouldWarn:  true,
			warnPattern: "ESC p",
		},
		{
			name:        "Safe command - LF",
			bytes:       []byte{0x0A},
			shouldWarn:  false,
			warnPattern: "",
		},
		{
			name:        "Safe command - Text",
			bytes:       []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F},
			shouldWarn:  false,
			warnPattern: "",
		},
		{
			name:        "Dangerous in middle of data",
			bytes:       []byte{0x00, 0x1B, 0x40, 0x00},
			shouldWarn:  true,
			warnPattern: "ESC @",
		},
		{
			name:        "Multiple dangerous commands",
			bytes:       []byte{0x1B, 0x40, 0x1B, 0x70},
			shouldWarn:  true,
			warnPattern: "",
		},
		{
			name:        "Empty bytes",
			bytes:       []byte{},
			shouldWarn:  false,
			warnPattern: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warning := document.CheckDangerousCommands(tt.bytes)

			if tt.shouldWarn {
				if warning == "" {
					t.Errorf("Expected warning for dangerous command, got none")
				} else if tt.warnPattern != "" && !strings.Contains(warning, tt.warnPattern) {
					t.Errorf("Warning %q doesn't contain expected pattern %q",
						warning, tt.warnPattern)
				}
			} else {
				if warning != "" {
					t.Errorf("Unexpected warning for safe command: %s", warning)
				}
			}
		})
	}
}

// ============================================================================
// Tests para ContainsBidirectionalCommand
// ============================================================================

func TestContainsBidirectionalCommand(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected bool
	}{
		{
			name:     "GS I (Transmit printer ID)",
			bytes:    []byte{0x1D, 0x49},
			expected: true,
		},
		{
			name:     "GS r (Transmit status)",
			bytes:    []byte{0x1D, 0x72},
			expected: true,
		},
		{
			name:     "DLE EOT (Real-time status)",
			bytes:    []byte{0x10, 0x04},
			expected: true,
		},
		{
			name:     "GS a (Auto status)",
			bytes:    []byte{0x1D, 0x61},
			expected: true,
		},
		{
			name:     "ESC u (Peripheral status)",
			bytes:    []byte{0x1B, 0x75},
			expected: true,
		},
		{
			name:     "ESC v (Paper sensor)",
			bytes:    []byte{0x1B, 0x76},
			expected: true,
		},
		{
			name:     "Write-only command - ESC @",
			bytes:    []byte{0x1B, 0x40},
			expected: false,
		},
		{
			name:     "Write-only command - Text",
			bytes:    []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F},
			expected: false,
		},
		{
			name:     "Bidirectional in middle",
			bytes:    []byte{0x00, 0x1D, 0x49, 0x00},
			expected: true,
		},
		{
			name:     "Empty bytes",
			bytes:    []byte{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := document.ContainsBidirectionalCommand(tt.bytes)
			if result != tt.expected {
				t.Errorf("ContainsBidirectionalCommand() = %v, want %v",
					result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Tests para ContainsSequence
// ============================================================================

func TestContainsSequence(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		pattern  []byte
		expected bool
	}{
		{
			name:     "Pattern at start",
			bytes:    []byte{0x1B, 0x40, 0x00},
			pattern:  []byte{0x1B, 0x40},
			expected: true,
		},
		{
			name:     "Pattern at end",
			bytes:    []byte{0x00, 0x1B, 0x40},
			pattern:  []byte{0x1B, 0x40},
			expected: true,
		},
		{
			name:     "Pattern in middle",
			bytes:    []byte{0x00, 0x1B, 0x40, 0x00},
			pattern:  []byte{0x1B, 0x40},
			expected: true,
		},
		{
			name:     "Pattern not present",
			bytes:    []byte{0x00, 0x1B, 0x41, 0x00},
			pattern:  []byte{0x1B, 0x40},
			expected: false,
		},
		{
			name:     "Pattern longer than bytes",
			bytes:    []byte{0x1B},
			pattern:  []byte{0x1B, 0x40},
			expected: false,
		},
		{
			name:     "Empty pattern",
			bytes:    []byte{0x1B, 0x40},
			pattern:  []byte{},
			expected: false,
		},
		{
			name:     "Empty bytes",
			bytes:    []byte{},
			pattern:  []byte{0x1B, 0x40},
			expected: false,
		},
		{
			name:     "Both empty",
			bytes:    []byte{},
			pattern:  []byte{},
			expected: false,
		},
		{
			name:     "Exact match",
			bytes:    []byte{0x1B, 0x40},
			pattern:  []byte{0x1B, 0x40},
			expected: true,
		},
		{
			name:     "Single byte match",
			bytes:    []byte{0x1B, 0x40, 0x00},
			pattern:  []byte{0x1B},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := document.ContainsSequence(tt.bytes, tt.pattern)
			if result != tt.expected {
				t.Errorf("ContainsSequence(%v, %v) = %v, want %v",
					tt.bytes, tt.pattern, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Tests de Integración para HandleRaw
// ============================================================================

func TestHandleRaw_ValidCommand(t *testing.T) {
	printer := newMockPrinterForRaw(t)
	executor := document.NewExecutor(printer)

	rawCmd := document.RawCommand{
		Hex:     "0A", // LF - safe command
		Format:  "hex",
		Comment: "Test line feed",
	}

	data, err := json.Marshal(rawCmd)
	if err != nil {
		t.Fatalf("Failed to marshal command: %v", err)
	}

	err = executor.HandleRaw(printer, data)
	if err != nil {
		t.Errorf("HandleRaw() unexpected error: %v", err)
	}
}

func TestHandleRaw_InvalidHex(t *testing.T) {
	printer := newMockPrinterForRaw(t)
	executor := document.NewExecutor(printer)

	rawCmd := document.RawCommand{
		Hex:    "ZZ ZZ", // Invalid hex - Z is not valid
		Format: "hex",
	}

	data, err := json.Marshal(rawCmd)
	if err != nil {
		t.Fatalf("Failed to marshal command: %v", err)
	}

	err = executor.HandleRaw(printer, data)
	if err == nil {
		t.Error("Expected error for invalid hex, got nil")
	}
}

func TestHandleRaw_EmptyCommand(t *testing.T) {
	printer := newMockPrinterForRaw(t)
	executor := document.NewExecutor(printer)

	rawCmd := document.RawCommand{
		Hex:    "",
		Format: "hex",
	}

	data, err := json.Marshal(rawCmd)
	if err != nil {
		t.Fatalf("Failed to marshal command: %v", err)
	}

	err = executor.HandleRaw(printer, data)
	if err == nil {
		t.Error("Expected error for empty command, got nil")
	}
}

func TestHandleRaw_SafeMode_BlocksDangerous(t *testing.T) {
	printer := newMockPrinterForRaw(t)
	executor := document.NewExecutor(printer)

	rawCmd := document.RawCommand{
		Hex:      "1B 40", // ESC @ (dangerous)
		Format:   "hex",
		SafeMode: true,
	}

	data, err := json.Marshal(rawCmd)
	if err != nil {
		t.Fatalf("Failed to marshal command: %v", err)
	}

	err = executor.HandleRaw(printer, data)
	if err == nil {
		t.Error("Expected error for dangerous command in safe mode, got nil")
	}

	if err != nil && !strings.Contains(err.Error(), "unsafe command blocked") {
		t.Errorf("Error message should mention blocking: %v", err)
	}
}

func TestHandleRaw_SafeMode_AllowsSafe(t *testing.T) {
	printer := newMockPrinterForRaw(t)
	executor := document.NewExecutor(printer)

	rawCmd := document.RawCommand{
		Hex:      "0A", // LF (safe)
		Format:   "hex",
		SafeMode: true,
	}

	data, err := json.Marshal(rawCmd)
	if err != nil {
		t.Fatalf("Failed to marshal command: %v", err)
	}

	err = executor.HandleRaw(printer, data)
	if err != nil {
		t.Errorf("HandleRaw() unexpected error for safe command: %v", err)
	}
}

func TestHandleRaw_Base64Format(t *testing.T) {
	printer := newMockPrinterForRaw(t)
	executor := document.NewExecutor(printer)

	// "Hello" in base64
	rawCmd := document.RawCommand{
		Hex:    "SGVsbG8=",
		Format: "base64",
	}

	data, err := json.Marshal(rawCmd)
	if err != nil {
		t.Fatalf("Failed to marshal command: %v", err)
	}

	err = executor.HandleRaw(printer, data)
	if err != nil {
		t.Errorf("HandleRaw() unexpected error for base64: %v", err)
	}
}

func TestHandleRaw_InvalidBase64(t *testing.T) {
	printer := newMockPrinterForRaw(t)
	executor := document.NewExecutor(printer)

	rawCmd := document.RawCommand{
		Hex:    "!!! invalid!!! ",
		Format: "base64",
	}

	data, err := json.Marshal(rawCmd)
	if err != nil {
		t.Fatalf("Failed to marshal command: %v", err)
	}

	err = executor.HandleRaw(printer, data)
	if err == nil {
		t.Error("Expected error for invalid base64, got nil")
	}
}

func TestHandleRaw_UnsupportedFormat(t *testing.T) {
	printer := newMockPrinterForRaw(t)
	executor := document.NewExecutor(printer)

	rawCmd := document.RawCommand{
		Hex:    "1B 40",
		Format: "unknown",
	}

	data, err := json.Marshal(rawCmd)
	if err != nil {
		t.Fatalf("Failed to marshal command: %v", err)
	}

	err = executor.HandleRaw(printer, data)
	if err == nil {
		t.Error("Expected error for unsupported format, got nil")
	}
}

func TestHandleRaw_MaxSizeLimit(t *testing.T) {
	printer := newMockPrinterForRaw(t)
	executor := document.NewExecutor(printer)

	// Generate hex string > 4096 bytes
	// Cada "FF " = 1 byte final, necesitamos > 4096 bytes
	largeHex := strings.Repeat("FF ", 4097)

	rawCmd := document.RawCommand{
		Hex:    largeHex,
		Format: "hex",
	}

	data, err := json.Marshal(rawCmd)
	if err != nil {
		t.Fatalf("Failed to marshal command: %v", err)
	}

	err = executor.HandleRaw(printer, data)
	if err == nil {
		t.Error("Expected error for oversized command, got nil")
	}

	if err != nil && !strings.Contains(err.Error(), "too large") {
		t.Errorf("Error should mention size: %v", err)
	}
}

func TestHandleRaw_ExactlyAtLimit(t *testing.T) {
	printer := newMockPrinterForRaw(t)
	executor := document.NewExecutor(printer)

	// Exactly 4096 bytes - should pass
	largeHex := strings.Repeat("FF ", 4096)

	rawCmd := document.RawCommand{
		Hex:    largeHex,
		Format: "hex",
	}

	data, err := json.Marshal(rawCmd)
	if err != nil {
		t.Fatalf("Failed to marshal command: %v", err)
	}

	err = executor.HandleRaw(printer, data)
	if err != nil {
		t.Errorf("HandleRaw() should accept exactly 4096 bytes: %v", err)
	}
}

func TestHandleRaw_DefaultFormatIsHex(t *testing.T) {
	printer := newMockPrinterForRaw(t)
	executor := document.NewExecutor(printer)

	// No format specified - should default to hex
	rawCmd := document.RawCommand{
		Hex: "0A",
		// Format not set
	}

	data, err := json.Marshal(rawCmd)
	if err != nil {
		t.Fatalf("Failed to marshal command: %v", err)
	}

	err = executor.HandleRaw(printer, data)
	if err != nil {
		t.Errorf("HandleRaw() should default to hex format: %v", err)
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkCleanHexString(b *testing.B) {
	input := "1B 70 00 19 FA"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = document.CleanHexString(input)
	}
}

func BenchmarkCleanHexString_Complex(b *testing.B) {
	input := "0x1B 0x70 0x00 0x19 0xFA"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = document.CleanHexString(input)
	}
}

func BenchmarkParseHexString(b *testing.B) {
	input := "1B 70 00 19 FA"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = document.ParseHexString(input)
	}
}

func BenchmarkCheckDangerousCommands(b *testing.B) {
	bytes := []byte{0x1B, 0x40, 0x00, 0x1B, 0x70}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = document.CheckDangerousCommands(bytes)
	}
}

func BenchmarkContainsSequence(b *testing.B) {
	bytes := []byte{0x00, 0x01, 0x02, 0x1B, 0x40, 0x03, 0x04}
	pattern := []byte{0x1B, 0x40}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = document.ContainsSequence(bytes, pattern)
	}
}
