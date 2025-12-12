package executor

import (
	"testing"
)

// ============================================================================
// CleanHexString Tests
// ============================================================================

func TestCleanHexString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "already clean",
			input:    "1B40",
			expected: "1B40",
		},
		{
			name:     "with spaces",
			input:    "1B 40",
			expected: "1B40",
		},
		{
			name:     "with 0x prefix",
			input:    "0x1B0x40",
			expected: "1B40",
		},
		{
			name:     "lowercase to uppercase",
			input:    "1b40ab",
			expected: "1B40AB",
		},
		{
			name:     "mixed separators",
			input:    "1B, 40: 0A- FF",
			expected: "1B400AFF",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only invalid chars",
			input:    "GHXYZ",
			expected: "",
		},
		{
			name:     "0x at different positions",
			input:    "0x1B 0x40 0x0A",
			expected: "1B400A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanHexString(tt.input)
			if result != tt.expected {
				t.Errorf("CleanHexString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// ContainsSequence Tests
// ============================================================================

func TestContainsSequence(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		pattern  []byte
		expected bool
	}{
		{
			name:     "pattern at start",
			bytes:    []byte{0x1B, 0x40, 0x00, 0x00},
			pattern:  []byte{0x1B, 0x40},
			expected: true,
		},
		{
			name:     "pattern at end",
			bytes:    []byte{0x00, 0x00, 0x1B, 0x40},
			pattern:  []byte{0x1B, 0x40},
			expected: true,
		},
		{
			name:     "pattern in middle",
			bytes:    []byte{0x00, 0x1B, 0x40, 0x00},
			pattern:  []byte{0x1B, 0x40},
			expected: true,
		},
		{
			name:     "pattern not found",
			bytes:    []byte{0x00, 0x01, 0x02, 0x03},
			pattern:  []byte{0x1B, 0x40},
			expected: false,
		},
		{
			name:     "empty pattern",
			bytes:    []byte{0x1B, 0x40},
			pattern:  []byte{},
			expected: false,
		},
		{
			name:     "empty bytes",
			bytes:    []byte{},
			pattern:  []byte{0x1B, 0x40},
			expected: false,
		},
		{
			name:     "pattern longer than bytes",
			bytes:    []byte{0x1B},
			pattern:  []byte{0x1B, 0x40},
			expected: false,
		},
		{
			name:     "exact match",
			bytes:    []byte{0x1B, 0x40},
			pattern:  []byte{0x1B, 0x40},
			expected: true,
		},
		{
			name:     "partial match not complete",
			bytes:    []byte{0x1B, 0x00, 0x40},
			pattern:  []byte{0x1B, 0x40},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsSequence(tt.bytes, tt.pattern)
			if result != tt.expected {
				t.Errorf("ContainsSequence(%v, %v) = %v, want %v",
					tt.bytes, tt.pattern, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// CheckCriticalCommands Tests
// ============================================================================

func TestCheckCriticalCommands(t *testing.T) {
	tests := []struct {
		name        string
		bytes       []byte
		expectEmpty bool // true if no warning expected
	}{
		{
			name:        "ESC @ (initialize)",
			bytes:       []byte{0x1B, 0x40},
			expectEmpty: false,
		},
		{
			name:        "ESC = 0 (disable printer)",
			bytes:       []byte{0x1B, 0x3D, 0x00},
			expectEmpty: false,
		},
		{
			name:        "ESC p (cash drawer)",
			bytes:       []byte{0x1B, 0x70, 0x00, 0x32, 0x64},
			expectEmpty: false,
		},
		{
			name:        "GS :  (macro)",
			bytes:       []byte{0x1D, 0x3A},
			expectEmpty: false,
		},
		{
			name:        "safe command - text",
			bytes:       []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F}, // "Hello"
			expectEmpty: true,
		},
		{
			name:        "safe command - feed",
			bytes:       []byte{0x1B, 0x64, 0x03}, // ESC d 3
			expectEmpty: true,
		},
		{
			name:        "critical embedded in larger sequence",
			bytes:       []byte{0x00, 0x00, 0x1B, 0x40, 0x00, 0x00},
			expectEmpty: false,
		},
		{
			name:        "empty bytes",
			bytes:       []byte{},
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckCriticalCommands(tt.bytes)
			if tt.expectEmpty && result != "" {
				t.Errorf("CheckCriticalCommands(%v) = %q, expected empty", tt.bytes, result)
			}
			if !tt.expectEmpty && result == "" {
				t.Errorf("CheckCriticalCommands(%v) = empty, expected warning", tt.bytes)
			}
		})
	}
}

// ============================================================================
// ContainsBidirectionalCommand Tests
// ============================================================================

func TestContainsBidirectionalCommand(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected bool
	}{
		{
			name:     "DLE ENQ 4 - real-time status",
			bytes:    []byte{0x10, 0x05, 0x04},
			expected: true,
		},
		{
			name:     "GS I - transmit printer ID",
			bytes:    []byte{0x1D, 0x49, 0x01},
			expected: true,
		},
		{
			name:     "GS r - transmit status",
			bytes:    []byte{0x1D, 0x72, 0x01},
			expected: true,
		},
		{
			name:     "DLE EOT - real-time status",
			bytes:    []byte{0x10, 0x04, 0x01},
			expected: true,
		},
		{
			name:     "safe print command",
			bytes:    []byte{0x1B, 0x40}, // ESC @
			expected: false,
		},
		{
			name:     "empty bytes",
			bytes:    []byte{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsBidirectionalCommand(tt.bytes)
			if result != tt.expected {
				t.Errorf("ContainsBidirectionalCommand(%v) = %v, want %v",
					tt.bytes, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkParseHexString(b *testing.B) {
	input := "1B 40 1B 61 01 1D 56 42 00"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseHexString(input)
	}
}

func BenchmarkCleanHexString(b *testing.B) {
	input := "0x1B 0x40, 0x1B:  0x61- 0x01"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CleanHexString(input)
	}
}

func BenchmarkContainsSequence(b *testing.B) {
	bytes := []byte{0x00, 0x00, 0x1B, 0x40, 0x00, 0x00, 0x1B, 0x61, 0x01}
	pattern := []byte{0x1B, 0x40}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ContainsSequence(bytes, pattern)
	}
}

func BenchmarkCheckCriticalCommands(b *testing.B) {
	bytes := []byte{0x1B, 0x61, 0x01, 0x48, 0x65, 0x6C, 0x6C, 0x6F}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CheckCriticalCommands(bytes)
	}
}
