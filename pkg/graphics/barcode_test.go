package graphics_test

import (
	"testing"

	"github.com/adcondev/poster/pkg/commands/barcode"
	"github.com/adcondev/poster/pkg/constants"
	"github.com/adcondev/poster/pkg/graphics"
)

func TestMapSymbology(t *testing.T) {
	tests := []struct {
		input    string
		expected barcode.Symbology
		wantErr  bool
	}{
		{"upca", barcode.UPCA, false},
		{"UPC-A", barcode.UPCA, false},
		{"upce", barcode.UPCE, false},
		{"UPC-E", barcode.UPCE, false},
		{"ean13", barcode.EAN13, false},
		{"EAN-13", barcode.EAN13, false},
		{"ean8", barcode.EAN8, false},
		{"EAN-8", barcode.EAN8, false},
		{"code39", barcode.CODE39, false},
		{"CODE-39", barcode.CODE39, false},
		{"itf", barcode.ITF, false},
		{"codabar", barcode.CODABAR, false},
		{"code128", barcode.CODE128, false},
		{"CODE-128", barcode.CODE128, false},
		{"jan13", barcode.EAN13, false},
		{"JAN-13", barcode.EAN13, false},
		{"jan8", barcode.EAN8, false},
		{"JAN-8", barcode.EAN8, false},
		{"code93", barcode.CODE93, false},
		{"code-93", barcode.CODE93, false},
		{"CODE_93", barcode.CODE93, false},
		{"gs1_128", barcode.GS1128, false},
		{"gs1-128", barcode.GS1128, false},
		{"gs1128", barcode.GS1128, false},
		{"code128_auto", barcode.CODE128Auto, false},
		{"code128-auto", barcode.CODE128Auto, false},
		{"code128auto", barcode.CODE128Auto, false},
		// Invalid / unsupported
		{"unsupported_sym", 0, true},
		{"", 0, true},
		{"qr_code", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := graphics.MapSymbology(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("MapSymbology(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("MapSymbology(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMapHriPosition(t *testing.T) {
	tests := []struct {
		input    string
		expected barcode.HRIPosition
		wantErr  bool
	}{
		{"none", barcode.HRINotPrinted, false},
		{"NONE", barcode.HRINotPrinted, false},
		{"above", barcode.HRIAbove, false},
		{"ABOVE", barcode.HRIAbove, false},
		{"below", barcode.HRIBelow, false},
		{"BELOW", barcode.HRIBelow, false},
		{"both", barcode.HRIBoth, false},
		{"BOTH", barcode.HRIBoth, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := graphics.MapHriPosition(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("MapHriPosition(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("MapHriPosition(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMapHriFont(t *testing.T) {
	tests := []struct {
		input    string
		expected barcode.HRIFont
		wantErr  bool
	}{
		{"A", barcode.HRIFontA, false},
		{"a", barcode.HRIFontA, false},
		{"B", barcode.HRIFontB, false},
		{"b", barcode.HRIFontB, false},
		{"C", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := graphics.MapHriFont(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("MapHriFont(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("MapHriFont(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIntToBarcodeModuleWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected barcode.Width
	}{
		{"below min", 0, barcode.MinWidth},
		{"negative", -5, barcode.MinWidth},
		{"min bound", constants.MinBarcodeModuleWidth, barcode.MinWidth},
		{"in range 3", 3, barcode.Width(3)},
		{"in range 4", 4, barcode.Width(4)},
		{"max bound", constants.MaxBarcodeModuleWidth, barcode.MaxWidth},
		{"above max", 10, barcode.MaxWidth},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := graphics.IntToBarcodeModuleWidth(tt.input)
			if got != tt.expected {
				t.Errorf("IntToBarcodeModuleWidth(%d) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIntToBarcodeHeight(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected barcode.Height
	}{
		{"below min", 0, barcode.MinHeight},
		{"negative", -10, barcode.MinHeight},
		{"min bound", constants.MinBarcodeHeight, barcode.MinHeight},
		{"in range 80", 80, barcode.Height(80)},
		{"in range 162", 162, barcode.Height(162)},
		{"max bound", constants.MaxBarcodeHeight, barcode.MaxHeight},
		{"above max", 300, barcode.MaxHeight},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := graphics.IntToBarcodeHeight(tt.input)
			if got != tt.expected {
				t.Errorf("IntToBarcodeHeight(%d) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestDefaultBarcodeConfig(t *testing.T) {
	cfg := graphics.DefaultBarcodeConfig()
	if cfg.Symbology != barcode.CODE128 {
		t.Errorf("expected default symbology CODE128, got %v", cfg.Symbology)
	}
	if cfg.Width != barcode.Width(constants.DefaultBarcodeModuleWidth) {
		t.Errorf("expected default width %d, got %v", constants.DefaultBarcodeModuleWidth, cfg.Width)
	}
	if cfg.Height != barcode.Height(constants.DefaultBarcodeHeight) {
		t.Errorf("expected default height %d, got %v", constants.DefaultBarcodeHeight, cfg.Height)
	}
	if cfg.HRIPosition != barcode.HRIBelow {
		t.Errorf("expected default HRIPosition HRIBelow, got %v", cfg.HRIPosition)
	}
	if cfg.HRIFont != barcode.HRIFontA {
		t.Errorf("expected default HRIFont HRIFontA, got %v", cfg.HRIFont)
	}
	if cfg.CodeSet != barcode.Code128SetB {
		t.Errorf("expected default CodeSet Code128SetB, got %v", cfg.CodeSet)
	}
}
