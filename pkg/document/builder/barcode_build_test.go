package builder

import (
	"encoding/json"
	"testing"

	"github.com/adcondev/poster/pkg/constants"
)

// ============================================================================
// Barcode Builder Tests
// ============================================================================

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
}

func TestBarcodeBuilderAlign(t *testing.T) {
	doc := NewDocument().
		SetProfile("Test", 80, "WPC1252").
		Barcode("EAN13", "1234567890123").
		Align(constants.Right).
		End().
		Build()

	var cmd barcodeCommand
	_ = json.Unmarshal(doc.Commands[0].Data, &cmd)

	if *cmd.Align != constants.Right.String() {
		t.Errorf("Expected align 'right', got '%s'", *cmd.Align)
	}
}
