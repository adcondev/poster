package profile_test

import (
	"bytes"
	"testing"

	"github.com/adcondev/poster/pkg/commands/character"
	"github.com/adcondev/poster/pkg/profile"
)

func TestNormalizeTypography(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
		{
			name:     "plain ascii unchanged",
			input:    "Hello, world! 123",
			expected: "Hello, world! 123",
		},
		{
			name:     "smart single quotes",
			input:    "\u2018single\u2019 and \u201Aquote\u201B",
			expected: "'single' and 'quote'",
		},
		{
			name:     "smart double quotes",
			input:    "\u201Cdouble\u201D and \u201Equote\u201F",
			expected: "\"double\" and \"quote\"",
		},
		{
			name:     "dashes and minus",
			input:    "en\u2013dash, em\u2014dash, bar\u2015bar, minus\u2212sign",
			expected: "en-dash, em-dash, bar-bar, minus-sign",
		},
		{
			name:     "ellipsis",
			input:    "Wait for it\u2026",
			expected: "Wait for it...",
		},
		{
			name:     "non-breaking spaces",
			input:    "NBSP\u00A0space and narrow\u202FNBSP",
			expected: "NBSP space and narrow NBSP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := profile.NormalizeTypography(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizeTypography(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseCodeTable(t *testing.T) {
	tests := []struct {
		input    string
		expected character.CodeTable
		found    bool
	}{
		// PC437
		{"PC437", character.PC437, true},
		{"CP437", character.PC437, true},
		{"437", character.PC437, true},
		{"pc-437", character.PC437, true},
		{"cp_437", character.PC437, true},

		// PC850
		{"PC850", character.PC850, true},
		{"CP850", character.PC850, true},
		{"850", character.PC850, true},

		// PC852
		{"PC852", character.PC852, true},
		{"CP852", character.PC852, true},
		{"852", character.PC852, true},

		// PC855
		{"PC855", character.PC855, true},
		{"CP855", character.PC855, true},
		{"855", character.PC855, true},

		// PC858
		{"PC858", character.PC858, true},
		{"CP858", character.PC858, true},
		{"858", character.PC858, true},

		// PC860
		{"PC860", character.PC860, true},
		{"CP860", character.PC860, true},
		{"860", character.PC860, true},

		// PC862
		{"PC862", character.PC862, true},
		{"CP862", character.PC862, true},
		{"862", character.PC862, true},

		// PC863
		{"PC863", character.PC863, true},
		{"CP863", character.PC863, true},
		{"863", character.PC863, true},

		// PC865
		{"PC865", character.PC865, true},
		{"CP865", character.PC865, true},
		{"865", character.PC865, true},

		// PC866
		{"PC866", character.PC866, true},
		{"CP866", character.PC866, true},
		{"866", character.PC866, true},

		// Windows code pages
		{"WPC1250", character.WPC1250, true},
		{"CP1250", character.WPC1250, true},
		{"WINDOWS1250", character.WPC1250, true},
		{"windows-1250", character.WPC1250, true},

		{"WPC1251", character.WPC1251, true},
		{"CP1251", character.WPC1251, true},
		{"WINDOWS1251", character.WPC1251, true},

		{"WPC1252", character.WPC1252, true},
		{"CP1252", character.WPC1252, true},
		{"WINDOWS1252", character.WPC1252, true},

		{"WPC1253", character.WPC1253, true},
		{"CP1253", character.WPC1253, true},
		{"WINDOWS1253", character.WPC1253, true},

		{"WPC1254", character.WPC1254, true},
		{"CP1254", character.WPC1254, true},
		{"WINDOWS1254", character.WPC1254, true},

		{"WPC1255", character.WPC1255, true},
		{"CP1255", character.WPC1255, true},
		{"WINDOWS1255", character.WPC1255, true},

		{"WPC1256", character.WPC1256, true},
		{"CP1256", character.WPC1256, true},
		{"WINDOWS1256", character.WPC1256, true},

		{"WPC1257", character.WPC1257, true},
		{"CP1257", character.WPC1257, true},
		{"WINDOWS1257", character.WPC1257, true},

		{"WPC1258", character.WPC1258, true},
		{"CP1258", character.WPC1258, true},
		{"WINDOWS1258", character.WPC1258, true},

		// ISO-8859
		{"ISO88592", character.ISO88592, true},
		{"iso-8859-2", character.ISO88592, true},
		{"ISO88597", character.ISO88597, true},
		{"iso-8859-7", character.ISO88597, true},
		{"ISO885915", character.ISO885915, true},
		{"iso-8859-15", character.ISO885915, true},

		// Asian
		{"KATAKANA", character.Katakana, true},
		{"katakana", character.Katakana, true},
		{"SHIFTJIS", character.Katakana, true},
		{"shift_jis", character.Katakana, true},
		{"HIRAGANA", character.Hiragana, true},
		{"hiragana", character.Hiragana, true},

		// Invalid
		{"UNKNOWN_CODEPAGE", 0, false},
		{"", 0, false},
		{"12345", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, ok := profile.ParseCodeTable(tt.input)
			if ok != tt.found {
				t.Fatalf("ParseCodeTable(%q) found = %v, want %v", tt.input, ok, tt.found)
			}
			if ok && got != tt.expected {
				t.Errorf("ParseCodeTable(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestEncodeBytesSafe_NonStrict(t *testing.T) {
	p := profile.CreateProfile80mm()
	p.CodeTable = character.PC437
	p.StrictEncoding = false
	p.AutoSwitchCodePage = false

	// '🚀' cannot be represented in PC437, should be replaced with '?'
	out, err := p.EncodeBytesSafe("Total: $10 🚀")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Total: $10 ?"
	if string(out) != expected {
		t.Errorf("EncodeBytesSafe() = %q, want %q", string(out), expected)
	}
}

func TestEncodeBytesSafe_Strict(t *testing.T) {
	p := profile.CreateProfile80mm()
	p.CodeTable = character.PC437
	p.StrictEncoding = true

	// PC437 does not support Euro '€'
	_, err := p.EncodeBytesSafe("Price: 10 €")
	if err == nil {
		t.Fatal("expected error in strict mode for unsupported rune '€', got nil")
	}

	// But valid PC437 characters should succeed
	out, err := p.EncodeBytesSafe("Price: $10")
	if err != nil {
		t.Fatalf("unexpected error for valid ASCII: %v", err)
	}
	if string(out) != "Price: $10" {
		t.Errorf("got %q, want 'Price: $10'", string(out))
	}
}

func TestEncodeBytesSafe_AutoSwitchCodePage(t *testing.T) {
	p := profile.CreateProfile80mm()
	p.CodeTable = character.PC437
	p.StrictEncoding = false
	p.AutoSwitchCodePage = true

	// String containing:
	// - Box drawing "┌─┐" (present in PC437)
	// - Euro sign "€" (requires PC858, code table 19)
	// - Cyrillic "Д" (requires PC866, code table 17)
	input := "┌─┐€Д"
	out, err := p.EncodeBytesSafe(input)
	if err != nil {
		t.Fatalf("EncodeBytesSafe failed: %v", err)
	}

	// Verify ESC t sequences were emitted
	// Switch to PC858: 0x1B, 0x74, 19
	escPC858 := []byte{0x1B, 0x74, byte(character.PC858)}
	if !bytes.Contains(out, escPC858) {
		t.Errorf("expected output to contain ESC t %d for PC858, got %x", character.PC858, out)
	}

	// Switch to PC866: 0x1B, 0x74, 17
	escPC866 := []byte{0x1B, 0x74, byte(character.PC866)}
	if !bytes.Contains(out, escPC866) {
		t.Errorf("expected output to contain ESC t %d for PC866, got %x", character.PC866, out)
	}

	// Verify final CodeTable on profile is PC866
	if p.CodeTable != character.PC866 {
		t.Errorf("expected profile.CodeTable to be updated to PC866 (%d), got %d", character.PC866, p.CodeTable)
	}
}

func TestEncodeBytesSafe_AutoSwitchWithFallbackList(t *testing.T) {
	p := profile.CreateProfile80mm()
	p.CodeTable = character.PC437
	p.StrictEncoding = false
	p.AutoSwitchCodePage = true
	// Restrict fallbacks to WPC1252 only
	p.FallbackCodePages = []character.CodeTable{character.WPC1252}

	// '€' is in WPC1252 (0x80)
	out, err := p.EncodeBytesSafe("€")
	if err != nil {
		t.Fatalf("EncodeBytesSafe failed: %v", err)
	}

	escWPC1252 := []byte{0x1B, 0x74, byte(character.WPC1252)}
	if !bytes.Contains(out, escWPC1252) {
		t.Errorf("expected ESC t %d for WPC1252, got %x", character.WPC1252, out)
	}
	if p.CodeTable != character.WPC1252 {
		t.Errorf("expected CodeTable = WPC1252, got %v", p.CodeTable)
	}

	// Now try Cyrillic 'Д', which is not in WPC1252
	out2, err := p.EncodeBytesSafe("Д")
	if err != nil {
		t.Fatalf("EncodeBytesSafe failed: %v", err)
	}
	// Since WPC1252 doesn't have 'Д' and FallbackCodePages only has WPC1252, it should fallback to '?'
	if string(out2) != "?" {
		t.Errorf("expected '?', got %q", string(out2))
	}
}
