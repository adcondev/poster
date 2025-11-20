package profile_test

import (
	"testing"

	"github.com/adcondev/pos-printer/pkg/commands/character"
	"github.com/adcondev/pos-printer/pkg/profile"
)

func TestEncodeString(t *testing.T) {
	tests := []struct {
		name          string
		codeTable     character.CodeTable
		input         string
		expectedError bool
		expectedBytes []byte
	}{
		{
			name:          "PC850 - Basic ASCII",
			codeTable:     character.PC850,
			input:         "Hello",
			expectedError: false,
			expectedBytes: []byte("Hello"),
		},
		{
			name:          "PC850 - Accented character (é)",
			codeTable:     character.PC850,
			input:         "café",
			expectedError: false,
			// é in PC850 is 0x82 (130)
			expectedBytes: []byte{'c', 'a', 'f', 0x82},
		},
		{
			name:          "Windows1252 - Euro sign (€)",
			codeTable:     character.WPC1252,
			input:         "€",
			expectedError: false,
			// € in Windows1252 is 0x80 (128)
			expectedBytes: []byte{0x80},
		},
		{
			name:          "Unsupported Character in Encoding",
			codeTable:     character.PC437, // PC437 doesn't have €
			input:         "€",
			expectedError: true, // Should return "encoding: rune not supported by encoding"
		},
		{
			name:      "Fallback for unsupported CodeTable",
			codeTable: character.CodeTable(99), // Unsupported
			input:     "test",
			// It should fallback to Windows1252
			expectedError: false,
			expectedBytes: []byte("test"),
		},
	}

	p := profile.CreateProfile58mm()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.CodeTable = tt.codeTable
			encoded, err := p.EncodeString(tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("EncodeString() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if !tt.expectedError && tt.expectedBytes != nil {
				if encoded != string(tt.expectedBytes) {
					t.Errorf("EncodeString() = %v, want %v", []byte(encoded), tt.expectedBytes)
				}
			}
		})
	}
}

func TestIsSupported(t *testing.T) {
	p := profile.CreateProfile58mm()

	supported := []character.CodeTable{
		character.PC437,
		character.PC850,
		character.WPC1252,
		character.Katakana,
	}

	unsupported := []character.CodeTable{
		character.CodeTable(99),
		character.CodeTable(255),
	}

	for _, ct := range supported {
		if !p.IsSupported(ct) {
			t.Errorf("expected code table %v to be supported", ct)
		}
	}

	for _, ct := range unsupported {
		if p.IsSupported(ct) {
			t.Errorf("expected code table %v to be unsupported", ct)
		}
	}
}

func TestFallbackBehavior(t *testing.T) {
	// This test specifically targets the fallback log/behavior in getEncoding
	// Since getEncoding is private, we test it via EncodeString
	p := profile.CreateProfile58mm()
	p.CodeTable = character.CodeTable(123) // Random unsupported table

	// Should fallback to Windows1252
	// Windows1252 maps '€' to 0x80

	input := "€"
	encoded, err := p.EncodeString(input)
	if err != nil {
		t.Fatalf("unexpected error during fallback encoding: %v", err)
	}

	// Check if it used Windows1252
	expected := string([]byte{0x80})
	if encoded != expected {
		t.Errorf("expected fallback to Windows1252 (result %x), got %x", expected, encoded)
	}
}
