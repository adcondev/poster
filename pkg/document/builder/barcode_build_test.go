package builder

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

func TestBarcodeBuilder(t *testing.T) {
	width := 3
	height := 100

	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Barcode("CODE128", "12345").
		Width(width).
		Height(height).
		HRIPosition("below").
		HRIFont("A").
		Center().
		End().
		Build()

	if doc.Commands[0].Type != "barcode" {
		t.Errorf("Expected type 'barcode', got '%s'", doc.Commands[0].Type)
	}

	var cmd barcodeCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Symbology != "CODE128" {
		t.Errorf("Expected symbology 'CODE128', got '%s'", cmd.Symbology)
	}

	if cmd.Data != "12345" {
		t.Errorf("Expected data '12345', got '%s'", cmd.Data)
	}

	if *cmd.Width != width {
		t.Errorf("Expected width %d, got %d", width, *cmd.Width)
	}

	if *cmd.Height != height {
		t.Errorf("Expected height %d, got %d", height, *cmd.Height)
	}

	if *cmd.HRIPosition != "below" {
		t.Errorf("Expected HRI position 'below', got '%s'", *cmd.HRIPosition)
	}

	if *cmd.HRIFont != "A" {
		t.Errorf("Expected HRI font 'A', got '%s'", *cmd.HRIFont)
	}
}

func TestBarcodeBuilderDefaults(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Barcode("EAN13", "1234567890123").
		End().
		Build()

	var cmd barcodeCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if cmd.Symbology != "EAN13" {
		t.Errorf("Expected symbology 'EAN13', got '%s'", cmd.Symbology)
	}

	if cmd.Width != nil {
		t.Errorf("Expected width to be nil, got %d", *cmd.Width)
	}

	if cmd.Height != nil {
		t.Errorf("Expected height to be nil, got %d", *cmd.Height)
	}

	if cmd.HRIPosition != nil {
		t.Errorf("Expected HRI position to be nil, got '%s'", *cmd.HRIPosition)
	}

	if cmd.Align != nil {
		t.Errorf("Expected align to be nil, got '%s'", *cmd.Align)
	}
}

func TestBarcodeBuilderAlignment(t *testing.T) {
	tests := []struct {
		name     string
		align    func(*BarcodeBuilder) *BarcodeBuilder
		expected string
	}{
		{"Left", func(bb *BarcodeBuilder) *BarcodeBuilder { return bb.Left() }, constants.Left.String()},
		{"Center", func(bb *BarcodeBuilder) *BarcodeBuilder { return bb.Center() }, constants.Center.String()},
		{"Right", func(bb *BarcodeBuilder) *BarcodeBuilder { return bb.Right() }, constants.Right.String()},
		{"Align", func(bb *BarcodeBuilder) *BarcodeBuilder { return bb.Align(constants.Right) }, constants.Right.String()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bb := newBarcodeBuilder(NewDocument(), "CODE128", "data")
			tt.align(bb)

			if *bb.align != tt.expected {
				t.Errorf("Expected align '%s', got '%s'", tt.expected, *bb.align)
			}
		})
	}
}
